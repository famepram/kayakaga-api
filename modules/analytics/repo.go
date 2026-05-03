package analytics

import (
	"kayakaga-api/domain/mysql"
	"time"

	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) GetBudgetBreakdown(userID uint, period, accountID string) (*BudgetData, error) {
	query := r.db.Table("transactions t").
		Select("t.category_id, c.name as category_name, SUM(t.amount) as total, COUNT(*) as count").
		Joins("LEFT JOIN m_transaction_categories c ON t.category_id = c.id").
		Where("t.user_id = ?", userID)

	query = r.applyPeriodFilter(query, period)
	if accountID != "" {
		query = query.Where("t.account_id = ?", accountID)
	}

	var results []struct {
		CategoryID   uint
		CategoryName string
		Total        int64
		Count        int
	}

	if err := query.Group("t.category_id, c.name").Scan(&results).Error; err != nil {
		return nil, err
	}

	income := int64(0)
	expenses := int64(0)
	breakdown := []CategoryBreakdown{}

	for _, r := range results {
		total := r.Total
		if total > 0 {
			income += total
		} else {
			expenses += total
		}
	}

	for _, r := range results {
		total := r.Total
		var percentage float64
		if expenses < 0 && total < 0 {
			percentage = (float64(total) / float64(expenses)) * 100
		} else if income > 0 && total > 0 {
			percentage = (float64(total) / float64(income)) * 100
		}

		breakdown = append(breakdown, CategoryBreakdown{
			CategoryID:   r.CategoryID,
			CategoryName: r.CategoryName,
			Total:        total,
			Percentage:   percentage,
		})
	}

	return &BudgetData{
		Income:     income,
		Expenses:   expenses,
		Categories: breakdown,
	}, nil
}

func (r *repo) GetCompareData(userID uint, accountID string) (*CompareData, error) {
	now := time.Now().UTC()
	thisMonth := r.getTransactionByMonth(userID, accountID, now.Year(), int(now.Month()))
	lastMonth := now.AddDate(0, -1, 0)
	lastMonthData := r.getTransactionByMonth(userID, accountID, lastMonth.Year(), int(lastMonth.Month()))

	comparison := []CategoryComparison{}
	categoryMap := make(map[uint]*CategoryComparison)

	for _, cat := range thisMonth {
		categoryMap[cat.CategoryID] = &CategoryComparison{
			CategoryID:   cat.CategoryID,
			CategoryName: cat.CategoryName,
			ThisMonth:    cat.Total,
			LastMonth:    0,
		}
	}

	for _, cat := range lastMonthData {
		if cc, ok := categoryMap[cat.CategoryID]; ok {
			cc.LastMonth = cat.Total
		} else {
			categoryMap[cat.CategoryID] = &CategoryComparison{
				CategoryID:   cat.CategoryID,
				CategoryName: cat.CategoryName,
				ThisMonth:    0,
				LastMonth:    cat.Total,
			}
		}
	}

	for _, cc := range categoryMap {
		if cc.ThisMonth > 0 && cc.LastMonth == 0 {
			cc.Trend = "new"
		} else if cc.ThisMonth == 0 && cc.LastMonth > 0 {
			cc.Trend = "gone"
		} else if cc.ThisMonth > cc.LastMonth {
			cc.Trend = "up"
		} else if cc.ThisMonth < cc.LastMonth {
			cc.Trend = "down"
		} else {
			cc.Trend = "stable"
		}

		if cc.LastMonth != 0 {
			cc.DeltaPct = ((float64(cc.ThisMonth) - float64(cc.LastMonth)) / float64(cc.LastMonth)) * 100
		}

		comparison = append(comparison, *cc)
	}

	return &CompareData{
		ThisMonth: comparison,
	}, nil
}

