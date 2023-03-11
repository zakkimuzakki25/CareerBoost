package handler

import (
	"CareerBoost/src/entity"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (h *handler) SuccessResponse(c *gin.Context, code int64, message string, data interface{}, pagination *entity.PaginationParam) {
	c.JSON(int(code), entity.HTTPResponse{
		Message:    message,
		IsSuccess:  true,
		Data:       data,
		Pagination: pagination,
	})
}

func (h *handler) ErrorResponse(c *gin.Context, code int64, message string, data interface{}) {
	c.JSON(int(code), entity.HTTPResponse{
		Message:    message,
		IsSuccess:  false,
		Data:       data,
		Pagination: nil,
	})
}

func (h *handler) BindBody(c *gin.Context, body interface{}) interface{} {
	return c.ShouldBindWith(body, binding.JSON)
}

func (h *handler) BindParam(c *gin.Context, param interface{}) error {
	if err := c.ShouldBindUri(param); err != nil {
		return err
	}

	return c.ShouldBindWith(param, binding.Query)
}
