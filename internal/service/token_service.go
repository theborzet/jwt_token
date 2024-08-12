package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/theborzet/jwt_token/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

func (s *ApiService) IssueTokens(userId, ipAddress string) (string, *auth.RefreshToken, error) {
	if userId == "" || ipAddress == "" {
		return "", nil, errors.New("invalid input")
	}
	accessToken, err := s.tokenManager.NewJWT(userId, ipAddress)
	if err != nil {
		return "", nil, fmt.Errorf("could not generate access token: %w", err)
	}

	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return "", nil, fmt.Errorf("could not generate refresh token: %w", err)
	}

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken.Token), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, fmt.Errorf("could not hash refresh token: %w", err)
	}

	if err := s.repo.SaveRefreshTokenHash(userId, string(hashedToken), refreshToken.ExpiresAt); err != nil {
		return "", nil, fmt.Errorf("could not save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *ApiService) RefreshTokens(userId, refreshToken, ipAddress string) (string, string, error) {
	if userId == "" || ipAddress == "" {
		return "", "", errors.New("invalid input")
	}

	storedRefreshTokenHash, expiresAt, err := s.repo.GetRefreshTokenHash(userId)
	if err != nil {
		return "", "", fmt.Errorf("could not get refresh token hash from database: %w", err)
	}

	if time.Now().After(expiresAt) {
		return "", "", errors.New("refresh token expired")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedRefreshTokenHash), []byte(refreshToken)); err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	_, storedIP, err := s.tokenManager.ParseJWT(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("could not parse JWT: %w", err)
	}

	if storedIP != ipAddress {
		s.logger.Printf("Simulated email sent to user %s: Your IP address has changed from %s to %s.", userId, storedIP, ipAddress)
	}

	accessToken, newRefreshToken, err := s.IssueTokens(userId, ipAddress)
	if err != nil {
		return "", "", fmt.Errorf("could not issue new tokens: %w", err)
	}

	return accessToken, newRefreshToken.Token, nil
}
