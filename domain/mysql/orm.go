package mysql

import "time"

type MAccountType struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Code        string    `gorm:"size:50;uniqueIndex"`
	Name        string    `gorm:"size:100"`
	Description string    `gorm:"size:255"`
	Icon        string    `gorm:"size:100"`
	IsActive    int8      `gorm:"default:1"`
	CreatedAt   time.Time
}

func (MAccountType) TableName() string { return "m_account_types" }

type MTransactionCategory struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Code        string    `gorm:"size:50;uniqueIndex"`
	Name        string    `gorm:"size:100"`
	Description string    `gorm:"size:255"`
	Icon        string    `gorm:"size:100"`
	Color       string    `gorm:"size:7"`
	IsExpense   int8      `gorm:"default:1"`
	IsActive    int8      `gorm:"default:1"`
	CreatedAt   time.Time
}

func (MTransactionCategory) TableName() string { return "m_transaction_categories" }

type MTransactionSource struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Code        string    `gorm:"size:50;uniqueIndex"`
	Name        string    `gorm:"size:100"`
	Description string    `gorm:"size:255"`
	IsActive    int8      `gorm:"default:1"`
	CreatedAt   time.Time
}

func (MTransactionSource) TableName() string { return "m_transaction_sources" }

type MGoalType struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Code        string    `gorm:"size:50;uniqueIndex"`
	Name        string    `gorm:"size:100"`
	Description string    `gorm:"size:255"`
	Icon        string    `gorm:"size:100"`
	IsActive    int8      `gorm:"default:1"`
	CreatedAt   time.Time
}

func (MGoalType) TableName() string { return "m_goal_types" }

type MRiskProfile struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Code        string    `gorm:"size:50;uniqueIndex"`
	Name        string    `gorm:"size:100"`
	Description string    `gorm:"size:255"`
	IsActive    int8      `gorm:"default:1"`
	CreatedAt   time.Time
}

func (MRiskProfile) TableName() string { return "m_risk_profiles" }

type MDependent struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Code        string    `gorm:"size:50;uniqueIndex"`
	Name        string    `gorm:"size:100"`
	Description string    `gorm:"size:255"`
	IsActive    int8      `gorm:"default:1"`
	CreatedAt   time.Time
}

func (MDependent) TableName() string { return "m_dependents" }

type User struct {
	ID                      uint        `gorm:"primaryKey;autoIncrement"`
	Name                    string      `gorm:"size:100"`
	City                    string      `gorm:"size:100"`
	Profession              string      `gorm:"size:100"`
	DependentID             uint        `gorm:"default:1"`
	MonthlyIncome           int64       `gorm:"default:0"`
	MonthlyExpensesEstimate int64       `gorm:"default:0"`
	CurrentSavings          int64       `gorm:"default:0"`
	RiskProfileID           uint        `gorm:"default:1"`
	Currency                string      `gorm:"size:3;default:IDR"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
	Dependent               MDependent  `gorm:"foreignKey:DependentID"`
	RiskProfile             MRiskProfile `gorm:"foreignKey:RiskProfileID"`
}

func (User) TableName() string { return "users" }

type UserCredential struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UserID       uint      `gorm:"uniqueIndex"`
	Email        string    `gorm:"size:255;uniqueIndex"`
	PasswordHash string    `gorm:"size:255"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (UserCredential) TableName() string { return "user_credentials" }

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint
	Token     string    `gorm:"size:255;uniqueIndex"`
	ExpiresAt time.Time
	Revoked   int8      `gorm:"default:0"`
	CreatedAt time.Time
}

func (RefreshToken) TableName() string { return "refresh_tokens" }

type Account struct {
	ID            uint         `gorm:"primaryKey;autoIncrement"`
	UserID        uint
	AccountTypeID uint
	Name          string       `gorm:"size:100"`
	Balance       int64        `gorm:"default:0"`
	Color         string       `gorm:"size:7"`
	IsPrimary     int8         `gorm:"default:0"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	AccountType   MAccountType `gorm:"foreignKey:AccountTypeID"`
}

func (Account) TableName() string { return "accounts" }

type Transaction struct {
	ID            uint                 `gorm:"primaryKey;autoIncrement"`
	UserID        uint
	AccountID     uint
	CategoryID    uint
	SourceID      uint                 `gorm:"default:1"`
	Date          time.Time            `gorm:"type:date"`
	Time          *string              `gorm:"type:time"`
	Merchant      string               `gorm:"size:200"`
	Amount        int64
	Notes         string               `gorm:"type:text"`
	AiCategorized int8                 `gorm:"default:0"`
	IsRecurring   int8                 `gorm:"default:0"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Category      MTransactionCategory `gorm:"foreignKey:CategoryID"`
	Source        MTransactionSource   `gorm:"foreignKey:SourceID"`
	Account       Account              `gorm:"foreignKey:AccountID"`
}

func (Transaction) TableName() string { return "transactions" }

type Goal struct {
	ID                  uint           `gorm:"primaryKey;autoIncrement"`
	UserID              uint
	AccountID           uint
	GoalTypeID          uint
	Name                string         `gorm:"size:200"`
	TargetAmount        int64          `gorm:"default:0"`
	CurrentAmount       int64          `gorm:"default:0"`
	MonthlyContribution int64          `gorm:"default:0"`
	TargetDate          *time.Time     `gorm:"type:date"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	GoalType            MGoalType      `gorm:"foreignKey:GoalTypeID"`
	Milestones          []GoalMilestone `gorm:"foreignKey:GoalID"`
}

func (Goal) TableName() string { return "goals" }

type GoalMilestone struct {
	ID        uint       `gorm:"primaryKey;autoIncrement"`
	GoalID    uint
	Amount    int64
	ReachedAt *time.Time `gorm:"type:date"`
	CreatedAt time.Time
}

func (GoalMilestone) TableName() string { return "goal_milestones" }
