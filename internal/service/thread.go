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
	"strconv"
)

type ThreadService interface {
	CreateThread(ctx context.Context, request *dto.CreateThreadRequest) (*dto.CreateThreadResponse, error)
	UpdateVote(ctx context.Context, slugOrID string, request *dto.UpdateVoteRequest) (*dto.UpdateVoteResponse, error)
	GetDetails(ctx context.Context, slugOrID string) (*dto.GetDetailsResponse, error)
	UpdateThread(ctx context.Context, slugOrID string, request *dto.UpdateThreadRequest) (*dto.UpdateThreadResponse, error)
}

type threadServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *threadServiceImpl) CreateThread(ctx context.Context, request *dto.CreateThreadRequest) (*dto.CreateThreadResponse, error) {
	user, err := svc.db.UserRepo.GetUserByNickname(ctx, request.Author)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			//svc.log.Errorf("err: %s", err)
			return &dto.CreateThreadResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.Author)}, Code: http.StatusNotFound}, nil
		}
	}
	request.Author = user.Nickname

	if forum, err := svc.db.ForumRepo.GetForumBySlug(ctx, request.Forum); err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			//svc.log.Errorf("err: %s", err)
			return &dto.CreateThreadResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", request.Forum)}, Code: http.StatusNotFound}, nil
		}
	} else {
		request.Forum = forum.Slug
	}

	if request.Slug != "" {
		if thread, err := svc.db.ThreadRepo.GetThread(ctx, request.Slug); err != nil {
			if !errors.Is(err, constants.ErrDBNotFound) {
				//svc.log.Errorf("err: %s", err)
				return nil, err
			}
		} else {
			//svc.log.Errorf("err: %s", err)
			return &dto.CreateThreadResponse{Value: thread, Code: http.StatusConflict}, nil
		}
	}

	reqThread := &core.Thread{Forum: request.Forum, Title: request.Title, Author: request.Author, Message: request.Message, Slug: request.Slug, Created: request.Created}
	thread, err := svc.db.ThreadRepo.CreateThread(ctx, reqThread)
	if err != nil {
		return nil, err
	}

	return &dto.CreateThreadResponse{Value: thread, Code: http.StatusCreated}, nil
}

func (svc *threadServiceImpl) UpdateVote(ctx context.Context, slugOrID string, request *dto.UpdateVoteRequest) (*dto.UpdateVoteResponse, error) {
	var id int
	var err error
	id, err = strconv.Atoi(slugOrID)

	var thread *core.Thread
	if err != nil {
		if thread, err = svc.db.ThreadRepo.GetThread(ctx, slugOrID); err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.UpdateVoteResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", slugOrID)}, Code: http.StatusNotFound}, nil
			}
		} else {
			id = int(thread.ID)
		}
	} else {
		if thread, err = svc.db.ThreadRepo.GetThreadByID(ctx, int64(id)); err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.UpdateVoteResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by id: %d", id)}, Code: http.StatusNotFound}, nil
			}
		}
	}

	user, err := svc.db.UserRepo.GetUserByNickname(ctx, request.Nickname)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.UpdateVoteResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", request.Nickname)}, Code: http.StatusNotFound}, nil
		}
	}
	request.Nickname = user.Nickname

	exists, err := svc.db.VoteRepo.VoteExists(ctx, request.Nickname, thread.ID)
	if err != nil {
		return nil, err
	}

	if exists {
		if ok, err := svc.db.VoteRepo.UpdateVote(ctx, thread.ID, request.Nickname, request.Voice); err != nil {
			return nil, err
		} else if ok {
			thread.Votes += request.Voice * 2
		}
	} else {
		newVote := &core.Vote{
			Nickname: request.Nickname,
			ThreadID: thread.ID,
			Voice:    request.Voice,
		}

		if err := svc.db.VoteRepo.CreateVote(ctx, newVote); err != nil {
			return nil, err
		}

		thread.Votes += request.Voice
	}

	return &dto.UpdateVoteResponse{Value: thread, Code: http.StatusOK}, nil
}

func (svc *threadServiceImpl) GetDetails(ctx context.Context, slugOrID string) (*dto.GetDetailsResponse, error) {
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		if thread, err := svc.db.ThreadRepo.GetThread(ctx, slugOrID); err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.GetDetailsResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread by slug: %s", slugOrID)}, Code: http.StatusNotFound}, nil
			}
			return nil, err
		} else {
			return &dto.GetDetailsResponse{Value: thread, Code: http.StatusOK}, nil
		}
	}
	thread, err := svc.db.ThreadRepo.GetThreadByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.GetDetailsResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread by id: %d", id)}, Code: http.StatusNotFound}, nil
		}
	}
	return &dto.GetDetailsResponse{Value: thread, Code: http.StatusOK}, nil
}

func (svc *threadServiceImpl) UpdateThread(ctx context.Context, slugOrID string, request *dto.UpdateThreadRequest) (*dto.UpdateThreadResponse, error) {
	var thread *core.Thread
	var err error
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		if thread, err = svc.db.ThreadRepo.GetThread(ctx, slugOrID); err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.UpdateThreadResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", slugOrID)}, Code: http.StatusNotFound}, nil
			}
			return nil, err
		} else {
			id = int(thread.ID)
		}
	}

	if thread, err = svc.db.ThreadRepo.GetThreadByID(ctx, int64(id)); err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.UpdateThreadResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by id: %d", id)}, Code: http.StatusNotFound}, nil
		}
	}

	if request.Title == "" {
		request.Title = thread.Title
	}

	if request.Message == "" {
		request.Message = thread.Message
	}

	thread, err = svc.db.ThreadRepo.UpdateThread(ctx, int64(id), request.Title, request.Message)
	return &dto.UpdateThreadResponse{Value: thread, Code: http.StatusOK}, err
}

func NewThreadService(log *logrus.Entry, db *db.Repository) ThreadService {
	return &threadServiceImpl{log: log, db: db}
}
