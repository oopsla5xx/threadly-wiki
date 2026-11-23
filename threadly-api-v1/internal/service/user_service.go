package service

import "threadly-api-v1/internal/repository"

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepository: repository.NewUserRepository(),
	}
}

func (s *UserService) GetUserByID(id string) (interface{}, error) {
	return s.userRepository.GetByID(id)
}
