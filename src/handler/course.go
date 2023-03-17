package handler

import (
	"CareerBoost/src/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) addNewCourse(ctx *gin.Context) {
	var courseBody entity.CourseAdd
	if err := h.BindBody(ctx, &courseBody); err != nil {
		fmt.Println(err)
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request", nil)
		return
	}

	var courseDB entity.Course
	courseDB.Foto = courseBody.Foto
	courseDB.Judul = courseBody.Judul
	courseDB.Deskripsi = courseBody.Deskripsi
	courseDB.Intro = courseBody.Intro
	courseDB.Rate = courseBody.Rate
	courseDB.Price = courseBody.Price
	courseDB.InterestID = courseBody.InterestID
	courseDB.Vote = courseBody.Vote

	if err := h.db.Create(&courseDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var skill []entity.Skill
	if err := h.db.Find(&skill, courseBody.Skill).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "interest not found", nil)
		return
	}

	if err := h.db.Model(&courseDB).Association("Skill").Append(skill); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	for _, playlist := range courseBody.Playlist {
		var playl entity.Playlist

		playl.Nama = playlist.Nama
		playl.CourseID = courseDB.ID
		playl.Durasi = playlist.Durasi

		var videos []entity.Video

		for _, video := range playlist.Video {

			videos = append(videos, entity.Video{
				Link:       video.Link,
				Judul:      video.Judul,
				Durasi:     video.Durasi,
				PlaylistID: playl.ID,
			})
		}

		playl.Video = videos

		if err := h.db.Create(&playl).Error; err != nil {
			h.ErrorResponse(ctx, http.StatusInternalServerError, "failed to add playlist", nil)
			return
		}
	}

	h.SuccessResponse(ctx, http.StatusOK, "Course berhasil ditambahkan", nil, nil)
}

func (h *handler) getAllCourse(ctx *gin.Context) {

	user, exist := ctx.Get("user")
	if !exist {
		h.ErrorResponse(ctx, http.StatusBadRequest, "Unauthorized", nil)
		return
	}

	claims, ok := user.(entity.UserClaims)
	if !ok {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid token", nil)
		return
	}

	userID := claims.ID

	var courseParam entity.CourseParam
	if err := h.BindParam(ctx, &courseParam); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	courseParam.FormatPagination()

	var courseSearch entity.CourseSearch
	if err := h.BindParam(ctx, &courseSearch); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "gagal init body course", nil)
		return
	}

	var courseBody []entity.Course

	db := h.db.Model(entity.Course{}).
		Limit(int(courseParam.Limit)).
		Offset(int(courseParam.Offset))

	if courseSearch.Key != "" {
		db = db.Where("judul LIKE ?", "%"+courseSearch.Key+"%")
	}

	if err := db.
		Joins("JOIN interests ON courses.interest_id = interests.id").
		Joins("JOIN users_interest ON interests.id = users_interest.interest_id").
		Where("users_interest.user_id = ? AND courses.interest_id = users_interest.interest_id", userID).
		Where("users_interest.user_id = ?", userID).
		Limit(int(courseParam.Limit)).
		Offset(int(courseParam.Offset)).
		Find(&courseBody).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var totalElements int64

	if err := h.db.
		Table("courses").
		Joins("JOIN interests ON courses.interest_id = interests.id").
		Joins("JOIN users_interest ON interests.id = users_interest.interest_id").
		Where("users_interest.user_id = ? AND courses.interest_id = users_interest.interest_id", userID).
		Count(&totalElements).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	courseParam.ProcessPagination(totalElements)

	type resp struct {
		ID        uint    `json:"id"`
		Foto      string  `json:"foto"`
		Judul     string  `json:"judul"`
		Deskripsi string  `json:"deskripsi"`
		Rate      float32 `json:"rate"`
		Vote      uint    `json:"vote"`
		Price     float32 `json:"price"`
	}

	var courses []resp
	for _, course := range courseBody {

		var resps resp
		resps.ID = course.ID
		resps.Vote = course.Vote
		resps.Foto = course.Foto
		resps.Judul = course.Judul
		resps.Deskripsi = course.Deskripsi
		resps.Rate = course.Rate
		resps.Price = course.Price

		courses = append(courses, resps)
	}

	h.SuccessResponse(ctx, http.StatusOK, "Success", courses, &courseParam.PaginationParam)
}

