package user

type RegisterUserResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
}

type GetUserResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}

func FormatUser(user User, token string) RegisterUserResponse {
	userResponse := RegisterUserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Email:      user.Email,
		Token:      token,
	}

	return userResponse
}

func FormatUsers(users []User) []GetUserResponse {
	var usersResponse []GetUserResponse
	for _, user := range users {
		userResponse := GetUserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Occupation: user.Occupation,
			Email:      user.Email,
			Role:       user.Role,
		}

		usersResponse = append(usersResponse, userResponse)
	}

	return usersResponse
}