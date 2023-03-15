package handler

import (
	"CareerBoost/src/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func (h *handler) getAllMentor(ctx *gin.Context) {
// 	mentorDB := []entity.Mentor{}
// 	err := h.db.Order("rate desc").Find(&mentorDB).Error
// 	if err != nil {
// 		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
// 		return
// 	}

// 	var mentors []entity.MentorRespHome
// 	for _, mentor := range mentorDB {
// 		var resp entity.MentorRespHome
// 		resp.ID = mentor.ID
// 		resp.ProfilePhoto = mentor.ProfilePhoto
// 		resp.Nama = mentor.FullName
// 		resp.Work = mentor.Work
// 		resp.Lokasi = mentor.Lokasi
// 		resp.Rate = mentor.Rate
// 		resp.Deskripsi = mentor.Deskripsi

// 		var skill []entity.Skill
// 		if err := h.db.Model(&mentor).Association("Skill").Find(&skill); err != nil {
// 			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
// 			return
// 		}
// 		var skills []entity.RespSkill
// 		for _, s := range skill {
// 			skills = append(skills, entity.RespSkill{Nama: s.Nama})
// 		}
// 		resp.Skill = skills

// 		var interest []entity.Interest
// 		if err := h.db.Model(&mentor).Association("Interest").Find(&interest); err != nil {
// 			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
// 			return
// 		}
// 		var interests []entity.RespInterest
// 		for _, s := range interest {
// 			interests = append(interests, entity.RespInterest{Nama: s.Nama})
// 		}
// 		resp.Bidang = interests

// 		if err := h.db.Preload("Exp").Where("id = ?", mentor.ID).Take(&mentor).Error; err != nil {
// 			h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
// 			return
// 		}

// 		var exps []entity.ExpResp
// 		for _, exp := range mentor.Exp {
// 			if exp.MentorID == mentor.ID {
// 				exps = append(exps, entity.ExpResp{
// 					Logo:       exp.Logo,
// 					Perusahaan: exp.Perusahaan,
// 					Skill:      exp.Skill,
// 				})
// 			}
// 		}
// 		resp.Exp = exps

// 		mentors = append(mentors, resp)
// 	}

// 	h.SuccessResponse(ctx, http.StatusOK, "Success", mentors, nil)
// }

func (h *handler) getMentorData(ctx *gin.Context) {
	var mentorBody entity.MentorParam
	if err := h.BindParam(ctx, &mentorBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "gagal bindbody", nil)
		return
	}

	var mentorDB entity.Mentor

	err := h.db.Where("id = ?", mentorBody.ID).Take(&mentorDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var resp entity.MentorRespData

	resp.ID = mentorDB.ID
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
	type mentorParam struct {
		ID uint `uri:"mentor_id" gorm:"column:id"`
	}

	var mentorBody mentorParam
	if err := h.BindParam(ctx, &mentorBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request", nil)
		return
	}

	fmt.Println("================================")
	fmt.Println(mentorBody.ID)
	fmt.Println("================================")

	var mentorDB entity.Mentor
	if err := h.db.Preload("Exp").Where("id = ?", mentorBody.ID).Take(&mentorDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var exps []entity.ExpResp
	for _, exp := range mentorDB.Exp {
		if exp.MentorID == mentorDB.ID {
			exps = append(exps, entity.ExpResp{
				Logo:       exp.Logo,
				Perusahaan: exp.Perusahaan,
				Skill:      exp.Skill,
			})
		}
	}

	h.SuccessResponse(ctx, http.StatusOK, "success", exps, nil)
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
	mentorDB.WA = mentorBody.WA
	mentorDB.IG = mentorBody.IG
	mentorDB.Email = mentorBody.Email

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
	var mentorBody entity.Filter
	if err := h.BindParam(ctx, &mentorBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
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
	if mentorBody.Key != "" {
		db = db.Where("full_name LIKE ?", "%"+mentorBody.Key+"%")
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
	if mentorBody.Key != "" {
		db = db.Where("full_name LIKE ?", "%"+mentorBody.Key+"%")
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

func (h *handler) MentorGetPhotoProfile(ctx *gin.Context) {
	type mentorPhoto struct {
		Profile string `json:"photo_profile"`
	}

	var mentorBody entity.MentorReqByID
	if err := h.BindBody(ctx, &mentorBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "gagal bindbody", nil)
		return
	}

	var mentorDB entity.Mentor

	err := h.db.Where("id = ?", mentorBody.ID).Take(&mentorDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	mentorResp := mentorPhoto{
		Profile: mentorDB.ProfilePhoto,
	}

	h.SuccessResponse(ctx, http.StatusOK, "Succes", mentorResp, nil)
}
