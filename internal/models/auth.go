package models

type SignInRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type SignInResponse struct {
	Message string `json:"message"`
	UserId  string `json:"user_id"`
}

type CheckCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

type CheckCodeResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
