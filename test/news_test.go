package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	config "github.com/furqonzt99/news-redis/configs"
	"github.com/furqonzt99/news-redis/constants"
	"github.com/furqonzt99/news-redis/delivery/common"
	"github.com/furqonzt99/news-redis/delivery/controllers/news"
	"github.com/furqonzt99/news-redis/domain/entity"
	"github.com/furqonzt99/news-redis/domain/repository"
	"github.com/furqonzt99/news-redis/seeder"
	"github.com/furqonzt99/news-redis/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	config := config.GetConfig()
	db := utils.InitDB(config)

	db.Migrator().DropTable(&entity.News{}, &entity.Tag{}, &entity.NewsTags{})

	utils.InitialMigrate(db)

	constants.Rdb = utils.InitRedis(config)

	seeder.TagSeeder(db)
	seeder.NewsSeeder(db)
	seeder.NewsTagsSeeder(db)

	m.Run()
}

func TestCreateNews(t *testing.T) {
	config := config.GetConfig()
	db := utils.InitDB(config)

	e := echo.New()

	e.Validator = &common.Validator{Validator: validator.New()}

	nr := repository.NewNewsRepository(db)

	nc := news.NewNewsController(nr)

	t.Run("Create news success", func(t *testing.T) {
		e.POST("/news", nc.Create)

		createNewsRequest, _ := json.Marshal(news.CreateNewsRequest{
			Title: "Test Title",
			Body:  "Test Body",
			Tags:  []int{1, 2},
		})

		req := httptest.NewRequest(echo.POST, "/news", bytes.NewBuffer(createNewsRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Create news bad request validator", func(t *testing.T) {
		e.POST("/news", nc.Create)

		createNewsRequest, _ := json.Marshal(news.CreateNewsRequest{
			Title: "Test Title",
			Body:  "Test Body",
		})

		req := httptest.NewRequest(echo.POST, "/news", bytes.NewBuffer(createNewsRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Create news bad request repository", func(t *testing.T) {
		e.POST("/news", nc.Create)

		createNewsRequest, _ := json.Marshal(news.CreateNewsRequest{
			Title: "Test Title",
			Body:  "Test Body",
			Tags:  []int{11, 12},
		})

		req := httptest.NewRequest(echo.POST, "/news", bytes.NewBuffer(createNewsRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestGetNews(t *testing.T) {
	config := config.GetConfig()
	db := utils.InitDB(config)

	e := echo.New()

	e.Validator = &common.Validator{Validator: validator.New()}

	nr := repository.NewNewsRepository(db)

	nc := news.NewNewsController(nr)

	t.Run("Get one news success from database", func(t *testing.T) {
		e.GET("/news/:id", nc.ReadOne)

		req := httptest.NewRequest(echo.GET, "/news/1", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "database", response.Source)
	})

	t.Run("Get one news success from cache", func(t *testing.T) {
		e.GET("/news/:id", nc.ReadOne)

		req := httptest.NewRequest(echo.GET, "/news/1", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "cache", response.Source)
	})

	t.Run("Get one news bad request", func(t *testing.T) {
		e.GET("/news/:id", nc.ReadOne)

		req := httptest.NewRequest(echo.GET, "/news/ejejd", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Get one news not found", func(t *testing.T) {
		e.GET("/news/:id", nc.ReadOne)

		req := httptest.NewRequest(echo.GET, "/news/9999", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("Get all news success from database", func(t *testing.T) {
		e.GET("/news", nc.ReadAll)

		req := httptest.NewRequest(echo.GET, "/news", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "database", response.Source)
	})

	t.Run("Get all news success from cache", func(t *testing.T) {
		e.GET("/news", nc.ReadAll)

		req := httptest.NewRequest(echo.GET, "/news", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "cache", response.Source)
	})

	t.Run("Get all news success (topic & status) from database", func(t *testing.T) {
		e.GET("/news", nc.ReadAll)

		req := httptest.NewRequest(echo.GET, "/news?topic=Topic1&status=draft", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "database", response.Source)
	})

	t.Run("Get all news success (topic & status) from cache", func(t *testing.T) {
		e.GET("/news", nc.ReadAll)

		req := httptest.NewRequest(echo.GET, "/news?topic=Topic1&status=draft", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "cache", response.Source)
	})

	t.Run("Get all news success (topic) from database", func(t *testing.T) {
		e.GET("/news", nc.ReadAll)

		req := httptest.NewRequest(echo.GET, "/news?topic=Topic1", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "database", response.Source)
	})

	t.Run("Get all news success (topic) from cache", func(t *testing.T) {
		e.GET("/news", nc.ReadAll)

		req := httptest.NewRequest(echo.GET, "/news?topic=Topic1", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "cache", response.Source)
	})

	t.Run("Get all news success (status) from database", func(t *testing.T) {
		e.GET("/news", nc.ReadAll)

		req := httptest.NewRequest(echo.GET, "/news?status=draft", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "database", response.Source)
	})

	t.Run("Get all news success (status) from cache", func(t *testing.T) {
		e.GET("/news", nc.ReadAll)

		req := httptest.NewRequest(echo.GET, "/news?status=draft", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "cache", response.Source)
	})
}

func TestEditNews(t *testing.T) {
	config := config.GetConfig()
	db := utils.InitDB(config)

	e := echo.New()

	e.Validator = &common.Validator{Validator: validator.New()}

	nr := repository.NewNewsRepository(db)

	nc := news.NewNewsController(nr)

	t.Run("Update news success", func(t *testing.T) {
		e.PUT("/news/:id", nc.Edit)

		updateNewsRequest, _ := json.Marshal(news.UpdateNewsRequest{
			Title: "Test Title New",
			Body:  "Test Body New",
			Tags:  []int{1, 2, 3},
		})

		req := httptest.NewRequest(echo.PUT, "/news/1", bytes.NewBuffer(updateNewsRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Update news bad request validator", func(t *testing.T) {
		e.PUT("/news/:id", nc.Edit)

		updateNewsRequest, _ := json.Marshal(news.UpdateNewsRequest{
			Title: "Test Title",
			Body:  "Test Body",
		})

		req := httptest.NewRequest(echo.PUT, "/news/1", bytes.NewBuffer(updateNewsRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Update news bad request param", func(t *testing.T) {
		e.PUT("/news/:id", nc.Edit)

		updateNewsRequest, _ := json.Marshal(news.UpdateNewsRequest{
			Title: "Test Title",
			Body:  "Test Body",
		})

		req := httptest.NewRequest(echo.PUT, "/news/abc", bytes.NewBuffer(updateNewsRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Update news not found", func(t *testing.T) {
		e.PUT("/news/:id", nc.Edit)

		updateNewsRequest, _ := json.Marshal(news.UpdateNewsRequest{
			Title: "Test Title",
			Body:  "Test Body",
			Tags:  []int{1, 2},
		})

		req := httptest.NewRequest(echo.PUT, "/news/9999", bytes.NewBuffer(updateNewsRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestSetPublishNews(t *testing.T) {
	config := config.GetConfig()
	db := utils.InitDB(config)

	e := echo.New()

	e.Validator = &common.Validator{Validator: validator.New()}

	nr := repository.NewNewsRepository(db)

	nc := news.NewNewsController(nr)

	t.Run("Set Publish news success", func(t *testing.T) {
		e.PUT("/news/:id/publish", nc.SetStatusPublish)

		req := httptest.NewRequest(echo.PUT, "/news/1/publish", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Set Publish news bad request", func(t *testing.T) {
		e.PUT("/news/:id/publish", nc.SetStatusPublish)

		req := httptest.NewRequest(echo.PUT, "/news/qwer/publish", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Set Publish news not found", func(t *testing.T) {
		e.PUT("/news/:id/publish", nc.SetStatusPublish)

		req := httptest.NewRequest(echo.PUT, "/news/9999/publish", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestSetDraftNews(t *testing.T) {
	config := config.GetConfig()
	db := utils.InitDB(config)

	e := echo.New()

	e.Validator = &common.Validator{Validator: validator.New()}

	nr := repository.NewNewsRepository(db)

	nc := news.NewNewsController(nr)

	t.Run("Set Draft news success", func(t *testing.T) {
		e.PUT("/news/:id/draft", nc.SetStatusDraft)

		req := httptest.NewRequest(echo.PUT, "/news/1/draft", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Set Draft news bad request", func(t *testing.T) {
		e.PUT("/news/:id/draft", nc.SetStatusDraft)

		req := httptest.NewRequest(echo.PUT, "/news/qwer/draft", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Set Draft news not found", func(t *testing.T) {
		e.PUT("/news/:id/draft", nc.SetStatusDraft)

		req := httptest.NewRequest(echo.PUT, "/news/9999/draft", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}
func TestSetDeletedNews(t *testing.T) {
	config := config.GetConfig()
	db := utils.InitDB(config)

	e := echo.New()

	e.Validator = &common.Validator{Validator: validator.New()}

	nr := repository.NewNewsRepository(db)

	nc := news.NewNewsController(nr)

	t.Run("Set Deleted news success", func(t *testing.T) {
		e.PUT("/news/:id/deleted", nc.SetStatusDeleted)

		req := httptest.NewRequest(echo.PUT, "/news/1/deleted", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Set Deleted news bad request", func(t *testing.T) {
		e.PUT("/news/:id/deleted", nc.SetStatusDeleted)

		req := httptest.NewRequest(echo.PUT, "/news/qwer/deleted", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Set Deleted news not found", func(t *testing.T) {
		e.PUT("/news/:id/deleted", nc.SetStatusDeleted)

		req := httptest.NewRequest(echo.PUT, "/news/9999/deleted", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestDeleteNews(t *testing.T) {
	config := config.GetConfig()
	db := utils.InitDB(config)

	e := echo.New()

	e.Validator = &common.Validator{Validator: validator.New()}

	nr := repository.NewNewsRepository(db)

	nc := news.NewNewsController(nr)

	t.Run("Delete news success", func(t *testing.T) {
		e.DELETE("/news/:id", nc.Delete)

		req := httptest.NewRequest(echo.DELETE, "/news/1", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Delete news bad request", func(t *testing.T) {
		e.DELETE("/news/:id", nc.Delete)

		req := httptest.NewRequest(echo.DELETE, "/news/qwer", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Delete news not found", func(t *testing.T) {
		e.DELETE("/news/:id", nc.Delete)

		req := httptest.NewRequest(echo.DELETE, "/news/9999", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}
