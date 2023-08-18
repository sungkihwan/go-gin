package usecases

import (
	"go-gin-postgre/domain"
)

type UserUsecase interface {
	GetUsers() ([]domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(id uint) error
}

type userUsecaseImpl struct {
	repo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) UserUsecase {
	return &userUsecaseImpl{repo}
}

func (u *userUsecaseImpl) GetUsers() ([]domain.User, error) {
	return u.repo.FindAll()
}

func (u *userUsecaseImpl) GetUserByID(id uint) (*domain.User, error) {
	return u.repo.FindByID(id)
}

func (u *userUsecaseImpl) CreateUser(user *domain.User) error {
	return u.repo.Create(user)
}

func (u *userUsecaseImpl) UpdateUser(user *domain.User) error {
	return u.repo.Update(user)
}

func (u *userUsecaseImpl) DeleteUser(id uint) error {
	return u.repo.Delete(id)
}
