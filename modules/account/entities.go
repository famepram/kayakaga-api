package account

type CreateAccountRequest struct {
	Name          string `json:"name" binding:"required"`
	AccountTypeID uint   `json:"account_type_id" binding:"required"`
	Balance       int64  `json:"balance"`
	Color         string `json:"color"`
	IsPrimary     int8   `json:"is_primary"`
}

type UpdateAccountRequest struct {
	Name           string `json:"name"`
	AccountTypeID  uint   `json:"account_type_id"`
	Balance        int64  `json:"balance"`
	BalanceUpdated bool   `json:"balance_updated"`
	Color          string `json:"color"`
	IsPrimary      int8   `json:"is_primary"`
}

type AccountResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	AccountTypeID uint   `json:"account_type_id"`
	Balance       int64  `json:"balance"`
	Color         string `json:"color"`
	IsPrimary     int8   `json:"is_primary"`
}

type AccountBalance struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	AccountTypeID uint   `json:"account_type_id"`
	Balance       int64  `json:"balance"`
	Color         string `json:"color"`
	IsPrimary     int8   `json:"is_primary"`
}

type BalancesResponse struct {
	Accounts []AccountBalance `json:"accounts"`
	Total    int64            `json:"total"`
}