func (r *repo) GetAnomalies(userID uint, period, accountID string) ([]AnomalyData, error) {
	query := r.db.Table("transactions t").
		Select("t.id, t.merchant, t.amount, t.date, c.name as category_name").
		Joins("LEFT JOIN m_transaction_categories c ON t.category_id = c.id").
		Where("t.user_id = ? AND t.amount < 0 AND t.category_id != 9", userID)

	query = r.applyPeriodFilter(query, period)
	if accountID != "" {
		query = query.Where("t.account_id = ?", accountID)
	}

	var transactions []struct {
		ID           uint
		Merchant     string
		Amount       int64
		Date         time.Time
		CategoryName string
	}

	if err := query.Scan(&transactions).Error; err != nil {
		return nil, err
	}

	categoryAvg := make(map[string][]int64)
	for _, t := range transactions {
		categoryAvg[t.CategoryName] = append(categoryAvg[t.CategoryName], t.Amount)
	}

	anomalies := []AnomalyData{}
	merchantDates := make(map[string]time.Time)

	for _, t := range transactions {
		amounts := categoryAvg[t.CategoryName]
		sum := int64(0)
		for _, a := range amounts {
			sum += a
		}
		avg := float64(sum) / float64(len(amounts))

		data := AnomalyData{
			TransactionID: t.ID,
			Merchant:      t.Merchant,
			Amount:        t.Amount,
			Date:          t.Date.Format("2006-01-02"),
			AvgAmount:     avg,
		}

		multiplier := float64(t.Amount) / avg
		if multiplier >= 5 {
			data.Reason = "Amount 5x+ category average"
			anomalies = append(anomalies, data)
		} else if multiplier >= 3 {
			data.Reason = "Amount 3-5x category average"
			anomalies = append(anomalies, data)
		} else if multiplier >= 2 {
			data.Reason = "Amount 2x category average"
			anomalies = append(anomalies, data)
		}

		if t.Amount > 100000 {
			if lastDate, ok := merchantDates[t.Merchant]; !ok {
				data.Reason = "First time merchant with high amount"
				anomalies = append(anomalies, data)
			} else {
				if t.Date.Sub(lastDate) < 24*time.Hour && t.Amount == int64(avg) {
					data.Reason = "Duplicate transaction (same merchant + amount within 24h)"
					anomalies = append(anomalies, data)
				}
			}
		}

		merchantDates[t.Merchant] = t.Date
	}

	return anomalies, nil
}

func (r *repo) GetRecurring(userID uint, accountID string) (*RecurringData, error) {
	query := r.db.Table("transactions t").
		Select("t.merchant, t.amount, t.account_id, t.date, c.name as category_name").
		Joins("LEFT JOIN m_transaction_categories c ON t.category_id = c.id").
		Where("t.user_id = ? AND t.is_recurring = 1", userID)

	if accountID != "" {
		query = query.Where("t.account_id = ?", accountID)
	}

	var results []struct {
		Merchant     string
		Amount       int64
		AccountID    uint
		Date         time.Time
		CategoryName string
	}

	if err := query.Order("date DESC").Scan(&results).Error; err != nil {
		return nil, err
	}

	items := []RecurringItem{}
	total := int64(0)

	for _, r := range results {
		lastCharged := r.Date.Format("2006-01-02")
		nextRenewal := r.Date.AddDate(0, 1, 0).Format("2006-01-02")

		items = append(items, RecurringItem{
			Merchant:     r.Merchant,
			Amount:       r.Amount,
			AccountID:    r.AccountID,
			LastCharged:  lastCharged,
			NextRenewal:  nextRenewal,
			CategoryName: r.CategoryName,
		})
		total += r.Amount
	}

	return &RecurringData{
		Items: items,
	}, nil
}

