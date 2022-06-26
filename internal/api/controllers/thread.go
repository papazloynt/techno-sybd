package controllers

import (
	"SYBD/internal/model/dto"
	"SYBD/internal/service"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ThreadController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *ThreadController) CreateThread(ctx echo.Context) error {
	request := &dto.CreateThreadRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}
	request.Forum = ctx.Param("slug")

	response, err := c.registry.ThreadService.CreateThread(context.Background(), request)
	if err != nil {
		return err
	}

	//c.log.Infof("response: %s", response)
	return ctx.JSON(response.Code, response.Value)
}

func (c *ThreadController) UpdateVote(ctx echo.Context) error {
	request := &dto.UpdateVoteRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	slugOrID := ctx.Param("slug_or_id")
	response, err := c.registry.ThreadService.UpdateVote(context.Background(), slugOrID, request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Value)
}

func (c *ThreadController) GetDetails(ctx echo.Context) error {
	slugOrID := ctx.Param("slug_or_id")
	response, err := c.registry.ThreadService.GetDetails(context.Background(), slugOrID)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Code, response.Value)
}

func (c *ThreadController) UpdateForumThread(ctx echo.Context) error {
	request := &dto.UpdateThreadRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}
	slugOrID := ctx.Param("slug_or_id")
	response, err := c.registry.ThreadService.UpdateThread(context.Background(), slugOrID, request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Value)
}

func NewThreadController(log *logrus.Entry, registry *service.Registry) *ThreadController {
	return &ThreadController{log: log, registry: registry}
}
