package db

import (
	"slices"
	"testing"
)

func TestEncryptionKey(t *testing.T) {
	// Test creating a new key
	key, err := NewEncryptionKey()
	if err != nil {
		t.Fatal(err)
	}

	if len(key.Key) != 32 {
		t.Fatalf("expected key length to be 32, got %d", len(key.Key))
	}

	if key.Base64Key() == "" {
		t.Fatal("expected base64 key to be non-empty")
	}

	// Test loading key from base64
	loadedKey, err := NewEncryptionKeyFromBase64(key.Base64Key())
	if err != nil {
		t.Fatal(err)
	}

	if key.Base64Key() != loadedKey.Base64Key() {
		t.Fatalf("expected base64 keys to match, got %s and %s", key.Base64Key(), loadedKey.Base64Key())
	}

	if !slices.Equal(key.Key, loadedKey.Key) {
		t.Fatalf("expected keys to match, got %v and %v", key.Key, loadedKey.Key)
	}

}

func TestEncryptDecrypt(t *testing.T) {
	key, err := NewEncryptionKey()
	if err != nil {
		t.Fatal(err)
	}

	plaintext := "hello, world!"
	ciphertext, err := key.Encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := key.Decrypt(ciphertext)
	if err != nil {
		t.Fatal(err)
	}

	if decrypted != plaintext {
		t.Fatalf("expected decrypted text to be %s, got %s", plaintext, decrypted)
	}
}

func TestEncryptDecryptInvalidKey(t *testing.T) {
	key, err := NewEncryptionKey()
	if err != nil {
		t.Fatal(err)
	}

	plaintext := "hello, world!"
	ciphertext, err := key.Encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}

	invalidKey, err := NewEncryptionKey()
	if err != nil {
		t.Fatal(err)
	}

	_, err = invalidKey.Decrypt(ciphertext)
	if err == nil {
		t.Fatal("expected decrypt to fail with invalid key")
	}
}

func TestEncryptDecryptInvalidBase64Key(t *testing.T) {
	if _, err := NewEncryptionKeyFromBase64("invalid"); err == nil {
		t.Fatal("expected loading key from invalid base64 to fail")
	}
}
