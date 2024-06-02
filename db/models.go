package db

import (
	"time"
)

type Paste struct {
	ID             string    `json:"id"`
	Text           string    `json:"-"`
	EncryptedBytes []byte    `json:"-"`
	Key            string    `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	ExpiresAt      time.Time `json:"expiresAt"`
}

func NewEncryptedPaste(text string, expiresAt time.Time) (*Paste, error) {
	key := randomKey()
	encryptedText, err := encrypt(text, key)
	if err != nil {
		return nil, err
	}

	return &Paste{
		ID:             randomKey(),
		Text:           text,
		EncryptedBytes: encryptedText,
		Key:            key,
		CreatedAt:      time.Now(),
		ExpiresAt:      expiresAt,
	}, nil
}
