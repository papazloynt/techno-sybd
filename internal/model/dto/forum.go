package dto

type CreateForumRequest struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

type CreateForumResponse struct {
	Value interface{}
	Code  int
}

type GetForumRequest struct {
	Slug string `path:"slug"`
}

type GetForumResponse struct {
	Value interface{}
	Code  int
}

type GetForumThreadRequest struct {
	Slug  string `path:"slug"`
	Limit int64  `query:"limit"`
	Since string `query:"since"`
	Desc  bool   `query:"desc"`
}

type GetForumThreadResponse struct {
	Value interface{}
	Code  int
}

type GetForumUsersRequest struct {
	Slug  string `path:"slug"`
	Limit int64  `query:"limit"`
	Since string `query:"since"`
	Desc  bool   `query:"desc"`
}

type GetForumUsersResponse struct {
	Value interface{}
	Code  int
}
