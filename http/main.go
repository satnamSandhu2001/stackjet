package main

import (
	"fmt"

	"github.com/satnamSandhu2001/stackjet/database"
	"github.com/satnamSandhu2001/stackjet/internal/routers"
	"github.com/satnamSandhu2001/stackjet/pkg"
	"github.com/satnamSandhu2001/stackjet/pkg/initializer"
)

func init() {
	initializer.InitializeApp(false)
}

func main() {
	conn := database.Connect()
	defer conn.Close()

	r := routers.InitRouter(conn)
	r.SetTrustedProxies(nil)

	if err := r.Run(fmt.Sprintf(":%v", pkg.Config().PORT)); err != nil {
		panic(err)
	}
}
