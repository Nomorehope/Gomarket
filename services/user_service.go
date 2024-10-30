package services

import (
	"database/sql"
	"final/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	} else {
		return false
	}
}

func RegisterUser(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := ctx.MustGet("db").(*sql.DB)

	var existingUser models.User
	err := db.QueryRow("SELECT user_id FROM users WHERE username=$1", user.Username).Scan(&existingUser.User_ID)
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = hashedPassword

	_, err = db.Exec("INSERT INTO users (user_id, username, email, password) VALUES ($1, $2, $3, $4)",
		user.User_ID, user.Username, user.Email, user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func LoginUser(ctx *gin.Context) {

}

func GetUserProfile(ctx *gin.Context) {

}
