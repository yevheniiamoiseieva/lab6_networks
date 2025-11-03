package routes

import (
	"github.com/gin-gonic/gin"
	"laba6/internal/handlers"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter(engine *gin.Engine) *Router {
	return &Router{engine: engine}
}

func (r *Router) SetupRoutes(h *handlers.Handler) {
	apiGroup := r.engine.Group("/api")
	{
		cryptoKeysGroup := apiGroup.Group("/crypto-keys")
		{
			cryptoKeysGroup.POST("/generate/rsa-keys", h.Rsa.GenerateRsaKeys)
			cryptoKeysGroup.GET("/rsa-public-key/:id", h.Rsa.GetRsaPublicKey)
		}

		v1 := apiGroup.Group("/v1")
		{
			v1.GET("/employees", h.GetEmployees)
			v1.GET("/employees/:id", h.GetEmployee)
			v1.POST("/employees", h.CreateEmployee)
			v1.PUT("/employees/:id", h.UpdateEmployee)
			v1.DELETE("/employees/:id", h.DeleteEmployee)

			cryptoTestGroup := v1.Group("/crypto")
			{
				cryptoTestGroup.POST("/rsa/encrypt", h.Rsa.Encrypt)
				cryptoTestGroup.POST("/rsa/decrypt", h.Rsa.Decrypt)

				cryptoTestGroup.POST("/aes/generate", h.Aes.GenerateKeys)
				cryptoTestGroup.POST("/aes/encrypt", h.Aes.Encrypt)
				cryptoTestGroup.POST("/aes/decrypt", h.Aes.Decrypt)
			}
		}
	}
}
