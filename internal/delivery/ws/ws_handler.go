package ws

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/iskhakmuhamad/mylogs-ws/internal/repository"
)

func NewWsHandler(app *fiber.App, hub *Hub, repo repository.NotificationRepository) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		userStr := conn.Query("user_id")
		uid64, _ := strconv.ParseUint(userStr, 10, 64)
		uid := uint(uid64)
		if uid == 0 {
			conn.Close()
			return
		}

		hub.Add(uid, conn)
		defer hub.Remove(uid, conn)

		// initial unread
		if list, err := repo.FindByUser(uid, true); err == nil {
			conn.WriteJSON(map[string]any{"type": "init", "unread": list})
		}

		conn.SetReadLimit(1 << 20)
		conn.SetReadDeadline(time.Now().Add(3 * time.Minute))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(3 * time.Minute))
			return nil
		})

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
}
