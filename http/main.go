package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/satnamSandhu2001/stackjet/database"
	"github.com/satnamSandhu2001/stackjet/internal/routers"
	"github.com/satnamSandhu2001/stackjet/pkg"
	"github.com/satnamSandhu2001/stackjet/pkg/initializer"
)

func init() {
	initializer.InitializeApp(false)
}

func main() {
	log.Println("GO_ENV mode:", pkg.Config().GO_ENV)
	conn := database.Connect()
	defer conn.Close()

	r := gin.Default()
	r.SetTrustedProxies(nil)

	routers.InitRouter(r, conn)

	if err := r.Run(fmt.Sprintf(":%v", pkg.Config().PORT)); err != nil {
		panic(err)
	}
}
