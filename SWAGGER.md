# Swagger Documentation for Kayakaga API

## Overview

The Kayakaga API includes Swagger/OpenAPI documentation that can be accessed via Swagger UI or imported into tools like Postman, Insomnia, or other API clients.

## Accessing Swagger UI

Once the API is running, access the interactive Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

This provides:
- Interactive API documentation
- Try-it-out functionality for all endpoints
- Request/response schemas
- Authentication handling

## Importing into Postman

### Option 1: Import from Running API

1. Start the API server: `go run main.go` or `docker-compose up`
2. Open Postman
3. Go to **File** → **Import**
4. Enter URL: `http://localhost:8080/swagger/doc.json`
5. Click **Import**

### Option 2: Import from File

1. Use the [`docs/swagger.yaml`](docs/swagger.yaml) file
2. In Postman: **File** → **Import** → **Upload Files**
3. Select `swagger.yaml`
4. Postman will create a collection with all endpoints

### Option 3: Import from URL (Production)

```
https://api.kayakaga.com/swagger/doc.json
```

## Importing into Other Tools

### Insomnia
- **File** → **Import** → **From File**
- Select `docs/swagger.yaml`

### HTTPie (CLI)
```bash
# Install httpie: pip install httpie
# Use with Swagger: httpie --auth-type=bearer --auth="YOUR_TOKEN" GET localhost:8080/api/v1/users/profile
```

### VS Code REST Client
- Install REST Client extension
- Create `http` files and reference the swagger schema

## Generating Full Documentation

The current Swagger setup includes basic annotations on key endpoints. To generate complete documentation for ALL endpoints:

### 1. Install Swag
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. Add Annotations to Handlers
Add Swagger annotations to your handler functions (examples already provided in `modules/auth/handler.go`, `modules/user/handler.go`, etc.):

```go
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
    // ...
}
```

### 3. Generate Documentation
```bash
swag init
```

This will scan all Go files and generate:
- `docs/docs.go` - Go documentation
- `docs/swagger.json` - JSON schema
- `docs/swagger.yaml` - YAML schema

### 4. Restart Server
```bash
go run main.go
```

## Authentication

Swagger UI supports testing authenticated endpoints:

1. Call `POST /api/v1/auth/login` to get tokens
2. Click **Authorize** button (lock icon)
3. Enter: `Bearer YOUR_ACCESS_TOKEN`
4. Click **Authorize**
5. Now you can test protected endpoints

## Available Endpoints

Currently documented with Swagger annotations:

### Authentication (No Auth Required)
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout

### Users (Bearer Token Required)
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update profile

### Accounts (Bearer Token Required)
- `GET /api/v1/accounts` - List accounts
- `GET /api/v1/accounts/balances` - Get balances with total
- `POST /api/v1/accounts` - Create account
- `PUT /api/v1/accounts/:id` - Update account
- `DELETE /api/v1/accounts/:id` - Delete account

### Transactions (Bearer Token Required)
- `GET /api/v1/transactions` - List transactions with filters
- `POST /api/v1/transactions` - Create transaction
- `PUT /api/v1/transactions/:id` - Update transaction
- `DELETE /api/v1/transactions/:id` - Delete transaction

### Goals (Bearer Token Required)
- `GET /api/v1/goals` - List goals with progress
- `GET /api/v1/goals/:id` - Get goal details
- `POST /api/v1/goals` - Create goal
- `PUT /api/v1/goals/:id` - Update goal
- `DELETE /api/v1/goals/:id` - Delete goal
- `PUT /api/v1/goals/:id/contribution` - Update monthly contribution

### Analytics (Bearer Token Required)
- `GET /api/v1/analytics/budget` - Budget breakdown
- `GET /api/v1/analytics/compare` - Month-over-month comparison
- `GET /api/v1/analytics/anomalies` - Detect unusual spending
- `GET /api/v1/analytics/recurring` - List recurring transactions
- `GET /api/v1/analytics/savings-suggestion` - Get savings suggestions
- `GET /api/v1/analytics/goal-recommendation` - Goal scenarios

### Simulation (Bearer Token Required)
- `GET /api/v1/simulate/investment` - Investment calculator

## Customizing Swagger

### Change API Info
Edit the annotations in [`main.go`](main.go):
```go
// @title Kayakaga API
// @version 1.0
// @description Your description
// @host localhost:8080
```

### Add Custom Schemas
Add model definitions in your handler files and reference them:
```go
// @Success 200 {object} helper.Response{data=YourCustomType}
```

## Troubleshooting

### Swagger UI Not Loading
- Ensure the API server is running
- Check port 8080 is available
- Verify `go mod tidy` has been run

### "Cannot find package docs"
- Run: `swag init`
- Or ensure `docs/docs.go` exists

### Authorization Not Working
- Make sure to include "Bearer " before the token
- Token must be valid (not expired)
- Check JWT_SECRET in .env matches

## Best Practices

1. **Keep annotations current** - Update Swagger annotations when modifying handlers
2. **Use specific types** - Define proper structs for request/response
3. **Document errors** - Include all possible error responses
4. **Add examples** - Use `@Example` annotations for clarity
5. **Group by tags** - Use `@Tags` to organize endpoints logically

## Resources

- [Swagger Specification](https://swagger.io/specification/)
- [Swag for Go](https://github.com/swaggo/swag)
- [Gin Swagger](https://github.com/swaggo/gin-swagger)
- [Postman Import](https://learning.postman.com/docs/getting-started/importing-and-exporting-data/)
