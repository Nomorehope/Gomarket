package main

import (
	"database/sql"
	"regexp"

	_ "github.com/lib/pq"
)

func ValidatePassword (password string) bool {
	re := regexp.MustCompile("^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])(?=.{8,})")
	 return re.MatchString(password)
}

func main () {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=study sslmode=disable"
	 db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

}