package db

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func newTestDB() (*DB, error) {
	testDbName := fmt.Sprintf(".test.%d.boltdb", rand.Int())
	return NewDB(testDbName, true)
}

func TestNewPaste(t *testing.T) {
	db, err := newTestDB()
	assert.NoError(t, err)
	defer db.reset()

	paste, err := db.NewPaste("test paste", time.Now().Add(time.Hour))
	assert.NoError(t, err)

	assert.NotNil(t, paste)
	assert.NotEmpty(t, paste.ID)
	assert.NotEmpty(t, paste.EncryptedBytes)

	_, err = db.Get(paste.ID)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	db, err := newTestDB()
	assert.NoError(t, err)
	defer db.reset()

	expirationTime := time.Now().Add(time.Hour)
	paste, err := db.NewPaste("test paste", expirationTime)
	assert.NoError(t, err)

	t.Run("returns correct attributes", func(t *testing.T) {
		// only metadata is returned for Get
		savedPaste, err := db.Get(paste.ID)
		assert.NoError(t, err)

		assert.Equal(t, paste.ID, savedPaste.ID)
		assert.Empty(t, savedPaste.Text)
		assert.Empty(t, savedPaste.EncryptedBytes)
		assert.Equal(t, expirationTime.Unix(), savedPaste.ExpiresAt.Unix())
	})

	t.Run("non existent paste", func(t *testing.T) {
		// Test nonexistent paste
		_, err = db.Get("nonexistent")
		assert.Error(t, err)
	})

	t.Run("expired paste", func(t *testing.T) {
		// Test expired paste
		paste, err = db.NewPaste("test paste", time.Now().Add(-time.Hour))
		assert.NoError(t, err)

		_, err = db.Get(paste.ID)
		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	db, err := newTestDB()
	assert.NoError(t, err)
	defer db.reset()

	paste, err := db.NewPaste("test paste", time.Now().Add(time.Hour))
	assert.NoError(t, err)

	err = db.Delete(paste.ID)
	assert.NoError(t, err)

	_, err = db.Get(paste.ID)
	assert.Error(t, err)
}

func TestDecrypt(t *testing.T) {
	db, err := newTestDB()
	assert.NoError(t, err)
	defer db.reset()

	t.Run("decrypt paste", func(t *testing.T) {
		paste, err := db.NewPaste("test paste", time.Now().Add(time.Hour))
		assert.NoError(t, err)

		// decrypt paste
		decryptedText, err := db.Decrypt(paste.ID, paste.Key)
		assert.NoError(t, err)
		assert.Equal(t, "test paste", decryptedText)

		// paste is deleted after decryption
		_, err = db.Get(paste.ID)
		assert.Error(t, err)
	})

	t.Run("invalid paste", func(t *testing.T) {
		paste, err := db.NewPaste("test paste", time.Now().Add(time.Hour))
		assert.NoError(t, err)

		// test wrong key
		_, err = db.Decrypt(paste.ID, "wrong key")
		assert.Error(t, err)

		_, err = db.Decrypt("nonexistent", "key")
		assert.Error(t, err)
	})

	t.Run("expired paste", func(t *testing.T) {
		paste, err := db.NewPaste("test paste", time.Now().Add(-time.Hour))
		assert.NoError(t, err)

		decryptedText, err := db.Decrypt(paste.ID, paste.Key)
		assert.Empty(t, decryptedText)
		assert.Error(t, err)
	})
}

func TestDeleteExpired(t *testing.T) {
	db, err := newTestDB()
	assert.NoError(t, err)
	defer db.reset()

	_, err = db.NewPaste("test paste", time.Now().Add(time.Hour))
	assert.NoError(t, err)

	_, err = db.NewPaste("test paste", time.Now().Add(-time.Hour))
	assert.NoError(t, err)

	err = db.DeleteExpired()
	assert.NoError(t, err)

	pasteCount := 0
	db.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(pastesBucketName)
		if bucket == nil {
			return nil
		}

		cursor := bucket.Cursor()
		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			pasteCount++
		}

		return nil
	})

	assert.Equal(t, 1, pasteCount)
}
