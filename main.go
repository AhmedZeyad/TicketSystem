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

	ApiRouter(router)

	router.Run(":9090")

}
