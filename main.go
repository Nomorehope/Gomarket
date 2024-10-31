package main

import (
	"database/sql"
	"final/services"
	"log"
	"regexp"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func ValidatePassword(password string) bool {
	re := regexp.MustCompile("^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])(?=.{8,})")
	return re.MatchString(password)
}

func main() {
	router := gin.Default()

	connStr := "host=localhost port=5432 user=postgres password=q1w2e3r4 dbname=study sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router.Use(func(ctx *gin.Context) {
		ctx.Set("db", db)
		ctx.Next()
	})

	router.POST("/register", services.RegisterUser)
	router.POST("/login", services.LoginUser)
	router.GET("/users/:id", services.GetUserProfile)

	router.Run(":8080")
}
