# Prompt for Claude Code — Kayakaga API (Golang + Gin)

---

Build a REST API called **kayakaga-api** for Kayakaga, a Personal Finance Advisor app.
Stack: Golang, Gin framework, GORM, MySQL.

---

## Project Structure

Follow this exact structure (adapted from existing codebase pattern):

```
kayakaga-api/
├── main.go
├── go.mod
├── go.sum
├── .env
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── README.md
│
├── di/
│   ├── domain.go          # Wire provider definitions for all domains
│   └── wire_gen.go        # Generated wire code (manual for now, no wire CLI needed)
│
├── domain/
│   └── mysql/
│       ├── mysql.go       # MySQL driver + GORM connection
│       ├── orm.go         # All GORM struct definitions
│       ├── source.go      # MySQL repository interfaces
│       └── providers/
│           ├── user.go        # User + credential queries
│           ├── account.go     # Account queries
│           ├── transaction.go # Transaction queries
│           ├── goal.go        # Goal + milestone queries
│           └── analytics.go   # Budget, anomaly, recurring queries
│
├── modules/
│   ├── auth/
│   │   ├── domain.go      # Repo & UseCase interfaces
│   │   ├── entities.go    # Request/response structs
│   │   ├── handler.go     # HTTP handlers
│   │   ├── repo.go        # Repository implementation
│   │   └── router.go      # Route definitions
│   │
│   ├── user/
│   │   ├── domain.go
│   │   ├── entities.go
│   │   ├── handler.go
│   │   ├── repo.go
│   │   └── router.go
│   │
│   ├── account/
│   │   ├── domain.go
│   │   ├── entities.go
│   │   ├── handler.go
│   │   ├── repo.go
│   │   └── router.go
│   │
│   ├── transaction/
│   │   ├── domain.go
│   │   ├── entities.go
│   │   ├── handler.go
│   │   ├── repo.go
│   │   └── router.go
│   │
│   ├── goal/
│   │   ├── domain.go
│   │   ├── entities.go
│   │   ├── handler.go
│   │   ├── repo.go
│   │   └── router.go
│   │
│   ├── analytics/
│   │   ├── domain.go
│   │   ├── entities.go
│   │   ├── handler.go
│   │   ├── repo.go
│   │   └── router.go
│   │
│   └── simulate/
│       ├── domain.go
│       ├── entities.go
│       ├── handler.go     # Pure logic, no DB
│       └── router.go
│
└── utils/
    ├── helper/
    │   └── common.go      # FormatRupiah, response helpers
    ├── router/
    │   ├── middleware.go  # JWT auth, CORS, logging
    │   └── router.go      # Register all module routes
    └── tokenizer/
        └── tokenizer.go   # JWT generate, validate, refresh
```

---

## Environment Variables (.env)

```
APP_PORT=8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=kayakaga_db

JWT_SECRET=finai-jwt-secret-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=720h

API_VERSION=v1
```

---

## Database — MySQL (GORM)

Database name: `kayakaga_db`

### GORM Structs (domain/mysql/orm.go)

All tables already exist in DB. Map them exactly:

```go
// Master tables — read only from API
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

// Core tables
type User struct {
    ID                       uint         `gorm:"primaryKey;autoIncrement"`
    Name                     string       `gorm:"size:100"`
    City                     string       `gorm:"size:100"`
    Profession               string       `gorm:"size:100"`
    DependentID              uint         `gorm:"default:1"`
    MonthlyIncome            int64        `gorm:"default:0"`
    MonthlyExpensesEstimate  int64        `gorm:"default:0"`
    CurrentSavings           int64        `gorm:"default:0"`
    RiskProfileID            uint         `gorm:"default:1"`
    Currency                 string       `gorm:"size:3;default:IDR"`
    CreatedAt                time.Time
    UpdatedAt                time.Time
    Dependent                MDependent   `gorm:"foreignKey:DependentID"`
    RiskProfile              MRiskProfile `gorm:"foreignKey:RiskProfileID"`
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
    ID            uint          `gorm:"primaryKey;autoIncrement"`
    UserID        uint
    AccountTypeID uint
    Name          string        `gorm:"size:100"`
    Balance       int64         `gorm:"default:0"`
    Color         string        `gorm:"size:7"`
    IsPrimary     int8          `gorm:"default:0"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
    AccountType   MAccountType  `gorm:"foreignKey:AccountTypeID"`
}
func (Account) TableName() string { return "accounts" }

