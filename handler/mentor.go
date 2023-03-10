package handler

import (
	"CareerBoost/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) getAllMentor(ctx *gin.Context) {
	var mentorDB []entity.Mentor

	err := h.db.Find(&mentorDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "Succes", mentorDB, nil)
}

func (h *handler) getMentorData(ctx *gin.Context) {
	var mentorBody entity.MentorReqByID
	if err := h.BindBody(ctx, &mentorBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "gagal init body", nil)
		return
	}

	var mentorDB entity.Mentor

	err := h.db.Where("id = ?", mentorBody.ID).Take(&mentorDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var resp entity.MentorRespData

	resp.FullName = mentorDB.FullName
	resp.Skill = mentorDB.Skill
	resp.Lokasi = mentorDB.Lokasi
	resp.Interest = mentorDB.Interest
	resp.Deskripsi = mentorDB.Deskripsi
	resp.Rate = mentorDB.Rate
	resp.Fee = mentorDB.Fee

	h.SuccessResponse(ctx, http.StatusOK, "Succes", resp, nil)
}

func (h *handler) getMentorExp(ctx *gin.Context) {
	var mentorBody entity.MentorReqByID
	if err := h.BindBody(ctx, &mentorBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request register", nil)
		return
	}

	var mentorDB entity.Mentor

	err := h.db.Where("id = ?", mentorBody.ID).Take(&mentorDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "Succes", mentorDB, nil)
}

func (h *handler) addNewMentor(ctx *gin.Context) {
	var mentorBody entity.MentorAdd
	if err := h.BindBody(ctx, &mentorBody); err != nil {
		fmt.Println(err)
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request register", nil)
		return
	}

	var mentorDB entity.Mentor
	mentorDB.ProfilePhoto = mentorBody.ProfilePhoto
	mentorDB.FullName = mentorBody.FullName
	mentorDB.Lokasi = mentorBody.Lokasi
	mentorDB.Deskripsi = mentorBody.Deskripsi
	mentorDB.Rate = mentorBody.Rate
	mentorDB.Fee = mentorBody.Fee

	var exps []entity.Exp
	for _, exp := range mentorBody.Exp {
		exps = append(exps, entity.Exp{
			Logo:       exp.Logo,
			Skill:      exp.Skill,
			Perusahaan: exp.Perusahaan,
		})
	}

	var skills []entity.Skill
	if err := h.db.Find(&skills, mentorBody.Skill).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "interest not found", nil)
		return
	}

	var interests []entity.Interest
	if err := h.db.Find(&interests, mentorBody.Interest).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "interest not found", nil)
		return
	}

	if err := h.db.Create(&mentorDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "add mentor gagal", nil)
		return
	}

	if err := h.db.Model(&mentorDB).Association("Exp").Append(exps); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal nambah exp", nil)
		return
	}

	if err := h.db.Model(&mentorDB).Association("Interest").Append(interests); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal nambah Interest", nil)
		return
	}

	if err := h.db.Model(&mentorDB).Association("Skill").Append(skills); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "skill not added", nil)
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "tambah mentor sukses", nil, nil)
}
