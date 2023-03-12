package handler

import (
	"CareerBoost/src/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) getAllMagang(ctx *gin.Context) {
	var magangParam entity.MagangParam
	if err := h.BindParam(ctx, &magangParam); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	magangParam.FormatPagination()

	var magangDB []entity.Magang

	if err := h.db.
		Model(entity.Magang{}).
		Limit(int(magangParam.Limit)).
		Offset(int(magangParam.Offset)).
		Find(&magangDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var totalElements int64

	if err := h.db.
		Model(entity.Magang{}).
		Limit(int(magangParam.Limit)).
		Offset(int(magangParam.Offset)).
		Count(&totalElements).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	magangParam.ProcessPagination(totalElements)

	var magangs []entity.MagangRespData
	for _, magang := range magangDB {

		var count int64
		err := h.db.Model(&entity.User{}).Where("magang_id = ?", magang.ID).Count(&count).Error
		if err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		var resp entity.MagangRespData
		resp.CreatedAt = magang.CreatedAt
		resp.Logo = magang.Logo
		resp.Perusahaan = magang.Perusahaan
		resp.Lokasi = magang.Lokasi
		resp.Deskripsi = magang.Deskripsi
		resp.Applied = uint(count)
		resp.Rate = magang.Rate
		resp.Fee = magang.Fee
		resp.JangkaWaktu = magang.JangkaWaktu
		resp.StatusMagang = magang.StatusMagang

		var skill []entity.Skill
		if err := h.db.Model(&magang).Association("Skill").Find(&skill); err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var skills []entity.RespSkill
		for _, s := range skill {
			skills = append(skills, entity.RespSkill{Nama: s.Nama})
		}
		resp.Skill = skills

		var interest []entity.Interest
		if err := h.db.Model(&magang).Association("Interest").Find(&interest); err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var interests []entity.RespInterest
		for _, s := range interest {
			interests = append(interests, entity.RespInterest{Nama: s.Nama})
		}
		resp.Interest = interests

		magangs = append(magangs, resp)
	}

	if magangs != nil {
		h.SuccessResponse(ctx, http.StatusOK, "Success", magangs, &magangParam.PaginationParam)
	} else {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "Lowongan magang tidak ditemukan", nil)
		return
	}
}

func (h *handler) getMagangRecomendation(ctx *gin.Context) {
	var magangParam entity.MagangParam
	if err := h.BindParam(ctx, &magangParam); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	magangParam.FormatPagination()

	var magangDB []entity.Magang

	if err := h.db.
		Model(entity.Magang{}).
		Order("rate desc").
		Limit(int(magangParam.Limit)).
		Offset(int(magangParam.Offset)).
		Find(&magangDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var totalElements int64

	if err := h.db.
		Model(entity.Magang{}).
		Limit(int(magangParam.Limit)).
		Offset(int(magangParam.Offset)).
		Count(&totalElements).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	magangParam.ProcessPagination(totalElements)

	var magangs []entity.MagangRespData
	for _, magang := range magangDB {

		var count int64
		err := h.db.Model(&entity.User{}).Where("magang_id = ?", magang.ID).Count(&count).Error
		if err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		var resp entity.MagangRespData
		resp.CreatedAt = magang.CreatedAt
		resp.Logo = magang.Logo
		resp.Perusahaan = magang.Perusahaan
		resp.Lokasi = magang.Lokasi
		resp.Deskripsi = magang.Deskripsi
		resp.Applied = uint(count)
		resp.Rate = magang.Rate
		resp.Fee = magang.Fee
		resp.JangkaWaktu = magang.JangkaWaktu
		resp.StatusMagang = magang.StatusMagang

		var skill []entity.Skill
		if err := h.db.Model(&magang).Association("Skill").Find(&skill); err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var skills []entity.RespSkill
		for _, s := range skill {
			skills = append(skills, entity.RespSkill{Nama: s.Nama})
		}
		resp.Skill = skills

		var interest []entity.Interest
		if err := h.db.Model(&magang).Association("Interest").Find(&interest); err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var interests []entity.RespInterest
		for _, s := range interest {
			interests = append(interests, entity.RespInterest{Nama: s.Nama})
		}
		resp.Interest = interests

		magangs = append(magangs, resp)
	}

	if magangs != nil {
		h.SuccessResponse(ctx, http.StatusOK, "Success", magangs, &magangParam.PaginationParam)
	} else {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "Lowongan magang tidak ditemukan", nil)
		return
	}
}

func (h *handler) addNewMagang(ctx *gin.Context) {
	var magangBody entity.MagangAdd
	if err := h.BindBody(ctx, &magangBody); err != nil {
		fmt.Println(err)
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request register", nil)
		return
	}

	var magangDB entity.Magang
	magangDB.Logo = magangBody.Logo
	magangDB.Perusahaan = magangBody.Perusahaan
	magangDB.Lokasi = magangBody.Lokasi
	magangDB.Deskripsi = magangBody.Deskripsi
	magangDB.Rate = magangBody.Rate
	magangDB.Fee = magangBody.Fee
	magangDB.JangkaWaktu = magangBody.JangkaWaktu
	magangDB.StatusMagang = magangBody.StatusMagang

	var skills []entity.Skill
	if err := h.db.Find(&skills, magangBody.Skill).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "interest not found", nil)
		return
	}

	var interests []entity.Interest
	if err := h.db.Find(&interests, magangBody.Interest).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "interest not found", nil)
		return
	}

	if err := h.db.Create(&magangDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "add mentor gagal", nil)
		return
	}

	if err := h.db.Model(&magangDB).Association("Interest").Append(interests); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal nambah Interest", nil)
		return
	}

	if err := h.db.Model(&magangDB).Association("Skill").Append(skills); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "skill not added", nil)
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "tambah magang sukses", nil, nil)
}

