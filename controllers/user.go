package controllers

import (
	"log"
	"net/http"

	"github.com/Divyshekhar/finsnap-go/initializers"
	"github.com/Divyshekhar/finsnap-go/models"
	"github.com/Divyshekhar/finsnap-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id,omitempty"`
	Name     string    `json:"name" validate:"required,min=3"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=6"`
}

func CreateUser(c *gin.Context) {
	var user User
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"message": "Wrong body"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		log.Fatal("Error creating hash")
	}

	user.Password = string(hashed)
	result := initializers.Db.Create(user)
	if result.Error != nil {
		c.JSON(400, gin.H{"message": "error creating the user"})
		return
	}
	tokenString, err := utils.GenerateJwt(user.Name, user.Email, user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error signing token"})
		return
	}

	c.JSON(200, gin.H{
		"token": tokenString,
	})

}

func LoginUser(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "Wrong Data sent",
		})
		return
	}
	if user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "invalid password or email"})
		return
	}
	var loginUser *User
	err := initializers.Db.First(&loginUser, "email = ?", user.Email).Error
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "wrong email or password",
		})
		return
	}
	errr := bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(user.Password))
	if errr != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "wrong email or password",
		})
		return
	}
	tokenString, err := utils.GenerateJwt(loginUser.Name, loginUser.Email, loginUser.ID.String())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error signing token"})
		return
	}

	ctx.JSON(200, gin.H{
		"token": tokenString,
	})

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
			"error": err,
		})
		return
	} else {
		ctx.JSON(200, gin.H{"updated record": user})
	}

}
