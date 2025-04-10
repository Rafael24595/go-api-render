package controller

import (
	"errors"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	auth "github.com/Rafael24595/go-api-render/src/commons/auth/Jwt.go"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

const COOKIE_NAME = "Go-Api-Client"

type ControllerLogin struct {
	router *router.Router
}

func NewControllerLogin(
	router *router.Router) ControllerLogin {
	instance := ControllerLogin{
		router: router,
	}

	router.
		Route(http.MethodPost, instance.login, "/api/v1/login")

	return instance
}

func (c *ControllerLogin) login(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	login, err := jsonDeserialize[RequestLogin](r)
	if err != nil {
		return err
	}

	sessions := repository.InstanceManagerSession()
	session, err := sessions.Authorize(login.Username, login.Password)
	if err != nil {
		return err
	}

	if session == nil {
		return errors.New("unautorized")
	}

	defineSession(w, login.Username)

	return nil
}

func defineSession(w http.ResponseWriter, username string) error {
	token, err := auth.GenerateJWT(username)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     COOKIE_NAME,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}
