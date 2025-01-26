package auth

type LoginRequest struct {
	Username string `jsonn:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
	User  string `json:"user,omitempty"`
	Error string `json:"error,omitempty"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type VerifyResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}
