package user

type RegisterUserResponse struct {
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
}

func FormatUser(user User, token string) RegisterUserResponse {
	userResponse := RegisterUserResponse{
		Name:       user.Name,
		Occupation: user.Occupation,
		Email:      user.Email,
		Token:      token,
	}

	return userResponse
}

type CheckEmailAvailabilityResponse struct {
	IsAvailable bool `json:"is_available"`
}