package simulate

import (
	"errors"
	"kayakaga-api/utils/helper"
	"math"

	"github.com/gin-gonic/gin"
)

type handler struct{}

func NewUseCase() UseCase {
	return &handler{}
}

func (h *handler) SimulateInvestment(monthlyAmount, annualReturnPct, years int) (*InvestmentResponse, error) {
	if monthlyAmount <= 0 || annualReturnPct <= 0 || years <= 0 {
		return nil, errors.New("all parameters must be positive")
	}

	monthlyRate := float64(annualReturnPct) / 100 / 12
	n := years * 12

	futureValue := float64(monthlyAmount) * ((math.Pow(1+monthlyRate, float64(n)) - 1) / monthlyRate)
	totalInvested := int64(monthlyAmount * n)
	profit := int64(futureValue) - totalInvested
	roiPct := (float64(profit) / float64(totalInvested)) * 100

	breakdown := []MonthlyBreakdown{}
	total := int64(0)

	for month := 1; month <= n; month++ {
		total += int64(monthlyAmount)
		interest := int64(float64(total) * monthlyRate)
		total += interest

		if month%12 == 0 || month == n {
			breakdown = append(breakdown, MonthlyBreakdown{
				Month: month,
				Total: total,
			})
		}
	}

	return &InvestmentResponse{
		FutureValue:      int64(futureValue),
		TotalInvested:    totalInvested,
		Profit:           profit,
		RoiPct:           roiPct,
		MonthlyBreakdown: breakdown,
	}, nil
}

// SimulateInvestmentHandler godoc
// @Summary Simulate investment growth
// @Description Calculate future value of monthly investment with compound interest
// @Tags Simulation
// @Accept json
// @Produce json
// @Security Bearer
// @Param monthly_amount query int true "Monthly investment amount"
// @Param annual_return_pct query int true "Annual return percentage"
// @Param years query int true "Investment duration in years"
// @Success 200 {object} helper.Response{data=InvestmentResponse}
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /simulate/investment [get]
func SimulateInvestmentHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		monthlyAmount := c.Query("monthly_amount")
		annualReturnPct := c.Query("annual_return_pct")
		years := c.Query("years")

		var ma, arp, y int
		if _, err := c.GetQuery("monthly_amount"); err && monthlyAmount != "" {
			ma = c.GetInt("monthly_amount")
		}
		if _, err := c.GetQuery("annual_return_pct"); err && annualReturnPct != "" {
			arp = c.GetInt("annual_return_pct")
		}
		if _, err := c.GetQuery("years"); err && years != "" {
			y = c.GetInt("years")
		}

		if ma <= 0 || arp <= 0 || y <= 0 {
			helper.ErrorResponse(c, 400, "INVALID_INPUT", "All parameters must be positive integers")
			return
		}

		resp, err := uc.SimulateInvestment(ma, arp, y)
		if err != nil {
			helper.ErrorResponse(c, 500, "SIMULATION_FAILED", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}
