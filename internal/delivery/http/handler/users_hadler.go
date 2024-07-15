package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/theborzet/time-tracker/internal/models"
)

type CreateUser struct {
	PassportNumber string `json:"passportNumber"`
	PassportSerie  string `json:"passportSerie"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

// GetUsers обрабатывает GET запрос для получения списка пользователей с учетом фильтров и пагинации.
// @Summary Получить список пользователей
// @Description Получает список пользователей с учетом заданных фильтров и страницы.
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы"
// @Param filters query string false "Фильтры"
// @Success 200 {object} CommonResponse{data=[]models.User}
// @Failure 400 {object} ErrorResponse
// @Router /user/ [get]
func (h *ApiHandler) GetUsers(ctx *fiber.Ctx) error {
	filters := make(map[string]string)

	for key, values := range ctx.Queries() {
		if len(values) > 0 {
			filters[key] = values
		}
	}

	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "invalid page number"})
	}
	users, paginator, err := h.serv.GetUsersWithPaginate(filters, page)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Server error",
		})
	}

	return ctx.JSON(CommonResponse{
		Data:      users,
		Paginator: paginator,
		Message:   "Request processed",
	})
}

// CreateUser обрабатывает POST запрос для создания нового пользователя.
// @Summary Создать пользователя
// @Description Создает нового пользователя на основе переданных данных.
// @Accept json
// @Produce json
// @Param user body CreateUser true "Данные пользователя"
// @Success 200 {object} CommonResponse
// @Failure 400 {object} ErrorResponse
// @Router /user/create [post]
func (h *ApiHandler) CreateUser(ctx *fiber.Ctx) error {
	var user *models.User

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Invalid user data"})
	}

	if err := h.serv.CreateUser(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Failed to create user"})
	}

	return ctx.JSON(CommonResponse{
		Message: "User created successfully",
	})

}

// UpdateUser обрабатывает PUT запрос для обновления данных пользователя.
// @Summary Обновить пользователя
// @Description Обновляет данные пользователя на основе переданных данных.
// @Accept json
// @Produce json
// @Param user body models.User true "Данные пользователя"
// @Success 200 {object} CommonResponse
// @Failure 400 {object} ErrorResponse
// @Router /user/update [put]
func (h *ApiHandler) UpdateUser(ctx *fiber.Ctx) error {
	var user *models.User

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Invalid user data"})
	}

	if err := h.serv.UpdateUser(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Failed to update user"})
	}

	return ctx.JSON(CommonResponse{
		Message: "user update successfully",
	})

}

// DeleteUser обрабатывает DELETE запрос для удаления пользователя по его идентификатору.
// @Summary Удалить пользователя
// @Description Удаляет пользователя на основе его идентификатора.
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} CommonResponse
// @Failure 400 {object} ErrorResponse
// @Router /user/{id} [delete]
func (h *ApiHandler) DeleteUser(ctx *fiber.Ctx) error {
	usrIdStr, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Invalid user ID"})
	}
	if err := h.serv.DeleteUser(usrIdStr); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   err.Error(),
			Message: "Server error with delete user"})
	}

	return ctx.JSON(CommonResponse{
		Message: "user delete successfully",
	})

}
