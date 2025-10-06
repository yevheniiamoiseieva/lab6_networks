package handlers

import (
	"laba6/internal/processors"
)

type Handler struct {
	processors *processors.Processors
}

func NewHandler(p *processors.Processors) *Handler {
	return &Handler{
		processors: p,
	}
}
