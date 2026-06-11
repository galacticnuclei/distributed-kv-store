# Distributed Key-Value Store in Go

A distributed key-value store built in Go featuring leader election, heartbeat-based failure detection, replicated logs, majority-acknowledged writes, write-ahead logging (WAL), and crash recovery.

## Overview

This project explores the core concepts behind distributed databases and consensus systems. It implements a multi-node key-value store where a leader node coordinates writes, replicates operations to follower nodes, and requires a majority of nodes to acknowledge a write before it is committed.

The system uses an in-memory datastore protected by Go's `sync.RWMutex` for concurrent access, while nodes communicate over HTTP to exchange heartbeats, election votes, and replication requests. A write-ahead log (WAL) provides durability and allows nodes to recover state after restarts.

## Features

* Thread-safe in-memory key-value storage
* RESTful HTTP API for CRUD operations
* Leader election using term-based voting
* Heartbeat-based failure detection
* Replicated operation log across cluster nodes
* Majority-acknowledged writes (quorum commits)
* Write-ahead logging (WAL)
* Crash recovery through WAL replay
* Multi-node cluster configuration
* Modular architecture separating storage, networking, and cluster management

## Architecture

### Store Layer

Responsible for maintaining key-value data in memory.

* Concurrent access protected using `sync.RWMutex`
* Supports Set, Get, and Delete operations
* Optimized for low-latency reads and writes

### Node Layer

Responsible for cluster coordination.

* Tracks node role (Leader/Follower)
* Maintains peer information
* Handles election state and voting
* Maintains replicated log entries
* Coordinates heartbeat monitoring

### API Layer

Exposes HTTP endpoints for both client and inter-node communication.

#### Client Endpoints

```text
PUT    /key/{key}
GET    /key/{key}
DELETE /key/{key}
```

#### Internal Cluster Endpoints

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

## Failure Handling

The cluster requires a majority of nodes to acknowledge a write before it is accepted.

| Cluster State       | Result         |
| ------------------- | -------------- |
| 3/3 nodes available | Write succeeds |
| 2/3 nodes available | Write succeeds |
| 1/3 nodes available | Write rejected |

This prevents writes from being committed when a quorum is unavailable.

## Persistence & Recovery

All write operations are appended to a write-ahead log before being applied.

On startup, nodes replay the WAL to rebuild the in-memory state and recover from crashes.

Example:

```text
SET name mihir
SET city mumbai
DELETE city
```

## Testing

### Automated Tests

* KV store CRUD unit test
* WAL recovery test

Both tests currently pass successfully.

### Manual Validation

* Leader election across three nodes
* Heartbeat-based failure detection
* Replicated log propagation
* Majority-acknowledged writes
* Crash recovery after restart

## Benchmarks

### Storage Engine Benchmarks

Measured using Go's benchmarking framework.

| Operation             | Result        |
| --------------------- | ------------- |
| Read                  | ~28M ops/sec  |
| Write                 | ~2.9M ops/sec |
| Delete                | ~20M ops/sec  |
| Concurrent Read/Write | ~5.5M ops/sec |

Hardware: AMD Ryzen 7 250 with Radeon 780M Graphics.

### HTTP API Benchmark

1000 PUT requests executed against a running node.

| Metric                | Result      |
| --------------------- | ----------- |
| Average Write Latency | 2.23 ms     |
| Throughput            | 448 req/sec |
| Requests Tested       | 1000        |

## Technologies Used

* Go
* net/http
* sync.RWMutex
* JSON
* REST APIs
* Git/GitHub

## Key Learnings

Through this project I gained hands-on experience with:

* Concurrent programming in Go
* HTTP server development
* Distributed systems fundamentals
* Leader election and failure detection
* Replicated logs and quorum-based writes
* Write-ahead logging and crash recovery
* Benchmarking and performance analysis

## Future Improvements

Potential future extensions include:

* Commit index and deferred log application
* Snapshotting and log compaction
* Stronger Raft-style consistency guarantees
* Leader redirection for client requests
* Containerized deployment using Docker
* Persistent replicated logs
