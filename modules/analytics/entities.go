package analytics

type BudgetResponse struct {
	Income        int64              `json:"income"`
	Expenses      int64              `json:"expenses"`
	SavingsRate   float64            `json:"savings_rate"`
	Breakdown     []CategoryBreakdown `json:"breakdown"`
	VsPrevious    *VsPrevious        `json:"vs_previous,omitempty"`
}

type CategoryBreakdown struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Total        int64   `json:"total"`
	Percentage   float64 `json:"percentage"`
}

type VsPrevious struct {
	IncomeDeltaPct    float64 `json:"income_delta_pct"`
	ExpensesDeltaPct  float64 `json:"expenses_delta_pct"`
}

type CompareResponse struct {
	Comparison []CategoryComparison `json:"comparison"`
}

type CategoryComparison struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	ThisMonth    int64   `json:"this_month"`
	LastMonth    int64   `json:"last_month"`
	DeltaPct     float64 `json:"delta_pct"`
	Trend        string  `json:"trend"`
}

type AnomaliesResponse struct {
	Anomalies []Anomaly `json:"anomalies"`
}

type Anomaly struct {
	TransactionID uint   `json:"transaction_id"`
	Merchant      string `json:"merchant"`
	Amount        int64  `json:"amount"`
	Date          string `json:"date"`
	Reason        string `json:"reason"`
	Severity      string `json:"severity"`
}

type RecurringResponse struct {
	Items        []RecurringItem `json:"items"`
	TotalMonthly int64           `json:"total_monthly"`
}

type RecurringItem struct {
	Merchant      string `json:"merchant"`
	Amount        int64  `json:"amount"`
	AccountID     uint   `json:"account_id"`
	LastCharged   string `json:"last_charged"`
	NextRenewal   string `json:"next_renewal"`
	CategoryName  string `json:"category_name"`
}

type SavingsSuggestionResponse struct {
	Suggestions           []Suggestion           `json:"suggestions"`
	TotalPotentialSaving  int64                  `json:"total_potential_saving"`
	ImpactOnGoals         string                 `json:"impact_on_goals"`
}

type Suggestion struct {
	CategoryID        uint   `json:"category_id"`
	CategoryName      string `json:"category_name"`
	CurrentSpend      int64  `json:"current_spend"`
	SuggestedLimit    int64  `json:"suggested_limit"`
	PotentialSaving   int64  `json:"potential_saving"`
	Reasoning         string `json:"reasoning"`
}

type GoalRecommendationResponse struct {
	GoalName              string                 `json:"goal_name"`
	RemainingAmount       int64                  `json:"remaining_amount"`
	CurrentContribution   int64                  `json:"current_contribution"`
	CurrentEtaMonths      int                    `json:"current_eta_months"`
	CurrentEtaDate        string                 `json:"current_eta_date"`
	Scenarios             []Scenario             `json:"scenarios"`
}

type Scenario struct {
	MonthlyContribution int64  `json:"monthly_contribution"`
	EtaMonths           int    `json:"eta_months"`
	MonthsFaster        int    `json:"months_faster"`
	EtaDate             string `json:"eta_date"`
}

type BudgetData struct {
	Income      int64
	Expenses    int64
	Categories  []CategoryBreakdown
}

type CompareData struct {
	ThisMonth []CategoryComparison
}

type AnomalyData struct {
	TransactionID uint
	Merchant      string
	Amount        int64
	Date          string
	AvgAmount     float64
	IsFirstTime   bool
	IsDuplicate   bool
	Reason        string
}

type RecurringData struct {
	Items        []RecurringItem
}

type SavingsSuggestionData struct {
	Suggestions          []Suggestion
	TotalPotentialSaving int64
}

type GoalRecommendationData struct {
	GoalName            string
	RemainingAmount     int64
	CurrentContribution int64
	CurrentEtaMonths    int
	CurrentEtaDate      string
	Scenarios           []Scenario
}
