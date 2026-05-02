package simulate

type UseCase interface {
	SimulateInvestment(monthlyAmount, annualReturnPct, years int) (*InvestmentResponse, error)
}
