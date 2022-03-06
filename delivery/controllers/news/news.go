package news

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/furqonzt99/news-redis/delivery/common"
	"github.com/furqonzt99/news-redis/domain/entity"
	"github.com/furqonzt99/news-redis/domain/repository"
	"github.com/furqonzt99/news-redis/services"
	"github.com/labstack/echo/v4"
)

var newsEntity string = "news"

type NewsController struct {
	Repository repository.NewsInterface
}

func NewNewsController(repository repository.NewsInterface) *NewsController {
	return &NewsController{Repository: repository}
}

func (nc NewsController) Create(c echo.Context) error {
	var newsRequest createNewsRequest

	c.Bind(&newsRequest)

	if err := c.Validate(&newsRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	news := entity.News{
		Title: newsRequest.Title,
		Body:  newsRequest.Body,
	}

	_, err := nc.Repository.Create(news, newsRequest.Tags)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	go services.DeleteCache(newsEntity)

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (nc NewsController) ReadAll(c echo.Context) error {

	newsFilter := entity.NewsFilter{
		Status: c.QueryParam("status"),
		Tags:   strings.Split(c.QueryParam("topic"), ","),
	}

	response := []newsResponse{}

	// get data from cache
	newsCache, err := services.GetCache(newsEntity, 0, newsFilter)
	if err == nil {
		// Unmarshal response
		_ = json.Unmarshal([]byte(newsCache), &response)
		return c.JSON(http.StatusOK, common.SuccessResponseWithData(response))
	}

	newsDB, err := nc.Repository.ReadAll(newsFilter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	for _, news := range newsDB {
		if len(news.Tags) > 0 {

			tags := []string{}

			for _, tag := range news.Tags {
				tags = append(tags, tag.Name)
			}

			response = append(response, newsResponse{
				ID:     int(news.ID),
				Title:  news.Title,
				Body:   news.Body,
				Status: news.Status,
				Tags:   tags,
			})
		}
	}

	// Marshal response
	resMarshal, _ := json.Marshal(response)

	// Create cache
	go services.CreateCache(newsEntity, 0, newsFilter, resMarshal)

	return c.JSON(http.StatusOK, common.SuccessResponseWithData(response))
}

func (nc NewsController) ReadOne(c echo.Context) error {

	newsID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	response := newsResponse{}

	// get data from cache
	newsCache, err := services.GetCache(newsEntity, newsID, "")
	if err == nil {
		// Unmarshal response
		_ = json.Unmarshal([]byte(newsCache), &response)
		return c.JSON(http.StatusOK, common.SuccessResponseWithData(response))
	}

	newsDB, err := nc.Repository.ReadOne(newsID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	tags := []string{}

	for _, tag := range newsDB.Tags {
		tags = append(tags, tag.Name)
	}

	response = newsResponse{
		ID:     newsID,
		Title:  newsDB.Title,
		Body:   newsDB.Body,
		Status: newsDB.Status,
		Tags:   tags,
	}

	// Marshal response
	resMarshal, _ := json.Marshal(response)

	// Create cache
	go services.CreateCache(newsEntity, newsID, "", resMarshal)

	return c.JSON(http.StatusOK, common.SuccessResponseWithData(response))
}

func (nc NewsController) Edit(c echo.Context) error {
	newsID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	var newsRequest updateNewsRequest

	c.Bind(&newsRequest)

	if err := c.Validate(&newsRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	news := entity.News{
		Title: newsRequest.Title,
		Body:  newsRequest.Body,
	}

	_, err = nc.Repository.Edit(newsID, news, newsRequest.Tags)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	go services.DeleteCache(newsEntity)

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (nc NewsController) Delete(c echo.Context) error {
	newsID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	_, err = nc.Repository.Delete(newsID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	go services.DeleteCache(newsEntity)

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (nc NewsController) SetStatusDeleted(c echo.Context) error {
	newsID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	_, err = nc.Repository.SetStatusDeleted(newsID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	go services.DeleteCache(newsEntity)

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (nc NewsController) SetStatusPublish(c echo.Context) error {
	newsID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	_, err = nc.Repository.SetStatusPublish(newsID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	go services.DeleteCache(newsEntity)

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (nc NewsController) SetStatusDraft(c echo.Context) error {
	newsID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	_, err = nc.Repository.SetStatusDraft(newsID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	go services.DeleteCache(newsEntity)

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}
