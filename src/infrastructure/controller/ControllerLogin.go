package controller

import (
	"errors"
	"net"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	auth "github.com/Rafael24595/go-api-render/src/commons/auth/Jwt.go"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
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
		Route(http.MethodPost, instance.login, "/api/v1/login").
		Route(http.MethodGet, instance.user, "/api/v1/user").
		Route(http.MethodGet, instance.identify, "/api/v1/identify")

	return instance
}

func (c *ControllerLogin) login(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	login, err := jsonDeserialize[requestLogin](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	sessions := repository.InstanceManagerSession()
	session, err := sessions.Authorize(login.Username, login.Password)
	if err != nil {
		return result.Err(http.StatusUnauthorized, err)
	}

	if session == nil {
		return result.Err(http.StatusUnprocessableEntity, nil)
	}

	defineSession(w, login.Username)

	return result.Ok(nil)
}

func (c *ControllerLogin) user(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	username := findUser(ctx)

	sessions := repository.InstanceManagerSession()

	user, exists := sessions.Find(username)
	if !exists {
		err := errors.New("user not found")
		return result.Err(http.StatusNotFound, err)
	}

	response := responseUserData{
		Username:  user.Username,
		Timestamp: user.Timestamp,
		Context:   user.Context,
	}

	return result.Ok(response)
}

func (c *ControllerLogin) identify(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}

	parsedIP := net.ParseIP(ip)

	response := responseClientIdentity{
		Ip:     ip,
		IsHost: parsedIP.IsLoopback(),
	}

	return result.Ok(response)
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
