package models

// User Struct for request body
type User struct {
	ID               int64
	Name             string `json:"name"`
	Mobile_No        int64  `json:"mob_no"`
	Gender           string `json:"gender"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Confirm_Password string `json:"confirm_password"`
}

type Login struct {
	ID       int64
	Email    string `json:"username"`
	Password string `json:"password"`
	UserID   int64
}


//Response struct for Userdetails
type UserResponse struct {
	Name      string
	Mobile_No int64
	Gender    string
	Email     string
}
