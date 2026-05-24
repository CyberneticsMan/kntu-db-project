package models

type CreateUserRequest struct {
	FirstName string   `json:"first_name" binding:"required"`
	LastName  string   `json:"last_name" binding:"required"`
	Email     string   `json:"email" binding:"required,email"`
	Phone     string   `json:"phone" binding:"required"`
	Role      UserRole `json:"role" binding:"omitempty,oneof=admin customer staff"`
	Password  string   `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	ID        int      `json:"id"`
	FirstName string   `json:"first_name" binding:"required"`
	LastName  string   `json:"last_name" binding:"required"`
	Email     string   `json:"email" binding:"required,email"`
	Phone     string   `json:"phone" binding:"required"`
	Role      UserRole `json:"role" binding:"required,oneof=admin customer staff"`
}

type UserResponse struct {
	ID        int      `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	Role      UserRole `json:"role"`
}

func (r *CreateUserRequest) ToModel() *User {
	return &User{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Phone:     r.Phone,
		Password:  r.Password,
	}
}

func (r *UpdateUserRequest) ToModel() *User {
	return &User{
		ID:        r.ID,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Phone:     r.Phone,
		Role:      r.Role,
	}
}

func NewUserResponse(u *User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
	}
}
