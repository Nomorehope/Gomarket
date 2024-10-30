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

	// Подключение к базе данных
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=study sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Пример Middleware для передачи базы данных в каждый запрос
	router.Use(func(ctx *gin.Context) {
		ctx.Set("db", db)
		ctx.Next()
	})

	// Настраиваем маршруты
	router.POST("/register", services.RegisterUser)
	// Другие маршруты...

	router.Run(":8080")
}