func (h *handler) getAllCourseHome(ctx *gin.Context) {
	var courseParam entity.CourseParam
	if err := h.BindParam(ctx, &courseParam); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	courseParam.FormatPagination()

	var courseSearch entity.CourseSearch
	if err := h.BindParam(ctx, &courseSearch); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "gagal init body course", nil)
		return
	}

	var courseBody []entity.Course

	db := h.db.Model(entity.Course{}).
		Limit(int(courseParam.Limit)).
		Offset(int(courseParam.Offset))

	if courseSearch.Key != "" {
		db = db.Where("judul LIKE ?", "%"+courseSearch.Key+"%")
	}

	if err := db.
		Order("rate desc").Limit(8).
		Limit(int(courseParam.Limit)).
		Offset(int(courseParam.Offset)).
		Find(&courseBody).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var totalElements int64

	if err := h.db.
		Table("courses").
		Order("rate desc").Limit(8).
		Count(&totalElements).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	courseParam.ProcessPagination(totalElements)

	type resp struct {
		ID        uint    `json:"id"`
		Foto      string  `json:"foto"`
		Judul     string  `json:"judul"`
		Deskripsi string  `json:"deskripsi"`
		Rate      float32 `json:"rate"`
		Vote      uint    `json:"vote"`
		Price     float32 `json:"price"`
	}

	var courses []resp
	for _, course := range courseBody {

		var resps resp
		resps.ID = course.ID
		resps.Vote = course.Vote
		resps.Foto = course.Foto
		resps.Judul = course.Judul
		resps.Deskripsi = course.Deskripsi
		resps.Rate = course.Rate
		resps.Price = course.Price

		courses = append(courses, resps)
	}

	h.SuccessResponse(ctx, http.StatusOK, "Success", courses, &courseParam.PaginationParam)
}

func (h *handler) getCourseRekomendasi(ctx *gin.Context) {

	user, exist := ctx.Get("user")
	if !exist {
		h.ErrorResponse(ctx, http.StatusBadRequest, "Unauthorized", nil)
		return
	}

	claims, ok := user.(entity.UserClaims)
	if !ok {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid token", nil)
		return
	}

	userID := claims.ID

	var courseParam entity.CourseParam
	if err := h.BindParam(ctx, &courseParam); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	courseParam.FormatPagination()

	var courseBody []entity.Course

	db := h.db.Model(entity.Course{}).
		Limit(int(courseParam.Limit)).
		Offset(int(courseParam.Offset))

	if err := db.
		Order("rate desc").Limit(4).
		Joins("JOIN interests ON courses.interest_id = interests.id").
		Joins("JOIN users_interest ON interests.id = users_interest.interest_id").
		Where("users_interest.user_id = ? AND courses.interest_id = users_interest.interest_id", userID).
		Where("users_interest.user_id = ?", userID).
		Limit(int(courseParam.Limit)).
		Offset(int(courseParam.Offset)).
		Find(&courseBody).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var totalElements int64

	if err := h.db.
		Table("courses").
		Order("rate desc").Limit(4).
		Joins("JOIN interests ON courses.interest_id = interests.id").
		Joins("JOIN users_interest ON interests.id = users_interest.interest_id").
		Where("users_interest.user_id = ? AND courses.interest_id = users_interest.interest_id", userID).
		Count(&totalElements).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	courseParam.ProcessPagination(totalElements)

	type resp struct {
		ID        uint    `json:"id"`
		Foto      string  `json:"foto"`
		Judul     string  `json:"judul"`
		Deskripsi string  `json:"deskripsi"`
		Rate      float32 `json:"rate"`
		Vote      uint    `json:"vote"`
		Price     float32 `json:"price"`
	}

	var courses []resp
	for _, course := range courseBody {

		var resps resp
		resps.ID = course.ID
		resps.Vote = course.Vote
		resps.Foto = course.Foto
		resps.Judul = course.Judul
		resps.Deskripsi = course.Deskripsi
		resps.Rate = course.Rate
		resps.Price = course.Price

		courses = append(courses, resps)
	}

	h.SuccessResponse(ctx, http.StatusOK, "Success", courses, &courseParam.PaginationParam)
}

