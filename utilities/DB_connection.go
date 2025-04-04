package utilities

import (


	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

)


var DB *sqlx.DB


func ConecteToDb() {
	dsn := os.Getenv("DSN")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		println("fill to conect")
		println(err.Error())

		return
	}
	DB =db
}