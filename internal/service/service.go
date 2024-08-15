package service

import (
	"log"

	"github.com/theborzet/jwt_token/internal/repository"
	"github.com/theborzet/jwt_token/pkg/auth"
)

type Service interface {
	IssueTokens(userId, ipAddress string) (string, *auth.RefreshToken, error)
	RefreshTokens(userId, accessToken, refreshToken, ipAddress string) (string, string, error)
}

type ApiService struct {
	repo         repository.Repository
	logger       *log.Logger
	tokenManager auth.TokenManager
}

func NewApiService(repo repository.Repository, logger *log.Logger, tokenMan auth.TokenManager) *ApiService {
	return &ApiService{
		repo:         repo,
		logger:       logger,
		tokenManager: tokenMan,
	}
}
