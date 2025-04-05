package ticket

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTickets(context *gin.Context) {
	println("hi from tickets")
	t, err := GetAllTickets()
	if err != nil {
		println(err.Error())
		context.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	if t != nil {
		var TicketsStatus json.RawMessage = []byte(t[0].TicketsStatus)

		var ticketStatus TicketStatus

		err := json.Unmarshal(TicketsStatus, &ticketStatus)
		if err != nil {
		context.JSON(500, gin.H{
			"message": err.Error(),
		})
		}

		context.JSON(200, gin.H{
			"status":  "Sucsess",
			"tickets": t,
		})

	}

}
func AddTicket(context *gin.Context) {
	var t Ticket
	err := context.BindJSON(&t)
	if err != nil {
		println(err.Error())
		context.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = CreateTicket(t)
	if err != nil {
		context.JSON(400, gin.H{
			"message": err.Error(),
		})
	}
	context.JSON(200, gin.H{
		"message": "sucsess",
		"ticket":  t,
	})
	println("sgsfadsfasfds")

	var jsonagentfeedback json.RawMessage = []byte(t.AgentFeedback)
	var aa Feedback
	err = json.Unmarshal(jsonagentfeedback, &aa)
	if err != nil {
		println(err.Error())
	}
	println(aa.Feedback)

}
func EditTicket(context *gin.Context) {
	ticketID := context.Param("id")

	var updatedTicket Ticket
	if err := context.ShouldBindJSON(&updatedTicket); err != nil {
		context.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	updatedTicket.ID, _ = strconv.Atoi(ticketID)

	if err := UnpdaeTicket(updatedTicket); err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "Ticket updated successfully"})
}
func DeleteTicket(context *gin.Context) {
	strticketID := context.Param("id")
	ticketID, _ := strconv.Atoi(strticketID)

	if err := DeleteTicketById(ticketID); err != nil {

		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, gin.H{"message": "Ticket deleted successfully"})
}
func TicketRoutes(rg *gin.RouterGroup) {
	api := rg.Group("tickets")
	api.GET("/", GetTickets)
	api.POST("/", AddTicket)
	api.PUT("/:id", EditTicket)
	api.DELETE("/:id", DeleteTicket)
}
