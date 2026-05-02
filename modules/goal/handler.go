package goal

import (
	"kayakaga-api/domain/mysql"
	"kayakaga-api/utils/helper"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type handler struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &handler{repo: repo}
}

func (h *handler) ListGoals(userID uint) ([]GoalResponse, error) {
	goals, err := h.repo.ListGoals(userID)
	if err != nil {
		return nil, err
	}

	resp := make([]GoalResponse, len(goals))
	for i, g := range goals {
		resp[i] = *h.buildResponse(&g)
	}
	return resp, nil
}

func (h *handler) GetGoal(id, userID uint) (*GoalDetailResponse, error) {
	goal, err := h.repo.GetGoalByID(id, userID)
	if err != nil {
		return nil, err
	}

	milestones := make([]MilestoneResp, len(goal.Milestones))
	for i, m := range goal.Milestones {
		var reachedAt *string
		if m.ReachedAt != nil {
			str := m.ReachedAt.Format("2006-01-02")
			reachedAt = &str
		}
		milestones[i] = MilestoneResp{
			ID:        m.ID,
			Amount:    m.Amount,
			ReachedAt: reachedAt,
		}
	}

	return &GoalDetailResponse{
		Goal:       *h.buildResponse(goal),
		Milestones: milestones,
	}, nil
}

func (h *handler) CreateGoal(userID uint, req *CreateGoalRequest) (*GoalResponse, error) {
	goal := &mysql.Goal{
		UserID:              userID,
		GoalTypeID:          req.GoalTypeID,
		AccountID:           req.AccountID,
		Name:                req.Name,
		TargetAmount:        req.TargetAmount,
		CurrentAmount:       req.CurrentAmount,
		MonthlyContribution: req.MonthlyContribution,
		TargetDate:          req.TargetDate,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}

	if err := h.repo.CreateGoal(goal); err != nil {
		return nil, err
	}

	for _, m := range req.Milestones {
		milestone := &mysql.GoalMilestone{
			GoalID:    goal.ID,
			Amount:    m.Amount,
			CreatedAt: time.Now().UTC(),
		}
		h.repo.(*repo).db.Create(milestone)
	}

	return h.buildResponse(goal), nil
}

func (h *handler) UpdateGoal(id, userID uint, req *UpdateGoalRequest) (*GoalResponse, error) {
	goal, err := h.repo.GetGoalByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		goal.Name = req.Name
	}
	if req.GoalTypeID > 0 {
		goal.GoalTypeID = req.GoalTypeID
	}
	if req.AccountID > 0 {
		goal.AccountID = req.AccountID
	}
	if req.TargetAmount > 0 {
		goal.TargetAmount = req.TargetAmount
	}
	if req.CurrentAmount > 0 {
		goal.CurrentAmount = req.CurrentAmount
	}
	if req.MonthlyContribution > 0 {
		goal.MonthlyContribution = req.MonthlyContribution
	}
	if req.TargetDate != nil {
		goal.TargetDate = req.TargetDate
	}
	goal.UpdatedAt = time.Now().UTC()

	if err := h.repo.UpdateGoal(goal); err != nil {
		return nil, err
	}

	return h.buildResponse(goal), nil
}

func (h *handler) DeleteGoal(id, userID uint) error {
	return h.repo.DeleteGoal(id, userID)
}

func (h *handler) UpdateContribution(id, userID uint, req *UpdateContributionRequest) (*GoalResponse, error) {
	goal, err := h.repo.GetGoalByID(id, userID)
	if err != nil {
		return nil, err
	}

	goal.MonthlyContribution = req.Amount
	goal.UpdatedAt = time.Now().UTC()

	if err := h.repo.UpdateGoal(goal); err != nil {
		return nil, err
	}

	return h.buildResponse(goal), nil
}

func (h *handler) buildResponse(goal *mysql.Goal) *GoalResponse {
	progressPct := 0.0
	if goal.TargetAmount > 0 {
		progressPct = (float64(goal.CurrentAmount) / float64(goal.TargetAmount)) * 100
	}

	remaining := goal.TargetAmount - goal.CurrentAmount
	etaMonths := 0
	var etaDate *string

	if goal.MonthlyContribution > 0 && remaining > 0 {
		etaMonths = int(math.Ceil(float64(remaining) / float64(goal.MonthlyContribution)))
		eta := time.Now().UTC().AddDate(0, etaMonths, 0)
		etaStr := eta.Format("2006-01-02")
		etaDate = &etaStr
	}

	var targetDate *string
	if goal.TargetDate != nil {
		str := goal.TargetDate.Format("2006-01-02")
		targetDate = &str
	}

	return &GoalResponse{
		ID:                  goal.ID,
		Name:                goal.Name,
		GoalTypeID:          goal.GoalTypeID,
		AccountID:           goal.AccountID,
		TargetAmount:        goal.TargetAmount,
		CurrentAmount:       goal.CurrentAmount,
		MonthlyContribution: goal.MonthlyContribution,
		TargetDate:          targetDate,
		ProgressPct:         progressPct,
		EtaMonths:           etaMonths,
		EtaDate:             etaDate,
	}
}

func ListGoalsHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		resp, err := uc.ListGoals(userID)
		if err != nil {
			helper.ErrorResponse(c, 500, "INTERNAL_ERROR", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

func GetGoalHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		idParam := c.Param("id")
		id, _ := strconv.ParseUint(idParam, 10, 32)

		resp, err := uc.GetGoal(uint(id), userID)
		if err != nil {
			helper.ErrorResponse(c, 404, "GOAL_NOT_FOUND", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

func CreateGoalHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		var req CreateGoalRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		resp, err := uc.CreateGoal(userID, &req)
		if err != nil {
			helper.ErrorResponse(c, 400, "CREATE_FAILED", err.Error())
			return
		}

		helper.CreatedResponse(c, resp)
	}
}

func UpdateGoalHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		idParam := c.Param("id")
		id, _ := strconv.ParseUint(idParam, 10, 32)

		var req UpdateGoalRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		resp, err := uc.UpdateGoal(uint(id), userID, &req)
		if err != nil {
			helper.ErrorResponse(c, 400, "UPDATE_FAILED", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

func DeleteGoalHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		idParam := c.Param("id")
		id, _ := strconv.ParseUint(idParam, 10, 32)

		if err := uc.DeleteGoal(uint(id), userID); err != nil {
			helper.ErrorResponse(c, 404, "DELETE_FAILED", err.Error())
			return
		}

		helper.SuccessResponse(c, gin.H{"message": "goal deleted successfully"})
	}
}

func UpdateContributionHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		idParam := c.Param("id")
		id, _ := strconv.ParseUint(idParam, 10, 32)

		var req UpdateContributionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		resp, err := uc.UpdateContribution(uint(id), userID, &req)
		if err != nil {
			helper.ErrorResponse(c, 400, "UPDATE_FAILED", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}
