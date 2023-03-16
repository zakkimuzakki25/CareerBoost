package handler

import (
	"CareerBoost/src/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) getAllInterest(ctx *gin.Context) {
	var interest []entity.Interest
	err := h.db.Model(&interest).Find(&interest).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "Interest not found", nil)
		return
	}

	var resp []entity.RespInterestWithID
	for _, s := range interest {
		resp = append(resp, entity.RespInterestWithID{
			Nama: s.Nama,
			ID:   s.ID,
		})
	}

	h.SuccessResponse(ctx, http.StatusOK, "Succes", resp, nil)
}
