package main

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/objx"
	"net/http"
	"strings"
	"crypto/md5"
	"io"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	_, err := req.Cookie("auth")
	if err == http.ErrNoCookie {
		res.Header().Set("Location", "/login")
		res.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	h.next.ServeHTTP(res, req)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	defer func() {
		if err := recover(); err != nil {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(res, "Auth path %s not valid", path)
		}
	}()
	segments := strings.Split(path, "/")
	action := segments[2]
	provider := segments[3]
	switch action {
	case "login":
		providerLoginManager(res, provider)
	case "callback":
		providerResponseManager(res, req, provider)
	default:
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "Auth action %s not supported", action)
	}
}

func providerResponseManager(res http.ResponseWriter, req *http.Request, providerName string) {
	provider, err := gomniauth.Provider(providerName)
	if err != nil {
		msg := fmt.Sprintf("Error when trying to get provider %s, %s", provider, err)
		http.Error(res, msg, http.StatusBadRequest)
		return
	}

	creds, err := provider.CompleteAuth(objx.MustFromURLQuery(req.URL.RawQuery))
	if err != nil {
		format := "Error when trying to complete auth for %s: %s"
		http.Error(res, fmt.Sprintf(format, provider, err), http.StatusInternalServerError)
		return
	}

	user, err := provider.GetUser(creds)
	if err != nil {
		format := "Error when trying to get user from %s: %s"
		http.Error(res, fmt.Sprintf(format, provider, err), http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, createCookie(createCookieValue(user)))
	res.Header().Set("Location", "/chat")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func createCookieValue(user common.User) string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(user.Email()))
	userId := fmt.Sprintf("%x", m.Sum(nil))
	authCookieValue := objx.New(map[string]interface{}{
		"userId": userId,
		"name": user.Name(),
		"avatarUrl": user.AvatarURL(),
		"email": user.Email(),
	}).MustBase64()
	return authCookieValue
}

func createCookie(authCookieValue string) *http.Cookie {
	return &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	}
}

func providerLoginManager(res http.ResponseWriter, providerName string) {
	provider, err := gomniauth.Provider(providerName)
	if err != nil {
		msg := fmt.Sprintf("Error when trying to get provider %s, %s", provider, err)
		http.Error(res, msg, http.StatusInternalServerError)
		return
	}
	loginUrl, err := provider.GetBeginAuthURL(nil, nil)
	if err != nil {
		msg := fmt.Sprintf("Error when trying to GetBeginAuthURL for %s: %s", provider, err)
		http.Error(res, msg, http.StatusInternalServerError)
		return
	}
	res.Header().Set("Location", loginUrl)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
