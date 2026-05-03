package simulate

import (
	"errors"
	"kayakaga-api/utils/helper"
	"math"
	"strconv"

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
		monthlyAmountStr := c.Query("monthly_amount")
		annualReturnPctStr := c.Query("annual_return_pct")
		yearsStr := c.Query("years")

		ma, err := strconv.Atoi(monthlyAmountStr)
		if err != nil || ma <= 0 {
			helper.ErrorResponse(c, 400, "INVALID_INPUT", "monthly_amount must be a positive integer")
			return
		}

		arp, err := strconv.Atoi(annualReturnPctStr)
		if err != nil || arp <= 0 {
			helper.ErrorResponse(c, 400, "INVALID_INPUT", "annual_return_pct must be a positive integer")
			return
		}

		y, err := strconv.Atoi(yearsStr)
		if err != nil || y <= 0 {
			helper.ErrorResponse(c, 400, "INVALID_INPUT", "years must be a positive integer")
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
