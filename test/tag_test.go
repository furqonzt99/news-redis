package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	config "github.com/furqonzt99/news-redis/configs"
	"github.com/furqonzt99/news-redis/delivery/common"
	"github.com/furqonzt99/news-redis/delivery/controllers/tags"
	"github.com/furqonzt99/news-redis/domain/repository"
	"github.com/furqonzt99/news-redis/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateTag(t *testing.T) {
	config := config.GetConfig()
	db := utils.InitDB(config)

	e := echo.New()

	e.Validator = &common.Validator{Validator: validator.New()}

	tr := repository.NewTagRepository(db)

	tc := tags.NewTagController(tr)

	t.Run("Create tag success", func(t *testing.T) {
		e.POST("/tags", tc.Create)

		createTagRequest, _ := json.Marshal(tags.TagRequest{
			Name: "Tags Test",
		})

		req := httptest.NewRequest(echo.POST, "/tags", bytes.NewBuffer(createTagRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Create tag bad request", func(t *testing.T) {
		e.POST("/tags", tc.Create)

		createTagRequest, _ := json.Marshal(tags.TagRequest{})

		req := httptest.NewRequest(echo.POST, "/tags", bytes.NewBuffer(createTagRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestGetTag(t *testing.T) {
	config := config.GetConfig()
	db := utils.InitDB(config)

	e := echo.New()

	e.Validator = &common.Validator{Validator: validator.New()}

	tr := repository.NewTagRepository(db)

	tc := tags.NewTagController(tr)

	t.Run("Get All tag success", func(t *testing.T) {
		e.GET("/tags", tc.ReadAll)

		req := httptest.NewRequest(echo.GET, "/tags", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func TestEditTag(t *testing.T) {
	config := config.GetConfig()
	db := utils.InitDB(config)

	e := echo.New()

	e.Validator = &common.Validator{Validator: validator.New()}

	tr := repository.NewTagRepository(db)

	tc := tags.NewTagController(tr)

	t.Run("Edit tag success", func(t *testing.T) {
		e.PUT("/tags/:id", tc.Edit)

		updateTagRequest, _ := json.Marshal(tags.TagRequest{
			Name: "Tags Test",
		})

		req := httptest.NewRequest(echo.PUT, "/tags/1", bytes.NewBuffer(updateTagRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Edit tag bad request validator", func(t *testing.T) {
		e.PUT("/tags/:id", tc.Edit)

		updateTagRequest, _ := json.Marshal(tags.TagRequest{})

		req := httptest.NewRequest(echo.PUT, "/tags/1", bytes.NewBuffer(updateTagRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Edit tag bad request params", func(t *testing.T) {
		e.PUT("/tags/:id", tc.Edit)

		updateTagRequest, _ := json.Marshal(tags.TagRequest{})

		req := httptest.NewRequest(echo.PUT, "/tags/qwer", bytes.NewBuffer(updateTagRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Edit tag not found", func(t *testing.T) {
		e.PUT("/tags/:id", tc.Edit)

		updateTagRequest, _ := json.Marshal(tags.TagRequest{
			Name: "Test Topic",
		})

		req := httptest.NewRequest(echo.PUT, "/tags/9999", bytes.NewBuffer(updateTagRequest))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestDeleteTag(t *testing.T) {
	config := config.GetConfig()
	db := utils.InitDB(config)

	e := echo.New()

	e.Validator = &common.Validator{Validator: validator.New()}

	tr := repository.NewTagRepository(db)

	tc := tags.NewTagController(tr)

	t.Run("Delete tag success", func(t *testing.T) {
		e.DELETE("/tags/:id", tc.Delete)

		req := httptest.NewRequest(echo.DELETE, "/tags/1", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Delete tag bad request params", func(t *testing.T) {
		e.DELETE("/tags/:id", tc.Delete)

		req := httptest.NewRequest(echo.DELETE, "/tags/qwer", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Delete tag not found", func(t *testing.T) {
		e.DELETE("/tags/:id", tc.Delete)

		req := httptest.NewRequest(echo.DELETE, "/tags/9999", nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		var response common.DefaultResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}
