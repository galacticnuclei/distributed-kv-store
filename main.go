package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"kvstore/api"
	"kvstore/node"
	"kvstore/store"
)

func main() {
	port := "8001"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	peers := []string{}
	if len(os.Args) > 2 {
		peers = os.Args[2:]
	}

	kv := store.NewKVStore()

	n := &node.Node{
		ID:    port,
		Port:  port,
		Peers: peers,
		Role:  node.Follower,
		Store: kv,
	}

	n.LastHeartbeat = time.Now().Add(-10 * time.Second)
	n.LastElection = time.Now()

	handler := &api.Handler{
		Store: n.Store,
		Node:  n,
	}

	http.HandleFunc("/key/", handler.PutHandler)
	http.HandleFunc("/heartbeat", handler.HeartbeatHandler)
	http.HandleFunc("/vote", handler.VoteHandler)
	http.HandleFunc("/replicate", handler.ReplicateHandler)

	fmt.Println("Node running on", port)

	// leader heartbeat
	go func() {
		for {
			time.Sleep(2 * time.Second)
			n.Mu.Lock()
			if n.Role == node.Leader {
				for _, p := range n.Peers {
					http.Post("http://"+p+"/heartbeat", "application/json", nil)
				}
			}
			n.Mu.Unlock()
		}
	}()

	// election loop
	go func() {
		for {
			// 🔥 RANDOMIZED timeout
			timeout := time.Duration(3+rand.Intn(5)) * time.Second
			time.Sleep(timeout)

			n.Mu.Lock()

			// if already leader, do nothing.
			if n.Role == node.Leader {
				n.Mu.Unlock()
				continue
			}

			if time.Since(n.LastHeartbeat) < timeout {
				n.Mu.Unlock()
				continue
			}

			n.Term++
			term := n.Term
			n.VotedFor = n.ID
			n.Mu.Unlock()

			votes := 1

			for _, peer := range n.Peers {
				req := map[string]interface{}{
					"term":      term,
					"candidate": n.ID,
				}
				data, _ := json.Marshal(req)

				resp, err := http.Post("http://"+peer+"/vote", "application/json", bytes.NewBuffer(data))
				if err == nil && resp.StatusCode == 200 {
					votes++
				}
			}

			if votes > (len(n.Peers)+1)/2 {
				n.Mu.Lock()
				if n.Role != node.Leader {
					fmt.Println("Leader elected:", n.ID)
					n.Role = node.Leader
				}
				n.Mu.Unlock()
			}
		}
	}()

	http.ListenAndServe(":"+port, nil)
}