type Transaction struct {
    ID            uint                 `gorm:"primaryKey;autoIncrement"`
    UserID        uint
    AccountID     uint
    CategoryID    uint
    SourceID      uint                 `gorm:"default:1"`
    Date          time.Time            `gorm:"type:date"`
    Time          *time.Time           `gorm:"type:time"`
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
    ID                  uint       `gorm:"primaryKey;autoIncrement"`
    UserID              uint
    AccountID           uint
    GoalTypeID          uint
    Name                string     `gorm:"size:200"`
    TargetAmount        int64      `gorm:"default:0"`
    CurrentAmount       int64      `gorm:"default:0"`
    MonthlyContribution int64      `gorm:"default:0"`
    TargetDate          *time.Time `gorm:"type:date"`
    CreatedAt           time.Time
    UpdatedAt           time.Time
    GoalType            MGoalType  `gorm:"foreignKey:GoalTypeID"`
    Milestones          []GoalMilestone `gorm:"foreignKey:GoalID"`
}
func (Goal) TableName() string { return "goals" }

type GoalMilestone struct {
    ID        uint       `gorm:"primaryKey;autoIncrement"`
    GoalID    uint
    Amount    int64
    ReachedAt *time.Time
    CreatedAt time.Time
}
func (GoalMilestone) TableName() string { return "goal_milestones" }
```

---

## Auth Flow

### Token Strategy
- Access Token: JWT, expire 15 minutes, payload: `{user_id, email, exp}`
- Refresh Token: UUID string, expire 30 days, stored in `refresh_tokens` table

### Endpoints
```
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
```

### Register Request/Response
```json
// Request
{
  "name": "Andi Pratama",
  "email": "andi@finai.dev",
  "password": "password123"
}

// Response
{
  "success": true,
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "uuid-string",
    "user": {
      "id": 1,
      "name": "Andi Pratama",
      "email": "andi@finai.dev"
    }
  }
}
```

### Login Request/Response
```json
// Request
{ "email": "andi@finai.dev", "password": "password123" }

// Response — same as register
```

### Refresh Request/Response
```json
// Request
{ "refresh_token": "uuid-string" }

