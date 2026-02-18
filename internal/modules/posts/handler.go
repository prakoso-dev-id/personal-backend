package posts

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

// GetPublicPosts godoc
// @Summary      Public - Get All Posts
// @Description  Retrieve a list of all published posts
// @Tags         Public - Posts
// @Produce      json
// @Success      200  {array}   Post
// @Failure      500  {object}  map[string]string
// @Router       /public/posts [get]
func (h *Handler) GetPublicPosts(c *gin.Context) {
	posts, err := h.service.GetAll(true)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch posts", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Posts fetched successfully", posts)
}

// GetAdminPosts godoc
// @Summary      Admin - Get All Posts
// @Description  Retrieve a paginated list of all posts (including unpublished)
// @Tags         Admin - Posts
// @Produce      json
// @Param        page   query    int  false  "Page number" default(1)
// @Param        limit  query    int  false  "Items per page" default(10)
// @Security     BearerAuth
// @Success      200  {object}  pagination.PaginatedResponse
// @Failure      500  {object}  map[string]string
// @Router       /admin/posts [get]
func (h *Handler) GetAdminPosts(c *gin.Context) {
	p := pagination.FromContext(c)
	posts, err := h.service.GetAllAdmin(p.Page, p.Limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch posts", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Posts fetched successfully", posts)
}

func (h *Handler) GetPublicPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	post, err := h.service.GetBySlug(slug)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch post", err.Error())
		return
	}
	if post == nil {
		response.Error(c, http.StatusNotFound, "Post not found", "post not found")
		return
	}
	if !post.IsPublished {
		response.Error(c, http.StatusNotFound, "Post not found", "post not found")
		return
	}
	response.Success(c, http.StatusOK, "Post fetched successfully", post)
}

// CreatePost godoc
// @Summary      Admin - Create Post
// @Description  Create a new post
// @Tags         Admin - Posts
// @Accept       json
// @Produce      json
// @Param        request body CreatePostRequest true "Post Data"
// @Security     BearerAuth
// @Success      201  {object}  Post
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/posts [post]
func (h *Handler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	post, err := h.service.Create(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create post", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "Post created successfully", post)
}

// UpdatePost godoc
// @Summary      Admin - Update Post
// @Description  Update an existing post
// @Tags         Admin - Posts
// @Accept       json
// @Produce      json
// @Param        id   path     string  true  "Post ID"
// @Param        request body UpdatePostRequest true "Post Data"
// @Security     BearerAuth
// @Success      200  {object}  Post
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/posts/{id} [put]
func (h *Handler) UpdatePost(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", "invalid id")
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	post, err := h.service.Update(id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update post", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Post updated successfully", post)
}

// DeletePost godoc
// @Summary      Admin - Delete Post
// @Description  Delete a post
// @Tags         Admin - Posts
// @Produce      json
// @Param        id   path     string  true  "Post ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/posts/{id} [delete]
func (h *Handler) DeletePost(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ID", "invalid id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete post", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Post deleted successfully", nil)
}
