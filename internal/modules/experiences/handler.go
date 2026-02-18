package experiences

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prakoso-id/personal-backend/internal/modules/profiles"
	"github.com/prakoso-id/personal-backend/internal/utils/response"
)

type Handler struct {
	service        Service
	profileService profiles.Service
}

func NewHandler(service Service, profileService profiles.Service) *Handler {
	return &Handler{
		service:        service,
		profileService: profileService,
	}
}

// Public

// GetPublicExperiences godoc
// @Summary      Public - Get All Experiences
// @Description  Retrieve a list of all experiences
// @Tags         Public - Experiences
// @Produce      json
// @Success      200  {array}   Experience
// @Failure      500  {object}  map[string]string
// @Router       /public/experiences [get]
func (h *Handler) GetPublicExperiences(c *gin.Context) {
	experiences, err := h.service.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch experiences", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Experiences fetched successfully", experiences)
}

// Admin

// GetAdminExperiences godoc
// @Summary      Admin - Get All Experiences
// @Description  Retrieve a list of all experiences for admin
// @Tags         Admin - Experiences
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   Experience
// @Failure      500  {object}  map[string]string
// @Router       /admin/experiences [get]
func (h *Handler) GetAdminExperiences(c *gin.Context) {
	experiences, err := h.service.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch experiences", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Experiences fetched successfully", experiences)
}

// CreateExperience godoc
// @Summary      Admin - Create Experience
// @Description  Create a new experience
// @Tags         Admin - Experiences
// @Accept       json
// @Produce      json
// @Param        request body CreateExperienceRequest true "Experience Data"
// @Security     BearerAuth
// @Success      201  {object}  Experience
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/experiences [post]
func (h *Handler) CreateExperience(c *gin.Context) {
	var req CreateExperienceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	// Fetch Profile ID
	profile, err := h.profileService.GetProfileByUserID(uuid.MustParse(userID.(string)))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch profile", err.Error())
		return
	}
	if profile == nil {
		response.Error(c, http.StatusBadRequest, "Profile not found", "User profile does not exist")
		return
	}

	req.ProfileID = profile.ID

	experience, err := h.service.Create(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create experience", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "Experience created successfully", experience)
}

// UpdateExperience godoc
// @Summary      Admin - Update Experience
// @Description  Update an existing experience
// @Tags         Admin - Experiences
// @Accept       json
// @Produce      json
// @Param        id   path     string  true  "Experience ID"
// @Param        request body UpdateExperienceRequest true "Experience Data"
// @Security     BearerAuth
// @Success      200  {object}  Experience
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/experiences/{id} [put]
func (h *Handler) UpdateExperience(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", "invalid id")
		return
	}

	var req UpdateExperienceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	experience, err := h.service.Update(id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update experience", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Experience updated successfully", experience)
}

// DeleteExperience godoc
// @Summary      Admin - Delete Experience
// @Description  Delete an experience
// @Tags         Admin - Experiences
// @Produce      json
// @Param        id   path     string  true  "Experience ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/experiences/{id} [delete]
func (h *Handler) DeleteExperience(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", "invalid id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete experience", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Experience deleted successfully", nil)
}
