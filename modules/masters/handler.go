package masters

import (
	"kayakaga-api/utils/helper"

	"github.com/gin-gonic/gin"
)

type handler struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &handler{repo: repo}
}

func (h *handler) GetAccountTypes() ([]interface{}, error) {
	return h.repo.GetAccountTypes()
}

func (h *handler) GetCategories() ([]interface{}, error) {
	return h.repo.GetCategories()
}

func (h *handler) GetGoalTypes() ([]interface{}, error) {
	return h.repo.GetGoalTypes()
}

func (h *handler) GetRiskProfiles() ([]interface{}, error) {
	return h.repo.GetRiskProfiles()
}

func (h *handler) GetDependents() ([]interface{}, error) {
	return h.repo.GetDependents()
}

func setCacheHeader(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=3600")
}

// GetAccountTypesHandler godoc
// @Summary Get account types
// @Description Get list of active account types (Bank, E-Wallet, Investment)
// @Tags Master Data
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response{data=[]AccountTypeResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /masters/account-types [get]
func GetAccountTypesHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		setCacheHeader(c)

		data, err := uc.GetAccountTypes()
		if err != nil {
			helper.ErrorResponse(c, 500, "FETCH_FAILED", "Failed to fetch account types")
			return
		}

		helper.SuccessResponse(c, data)
	}
}

// GetCategoriesHandler godoc
// @Summary Get transaction categories
// @Description Get list of active transaction categories (Food, Transport, Entertainment, etc.)
// @Tags Master Data
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response{data=[]CategoryResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /masters/categories [get]
func GetCategoriesHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		setCacheHeader(c)

		data, err := uc.GetCategories()
		if err != nil {
			helper.ErrorResponse(c, 500, "FETCH_FAILED", "Failed to fetch categories")
			return
		}

		helper.SuccessResponse(c, data)
	}
}

// GetGoalTypesHandler godoc
// @Summary Get goal types
// @Description Get list of active goal types (House Down Payment, Car, Emergency Fund, etc.)
// @Tags Master Data
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response{data=[]GoalTypeResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /masters/goal-types [get]
func GetGoalTypesHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		setCacheHeader(c)

		data, err := uc.GetGoalTypes()
		if err != nil {
			helper.ErrorResponse(c, 500, "FETCH_FAILED", "Failed to fetch goal types")
			return
		}

		helper.SuccessResponse(c, data)
	}
}

// GetRiskProfilesHandler godoc
// @Summary Get risk profiles
// @Description Get list of active risk profiles (Conservative, Moderate, Aggressive, etc.)
// @Tags Master Data
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response{data=[]RiskProfileResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /masters/risk-profiles [get]
func GetRiskProfilesHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		setCacheHeader(c)

		data, err := uc.GetRiskProfiles()
		if err != nil {
			helper.ErrorResponse(c, 500, "FETCH_FAILED", "Failed to fetch risk profiles")
			return
		}

		helper.SuccessResponse(c, data)
	}
}

// GetDependentsHandler godoc
// @Summary Get dependents types
// @Description Get list of active dependent types (Single, Married with 1 child, etc.)
// @Tags Master Data
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response{data=[]DependentResponse}
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /masters/dependents [get]
func GetDependentsHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		setCacheHeader(c)

		data, err := uc.GetDependents()
		if err != nil {
			helper.ErrorResponse(c, 500, "FETCH_FAILED", "Failed to fetch dependents")
			return
		}

		helper.SuccessResponse(c, data)
	}
}
