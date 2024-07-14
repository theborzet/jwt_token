package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetUserTasks обрабатывает GET запрос для получения задач пользователя.
// @Summary Получить задачи пользователя
// @Description Получает задачи пользователя с заданными параметрами пагинации и временем.
// @Accept json
// @Produce json
// @Param userId query int true "ID пользователя"
// @Param startTime query string false "Начальное время"
// @Param endTime query string false "Конечное время"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /user/tasks [get]
func (h *ApiHandler) GetUserTasks(ctx *fiber.Ctx) error {
	userId, err := strconv.Atoi(ctx.Query("userId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}
	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	tasks, err := h.serv.GetUserTasks(userId, startTime, endTime)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "server error",
		})
	}

	return ctx.JSON(fiber.Map{
		"tasks": tasks,
	})
}

// StartTask обрабатывает POST запрос для начала задачи пользователя.
// @Summary Начать задачу
// @Description Начинает задачу для указанного пользователя с заданным названием задачи.
// @Accept json
// @Produce json
// @Param userId query int true "ID пользователя"
// @Param taskName query string true "Название задачи"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /task/start [post]
func (h *ApiHandler) StartTask(ctx *fiber.Ctx) error {
	userId, err := strconv.Atoi(ctx.Query("userId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}

	taskName := ctx.Query("taskName")

	if err := h.serv.StartTask(userId, taskName); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "server error",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "start success",
	})
}

// EndTask обрабатывает POST запрос для завершения задачи пользователя.
// @Summary Завершить задачу
// @Description Завершает задачу для указанного пользователя с заданным названием задачи.
// @Accept json
// @Produce json
// @Param userId query int true "ID пользователя"
// @Param taskName query string true "Название задачи"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Router /task/end [post]
func (h *ApiHandler) EndTask(ctx *fiber.Ctx) error {
	userId, err := strconv.Atoi(ctx.Query("userId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid userId"})
	}

	taskName := ctx.Query("taskName")

	if err := h.serv.EndTask(userId, taskName); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "server error",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "end success",
	})
}
