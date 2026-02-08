package dto


type RegisterRequest struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	PushConsent bool   `json:"push_consent"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UpdateEmailRequest struct {
	NewEmail string `json:"new_email"`
}

type RegisterResponse struct {
	Success string `json:"success"`
	Message string `json:"message"`
	User_id  string `json:"user_id"`
}