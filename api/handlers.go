package api

import (
	"fmt"
	"encoding/json"
	"net/http"

	"kvstore/store"
)

type Handler struct {
	Store *store.KVStore
}

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (h *Handler) SetHandler(w http.ResponseWriter, r *http.Request) {
	var req SetRequest
	json.NewDecoder(r.Body).Decode(&req)

	h.Store.Set(req.Key, req.Value)
	w.Write([]byte("OK"))
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	val, ok := h.Store.Get(key)
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Write([]byte(val))
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	h.Store.Delete(key)
	w.Write([]byte("Deleted"))
}

func (h *Handler) HeartbeatHandler(w http.ResponseWriter,r *http.Request) {
	fmt.Println("Received heartbeat")
	w.Write([]byte("OK"))

}
