package web_test

import (
	"errors"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lonepeon/golib/testutils"
	"github.com/lonepeon/golib/web"
	"github.com/lonepeon/golib/web/webtest"
)

func TestShowLoginPageNotAuthenticated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/login", nil)
	auth := web.NewAuthentication(frontend, nil, "login/new.html")

	frontend.EXPECT().CurrentUserID(r).Return(web.AuthenticationUserID(""), nil)
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().Response(200, "login/new.html", nil).Return(expectedResponse)

	response := auth.ShowLoginPage("/home")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
}

func TestShowLoginPageReturnsPageWhenItCantGetUsername(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/login", nil)
	auth := web.NewAuthentication(frontend, nil, "login/new.html")

	frontend.EXPECT().CurrentUserID(r).Return(web.AuthenticationUserID(""), errors.New("boom"))
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().Response(200, "login/new.html", nil).Return(expectedResponse)

	response := auth.ShowLoginPage("/home")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
}

func TestShowLoginPageRedirectToSuccessPathWhenAlreadyLoggedInAndNoPathOverride(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/login", nil)
	auth := web.NewAuthentication(frontend, nil, "login/new.html")

	frontend.EXPECT().CurrentUserID(r).Return(web.AuthenticationUserID("42"), nil)
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().Redirect(w, 302, "/home").Return(expectedResponse)

	response := auth.ShowLoginPage("/home")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
}

func TestShowLoginPageRedirectToSuccessPathWhenAlreadyLoggedInAndPathOverride(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/login?to=/some-other-path", nil)
	auth := web.NewAuthentication(frontend, nil, "login/new.html")

	frontend.EXPECT().CurrentUserID(r).Return(web.AuthenticationUserID("42"), nil)
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().Redirect(w, 302, "/some-other-path").Return(expectedResponse)

	response := auth.ShowLoginPage("/home")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
}

func TestLogoutRedictWhenStorageFails(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/logout", nil)
	auth := web.NewAuthentication(frontend, nil, "login/new.html")

	frontend.EXPECT().Clear(w, r).Return(errors.New("boom"))
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().AddFlash(web.NewFlashMessageError("We can't log you out. Please retry"))
	ctx.EXPECT().Redirect(w, 302, "/login").Return(expectedResponse)

	response := auth.Logout("/login")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
	testutils.AssertContainsString(t, "can't remove user", response.LogMessage, "unexpected log message")
	testutils.AssertContainsString(t, "boom", response.LogMessage, "unexpected log message")
}

func TestLogoutRedictWhenStorageSucceed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/logout", nil)
	auth := web.NewAuthentication(frontend, nil, "login/new.html")

	frontend.EXPECT().Clear(w, r).Return(nil)
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().Redirect(w, 302, "/login").Return(expectedResponse)

	response := auth.Logout("/login")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
	testutils.AssertContainsString(t, "redirecting to /login", response.LogMessage, "unexpected log message")
}

func TestLogoutRedictWhenStorageSucceedWithPathOverride(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/logout?to=/some-other-path", nil)
	auth := web.NewAuthentication(frontend, nil, "login/new.html")

	frontend.EXPECT().Clear(w, r).Return(nil)
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().Redirect(w, 302, "/some-other-path").Return(expectedResponse)

	response := auth.Logout("/login")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
	testutils.AssertContainsString(t, "redirecting to /some-other-path", response.LogMessage, "unexpected log message")
}

func TestLoginWhenNoUsername(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/login", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := web.NewAuthentication(nil, nil, "login/new.html")

	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().AddFlash(web.NewFlashMessageError("username/password combination is required"))
	ctx.EXPECT().Response(200, "login/new.html", nil).Return(expectedResponse)

	response := auth.Login("/dashboard")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
}

