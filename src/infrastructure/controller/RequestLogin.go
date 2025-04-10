package controller

type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