func (h *handler) getMagangFilter(ctx *gin.Context) {
	var magangBody entity.Filter
	if err := h.BindBody(ctx, &magangBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "gagal init body magang", nil)
		return
	}

	var magangParam entity.MagangParam
	if err := h.BindParam(ctx, &magangParam); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	magangParam.FormatPagination()

	var magangDB []entity.Magang

	db := h.db.Model(entity.Magang{}).
		Limit(int(magangParam.Limit)).
		Offset(int(magangParam.Offset))

	if len(magangBody.InterestID) > 0 {
		db = db.Joins("JOIN magangs_interest ON magangs_interest.magang_id = magangs.id").
			Where("magangs_interest.interest_id IN (?)", magangBody.InterestID)
	}
	if magangBody.Key != "" {
		db = db.Where("perusahaan LIKE ?", "%"+magangBody.Key+"%")
	}

	if err := db.Find(&magangDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var totalElements int64

	db = h.db.Model(entity.Magang{}).
		Limit(int(magangParam.Limit)).
		Offset(int(magangParam.Offset))

	if len(magangBody.InterestID) > 0 {
		db = db.Joins("JOIN magangs_interest ON magangs_interest.magang_id = magangs.id").
			Where("magangs_interest.interest_id IN (?)", magangBody.InterestID)
	}
	if magangBody.Key != "" {
		db = db.Where("perusahaan LIKE ?", "%"+magangBody.Key+"%")
	}

	if err := db.Count(&totalElements).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	magangParam.ProcessPagination(totalElements)

	var magangs []entity.MagangRespData
	for _, magang := range magangDB {

		var count int64
		err := h.db.Model(&entity.User{}).Where("magang_id = ?", magang.ID).Count(&count).Error
		if err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		var resp entity.MagangRespData
		resp.CreatedAt = magang.CreatedAt
		resp.Logo = magang.Logo
		resp.Perusahaan = magang.Perusahaan
		resp.Lokasi = magang.Lokasi
		resp.Deskripsi = magang.Deskripsi
		resp.Applied = uint(count)
		resp.Rate = magang.Rate
		resp.Fee = magang.Fee
		resp.JangkaWaktu = magang.JangkaWaktu
		resp.StatusMagang = magang.StatusMagang

		var skill []entity.Skill
		if err := h.db.Model(&magang).Association("Skill").Find(&skill); err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var skills []entity.RespSkill
		for _, s := range skill {
			skills = append(skills, entity.RespSkill{Nama: s.Nama})
		}
		resp.Skill = skills

		var interest []entity.Interest
		if err := h.db.Model(&magang).Association("Interest").Find(&interest); err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var interests []entity.RespInterest
		for _, s := range interest {
			interests = append(interests, entity.RespInterest{Nama: s.Nama})
		}
		resp.Interest = interests

		magangs = append(magangs, resp)
	}

	if magangs != nil {
		h.SuccessResponse(ctx, http.StatusOK, "Success", magangs, &magangParam.PaginationParam)
	} else {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "Lowongan magang tidak ditemukan", nil)
		return
	}
}

func (h *handler) getMagangData(ctx *gin.Context) {
	var magangBody entity.MagangReqByID
	if err := h.BindBody(ctx, &magangBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "gagal init body", nil)
		return
	}

	var magangDB entity.Magang

	err := h.db.Where("id = ?", magangBody.ID).Take(&magangDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var resp entity.MagangRespData

	errrr := h.db.Where("id = ?", magangBody.ID).Take(&magangDB).Error
	if errrr != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var count int64
	errr := h.db.Model(&entity.User{}).Where("magang_id = ?", magangDB.ID).Count(&count).Error
	if errr != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp.CreatedAt = magangDB.CreatedAt
	resp.Logo = magangDB.Logo
	resp.Perusahaan = magangDB.Perusahaan
	resp.Lokasi = magangDB.Lokasi
	resp.Deskripsi = magangDB.Deskripsi
	resp.Applied = uint(count)
	resp.Rate = magangDB.Rate
	resp.Fee = magangDB.Fee
	resp.JangkaWaktu = magangDB.JangkaWaktu
	resp.StatusMagang = magangDB.StatusMagang

	var skill []entity.Skill
	if err := h.db.Model(&magangDB).Association("Skill").Find(&skill); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	var skills []entity.RespSkill
	for _, s := range skill {
		skills = append(skills, entity.RespSkill{Nama: s.Nama})
	}
	resp.Skill = skills

	var interest []entity.Interest
	if err := h.db.Model(&magangDB).Association("Interest").Find(&interest); err != nil {
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
