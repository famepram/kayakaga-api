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

// GetBudgetHandler godoc
// @Summary Get budget analytics
// @Description Get budget breakdown with income vs expenses and savings rate
// @Tags Analytics
// @Accept json
// @Produce json
// @Security Bearer
// @Param period query string false "Period filter (today, week, month, last_month, year)" Enums(today, week, month, last_month, year)
// @Param account_id query string false "Filter by account ID"
// @Success 200 {object} helper.Response{data=BudgetResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /analytics/budget [get]
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

// GetCompareHandler godoc
// @Summary Get period comparison
// @Description Compare financial metrics between current and previous period
// @Tags Analytics
// @Accept json
// @Produce json
// @Security Bearer
// @Param period query string false "Period filter (month, year)" Enums(month, year)
// @Success 200 {object} helper.Response{data=CompareResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /analytics/compare [get]
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

// GetAnomaliesHandler godoc
// @Summary Get spending anomalies
// @Description Detect unusual spending patterns or large transactions
// @Tags Analytics
// @Accept json
// @Produce json
// @Security Bearer
// @Param period query string false "Period filter (week, month, year)" Enums(week, month, year)
// @Param threshold query number false "Threshold amount (default: 1000000)"
// @Success 200 {object} helper.Response{data=[]Anomaly}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /analytics/anomalies [get]
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

// GetRecurringHandler godoc
// @Summary Get recurring transactions
// @Description Get list of recurring transactions
// @Tags Analytics
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response{data=[]RecurringResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /analytics/recurring [get]
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

// GetSavingsSuggestionHandler godoc
// @Summary Get savings suggestions
// @Description Get AI-powered savings suggestions based on spending patterns
// @Tags Analytics
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response{data=SavingsSuggestionResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /analytics/savings-suggestions [get]
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
