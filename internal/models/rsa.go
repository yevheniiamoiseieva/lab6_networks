package models

type RsaKeys struct {
	ID         int    `json:"-" db:"id"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}
