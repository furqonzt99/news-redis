package routes

import (
	"github.com/furqonzt99/news-redis/delivery/controllers/news"
	"github.com/labstack/echo/v4"
)

func RegisterNewsPath(e *echo.Echo, newsController *news.NewsController) {
	e.POST("/news", newsController.Create)
	e.GET("/news", newsController.ReadAll)
	e.GET("/news/:id", newsController.ReadOne)
	e.PUT("/news/:id", newsController.Edit)
	e.PUT("/news/:id/publish", newsController.SetStatusPublish)
	e.PUT("/news/:id/draft", newsController.SetStatusDraft)
	e.PUT("/news/:id/deleted", newsController.SetStatusDeleted)
	e.DELETE("/news/:id", newsController.Delete)
}
