package main

import (
	"net/http"
	"strings"
	"log"
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
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
		provider, err := gomniauth.Provider(provider)
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
		log.Println("Todo hadler login for", provider)
	case "callback":
		provider, err := gomniauth.Provider(provider)
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

		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()

		http.SetCookie(res, &http.Cookie{
			Name: "auth",
			Value: authCookieValue,
			Path: "/",
		})
		res.Header().Set("Location", "/chat")
		res.WriteHeader(http.StatusTemporaryRedirect)
	default:
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "Auth action %s not supported", action)
	}
}