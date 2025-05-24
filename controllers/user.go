package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Divyshekhar/finsnap-go/initializers"
	"github.com/Divyshekhar/finsnap-go/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id,omitempty"`
	Name     string    `json:"name" validate:"required,min=3"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=6"`
}

func CreateUser(c *gin.Context) {
	var user *User
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"message": "Wrong body"})
		return
	}
	result := initializers.Db.Create(user)
	if result.Error != nil {
		c.JSON(400, gin.H{"message": "error creating the user"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":   user.Name,
		"userId": user.ID,
		"email":  user.Email,
		"exp":    time.Now().Add(24*time.Hour).Unix(),
	})
	secret := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error signing token"})
		return
	}
	c.JSON(200, gin.H{"token": tokenString})

}

func UpdateUser(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
		return
	}
	userId := userIdRaw.(string)

	var inputs models.User
	if err := ctx.ShouldBindBodyWithJSON(&inputs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}
	var user models.User
	if err := initializers.Db.First(&user, "id = ?", userId).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "could not find the user",
		})
		return
	}
	if inputs.Email != "" {
		user.Email = inputs.Email
	}
	if inputs.Name != "" {

		user.Name = inputs.Name
	}

	if err := initializers.Db.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not update",
		})
		return
	}else{
		ctx.JSON(200, gin.H{"updated record": user})
	}

}