func (r *repo) GetSavingsSuggestions(userID uint, targetSavings int64) (*SavingsSuggestionData, error) {
	now := time.Now().UTC()
	query := r.db.Table("transactions t").
		Select("t.category_id, c.name as category_name, SUM(ABS(t.amount)) as total").
		Joins("LEFT JOIN m_transaction_categories c ON t.category_id = c.id").
		Where("t.user_id = ? AND t.amount < 0", userID).
		Where("YEAR(t.date) = ? AND MONTH(t.date) = ?", now.Year(), int(now.Month())).
		Group("t.category_id, c.name").
		Order("total DESC")

	var results []struct {
		CategoryID   uint
		CategoryName string
		Total        int64
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	suggestions := []Suggestion{}
	totalPotential := int64(0)

	for _, r := range results {
		suggestedLimit := int64(float64(r.Total) * 0.8)
		potential := r.Total - suggestedLimit

		suggestions = append(suggestions, Suggestion{
			CategoryID:     r.CategoryID,
			CategoryName:   r.CategoryName,
			CurrentSpend:   r.Total,
			SuggestedLimit: suggestedLimit,
			PotentialSaving: potential,
			Reasoning:      "Reduce spending by 20% in this category",
		})
		totalPotential += potential
	}

	return &SavingsSuggestionData{
		Suggestions:          suggestions,
		TotalPotentialSaving: totalPotential,
	}, nil
}

func (r *repo) GetGoalRecommendation(goalID, userID uint, targetMonths *int, newContribution *int64) (*GoalRecommendationData, error) {
	var goal mysql.Goal
	if err := r.db.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		return nil, err
	}

	remaining := goal.TargetAmount - goal.CurrentAmount
	currentEta := 0
	if goal.MonthlyContribution > 0 {
		currentEta = int((remaining + goal.MonthlyContribution - 1) / goal.MonthlyContribution)
	}

	currentEtaDate := time.Now().UTC().AddDate(0, currentEta, 0).Format("2006-01-02")

	// Generate 3 automatic scenarios
	scenarios := []Scenario{}

	// Scenario 1: Current contribution
	scenarios = append(scenarios, Scenario{
		MonthlyContribution: goal.MonthlyContribution,
		EtaMonths:           currentEta,
		MonthsFaster:        0,
		EtaDate:             currentEtaDate,
	})

	// Scenario 2: Current + 500,000
	if goal.MonthlyContribution > 0 {
		increasedAmount := goal.MonthlyContribution + 500000
		newEta := int((remaining + increasedAmount - 1) / increasedAmount)
		newEtaDate := time.Now().UTC().AddDate(0, newEta, 0)
		scenarios = append(scenarios, Scenario{
			MonthlyContribution: increasedAmount,
			EtaMonths:           newEta,
			MonthsFaster:        currentEta - newEta,
			EtaDate:             newEtaDate.Format("2006-01-02"),
		})
	}

	// Scenario 3: Current + 1,000,000
	if goal.MonthlyContribution > 0 {
		increasedAmount := goal.MonthlyContribution + 1000000
		newEta := int((remaining + increasedAmount - 1) / increasedAmount)
		newEtaDate := time.Now().UTC().AddDate(0, newEta, 0)
		scenarios = append(scenarios, Scenario{
			MonthlyContribution: increasedAmount,
			EtaMonths:           newEta,
			MonthsFaster:        currentEta - newEta,
			EtaDate:             newEtaDate.Format("2006-01-02"),
		})
	}

	return &GoalRecommendationData{
		GoalName:            goal.Name,
		RemainingAmount:     remaining,
		CurrentContribution: goal.MonthlyContribution,
		CurrentEtaMonths:    currentEta,
		CurrentEtaDate:      currentEtaDate,
		Scenarios:           scenarios,
	}, nil
}

func (r *repo) getTransactionByMonth(userID uint, accountID string, year, month int) []struct {
	CategoryID   uint
	CategoryName string
	Total        int64
} {
	query := r.db.Table("transactions t").
		Select("t.category_id, c.name as category_name, SUM(t.amount) as total").
		Joins("LEFT JOIN m_transaction_categories c ON t.category_id = c.id").
		Where("t.user_id = ? AND YEAR(t.date) = ? AND MONTH(t.date) = ?", userID, year, month)

	if accountID != "" {
		query = query.Where("t.account_id = ?", accountID)
	}

	var results []struct {
		CategoryID   uint
		CategoryName string
		Total        int64
	}
	query.Group("t.category_id, c.name").Scan(&results)
	return results
}

func (r *repo) applyPeriodFilter(query *gorm.DB, period string) *gorm.DB {
	now := time.Now().UTC()

	switch period {
	case "month":
		return query.Where("YEAR(t.date) = ? AND MONTH(t.date) = ?", now.Year(), int(now.Month()))
	case "last_month":
		lastMonth := now.AddDate(0, -1, 0)
		return query.Where("YEAR(t.date) = ? AND MONTH(t.date) = ?", lastMonth.Year(), int(lastMonth.Month()))
	case "year":
		return query.Where("YEAR(t.date) = ?", now.Year())
	default:
		return query.Where("YEAR(t.date) = ? AND MONTH(t.date) = ?", now.Year(), int(now.Month()))
	}
}
