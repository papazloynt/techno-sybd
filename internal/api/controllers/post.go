package controllers

import (
	"SYBD/internal/model/dto"
	"SYBD/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"strconv"
)

type PostController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *PostController) CreatePost(ctx echo.Context) error {
	var request []*dto.Post
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(ctx.Request().Body)
	if err != nil {
		//log.Errorf("Unmarshal error: %s", err)
		return err
	}
	err = json.Unmarshal(buf.Bytes(), &request)

	slugOrID := ctx.Param("slug_or_id")
	response, err := c.registry.PostService.CreatePost(context.Background(), slugOrID, request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Value)
}

func (c *PostController) GetPost(ctx echo.Context) error {
	slugOrID := ctx.Param("slug_or_id")

	sort := ctx.QueryParam("sort")
	if sort == "" {
		sort = "flat"
	}

	since := ctx.QueryParam("since")
	if since == "" {
		since = "-1"
	}

	limit := ctx.QueryParam("limit")
	if limit == "" {
		limit = "100"
	}

	sinceInt, _ := strconv.ParseInt(since, 10, 64)
	limitInt, _ := strconv.ParseInt(limit, 10, 64)
	descBool, _ := strconv.ParseBool(ctx.QueryParam("desc"))

	response, err := c.registry.PostService.GetPost(context.Background(),
		slugOrID,
		sort,
		sinceInt,
		descBool,
		limitInt)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Value)
}

func (c *PostController) GetPostDetails(ctx echo.Context) error {
	request := &dto.GetPostDetailsRequest{}
	if err := ctx.Bind(request); err != nil {
		//c.log.Errorf("Bind error: %s", err)
		return err
	}
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	request.ID = id

	response, err := c.registry.PostService.GetPostDetails(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Value)
}

func (c *PostController) UpdatePost(ctx echo.Context) error {
	request := &dto.UpdatePostRequest{}
	request.ID, _ = strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err := ctx.Bind(request); err != nil {
		return err
	}
	response, err := c.registry.PostService.UpdatePost(context.Background(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Code, response.Value)
}

func NewPostController(log *logrus.Entry, registry *service.Registry) *PostController {
	return &PostController{log: log, registry: registry}
}
