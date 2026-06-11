package store

import (
	"strconv"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	kv := NewKVStore()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		kv.Set(strconv.Itoa(i), "value")
	}
}

func BenchmarkGet(b *testing.B) {
	kv := NewKVStore()

	// preload data
	for i := 0; i < 10000; i++ {
		kv.Set(strconv.Itoa(i), "value")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		kv.Get(strconv.Itoa(i % 10000))
	}
}

func BenchmarkDelete(b *testing.B) {
	kv := NewKVStore()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)

		kv.Set(key, "value")
		kv.Delete(key)
	}
}

func BenchmarkConcurrentReadWrite(b *testing.B) {
	kv := NewKVStore()

	b.RunParallel(func(pb *testing.PB) {
		i := 0

		for pb.Next() {
			key := strconv.Itoa(i)

			kv.Set(key, "value")
			kv.Get(key)

			i++
		}
	})
}

func TestSetGetDelete(t *testing.T) {
	kv := NewKVStore()

	// SET
	kv.Set("name", "mihir")

	// GET
	val, ok := kv.Get("name")

	if !ok {
		t.Fatal("expected key 'name' to exist")
	}

	if val != "mihir" {
		t.Fatalf("expected value 'mihir', got '%s'", val)
	}

	// DELETE
	kv.Delete("name")

	_, ok = kv.Get("name")

	if ok {
		t.Fatal("expected key 'name' to be deleted")
	}
}
