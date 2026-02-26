package service

import (
	"golang/internal/core/models"
	"golang/internal/core/repository"
)
type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(user models.User) error {
	return s.Repo.Create(user)
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) SearchUsers(keyword string) ([]models.User, error) {
	return s.Repo.Search(keyword)
}
