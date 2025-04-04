package main

import (

	"github.com/AhmedZeyad/TicketSystem/utilities"
	"github.com/gin-gonic/gin"
)

func init() {
	println("Hello from init")

	utilities.LoadEnvVariables()
	utilities.ConecteToDb()
}
func main() {
	println("Hello World")

	router := gin.Default()
	router.GET("/status", func(ctx *gin.Context) {
		_, err := utilities.DB.Exec("USE TicketSys")
		if err != nil {

			println("DB is'n working ‼️")
			ctx.JSON(400, gin.H{
				"message": "DB is'n working ",
			})
			return
		}
		println("DB is working ✅")
		println("API is working ✅")
		ctx.JSON(200, gin.H{
			"message": "Hello From API",
			"APi":     "DB is working ✅",
			"DB":      "API is working ✅",
		})
	})
	

	router.Run(":9090")

}