func (h *handler) getCourseData(ctx *gin.Context) {
	var courseBody entity.CourseParam
	if err := h.BindParam(ctx, &courseBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "failed to bind body", nil)
		return
	}

	var courseDB entity.Course

	err := h.db.Where("id = ?", courseBody.ID).Take(&courseDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var resp entity.CourseRespData

	resp.Foto = courseDB.Foto
	resp.Judul = courseDB.Judul
	resp.Deskripsi = courseDB.Deskripsi
	resp.Intro = courseDB.Intro
	resp.Price = courseDB.Price
	resp.Rate = courseDB.Rate

	var playlists []entity.Playlist
	if err := h.db.Where("course_id = ?", courseDB.ID).Find(&playlists).Error; err != nil {
		fmt.Println(err)
		h.ErrorResponse(ctx, http.StatusInternalServerError, "error occurred", nil)
		return
	}

	for _, playlist := range playlists {
		var videos []entity.Video
		if err := h.db.Where("playlist_id = ?", playlist.ID).Find(&videos).Error; err != nil {
			fmt.Println(err)
			h.ErrorResponse(ctx, http.StatusInternalServerError, "error occurred", nil)
			return
		}

		var respVideos []entity.RespVideo
		for _, v := range videos {
			respVideos = append(respVideos, entity.RespVideo{
				Link:       v.Link,
				Judul:      v.Judul,
				Durasi:     v.Durasi,
				PlaylistID: v.PlaylistID,
			})
		}

		resp.Playlist = append(resp.Playlist, entity.RespPlaylist{
			Nama:     playlist.Nama,
			Durasi:   playlist.Durasi,
			CourseID: playlist.CourseID,
			Video:    respVideos,
		})
	}

	h.SuccessResponse(ctx, http.StatusOK, "Success", resp, nil)
}

func (h *handler) getCourseInfo(ctx *gin.Context) {
	var courseBody entity.CourseParam
	if err := h.BindParam(ctx, &courseBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "failed to bind body", nil)
		return
	}

	var courseDB entity.Course

	err := h.db.Where("id = ?", courseBody.ID).Take(&courseDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	type RespVideo struct {
		Judul  string `json:"judul" gorm:"type:varchar(255)"`
		Durasi string `json:"durasi"`
	}

	type RespPlaylist struct {
		Nama   string
		Durasi string
		Video  []RespVideo
	}

	type RespSkill struct {
		Nama string
	}

	type CourseRespData struct {
		Judul     string         `json:"judul"`
		Deskripsi string         `json:"deskripsi"`
		Intro     string         `json:"intro"`
		Playlist  []RespPlaylist `json:"playlist"`
		Rate      float32        `json:"rate"`
		Vote      uint           `json:"vote"`
		Price     float32        `json:"price"`
		Skill     []RespSkill    `json:"skill"`
	}

	var resp CourseRespData

	resp.Vote = courseDB.Vote
	resp.Judul = courseDB.Judul
	resp.Deskripsi = courseDB.Deskripsi
	resp.Intro = courseDB.Intro
	resp.Price = courseDB.Price
	resp.Rate = courseDB.Rate

	var playlists []entity.Playlist
	if err := h.db.Where("course_id = ?", courseDB.ID).Find(&playlists).Error; err != nil {
		fmt.Println(err)
		h.ErrorResponse(ctx, http.StatusInternalServerError, "error occurred", nil)
		return
	}

	for _, playlist := range playlists {
		var videos []entity.Video
		if err := h.db.Where("playlist_id = ?", playlist.ID).Find(&videos).Error; err != nil {
			fmt.Println(err)
			h.ErrorResponse(ctx, http.StatusInternalServerError, "error occurred", nil)
			return
		}

		var respVideos []RespVideo
		for _, v := range videos {
			respVideos = append(respVideos, RespVideo{
				Judul:  v.Judul,
				Durasi: v.Durasi,
			})
		}

		resp.Playlist = append(resp.Playlist, RespPlaylist{
			Nama:   playlist.Nama,
			Durasi: playlist.Durasi,
			Video:  respVideos,
		})
	}

	var skills []entity.Skill
	if err := h.db.Joins("JOIN coursess_skill ON coursess_skill.skill_id = skills.id").Where("coursess_skill.course_id = ?", courseDB.ID).Find(&skills).Error; err != nil {
		fmt.Println(err)
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	var skillresp []RespSkill
	for _, s := range skills {
		skillresp = append(skillresp, RespSkill{Nama: s.Nama})
	}

	resp.Skill = skillresp

	h.SuccessResponse(ctx, http.StatusOK, "Success", resp, nil)
}
