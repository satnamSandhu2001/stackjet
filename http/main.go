package main

import (
	"fmt"

	"github.com/satnamSandhu2001/stackjet/database"
	"github.com/satnamSandhu2001/stackjet/internal/routers"
	"github.com/satnamSandhu2001/stackjet/pkg"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	conn := database.Connect()
	database.RunMigrations(conn)

	r := routers.InitRouter(conn)
	r.SetTrustedProxies(nil)

	if err := r.Run(fmt.Sprintf(":%v", pkg.Config().PORT)); err != nil {
		panic(err)
	}
}
