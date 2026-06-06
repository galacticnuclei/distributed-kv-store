# Distributed Key-Value Store in Go

A distributed key-value store built in Go featuring leader election, heartbeat-based failure detection, write replication, and a thread-safe in-memory storage engine.

## Overview

This project began as an exploration of distributed systems concepts and evolved into a multi-node key-value store capable of coordinating writes through a leader node and replicating data across followers.

The system uses an in-memory datastore protected by Go's `sync.RWMutex` for concurrent access, while nodes communicate over HTTP to exchange heartbeats, election votes, and replication requests.

## Features

* Thread-safe in-memory key-value storage
* RESTful HTTP API for key-value operations
* Leader election using a simplified term-based voting mechanism
* Heartbeat-based failure detection
* Write replication from leader to follower nodes
* Multi-node cluster configuration through peer discovery
* Modular architecture separating storage, networking, and node management

## Architecture

The project is organized into three primary components:

### Store Layer

Responsible for maintaining key-value data in memory.

* Concurrent access protected using `sync.RWMutex`
* Supports Set, Get, and Delete operations
* Optimized for low-latency reads and writes

### Node Layer

Responsible for cluster coordination.

* Tracks node role (Leader/Follower)
* Maintains peer information
* Handles election state, voting, and heartbeat tracking

### API Layer

Exposes HTTP endpoints for both client and inter-node communication.

Client-facing endpoints:

```text
PUT    /key/{key}
GET    /key/{key}
DELETE /key/{key}
```

Internal cluster endpoints:

```text
POST /heartbeat
POST /vote
POST /replicate
```

## Running the Project

### Single Node

```bash
go run main.go 8001
```

### Three-Node Cluster

Terminal 1:

```bash
go run main.go 8001 localhost:8002 localhost:8003
```

Terminal 2:

```bash
go run main.go 8002 localhost:8001 localhost:8003
```

Terminal 3:

```bash
go run main.go 8003 localhost:8001 localhost:8002
```

After startup, one node will be elected leader and begin sending heartbeats to follower nodes.

## Example Usage

Store a value:

```bash
curl -X PUT http://localhost:8001/key/name \
-H "Content-Type: application/json" \
-d "{\"value\":\"mihir\"}"
```

Retrieve a value:

```bash
curl http://localhost:8001/key/name
```

Delete a value:

```bash
curl -X DELETE http://localhost:8001/key/name
```

## Benchmarks

### Storage Engine Benchmarks

Measured using Go's benchmarking framework.

| Operation             | Performance          |
| --------------------- | -------------------- |
| Read                  | ~28 million ops/sec  |
| Write                 | ~2.9 million ops/sec |
| Delete                | ~20 million ops/sec  |
| Concurrent Read/Write | ~5.5 million ops/sec |

### HTTP API Benchmark

1000 PUT requests executed against a running node.

| Metric                | Result           |
| --------------------- | ---------------- |
| Average Write Latency | 2.23 ms          |
| Throughput            | 448 requests/sec |
| Requests Tested       | 1000             |

## Future Improvements

This project intentionally focuses on the core mechanics of a distributed key-value store. Possible future extensions include:

* Persistent storage using a write-ahead log
* Snapshotting and recovery
* Stronger consensus guarantees
* Leader redirection for client requests
* Majority acknowledgement before confirming writes
* Containerized deployment using Docker

## Technologies Used

* Go
* net/http
* sync.RWMutex
* JSON
* REST APIs
* Git/GitHub

## Key Learnings

Through this project I gained hands-on experience with:

* Concurrency control in Go
* HTTP server development
* Distributed systems fundamentals
* Leader election and failure detection
* Data replication strategies
* Benchmarking and performance measurement
