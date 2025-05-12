package controller

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/Rafael24595/go-api-core/src/domain"
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
		Route(http.MethodDelete, instance.logout, "/api/v1/login").
		Route(http.MethodGet, instance.user, "/api/v1/user").
		Route(http.MethodPost, instance.signin, "/api/v1/user").
		Route(http.MethodDelete, instance.delete, "/api/v1/user").
		Route(http.MethodPut, instance.verify, "/api/v1/user/verify").
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

	ctx.Put(USER, login.Username)

	return c.user(w, r, ctx)
}

func (c *ControllerLogin) logout(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	closeSession(w)
	ctx.Put(USER, domain.ANONYMOUS_OWNER)
	return c.user(w, r, ctx)
}

func (c *ControllerLogin) user(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	username := findUser(ctx)

	sessions := repository.InstanceManagerSession()

	user, exists := sessions.Find(username)
	if !exists {
		err := errors.New("user not found")
		return result.Err(http.StatusNotFound, err)
	}

	collection, err := sessions.FindUserCollection(user.Username)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}

	response := responseUserData{
		Username:    user.Username,
		Timestamp:   user.Timestamp,
		History:     user.History,
		Collection:  user.Collection,
		Context:     collection.Context,
		IsProtected: user.IsProtected,
		IsAdmin:     user.IsAdmin,
		FirstTime:   user.Count < 0,
	}

	return result.Ok(response)
}

func (c *ControllerLogin) signin(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	username := findUser(ctx)

	request, err := jsonDeserialize[requestSigninUser](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	sessions := repository.InstanceManagerSession()

	user, exists := sessions.Find(username)
	if !exists {
		err := errors.New("user not found")
		return result.Err(http.StatusNotFound, err)
	}

	session, err := sessions.Insert(user, request.Username, request.Password1, request.IsAdmin)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	ctx.Put(USER, session.Username)

	return c.user(w, r, ctx)
}

func (c *ControllerLogin) verify(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	username := findUser(ctx)

	verify, err := jsonDeserialize[requestVerify](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	sessions := repository.InstanceManagerSession()
	session, err := sessions.Verify(username, verify.OldPassword, verify.NewPassword1, verify.NewPassword2)
	if err != nil {
		return result.Err(http.StatusUnauthorized, err)
	}

	if session == nil {
		return result.Err(http.StatusInternalServerError, nil)
	}

	defineSession(w, session.Username)

	ctx.Put(USER, session.Username)

	return c.user(w, r, ctx)
}

func (c *ControllerLogin) delete(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	username := findUser(ctx)

	sessions := repository.InstanceManagerSession()

	user, exists := sessions.Find(username)
	if !exists {
		err := errors.New("user not found")
		return result.Err(http.StatusNotFound, err)
	}

	if _, err := sessions.Delete(user); err != nil {
		return result.Err(http.StatusForbidden, err)
	}

	closeSession(w)

	ctx.Put(USER, domain.ANONYMOUS_OWNER)

	return c.user(w, r, ctx)
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

func closeSession(w http.ResponseWriter) error {
	http.SetCookie(w, &http.Cookie{
		Name:     COOKIE_NAME,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}
