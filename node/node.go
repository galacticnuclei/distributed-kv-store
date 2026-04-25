package node
import "kvstore/store"

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
}