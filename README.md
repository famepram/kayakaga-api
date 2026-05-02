# Kayakaga API

Personal Finance Advisor REST API built with Go, Gin framework, GORM, and MySQL.

## Tech Stack

- **Go** 1.21+
- **Gin** - Web framework
- **GORM** - ORM for database operations
- **MySQL** 8.0 - Database
- **JWT** - Authentication

## Features

- User authentication with JWT (access + refresh tokens)
- Profile management
- Account management with balance tracking
- Transaction tracking with advanced filtering
- Goal tracking with progress & ETA calculations
- Advanced analytics (budget, anomalies, recurring transactions)
- Investment simulation calculator

## Prerequisites

- Go 1.21+
- MySQL 8.0+
- Docker & Docker Compose (optional)

## Quick Start with Docker

### 1. Clone & Setup
```bash
cd kayakaga-api
```

### 2. Start Services
```bash
docker-compose up -d
```

This will start:
- **MySQL** on port 3306
- **Kayakaga API** on port 8080

### 3. Check Status
```bash
docker-compose ps
docker-compose logs api
```

### 4. Stop Services
```bash
docker-compose down
```

### 5. Reset Everything (including database)
```bash
docker-compose down -v
docker-compose up -d
```

## Manual Setup (without Docker)

### 1. Install Dependencies
```bash
go mod download
```

### 2. Setup Database
Create MySQL database named `kayakaga_db` and run the schema scripts.

### 3. Configure Environment
Copy `.env` file and update with your credentials:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=kayakaga_db
```

### 4. Run
```bash
go run main.go
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout

### Users (Protected)
- `GET /api/v1/users/profile` - Get current user profile
- `PUT /api/v1/users/profile` - Update profile

### Accounts (Protected)
- `GET /api/v1/accounts` - List all accounts
- `GET /api/v1/accounts/balances` - Get balances with total
- `POST /api/v1/accounts` - Create account
- `PUT /api/v1/accounts/:id` - Update account
- `DELETE /api/v1/accounts/:id` - Delete account

### Transactions (Protected)
- `GET /api/v1/transactions` - List transactions with filters
- `POST /api/v1/transactions` - Create transaction
- `PUT /api/v1/transactions/:id` - Update transaction
- `DELETE /api/v1/transactions/:id` - Delete transaction

### Goals (Protected)
- `GET /api/v1/goals` - List goals with progress
- `GET /api/v1/goals/:id` - Get goal details
- `POST /api/v1/goals` - Create goal
- `PUT /api/v1/goals/:id` - Update goal
- `DELETE /api/v1/goals/:id` - Delete goal
- `PUT /api/v1/goals/:id/contribution` - Update monthly contribution

### Analytics (Protected)
- `GET /api/v1/analytics/budget` - Budget breakdown
- `GET /api/v1/analytics/compare` - Month-over-month comparison
- `GET /api/v1/analytics/anomalies` - Detect unusual spending
- `GET /api/v1/analytics/recurring` - List recurring transactions
- `GET /api/v1/analytics/savings-suggestion` - Get savings suggestions
- `GET /api/v1/analytics/goal-recommendation` - Goal scenarios

### Simulation (Protected)
- `GET /api/v1/simulate/investment` - Investment calculator

## Testing

### Register & Login
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Andi Pratama","email":"andi@test.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"andi@test.com","password":"password123"}'
```

### Access Protected Endpoint
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Project Structure

```
kayakaga-api/
├── main.go                 # Application entry point
├── go.mod / go.sum        # Go dependencies
├── .env                   # Environment variables
├── Dockerfile             # Docker build config
├── docker-compose.yml     # Docker services
├── di/                    # Dependency injection
├── domain/                # Domain layer (database)
│   └── mysql/            # MySQL & GORM
├── modules/              # Feature modules
│   ├── auth/
│   ├── user/
│   ├── account/
│   ├── transaction/
│   ├── goal/
│   ├── analytics/
│   └── simulate/
└── utils/                # Utilities
    ├── helper/
    ├── router/
    └── tokenizer/
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| APP_PORT | Server port | 8080 |
| APP_ENV | Environment | development |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 3306 |
| DB_USER | Database user | root |
| DB_PASSWORD | Database password | |
| DB_NAME | Database name | kayakaga_db |
| JWT_SECRET | JWT signing secret | |
| JWT_ACCESS_EXPIRY | Access token expiry | 15m |
| JWT_REFRESH_EXPIRY | Refresh token expiry | 720h |

## License

MIT
