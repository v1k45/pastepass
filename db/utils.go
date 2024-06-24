package db

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

func NewDB(path string, reset bool) (*DB, error) {
	if reset {
		removeDB(path)
	}

	boltDB, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &DB{boltDB: boltDB, path: path}, nil
}

func removeDB(path string) {
	slog.Info("resetting_db", "path", path)
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		slog.Error("db_does_not_exist", "path", path, "error", err)
		return
	}

	if err := os.Remove(path); err != nil {
		slog.Error("error_removing_db", "path", path, "error", err)
		return
	}

	slog.Info("db_removed", "path", path)
}

func NewTestDB() (*DB, error) {
	testDbName := fmt.Sprintf(".test.%d.boltdb", rand.Int())
	return NewDB(testDbName, true)
}
