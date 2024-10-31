package services

import (
	"database/sql"
	"errors"
	"final/models"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("secret")

type Claims struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

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

	// Инициализация User_ID
	user.User_ID = uuid.New()

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = hashedPassword

	_, err = db.Exec("INSERT INTO users (user_id, username, email, password) VALUES ($1, $2, $3, $4)",
		user.User_ID, user.Username, user.Email, user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func LoginUser(ctx *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json: "password"`
	}
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := ctx.MustGet("db").(*sql.DB)

	var user models.User

	err := db.QueryRow("SELECT user_id, username, password FROM users WHERE username=$1", loginData.Username).Scan(&user.User_ID, &user.Username, &user.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	if !VerifyPassword(loginData.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := GenerateJWT(user.User_ID.String())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func GenerateJWT(ID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GetUserProfile(ctx *gin.Context) {
	var user models.User
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := ValidateToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	db := ctx.MustGet("db").(*sql.DB)

	err = db.QueryRow("SELECT user_id, username, email FROM users WHERE user_id=$1", claims.UID).Scan(&user.User_ID, &user.Username, &user.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func ValidateToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
