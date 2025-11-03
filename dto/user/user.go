package user

// User represents a user in the system.
type User struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john.doe@example.com"`
}
