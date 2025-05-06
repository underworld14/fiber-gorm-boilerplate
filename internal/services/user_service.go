package services

import (
	"fiber-gorm/internal/models"
	"fiber-gorm/internal/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(user *models.User) error {
	// Add business logic here (validation, etc.)
	return s.Repo.CreateUser(user)
}

func (s *UserService) FindAllUsers() ([]models.User, error) {
	return s.Repo.FindAllUsers()
}

func (s *UserService) FindUserById(id string) (*models.User, error) {
	return s.Repo.FindUserById(id)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.Repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(user *models.User) error {
	return s.Repo.DeleteUser(user)
}
