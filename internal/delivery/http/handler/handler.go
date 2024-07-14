package handler

import "github.com/theborzet/time-tracker/internal/service"

type ApiHandler struct {
	serv *service.ApiService
}

func NewApiHandler(serv *service.ApiService) *ApiHandler {
	return &ApiHandler{serv: serv}
}
