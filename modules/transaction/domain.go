package transaction

import (
	"kayakaga-api/domain/mysql"
	"time"
)

type Repository interface {
	ListTransactions(userID uint, filters *ListFilters) ([]mysql.Transaction, error)
	GetTransactionByID(id, userID uint) (*mysql.Transaction, error)
	CreateTransaction(txn *mysql.Transaction) error
	UpdateTransaction(txn *mysql.Transaction) error
	DeleteTransaction(id, userID uint) error
	BulkInsert(txns []mysql.Transaction) error
	CheckDuplicates(userID, accountID uint, dates []time.Time, amounts []int64, merchants []string) (map[int]bool, error)
}

type UseCase interface {
	ListTransactions(userID uint, filters *ListFilters) (*ListResponse, error)
	CreateTransaction(userID uint, req *CreateTransactionRequest) (*TransactionResponse, error)
	UpdateTransaction(id, userID uint, req *UpdateTransactionRequest) (*TransactionResponse, error)
	DeleteTransaction(id, userID uint) error
	ImportCSV(userID uint, accountID uint, csvData [][]string) (*ImportResult, error)
	ProcessReceipt(imageData []byte, accountID uint) (*ReceiptData, error)
	ConfirmReceipt(userID uint, req *ConfirmReceiptRequest) (*TransactionResponse, error)
}

type ListFilters struct {
	Period      string
	AccountID   *uint
	CategoryID  *uint
	Merchant    string
	IsRecurring *int8
}
