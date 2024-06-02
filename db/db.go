package db

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"encoding/json"

	"github.com/boltdb/bolt"
)

var (
	pastesBucketName   = []byte("pastes")
	metadataBucketName = []byte("metadata")
)

var (
	ErrPasteNotFound  = errors.New("paste not found")
	ErrBucketNotFound = errors.New("bucket not found")
	ErrPasteExpired   = errors.New("paste expired")
)

type DB struct {
	boltDB *bolt.DB
}

func (d *DB) Close() error {
	return d.boltDB.Close()
}

func (d *DB) NewPaste(text string, expiresAt time.Time) (*Paste, error) {
	paste, err := NewEncryptedPaste(text, expiresAt)
	if err != nil {
		return nil, err
	}
	return paste, d.save(paste)
}

func (d *DB) save(paste *Paste) error {
	return d.boltDB.Update(func(tx *bolt.Tx) error {
		// Save encrypted paste
		pasteBucket, err := tx.CreateBucketIfNotExists(pastesBucketName)
		if err != nil {
			return err
		}

		if err = pasteBucket.Put([]byte(paste.ID), paste.EncryptedBytes); err != nil {
			return err
		}

		// Save metadata to check expiration
		pasteJson, err := json.Marshal(paste)
		if err != nil {
			return err
		}

		metadataBucket, err := tx.CreateBucketIfNotExists(metadataBucketName)
		if err != nil {
			return err
		}
		return metadataBucket.Put([]byte(paste.ID), pasteJson)
	})
}

func (d *DB) Get(id string) (*Paste, error) {
	var paste Paste

	err := d.boltDB.View(func(tx *bolt.Tx) error {
		// get metadata
		bucket := tx.Bucket(metadataBucketName)
		if bucket == nil {
			return ErrBucketNotFound
		}

		// unmarshal metadata
		jsonPaste := bucket.Get([]byte(id))
		if jsonPaste == nil {
			return ErrPasteNotFound
		}

		if err := json.Unmarshal(jsonPaste, &paste); err != nil {
			return err
		}

		// ensure paste is not expired
		if time.Now().After(paste.ExpiresAt) {
			return ErrPasteExpired
		}

		return nil
	})

	return &paste, err
}

func (d *DB) Decrypt(id string, key string) (string, error) {
	// delete paste if expired
	if _, err := d.Get(id); errors.Is(err, ErrPasteExpired) {
		return "", d.Delete(id)
	}

	var decryptedText string
	err := d.boltDB.Update(func(tx *bolt.Tx) error {
		pasteBucket := tx.Bucket(pastesBucketName)
		if pasteBucket == nil {
			return ErrBucketNotFound
		}

		encryptedPaste := pasteBucket.Get([]byte(id))
		if encryptedPaste == nil {
			return ErrPasteNotFound
		}

		var err error
		decryptedText, err = decrypt(encryptedPaste, key)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return decryptedText, d.Delete(id)
}

func (d *DB) Delete(id string) error {
	return d.boltDB.Update(func(tx *bolt.Tx) error {
		pasteBucket := tx.Bucket(pastesBucketName)
		if pasteBucket == nil {
			return ErrBucketNotFound
		}

		if err := pasteBucket.Delete([]byte(id)); err != nil {
			return err
		}

		metadataBucket := tx.Bucket(metadataBucketName)
		if metadataBucket == nil {
			return ErrBucketNotFound
		}

		return metadataBucket.Delete([]byte(id))
	})
}

func (d *DB) DeleteExpired() error {
	var expiredPastes []string
	err := d.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(metadataBucketName)
		if bucket == nil {
			return ErrBucketNotFound
		}

		return bucket.ForEach(func(k, v []byte) error {
			var paste Paste
			if err := json.Unmarshal(v, &paste); err != nil {
				return err
			}

			if time.Now().After(paste.ExpiresAt) {
				expiredPastes = append(expiredPastes, string(k))
			}

			return nil
		})
	})

	if err != nil {
		return fmt.Errorf("error getting expired pastes: %v", err)
	}

	for _, id := range expiredPastes {
		if err := d.Delete(id); err != nil {
			slog.Error("error_deleting_expired_paste", "id", id, "error", err)
		}
	}

	return nil
}

func (d *DB) DeleteExpiredPeriodically(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := d.DeleteExpired(); err != nil {
			slog.Error("error_starting_expired_paste_job", "error", err)
		}
	}
}
