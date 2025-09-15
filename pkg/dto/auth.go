package dto

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignupRequest struct {
	Email string `json:"email"    validate:"required,email"`
	Name  string `json:"name"     validate:"required"`
	RegNo string `json:"reg_no"   validate:"required"`
}
