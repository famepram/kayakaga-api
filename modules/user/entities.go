package user

type UpdateProfileRequest struct {
	City                   string `json:"city"`
	Profession             string `json:"profession"`
	DependentID            uint   `json:"dependent_id"`
	MonthlyIncome          int64  `json:"monthly_income"`
	MonthlyExpensesEstimate int64  `json:"monthly_expenses_estimate"`
	CurrentSavings         int64  `json:"current_savings"`
	RiskProfileID          uint   `json:"risk_profile_id"`
}

type ProfileResponse struct {
	ID                     uint   `json:"id"`
	Name                   string `json:"name"`
	Email                  string `json:"email"`
	City                   string `json:"city"`
	Profession             string `json:"profession"`
	DependentID            uint   `json:"dependent_id"`
	MonthlyIncome          int64  `json:"monthly_income"`
	MonthlyExpensesEstimate int64  `json:"monthly_expenses_estimate"`
	CurrentSavings         int64  `json:"current_savings"`
	RiskProfileID          uint   `json:"risk_profile_id"`
	Currency               string `json:"currency"`
	CreatedAt              string `json:"created_at"`
}
