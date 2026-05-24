package models

type CreateUserRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Phone    string   `json:"phone" binding:"required"`
	Role     UserRole `json:"role" binding:"required,oneof=admin customer staff"`
}

type UpdateUserRequest struct {
	ID       int      `json:"id"`
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Phone    string   `json:"phone" binding:"required"`
	Role     UserRole `json:"role" binding:"required,oneof=admin customer staff"`
}

type UserResponse struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Role     UserRole `json:"role"`
}

func (r *CreateUserRequest) ToModel() *User {
	return &User{
		Username: r.Username,
		Email:    r.Email,
		Phone:    r.Phone,
		Role:     r.Role,
	}
}

func (r *UpdateUserRequest) ToModel() *User {
	return &User{
		ID:       r.ID,
		Username: r.Username,
		Email:    r.Email,
		Phone:    r.Phone,
		Role:     r.Role,
	}
}

func NewUserResponse(u *User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Phone:    u.Phone,
		Role:     u.Role,
	}
}
