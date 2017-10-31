package main

import (
	"net/http"
	"strings"
	"log"
	"fmt"
	"github.com/stretchr/gomniauth"
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
	default:
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "Auth action %s not supported", action)
	}
}