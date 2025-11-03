package user

// CreateUserRequest represents the request body for creating a new user.
type CreateUserRequest struct {
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john.doe@example.com"`
}
