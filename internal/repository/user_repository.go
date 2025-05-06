package repository

import (
	"fiber-gorm/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindAllUsers() ([]models.User, error) {
	var users []models.User
	return users, r.DB.Find(&users).Error
}

func (r *UserRepository) FindUserById(id string) (*models.User, error) {
	var user models.User
	return &user, r.DB.Where("id = ?", id).First(&user).Error
}

func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	return &user, r.DB.First(&user, "email = ?", email).Error
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) DeleteUser(user *models.User) error {
	return r.DB.Delete(user).Error
}
