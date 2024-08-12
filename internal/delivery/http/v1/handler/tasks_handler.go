package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type request struct {
	RefreshToken string `json:"refresh_token"`
}

type successResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// IssueTokensHandler обрабатывает POST запрос для выдачи Access и Refresh токенов.
// @Summary Выдать токены
// @Description Выдает пару Access и Refresh токенов для пользователя с указанным ID.
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя (GUID)"
// @Success 200 {object} successResponse{access_token=string, refresh_token=string}
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /issue-tokens/{id} [post]
func (h *ApiHandler) IssueTokensHandler(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	ip := ctx.IP()

	accessToken, refreshToken, err := h.serv.IssueTokens(userId, ip)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(errorResponse{
			Error:   err.Error(),
			Message: "Could not issue tokens"})
	}

	return ctx.JSON(successResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	})

}

// RefreshTokensHandler обрабатывает POST запрос для обновления Access и Refresh токенов.
// @Summary Обновить токены
// @Description Обновляет пару Access и Refresh токенов, используя предоставленный Refresh токен.
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя (GUID)"
// @Param request body request true "Refresh токен"
// @Success 200 {object} successResponse{access_token=string, refresh_token=string}
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /refresh-tokens/{id} [post]
func (h *ApiHandler) RefreshTokensHandler(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	var req request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Error:   err.Error(),
			Message: "invalid request",
		})
	}
	ipAddress := ctx.IP()

	accessToken, refreshToken, err := h.serv.RefreshTokens(userId, req.RefreshToken, ipAddress)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse{
			Error:   err.Error(),
			Message: "invalid request",
		})
	}

	return ctx.JSON(successResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}
