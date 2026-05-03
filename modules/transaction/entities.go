package transaction

import "time"

type TransactionResponse struct {
	ID          uint       `json:"id"`
	AccountID   uint       `json:"account_id"`
	CategoryID  uint       `json:"category_id"`
	SourceID    uint       `json:"source_id"`
	Date        string     `json:"date"`
	Time        *string    `json:"time"`
	Merchant    string     `json:"merchant"`
	Amount      int64      `json:"amount"`
	Notes       string     `json:"notes"`
	IsRecurring int8       `json:"is_recurring"`
}

type CreateTransactionRequest struct {
	AccountID   uint       `json:"account_id" binding:"required"`
	CategoryID  uint       `json:"category_id" binding:"required"`
	SourceID    uint       `json:"source_id"`
	Date        time.Time  `json:"date" binding:"required"`
	Time        *string    `json:"time"`
	Merchant    string     `json:"merchant" binding:"required"`
	Amount      int64      `json:"amount" binding:"required"`
	Notes       string     `json:"notes"`
	IsRecurring int8       `json:"is_recurring"`
}

type UpdateTransactionRequest struct {
	AccountID   uint
	CategoryID  uint
	SourceID    uint
	Date        time.Time
	Time        *string
	Merchant    string
	Amount      int64
	Notes       string
	IsRecurring int8
}

type ListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Summary      Summary              `json:"summary"`
}

type Summary struct {
	TotalIn  int64 `json:"total_in"`
	TotalOut int64 `json:"total_out"`
	Count    int   `json:"count"`
}

type ImportError struct {
	Row    int    `json:"row"`
	Reason string `json:"reason"`
}

type ImportResult struct {
	Imported          int           `json:"imported"`
	SkippedDuplicates int           `json:"skipped_duplicates"`
	SkippedErrors     int           `json:"skipped_errors"`
	TotalRows         int           `json:"total_rows"`
	Errors            []ImportError `json:"errors"`
}
