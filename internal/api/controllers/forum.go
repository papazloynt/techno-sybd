package controllers

import (
	"SYBD/internal/model/dto"
	"SYBD/internal/service"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ForumController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *ForumController) CreateForum(ctx echo.Context) error {
	request := &dto.CreateForumRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	response, err := c.registry.ForumService.CreateForum(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Value)
}

func (c *ForumController) GetForum(ctx echo.Context) error {
	request := &dto.GetForumRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}
	request.Slug = ctx.Param("slug")

	response, err := c.registry.ForumService.GetForum(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Value)
}

func (c *ForumController) GetForumThreads(ctx echo.Context) error {
	request := new(dto.GetForumThreadRequest)

	if err := ctx.Bind(request); err != nil {
		//c.log.Errorf("Bind error: %s", err)
		return err
	}
	request.Slug = ctx.Param("slug")
	if request.Limit < -1 || request.Limit == 0 {
		request.Limit = 100
	}

	response, err := c.registry.ForumService.GetThread(context.Background(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Code, response.Value)
}

func (c *ForumController) GetUsers(ctx echo.Context) error {
	request := new(dto.GetForumUsersRequest)

	if err := ctx.Bind(request); err != nil {
		//c.log.Errorf("Bind error: %s", err)
		return err
	}
	request.Slug = ctx.Param("slug")
	if request.Limit < -1 || request.Limit == 0 {
		request.Limit = 100
	}

	response, err := c.registry.ForumService.GetUsers(context.Background(), request)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Code, response.Value)
}

func NewForumController(log *logrus.Entry, registry *service.Registry) *ForumController {
	return &ForumController{log: log, registry: registry}
}
