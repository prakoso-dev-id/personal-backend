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
// @Description  Update details of the user profile. Avatar and resume are uploaded as files.
// @Tags         Admin - Profile
// @Accept       multipart/form-data
// @Produce      json
// @Param        full_name formData string false "Full Name"
// @Param        bio       formData string false "Bio"
// @Param        avatar    formData file   false "Avatar Image (jpg, jpeg, png, webp - max 5MB)"
// @Param        resume    formData file   false "Resume File (pdf, doc, docx - max 10MB)"
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
	req.FullName = c.PostForm("full_name")
	req.Bio = c.PostForm("bio")

	// Get avatar file (optional)
	avatarFile, err := c.FormFile("avatar")
	if err == nil {
		req.AvatarFile = avatarFile
	}

	// Get resume file (optional)
	resumeFile, err := c.FormFile("resume")
	if err == nil {
		req.ResumeFile = resumeFile
	}

	profile, err := h.service.CreateOrUpdateProfile(userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update profile", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Profile updated successfully", profile)
}
