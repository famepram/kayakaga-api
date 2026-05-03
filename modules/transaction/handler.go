package transaction

import (
	"encoding/csv"
	"errors"
	"fmt"
	"kayakaga-api/domain/mysql"
	"kayakaga-api/utils/helper"
	"strconv"
	"strings"
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
		txnResp[i] = TransactionResponse{
			ID:          t.ID,
			AccountID:   t.AccountID,
			CategoryID:  t.CategoryID,
			SourceID:    t.SourceID,
			Date:        t.Date.Format("2006-01-02"),
			Time:        t.Time,
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

func (h *handler) ImportCSV(userID uint, accountID uint, csvData [][]string) (*ImportResult, error) {
	result := &ImportResult{
		TotalRows: len(csvData) - 1,
		Errors:    []ImportError{},
	}

	if len(csvData) < 2 {
		return nil, errors.New("CSV file is empty or has no data rows")
	}

	format, err := detectFormat(csvData[0])
	if err != nil {
		return nil, err
	}

	validTxns := []mysql.Transaction{}
	dates := []time.Time{}
	amounts := []int64{}
	merchants := []string{}

	for i := 1; i < len(csvData); i++ {
		row := csvData[i]
		if len(row) == 0 || (len(row) == 1 && strings.TrimSpace(row[0]) == "") {
			continue
		}

		txn, err := parseRow(row, format, accountID)
		if err != nil {
			result.SkippedErrors++
			result.Errors = append(result.Errors, ImportError{
				Row:    i + 1,
				Reason: err.Error(),
			})
			continue
		}

		dates = append(dates, txn.Date)
		amounts = append(amounts, txn.Amount)
		merchants = append(merchants, txn.Merchant)
		validTxns = append(validTxns, *txn)
	}

	if len(validTxns) == 0 {
		return result, nil
	}

	duplicates, err := h.repo.CheckDuplicates(userID, accountID, dates, amounts, merchants)
	if err != nil {
		return nil, err
	}

	finalTxns := []mysql.Transaction{}
	for i, txn := range validTxns {
		if duplicates[i] {
			result.SkippedDuplicates++
			continue
		}
		txn.UserID = userID
		finalTxns = append(finalTxns, txn)
	}

	if len(finalTxns) > 0 {
		if err := h.repo.BulkInsert(finalTxns); err != nil {
			return nil, err
		}
		result.Imported = len(finalTxns)
	}

	return result, nil
}

type csvFormat int

const (
	FormatStandard csvFormat = iota
	FormatSimple
)

func detectFormat(headers []string) (csvFormat, error) {
	normalized := make([]string, len(headers))
	for i, h := range headers {
		normalized[i] = strings.ToLower(strings.TrimSpace(h))
	}

	if containsAll(normalized, []string{"tanggal", "keterangan", "debit", "kredit", "saldo"}) {
		return FormatStandard, nil
	}

	if containsAll(normalized, []string{"date", "description", "amount"}) {
		return FormatSimple, nil
	}

	return 0, errors.New("CSV format not recognized. Expected headers: tanggal,keterangan,debit,kredit,saldo OR date,description,amount")
}

func containsAll(row []string, required []string) bool {
	for _, req := range required {
		found := false
		for _, r := range row {
			if r == req {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func parseRow(row []string, format csvFormat, accountID uint) (*mysql.Transaction, error) {
	txn := &mysql.Transaction{
		AccountID:     accountID,
		SourceID:      1,
		IsRecurring:   0,
		AiCategorized: 1,
	}

	var dateStr, merchant string
	var amount int64

	if format == FormatStandard {
		if len(row) < 5 {
			return nil, errors.New("invalid row format")
		}

		dateStr = strings.TrimSpace(row[0])
		merchant = strings.TrimSpace(row[1])

		debitStr := strings.TrimSpace(row[2])
		kreditStr := strings.TrimSpace(row[3])

		if kreditStr != "" {
			kredit, err := strconv.ParseInt(strings.ReplaceAll(kreditStr, ",", ""), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid kredit value: %s", kreditStr)
			}
			amount = kredit
		} else if debitStr != "" {
			debit, err := strconv.ParseInt(strings.ReplaceAll(debitStr, ",", ""), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid debit value: %s", debitStr)
			}
			amount = -debit
		} else {
			return nil, errors.New("both debit and kredit are empty")
		}

		date, err := time.Parse("02/01/2006", dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %s (expected dd/mm/yyyy)", dateStr)
		}
		txn.Date = date

	} else {
		if len(row) < 3 {
			return nil, errors.New("invalid row format")
		}

		dateStr = strings.TrimSpace(row[0])
		merchant = strings.TrimSpace(row[1])
		amountStr := strings.TrimSpace(row[2])

		amt, err := strconv.ParseInt(strings.ReplaceAll(amountStr, ",", ""), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount value: %s", amountStr)
		}
		amount = amt

		date, err := tryParseDate(dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %s", dateStr)
		}
		txn.Date = date
	}

	txn.Merchant = merchant
	txn.Amount = amount
	txn.CategoryID = categorizeTransaction(merchant, amount)

	return txn, nil
}

func tryParseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"02/01/2006",
		"02-01-2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("date format not recognized")
}

func categorizeTransaction(merchant string, amount int64) uint {
	lowerMerchant := strings.ToLower(merchant)

	keywordMap := map[string]uint{
		"gaji":           9,
		"salary":         9,
		"transfer masuk": 9,
		"grab":           2,
		"gojek":          2,
		"ojek":           2,
		"taxi":           2,
		"parkir":         2,
		"gopay":          8,
		"ovo":            8,
		"dana":           8,
		"shopee":         5,
		"netflix":        3,
		"spotify":        3,
		"youtube":        3,
		"steam":          3,
		"game":           3,
		"pln":            4,
		"listrik":        4,
		"pdam":           4,
		"air":            4,
		"internet":       4,
		"wifi":           4,
		"bpjs":           4,
		"telkom":         4,
		"indihome":       4,
		"indomaret":      5,
		"alfamart":       5,
		"supermarket":    5,
		"mall":           5,
		"tokopedia":      5,
		"lazada":         5,
		"apotek":         6,
		"klinik":         6,
		"dokter":         6,
		"rumah sakit":    6,
		"obat":           6,
		"medis":          6,
		"investasi":      7,
		"saham":          7,
		"reksa":          7,
		"tabungan":       7,
	}

	for keyword, categoryID := range keywordMap {
		if strings.Contains(lowerMerchant, keyword) {
			return categoryID
		}
	}

	if amount > 0 {
		return 9
	}

	return 8
}

func (h *handler) buildResponse(txn *mysql.Transaction) *TransactionResponse {
	return &TransactionResponse{
		ID:          txn.ID,
		AccountID:   txn.AccountID,
		CategoryID:  txn.CategoryID,
		SourceID:    txn.SourceID,
		Date:        txn.Date.Format("2006-01-02"),
		Time:        txn.Time,
		Merchant:    txn.Merchant,
		Amount:      txn.Amount,
		Notes:       txn.Notes,
		IsRecurring: txn.IsRecurring,
	}
}

// ListTransactionsHandler godoc
// @Summary List transactions
// @Description Get list of user transactions with filtering and summary
// @Tags Transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param period query string false "Period filter (today, week, month, last_month, year)" Enums(today, week, month, last_month, year)
// @Param account_id query int false "Filter by account ID"
// @Param category_id query int false "Filter by category ID"
// @Param merchant query string false "Filter by merchant name (partial match)"
// @Param is_recurring query int false "Filter by recurring status (0 or 1)"
// @Success 200 {object} helper.Response{data=ListResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /transactions [get]
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

// CreateTransactionHandler godoc
// @Summary Create transaction
// @Description Create a new transaction
// @Tags Transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateTransactionRequest true "Transaction details"
// @Success 201 {object} helper.Response{data=TransactionResponse}
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /transactions [post]
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

// UpdateTransactionHandler godoc
// @Summary Update transaction
// @Description Update existing transaction
// @Tags Transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Transaction ID"
// @Param request body UpdateTransactionRequest true "Transaction updates"
// @Success 200 {object} helper.Response{data=TransactionResponse}
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Router /transactions/{id} [put]
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

// DeleteTransactionHandler godoc
// @Summary Delete transaction
// @Description Delete a transaction
// @Tags Transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Transaction ID"
// @Success 200 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Router /transactions/{id} [delete]
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

// ImportCSVHandler godoc
// @Summary Import transactions from CSV
// @Description Bulk import transactions from CSV file. Supports Indonesian bank format (tanggal,keterangan,debit,kredit,saldo) and simple format (date,description,amount)
// @Tags Transactions
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "CSV file (max 5MB)"
// @Param account_id formData int true "Target Account ID"
// @Success 200 {object} helper.Response{data=ImportResult}
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /transactions/import/csv [post]
func ImportCSVHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		accountIDStr := c.PostForm("account_id")
		if accountIDStr == "" {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", "account_id is required")
			return
		}
		accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
		if err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", "invalid account_id")
			return
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			helper.ErrorResponse(c, 400, "INVALID_FILE", "file is required")
			return
		}

		if fileHeader.Size > 5*1024*1024 {
			helper.ErrorResponse(c, 400, "INVALID_FILE", "File must be a CSV with max size 5MB")
			return
		}

		if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".csv") {
			helper.ErrorResponse(c, 400, "INVALID_FILE", "File must be a CSV")
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			helper.ErrorResponse(c, 500, "FILE_READ_ERROR", "Failed to read file")
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			helper.ErrorResponse(c, 400, "INVALID_CSV", "Failed to parse CSV file")
			return
		}

		result, err := uc.ImportCSV(userID, uint(accountID), records)
		if err != nil {
			if strings.Contains(err.Error(), "format not recognized") {
				helper.ErrorResponse(c, 400, "INVALID_FORMAT", err.Error())
			} else {
				helper.ErrorResponse(c, 500, "IMPORT_FAILED", err.Error())
			}
			return
		}

		helper.SuccessResponse(c, result)
	}
}

// CreateTransactionHandler godoc
