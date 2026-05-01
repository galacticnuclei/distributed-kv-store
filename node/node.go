package node

import (
	"sync"
	"time"

	"kvstore/store"
)

type Role string

const (
	Leader   Role = "leader"
	Follower Role = "follower"
)

type Node struct {
	ID    string
	Port  string
	Peers []string
	Role  Role

	Store *store.KVStore

	Mu            sync.Mutex
	LastHeartbeat time.Time
	LastElection  time.Time

	Term     int
	VotedFor string
}
