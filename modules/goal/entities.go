package goal

import "time"

type CreateGoalRequest struct {
	Name                string      `json:"name" binding:"required"`
	GoalTypeID          uint        `json:"goal_type_id" binding:"required"`
	AccountID           uint        `json:"account_id" binding:"required"`
	TargetAmount        int64       `json:"target_amount" binding:"required"`
	CurrentAmount       int64       `json:"current_amount"`
	MonthlyContribution int64       `json:"monthly_contribution"`
	TargetDate          *time.Time  `json:"target_date"`
	Milestones          []Milestone `json:"milestones"`
}

type UpdateGoalRequest struct {
	Name                string
	GoalTypeID          uint
	AccountID           uint
	TargetAmount        int64
	CurrentAmount       int64
	MonthlyContribution int64
	TargetDate          *time.Time
}

type UpdateContributionRequest struct {
	Amount int64 `json:"amount" binding:"required"`
}

type Milestone struct {
	Amount int64 `json:"amount"`
}

type GoalResponse struct {
	ID                  uint       `json:"id"`
	Name                string     `json:"name"`
	GoalTypeID          uint       `json:"goal_type_id"`
	AccountID           uint       `json:"account_id"`
	TargetAmount        int64      `json:"target_amount"`
	CurrentAmount       int64      `json:"current_amount"`
	MonthlyContribution int64      `json:"monthly_contribution"`
	TargetDate          *string    `json:"target_date"`
	ProgressPct         float64    `json:"progress_pct"`
	EtaMonths           int        `json:"eta_months"`
	EtaDate             *string    `json:"eta_date"`
}

type GoalDetailResponse struct {
	Goal       GoalResponse  `json:"goal"`
	Milestones []MilestoneResp `json:"milestones"`
}

type MilestoneResp struct {
	ID        uint    `json:"id"`
	Amount    int64   `json:"amount"`
	ReachedAt *string `json:"reached_at"`
}
