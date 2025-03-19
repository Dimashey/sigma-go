package items

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Transport struct {
	repository *ItemsRepositry
}

var (
	ErrInvalidJSON = errors.New("INVALID JSON")
	ErrInvalidID   = errors.New("INVALID ID")
)

func NewTransport() *Transport {
	repository := NewRepository()
  return &Transport{repository}
}

func (t Transport) Create(c *gin.Context) {
	var item Item

	if err := c.ShouldBindJSON(&item); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON.Error()})
		return
	}

	item, err := t.repository.Create(item)

	if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

  c.JSON(http.StatusCreated, item)
}

func (t Transport) GetMany(c *gin.Context) {
	items := t.repository.GetMany()

  c.JSON(http.StatusOK, items)
}

func (t Transport) GetOne(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidID.Error()})
		return
	}

	item, err := t.repository.GetOne(id)

	if err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

  c.JSON(http.StatusOK, item)
}

func (t Transport) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidID.Error()})
		return
	}

	var updatedItem Item

	if err := c.ShouldBindJSON(&updatedItem); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON.Error()})
		return
	}

	item, err := t.repository.Update(id, updatedItem)
	if err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

  c.JSON(http.StatusOK, item)
}

func (t Transport) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidID.Error()})
		return
	}

	err = t.repository.Delete(id)
	if err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

  c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}
