package auth

import (
    "net/http"
    "strings"
	"errors"
	"strconv"
	"time"
	
    "github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type registerRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

type loginRequest struct {
	Email        string `json:"Email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

func (_handle *Handler) Register(_context *gin.Context) {
	var req registerRequest
	if err := _context.ShouldBindJSON(&req); err != nil {
		_context.JSON(http.StatusBadRequest, gin.H{
			"error": "username, a valid email, password and password_confirm are required",
		})
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		_context.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	if req.Password != req.PasswordConfirm {
		_context.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	err := _handle.service.Register(_context.Request.Context(), email, req.Password)
	if errors.Is(err, ErrEmailTaken) {
		_context.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}
	if err != nil {
		_context.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	_context.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (_handle *Handler) Login(_context *gin.Context){
	var req loginRequest
	if err := _context.ShouldBindJSON(&req); err != nil {
		_context.JSON(http.StatusBadRequest, gin.H{
			"error": "username, password are required",
		})
		return
	}
	email := strings.TrimSpace(req.Email)
	if email == "" {
		_context.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	password := strings.TrimSpace(req.Password)
	if password == "" {
		_context.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	token, err := _handle.service.Login(_context.Request.Context(), req.Email, req.Password)
	if errors.Is(err, ErrInvalidCredentials) {
		_context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err != nil {
		_context.JSON(http.StatusInternalServerError, gin.H{"error": "could not log in"})
		return
	}

	_context.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
	})

}

func (_handle *Handler) Logout(_context *gin.Context) {
    jti := _context.GetString("token_jti")
    value, _ := _context.Get("token_expires_at")
    expiresAt, ok := value.(time.Time)

    if jti == "" || !ok || expiresAt.IsZero() {
        _context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
        return
    }

    if err := _handle.service.RevokeToken(_context.Request.Context(), jti, expiresAt); err != nil {
        _context.JSON(http.StatusInternalServerError, gin.H{"error": "could not log out"})
        return
    }

    _context.JSON(http.StatusOK, gin.H{"message": "logout success"})
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userID, err := strconv.ParseInt(c.GetString("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "old_password and new_password are required",
		})
		return
	}

	err = h.service.ChangePassword(
		c.Request.Context(), userID, req.OldPassword, req.NewPassword,
	)
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		c.JSON(http.StatusBadRequest, gin.H{"error": "current password is incorrect"})
	case errors.Is(err, ErrPasswordChangedConcurrently):
		c.JSON(http.StatusConflict, gin.H{"error": "password changed during request"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not change password"})
	default:
		c.Status(http.StatusNoContent)
	}
}