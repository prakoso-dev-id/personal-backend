package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Meta       Meta        `json:"meta"`
}

type Meta struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
	TotalData   int64 `json:"total_data"`
	Limit       int   `json:"limit"`
}

func FromContext(c *gin.Context) Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return Pagination{
		Page:  page,
		Limit: limit,
	}
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

func Paginate(p Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.Offset()).Limit(p.Limit)
	}
}

func NewResponse(data interface{}, total int64, p Pagination) PaginatedResponse {
	totalPage := int(total) / p.Limit
	if int(total)%p.Limit != 0 {
		totalPage++
	}

	return PaginatedResponse{
		Data: data,
		Meta: Meta{
			CurrentPage: p.Page,
			TotalPage:   totalPage,
			TotalData:   total,
			Limit:       p.Limit,
		},
	}
}
