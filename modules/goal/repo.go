package goal

import (
	"errors"
	"kayakaga-api/domain/mysql"

	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) ListGoals(userID uint) ([]mysql.Goal, error) {
	var goals []mysql.Goal
	err := r.db.Preload("GoalType").Preload("Milestones").
		Where("user_id = ?", userID).
		Order("id ASC").
		Find(&goals).Error
	return goals, err
}

func (r *repo) GetGoalByID(id, userID uint) (*mysql.Goal, error) {
	var goal mysql.Goal
	err := r.db.Preload("GoalType").Preload("Milestones").
		Where("id = ? AND user_id = ?", id, userID).
		First(&goal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("goal not found")
		}
		return nil, err
	}
	return &goal, nil
}

func (r *repo) CreateGoal(goal *mysql.Goal) error {
	return r.db.Create(goal).Error
}

func (r *repo) UpdateGoal(goal *mysql.Goal) error {
	return r.db.Save(goal).Error
}

func (r *repo) DeleteGoal(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&mysql.Goal{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("goal not found")
	}
	return nil
}
