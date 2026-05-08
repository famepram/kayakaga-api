package user

import (
	"errors"
	"kayakaga-api/domain/mysql"

	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) GetUserByID(userID uint) (*mysql.User, error) {
	var user mysql.User
	err := r.db.Preload("Dependent").Preload("RiskProfile").
		First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *repo) GetEmailByUserID(userID uint) (string, error) {
	var cred mysql.UserCredential
	err := r.db.Where("user_id = ?", userID).First(&cred).Error
	if err != nil {
		return "", err
	}
	return cred.Email, nil
}

func (r *repo) UpdateUser(user *mysql.User) error {
	return r.db.Save(user).Error
}
