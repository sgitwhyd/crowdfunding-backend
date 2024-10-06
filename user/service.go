package user

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginUserInput) (User, error)
	GetUsers() ([]User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserByID(ID int)(User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if(err != nil){
		log.Println("error hashing", err)
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)

	if err != nil {
		log.Println("error save", err)
		return newUser, err
	}

	return newUser, nil
}

func (s *service) Login(input LoginUserInput) (User, error) {

	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, errors.New("password not match")
	}

	return user, nil
}

func (s *service) GetUsers() ([]User, error){
	users, err := s.repository.GetUsers()
	if err != nil {
		return users, errors.New("no users found")
	}

	return users, nil

}

func (s *service) IsEmailAvailable(input CheckEmailInput)(bool, error){
	// check email is exist on database based on input
	emailInput := input.Email
	user, err := s.repository.FindByEmail(emailInput)
	if err != nil {
		return false, errors.New("email cannot be found")
	}

	if(user.ID == 0){
		return true, nil
	}

	return true, nil
}

func (s *service) SaveAvatar(ID int, fileLocation string)(User, error){
	// get user by ID
	user, err := s.repository.FindById(ID)
	if(err != nil){
		return user, err
	}

	// update attribute avatar file location
	user.AvatarFileName = fileLocation
	// save user
	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil

}

func (s *service) GetUserByID(ID int)(User,error){
	user, err := s.repository.FindById(ID)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("User not found")
	}

	return user, nil
}