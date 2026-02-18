package contact

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prakoso-id/personal-backend/internal/utils/response"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

type CreateMessageRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject"`
	Message string `json:"message" binding:"required"`
}

// CreateMessage godoc
// @Summary      Public - Send Contact Message
// @Description  Send a contact message
// @Tags         Public - Contact
// @Accept       json
// @Produce      json
// @Param        request body CreateMessageRequest true "Message"
// @Success      201  {object}  ContactMessage
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /public/contact [post]
func (h *Handler) CreateMessage(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	msg := ContactMessage{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Message: req.Message,
	}

	if err := h.db.Create(&msg).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to send message", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Message sent successfully", msg)
}

// GetAllMessages godoc
// @Summary      Admin - Get All Messages
// @Description  Retrieve all contact messages
// @Tags         Admin - Contact
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   ContactMessage
// @Failure      500  {object}  map[string]string
// @Router       /admin/messages [get]
func (h *Handler) GetAllMessages(c *gin.Context) {
	var msgs []ContactMessage
	if err := h.db.Order("created_at DESC").Find(&msgs).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch messages", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Messages fetched successfully", msgs)
}
