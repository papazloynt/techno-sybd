package core

type User struct {
	FullName string `json:"fullname"`
	Nickname string `json:"nickname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}
