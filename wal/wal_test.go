package wal

import (
	"os"
	"testing"

	"kvstore/store"
)

func TestRecover(t *testing.T) {

	// clean old WAL if it exists
	os.Remove("wal.log")

	// write commands to WAL
	err := Append("SET name mihir")
	if err != nil {
		t.Fatalf("failed to append: %v", err)
	}

	err = Append("SET city mumbai")
	if err != nil {
		t.Fatalf("failed to append: %v", err)
	}

	err = Append("DELETE city")
	if err != nil {
		t.Fatalf("failed to append: %v", err)
	}

	// simulate restart
	kv := store.NewKVStore()

	err = Recover(kv)
	if err != nil {
		t.Fatalf("failed to recover: %v", err)
	}

	// name should exist
	val, ok := kv.Get("name")

	if !ok {
		t.Fatal("expected key 'name' to exist")
	}

	if val != "mihir" {
		t.Fatalf("expected 'mihir', got '%s'", val)
	}

	// city should have been deleted
	_, ok = kv.Get("city")

	if ok {
		t.Fatal("expected key 'city' to be deleted")
	}

	// cleanup
	os.Remove("wal.log")
}
