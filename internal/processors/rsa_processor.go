package processors

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"laba6/internal/models"
)

type IRsaService interface {
	GenerateCryptoKeys() (models.RsaKeys, error)
	Encrypt(publicKey, plainText string) (string, error)
	Decrypt(privateKey, cipherTextBase64 string) (string, error)
}

type RsaService struct {
	bits int // Key size in bits (e.g., 2048, 4096)
}

func NewRsaService(bits int) *RsaService {
	return &RsaService{bits: bits}
}

// GenerateCryptoKeys generates a new RSA key pair.
func (s *RsaService) GenerateCryptoKeys() (models.RsaKeys, error) {
	// 1. Generate Private Key
	privateKey, err := rsa.GenerateKey(rand.Reader, s.bits)
	if err != nil {
		return models.RsaKeys{}, fmt.Errorf("failed to generate private key: %w", err)
	}

	// 2. Encode Private Key to PEM format
	privatePEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	// 3. Encode Public Key to PEM format
	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return models.RsaKeys{}, fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicBytes,
		},
	)

	return models.RsaKeys{
		PrivateKey: string(privatePEM),
		PublicKey:  string(publicPEM),
	}, nil
}

func (s *RsaService) Encrypt(publicKeyPEM, plainText string) (string, error) {
	// 1. Decode Public Key from PEM
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode public key PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %w", err)
	}
	rsaPubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("key is not an RSA public key")
	}

	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		rsaPubKey,
		[]byte(plainText),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *RsaService) Decrypt(privateKeyPEM, cipherTextBase64 string) (string, error) {
	// 1. Decode Private Key from PEM
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode private key PEM block")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// 2. Decode the Base64 ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	plaintextBytes, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		priv,
		ciphertext,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintextBytes), nil
}
