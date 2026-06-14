package wsclients

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Registry is a thread-safe map of WebSocket connections keyed by user ID.
type Registry struct {
	mu      sync.RWMutex
	clients map[string]*websocket.Conn
}

func NewRegistry() *Registry {
	return &Registry{clients: make(map[string]*websocket.Conn)}
}

func (r *Registry) Set(id string, conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[id] = conn
}

func (r *Registry) Get(id string) (*websocket.Conn, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	conn, ok := r.clients[id]
	return conn, ok
}

func (r *Registry) Delete(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, id)
}
