package tags

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/furqonzt99/news-redis/delivery/common"
	"github.com/furqonzt99/news-redis/domain/entity"
	"github.com/furqonzt99/news-redis/domain/repository"
	"github.com/furqonzt99/news-redis/services"
	"github.com/labstack/echo/v4"
)

var tagEntity string = "tag"

type TagController struct {
	Repository repository.TagInterface
}

func NewTagController(tagRepository repository.TagInterface) *TagController {
	return &TagController{Repository: tagRepository}
}

func (tc TagController) Create(c echo.Context) error {
	var tagRequest tagRequest

	c.Bind(&tagRequest)

	if err := c.Validate(&tagRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	tag := entity.Tag{
		Name: tagRequest.Name,
	}

	_, err := tc.Repository.Create(tag)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	go services.DeleteCache(tagEntity)

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (tc TagController) ReadAll(c echo.Context) error {

	response := []tagResponse{}

	// get data from cache
	newsCache, err := services.GetCache(tagEntity, 0, "")
	if err == nil {
		// Unmarshal response
		_ = json.Unmarshal([]byte(newsCache), &response)
		return c.JSON(http.StatusOK, common.SuccessResponseWithData(response))
	}

	tagsDB, err := tc.Repository.ReadAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	for _, tag := range tagsDB {
		response = append(response, tagResponse{
			ID:   int(tag.ID),
			Name: tag.Name,
		})
	}

	// Marshal response
	resMarshal, _ := json.Marshal(response)

	// Create cache
	go services.CreateCache(tagEntity, 0, "", resMarshal)

	return c.JSON(http.StatusOK, common.SuccessResponseWithData(response))
}

func (tc TagController) Edit(c echo.Context) error {
	tagID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	var tagRequest tagRequest

	c.Bind(&tagRequest)

	if err := c.Validate(&tagRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	tag := entity.Tag{
		Name: tagRequest.Name,
	}

	_, err = tc.Repository.Edit(tagID, tag)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	go services.DeleteCache(tagEntity)

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (tc TagController) Delete(c echo.Context) error {
	tagID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	_, err = tc.Repository.Delete(tagID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	go services.DeleteCache(tagEntity)

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}
