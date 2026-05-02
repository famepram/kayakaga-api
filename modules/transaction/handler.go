package transaction

import (
	"kayakaga-api/domain/mysql"
	"kayakaga-api/utils/helper"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type handler struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &handler{repo: repo}
}

func (h *handler) ListTransactions(userID uint, filters *ListFilters) (*ListResponse, error) {
	txnRecords, err := h.repo.ListTransactions(userID, filters)
	if err != nil {
		return nil, err
	}

	txnResp := make([]TransactionResponse, len(txnRecords))
	totalIn := int64(0)
	totalOut := int64(0)

	for i, t := range txnRecords {
		var timeStr *string
		if t.Time != nil {
			str := t.Time.Format("15:04:05")
			timeStr = &str
		}

		txnResp[i] = TransactionResponse{
			ID:          t.ID,
			AccountID:   t.AccountID,
			CategoryID:  t.CategoryID,
			SourceID:    t.SourceID,
			Date:        t.Date.Format("2006-01-02"),
			Time:        timeStr,
			Merchant:    t.Merchant,
			Amount:      t.Amount,
			Notes:       t.Notes,
			IsRecurring: t.IsRecurring,
		}

		if t.Amount > 0 {
			totalIn += t.Amount
		} else {
			totalOut += t.Amount
		}
	}

	return &ListResponse{
		Transactions: txnResp,
		Summary: Summary{
			TotalIn:  totalIn,
			TotalOut: totalOut,
			Count:    len(txnRecords),
		},
	}, nil
}

func (h *handler) CreateTransaction(userID uint, req *CreateTransactionRequest) (*TransactionResponse, error) {
	txn := &mysql.Transaction{
		UserID:      userID,
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		SourceID:    req.SourceID,
		Date:        req.Date,
		Time:        req.Time,
		Merchant:    req.Merchant,
		Amount:      req.Amount,
		Notes:       req.Notes,
		IsRecurring: req.IsRecurring,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := h.repo.CreateTransaction(txn); err != nil {
		return nil, err
	}

	return h.buildResponse(txn), nil
}

func (h *handler) UpdateTransaction(id, userID uint, req *UpdateTransactionRequest) (*TransactionResponse, error) {
	txn, err := h.repo.GetTransactionByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.AccountID > 0 {
		txn.AccountID = req.AccountID
	}
	if req.CategoryID > 0 {
		txn.CategoryID = req.CategoryID
	}
	if req.SourceID > 0 {
		txn.SourceID = req.SourceID
	}
	if !req.Date.IsZero() {
		txn.Date = req.Date
	}
	if req.Time != nil {
		txn.Time = req.Time
	}
	if req.Merchant != "" {
		txn.Merchant = req.Merchant
	}
	if req.Amount != 0 {
		txn.Amount = req.Amount
	}
	if req.Notes != "" {
		txn.Notes = req.Notes
	}
	txn.IsRecurring = req.IsRecurring
	txn.UpdatedAt = time.Now().UTC()

	if err := h.repo.UpdateTransaction(txn); err != nil {
		return nil, err
	}

	return h.buildResponse(txn), nil
}

func (h *handler) DeleteTransaction(id, userID uint) error {
	return h.repo.DeleteTransaction(id, userID)
}

func (h *handler) buildResponse(txn *mysql.Transaction) *TransactionResponse {
	var timeStr *string
	if txn.Time != nil {
		str := txn.Time.Format("15:04:05")
		timeStr = &str
	}

	return &TransactionResponse{
		ID:          txn.ID,
		AccountID:   txn.AccountID,
		CategoryID:  txn.CategoryID,
		SourceID:    txn.SourceID,
		Date:        txn.Date.Format("2006-01-02"),
		Time:        timeStr,
		Merchant:    txn.Merchant,
		Amount:      txn.Amount,
		Notes:       txn.Notes,
		IsRecurring: txn.IsRecurring,
	}
}

func ListTransactionsHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		filters := &ListFilters{
			Period: c.DefaultQuery("period", "month"),
		}

		if accountID := c.Query("account_id"); accountID != "" {
			id := uint(0)
			if _, err := c.GetQuery("account_id"); err {
				id = c.GetUint("account_id")
			}
			filters.AccountID = &id
		}
		if categoryID := c.Query("category_id"); categoryID != "" {
			id := uint(0)
			if _, err := c.GetQuery("category_id"); err {
				id = c.GetUint("category_id")
			}
			filters.CategoryID = &id
		}
		if merchant := c.Query("merchant"); merchant != "" {
			filters.Merchant = merchant
		}
		if isRecurring := c.Query("is_recurring"); isRecurring != "" {
			val := int8(1)
			if isRecurring == "0" {
				val = 0
			}
			filters.IsRecurring = &val
		}

		resp, err := uc.ListTransactions(userID, filters)
		if err != nil {
			helper.ErrorResponse(c, 500, "INTERNAL_ERROR", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

func CreateTransactionHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		var req CreateTransactionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		if req.SourceID == 0 {
			req.SourceID = 1
		}

		resp, err := uc.CreateTransaction(userID, &req)
		if err != nil {
			helper.ErrorResponse(c, 400, "CREATE_FAILED", err.Error())
			return
		}

		helper.CreatedResponse(c, resp)
	}
}

func UpdateTransactionHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		idParam := c.Param("id")
		id, _ := strconv.ParseUint(idParam, 10, 32)

		var req UpdateTransactionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		resp, err := uc.UpdateTransaction(uint(id), userID, &req)
		if err != nil {
			helper.ErrorResponse(c, 400, "UPDATE_FAILED", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

func DeleteTransactionHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		idParam := c.Param("id")
		id, _ := strconv.ParseUint(idParam, 10, 32)

		if err := uc.DeleteTransaction(uint(id), userID); err != nil {
			helper.ErrorResponse(c, 404, "DELETE_FAILED", err.Error())
			return
		}

		helper.SuccessResponse(c, gin.H{"message": "transaction deleted successfully"})
	}
}
