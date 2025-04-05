package ticket

import "encoding/json"


type Ticket struct {
	ID                        int             `json:"id" db:"id"`
	OrderID                   int             `json:"OrderId" db:"OrderId"`
	VendorID                  int             `json:"vendor_id" db:"vendor_id"`
	Reason                    string          `json:"reason" db:"reason"`
	RubReason                 string          `json:"rub_reason" db:"rub_reason"`
	ItemSelectedID            int             `json:"item_selected_id" db:"item_selected_id"`
	ImageURL                  string          `json:"image_url" db:"image_url"`
	Description               string          `json:"description" db:"description"`
	TicketsStatus             json.RawMessage `json:"tickets_status" db:"tickets_status"` // Serialized JSON
	TicketLastStatus          json.RawMessage `json:"ticket_last_status" db:"ticket_last_status"`
	CountOfTicketsCustomerHas int             `json:"count_of_ticket_customer_has" db:"count_of_ticket_customer_has"`
	IsTicketDuplicate         bool            `json:"is_ticket_duplicate" db:"is_ticket_duplicate"`

	CustomerFeedback          json.RawMessage `json:"customer_feedback" db:"customer_feedback"`                       // Serialized JSON
	AgentFeedback             json.RawMessage `json:"agent_feedback" db:"agent_feedback"`                             // Serialized JSON
	HistoryOfCustomerFeedback json.RawMessage `json:"history_of_customer_feedback" db:"history_of_customer_feedback"` // Serialized JSON
	HistoryOfAgentFeedback    json.RawMessage `json:"history_of_agent_feedback" db:"history_of_agent_feedback"`       // Serialized JSON

	AgentInternalNotes  json.RawMessage `json:"agent_internal_notes" db:"agent_internal_notes"`     // Serialized JSON
	FeedbackSentByAgent json.RawMessage `json:"feedback_sent_by_agent" db:"feedback_sent_by_agent"` // Serialized JSON
	HistoryOfAgentNote  json.RawMessage `json:"history_of_agent_note" db:"history_of_agent_note"`   // Serialized JSON
	TicketRating        int             `json:"ticket_rating" db:"ticket_rating"`
	TicketReview        string          `json:"ticket_review" db:"ticket_review"`
	CreatedAtBy         json.RawMessage `json:"created" db:"created_at_and_by"` // Serialized JSON
}

type EventAtBy struct {
	At string `json:"at" db:"at"`
	By int    `json:"by" db:"by"`
}

type TicketStatus struct {
	States []TicketState `json:"States"`
}

type TicketState struct {
	TicketState string    `json:"status"`
	CreatedAtBy EventAtBy `json:"created"`
	UpdatedAtBy EventAtBy `json:"updated"`
}

type HistoryOfAgentNote struct {
	Note []note `json:"note"`
}

type note struct {
	Notes       string    `json:"notes"`
	CreatedAtBy EventAtBy `json:"created_at_and_by"`
}

type History_of_feedback struct {
	Feedbacks []Feedback `json:"feedbacks"`
}

type Feedback struct {
	Feedback    string    `json:"feedback"`
	CreatedAtBy EventAtBy `json:"created"`
}

type HistoryOfAgentFeedback struct {
	Feedback    string    `json:"feedback"`
	CreatedAtBy EventAtBy `json:"created_at_and_by"`
}

type HistoryOfCustomerFeedback struct{}

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
