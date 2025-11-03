package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"laba6/internal/models"
)

type IKeyStorage interface {
	SaveRsaKeys(keys models.RsaKeys) (int, error)
	GetRsaPublicKey(id int) (string, error)
}

type PostgresKeyStorage struct {
	DB *sql.DB
}

func NewPostgresKeyStorage(db *sql.DB) *PostgresKeyStorage {
	return &PostgresKeyStorage{DB: db}
}

func (s *PostgresKeyStorage) SaveRsaKeys(keys models.RsaKeys) (int, error) {
	query := `INSERT INTO rsa_keys (public_key, private_key) VALUES ($1, $2) RETURNING id`

	var id int
	err := s.DB.QueryRowContext(context.Background(), query, keys.PublicKey, keys.PrivateKey).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to insert keys into postgres: %w", err)
	}

	return id, nil
}

func (s *PostgresKeyStorage) GetRsaPublicKey(id int) (string, error) {
	query := `SELECT public_key FROM rsa_keys WHERE id = $1`

	var publicKey string
	err := s.DB.QueryRowContext(context.Background(), query, id).Scan(&publicKey)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("RSA key pair with ID %d not found: %w", id, sql.ErrNoRows)
		}
		return "", fmt.Errorf("failed to retrieve public key from postgres: %w", err)
	}

	return publicKey, nil
}
