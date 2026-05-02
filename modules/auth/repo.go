package auth

import (
	"errors"
	"kayakaga-api/domain/mysql"
	"time"

	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) CreateUser(user *mysql.User) error {
	return r.db.Create(user).Error
}

func (r *repo) CreateUserCredential(cred *mysql.UserCredential) error {
	return r.db.Create(cred).Error
}

func (r *repo) CreateRefreshToken(token *mysql.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *repo) GetUserByEmail(email string) (*mysql.User, *mysql.UserCredential, error) {
	var cred mysql.UserCredential
	err := r.db.Where("email = ?", email).First(&cred).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("user not found")
		}
		return nil, nil, err
	}

	var user mysql.User
	err = r.db.Where("id = ?", cred.UserID).First(&user).Error
	if err != nil {
		return nil, nil, err
	}

	return &user, &cred, nil
}

func (r *repo) GetRefreshToken(token string) (*mysql.RefreshToken, error) {
	var rt mysql.RefreshToken
	err := r.db.Where("token = ? AND revoked = 0", token).
		Where("expires_at > ?", time.Now().UTC()).
		First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired refresh token")
		}
		return nil, err
	}
	return &rt, nil
}

func (r *repo) RevokeRefreshToken(token string) error {
	return r.db.Model(&mysql.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked", 1).Error
}
