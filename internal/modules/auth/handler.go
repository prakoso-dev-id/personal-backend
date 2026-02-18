package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prakoso-id/personal-backend/internal/utils/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary      Admin - Login
// @Description  Authenticates an admin user and returns a JWT token
// @Tags         Admin - Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login Credentials"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /admin/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	fullname := ""
	if user.Profile != nil {
		fullname = user.Profile.FullName
	}

	response.Success(c, http.StatusOK, "Login successful", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"fullname": fullname,
		},
	})
}

type UpdateEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// UpdateEmail godoc
// @Summary      Admin - Update Email
// @Description  Update the authenticated admin's email address
// @Tags         Admin - Auth
// @Accept       json
// @Produce      json
// @Param        request body UpdateEmailRequest true "New Email"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/update-email [put]
func (h *Handler) UpdateEmail(c *gin.Context) {
	var req UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
		return
	}
	userID, _ := uuid.Parse(userIDVal.(string))

	if err := h.service.UpdateEmail(userID, req.Email); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update email", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Email updated successfully", nil)
}

type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

// UpdatePassword godoc
// @Summary      Admin - Update Password
// @Description  Update the authenticated admin's password
// @Tags         Admin - Auth
// @Accept       json
// @Produce      json
// @Param        request body UpdatePasswordRequest true "New Password"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/update-password [put]
func (h *Handler) UpdatePassword(c *gin.Context) {
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
		return
	}
	userID, _ := uuid.Parse(userIDVal.(string))

	if err := h.service.UpdatePassword(userID, req.Password); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update password", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password updated successfully", nil)
}
