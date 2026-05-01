package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"kvstore/node"
	"kvstore/store"
)

type Handler struct {
	Store *store.KVStore
	Node  *node.Node
}

type SetRequest struct {
	Value string `json:"value"`
}

type ReplicateRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type VoteRequest struct {
	Term      int    `json:"term"`
	Candidate string `json:"candidate"`
}

func getKeyFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}

// PUT
func (h *Handler) PutHandler(w http.ResponseWriter, r *http.Request) {
	key := getKeyFromPath(r.URL.Path)

	if key == "" {
		http.Error(w, "Invalid key", http.StatusBadRequest)
		return
	}

	if h.Node.Role != node.Leader {
		http.Error(w, "Not leader", http.StatusForbidden)
		return
	}

	var req SetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Value == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	h.Store.Set(key, req.Value)
	replicateToFollowers(h.Node, key, req.Value)

	w.Write([]byte("OK"))
}

// heartbeat
func (h *Handler) HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	h.Node.Mu.Lock()
	h.Node.LastHeartbeat = time.Now()
	h.Node.Role = node.Follower
	h.Node.Mu.Unlock()

	w.Write([]byte("OK"))
}

// fixed voting logic
func (h *Handler) VoteHandler(w http.ResponseWriter, r *http.Request) {
	var req VoteRequest
	json.NewDecoder(r.Body).Decode(&req)

	h.Node.Mu.Lock()
	defer h.Node.Mu.Unlock()

	// if incoming term is newer → accept it
	if req.Term > h.Node.Term {
		h.Node.Term = req.Term
		h.Node.VotedFor = ""
		h.Node.Role = node.Follower
	}

	// reject old terms
	if req.Term < h.Node.Term {
		http.Error(w, "Old term", http.StatusForbidden)
		return
	}

	// vote only once per term
	if h.Node.VotedFor == "" || h.Node.VotedFor == req.Candidate {
		h.Node.VotedFor = req.Candidate
		w.Write([]byte("OK"))
		return
	}

	http.Error(w, "Already voted", http.StatusForbidden)
}

func (h *Handler) ReplicateHandler(w http.ResponseWriter, r *http.Request) {
	var req ReplicateRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.Key != "" {
		h.Store.Set(req.Key, req.Value)
	}

	w.Write([]byte("OK"))
}

func replicateToFollowers(n *node.Node, key, value string) {
	body := map[string]string{"key": key, "value": value}
	data, _ := json.Marshal(body)

	for _, peer := range n.Peers {
		http.Post("http://"+peer+"/replicate", "application/json", bytes.NewBuffer(data))
	}
}
