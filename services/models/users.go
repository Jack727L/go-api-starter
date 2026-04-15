package models

// UpdateUserRequest is the payload for PUT /users/me.
type UpdateUserRequest struct {
	Name *string `json:"name" validate:"omitempty,max=255"`
}

// UserResponse is the shape returned by GET /users/me and PUT /users/me.
type UserResponse struct {
	ID    int32   `json:"id"`
	Email string  `json:"email"`
	Name  *string `json:"name"`
}
