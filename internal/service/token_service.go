package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/theborzet/jwt_token/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

func (s *ApiService) IssueTokens(userId, ipAddress string) (string, *auth.RefreshToken, error) {
	if userId == "" || ipAddress == "" {
		return "", nil, errors.New("invalid input")
	}
	tokenId := uuid.New().String()

	accessToken, err := s.tokenManager.NewJWT(userId, ipAddress, tokenId)
	if err != nil {
		return "", nil, fmt.Errorf("could not generate access token: %w", err)
	}

	refreshToken, err := s.tokenManager.NewRefreshToken(tokenId)
	if err != nil {
		return "", nil, fmt.Errorf("could not generate refresh token: %w", err)
	}

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken.Token), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, fmt.Errorf("could not hash refresh token: %w", err)
	}

	if err := s.repo.SaveRefreshTokenHash(userId, tokenId, string(hashedToken), refreshToken.ExpiresAt); err != nil {
		return "", nil, fmt.Errorf("could not save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *ApiService) RefreshTokens(userId, accessToken, refreshToken, ipAddress string) (string, string, error) {
	if userId == "" || accessToken == "" || refreshToken == "" || ipAddress == "" {
		return "", "", errors.New("invalid input")
	}

	claimsUserID, storedIP, accessTokenID, err := s.tokenManager.ParseJWT(accessToken)
	if err != nil {
		return "", "", fmt.Errorf("could not parse access token: %w", err)
	}

	if claimsUserID != userId {
		return "", "", errors.New("access token does not belong to this user")
	}

	storedRefreshTokenHash, refreshTokenID, expiresAt, err := s.repo.GetRefreshTokenHash(userId)
	if err != nil {
		return "", "", fmt.Errorf("could not get refresh token hash from database: %w", err)
	}

	if time.Now().After(expiresAt) {
		return "", "", errors.New("refresh token expired")
	}

	if accessTokenID != refreshTokenID {
		return "", "", errors.New("mismatched token ID between access and refresh tokens")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedRefreshTokenHash), []byte(refreshToken)); err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	if storedIP != ipAddress {
		//Email-адрес должен храниться в таблице user, и мы бы по guid'у доставали его отттуда
		//и уже прикрутили бы не моковую логику для отправки писем
		s.logger.Printf("Simulated email sent to user %s: Your IP address has changed from %s to %s.", userId, storedIP, ipAddress)
	}

	newAccessToken, newRefreshToken, err := s.IssueTokens(userId, ipAddress)
	if err != nil {
		return "", "", fmt.Errorf("could not issue new tokens: %w", err)
	}

	return newAccessToken, newRefreshToken.Token, nil
}
