package dto

import "time"

type CreateThreadRequest struct {
	Author  string    `json:"author"`
	Forum   string    `path:"slug"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	Created time.Time `json:"created,omitempty"`
}

type CreateThreadResponse struct {
	Value interface{}
	Code  int
}

type UpdateVoteRequest struct {
	Nickname string `json:"nickname"`
	Voice    int64  `json:"voice"`
}

type UpdateVoteResponse struct {
	Value interface{}
	Code  int
}

type GetDetailsResponse struct {
	Value interface{}
	Code  int
}

type UpdateThreadRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type UpdateThreadResponse struct {
	Value interface{}
	Code  int
}