// Response
{
  "success": true,
  "data": { "access_token": "eyJ..." }
}
```

### Logout Request
```json
// Header: Authorization: Bearer <access_token>
// Body: { "refresh_token": "uuid-string" }
// Action: revoke refresh token in DB
```

---

## JWT Middleware (utils/router/middleware.go)

All routes except `/api/v1/auth/*` require:
```
Header: Authorization: Bearer <access_token>
```

Middleware must:
1. Extract token from Authorization header
2. Validate JWT signature and expiry
3. Extract `user_id` from claims
4. Inject `user_id` into Gin context: `c.Set("user_id", userID)`
5. Handlers get user_id via: `userID := c.GetUint("user_id")`

NEVER trust user_id from request body or query params for ownership checks.

---

## Standard Response Format

All endpoints use this format:

```go
// utils/helper/common.go

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
    Meta    *MetaInfo   `json:"meta,omitempty"`
}

type ErrorInfo struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type MetaInfo struct {
    Timestamp string `json:"timestamp"`
    Total     *int64 `json:"total,omitempty"`
    Page      *int   `json:"page,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Success: true,
        Data:    data,
        Meta:    &MetaInfo{Timestamp: time.Now().UTC().Format(time.RFC3339)},
    })
}

func ErrorResponse(c *gin.Context, status int, code string, message string) {
    c.JSON(status, Response{
        Success: false,
        Error:   &ErrorInfo{Code: code, Message: message},
    })
}
```

---

## All Endpoints

### Users (protected)
```
GET  /api/v1/users/profile        → get current user profile
PUT  /api/v1/users/profile        → update profile
                                    body: {city, profession, dependent_id,
                                           monthly_income, monthly_expenses_estimate,
                                           current_savings, risk_profile_id}
```

### Accounts (protected)
```
GET    /api/v1/accounts           → list all accounts for current user
GET    /api/v1/accounts/balances  → balances per account + grand total
POST   /api/v1/accounts           → create account
                                    body: {name, account_type_id, balance, color, is_primary}
PUT    /api/v1/accounts/:id       → update account
DELETE /api/v1/accounts/:id       → delete account
```

### Transactions (protected)
```
GET    /api/v1/transactions       → list transactions
                                    query: period (today|week|month|last_month|year),
                                           account_id?, category_id?, merchant?, 
                                           is_recurring?
                                    response includes: list + summary {total_in, total_out, count}
POST   /api/v1/transactions       → create transaction
                                    body: {account_id, category_id, source_id,
                                           date, time?, merchant, amount, notes?,
                                           is_recurring}
PUT    /api/v1/transactions/:id   → update transaction
DELETE /api/v1/transactions/:id   → delete transaction
```

### Goals (protected)
```
GET    /api/v1/goals              → list goals with progress % and ETA
GET    /api/v1/goals/:id          → single goal with milestones
POST   /api/v1/goals              → create goal
                                    body: {name, goal_type_id, account_id,
                                           target_amount, monthly_contribution,
                                           target_date?, milestones[]}
PUT    /api/v1/goals/:id          → update goal
DELETE /api/v1/goals/:id          → delete goal
PUT    /api/v1/goals/:id/contribution  → update monthly contribution
                                         body: {amount}
```

### Analytics (protected)
```
GET /api/v1/analytics/budget
    query: period (month|last_month|year), account_id?
    response: {
      income, expenses, savings_rate,
      breakdown: [{category_id, category_name, total, percentage}],
      vs_previous: {income_delta_pct, expenses_delta_pct}
    }

GET /api/v1/analytics/compare
    query: account_id?
    response: [{
      category_id, category_name,
      this_month, last_month,
      delta_pct, trend (up|down|new|gone)
    }]

GET /api/v1/analytics/anomalies
    query: period (week|month), account_id?
    response: [{
      transaction_id, merchant, amount, date,
      reason, severity (high|medium|low)
    }]

GET /api/v1/analytics/recurring
    query: account_id?
    response: {
      items: [{merchant, amount, account_id, last_charged, next_renewal, category_name}],
      total_monthly: int64
    }

GET /api/v1/analytics/savings-suggestion
    query: target_savings? (int64)
    response: {
      suggestions: [{category_id, category_name, current_spend,
                     suggested_limit, potential_saving, reasoning}],
      total_potential_saving: int64,
      impact_on_goals: string
    }

GET /api/v1/analytics/goal-recommendation
    query: goal_id (required), target_months?, new_monthly_contribution?
    response: {
      goal_name, remaining_amount,
      current_contribution, current_eta_months, current_eta_date,
      scenarios: [{monthly_contribution, eta_months, months_faster, eta_date}]
    }
```

### Simulate (protected)
```
GET /api/v1/simulate/investment
    query: monthly_amount, annual_return_pct, years
    response: {
      future_value, total_invested, profit,
      roi_pct, monthly_breakdown: [{month, total}]
    }
```

### Master Data (protected — for dropdowns)
```
GET /api/v1/masters/account-types
GET /api/v1/masters/categories
GET /api/v1/masters/goal-types
GET /api/v1/masters/risk-profiles
GET /api/v1/masters/dependents
```

---

## Key Implementation Notes

### Period Filter Logic (transactions)
```go
// today     → date = CURDATE()
// week      → date >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)
// month     → YEAR(date) = YEAR(CURDATE()) AND MONTH(date) = MONTH(CURDATE())
// last_month→ YEAR(date) = YEAR(DATE_SUB(CURDATE(),INTERVAL 1 MONTH))
//             AND MONTH(date) = MONTH(DATE_SUB(CURDATE(),INTERVAL 1 MONTH))
// year      → YEAR(date) = YEAR(CURDATE())
```

### Goal Progress & ETA Calculation
```go
// progress_pct = (current_amount / target_amount) * 100
// remaining    = target_amount - current_amount
// eta_months   = ceil(remaining / monthly_contribution)
// eta_date     = time.Now().AddDate(0, eta_months, 0)
```

### Anomaly Detection Logic
```go
// 1. Calculate avg amount per category for the period
// 2. Flag if amount > 2x category average → severity based on multiplier:
//    2x-3x = low, 3x-5x = medium, >5x = high
// 3. Flag merchant appearing for first time with amount > 100000
// 4. Flag duplicate: same merchant + same amount within 24 hours
```

### Investment Simulation (compound interest)
```go
// monthly rate = annual_return_pct / 12 / 100
// future_value = monthly_amount * ((1 + rate)^n - 1) / rate
// where n = years * 12
```

### Ownership Check Pattern
```go
// Always verify resource belongs to current user
// Example for account:
func (r *accountRepo) GetByID(id, userID uint) (*Account, error) {
    var account orm.Account
    err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&account).Error
    // ...
}
```

---

## go.mod Dependencies

```
github.com/gin-gonic/gin
gorm.io/gorm
gorm.io/driver/mysql
github.com/golang-jwt/jwt/v5
github.com/joho/godotenv
golang.org/x/crypto          ← bcrypt
github.com/google/uuid        ← for refresh token generation
```

---

## Build Order

Build in this order — each step must compile before moving to next:

1. `go.mod` + dependencies
2. `.env` + `domain/mysql/mysql.go` (DB connection)
3. `domain/mysql/orm.go` (all GORM structs)
4. `utils/helper/common.go` (response helpers)
5. `utils/tokenizer/tokenizer.go` (JWT logic)
6. `utils/router/middleware.go` (auth middleware)
7. `modules/auth/` (full module — register + login + refresh + logout)
8. `modules/user/` (profile get + update)
9. `modules/account/` (CRUD + balances)
10. `modules/transaction/` (CRUD + filtered list)
11. `modules/goal/` (CRUD + contribution update)
12. `modules/analytics/` (all 6 analytics endpoints)
13. `modules/simulate/` (investment simulation)
14. `di/domain.go` (wire everything together)
15. `utils/router/router.go` (register all routes)
16. `main.go` (entry point)

---

## Testing Checklist

After build, verify these work with curl or Postman:

1. `POST /api/v1/auth/register` → get tokens
2. `POST /api/v1/auth/login` → get tokens
3. `GET  /api/v1/accounts/balances` → 3 accounts, total 23.4jt
4. `GET  /api/v1/transactions?period=month` → 30 transactions
5. `GET  /api/v1/transactions?period=month&account_id=3` → GoPay only
6. `GET  /api/v1/analytics/budget?period=month` → income/expense breakdown
7. `GET  /api/v1/analytics/anomalies?period=month` → Steam Games + Seafood Ancol
8. `GET  /api/v1/analytics/recurring` → Netflix, Spotify, PLN, etc
9. `GET  /api/v1/goals` → 2 goals with progress %
10. `GET  /api/v1/simulate/investment?monthly_amount=2000000&annual_return_pct=10&years=10`
    → future_value ~413jt

---

## Important Constraints

- Every file max 200 lines — split if needed
- All monetary values as int64 (IDR, no decimals)
- Dates as `time.Time` with `gorm:"type:date"`
- Always filter by `user_id` from JWT context — never from request params
- Return proper HTTP status codes:
  200 OK, 201 Created, 400 Bad Request,
  401 Unauthorized, 403 Forbidden, 404 Not Found, 500 Internal Server Error
- No hardcoded values — all config from .env
- Add logging for each request (use Gin's built-in logger)
