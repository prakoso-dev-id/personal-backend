package skills

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prakoso-id/personal-backend/internal/utils/pagination"
	"github.com/prakoso-id/personal-backend/internal/utils/response"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// GetAll godoc
// @Summary      Public - Get All Skills
// @Description  Retrieve a list of all skills
// @Tags         Public - Skills
// @Produce      json
// @Success      200  {array}   Skill
// @Failure      500  {object}  map[string]string
// @Router       /public/skills [get]
func (h *Handler) GetAll(c *gin.Context) {
	p := pagination.FromContext(c)
	var skills []Skill
	var total int64

	if err := h.db.Model(&Skill{}).Count(&total).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to count skills", err.Error())
		return
	}

	if err := h.db.Scopes(pagination.Paginate(p)).Find(&skills).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch skills", err.Error())
		return
	}

	res := pagination.NewResponse(skills, total, p)
	response.Success(c, http.StatusOK, "Skills fetched successfully", res)
}

type CreateSkillRequest struct {
	Name     string `json:"name" binding:"required"`
	Category string `json:"category" binding:"required"`
	IconURL  string `json:"icon_url"`
}

// Create godoc
// @Summary      Admin - Create Skill
// @Description  Create a new skill
// @Tags         Admin - Skills
// @Accept       json
// @Produce      json
// @Param        request body CreateSkillRequest true "Skill Data"
// @Security     BearerAuth
// @Success      201  {object}  Skill
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/skills [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	skill := Skill{
		Name:     req.Name,
		Category: req.Category,
		IconURL:  req.IconURL,
	}

	if err := h.db.Create(&skill).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create skill", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Skill created successfully", skill)
}

type UpdateSkillRequest struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	IconURL  string `json:"icon_url"`
}

// Update godoc
// @Summary      Admin - Update Skill
// @Description  Update an existing skill
// @Tags         Admin - Skills
// @Accept       json
// @Produce      json
// @Param        id   path     string  true  "Skill ID"
// @Param        request body UpdateSkillRequest true "Skill Data"
// @Security     BearerAuth
// @Success      200  {object}  Skill
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/skills/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	var skill Skill
	if err := h.db.First(&skill, "id = ?", id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Skill not found", "skill not found")
		return
	}

	if req.Name != "" {
		skill.Name = req.Name
	}
	if req.Category != "" {
		skill.Category = req.Category
	}
	if req.IconURL != "" {
		skill.IconURL = req.IconURL
	}

	if err := h.db.Save(&skill).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update skill", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Skill updated successfully", skill)
}

// Delete godoc
// @Summary      Admin - Delete Skill
// @Description  Delete a skill
// @Tags         Admin - Skills
// @Produce      json
// @Param        id   path     string  true  "Skill ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/skills/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&Skill{}, "id = ?", id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete skill", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Skill deleted successfully", nil)
}
