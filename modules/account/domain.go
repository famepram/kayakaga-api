package account

import (
	"kayakaga-api/domain/mysql"
)

type Repository interface {
	ListAccounts(userID uint) ([]mysql.Account, error)
	GetAccountByID(id, userID uint) (*mysql.Account, error)
	CreateAccount(account *mysql.Account) error
	UpdateAccount(account *mysql.Account) error
	DeleteAccount(id, userID uint) error
}

type UseCase interface {
	ListAccounts(userID uint) ([]AccountResponse, error)
	GetBalances(userID uint) (*BalancesResponse, error)
	CreateAccount(userID uint, req *CreateAccountRequest) (*AccountResponse, error)
	UpdateAccount(id, userID uint, req *UpdateAccountRequest) (*AccountResponse, error)
	DeleteAccount(id, userID uint) error
}
