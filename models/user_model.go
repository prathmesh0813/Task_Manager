package models

// User Struct for request body
type User struct {
	ID        int64
	Name      string `json:"name"`
	Mobile_No string `json:"mob_no"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Login struct {
	ID       int64
	Email    string `json:"username"`
	Password string `json:"password"`
	UserID   int64
}

// Response struct for Userdetails
type UserResponse struct {
	Name      string
	Mobile_No int64
	Gender    string
	Email     string
	Avatar    string
}

// Request struct to update user details
type UpdateUserRequest struct {
	Name      string `json:"name" binding:"required"`
	Mobile_No string `json:"mobile_no" binding:"required"`
}

// Request struct to update user password
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
