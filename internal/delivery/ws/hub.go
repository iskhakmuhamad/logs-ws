package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Hub struct {
	mu        sync.RWMutex
	conns     map[uint]map[*websocket.Conn]struct{}
	broadcast chan broadcastMessage
}

type broadcastMessage struct {
	userID  uint
	payload any
}

func NewHub() *Hub {
	return &Hub{
		conns:     make(map[uint]map[*websocket.Conn]struct{}),
		broadcast: make(chan broadcastMessage),
	}
}

func (h *Hub) Run() {
	for msg := range h.broadcast {
		h.send(msg.userID, msg.payload)
	}
}

func (h *Hub) Add(userID uint, c *websocket.Conn) {
	h.mu.Lock()
	if h.conns[userID] == nil {
		h.conns[userID] = make(map[*websocket.Conn]struct{})
	}
	h.conns[userID][c] = struct{}{}
	h.mu.Unlock()
}

func (h *Hub) Remove(userID uint, c *websocket.Conn) {
	h.mu.Lock()
	delete(h.conns[userID], c)
	if len(h.conns[userID]) == 0 {
		delete(h.conns, userID)
	}
	h.mu.Unlock()
}

func (h *Hub) Broadcast(userID uint, payload any) {
	h.broadcast <- broadcastMessage{userID: userID, payload: payload}
}

func (h *Hub) send(userID uint, payload any) {
	msg, err := json.Marshal(payload)
	if err != nil {
		log.Println("broadcast marshal:", err)
		return
	}
	h.mu.RLock()
	conns := h.conns[userID]
	h.mu.RUnlock()
	for c := range conns {
		if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("ws write err:", err)
			c.Close()
			h.Remove(userID, c)
		}
	}
}
