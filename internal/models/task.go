package models

type Task struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	DueDate     *string `json:"dueDate"`
	Completed   bool    `json:"completed"`
	CompletedAt *string `json:"completedAt"`
	Importance  int     `json:"importance"`
	Tags        []Tag   `json:"tags"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type CreateTaskInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	DueDate     *string `json:"dueDate"`
	Importance  int     `json:"importance"`
	TagIDs      []int64 `json:"tagIds"`
}

type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Date        *string `json:"date"`
	DueDate     *string `json:"dueDate"`
	Completed   *bool   `json:"completed"`
	Importance  *int    `json:"importance"`
	TagIDs      []int64 `json:"tagIds"`
}

type TaskFilter struct {
	Date       *string `json:"date"`
	TagID      *int64  `json:"tagId"`
	Importance *int    `json:"importance"`
	Completed  *bool   `json:"completed"`
	Query      *string `json:"query"`
}
