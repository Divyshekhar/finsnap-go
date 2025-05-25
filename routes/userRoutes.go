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
		userGroup.POST("/singup", controllers.CreateUser)
		userGroup.POST("/signin", controllers.LoginUser)
		userGroup.GET("/", controllers.GetAllUser)
		userGroup.PUT("/update", middleware.RequireAuth(), controllers.UpdateUser)
	}
}
