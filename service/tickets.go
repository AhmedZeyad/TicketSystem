package service

import (
	"TicketSystem/engine"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Ticket struct {
	ID          int       `json:"id" db:"id"`
	UserId      int       `json:"userId" db:"userId"`
	Reason      string    `json:"reason" db:"reason"`
	RubReason   string    `json:"rub_reason" db:"rub_reason"`
	Description string    `json:"discreption" db:"discreption"`
	AssignTo    *int      `json:"assignTo" db:"assignTo"`
	Status      string    `json:"status" db:"status"` // Serialized JSON
	CreatedAt   time.Time `json:"_" db:"created_at"`
	CreatedBy   int       `json:"_" db:"created_by"`
	UpdatedAt   time.Time `json:"_" db:"updated_at"`
	UpdatedBy   int       `json:"_" db:"updated_by"`
	DeletedAt   time.Time `json:"_" db:"deleted_at"`
	DeletedBy   int       `json:"_" db:"deleted_by"`
}

const (
	Pending            string = "pending"
	InProgress         string = "in_progress"
	WaitingForCustomer string = "waiting_for_customer_response"
	Resolved           string = "resolved"
	NotResolved        string = "not_resolved"
)

const (
	Food      string = "food"
	Delivery  string = "delivery"
	Delivered string = "delivered"
	Canceled  string = "canceled"
)

func GetAllTickets(page int) ([]Ticket, error) {
	var tickets []Ticket
	page = (page - 1) * 10
	err := engine.DB.Select(&tickets, "SELECT id,userId,reason,rub_reason,discreption,assignTo,status FROM tickets limit 10 offset ?", page)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	return tickets, nil
}

func (t Ticket) CreateTicket() error {

	_, err := engine.DB.Exec(`
 	INSERT INTO tickets (userId,reason,rub_reason,discreption,updated_by,created_by)
 	VALUES (?,?,?,?,?,?)`, t.UserId, t.Reason, t.RubReason, t.Description, t.UserId, t.UserId)
	if err != nil {
		return err
	}
	return nil
}

func UnpdaeTicket(t Ticket) error {

	_, err := engine.DB.Exec(`UPDATE tickets SET assignTo = ?,status = ?
			WHERE id = ?`,
		t.AssignTo,
		t.Status,
		t.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
func DeleteTicketById(id int) error {

	result, err := engine.DB.Exec("UPDATE tickets SET deleted_at = ? WHERE id = ?", time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete ticket with id %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not determine rows affected for ticket id %d: %w", id, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no ticket found with id %d", id)
	}

	return nil
}

func GetTickets(context *gin.Context) {
	// strPage := context.Params.ByName("page")
	strPage  := context.DefaultQuery("page", "1") // default to page 1

	log.Println("gglog")
	log.Println(strPage)
	
	// strPage := context.Param("page")

	if strPage == "" || strPage == "0" {
		strPage = "1"
	}
	page, err := strconv.Atoi(strPage)

	if err != nil {
		context.JSON(400, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	t, err := GetAllTickets(page)
	if err != nil {
		println(err.Error())
		context.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	if t != nil {

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
		context.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = t.CreateTicket()
	if err != nil {
		context.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	context.JSON(200, gin.H{
		"message": "sucsess",
		"ticket":  t,
	})

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
	// api.GET("/:page", GetTickets)
	api.GET("/", GetTickets)
	api.POST("/", AddTicket)
	api.PUT("/:id", EditTicket)
	api.DELETE("/:id", DeleteTicket)
}
