package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type request struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type successResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// @Summary Выдать токены
// @Description Выдает пару Access и Refresh токенов для пользователя с указанным ID.
// Токены возвращаются в теле ответа и должны быть сохранены вручную на клиентской стороне для последующего использования.
// @Accept json
// @Produce json
// @Param userID path string true "ID пользователя (GUID)"
// @Success 200 {object} successResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/auth/token/{userID} [post]
func (h *ApiHandler) IssueTokensHandler(ctx *fiber.Ctx) error {
	userId := ctx.Params("userID")
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

// @Summary Обновить токены
// @Description Обновляет пару Access и Refresh токенов, используя предоставленные Refresh и Access токены.
// Токены должны быть переданы в теле запроса. Обновленные токены возвращаются в теле ответа и должны быть сохранены вручную на клиентской стороне.
// Это все для упрощения, по-хорошему их нужно хранить в кэше или на крайний случай в куках
// @Accept json
// @Produce json
// @Param userID path string true "ID пользователя (GUID)"
// @Param request body request true "Тело запроса с токенами"
// @Success 200 {object} successResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/auth/token/refresh/{userID} [post]
func (h *ApiHandler) RefreshTokensHandler(ctx *fiber.Ctx) error {
	userId := ctx.Params("userID")
	var req request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Error:   err.Error(),
			Message: "invalid request",
		})
	}
	ipAddress := ctx.IP()

	accessToken, refreshToken, err := h.serv.RefreshTokens(userId, req.AccessToken, req.RefreshToken, ipAddress)
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
