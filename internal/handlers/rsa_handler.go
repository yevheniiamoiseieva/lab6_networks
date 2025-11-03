package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"laba6/internal/processors"
	"laba6/internal/repositories"
)

type RsaHandler struct {
	RsaService processors.IRsaService
	KeyStorage repositories.IKeyStorage
}

func NewRsaHandler(service processors.IRsaService, storage repositories.IKeyStorage) *RsaHandler {
	return &RsaHandler{RsaService: service, KeyStorage: storage}
}

func (h *RsaHandler) GenerateRsaKeys(c *gin.Context) {
	keys, err := h.RsaService.GenerateCryptoKeys()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate RSA keys", "details": err.Error()})
		return
	}

	id, err := h.KeyStorage.SaveRsaKeys(keys)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save keys to storage", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *RsaHandler) GetRsaPublicKey(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format. Must be an integer."})
		return
	}

	publicKey, err := h.KeyStorage.GetRsaPublicKey(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("RSA key pair with ID %d not found.", id)})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal storage error: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"publicKey": publicKey})
}

type EncryptionRequest struct {
	PublicKey string `json:"publicKey"`
	PlainText string `json:"plainText"`
}

type EncryptionResponse struct {
	CipherText string `json:"cipherText"`
}

func (h *RsaHandler) Encrypt(c *gin.Context) {
	var req EncryptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	cipherText, err := h.RsaService.Encrypt(req.PublicKey, req.PlainText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Encryption failed: %s", err)})
		return
	}

	resp := EncryptionResponse{CipherText: cipherText}
	c.JSON(http.StatusOK, resp)
}

type DecryptionRequest struct {
	PrivateKey       string `json:"privateKey"`
	CipherTextBase64 string `json:"cipherText"`
}

type DecryptionResponse struct {
	PlainText string `json:"plainText"`
}

func (h *RsaHandler) Decrypt(c *gin.Context) {
	var req DecryptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	plainText, err := h.RsaService.Decrypt(req.PrivateKey, req.CipherTextBase64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Decryption failed: %s", err)})
		return
	}

	resp := DecryptionResponse{PlainText: plainText}
	c.JSON(http.StatusOK, resp)
}
