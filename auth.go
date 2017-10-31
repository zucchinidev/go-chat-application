package main

import "net/http"

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
