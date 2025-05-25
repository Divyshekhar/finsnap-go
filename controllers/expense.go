package controllers

import (
	"net/http"
	"time"

	"github.com/Divyshekhar/finsnap-go/initializers"
	"github.com/Divyshekhar/finsnap-go/models"
	"github.com/gin-gonic/gin"
)

type UpdateExpenseInput struct {
	Title       *string    `json:"title"`
	Amount      *float64   `json:"amount"`
	Type        *string    `json:"type"`
	Date        *time.Time `json:"date"`
	Category    *string    `json:"category"`
	Description *string    `json:"description"`
}

type ExpenseCategorySummary struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

func CreateExpense(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userId := userIdRaw.(string)

	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not find user"})
		return
	}

	var input models.Expense
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Body"})
		return
	}

	input.UserID = userId
	if err := initializers.Db.Create(&input).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error creating expense"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Created Successfully",
		"record":  input,
	})
}

func EditExpense(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("user_id")
	expenseId := ctx.Param("expense_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userId := userIdRaw.(string)

	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not find user"})
		return
	}

	var input UpdateExpenseInput
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Body"})
		return
	}

	var foundExpense models.Expense
	if err := initializers.Db.First(&foundExpense, "id = ? AND user_id = ?", expenseId, userId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Expense not found or access denied"})
		return
	}

	if input.Title != nil {
		foundExpense.Title = *input.Title
	}
	if input.Amount != nil {
		foundExpense.Amount = *input.Amount
	}
	if input.Type != nil {
		foundExpense.Type = *input.Type
	}
	if input.Date != nil {
		foundExpense.Date = *input.Date
	}
	if input.Category != nil {
		foundExpense.Category = *input.Category
	}
	if input.Description != nil {
		foundExpense.Description = *input.Description
	}

	if err := initializers.Db.Save(&foundExpense).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not update"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":        "Expense Updated Successfully",
		"updated record": foundExpense,
	})
}

func DeleteExpense(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("user_id")
	expenseId := ctx.Param("expense_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userId := userIdRaw.(string)

	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not find user"})
		return
	}

	var toBeDeleted models.Expense
	if err := initializers.Db.First(&toBeDeleted, "id = ? AND user_id = ?", expenseId, userId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Could not find the record to be deleted or access denied"})
		return
	}

	if err := initializers.Db.Delete(&toBeDeleted).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete record"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Deleted the record successfully"})
}

func GetExpenseByCategory(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userId := userIdRaw.(string)

	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not find user"})
		return
	}

	var results []ExpenseCategorySummary
	if err := initializers.Db.
		Model(&models.Expense{}).
		Select("category, SUM(amount) as total").
		Where("user_id = ?", userId).
		Group("category").
		Scan(&results).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not retrieve"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Found",
		"found":   results,
	})
}

func GetTotalExpense(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userId := userIdRaw.(string)

	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not find user"})
		return
	}

	var total float64
	if err := initializers.Db.
		Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ?", userId).
		Scan(&total).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error calculating the total expense"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"total": total})
}

func GetExpenseHistory(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userId := userIdRaw.(string)

	category := ctx.Param("category")
	if category == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Category is required"})
		return
	}

	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user"})
		return
	}

	var expenses []models.Expense
	if err := initializers.Db.
		Select("id", "title", "amount", "date").
		Where("user_id = ? AND category = ?", userId, category).
		Order("created_at DESC").
		Find(&expenses).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving expenses"})
		return
	}

	if len(expenses) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "No expenses found for this category",
			"data":    []models.Expense{},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Expenses fetched",
		"data":    expenses,
	})
}

func GetExpensesByUserID(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	var expenses []models.Expense
	if err := initializers.Db.Where("user_id = ?", userID).Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	c.JSON(http.StatusOK, expenses)
}
