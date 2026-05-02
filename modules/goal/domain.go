package goal

import (
	"kayakaga-api/domain/mysql"
)

type Repository interface {
	ListGoals(userID uint) ([]mysql.Goal, error)
	GetGoalByID(id, userID uint) (*mysql.Goal, error)
	CreateGoal(goal *mysql.Goal) error
	UpdateGoal(goal *mysql.Goal) error
	DeleteGoal(id, userID uint) error
}

type UseCase interface {
	ListGoals(userID uint) ([]GoalResponse, error)
	GetGoal(id, userID uint) (*GoalDetailResponse, error)
	CreateGoal(userID uint, req *CreateGoalRequest) (*GoalResponse, error)
	UpdateGoal(id, userID uint, req *UpdateGoalRequest) (*GoalResponse, error)
	DeleteGoal(id, userID uint) error
	UpdateContribution(id, userID uint, req *UpdateContributionRequest) (*GoalResponse, error)
}
