package ticket

import (

	"fmt"

	"github.com/AhmedZeyad/TicketSystem/utilities"
)

func GetAllTickets() ([]Ticket, error) {
	var tickets []Ticket

	err := utilities.DB.Select(&tickets, "SELECT * FROM Ticket")
	if err != nil {
		println(err.Error())
		return nil, err
	}
	return tickets, nil
}

func CreateTicket(t Ticket) error {

	_, err := utilities.DB.Exec(`
 	INSERT INTO Ticket (
 		OrderId, vendor_id, reason, rub_reason, item_selected_id, image_url, description,
 		tickets_status, ticket_last_status, count_of_ticket_customer_has, is_ticket_duplicate,
 		customer_feedback, agent_feedback, history_of_customer_feedback, history_of_agent_feedback,
 		agent_internal_notes, feedback_sent_by_agent, history_of_agent_note,
 		ticket_rating, ticket_review, created_at_and_by
 	)
 	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		t.OrderID,
		t.VendorID,
		t.Reason,
		t.RubReason,
		t.ItemSelectedID,
		t.ImageURL,
		t.Description,
		t.TicketsStatus,
		t.TicketLastStatus,
		t.CountOfTicketsCustomerHas,
		t.IsTicketDuplicate,
		t.CustomerFeedback,
		t.AgentFeedback,
		t.HistoryOfCustomerFeedback,
		t.HistoryOfAgentFeedback,
		t.AgentInternalNotes,
		t.FeedbackSentByAgent,
		t.HistoryOfAgentNote,
		t.TicketRating,
		t.TicketReview,
		t.CreatedAtBy,
	)
	if err != nil {
		return err
	}
	return nil
}

func UnpdaeTicket(t Ticket) error {
	query := `
		UPDATE Ticket SET
			OrderId = ?, vendor_id = ?, reason = ?, rub_reason = ?, item_selected_id = ?, image_url = ?, description = ?,
			tickets_status = ?, ticket_last_status = ?, count_of_ticket_customer_has = ?, is_ticket_duplicate = ?,
			customer_feedback = ?, agent_feedback = ?, history_of_customer_feedback = ?, history_of_agent_feedback = ?,
			agent_internal_notes = ?, feedback_sent_by_agent = ?, history_of_agent_note = ?,
			ticket_rating = ?, ticket_review = ?, created_at_and_by = ?
		WHERE id = ?
	`

	_, err := utilities.DB.Exec(query,
		t.OrderID,
		t.VendorID,
		t.Reason,
		t.RubReason,
		t.ItemSelectedID,
		t.ImageURL,
		t.Description,
		t.TicketsStatus,
		t.TicketLastStatus,
		t.CountOfTicketsCustomerHas,
		t.IsTicketDuplicate,
		t.CustomerFeedback,
		t.AgentFeedback,
		t.HistoryOfCustomerFeedback,
		t.HistoryOfAgentFeedback,
		t.AgentInternalNotes,
		t.FeedbackSentByAgent,
		t.HistoryOfAgentNote,
		t.TicketRating,
		t.TicketReview,
		t.CreatedAtBy,
		t.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
func DeleteTicketById(id int) error {

	result, err := utilities.DB.Exec("DELETE FROM Ticket WHERE id = ?", id)
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
