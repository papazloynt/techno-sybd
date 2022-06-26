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

type PostService interface {
	CreatePost(ctx context.Context, slugOrID string, posts []*dto.Post) (*dto.CreatePostResponse, error)
	GetPost(ctx context.Context, slugOrID string, sort string, since int64, desc bool, limit int64) (*dto.GetPostResponse, error)
	GetPostDetails(ctx context.Context, request *dto.GetPostDetailsRequest) (*dto.GetPostDetailsResponse, error)
	UpdatePost(ctx context.Context, request *dto.UpdatePostRequest) (*dto.UpdatePostResponse, error)
}

type postServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *postServiceImpl) CreatePost(ctx context.Context, slugOrID string, posts []*dto.Post) (*dto.CreatePostResponse, error) {
	var id int
	var err error
	id, err = strconv.Atoi(slugOrID)

	var thread *core.Thread
	if err != nil {
		if thread, err = svc.db.ThreadRepo.GetThread(ctx, slugOrID); err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.CreatePostResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", slugOrID)}, Code: http.StatusNotFound}, nil
			}
		} else {
			id = int(thread.ID)
		}
	} else {
		if thread, err = svc.db.ThreadRepo.GetThreadByID(ctx, int64(id)); err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.CreatePostResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by id: %d", id)}, Code: http.StatusNotFound}, nil
			}
		}
	}

	if len(posts) == 0 {
		return &dto.CreatePostResponse{Value: []struct{}{}, Code: http.StatusCreated}, nil
	}

	if posts[0].Parent != 0 {
		parentThreadID, err := svc.db.PostRepo.CheckPredPost(ctx, int(posts[0].Parent))
		if err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.CreatePostResponse{Value: dto.ErrorResponse{Message: "Parent post was created in another thread"}, Code: http.StatusConflict}, nil
			}
		}

		if parentThreadID != id {
			return &dto.CreatePostResponse{Value: dto.ErrorResponse{Message: "Parent post was created in another thread"}, Code: http.StatusConflict}, nil
		}
	}

	if _, err := svc.db.UserRepo.GetUserByNickname(ctx, posts[0].Author); err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.CreatePostResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find user by nickname: %s", posts[0].Author)}, Code: http.StatusNotFound}, nil
		}
	}

	insertedPosts, err := svc.db.PostRepo.CreatePost(ctx, thread.Forum, int64(id), posts)
	if err != nil {
		return nil, err
	}

	return &dto.CreatePostResponse{Value: insertedPosts, Code: http.StatusCreated}, nil
}

func (svc *postServiceImpl) GetPost(ctx context.Context, slugOrID string, sort string, since int64, desc bool, limit int64) (*dto.GetPostResponse, error) {
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		if thread, err := svc.db.ThreadRepo.GetThread(ctx, slugOrID); err != nil {
			if errors.Is(err, constants.ErrDBNotFound) {
				return &dto.GetPostResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by slug: %s", slugOrID)}, Code: http.StatusNotFound}, nil
			}
		} else {
			id = int(thread.ID)
		}
	}

	if _, err := svc.db.ThreadRepo.GetThreadByID(ctx, int64(id)); err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.GetPostResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find thread forum by id: %d", id)}, Code: http.StatusNotFound}, nil
		}
	}
	var posts []*core.Post
	switch sort {
	case "flat":
		posts, err = svc.db.PostRepo.GetPost(ctx, id, since, desc, limit)
	case "tree":
		posts, err = svc.db.PostRepo.GetPostTree(ctx, id, since, desc, limit)
	case "parent_tree":
		posts, err = svc.db.PostRepo.GetPostPredTree(ctx, id, since, desc, limit)
	default:
		posts, err = svc.db.PostRepo.GetPost(ctx, id, since, desc, limit)
	}
	if err != nil {
		return nil, err
	}
	return &dto.GetPostResponse{Value: posts, Code: http.StatusOK}, nil
}

func (svc *postServiceImpl) GetPostDetails(ctx context.Context, request *dto.GetPostDetailsRequest) (*dto.GetPostDetailsResponse, error) {
	//svc.log.Infof("request: %v", request)
	post, err := svc.db.PostRepo.GetPostByID(ctx, request.ID)
	//svc.log.Infof("post: %v \n err: %s", post, err)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.GetPostDetailsResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find post by id: %d", request.ID)}, Code: http.StatusNotFound}, nil
		}
		//svc.log.Errorf("post: %v \n err: %s", post, err)
		return nil, err
	}

	postDetails, err := svc.db.PostRepo.GetPostDetails(ctx, request.ID, request.Related)
	if err != nil {
		//svc.log.Errorf("post: %v \n err: %s", post, err)
		return nil, err
	}
	postDetails.Post = post

	return &dto.GetPostDetailsResponse{Value: postDetails, Code: http.StatusOK}, nil
}

func (svc *postServiceImpl) UpdatePost(ctx context.Context, request *dto.UpdatePostRequest) (*dto.UpdatePostResponse, error) {
	post, err := svc.db.PostRepo.GetPostByID(ctx, request.ID)
	if err != nil {
		if errors.Is(err, constants.ErrDBNotFound) {
			return &dto.UpdatePostResponse{Value: dto.ErrorResponse{Message: fmt.Sprintf("Can't find post by id: %d", request.ID)}, Code: http.StatusNotFound}, nil
		}
		return nil, err
	}

	if len(request.Message) == 0 || request.Message == post.Message {
		return &dto.UpdatePostResponse{Value: post, Code: http.StatusOK}, nil
	}

	updatedPost, err := svc.db.PostRepo.UpdatePost(ctx, request.ID, request.Message)
	if err != nil {
		return nil, err
	}

	return &dto.UpdatePostResponse{Value: updatedPost, Code: http.StatusOK}, nil
}

func NewPostService(log *logrus.Entry, db *db.Repository) PostService {
	return &postServiceImpl{log: log, db: db}
}
