package dto

import "SYBD/internal/model/core"

type Post struct {
	Parent  int64  `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

type CreatePostResponse struct {
	Value interface{}
	Code  int
}

type GetPostResponse struct {
	Value interface{}
	Code  int
}

type PostInfo struct {
	Post   *core.Post   `json:"post"`
	Author *core.User   `json:"author,omitempty"`
	Thread *core.Thread `json:"thread,omitempty"`
	Forum  *core.Forum  `json:"forum,omitempty"`
}

type GetPostDetailsRequest struct {
	ID      int64  `path:"id"`
	Related string `query:"related"`
}

type GetPostDetailsResponse struct {
	Value interface{}
	Code  int
}

type UpdatePostRequest struct {
	ID      int64  `path:"id"`
	Message string `json:"message"`
}

type UpdatePostResponse struct {
	Value interface{}
	Code  int
}
