package main

import (
	"fmt"
	"net/http"
	"os"

	"kvstore/api"
	"kvstore/store"
)

func main() {
	port := "8001"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	kv := store.NewKVStore()
	handler := &api.Handler{Store: kv}

	http.HandleFunc("/set", handler.SetHandler)
	http.HandleFunc("/get", handler.GetHandler)
	http.HandleFunc("/delete", handler.DeleteHandler)

	fmt.Println("Server running on port", port)
	http.ListenAndServe(":"+port, nil)
}
