package auth

import (
	"kayakaga-api/domain/mysql"
)

type Repository interface {
	CreateUser(user *mysql.User) error
	CreateUserCredential(cred *mysql.UserCredential) error
	CreateRefreshToken(token *mysql.RefreshToken) error
	GetUserByEmail(email string) (*mysql.User, *mysql.UserCredential, error)
	GetRefreshToken(token string) (*mysql.RefreshToken, error)
	RevokeRefreshToken(token string) error
}

type UseCase interface {
	Register(req *RegisterRequest) (*AuthResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	Refresh(req *RefreshRequest) (*RefreshResponse, error)
	Logout(req *LogoutRequest) error
}
