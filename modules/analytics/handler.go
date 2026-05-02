package analytics

import (
	"kayakaga-api/utils/helper"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type handler struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &handler{repo: repo}
}

func (h *handler) GetBudget(userID uint, period, accountID string) (*BudgetResponse, error) {
	data, err := h.repo.GetBudgetBreakdown(userID, period, accountID)
	if err != nil {
		return nil, err
	}

	savingsRate := 0.0
	if data.Income > 0 {
		savingsRate = ((float64(data.Income) - float64(data.Expenses)) / float64(data.Income)) * 100
	}

	return &BudgetResponse{
		Income:      data.Income,
		Expenses:    data.Expenses,
		SavingsRate: savingsRate,
		Breakdown:   data.Categories,
	}, nil
}

func (h *handler) GetCompare(userID uint, accountID string) (*CompareResponse, error) {
	data, err := h.repo.GetCompareData(userID, accountID)
	if err != nil {
		return nil, err
	}

	return &CompareResponse{
		Comparison: data.ThisMonth,
	}, nil
}

func (h *handler) GetAnomalies(userID uint, period, accountID string) (*AnomaliesResponse, error) {
	anomaliesData, err := h.repo.GetAnomalies(userID, period, accountID)
	if err != nil {
		return nil, err
	}

	anomalies := []Anomaly{}
	for _, a := range anomaliesData {
		severity := "low"
		multiplier := float64(a.Amount) / a.AvgAmount
		if multiplier >= 5 {
			severity = "high"
		} else if multiplier >= 3 {
			severity = "medium"
		}

		anomalies = append(anomalies, Anomaly{
			TransactionID: a.TransactionID,
			Merchant:      a.Merchant,
			Amount:        a.Amount,
			Date:          a.Date,
			Reason:        a.Reason,
			Severity:      severity,
		})
	}

	return &AnomaliesResponse{
		Anomalies: anomalies,
	}, nil
}

func (h *handler) GetRecurring(userID uint, accountID string) (*RecurringResponse, error) {
	data, err := h.repo.GetRecurring(userID, accountID)
	if err != nil {
		return nil, err
	}

	totalMonthly := int64(0)
	for _, item := range data.Items {
		totalMonthly += item.Amount
	}

	return &RecurringResponse{
		Items:        data.Items,
		TotalMonthly: totalMonthly,
	}, nil
}

func (h *handler) GetSavingsSuggestion(userID uint, targetSavings *int64) (*SavingsSuggestionResponse, error) {
	data, err := h.repo.GetSavingsSuggestions(userID, 0)
	if err != nil {
		return nil, err
	}

	impact := "With these changes, you could save money towards your financial goals faster."
	if targetSavings != nil && *targetSavings > 0 {
		months := math.Ceil(float64(*targetSavings-data.TotalPotentialSaving) / float64(data.TotalPotentialSaving))
		impact = "This could help you reach your savings goal " + string(int(months)) + " months faster."
	}

	return &SavingsSuggestionResponse{
		Suggestions:          data.Suggestions,
		TotalPotentialSaving: data.TotalPotentialSaving,
		ImpactOnGoals:        impact,
	}, nil
}

func (h *handler) GetGoalRecommendation(goalID, userID uint, targetMonths *int, newContribution *int64) (*GoalRecommendationResponse, error) {
	data, err := h.repo.GetGoalRecommendation(goalID, userID, targetMonths, newContribution)
	if err != nil {
		return nil, err
	}

	return &GoalRecommendationResponse{
		GoalName:            data.GoalName,
		RemainingAmount:     data.RemainingAmount,
		CurrentContribution: data.CurrentContribution,
		CurrentEtaMonths:    data.CurrentEtaMonths,
		CurrentEtaDate:      data.CurrentEtaDate,
		Scenarios:           data.Scenarios,
	}, nil
}

func GetBudgetHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		period := c.DefaultQuery("period", "month")
		accountID := c.Query("account_id")

		resp, err := uc.GetBudget(userID, period, accountID)
		if err != nil {
			helper.ErrorResponse(c, 500, "INTERNAL_ERROR", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

func GetCompareHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		accountID := c.Query("account_id")

		resp, err := uc.GetCompare(userID, accountID)
		if err != nil {
			helper.ErrorResponse(c, 500, "INTERNAL_ERROR", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

func GetAnomaliesHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		period := c.DefaultQuery("period", "month")
		accountID := c.Query("account_id")

		resp, err := uc.GetAnomalies(userID, period, accountID)
		if err != nil {
			helper.ErrorResponse(c, 500, "INTERNAL_ERROR", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

func GetRecurringHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		accountID := c.Query("account_id")

		resp, err := uc.GetRecurring(userID, accountID)
		if err != nil {
			helper.ErrorResponse(c, 500, "INTERNAL_ERROR", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

func GetSavingsSuggestionHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var targetSavings *int64
		if ts := c.Query("target_savings"); ts != "" {
			val := int64(0)
			targetSavings = &val
		}

		resp, err := uc.GetSavingsSuggestion(userID, targetSavings)
		if err != nil {
			helper.ErrorResponse(c, 500, "INTERNAL_ERROR", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

func GetGoalRecommendationHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		goalIDParam := c.Query("goal_id")
		goalID, _ := strconv.ParseUint(goalIDParam, 10, 32)

		var targetMonths *int
		var newContribution *int64

		if tm := c.Query("target_months"); tm != "" {
			val := 0
			targetMonths = &val
		}
		if nc := c.Query("new_monthly_contribution"); nc != "" {
			val := int64(0)
			newContribution = &val
		}

		resp, err := uc.GetGoalRecommendation(uint(goalID), userID, targetMonths, newContribution)
		if err != nil {
			helper.ErrorResponse(c, 500, "INTERNAL_ERROR", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}
