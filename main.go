package main

import (
	"TicketSystem/engine"
	"TicketSystem/service"

	"github.com/gin-gonic/gin"
)

func init() {
	println("Hello from init")

	engine.LoadEnvVariables()
	engine.ConecteToDb()

}
func main() {
	println("Hello World")

	router := gin.Default()

	ApiRouter(router)

	router.Run(":9090")

}
func ApiRouter(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/status", func(ctx *gin.Context) {
		_, err := engine.DB.Exec("USE TicketSys")
		if err != nil {

			println("DB is'n working !!")
			ctx.JSON(400, gin.H{
				"message": "DB is'n working ",
			})
			return
		}
		println("DB is working ")
		println("API is working ")
		ctx.JSON(200, gin.H{
			"message": "Hello From API",
			"APi":     "DB is working ",
			"DB":      "API is working ",
		})
	})
	service.TicketRoutes(api)
	service.UserRoutes(api)
	service.TicketReasonRoutes(api)

}
