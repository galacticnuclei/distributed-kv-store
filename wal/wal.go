package wal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"kvstore/store"
)

func Append(command string) error {
	f, err := os.OpenFile(
		"wal.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = fmt.Fprintln(f, command)
	return err
}

func Recover(kv *store.KVStore) error {
	f, err := os.Open("wal.log")

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, " ", 3)

		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "SET":
			if len(parts) == 3 {
				kv.Set(parts[1], parts[2])
			}

		case "DELETE":
			kv.Delete(parts[1])
		}
	}

	return scanner.Err()
}