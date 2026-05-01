package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"kvstore/api"
	"kvstore/store"
	"kvstore/node"
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

	n.LastHeartbeat = time.Now()

	handler := &api.Handler{
		Store: n.Store,
		Node:  n,
	}

	fmt.Printf("Node %s started on port %s\n", n.ID, n.Port)
	fmt.Println("Peers:", n.Peers)
	fmt.Println("Role:", n.Role)

	handler = &api.Handler{Store: n.Store}

	http.HandleFunc("/key/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case http.MethodPut:
				handler.PutHandler(w, r)
			case http.MethodGet:
				handler.GetHandler(w, r)
			case http.MethodDelete:
				handler.DeleteHandler(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/heartbeat", handler.HeartbeatHandler)
	fmt.Println("Server running on port", port)

	go func() {
		for {
			time.Sleep(2 * time.Second)

			for _,peer := range n.Peers {
				url := "http://" + peer + "/heartbeat"
				http.Post(url, "application/json", nil)
			}
		}
	}()
	http.ListenAndServe(":"+port, nil)
}
