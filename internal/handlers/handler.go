package handlers

import (
	"laba6/internal/processors"
	"laba6/internal/repositories"
)

type Handler struct {
	processors *processors.Processors
	Rsa        *RsaHandler
	Aes        *AesHandler
}

func NewHandler(p *processors.Processors, keyStorage repositories.IKeyStorage) *Handler {
	return &Handler{
		processors: p,
		Rsa:        NewRsaHandler(p.Rsa, keyStorage),
		Aes:        NewAesHandler(p.Aes),
	}
}
