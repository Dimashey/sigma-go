package items

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Transport struct {
	repository *ItemsRepositry
}

var (
	ErrInvalidJSON = errors.New("INVALID JSON")
	ErrInvalidID   = errors.New("INVALID ID")
)

func sendJSONResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func NewTransport() *Transport {
	repository := NewRepository()
  return &Transport{repository}
}

func (t Transport) Create(w http.ResponseWriter, r *http.Request) {
	var item Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		sendJSONResponse(w, map[string]string{"error": ErrInvalidJSON.Error()}, http.StatusBadRequest)
		return
	}

	item, err := t.repository.Create(item)
	if err != nil {
		sendJSONResponse(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, item, http.StatusCreated)
}

func (t Transport) GetMany(w http.ResponseWriter, r *http.Request) {
	items := t.repository.GetMany()

	sendJSONResponse(w, items, http.StatusOK)
}

func (t Transport) GetOne(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/v1/items/"):])
	if err != nil {
		sendJSONResponse(w, map[string]string{"error": ErrInvalidID.Error()}, http.StatusBadRequest)
		return
	}

	item, err := t.repository.GetOne(id)
	if err != nil {
		sendJSONResponse(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	sendJSONResponse(w, item, http.StatusOK)
}

func (t Transport) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/v1/items/"):])
	if err != nil {
		sendJSONResponse(w, map[string]string{"error": ErrInvalidID.Error()}, http.StatusBadRequest)
		return
	}

	var updatedItem Item

	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		sendJSONResponse(w, map[string]string{"error": ErrInvalidJSON.Error()}, http.StatusBadRequest)
		return
	}

	item, err := t.repository.Update(id, updatedItem)
	if err != nil {
		sendJSONResponse(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	sendJSONResponse(w, item, http.StatusOK)
}

func (t Transport) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/v1/items/"):])
	if err != nil {
		sendJSONResponse(w, map[string]string{"error": ErrInvalidID.Error()}, http.StatusBadRequest)
		return
	}

	err = t.repository.Delete(id)
	if err != nil {
		sendJSONResponse(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	sendJSONResponse(w, map[string]string{"message": "Item deleted"}, http.StatusOK)
}
