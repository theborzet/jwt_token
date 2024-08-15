package handler

import (
	"log"

	"github.com/theborzet/jwt_token/internal/service"
)

type ApiHandler struct {
	serv   service.Service
	logger *log.Logger
}

func NewApiHandler(serv service.Service, logger *log.Logger) *ApiHandler {
	return &ApiHandler{
		serv:   serv,
		logger: logger}
}
