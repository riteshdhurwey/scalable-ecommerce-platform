package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/handler"
	"github.com/riteshdhurwey/scalable-ecommerce-platform/services/user/internal/middleware"
)

type Router struct {
	engine  *gin.Engine
	handler *handler.UserHandler
}

func InitRouter(engine *gin.Engine, handler *handler.UserHandler) *Router {
	return &Router{
		engine:  engine,
		handler: handler,
	}
}

func (r *Router) UserRoutes() {
	userRoutes := r.engine.Group("/api/v1/users")
	userRoutes.POST("", r.handler.Create)

	// Public route
	userRoutes.POST("/login", r.handler.Login)

	// Protected routes
	protected := userRoutes.Group("")
	protected.Use(middleware.AuthMiddleware)

	protected.GET("/me", r.handler.GetProfile)
	protected.DELETE("/me/delete",r.handler.DeleteUser)
}
