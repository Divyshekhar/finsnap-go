package main

import (
	"github.com/Divyshekhar/finsnap-go/initializers"
	"github.com/Divyshekhar/finsnap-go/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectDB()
	// initializers.SyncDb()
}

func main() {
	router := gin.Default()
	routes.RegisterUserRoutes(router)
	router.Run(":3000")

}
