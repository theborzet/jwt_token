package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "github.com/theborzet/jwt_token/internal/service/mocks"
	"github.com/theborzet/jwt_token/pkg/auth"
)

func TestIssueTokensHandler(t *testing.T) {
	tests := []struct {
		name         string
		userId       string
		mockService  func(ctrl *gomock.Controller) *mocks.MockService
		expectedCode int
	}{
		{
			name:   "Success",
			userId: "user-1234",
			mockService: func(ctrl *gomock.Controller) *mocks.MockService {
				s := mocks.NewMockService(ctrl)
				s.EXPECT().IssueTokens("user-1234", gomock.Any()).Return("access-token", &auth.RefreshToken{Token: "refresh-token"}, nil)
				return s
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "Service Error",
			userId: "user-1234",
			mockService: func(ctrl *gomock.Controller) *mocks.MockService {
				s := mocks.NewMockService(ctrl)
				s.EXPECT().IssueTokens("user-1234", gomock.Any()).Return("", nil, fmt.Errorf("service error"))
				return s
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := fiber.New()
			mockService := tt.mockService(ctrl)
			apiHandler := &ApiHandler{serv: mockService}
			app.Post("/api/v1/auth/token/:userID", apiHandler.IssueTokensHandler)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/auth/token/%s", tt.userId), nil)

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}

func TestRefreshTokensHandler(t *testing.T) {
	tests := []struct {
		name         string
		userId       string
		accessToken  string
		refreshToken string
		ip           string
		mockService  func(ctrl *gomock.Controller) *mocks.MockService
		expectedCode int
	}{
		{
			name:         "Success",
			userId:       "user-1234",
			accessToken:  "valid-access-token",
			refreshToken: "valid-refresh-token",
			ip:           "192.168.1.1",
			mockService: func(ctrl *gomock.Controller) *mocks.MockService {
				s := mocks.NewMockService(ctrl)
				s.EXPECT().RefreshTokens("user-1234", "valid-access-token", "valid-refresh-token", gomock.Any()).
					Return("new-access-token", "new-refresh-token", nil)
				return s
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Service Error",
			userId:       "user-1234",
			accessToken:  "valid-access-token",
			refreshToken: "valid-refresh-token",
			ip:           "192.168.1.1",
			mockService: func(ctrl *gomock.Controller) *mocks.MockService {
				s := mocks.NewMockService(ctrl)
				s.EXPECT().RefreshTokens("user-1234", "valid-access-token", "valid-refresh-token", gomock.Any()).
					Return("", "", fmt.Errorf("service error"))
				return s
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := fiber.New()
			mockService := tt.mockService(ctrl)
			apiHandler := &ApiHandler{serv: mockService}
			app.Post("/api/v1/auth/token/refresh/:userID", apiHandler.RefreshTokensHandler)

			body := fmt.Sprintf(`{"refresh_token": "%s", "access_token": "%s"}`, tt.refreshToken, tt.accessToken)
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/auth/token/refresh/%s", tt.userId), bytes.NewBuffer([]byte(body)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Forwarded-For", tt.ip)

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
