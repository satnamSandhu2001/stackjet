package routers

import (
	"github.com/satnamSandhu2001/stackjet/internal/handlers"
	"github.com/satnamSandhu2001/stackjet/internal/middlewares"
	"github.com/satnamSandhu2001/stackjet/internal/services"

	"github.com/gin-gonic/gin"

	"github.com/jmoiron/sqlx"
)

func RegisterUserRouter(rg *gin.RouterGroup, db *sqlx.DB) {

	s := services.NewUserService(db)
	h := handlers.NewUserHandler(s)

	usersGroup := rg.Group("/users")
	authGroup := usersGroup.Group("", middlewares.AuthMiddleware(s))
	{
		authGroup.GET("me", h.GetMyDetails)
		authGroup.GET("", h.ListUsers)
	}

}
