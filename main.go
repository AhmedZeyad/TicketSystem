package main

import (
	"github.com/gin-gonic/gin"
	"github.com/AhmedZeyad/TicketSystem/utilities"
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
