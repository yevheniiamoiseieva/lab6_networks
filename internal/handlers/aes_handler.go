// internal/handlers/aes_handler.go
package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"laba6/internal/models"
	"laba6/internal/processors"
	"net/http"
)

type AesHandler struct {
	AesService processors.IAesService
}

func NewAesHandler(service processors.IAesService) *AesHandler {
	return &AesHandler{AesService: service}
}

// GenerateKeys handles the request to generate AES keys.
func (h *AesHandler) GenerateKeys(c *gin.Context) {
	keys, err := h.AesService.GenerateSecretKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate keys: %s", err)})
		return
	}
	c.JSON(http.StatusOK, keys)
}

// Encrypt handles the request to encrypt data.
func (h *AesHandler) Encrypt(c *gin.Context) {
	var req struct {
		AesKey    models.AesKey `json:"aesKey"`
		PlainText string        `json:"plainText"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	cipherText, err := h.AesService.Encrypt(req.AesKey, req.PlainText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Encryption failed: %s", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cipherText": cipherText})
}

// Decrypt handles the request to decrypt data.
func (h *AesHandler) Decrypt(c *gin.Context) {
	var req struct {
		AesKey           models.AesKey `json:"aesKey"`
		CipherTextBase64 string        `json:"cipherText"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	plainText, err := h.AesService.Decrypt(req.AesKey, req.CipherTextBase64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Decryption failed: %s", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plainText": plainText})
}
