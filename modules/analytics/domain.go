package analytics

type Repository interface {
	GetBudgetBreakdown(userID uint, period, accountID string) (*BudgetData, error)
	GetCompareData(userID uint, accountID string) (*CompareData, error)
	GetAnomalies(userID uint, period, accountID string) ([]AnomalyData, error)
	GetRecurring(userID uint, accountID string) (*RecurringData, error)
	GetSavingsSuggestions(userID uint, targetSavings int64) (*SavingsSuggestionData, error)
	GetGoalRecommendation(goalID, userID uint, targetMonths *int, newContribution *int64) (*GoalRecommendationData, error)
}

type UseCase interface {
	GetBudget(userID uint, period, accountID string) (*BudgetResponse, error)
	GetCompare(userID uint, accountID string) (*CompareResponse, error)
	GetAnomalies(userID uint, period, accountID string) (*AnomaliesResponse, error)
	GetRecurring(userID uint, accountID string) (*RecurringResponse, error)
	GetSavingsSuggestion(userID uint, targetSavings *int64) (*SavingsSuggestionResponse, error)
	GetGoalRecommendation(goalID, userID uint, targetMonths *int, newContribution *int64) (*GoalRecommendationResponse, error)
}
