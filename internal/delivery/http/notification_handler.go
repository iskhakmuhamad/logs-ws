package http

import (
	"net/http"
	"strconv"

	"github.com/iskhakmuhamad/mylogs-ws/internal/entity"
	"github.com/iskhakmuhamad/mylogs-ws/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	uc usecase.NotificationUsecase
}

func NewNotificationHandler(app *fiber.App, uc usecase.NotificationUsecase) {
	h := &NotificationHandler{uc}
	app.Post("/notifications", h.Create)
	app.Get("/notifications", h.List)
	app.Patch("/notifications/:id/read", h.MarkAsRead)
}

func (h *NotificationHandler) Create(c *fiber.Ctx) error {
	var n entity.Notification
	if err := c.BodyParser(&n); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid body")
	}
	if n.UserID == 0 || n.Title == "" {
		return fiber.NewError(http.StatusBadRequest, "user_id and title required")
	}
	if err := h.uc.Create(&n); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.Status(http.StatusCreated).JSON(n)
}

func (h *NotificationHandler) List(c *fiber.Ctx) error {
	userStr := c.Query("user_id")
	uid64, _ := strconv.ParseUint(userStr, 10, 64)
	uid := uint(uid64)
	if uid == 0 {
		return fiber.NewError(http.StatusBadRequest, "user_id required")
	}
	unreadOnly := c.Query("unread") == "1"
	list, err := h.uc.List(uid, unreadOnly)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(list)
}

func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	if err := h.uc.MarkAsRead(uint(id)); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{"id": id, "read": true})
}
