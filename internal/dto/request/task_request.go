package request

type TaskCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TaskUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsCompleted *bool   `json:"isCompleted"`
}