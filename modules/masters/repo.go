package masters

import (
	"kayakaga-api/domain/mysql"

	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) GetAccountTypes() ([]interface{}, error) {
	var types []mysql.MAccountType
	err := r.db.Where("is_active = ?", 1).Find(&types).Error
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(types))
	for i, t := range types {
		result[i] = AccountTypeResponse{
			ID:          t.ID,
			Code:        t.Code,
			Name:        t.Name,
			Description: t.Description,
			Icon:        t.Icon,
		}
	}
	return result, nil
}

func (r *repo) GetCategories() ([]interface{}, error) {
	var categories []mysql.MTransactionCategory
	err := r.db.Where("is_active = ?", 1).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(categories))
	for i, c := range categories {
		result[i] = CategoryResponse{
			ID:          c.ID,
			Code:        c.Code,
			Name:        c.Name,
			Description: c.Description,
			Icon:        c.Icon,
			Color:       c.Color,
			IsExpense:   c.IsExpense == 1,
		}
	}
	return result, nil
}

func (r *repo) GetGoalTypes() ([]interface{}, error) {
	var types []mysql.MGoalType
	err := r.db.Where("is_active = ?", 1).Find(&types).Error
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(types))
	for i, t := range types {
		result[i] = GoalTypeResponse{
			ID:          t.ID,
			Code:        t.Code,
			Name:        t.Name,
			Description: t.Description,
			Icon:        t.Icon,
		}
	}
	return result, nil
}

func (r *repo) GetRiskProfiles() ([]interface{}, error) {
	var profiles []mysql.MRiskProfile
	err := r.db.Where("is_active = ?", 1).Find(&profiles).Error
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(profiles))
	for i, p := range profiles {
		result[i] = RiskProfileResponse{
			ID:          p.ID,
			Code:        p.Code,
			Name:        p.Name,
			Description: p.Description,
		}
	}
	return result, nil
}

func (r *repo) GetDependents() ([]interface{}, error) {
	var dependents []mysql.MDependent
	err := r.db.Where("is_active = ?", 1).Find(&dependents).Error
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(dependents))
	for i, d := range dependents {
		result[i] = DependentResponse{
			ID:          d.ID,
			Code:        d.Code,
			Name:        d.Name,
			Description: d.Description,
		}
	}
	return result, nil
}
