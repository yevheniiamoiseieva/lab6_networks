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
	v1 := r.engine.Group("/v1")
	{
		v1.GET("/employees", h.GetEmployees)
		v1.GET("/employees/:id", h.GetEmployee)
		v1.POST("/employees", h.CreateEmployee)
		v1.PUT("/employees/:id", h.UpdateEmployee)
		v1.DELETE("/employees/:id", h.DeleteEmployee)
	}
}
