package profiles

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

// GetProfile godoc
// @Summary      Public - Get Profile
// @Description  Retrieve the user profile
// @Tags         Public - Profile
// @Produce      json
// @Success      200  {object}  Profile
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /public/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	profile, err := h.service.GetProfile()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch profile", err.Error())
		return
	}
	if profile == nil {
		response.Error(c, http.StatusNotFound, "Profile not found", "Profile not found")
		return
	}
	response.Success(c, http.StatusOK, "Profile fetched successfully", profile)
}

// UpdateProfile godoc
// @Summary      Admin - Update Profile
// @Description  Update details of the user profile
// @Tags         Admin - Profile
// @Accept       json
// @Produce      json
// @Param        request body UpdateProfileRequest true "Profile Data"
// @Security     BearerAuth
// @Success      200  {object}  Profile
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
		return
	}
	
	// Handle both string and uuid.UUID types for user_id
	var userID uuid.UUID
	var err error
	
	switch v := userIDVal.(type) {
	case string:
		userID, err = uuid.Parse(v)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid user ID format", "Invalid user ID format")
			return
		}
	case uuid.UUID:
		userID = v
	default:
		response.Error(c, http.StatusInternalServerError, "Invalid user ID type", "Invalid user ID type")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	profile, err := h.service.CreateOrUpdateProfile(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update profile", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Profile updated successfully", profile)
}
