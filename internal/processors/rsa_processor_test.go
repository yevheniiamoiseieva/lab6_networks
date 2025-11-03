package processors_test

import (
	"laba6/internal/processors"
	"testing"
)

var rsaService = processors.NewRsaService(2048)

func TestRsaService_GenerateCryptoKeys(t *testing.T) {
	keys, err := rsaService.GenerateCryptoKeys()

	if err != nil {
		t.Fatalf("GenerateCryptoKeys failed unexpectedly: %v", err)
	}

	if keys.PublicKey == "" {
		t.Error("PublicKey should not be empty")
	}
	if keys.PrivateKey == "" {
		t.Error("PrivateKey should not be empty")
	}

	if len(keys.PublicKey) < 100 || len(keys.PrivateKey) < 1000 {
		t.Errorf("Keys seem too short (Key length: %d, Private length: %d)", len(keys.PublicKey), len(keys.PrivateKey))
	}
}

func TestRsaService_EncryptDecrypt(t *testing.T) {
	const originalMessage = "Hello, world! This is a secret test message."

	keys, err := rsaService.GenerateCryptoKeys()
	if err != nil {
		t.Fatalf("Setup failed: Could not generate keys: %v", err)
	}

	ciphertext, err := rsaService.Encrypt(keys.PublicKey, originalMessage)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if ciphertext == "" || ciphertext == originalMessage {
		t.Error("Ciphertext is empty or identical to the plaintext, encryption failed.")
	}

	decryptedMessage, err := rsaService.Decrypt(keys.PrivateKey, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decryptedMessage != originalMessage {
		t.Errorf("Decryption mismatch:\nExpected: %s\nActual: %s", originalMessage, decryptedMessage)
	}
}

func TestRsaService_Decrypt_InvalidKey(t *testing.T) {
	const originalMessage = "Test"

	keys1, err := rsaService.GenerateCryptoKeys()
	if err != nil {
		t.Fatalf("Setup failed: Could not generate keys: %v", err)
	}

	keys2, err := rsaService.GenerateCryptoKeys()
	if err != nil {
		t.Fatalf("Setup failed: Could not generate second keyset: %v", err)
	}

	ciphertext, err := rsaService.Encrypt(keys1.PublicKey, originalMessage)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = rsaService.Decrypt(keys2.PrivateKey, ciphertext)

	if err == nil {
		t.Error("Decrypt should have failed when using a wrong private key, but it succeeded.")
	}
}
