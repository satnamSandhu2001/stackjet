package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/satnamSandhu2001/stackjet/internal/handlers"
	"github.com/satnamSandhu2001/stackjet/internal/middlewares"
	"github.com/satnamSandhu2001/stackjet/internal/services"
)

func InitRouter(router *gin.Engine, db *sqlx.DB) {

	v1 := router.Group("/api/v1")

	userService := services.NewUserService(db)

	// auth routes
	authHandler := handlers.NewAuthHandler(userService)
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.Signup)
		authGroup.POST("/login", authHandler.Login)
	}

	// user routes
	userHandler := handlers.NewUserHandler(userService)
	userGroup := v1.Group("/users", middlewares.AuthMiddleware(userService))
	{
		userGroup.GET("me", userHandler.GetMyDetails)
		userGroup.GET("", userHandler.ListUsers)
	}

	// stack routes
	stackService := services.NewStackService(db)
	stackHandler := handlers.NewStackHandler(stackService)
	stackGroup := v1.Group("/stack", middlewares.AuthMiddleware(userService))
	{
		stackGroup.GET("/list", stackHandler.ListStacks)
		stackGroup.POST("/new", stackHandler.CreateNewStack)
		stackGroup.POST("/deploy/:id", stackHandler.DeployStack)
	}

}
