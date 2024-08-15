package service

import (
	"bytes"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_repository "github.com/theborzet/jwt_token/internal/repository/mocks"
	"github.com/theborzet/jwt_token/pkg/auth"
	mock_auth "github.com/theborzet/jwt_token/pkg/auth/mocks"
	"golang.org/x/crypto/bcrypt"
)

var errInternalServErr = errors.New("test: internal server error")

func mockApiService(t *testing.T) (*ApiService, *mock_repository.MockRepository, *mock_auth.MockTokenManager) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	repo := mock_repository.NewMockRepository(mockCtl)
	tokenManager := mock_auth.NewMockTokenManager(mockCtl)

	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "TEST: ", log.LstdFlags)

	apiService := NewApiService(repo, logger, tokenManager)

	return apiService, repo, tokenManager
}

func TestIssueTokens_ErrorGeneratingAccessToken(t *testing.T) {
	apiService, _, tokenManager := mockApiService(t)

	userId := "1"
	ipAddress := "192.168.1.1"
	var capturedTokenID string

	tokenManager.EXPECT().NewJWT(userId, ipAddress, gomock.Any()).DoAndReturn(
		func(uid, ip, tid string) (string, error) {
			capturedTokenID = tid
			return "", errInternalServErr
		},
	)

	accessToken, refreshToken, err := apiService.IssueTokens(userId, ipAddress)

	require.True(t, errors.Is(err, errInternalServErr))
	require.Empty(t, accessToken)
	require.Empty(t, refreshToken)
	require.NotEmpty(t, capturedTokenID)
}

func TestIssueTokens_ErrorGeneratingRefreshToken(t *testing.T) {
	apiService, _, tokenManager := mockApiService(t)

	userId := "1"
	ipAddress := "192.168.1.1"
	var capturedTokenID string
	refreshToken := "refresh-token"

	tokenManager.EXPECT().NewJWT(userId, ipAddress, gomock.Any()).DoAndReturn(
		func(uid, ip, tid string) (string, error) {
			capturedTokenID = tid
			return "access-token", nil
		},
	)

	tokenManager.EXPECT().NewRefreshToken(gomock.Any()).Return(&auth.RefreshToken{Token: refreshToken, ExpiresAt: time.Time{}}, errInternalServErr)

	accessToken, generatedRefreshToken, err := apiService.IssueTokens(userId, ipAddress)

	require.True(t, errors.Is(err, errInternalServErr))
	require.Empty(t, accessToken)
	require.Empty(t, generatedRefreshToken)
	require.NotEmpty(t, capturedTokenID)
}

