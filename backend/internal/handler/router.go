package handler

import (
	"net/http"
)

type UserHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

func RegisterRoutes(mux *http.ServeMux, userHandler UserHandler) {
	mux.HandleFunc("POST /api/register", userHandler.Register)
	mux.HandleFunc("GET /api/login", userHandler.Login)
}