func TestLoginWhenInvalidCredentials(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	backend := webtest.NewMockAuthenticationBackendStorer(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	body := strings.NewReader(url.Values{"username": []string{"jane"}, "password": []string{"doe"}}.Encode())
	r := httptest.NewRequest("POST", "/login", body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := web.NewAuthentication(nil, backend, "login/new.html")

	backend.EXPECT().Authenticate("jane", "doe").Return(web.AuthenticationUserID(""), web.ErrUserInvalidCredentials)
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().AddFlash(web.NewFlashMessageError("username/password combination is invalid"))
	ctx.EXPECT().Response(200, "login/new.html", nil).Return(expectedResponse)

	response := auth.Login("/dashboard")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
}

func TestLoginWithValidCredentialsButCantPersistUsername(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	backend := webtest.NewMockAuthenticationBackendStorer(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	body := strings.NewReader(url.Values{"username": []string{"jdoe"}, "password": []string{"secret"}}.Encode())
	r := httptest.NewRequest("POST", "/login?to=/my-other-path", body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := web.NewAuthentication(frontend, backend, "login/new.html")

	backend.EXPECT().Authenticate("jdoe", "secret").Return(web.AuthenticationUserID("42"), nil)
	backend.EXPECT().Lookup(web.AuthenticationUserID("42")).Return(web.AuthenticationUser{ID: web.AuthenticationUserID("42"), Username: "jdoe"}, nil)
	frontend.EXPECT().StoreUserID(w, r, web.AuthenticationUserID("42")).Return(errors.New("boom"))
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().AddFlash(web.NewFlashMessageError("something wrong happened. Please try again."))
	ctx.EXPECT().Response(200, "login/new.html", nil).Return(expectedResponse)

	response := auth.Login("/dashboard")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
	testutils.AssertContainsString(t, "storage error", response.LogMessage, "unexpected log message")
	testutils.AssertContainsString(t, "boom", response.LogMessage, "unexpected log message")
}

func TestLoginWithValidCredentials(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	backend := webtest.NewMockAuthenticationBackendStorer(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	body := strings.NewReader(url.Values{"username": []string{"jdoe"}, "password": []string{"secret"}}.Encode())
	r := httptest.NewRequest("POST", "/login", body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	auth := web.NewAuthentication(frontend, backend, "login/new.html")

	backend.EXPECT().Authenticate("jdoe", "secret").Return(web.AuthenticationUserID("42"), nil)
	backend.EXPECT().Lookup(web.AuthenticationUserID("42")).Return(web.AuthenticationUser{ID: web.AuthenticationUserID("42"), Username: "jdoe"}, nil)
	frontend.EXPECT().StoreUserID(w, r, web.AuthenticationUserID("42")).Return(nil)
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().Redirect(w, 302, "/dashboard").Return(expectedResponse)

	response := auth.Login("/dashboard")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
}

func TestLoginWithValidCredentialsAndPathOverride(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	backend := webtest.NewMockAuthenticationBackendStorer(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	body := strings.NewReader(url.Values{"username": []string{"jdoe"}, "password": []string{"secret"}}.Encode())
	r := httptest.NewRequest("POST", "/login?to=/my-other-path", body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := web.NewAuthentication(frontend, backend, "login/new.html")

	backend.EXPECT().Authenticate("jdoe", "secret").Return(web.AuthenticationUserID("42"), nil)
	backend.EXPECT().Lookup(web.AuthenticationUserID("42")).Return(web.AuthenticationUser{ID: web.AuthenticationUserID("42"), Username: "jdoe"}, nil)
	frontend.EXPECT().StoreUserID(w, r, web.AuthenticationUserID("42")).Return(nil)
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().Redirect(w, 302, "/my-other-path").Return(expectedResponse)

	response := auth.Login("/dashboard")(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
}

func TestEnsureAuthenticationNotLoggedIn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	handler := webtest.NewMockHandler(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/dashboard", nil)
	auth := web.NewAuthentication(frontend, nil, "login/new.html")

	frontend.EXPECT().CurrentUserID(r).Return(web.AuthenticationUserID(""), nil)
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().Redirect(w, 302, "/login").Return(expectedResponse)
	handler.EXPECT().Handle(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(0)

	response := auth.EnsureAuthentication("/login", handler.Handle)(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
	testutils.AssertContainsString(t, "/login", response.LogMessage, "unexpected log message")
	testutils.AssertContainsString(t, "not authenticated", response.LogMessage, "unexpected log message")
}

func TestEnsureAuthenticationCantGetUsername(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	handler := webtest.NewMockHandler(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/dashboard", nil)
	auth := web.NewAuthentication(frontend, nil, "login/new.html")

	frontend.EXPECT().CurrentUserID(r).Return(web.AuthenticationUserID(""), errors.New("boom"))
	expectedResponse := webtest.MockedResponse("expected response")
	ctx.EXPECT().Redirect(w, 302, "/login").Return(expectedResponse)
	handler.EXPECT().Handle(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(0)

	response := auth.EnsureAuthentication("/login", handler.Handle)(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
	testutils.AssertContainsString(t, "/login", response.LogMessage, "unexpected log message")
	testutils.AssertContainsString(t, "can't load current user", response.LogMessage, "unexpected log message")
	testutils.AssertContainsString(t, "boom", response.LogMessage, "unexpected log message")
}

func TestEnsureAuthenticationIsLoggedIn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	backend := webtest.NewMockAuthenticationBackendStorer(mockCtrl)
	handler := webtest.NewMockHandler(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/dashboard", nil)
	auth := web.NewAuthentication(frontend, backend, "login/new.html")

	frontend.EXPECT().CurrentUserID(r).Return(web.AuthenticationUserID("42"), nil)
	backend.EXPECT().Lookup(web.AuthenticationUserID("42")).Return(web.AuthenticationUser{ID: web.AuthenticationUserID("42"), Username: "jdoe"}, nil)

	ctx.EXPECT().AddData("Authentication", map[string]interface{}{
		"IsLoggedIn": true,
		"Username":   "jdoe",
		"User":       web.AuthenticationUser{ID: web.AuthenticationUserID("42"), Username: "jdoe"},
	})
	expectedResponse := webtest.MockedResponse("expected response")
	handler.EXPECT().Handle(ctx, w, r).Return(expectedResponse)

	response := auth.EnsureAuthentication("/login", handler.Handle)(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
}

func TestIdentifyCurrentUserWhenStorageFailed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	handler := webtest.NewMockHandler(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/dashboard", nil)
	auth := web.NewAuthentication(frontend, nil, "login/new.html")

	frontend.EXPECT().CurrentUserID(r).Return(web.AuthenticationUserID(""), errors.New("boom"))
	ctx.EXPECT().AddData("Authentication", map[string]interface{}{
		"IsLoggedIn": false,
		"User":       web.AuthenticationUser{},
		"Username":   "",
	})
	expectedResponse := webtest.MockedResponse("expected response")
	handler.EXPECT().Handle(ctx, w, r).Return(expectedResponse)

	response := auth.IdentifyCurrentUser(handler.Handle)(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
}

func TestIdentifyCurrentUserWhenUserNotLoggedIn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	handler := webtest.NewMockHandler(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/dashboard", nil)
	auth := web.NewAuthentication(frontend, nil, "login/new.html")

	frontend.EXPECT().CurrentUserID(r).Return(web.AuthenticationUserID(""), nil)
	ctx.EXPECT().AddData("Authentication", map[string]interface{}{
		"IsLoggedIn": false,
		"User":       web.AuthenticationUser{},
		"Username":   "",
	})
	expectedResponse := webtest.MockedResponse("expected response")
	handler.EXPECT().Handle(ctx, w, r).Return(expectedResponse)

	response := auth.IdentifyCurrentUser(handler.Handle)(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
}

func TestIdentifyCurrentUserWhenHasUserLoggedIn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	frontend := webtest.NewMockAuthenticationFrontendStorer(mockCtrl)
	backend := webtest.NewMockAuthenticationBackendStorer(mockCtrl)
	handler := webtest.NewMockHandler(mockCtrl)
	ctx := webtest.NewMockContext(mockCtrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/dashboard", nil)
	auth := web.NewAuthentication(frontend, backend, "login/new.html")

	frontend.EXPECT().CurrentUserID(r).Return(web.AuthenticationUserID("42"), nil)
	backend.EXPECT().Lookup(web.AuthenticationUserID("42")).Return(web.AuthenticationUser{ID: web.AuthenticationUserID("42"), Username: "jdoe"}, nil)

	ctx.EXPECT().AddData("Authentication", map[string]interface{}{
		"IsLoggedIn": true,
		"Username":   "jdoe",
		"User":       web.AuthenticationUser{ID: web.AuthenticationUserID("42"), Username: "jdoe"},
	})
	expectedResponse := webtest.MockedResponse("expected response")
	handler.EXPECT().Handle(ctx, w, r).Return(expectedResponse)

	response := auth.IdentifyCurrentUser(handler.Handle)(ctx, w, r)

	webtest.AssertResponse(t, expectedResponse, response, "unexpected web response")
}
