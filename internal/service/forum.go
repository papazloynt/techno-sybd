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

type ForumService interface {
	CreateForum(ctx context.Context, request *dto.CreateForumRequest) (*dto.CreateForumResponse, error)
	GetForum(ctx context.Context, request *dto.GetForumRequest) (*dto.GetForumResponse, error)
	GetThread(ctx context.Context, request *dto.GetForumThreadRequest) (*dto.GetForumThreadResponse, error)
	GetUsers(ctx context.Context, request *dto.GetForumUsersRequest) (*dto.GetForumUsersResponse, error)
}

type forumServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *forumServiceImpl) CreateForum(ctx context.Context, request *dto.CreateForumRequest) (*dto.CreateForumResponse, error) {
	if forum, err := svc.db.ForumRepo.GetForumBySlug(ctx, request.Slug); err != nil {
		if !errors.Is(err, constants.ErrDBNotFound) {
			return nil, err
		}
	} else {
		return &dto.CreateForumResponse{Value: forum, Code: http.StatusConflict}, nil
	}

	user, err := svc.db.UserRepo.GetUserByNickname(ctx, request.User)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.CreateForumResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.User)}, Code: http.StatusNotFound}, nil
		}
	}
	request.User = user.Nickname

	if err := svc.db.ForumRepo.CreateForum(ctx, &core.Forum{Title: request.Title, User: request.User, Slug: request.Slug}); err != nil {
		return nil, err
	}

	forum, err := svc.db.ForumRepo.GetForumBySlug(ctx, request.Slug)
	if err != nil {
		return nil, err
	}

	return &dto.CreateForumResponse{Value: forum, Code: http.StatusCreated}, nil
}

func (svc *forumServiceImpl) GetForum(ctx context.Context, request *dto.GetForumRequest) (*dto.GetForumResponse, error) {
	forum, err := svc.db.ForumRepo.GetForumBySlug(ctx, request.Slug)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			//svc.log.Errorf("err: %s", err)
			return &dto.GetForumResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find forum with slug: %s", request.Slug)}, Code: http.StatusNotFound}, nil
		}
		return nil, err
	}
	//svc.log.Infof("Success")
	return &dto.GetForumResponse{Value: forum, Code: http.StatusOK}, nil
}

func (svc *forumServiceImpl) GetThread(ctx context.Context, request *dto.GetForumThreadRequest) (*dto.GetForumThreadResponse, error) {
	if forum, err := svc.db.ForumRepo.GetForumBySlug(ctx, request.Slug); err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.GetForumThreadResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find forum with slug: %s", request.Slug)}, Code: http.StatusNotFound}, nil
		}
	} else {
		request.Slug = forum.Slug
	}

	threads, err := svc.db.ForumRepo.GetThreadsFromForum(ctx,
		request.Slug,
		request.Limit,
		request.Since,
		request.Desc)
	if err != nil {
		return nil, err
	}

	return &dto.GetForumThreadResponse{Value: threads, Code: http.StatusOK}, nil
}

func (svc *forumServiceImpl) GetUsers(ctx context.Context, request *dto.GetForumUsersRequest) (*dto.GetForumUsersResponse, error) {
	if forum, err := svc.db.ForumRepo.GetForumBySlug(ctx, request.Slug); err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.GetForumUsersResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find forum with slug: %s", request.Slug)}, Code: http.StatusNotFound}, nil
		}
	} else {
		request.Slug = forum.Slug
	}

	threads, err := svc.db.ForumRepo.GetUsersFromForum(ctx,
		request.Slug,
		request.Limit,
		request.Since,
		request.Desc)
	if err != nil {
		return nil, err
	}

	return &dto.GetForumUsersResponse{Value: threads, Code: http.StatusOK}, nil
}

func NewForumService(log *logrus.Entry, db *db.Repository) ForumService {
	return &forumServiceImpl{log: log, db: db}
}
