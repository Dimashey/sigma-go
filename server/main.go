package main

import (
	"github.com/Dimashey/sigma-go/server/items"
	"github.com/Dimashey/sigma-go/server/middlewares"
	"github.com/gin-gonic/gin"
)



func main() {
	itemsTransport := items.NewTransport()

  r := gin.Default()

  r.Use(middlewares.CorsConfig())

  v1 := r.Group("v1", middlewares.ApiKeyAuth(), middlewares.RateLimitMiddleware())

  v1.POST("/items", itemsTransport.Create)
  v1.GET("/items", itemsTransport.GetMany)
  v1.GET("/items/:id", itemsTransport.GetOne)
  v1.PUT("/items/:id", itemsTransport.Update)
  v1.DELETE("/items/:id", itemsTransport.Delete)

  r.Run(":8080")
}
