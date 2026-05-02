package auth

import (
	"errors"
	"kayakaga-api/domain/mysql"
	"kayakaga-api/utils/helper"
	"kayakaga-api/utils/tokenizer"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &handler{repo: repo}
}

func (h *handler) Register(req *RegisterRequest) (*AuthResponse, error) {
	_, _, err := h.repo.GetUserByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &mysql.User{
		Name:      req.Name,
		Currency:  "IDR",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := h.repo.CreateUser(user); err != nil {
		return nil, err
	}

	cred := &mysql.UserCredential{
		UserID:       user.ID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := h.repo.CreateUserCredential(cred); err != nil {
		return nil, err
	}

	accessToken, err := tokenizer.GenerateAccessToken(user.ID, req.Email)
	if err != nil {
		return nil, err
	}

	refreshToken := tokenizer.GenerateRefreshToken()
	rt := &mysql.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: tokenizer.GetRefreshTokenExpiry(),
		CreatedAt: time.Now().UTC(),
	}

	if err := h.repo.CreateRefreshToken(rt); err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Email: req.Email,
		},
	}, nil
}

func (h *handler) Login(req *LoginRequest) (*AuthResponse, error) {
	user, cred, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(cred.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := tokenizer.GenerateAccessToken(user.ID, cred.Email)
	if err != nil {
		return nil, err
	}

	refreshToken := tokenizer.GenerateRefreshToken()
	rt := &mysql.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: tokenizer.GetRefreshTokenExpiry(),
		CreatedAt: time.Now().UTC(),
	}

	if err := h.repo.CreateRefreshToken(rt); err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Email: cred.Email,
		},
	}, nil
}

func (h *handler) Refresh(req *RefreshRequest) (*RefreshResponse, error) {
	rt, err := h.repo.GetRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	var cred mysql.UserCredential
	err = h.repo.(*repo).db.Where("user_id = ?", rt.UserID).First(&cred).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := tokenizer.GenerateAccessToken(rt.UserID, cred.Email)
	if err != nil {
		return nil, err
	}

	return &RefreshResponse{
		AccessToken: accessToken,
	}, nil
}

func (h *handler) Logout(req *LogoutRequest) error {
	return h.repo.RevokeRefreshToken(req.RefreshToken)
}

// RegisterHandler godoc
// @Summary Register new user
// @Description Register a new user account with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} helper.Response{data=AuthResponse}
// @Failure 400 {object} helper.Response
// @Router /auth/register [post]
func RegisterHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		resp, err := uc.Register(&req)
		if err != nil {
			helper.ErrorResponse(c, 400, "REGISTER_FAILED", err.Error())
			return
		}

		helper.CreatedResponse(c, resp)
	}
}

// LoginHandler godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} helper.Response{data=AuthResponse}
// @Failure 401 {object} helper.Response
// @Router /auth/login [post]
func LoginHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		resp, err := uc.Login(&req)
		if err != nil {
			helper.ErrorResponse(c, 401, "LOGIN_FAILED", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

// RefreshHandler godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} helper.Response{data=RefreshResponse}
// @Failure 401 {object} helper.Response
// @Router /auth/refresh [post]
func RefreshHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		resp, err := uc.Refresh(&req)
		if err != nil {
			helper.ErrorResponse(c, 401, "REFRESH_FAILED", err.Error())
			return
		}

		helper.SuccessResponse(c, resp)
	}
}

// LogoutHandler godoc
// @Summary Logout user
// @Description Logout user and revoke refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "Refresh token to revoke"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /auth/logout [post]
func LogoutHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LogoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		if err := uc.Logout(&req); err != nil {
			helper.ErrorResponse(c, 400, "LOGOUT_FAILED", err.Error())
			return
		}

		helper.SuccessResponse(c, gin.H{"message": "logged out successfully"})
	}
}
