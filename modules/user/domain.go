package user

import (
	"kayakaga-api/domain/mysql"
)

type Repository interface {
	GetUserByID(userID uint) (*mysql.User, error)
	GetEmailByUserID(userID uint) (string, error)
	UpdateUser(user *mysql.User) error
}

type UseCase interface {
	GetProfile(userID uint) (*ProfileResponse, error)
	UpdateProfile(userID uint, req *UpdateProfileRequest) (*ProfileResponse, error)
}
