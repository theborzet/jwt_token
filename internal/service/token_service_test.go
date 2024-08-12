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
	defer mockCtl.Finish()

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

	tokenManager.EXPECT().NewJWT(userId, ipAddress).Return("", errInternalServErr)

	accessToken, refreshToken, err := apiService.IssueTokens(userId, ipAddress)

	require.True(t, errors.Is(err, errInternalServErr))
	require.Empty(t, accessToken)
	require.Empty(t, refreshToken)
}

func TestIssueTokens_ErrorGeneratingRefreshToken(t *testing.T) {
	apiService, _, tokenManager := mockApiService(t)

	userId := "1"
	ipAddress := "192.168.1.1"
	refreshToken := "refresh-token"

	tokenManager.EXPECT().NewJWT(userId, ipAddress).Return("access-token", nil)
	tokenManager.EXPECT().NewRefreshToken().Return(&auth.RefreshToken{Token: refreshToken, ExpiresAt: time.Time{}}, errInternalServErr)

	accessToken, generatedRefreshToken, err := apiService.IssueTokens(userId, ipAddress)

	require.True(t, errors.Is(err, errInternalServErr))
	require.Empty(t, accessToken)
	require.Empty(t, generatedRefreshToken)
}

func TestIssueTokens_Success(t *testing.T) {
	apiService, repo, tokenManager := mockApiService(t)

	userId := "1"
	ipAddress := "192.168.1.1"
	refreshToken := "refresh-token"

	tokenManager.EXPECT().NewJWT(userId, ipAddress).Return("access-token", nil)
	tokenManager.EXPECT().NewRefreshToken().Return(&auth.RefreshToken{Token: refreshToken, ExpiresAt: time.Time{}}, nil)

	repo.EXPECT().SaveRefreshTokenHash(userId, gomock.Any(), time.Time{}).Return(nil)

	accessToken, generatedRefreshToken, err := apiService.IssueTokens(userId, ipAddress)

	require.NoError(t, err)
	require.Equal(t, "access-token", accessToken)
	require.Equal(t, refreshToken, generatedRefreshToken.Token)
}

func TestRefreshTokens_ErrorGettingStoredHash(t *testing.T) {
	apiService, repo, _ := mockApiService(t)

	userId := "1"
	ipAddress := "192.168.1.1"

	repo.EXPECT().GetRefreshTokenHash(userId).Return("", time.Time{}, errInternalServErr)

	accessToken, refreshToken, err := apiService.RefreshTokens(userId, "refresh-token", ipAddress)

	require.True(t, errors.Is(err, errInternalServErr))
	require.Empty(t, accessToken)
	require.Empty(t, refreshToken)
}

func TestRefreshTokens_ErrorComparingHash(t *testing.T) {
	apiService, repo, _ := mockApiService(t)

	userId := "1"
	ipAddress := "192.168.1.1"
	refreshToken := "refresh-token"
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	expiresAt := time.Now().Add(24 * time.Hour)
	repo.EXPECT().GetRefreshTokenHash(userId).Return(string(hashedToken), expiresAt, nil)

	accessToken, refreshToken, err := apiService.RefreshTokens(userId, "wrong-refresh-token", ipAddress)

	require.EqualError(t, err, "invalid refresh token")
	require.Empty(t, accessToken)
	require.Empty(t, refreshToken)
}

func TestRefreshTokens_IPChanged(t *testing.T) {
	apiService, repo, tokenManager := mockApiService(t)

	userId := "1"
	oldIP := "192.168.1.1"
	newIP := "192.168.1.2"
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte("refresh-token"), bcrypt.DefaultCost)
	expiresAt := time.Now().Add(24 * time.Hour)

	repo.EXPECT().GetRefreshTokenHash(userId).Return(string(hashedToken), expiresAt, nil)
	tokenManager.EXPECT().ParseJWT(gomock.Any()).Return(userId, oldIP, nil)
	tokenManager.EXPECT().NewJWT(userId, newIP).Return("new-access-token", nil)

	tokenManager.EXPECT().NewRefreshToken().Return(&auth.RefreshToken{Token: "new-refresh-token", ExpiresAt: expiresAt}, nil)
	repo.EXPECT().SaveRefreshTokenHash(userId, gomock.Any(), expiresAt).Return(nil)

	accessToken, generatedRefreshToken, err := apiService.RefreshTokens(userId, "refresh-token", newIP)

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
	rt := &auth.RefreshToken{
		Token:     "new-refresh-token",
		ExpiresAt: expiresAt,
	}

	repo.EXPECT().GetRefreshTokenHash(userId).Return(string(hashedToken), expiresAt, nil)
	tokenManager.EXPECT().ParseJWT(gomock.Any()).Return(userId, ipAddress, nil)
	tokenManager.EXPECT().NewJWT(userId, ipAddress).Return("new-access-token", nil)

	tokenManager.EXPECT().NewRefreshToken().Return(rt, nil)
	repo.EXPECT().SaveRefreshTokenHash(userId, gomock.Any(), expiresAt).Return(nil)

	accessToken, generatedRefreshToken, err := apiService.RefreshTokens(userId, "refresh-token", ipAddress)

	require.NoError(t, err)
	require.Equal(t, "new-access-token", accessToken)
	require.Equal(t, rt.Token, generatedRefreshToken)
}

func TestRefreshTokens_ExpiredRefreshToken(t *testing.T) {
	apiService, repo, _ := mockApiService(t)

	userId := "1"
	ipAddress := "192.168.1.1"
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte("refresh-token"), bcrypt.DefaultCost)
	expiredAt := time.Now().Add(-1 * time.Hour)

	repo.EXPECT().GetRefreshTokenHash(userId).Return(string(hashedToken), expiredAt, nil)

	accessToken, refreshToken, err := apiService.RefreshTokens(userId, "refresh-token", ipAddress)

	require.EqualError(t, err, "refresh token expired")
	require.Empty(t, accessToken)
	require.Empty(t, refreshToken)
}
