package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Rafael24595/go-api-core/src/commons/session"
	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	auth "github.com/Rafael24595/go-api-render/src/commons/auth/Jwt.go"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const AUTH_COOKIE = "go_api_token"
const AUTH_COOKIE_DESCRIPTION = "Session cookie token"
const REFRESH_COOKIE = "go_api_refresh"
const REFRESH_COOKIE_DESCRIPTION = "User refresh token"

type ControllerLogin struct {
	router *router.Router
}

func NewControllerLogin(
	router *router.Router) ControllerLogin {
	instance := ControllerLogin{
		router: router,
	}

	router.
		RouteDocument(http.MethodPost, instance.login, "login", instance.docLogin()).
		RouteDocument(http.MethodDelete, instance.logout, "login", instance.docLogout()).
		RouteDocument(http.MethodGet, instance.user, "user", instance.docUser()).
		RouteDocument(http.MethodPost, instance.signin, "user", instance.docSignin()).
		RouteDocument(http.MethodDelete, instance.delete, "user", instance.docDelete()).
		RouteDocument(http.MethodPut, instance.verify, "user/verify", instance.docVerify()).
		RouteDocument(http.MethodGet, instance.refresh, "token/refresh", instance.docRefresh())

	return instance
}

func (c *ControllerLogin) docLogin() docs.DocRoute {
	return docs.DocRoute{
		Description: "Authenticate user and establish a session with JWT and refresh token cookies.",
		Request:     docs.DocJsonPayload(requestLogin{}),
		Responses:   c.docUser().Responses,
	}
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
		return result.Reject(http.StatusUnprocessableEntity)
	}

	_, _, err = defineSession(w, session)
	if err != nil {
		return result.Err(http.StatusUnauthorized, err)
	}

	ctx.Put(USER, login.Username)

	return c.user(w, r, ctx)
}

func (c *ControllerLogin) docLogout() docs.DocRoute {
	return docs.DocRoute{
		Description: "Logout the current user and invalidate session cookies.",
		Responses:   c.docUser().Responses,
	}
}

func (c *ControllerLogin) logout(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	eraseSession(w)
	ctx.Put(USER, domain.ANONYMOUS_OWNER)
	return c.user(w, r, ctx)
}

func (c *ControllerLogin) docUser() docs.DocRoute {
	return docs.DocRoute{
		Description: "Get the currently authenticated user's information and context.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload(responseUserData{}),
		},
	}
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

	return result.JsonOk(response)
}

func (c *ControllerLogin) docSignin() docs.DocRoute {
	return docs.DocRoute{
		Description: "Register a new user using the current user's session context.",
		Request:     docs.DocJsonPayload(requestSigninUser{}),
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload(responseUserData{}),
		},
	}
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

func (c *ControllerLogin) docDelete() docs.DocRoute {
	return docs.DocRoute{
		Description: "Delete the current user account and clear session data.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload(responseUserData{}),
		},
	}
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

	eraseSession(w)

	ctx.Put(USER, domain.ANONYMOUS_OWNER)

	return c.user(w, r, ctx)
}

func (c *ControllerLogin) docVerify() docs.DocRoute {
	return docs.DocRoute{
		Description: "Change the user's password by verifying the old one and setting a new one.",
		Request:     docs.DocJsonPayload(requestVerify{}),
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload(responseUserData{}),
		},
	}
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
		return result.Reject(http.StatusInternalServerError)
	}

	_, _, err = defineSession(w, session)
	if err != nil {
		return result.Err(401, err)
	}

	ctx.Put(USER, session.Username)

	return c.user(w, r, ctx)
}

func (c *ControllerLogin) docRefresh() docs.DocRoute {
	return docs.DocRoute{
		Description: "Refresh the session token using the refresh cookie.",
		Cookies: docs.DocParameters{
			REFRESH_COOKIE: REFRESH_COOKIE_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload(responseUserData{}),
		},
	}
}

func (c *ControllerLogin) refresh(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	cookie, err := r.Cookie(REFRESH_COOKIE)
	if err != nil {
		return result.Err(http.StatusUnauthorized, err)
	}

	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		if claims != nil && 0 >= time.Until(claims.ExpiresAt.Time) {
			eraseSession(w)
		}
		return result.Err(http.StatusUnauthorized, err)
	}

	user := claims.Username

	sessions := repository.InstanceManagerSession()
	session, exists := sessions.Find(user)
	if !exists {
		err = errors.New("user not exists")
		return result.Err(http.StatusNotFound, err)
	}

	if cookie.Value != session.Refresh {
		return result.Reject(http.StatusUnauthorized)
	}

	_, _, err = defineSession(w, session)
	if err != nil {
		return result.Err(http.StatusBadRequest, err)
	}

	ctx.Put(USER, user)

	return c.user(w, r, ctx)
}

func defineSession(w http.ResponseWriter, session *session.Session) (string, string, error) {
	sessions := repository.InstanceManagerSession()

	token, err := auth.GenerateJWT(session.Username)
	if err != nil {
		return "", "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     AUTH_COOKIE,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	refresh := session.Refresh
	if _, err := auth.ValidateJWT(refresh); err != nil {
		refresh, err = auth.GenerateRefreshJWT(session.Username)
		if err != nil {
			return "", "", err
		}
		sessions.Refresh(session, refresh)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     REFRESH_COOKIE,
		Value:    refresh,
		Path:     fmt.Sprintf("%stoken/refresh", BASE_PATH),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return token, refresh, nil
}

func closeSession(w http.ResponseWriter) http.ResponseWriter {
	http.SetCookie(w, &http.Cookie{
		Name:     AUTH_COOKIE,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return w
}

func eraseSession(w http.ResponseWriter) http.ResponseWriter {
	closeSession(w)
	http.SetCookie(w, &http.Cookie{
		Name:     REFRESH_COOKIE,
		Value:    "",
		Path:     fmt.Sprintf("%stoken/refresh", BASE_PATH),
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return w
}
