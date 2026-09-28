package response

type UserResponse struct {
	Id       uint
	Username string
	Role     string
}

type LoginResponse struct {
	Token     string
	ExpiresIn int64 // detik
	User      UserResponse
}
