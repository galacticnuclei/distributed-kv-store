package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"kvstore/node"
	"kvstore/store"
	"kvstore/wal"
)

type Handler struct {
	Store *store.KVStore
	Node  *node.Node
}

type SetRequest struct {
	Value string `json:"value"`
}

type ReplicateRequest struct {
	Term  int    `json:"term"`
	Op    string `json:"op"`
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

	entry := node.LogEntry{
		Term:  h.Node.Term,
		Op:    "SET",
		Key:   key,
		Value: req.Value,
	}

	h.Node.Mu.Lock()
	h.Node.Log = append(h.Node.Log, entry)
	h.Node.Mu.Unlock()

	wal.Append("SET " + key + " " + req.Value)

	h.Store.Set(key, req.Value)

	acks := replicateToFollowers(h.Node, entry)
	majority := (len(h.Node.Peers)+1)/2 + 1

	if acks < majority {
		http.Error(
			w,
			"Failed to reach majority",
			http.StatusInternalServerError,
		)
		return
	}

	w.Write([]byte("OK"))
}

// GET
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

// DELETE
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := getKeyFromPath(r.URL.Path)

	if key == "" {
		http.Error(w, "Invalid key", http.StatusBadRequest)
		return
	}

	if h.Node.Role != node.Leader {
		http.Error(w, "Not leader", http.StatusForbidden)
		return
	}

	wal.Append("DELETE " + key)

	h.Store.Delete(key)

	w.Write([]byte("Deleted"))
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

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	entry := node.LogEntry{
		Term:  req.Term,
		Op:    req.Op,
		Key:   req.Key,
		Value: req.Value,
	}

	h.Node.Mu.Lock()
	h.Node.Log = append(h.Node.Log, entry)
	h.Node.Mu.Unlock()

	switch req.Op {
	case "SET":
		h.Store.Set(req.Key, req.Value)

	case "DELETE":
		h.Store.Delete(req.Key)
	}

	w.Write([]byte("OK"))
}

func replicateToFollowers(
	n *node.Node,
	entry node.LogEntry,
) int {

	body := ReplicateRequest{
		Term:  entry.Term,
		Op:    entry.Op,
		Key:   entry.Key,
		Value: entry.Value,
	}

	data, _ := json.Marshal(body)

	acks := 1 // leader counts as an ACK

	for _, peer := range n.Peers {

		resp, err := http.Post(
			"http://"+peer+"/replicate",
			"application/json",
			bytes.NewBuffer(data),
		)

		if err != nil {
			continue
		}

		if resp.StatusCode == http.StatusOK {
			acks++
		}
	}

	return acks
}
