package masters

type Repository interface {
	GetAccountTypes() ([]interface{}, error)
	GetCategories() ([]interface{}, error)
	GetGoalTypes() ([]interface{}, error)
	GetRiskProfiles() ([]interface{}, error)
	GetDependents() ([]interface{}, error)
}

type UseCase interface {
	GetAccountTypes() ([]interface{}, error)
	GetCategories() ([]interface{}, error)
	GetGoalTypes() ([]interface{}, error)
	GetRiskProfiles() ([]interface{}, error)
	GetDependents() ([]interface{}, error)
}
