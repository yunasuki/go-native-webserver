package auth

import (
	"go-native-webserver/internal/repositories"
	"go-native-webserver/internal/service"
)

type UserService interface {
	Login(username, password string) (userID int64, userToken int64, err error)
}

type userService struct {
	userRepository repositories.UserRepository
	service.BaseService
}

func NewUserService(ctx service.BaseService) *userService {
	return &userService{
		BaseService: ctx,
	}
}

func (s *userService) Login(username, password string) (userID int64, userToken int64, err error) {
	// Implement login logic here
	return
}
