package routes

import (
	"github.com/Divyshekhar/finsnap-go/controllers"
	"github.com/Divyshekhar/finsnap-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/user")
	{
		userGroup.GET("/:userid", controllers.GetUserById)
		userGroup.POST("/create", controllers.CreateUser)
		userGroup.POST("/login", controllers.LoginUser)
		userGroup.GET("/all", controllers.GetAllUser)
		userGroup.PUT("/update", middleware.RequireAuth(), controllers.UpdateUser)
	}
}
