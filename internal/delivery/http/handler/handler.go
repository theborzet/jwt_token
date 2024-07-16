package handler

import (
	"log"

	"github.com/theborzet/time-tracker/internal/pagination"
	"github.com/theborzet/time-tracker/internal/service"
)

type CommonResponse struct {
	Message   string                `json:"message"`
	Data      interface{}           `json:"data,omitempty"`
	Paginator *pagination.Paginator `json:"paginator,omitempty"`
}
type SuccessResponse struct {
	Message string `json:"message"`
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
