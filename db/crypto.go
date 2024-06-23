package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	random "math/rand"
)

func encrypt(text, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, text, nil), nil
}

func decrypt(ciphertext, key []byte) (string, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

type EncryptionKey struct {
	Key []byte
}

func (k *EncryptionKey) Base64Key() string {
	return base64.RawURLEncoding.EncodeToString(k.Key)
}

func NewEncryptionKey() (*EncryptionKey, error) {
	key := make([]byte, keyLength)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}

	return &EncryptionKey{
		Key: key,
	}, nil
}

func NewEncryptionKeyFromBase64(base64Key string) (*EncryptionKey, error) {
	key, err := base64.RawURLEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}

	if len(key) != 32 {
		return nil, errors.New("invalid key length")
	}

	return &EncryptionKey{
		Key: key,
	}, nil
}

func (k *EncryptionKey) Encrypt(text string) ([]byte, error) {
	return encrypt([]byte(text), k.Key)
}

func (k *EncryptionKey) Decrypt(ciphertext []byte) (string, error) {
	return decrypt(ciphertext, k.Key)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const keyLength = 32

func randomKey() string {
	b := make([]rune, keyLength)
	for i := range b {
		b[i] = letterRunes[random.Intn(len(letterRunes))]
	}
	return string(b)
}
