package controllers

import (
	"net/http"
	"time"

	"github.com/Divyshekhar/finsnap-go/initializers"
	"github.com/Divyshekhar/finsnap-go/models"
	"github.com/gin-gonic/gin"
)

type UpdateIncomeInput struct {
	Title       *string    `json:"title"`
	Amount      *float64   `json:"amount"`
	Type        *string    `json:"type"`
	Date        *time.Time `json:"date"`
	Category    *string    `json:"category"`
	Description *string    `json:"description"`
}
type IncomeCategorySummary struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

func CreateIncome(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	userId := userIdRaw.(string)

	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not find user",
		})
		return
	}
	var input models.Income
	err := ctx.ShouldBindBodyWithJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Body",
		})
		return
	}
	input.UserID = userId
	if err := initializers.Db.Create(&input).Error; err != nil {
		ctx.JSON(400, gin.H{
			"message": "Error creating income",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "Created Successfully",
		"record":  input,
	})
}

func EditIncome(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("user_id")
	incomeId := ctx.Param("income_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	userId := userIdRaw.(string)

	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not find user",
		})
		return
	}
	var input UpdateIncomeInput
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "bad gateway",
		})
		return
	}

	var foundIncome models.Income
	if err := initializers.Db.First(&foundIncome, "id = ? AND user_id = ?", incomeId, userId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Income not found or access denied",
		})
		return
	}
	if input.Title != nil {
		foundIncome.Title = *input.Title
	}
	if input.Amount != nil {
		foundIncome.Amount = *input.Amount
	}
	if input.Type != nil {
		foundIncome.Type = *input.Type
	}
	if input.Date != nil {
		foundIncome.Date = *input.Date
	}
	if input.Category != nil {
		foundIncome.Category = *input.Category
	}
	if input.Description != nil {
		foundIncome.Description = *input.Description
	}
	if err := initializers.Db.Save(&foundIncome).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "could not update",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message":        "Income Updated Successfully",
		"updated record": foundIncome,
	})

}
func DeleteIncome(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("user_id")
	incomeId := ctx.Param("income_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	userId := userIdRaw.(string)

	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not find user",
		})
		return
	}
	var toBeDeleted models.Income

	if err := initializers.Db.First(&toBeDeleted, "id = ? AND user_id = ?", incomeId, userId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Could not find the record to be deleted or access denied",
		})
		return
	}
	if err := initializers.Db.Delete(&toBeDeleted).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not delete record",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "Deleted the record successfully",
	})
}

func GetIncomeByCategory(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	userId := userIdRaw.(string)

	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not find user",
		})
		return
	}
	var results []IncomeCategorySummary
	if err := initializers.Db.
		Model(&models.Income{}).
		Select("category, SUM(amount) as total").
		Where("user_id = ?", userId).
		Group("category").
		Scan(&results).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "coud not retrieve",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "found",
		"found":   results,
	})
}
func GetToalIncome(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	userId := userIdRaw.(string)

	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not find user",
		})
		return
	}
	var total float64
	if err := initializers.Db.
		Model(&models.Income{}).
		Select("COALESCE(SUM(amount), 0)"). // Ensures 0 if nothing is found
		Where("user_id = ?", userId).
		Scan(&total).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error calculating the total income",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": total,
	})

}
func GetIncomeHistory(ctx *gin.Context) {
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

	// Check if user exists
	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user"})
		return
	}

	var incomes []models.Income
	if err := initializers.Db.
		Select("id", "title", "amount", "date").
		Where("user_id = ? AND category = ?", userId, category).
		Order("created_at DESC").
		Find(&incomes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving incomes"})
		return
	}

	// Check if no records found
	if len(incomes) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "No incomes found for this category",
			"data":    []models.Income{},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Incomes fetched",
		"data":    incomes,
	})
}
