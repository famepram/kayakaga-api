package simulate

type InvestmentResponse struct {
	FutureValue       int64            `json:"future_value"`
	TotalInvested     int64            `json:"total_invested"`
	Profit            int64            `json:"profit"`
	RoiPct            float64          `json:"roi_pct"`
	MonthlyBreakdown  []MonthlyBreakdown `json:"monthly_breakdown"`
}

type MonthlyBreakdown struct {
	Month int   `json:"month"`
	Total int64 `json:"total"`
}
