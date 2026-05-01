package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"kvstore/store"
	"kvstore/node"
)

type Handler struct {
	Store *store.KVStore
	Node *node.Node
}

type SetRequest struct {
	Value string `json:"value"`
}

// helper to extract key from /key/{key}
func getKeyFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}

// PUT /key/{key}
func (h *Handler) PutHandler(w http.ResponseWriter, r *http.Request) {
	key := getKeyFromPath(r.URL.Path)

	// ✅ key validation
	if key == "" {
		http.Error(w, "Invalid key", http.StatusBadRequest)
		return
	}

	var req SetRequest
	json.NewDecoder(r.Body).Decode(&req)

	// ✅ PUT VALUE VALIDATION — ADD IT RIGHT HERE
	if req.Value == "" {
		http.Error(w, "Value required", http.StatusBadRequest)
		return
	}

	h.Store.Set(key, req.Value)
	w.Write([]byte("OK"))
}

// GET /key/{key}
func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := getKeyFromPath(r.URL.Path)
	if key == "" {
		http.Error(w, "Invalid key", http.StatusBadRequest)
		return
	}
	val, ok := h.Store.Get(key)
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Write([]byte(val))
}

// DELETE /key/{key}
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := getKeyFromPath(r.URL.Path)
	if key == "" {
		http.Error(w, "Invalid key", http.StatusBadRequest)
		return
	}
	h.Store.Delete(key)
	w.Write([]byte("Deleted"))
}

func (h *Handler) HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	// If you're using Node, update heartbeat time
	if h.Node != nil {
		h.Node.Mu.Lock()
		h.Node.LastHeartbeat = time.Now()
		h.Node.Role = node.Follower
		h.Node.Mu.Unlock()
	}

	w.Write([]byte("OK"))
}