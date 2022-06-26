package service

import (
	"SYBD/internal/constants"
	"SYBD/internal/db"
	"SYBD/internal/model/core"
	"SYBD/internal/model/dto"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserService interface {
	CreateUser(ctx context.Context, request *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	GetProfile(ctx context.Context, request *dto.GetProfileRequest) (*dto.GetProfileResponse, error)
	UpdateProfile(ctx context.Context, request *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error)
}

type userServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *userServiceImpl) CreateUser(ctx context.Context, request *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	if users, err := svc.db.UserRepo.GetSimilaryUsers(ctx, request.Email, request.Nickname); err != nil {
		return nil, err
	} else if len(users) > 0 {
		return &dto.CreateUserResponse{Value: users, Code: http.StatusConflict}, nil
	}

	user := &core.User{Nickname: request.Nickname, FullName: request.FullName, About: request.About, Email: request.Email}
	if err := svc.db.UserRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	//svc.log.Infof("user:  %s ", user)
	return &dto.CreateUserResponse{Value: user, Code: http.StatusCreated}, nil
}

func (svc *userServiceImpl) GetProfile(ctx context.Context, request *dto.GetProfileRequest) (*dto.GetProfileResponse, error) {
	user, err := svc.db.UserRepo.GetUserByNickname(ctx, request.Nickname)
	//svc.log.Infof("user:  %s \n err: %s", user, err)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.GetProfileResponse{Value: constants.CreateNewError(fmt.Sprintf("Can't find user with that nickname: %s", request.Nickname), http.StatusNotFound), Code: http.StatusNotFound}, nil
		}
		return nil, err
	}
	return &dto.GetProfileResponse{Value: *user, Code: http.StatusOK}, nil
}

func (svc *userServiceImpl) UpdateProfile(ctx context.Context, request *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error) {
	if len(request.Email) > 0 {
		if user, err := svc.db.UserRepo.GetUserByEmail(ctx, request.Email); err != nil {
			if !errors.Is(err, constants.ErrDBNotFound) {
				return nil, err
			}
		} else if user.Nickname != request.Nickname {
			return &dto.UpdateProfileResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("This email is already registered by user: %s", user.Nickname)}, Code: http.StatusConflict}, nil
		}
	}

	user := &core.User{Nickname: request.Nickname, FullName: request.FullName, About: request.About, Email: request.Email}
	updatedUser, err := svc.db.UserRepo.UpdateUser(ctx, user)
	if err != nil {
		return &dto.UpdateProfileResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.Nickname)}, Code: http.StatusNotFound}, nil
	}
	return &dto.UpdateProfileResponse{Value: updatedUser, Code: http.StatusOK}, nil
}

func NewUserService(log *logrus.Entry, db *db.Repository) UserService {
	return &userServiceImpl{log: log, db: db}
}
