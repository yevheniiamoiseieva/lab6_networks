package processors_test

import (
	"laba6/internal/processors"
	"testing"
)

var aesService = processors.NewAesService(32)

func TestAesService_GenerateSecretKey(t *testing.T) {
	keys, err := aesService.GenerateSecretKey()

	if err != nil {
		t.Fatalf("GenerateSecretKey failed unexpectedly: %v", err)
	}

	if keys.Key == "" {
		t.Error("Secret Key should not be empty")
	}
	if keys.IV == "" {
		t.Error("Initialization Vector (IV) should not be empty")
	}

	if len(keys.Key) != 44 {
		t.Errorf("Expected key length 44 (32 bytes Base64), got %d", len(keys.Key))
	}
	if len(keys.IV) != 24 {
		t.Errorf("Expected IV length 24 (16 bytes Base64), got %d", len(keys.IV))
	}
}

func TestAesService_EncryptDecrypt(t *testing.T) {
	const originalMessage = "The quick brown fox jumps over the lazy dog. AES test message."

	aesKey, err := aesService.GenerateSecretKey()
	if err != nil {
		t.Fatalf("Setup failed: Could not generate keys: %v", err)
	}

	ciphertext, err := aesService.Encrypt(aesKey, originalMessage)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if ciphertext == "" {
		t.Fatal("Ciphertext is empty, encryption failed.")
	}

	decryptedMessage, err := aesService.Decrypt(aesKey, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decryptedMessage != originalMessage {
		t.Errorf("Decryption mismatch:\nExpected: %s\nActual: %s", originalMessage, decryptedMessage)
	}
}

func TestAesService_Decrypt_InvalidKey(t *testing.T) {
	const originalMessage = "Test data"

	keys1, err := aesService.GenerateSecretKey()
	if err != nil {
		t.Fatalf("Setup failed: Could not generate keys: %v", err)
	}

	keys2, err := aesService.GenerateSecretKey()
	if err != nil {
		t.Fatalf("Setup failed: Could not generate second keyset: %v", err)
	}

	ciphertext, err := aesService.Encrypt(keys1, originalMessage)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	badKeySet := keys2
	badKeySet.IV = keys1.IV

	_, err = aesService.Decrypt(badKeySet, ciphertext)

	if err == nil {
		invalidFormatKey := keys1
		invalidFormatKey.Key = "This-is-not-base64-key"
		_, err = aesService.Decrypt(invalidFormatKey, ciphertext)

		if err == nil {
			t.Error("Decrypt should have failed due to invalid key format, but it succeeded.")
		}
	}
}
