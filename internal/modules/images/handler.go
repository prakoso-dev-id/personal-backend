package images

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

type UploadRequest struct {
	// No other fields required for now
}

// Upload godoc
// @Summary      Admin - Upload Image
// @Description  Upload an image file to storage. Returns file details to be used in other endpoints.
// @Tags         Admin - Images
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Image File"
// @Security     BearerAuth
// @Success      201  {object}  ImageUploadResult
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/images/upload [post]
func (h *Handler) Upload(c *gin.Context) {
	// Just check if file exists
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "File is required", "file is required")
		return
	}

	result, err := h.service.UploadFile(file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to upload image", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Image uploaded successfully", result)
}

// Delete godoc
// @Summary      Admin - Delete Image
// @Description  Delete an image
// @Tags         Admin - Images
// @Produce      json
// @Param        id   path     string  true  "Image ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/images/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid image ID", "invalid image id")
		return
	}

	if err := h.service.DeleteImage(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete image", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Image deleted successfully", nil)
}