func TestIssueTokens_Success(t *testing.T) {
	apiService, repo, tokenManager := mockApiService(t)

	userId := "1"
	ipAddress := "192.168.1.1"
	refreshToken := "refresh-token"

	tokenManager.EXPECT().NewJWT(userId, ipAddress, gomock.Any()).DoAndReturn(func(uid, ip, tokenID string) (string, error) {
		require.NotEmpty(t, tokenID, "tokenID should not be empty")
		return "access-token", nil
	})

	tokenManager.EXPECT().NewRefreshToken(gomock.Any()).DoAndReturn(func(tokenID string) (*auth.RefreshToken, error) {
		require.NotEmpty(t, tokenID, "tokenID should not be empty")
		return &auth.RefreshToken{Token: refreshToken, ExpiresAt: time.Time{}}, nil
	})

	repo.EXPECT().SaveRefreshTokenHash(userId, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	accessToken, generatedRefreshToken, err := apiService.IssueTokens(userId, ipAddress)

	require.NoError(t, err)
	require.Equal(t, "access-token", accessToken)
	require.Equal(t, refreshToken, generatedRefreshToken.Token)
}

func TestRefreshTokens_ErrorGettingStoredHash(t *testing.T) {
	apiService, repo, tokenManager := mockApiService(t)

	userId := "1"
	ipAddress := "192.168.1.1"
	refreshToken := "refresh-token"

	tokenManager.EXPECT().ParseJWT("access-token").Return(userId, ipAddress, "token-id", nil)

	repo.EXPECT().GetRefreshTokenHash(userId).Return("", "", time.Time{}, errors.New("some error"))

	_, _, err := apiService.RefreshTokens(userId, "access-token", refreshToken, ipAddress)

	require.Error(t, err)
	require.EqualError(t, err, "could not get refresh token hash from database: some error")
}

func TestRefreshTokens_ErrorComparingHash(t *testing.T) {
	apiService, repo, tokenManager := mockApiService(t)

	userId := "1"
	ipAddress := "192.198.1.1"
	refreshToken := "refresh_token"
	capturedTokenID := "token-id"
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	expiresAt := time.Now().Add(24 * time.Hour)

	repo.EXPECT().GetRefreshTokenHash(userId).Return(capturedTokenID, string(hashedToken), expiresAt, nil)
	tokenManager.EXPECT().ParseJWT("access-token").Return(userId, ipAddress, capturedTokenID, nil)

	accessToken, generatedRefreshToken, err := apiService.RefreshTokens(userId, "access-token", "wrong-refresh-token", ipAddress)

	require.EqualError(t, err, "mismatched token ID between access and refresh tokens")
	require.Empty(t, accessToken)
	require.Empty(t, generatedRefreshToken)
}

func TestRefreshTokens_IPChanged(t *testing.T) {
	apiService, repo, tokenManager := mockApiService(t)

	userId := "1"
	oldIP := "192.198.1.1"
	newIP := "192.198.1.2"
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte("refresh-token"), bcrypt.DefaultCost)
	expiresAt := time.Now().Add(24 * time.Hour)

	repo.EXPECT().GetRefreshTokenHash(userId).Return(string(hashedToken), "token-id", expiresAt, nil)

	tokenManager.EXPECT().ParseJWT(gomock.Any()).Return(userId, oldIP, "token-id", nil)
	tokenManager.EXPECT().NewJWT(userId, newIP, gomock.Any()).Return("new-access-token", nil)
	tokenManager.EXPECT().NewRefreshToken(gomock.Any()).Return(&auth.RefreshToken{Token: "new-refresh-token", ExpiresAt: expiresAt, TokenID: "token-id"}, nil)

	repo.EXPECT().SaveRefreshTokenHash(userId, gomock.Any(), gomock.Any(), expiresAt).Return(nil)

	accessToken, generatedRefreshToken, err := apiService.RefreshTokens(userId, "access-token", "refresh-token", newIP)

	require.NoError(t, err)
	require.Equal(t, "new-access-token", accessToken)
	require.Equal(t, "new-refresh-token", generatedRefreshToken)
}

func TestRefreshTokens_Success(t *testing.T) {
	apiService, repo, tokenManager := mockApiService(t)

	userId := "1"
	ipAddress := "192.168.1.1"
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte("refresh-token"), bcrypt.DefaultCost)
	expiresAt := time.Now().Add(24 * time.Hour)
	tokenID := "token-id"
	rt := &auth.RefreshToken{
		Token:     "new-refresh-token",
		ExpiresAt: expiresAt,
		TokenID:   tokenID,
	}

	repo.EXPECT().GetRefreshTokenHash(userId).Return(string(hashedToken), tokenID, expiresAt, nil)
	tokenManager.EXPECT().ParseJWT(gomock.Any()).Return(userId, ipAddress, tokenID, nil)

	tokenManager.EXPECT().NewJWT(userId, ipAddress, gomock.Any()).Return("new-access-token", nil)

	tokenManager.EXPECT().NewRefreshToken(gomock.Any()).Return(rt, nil)

	repo.EXPECT().SaveRefreshTokenHash(userId, gomock.Any(), gomock.Any(), expiresAt).Return(nil)

	accessToken, generatedRefreshToken, err := apiService.RefreshTokens(userId, "access-token", "refresh-token", ipAddress)

	require.NoError(t, err)
	require.Equal(t, "new-access-token", accessToken)
	require.Equal(t, rt.Token, generatedRefreshToken)
}

func TestRefreshTokens_ExpiredRefreshToken(t *testing.T) {
	apiService, repo, tokenManager := mockApiService(t)

	userId := "1"
	ipAddress := "192.168.1.1"
	tokenID := "token-id"
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte("refresh-token"), bcrypt.DefaultCost)
	expiredAt := time.Now().Add(-1 * time.Hour)

	repo.EXPECT().GetRefreshTokenHash(userId).Return(string(hashedToken), tokenID, expiredAt, nil)

	tokenManager.EXPECT().ParseJWT(gomock.Any()).Return(userId, ipAddress, tokenID, nil)

	accessToken, refreshToken, err := apiService.RefreshTokens(userId, "access-token", "refresh-token", ipAddress)

	require.EqualError(t, err, "refresh token expired")
	require.Empty(t, accessToken)
	require.Empty(t, refreshToken)
}
