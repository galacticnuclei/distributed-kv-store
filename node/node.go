package node
import "kvstore/store"
import "sync"
import "time"

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
    Mu  sync.Mutex
    LastHeartbeat time.Time
}