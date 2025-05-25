package routes

import (
	"github.com/Divyshekhar/finsnap-go/controllers"
	"github.com/Divyshekhar/finsnap-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterIncomeRoutes(router *gin.Engine) {
	incomeGroup := router.Group("/income")
	{
		incomeGroup.POST("/create", middleware.RequireAuth(), controllers.CreateIncome)
		incomeGroup.POST("/edit/:income_id", middleware.RequireAuth(), controllers.EditIncome)
		incomeGroup.DELETE("/delete/:income_id", middleware.RequireAuth(), controllers.DeleteIncome)
		incomeGroup.GET("/category", middleware.RequireAuth(), controllers.GetIncomeByCategory)
		incomeGroup.GET("/total", middleware.RequireAuth(), controllers.GetToalIncome)
		incomeGroup.GET("/history/:category", middleware.RequireAuth(), controllers.GetIncomeHistory)


	}
}
