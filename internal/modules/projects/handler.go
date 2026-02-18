package projects

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prakoso-id/personal-backend/internal/utils/pagination"
	"github.com/prakoso-id/personal-backend/internal/utils/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetPublicProjects godoc
// @Summary      Public - Get All Projects
// @Description  Retrieve a list of all projects
// @Tags         Public - Projects
// @Produce      json
// @Success      200  {array}   Project
// @Failure      500  {object}  map[string]string
// @Router       /public/projects [get]
func (h *Handler) GetPublicProjects(c *gin.Context) {
	projects, err := h.service.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch projects", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Projects fetched successfully", projects)
}

// GetAdminProjects godoc
// @Summary      Admin - Get All Projects
// @Description  Retrieve a paginated list of all projects
// @Tags         Admin - Projects
// @Produce      json
// @Param        page   query    int  false  "Page number" default(1)
// @Param        limit  query    int  false  "Items per page" default(10)
// @Security     BearerAuth
// @Success      200  {object}  pagination.PaginatedResponse
// @Failure      500  {object}  map[string]string
// @Router       /admin/projects [get]
func (h *Handler) GetAdminProjects(c *gin.Context) {
	p := pagination.FromContext(c)
	projects, err := h.service.GetAllAdmin(p.Page, p.Limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch projects", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Projects fetched successfully", projects)
}

func (h *Handler) GetPublicProjectByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", "invalid id")
		return
	}

	project, err := h.service.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch project", err.Error())
		return
	}
	if project == nil {
		response.Error(c, http.StatusNotFound, "Project not found", "project not found")
		return
	}
	response.Success(c, http.StatusOK, "Project fetched successfully", project)
}

// CreateProject godoc
// @Summary      Admin - Create Project
// @Description  Create a new project
// @Tags         Admin - Projects
// @Accept       json
// @Produce      json
// @Param        request body CreateProjectRequest true "Project Data"
// @Security     BearerAuth
// @Success      201  {object}  Project
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/projects [post]
func (h *Handler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	project, err := h.service.Create(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create project", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "Project created successfully", project)
}

// UpdateProject godoc
// @Summary      Admin - Update Project
// @Description  Update an existing project
// @Tags         Admin - Projects
// @Accept       json
// @Produce      json
// @Param        id   path     string  true  "Project ID"
// @Param        request body UpdateProjectRequest true "Project Data"
// @Security     BearerAuth
// @Success      200  {object}  Project
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/projects/{id} [put]
func (h *Handler) UpdateProject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", "invalid id")
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	project, err := h.service.Update(id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update project", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Project updated successfully", project)
}

// Delete godoc
// @Summary      Admin - Delete Project
// @Description  Delete a project
// @Tags         Admin - Projects
// @Produce      json
// @Param        id   path     string  true  "Project ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/projects/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", "invalid id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete project", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Project deleted successfully", nil)
}
