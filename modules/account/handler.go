package account

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

func (h *handler) ListAccounts(userID uint) ([]AccountResponse, error) {
	accounts, err := h.repo.ListAccounts(userID)
	if err != nil {
		return nil, err
	}

	resp := make([]AccountResponse, len(accounts))
	for i, a := range accounts {
		resp[i] = AccountResponse{
			ID:            a.ID,
			Name:          a.Name,
			AccountTypeID: a.AccountTypeID,
			Balance:       a.Balance,
			Color:         a.Color,
			IsPrimary:     a.IsPrimary,
		}
	}
	return resp, nil
}

func (h *handler) GetBalances(userID uint) (*BalancesResponse, error) {
	accounts, err := h.repo.ListAccounts(userID)
	if err != nil {
		return nil, err
	}

	balances := make([]AccountBalance, len(accounts))
	total := int64(0)
	for i, a := range accounts {
		balances[i] = AccountBalance{
			ID:            a.ID,
			Name:          a.Name,
			AccountTypeID: a.AccountTypeID,
			Balance:       a.Balance,
			Color:         a.Color,
			IsPrimary:     a.IsPrimary,
		}
		total += a.Balance
	}

	return &BalancesResponse{
		Accounts: balances,
		Total:    total,
	}, nil
}

func (h *handler) CreateAccount(userID uint, req *CreateAccountRequest) (*AccountResponse, error) {
	account := &mysql.Account{
		UserID:        userID,
		AccountTypeID: req.AccountTypeID,
		Name:          req.Name,
		Balance:       req.Balance,
		Color:         req.Color,
		IsPrimary:     req.IsPrimary,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := h.repo.CreateAccount(account); err != nil {
		return nil, err
	}

	return &AccountResponse{
		ID:            account.ID,
		Name:          account.Name,
		AccountTypeID: account.AccountTypeID,
		Balance:       account.Balance,
		Color:         account.Color,
		IsPrimary:     account.IsPrimary,
	}, nil
}

func (h *handler) UpdateAccount(id, userID uint, req *UpdateAccountRequest) (*AccountResponse, error) {
	account, err := h.repo.GetAccountByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		account.Name = req.Name
	}
	if req.AccountTypeID > 0 {
		account.AccountTypeID = req.AccountTypeID
	}
	if req.BalanceUpdated {
		account.Balance = req.Balance
	}
	if req.Color != "" {
		account.Color = req.Color
	}
	account.IsPrimary = req.IsPrimary
	account.UpdatedAt = time.Now().UTC()

	if err := h.repo.UpdateAccount(account); err != nil {
		return nil, err
	}

	return &AccountResponse{
		ID:            account.ID,
		Name:          account.Name,
		AccountTypeID: account.AccountTypeID,
		Balance:       account.Balance,
		Color:         account.Color,
		IsPrimary:     account.IsPrimary,
	}, nil
}

func (h *handler) DeleteAccount(id, userID uint) error {
	return h.repo.DeleteAccount(id, userID)
}

// ListAccountsHandler godoc
// @Summary List accounts
// @Description Get list of user accounts
// @Tags Accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response{data=[]AccountResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /accounts [get]
func ListAccountsHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		resp, err := uc.ListAccounts(userID)
		if err != nil {
			helper.ErrorResponse(c, 500, "INTERNAL_ERROR", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

// GetBalancesHandler godoc
// @Summary Get account balances
// @Description Get balances for all accounts with grand total
// @Tags Accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response{data=BalancesResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /accounts/balances [get]
// GetBalancesHandler godoc
// @Summary Get account balances
// @Description Get balances for all accounts with grand total
// @Tags Accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response{data=BalancesResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /accounts/balances [get]
func GetBalancesHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		resp, err := uc.GetBalances(userID)
		if err != nil {
			helper.ErrorResponse(c, 500, "INTERNAL_ERROR", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

// CreateAccountHandler godoc
// @Summary Create new account
// @Description Create a new account for the user
// @Tags Accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateAccountRequest true "Account details"
// @Success 201 {object} helper.Response{data=AccountResponse}
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Router /accounts [post]
// CreateAccountHandler godoc
// @Summary Create account
// @Description Create a new account
// @Tags Accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateAccountRequest true "Account details"
// @Success 201 {object} helper.Response{data=AccountResponse}
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /accounts [post]
func CreateAccountHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		var req CreateAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		resp, err := uc.CreateAccount(userID, &req)
		if err != nil {
			helper.ErrorResponse(c, 400, "CREATE_FAILED", err.Error())
			return
		}

		helper.CreatedResponse(c, resp)
	}
}

// UpdateAccountHandler godoc
// @Summary Update account
// @Description Update existing account
// @Tags Accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Account ID"
// @Param request body UpdateAccountRequest true "Account updates"
// @Success 200 {object} helper.Response{data=AccountResponse}
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Router /accounts/{id} [put]
func UpdateAccountHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		idParam := c.Param("id")
		id, _ := strconv.ParseUint(idParam, 10, 32)

		var req UpdateAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		resp, err := uc.UpdateAccount(uint(id), userID, &req)
		if err != nil {
			helper.ErrorResponse(c, 400, "UPDATE_FAILED", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

// DeleteAccountHandler godoc
// @Summary Delete account
// @Description Delete an account
// @Tags Accounts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Account ID"
// @Success 200 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Router /accounts/{id} [delete]
func DeleteAccountHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		idParam := c.Param("id")
		id, _ := strconv.ParseUint(idParam, 10, 32)

		if err := uc.DeleteAccount(uint(id), userID); err != nil {
			helper.ErrorResponse(c, 404, "DELETE_FAILED", err.Error())
			return
		}

		helper.SuccessResponse(c, gin.H{"message": "account deleted successfully"})
	}
}
