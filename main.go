package main

import (
	config "github.com/furqonzt99/news-redis/configs"
	"github.com/furqonzt99/news-redis/constants"
	"github.com/furqonzt99/news-redis/delivery/common"
	"github.com/furqonzt99/news-redis/delivery/controllers/news"
	"github.com/furqonzt99/news-redis/delivery/controllers/tags"
	"github.com/furqonzt99/news-redis/delivery/middlewares"
	"github.com/furqonzt99/news-redis/delivery/routes"
	"github.com/furqonzt99/news-redis/domain/repository"
	"github.com/furqonzt99/news-redis/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := config.GetConfig()

	db := utils.InitDB(config)

	utils.InitialMigrate(db)

	constants.Rdb = utils.InitRedis(config)

	e := echo.New()

	// CORS
	e.Use(middleware.CORS())

	// logger
	middlewares.LogMiddleware(e)

	// remove trailing slash
	e.Pre(middleware.RemoveTrailingSlash())

	// validator
	e.Validator = &common.Validator{Validator: validator.New()}

	// repository
	tr := repository.NewTagRepository(db)
	nr := repository.NewNewsRepository(db)

	// controller
	tc := tags.NewTagController(tr)
	nc := news.NewNewsController(nr)

	// routes
	routes.RegisterTagPath(e, tc)
	routes.RegisterNewsPath(e, nc)

	e.Logger.Fatal(e.Start(":" + config.Port))
}
