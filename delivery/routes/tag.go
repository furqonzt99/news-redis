package routes

import (
	"github.com/furqonzt99/news-redis/delivery/controllers/tags"
	"github.com/labstack/echo/v4"
)

func RegisterTagPath(e *echo.Echo, tagController *tags.TagController) {
	e.POST("/tags", tagController.Create)
	e.GET("/tags", tagController.ReadAll)
	e.PUT("/tags/:id", tagController.Edit)
	e.DELETE("/tags/:id", tagController.Delete)
}
