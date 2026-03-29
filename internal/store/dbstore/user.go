package dbstore

import (
	"goth/internal/hash"
	"goth/internal/store"

	"gorm.io/gorm"
)

type UserStore struct {
	db           *gorm.DB
	passwordhash hash.PasswordHash
}

type NewUserStoreParams struct {
	DB           *gorm.DB
	PasswordHash hash.PasswordHash
}

func NewUserStore(params NewUserStoreParams) *UserStore {
	return &UserStore{
		db:           params.DB,
		passwordhash: params.PasswordHash,
	}
}

func (s *UserStore) CreateUser(name string, password string) error {

	hashedPassword, err := s.passwordhash.GenerateFromPassword(password)
	if err != nil {
		return err
	}

	return s.db.Create(&store.User{
		Name:    name,
		Password: hashedPassword,
	}).Error
}

func (s *UserStore) GetUser(name string) (*store.User, error) {

	var user store.User
	err := s.db.Where("name = ?", name).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}
