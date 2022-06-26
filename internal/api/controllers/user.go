package controllers

import (
	"SYBD/internal/model/dto"
	"SYBD/internal/service"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	log      *logrus.Entry
	registry *service.Registry
}

func (c *UserController) CreateUser(ctx echo.Context) error {
	request := &dto.CreateUserRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}
	request.Nickname = ctx.Param("nickname")
	//c.log.Infof("request nick: %s", request.Nickname)
	response, err := c.registry.UserService.CreateUser(context.Background(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Code, response.Value)
}

func (c *UserController) GetProfile(ctx echo.Context) error {
	request := &dto.GetProfileRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}
	request.Nickname = ctx.Param("nickname")
	//c.log.Infof("request nick: %s", request.Nickname)

	response, err := c.registry.UserService.GetProfile(context.Background(), request)
	if err != nil {
		//c.log.Infof("err: %s", err)
		return err
	}
	//c.log.Infof("response: %s", response)

	return ctx.JSON(response.Code, response.Value)
}

func (c *UserController) UpdateProfile(ctx echo.Context) error {
	request := &dto.UpdateProfileRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}
	request.Nickname = ctx.Param("nickname")
	//c.log.Infof("request nick: %s", request.Nickname)
	response, err := c.registry.UserService.UpdateProfile(context.Background(), request)
	if err != nil {
		return err
	}
	//c.log.Infof("response : %s", response.Value)
	return ctx.JSON(response.Code, response.Value)
}

func NewUserController(log *logrus.Entry, registry *service.Registry) *UserController {
	return &UserController{log: log, registry: registry}
}
