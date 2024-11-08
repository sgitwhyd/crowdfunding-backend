package user

type RegisterUserInput struct {
	Name       string `json:"name" binding:"required"`
	Occupation string `json:"occupation" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required,email" example:"testing@developer.com"`
	Password string `json:"password" binding:"required" example:"password"`
}

type CheckEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

type FormUpdateUserInput struct {
	Name           string `json:"name" binding:"required"`
	Occupation     string `json:"occupation" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	AvatarFileName string `json:"avatar_file_name"`
}