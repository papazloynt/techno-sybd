package dto

type CreateUserRequest struct {
	Nickname string `path:"nickname"`
	About    string `json:"about"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
}

type CreateUserResponse struct {
	Value interface{}
	Code  int
}

type GetProfileRequest struct {
	Nickname string `path:"nickname"`
}

type GetProfileResponse struct {
	Value interface{}
	Code  int
}

type UpdateProfileRequest struct {
	Nickname string `path:"nickname"`
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type UpdateProfileResponse struct {
	Value interface{}
	Code  int
}
