package account

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

func (r *repo) ListAccounts(userID uint) ([]mysql.Account, error) {
	var accounts []mysql.Account
	err := r.db.Preload("AccountType").
		Where("user_id = ?", userID).
		Order("is_primary DESC, id ASC").
		Find(&accounts).Error
	return accounts, err
}

func (r *repo) GetAccountByID(id, userID uint) (*mysql.Account, error) {
	var account mysql.Account
	err := r.db.Preload("AccountType").
		Where("id = ? AND user_id = ?", id, userID).
		First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return &account, nil
}

func (r *repo) CreateAccount(account *mysql.Account) error {
	return r.db.Create(account).Error
}

func (r *repo) UpdateAccount(account *mysql.Account) error {
	return r.db.Save(account).Error
}

func (r *repo) DeleteAccount(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&mysql.Account{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("account not found")
	}
	return nil
}
