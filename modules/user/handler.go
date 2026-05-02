package user

import (
	"kayakaga-api/utils/helper"
	"time"

	"github.com/gin-gonic/gin"
)

type handler struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &handler{repo: repo}
}

func (h *handler) GetProfile(userID uint) (*ProfileResponse, error) {
	user, err := h.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &ProfileResponse{
		ID:                     user.ID,
		Name:                   user.Name,
		City:                   user.City,
		Profession:             user.Profession,
		DependentID:            user.DependentID,
		MonthlyIncome:          user.MonthlyIncome,
		MonthlyExpensesEstimate: user.MonthlyExpensesEstimate,
		CurrentSavings:         user.CurrentSavings,
		RiskProfileID:          user.RiskProfileID,
		Currency:               user.Currency,
		CreatedAt:              user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (h *handler) UpdateProfile(userID uint, req *UpdateProfileRequest) (*ProfileResponse, error) {
	user, err := h.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if req.City != "" {
		user.City = req.City
	}
	if req.Profession != "" {
		user.Profession = req.Profession
	}
	if req.DependentID > 0 {
		user.DependentID = req.DependentID
	}
	if req.MonthlyIncome > 0 {
		user.MonthlyIncome = req.MonthlyIncome
	}
	if req.MonthlyExpensesEstimate > 0 {
		user.MonthlyExpensesEstimate = req.MonthlyExpensesEstimate
	}
	if req.CurrentSavings > 0 {
		user.CurrentSavings = req.CurrentSavings
	}
	if req.RiskProfileID > 0 {
		user.RiskProfileID = req.RiskProfileID
	}
	user.UpdatedAt = time.Now().UTC()

	if err := h.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	return h.GetProfile(userID)
}

// GetProfileHandler godoc
// @Summary Get user profile
// @Description Get current user profile information
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} helper.Response{data=ProfileResponse}
// @Failure 401 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Router /users/profile [get]
func GetProfileHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		resp, err := uc.GetProfile(userID)
		if err != nil {
			helper.ErrorResponse(c, 404, "USER_NOT_FOUND", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

// UpdateProfileHandler godoc
// @Summary Update user profile
// @Description Update current user profile information
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body UpdateProfileRequest true "Profile updates"
// @Success 200 {object} helper.Response{data=ProfileResponse}
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Router /users/profile [put]
func UpdateProfileHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		var req UpdateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		resp, err := uc.UpdateProfile(userID, &req)
		if err != nil {
			helper.ErrorResponse(c, 400, "UPDATE_FAILED", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}
