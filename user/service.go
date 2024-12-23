package user

import (
	"errors"
	"log"
	"mime/multipart"
	"strconv"
	"strings"

	cloud "be-bwastartup/cloudinary"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	GetUserByID(ID int) (User, error)
	Login(input LoginUserInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	UploadAvatar(ID int, file *multipart.FileHeader) (User, error)
	UpdateUser(userID int, input FormUpdateUserInput) (User, error)
	GetAllUsers() ([]User, error)
}

type service struct {
	repository Repository
	cloudinary cloud.Service
}

func NewService(repository Repository, cloudinary cloud.Service) *service {
	return &service{repository, cloudinary}
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

func (s *service) GetUserByID(ID int) (User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) Login(input LoginUserInput) (User, error){
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, errors.New("credential doesnt match")
	}

	return user, nil
}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error){
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return true, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) UploadAvatar(ID int, file *multipart.FileHeader) (User, error){
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	if user.AvatarFileName != "" {
		url := strings.Split(user.AvatarFileName, "/")
		path := strings.Join(url[len(url)-3:], "/")
		publicID := strings.Split(path, ".")[0]

		_, err := s.cloudinary.DeleteImage(publicID)
		if err != nil {
			return User{}, err
		}
	}

	ImageUrl, err := s.cloudinary.UploadImage(file, "user", strconv.Itoa(user.ID))
	if err != nil {
		return User{}, err
	}

	user.AvatarFileName = ImageUrl


	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *service) UpdateUser(userID int, input FormUpdateUserInput)(User, error){
	user, err := s.repository.FindByID(userID)
	if err != nil {
		return user, err
	}

	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *service) GetAllUsers()([]User, error){
	users, err := s.repository.FindAll()
	if err != nil {
		return users, err
	}

	return users, nil
}