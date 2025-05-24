package main

import (
	"github.com/Divyshekhar/finsnap-go/controllers"
	"github.com/Divyshekhar/finsnap-go/initializers"
	"github.com/Divyshekhar/finsnap-go/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectDB()
	// initializers.SyncDb()
}

func main() {
	router := gin.Default()
	router.POST("/create", controllers.CreateUser)

	router.PUT("/update", middleware.RequireAuth(), controllers.UpdateUser)

	router.Run(":3000")

}
