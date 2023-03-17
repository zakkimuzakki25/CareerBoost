package handler

import (
	"CareerBoost/src/entity"
	"net/http"
	"strconv"

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
			ID:   strconv.FormatUint(uint64(s.ID), 10),
		})
	}

	h.SuccessResponse(ctx, http.StatusOK, "Succes", resp, nil)
}
