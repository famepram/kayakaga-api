package di

import (
	"kayakaga-api/modules/account"
	"kayakaga-api/modules/analytics"
	"kayakaga-api/modules/auth"
	"kayakaga-api/modules/goal"
	"kayakaga-api/modules/simulate"
	"kayakaga-api/modules/transaction"
	"kayakaga-api/modules/user"

	"gorm.io/gorm"
)

type Container struct {
	Auth       auth.UseCase
	User       user.UseCase
	Account    account.UseCase
	Transaction transaction.UseCase
	Goal       goal.UseCase
	Analytics  analytics.UseCase
	Simulate   simulate.UseCase
}

func InitializeContainer(db *gorm.DB) *Container {
	authRepo := auth.NewRepository(db)
	userRepo := user.NewRepository(db)
	accountRepo := account.NewRepository(db)
	transactionRepo := transaction.NewRepository(db)
	goalRepo := goal.NewRepository(db)
	analyticsRepo := analytics.NewRepository(db)

	return &Container{
		Auth:       auth.NewUseCase(authRepo),
		User:       user.NewUseCase(userRepo),
		Account:    account.NewUseCase(accountRepo),
		Transaction: transaction.NewUseCase(transactionRepo),
		Goal:       goal.NewUseCase(goalRepo),
		Analytics:  analytics.NewUseCase(analyticsRepo),
		Simulate:   simulate.NewUseCase(),
	}
}
