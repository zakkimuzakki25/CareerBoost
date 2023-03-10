package handler

import (
	"CareerBoost/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) getAllMentor(ctx *gin.Context) {
	var postParam entity.MentorParam
	if err := h.BindParam(ctx, &postParam); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	postParam.FormatPagination()

	var mentorDB []entity.Mentor

	if err := h.db.
		Model(entity.Mentor{}).
		Limit(int(postParam.Limit)).
		Offset(int(postParam.Offset)).
		Find(&mentorDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var totalElements int64

	if err := h.db.
		Model(entity.Mentor{}).
		Limit(int(postParam.Limit)).
		Offset(int(postParam.Offset)).
		Count(&totalElements).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	postParam.ProcessPagination(totalElements)

	var mentors []entity.MentorRespData
	for _, mentor := range mentorDB {

		var resp entity.MentorRespData
		resp.ProfilePhoto = mentor.ProfilePhoto
		resp.FullName = mentor.FullName
		resp.Lokasi = mentor.Lokasi
		resp.Deskripsi = mentor.Deskripsi
		resp.Rate = mentor.Rate
		resp.Fee = mentor.Fee

		var skill []entity.Skill
		if err := h.db.Model(&mentor).Association("Skill").Find(&skill); err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var skills []entity.RespSkill
		for _, s := range skill {
			skills = append(skills, entity.RespSkill{Nama: s.Nama})
		}
		resp.Skill = skills

		var interest []entity.Interest
		if err := h.db.Model(&mentor).Association("Interest").Find(&interest); err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var interests []entity.RespInterest
		for _, s := range interest {
			interests = append(interests, entity.RespInterest{Nama: s.Nama})
		}
		resp.Interest = interests

		mentors = append(mentors, resp)
	}

	h.SuccessResponse(ctx, http.StatusOK, "Success", mentors, &postParam.PaginationParam)
}

func (h *handler) getMentorRekomendation(ctx *gin.Context) {
	var mentorDB []entity.Mentor

	err := h.db.Order("rate desc").Limit(4).Find(&mentorDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var mentors []entity.MentorRespData
	for _, mentor := range mentorDB {

		var resp entity.MentorRespData
		resp.ProfilePhoto = mentor.ProfilePhoto
		resp.FullName = mentor.FullName
		resp.Lokasi = mentor.Lokasi
		resp.Deskripsi = mentor.Deskripsi
		resp.Rate = mentor.Rate
		resp.Fee = mentor.Fee

		var skill []entity.Skill
		if err := h.db.Model(&mentor).Association("Skill").Find(&skill); err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var skills []entity.RespSkill
		for _, s := range skill {
			skills = append(skills, entity.RespSkill{Nama: s.Nama})
		}
		resp.Skill = skills

		var interest []entity.Interest
		if err := h.db.Model(&mentor).Association("Interest").Find(&interest); err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var interests []entity.RespInterest
		for _, s := range interest {
			interests = append(interests, entity.RespInterest{Nama: s.Nama})
		}
		resp.Interest = interests

		mentors = append(mentors, resp)
	}

	h.SuccessResponse(ctx, http.StatusOK, "Success", mentors, nil)
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

	resp.ProfilePhoto = mentorDB.ProfilePhoto
	resp.FullName = mentorDB.FullName
	resp.Lokasi = mentorDB.Lokasi
	resp.Deskripsi = mentorDB.Deskripsi
	resp.Rate = mentorDB.Rate
	resp.Fee = mentorDB.Fee

	var skill []entity.Skill
	if err := h.db.Model(&mentorDB).Association("Skill").Find(&skill); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	var skills []entity.RespSkill
	for _, s := range skill {
		skills = append(skills, entity.RespSkill{Nama: s.Nama})
	}
	resp.Skill = skills

	var interest []entity.Interest
	if err := h.db.Model(&mentorDB).Association("Interest").Find(&interest); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	var interests []entity.RespInterest
	for _, s := range interest {
		interests = append(interests, entity.RespInterest{Nama: s.Nama})
	}
	resp.Interest = interests

	h.SuccessResponse(ctx, http.StatusOK, "Success", resp, nil)
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

func (h *handler) getMentorFilter(ctx *gin.Context) {
	var mentorBody entity.MentorFilter
	if err := h.BindBody(ctx, &mentorBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "gagal init body", nil)
		return
	}

	var postParam entity.MentorParam
	if err := h.BindParam(ctx, &postParam); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	postParam.FormatPagination()

	var mentorDB []entity.Mentor

	db := h.db.Model(entity.Mentor{}).
		Limit(int(postParam.Limit)).
		Offset(int(postParam.Offset))

	if len(mentorBody.InterestID) > 0 {
		db = db.Joins("JOIN mentors_interest ON mentors_interest.mentor_id = mentors.id").
			Where("mentors_interest.interest_id IN (?)", mentorBody.InterestID)
	}

	if err := db.Find(&mentorDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var totalElements int64

	db = h.db.Model(entity.Mentor{}).
		Limit(int(postParam.Limit)).
		Offset(int(postParam.Offset))

	if len(mentorBody.InterestID) > 0 {
		db = db.Joins("JOIN mentors_interest ON mentors_interest.mentor_id = mentors.id").
			Where("mentors_interest.interest_id IN (?)", mentorBody.InterestID)
	}

	if err := db.Count(&totalElements).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	postParam.ProcessPagination(totalElements)

	var mentors []entity.MentorRespData
	for _, mentor := range mentorDB {

		var resp entity.MentorRespData
		resp.ProfilePhoto = mentor.ProfilePhoto
		resp.FullName = mentor.FullName
		resp.Lokasi = mentor.Lokasi
		resp.Deskripsi = mentor.Deskripsi
		resp.Rate = mentor.Rate
		resp.Fee = mentor.Fee

		var skill []entity.Skill
		if err := h.db.Model(&mentor).Association("Skill").Find(&skill); err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var skills []entity.RespSkill
		for _, s := range skill {
			skills = append(skills, entity.RespSkill{Nama: s.Nama})
		}
		resp.Skill = skills

		var interest []entity.Interest
		if err := h.db.Model(&mentor).Association("Interest").Find(&interest); err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var interests []entity.RespInterest
		for _, s := range interest {
			interests = append(interests, entity.RespInterest{Nama: s.Nama})
		}
		resp.Interest = interests

		mentors = append(mentors, resp)
	}

	if mentors != nil {
		h.SuccessResponse(ctx, http.StatusOK, "Success", mentors, &postParam.PaginationParam)
	} else {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "Mentor Tidak Ditemukan", nil)
		return
	}
}
