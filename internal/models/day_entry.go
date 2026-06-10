package models

type DayEntry struct {
	ID         int64  `json:"id"`
	Date       string `json:"date"`
	PlanText   string `json:"planText"`
	ResultText string `json:"resultText"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}
