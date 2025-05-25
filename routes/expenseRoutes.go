package routes

import (
	"github.com/Divyshekhar/finsnap-go/controllers"
	"github.com/Divyshekhar/finsnap-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterExpenseRoutes(router *gin.Engine) {
	expenseGroup := router.Group("/expense")
	{
		expenseGroup.POST("/create", middleware.RequireAuth(), controllers.CreateExpense)
		expenseGroup.GET("/", middleware.RequireAuth(), controllers.GetExpensesByUserID)
		expenseGroup.GET("/category", middleware.RequireAuth(), controllers.GetExpenseByCategory)
		expenseGroup.PUT("/:id", middleware.RequireAuth(), controllers.EditExpense)
		expenseGroup.DELETE("/delete/:id", middleware.RequireAuth(), controllers.DeleteExpense)
		expenseGroup.GET("/total-expense", middleware.RequireAuth(), controllers.GetTotalExpense)
		expenseGroup.GET("/history/:category", middleware.RequireAuth(), controllers.GetExpenseHistory)

	}
}
