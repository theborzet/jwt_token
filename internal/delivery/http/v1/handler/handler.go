package handler

import (
	"log"

	"github.com/theborzet/jwt_token/internal/service"
)

type SuccessResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ApiHandler struct {
	serv   *service.ApiService
	logger *log.Logger
}

func NewApiHandler(serv *service.ApiService, logger *log.Logger) *ApiHandler {
	return &ApiHandler{
		serv:   serv,
		logger: logger}
}
