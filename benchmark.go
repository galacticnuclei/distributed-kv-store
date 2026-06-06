package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{}

	const requests = 1000

	start := time.Now()

	for i := 0; i < requests; i++ {
		body := []byte(`{"value":"test"}`)

		req, err := http.NewRequest(
			http.MethodPut,
			"http://localhost:8001/key/benchmark",
			bytes.NewBuffer(body),
		)

		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		resp.Body.Close()
	}

	elapsed := time.Since(start)

	fmt.Println("Requests:", requests)
	fmt.Println("Total time:", elapsed)

	fmt.Printf(
		"Average latency: %.3f ms\n",
		float64(elapsed.Microseconds())/1000.0/requests,
	)

	fmt.Printf(
		"Throughput: %.2f req/sec\n",
		float64(requests)/elapsed.Seconds(),
	)
}