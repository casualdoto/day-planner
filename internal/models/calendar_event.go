package models

type CalendarEvent struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	StartDate   string  `json:"startDate"`
	EndDate     *string `json:"endDate"`
	Completed   bool    `json:"completed"`
	CompletedAt *string `json:"completedAt"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type CreateCalendarEventInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	StartDate   string  `json:"startDate"`
	EndDate     *string `json:"endDate"`
}
