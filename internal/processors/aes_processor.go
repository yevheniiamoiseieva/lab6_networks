// internal/processors/aes_processor.go
package processors

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"laba6/internal/models" // Adjust the path as needed
)

// IAesService defines the interface for AES operations.
type IAesService interface {
	GenerateSecretKey() (models.AesKey, error)
	Encrypt(aesKey models.AesKey, plainText string) (string, error)        // Returns base64 encoded string
	Decrypt(aesKey models.AesKey, cipherTextBase64 string) (string, error) // cipherText is base64 encoded
}

// AesService implements the IAesService interface.
type AesService struct {
	keySize int // Typically 16, 24, or 32 bytes for AES-128, AES-192, AES-256
}

// NewAesService creates a new AesService instance.
func NewAesService(keySize int) *AesService {
	return &AesService{keySize: keySize}
}

// GenerateSecretKey generates a new AES secret key and IV.
func (s *AesService) GenerateSecretKey() (models.AesKey, error) {
	// 1. Key Generation
	key := make([]byte, s.keySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return models.AesKey{}, fmt.Errorf("failed to generate key: %w", err)
	}

	// 2. IV Generation (IV must be the same size as AES block size, which is always 16 bytes)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return models.AesKey{}, fmt.Errorf("failed to generate IV: %w", err)
	}

	return models.AesKey{
		Key: base64.StdEncoding.EncodeToString(key),
		IV:  base64.StdEncoding.EncodeToString(iv),
	}, nil
}

// Encrypt encrypts the plainText using the AES key and IV.
func (s *AesService) Encrypt(aesKey models.AesKey, plainText string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(aesKey.Key)
	if err != nil {
		return "", fmt.Errorf("invalid key encoding: %w", err)
	}
	iv, err := base64.StdEncoding.DecodeString(aesKey.IV)
	if err != nil {
		return "", fmt.Errorf("invalid IV encoding: %w", err)
	}

	plaintext := []byte(plainText)

	// 1. Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// 2. Create a CFB Encrypter
	stream := cipher.NewCFBEncrypter(block, iv)

	// 3. Encrypt the data
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	// 4. Return Base64 encoded ciphertext
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the base64-encoded cipherText using the AES key and IV.
func (s *AesService) Decrypt(aesKey models.AesKey, cipherTextBase64 string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(aesKey.Key)
	if err != nil {
		return "", fmt.Errorf("invalid key encoding: %w", err)
	}
	iv, err := base64.StdEncoding.DecodeString(aesKey.IV)
	if err != nil {
		return "", fmt.Errorf("invalid IV encoding: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	// 1. Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// 2. Create a CFB Decrypter (must match the mode used for encryption)
	stream := cipher.NewCFBDecrypter(block, iv)

	// 3. Decrypt the data
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	// 4. Return plaintext string
	return string(plaintext), nil
}
