package handler

import (
	"CareerBoost/sdk/config"
	"CareerBoost/src/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// function register
func (h *handler) userRegister(ctx *gin.Context) {
	var userBody entity.UserRegister

	if err := h.BindBody(ctx, &userBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "Kolom harus diisi", nil)
		return
	}

	var userDB entity.User

	userDB.FullName = userBody.FullName
	userDB.Username = userBody.Username
	userDB.Email = userBody.Email
	userDB.Password = userBody.Password

	var interest []entity.Interest
	if err := h.db.Find(&interest, userBody.InterestID).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "interest not found", nil)
		return
	}

	hashPW, _ := bcrypt.GenerateFromPassword([]byte(userDB.Password), bcrypt.DefaultCost)
	userDB.Password = string(hashPW)

	if err := h.db.Create(&userDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "Username atau Email tidak tersedia", nil)
		fmt.Println(err)
		return
	}

	if err := h.db.Model(&userDB).Association("Interest").Append(interest); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "interest not added", nil)
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "Registrasi sukses", nil, nil)
}

// function login
func (h *handler) userLogin(ctx *gin.Context) {
	var userBody entity.UserLogin

	if err := h.BindBody(ctx, &userBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "Kolom harus diisi", nil)
		return
	}

	var user entity.User

	if err := h.db.Where("email = ?", userBody.Email).First(&user).Error; err != nil {

		if err := h.db.Where("username = ?", userBody.Email).First(&user).Error; err != nil {
			h.ErrorResponse(ctx, http.StatusBadRequest, "Akun tidak ditemukan", nil)
			return
		}

	}

	//cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userBody.Password)); err != nil {
		h.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	tokenJwt, err := config.GenerateToken(user)
	if err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "create token failed", nil)
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    tokenJwt,
		HttpOnly: true,
	})

	ctx.Header("Authorization", "Bearer "+tokenJwt)
	h.SuccessResponse(ctx, http.StatusOK, "Login Berhasil", gin.H{
		"token": tokenJwt,
	}, nil)
}

func (h *handler) userUpdateProfile(ctx *gin.Context) {
	var userBody entity.UserProfilePage

	if err := h.BindBody(ctx, &userBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "Kolom harus diisi", nil)
		return
	}

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

	var userDB entity.User

	if err := h.db.Model(&userDB).Where("id = ?", userID).First(&userDB).Updates(entity.User{
		FullName:     userBody.FullName,
		Lokasi:       userBody.Lokasi,
		TanggalLahir: userBody.TanggalLahir,
		TempatLahir:  userBody.TempatLahir,
		Deskripsi:    userBody.Deskripsi,
	}).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "error sini", nil)
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "Update berhasil", nil, nil)
}

// upload foto
func (h *handler) userUploadPhotoProfile(ctx *gin.Context) {
	file, err := ctx.FormFile("profile")
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	link, err := h.supClient.Upload(file)
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

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

	var userDB entity.User

	if err := h.db.Model(&userDB).Where("id = ?", userID).First(&userDB).Updates(entity.User{
		ProfilePhoto: link,
	}).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "Update berhasil", nil, nil)
}

func (h *handler) userGetProfile(ctx *gin.Context) {
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

	var userDB entity.User
	err := h.db.Where("id = ?", userID).Take(&userDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	errr := h.db.Preload("Interest").First(&userDB, userID).Error
	if errr != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var interest []entity.Interest
	if err := h.db.Model(&userDB).Association("Interest").Find(&interest); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	var interests []entity.RespInterest
	for _, s := range interest {
		interests = append(interests, entity.RespInterest{Nama: s.Nama})
	}

	userResp := entity.UserProfilePage{
		Email:        userDB.Email,
		FullName:     userDB.FullName,
		Lokasi:       userDB.Lokasi,
		ProfilePhoto: userDB.ProfilePhoto,
		Deskripsi:    userDB.Deskripsi,
		TanggalLahir: userDB.TanggalLahir,
		TempatLahir:  userDB.TempatLahir,
		InterestID:   interests,
	}

	h.SuccessResponse(ctx, http.StatusOK, "Succes", userResp, nil)
}

func (h *handler) userGetHome(ctx *gin.Context) {
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

	var userDB entity.User
	err := h.db.Where("id = ?", userID).Take(&userDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userResp := entity.UserHome{
		FullName:     userDB.FullName,
		ProfilePhoto: userDB.ProfilePhoto,
	}

	h.SuccessResponse(ctx, http.StatusOK, "Succes", userResp, nil)
}

func (h *handler) UserAddMentor(ctx *gin.Context) {
	var reqBody entity.MentorReqByID

	if err := h.BindBody(ctx, &reqBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request update", nil)
		return
	}

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

	var userDB entity.User
	if err := h.db.Preload("Mentor").Where("id = ?", userID).First(&userDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var mentor entity.Mentor
	if err := h.db.Where("id = ?", reqBody.ID).First(&mentor).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "mentor not found", nil)
		return
	}

	if h.db.Model(&userDB).Where("id = ?", mentor.ID).Association("Mentor").Count() > 0 {
		h.ErrorResponse(ctx, http.StatusBadRequest, "Sudah berlangganan", nil)
		return
	}

	if err := h.db.Model(&userDB).Association("Mentor").Append(&mentor); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "failed to add mentor", nil)
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "Berhasil berlangganan", nil, nil)
}

func (h *handler) userGetMentors(ctx *gin.Context) {
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

	type respHistory struct {
		ID       uint                `json:"id"`
		Nama     string              `json:"nama"`
		Interest entity.RespInterest `json:"interest"`
	}

	var userDB entity.User
	err := h.db.Preload("Mentor.Interest").Where("id = ?", userID).Take(&userDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	mentors := userDB.Mentor
	if len(mentors) == 0 {
		h.SuccessResponse(ctx, http.StatusOK, "Success", nil, nil)
		return
	}

	var histories []respHistory
	for _, mentor := range mentors {
		var interest entity.RespInterest
		if len(mentor.Interest) > 0 {
			interest = entity.RespInterest{
				Nama: mentor.Interest[0].Nama,
			}
		}
		histories = append(histories, respHistory{
			ID:       mentor.ID,
			Nama:     mentor.FullName,
			Interest: interest,
		})
	}

	h.SuccessResponse(ctx, http.StatusOK, "ini mentor", histories, nil)
}

func (h *handler) UserAddMagang(ctx *gin.Context) {
	var reqBody entity.MagangReqByID

	if err := h.BindBody(ctx, &reqBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request update", nil)
		return
	}

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

	var userDB entity.User
	if err := h.db.Preload("Magang").Where("id = ?", userID).First(&userDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var magang entity.Magang
	if err := h.db.Where("id = ?", reqBody.ID).First(&magang).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "magang not found", nil)
		return
	}

	if h.db.Model(&userDB).Where("id = ?", magang.ID).Association("Magang").Count() > 0 {
		h.ErrorResponse(ctx, http.StatusBadRequest, "Sudah apply", nil)
		return
	}

	if err := h.db.Model(&userDB).Association("Magang").Append(&magang); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "failed to add magang", nil)
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "Berhasil apply magang", nil, nil)
}

func (h *handler) userGetMagangs(ctx *gin.Context) {
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

	type respHistory struct {
		ID         uint                `json:"id"`
		Logo       string              `json:"logo"`
		Perusahaan string              `json:"perusahaan"`
		Interest   entity.RespInterest `json:"interest"`
	}

	var userDB entity.User
	err := h.db.Preload("Magang.Interest").Where("id = ?", userID).Take(&userDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if len(userDB.Magang) == 0 {
		h.SuccessResponse(ctx, http.StatusOK, "Success", nil, nil)
		return
	}

	magangs := userDB.Magang
	var histories []respHistory
	for _, magang := range magangs {
		var interest entity.RespInterest
		if len(magang.Interest) > 0 {
			interest = entity.RespInterest{
				Nama: magang.Interest[0].Nama,
			}
		}
		histories = append(histories, respHistory{
			ID:         magang.ID,
			Logo:       magang.Logo,
			Perusahaan: magang.Perusahaan,
			Interest:   interest,
		})
	}

	h.SuccessResponse(ctx, http.StatusOK, "ini magang", histories, nil)
}

func (h *handler) UserAddCourse(ctx *gin.Context) {
	var reqBody entity.CourseReqByID

	if err := h.BindBody(ctx, &reqBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request update", nil)
		return
	}

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

	var userDB entity.User
	if err := h.db.Preload("Magang").Where("id = ?", userID).First(&userDB).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var course entity.Course
	if err := h.db.Where("id = ?", reqBody.ID).First(&course).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "course not found", nil)
		return
	}

	if h.db.Model(&userDB).Where("id = ?", course.ID).Association("Course").Count() > 0 {
		h.ErrorResponse(ctx, http.StatusBadRequest, "Sudah berlangganan", nil)
		return
	}

	if err := h.db.Model(&userDB).Association("Course").Append(&course); err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "failed to add course", nil)
		return
	}

	h.SuccessResponse(ctx, http.StatusOK, "Berhasil berlangganan", nil, nil)
}

// func (h *handler) userGetCourses(ctx *gin.Context) {
// 	user, exist := ctx.Get("user")
// 	if !exist {
// 		h.ErrorResponse(ctx, http.StatusBadRequest, "Unauthorized", nil)
// 		return
// 	}

// 	claims, ok := user.(entity.UserClaims)
// 	if !ok {
// 		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid token", nil)
// 		return
// 	}

// 	userID := claims.ID

// 	type respHistory struct {
// 		ID    uint   `json:"id"`
// 		Judul string `json:"judul"`
// 		Intro string `json:"intro"`
// 	}

// 	var userDB entity.User
// 	err := h.db.Preload("Magang.Interest").Where("id = ?", userID).Take(&userDB).Error
// 	if err != nil {
// 		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
// 		return
// 	}

// 	if len(userDB.Magang) == 0 {
// 		h.SuccessResponse(ctx, http.StatusOK, "Success", nil, nil)
// 		return
// 	}

// 	magangs := userDB.Magang
// 	var histories []respHistory
// 	for _, magang := range magangs {
// 		var interest entity.RespInterest
// 		if len(magang.Interest) > 0 {
// 			interest = entity.RespInterest{
// 				Nama: magang.Interest[0].Nama,
// 			}
// 		}
// 		histories = append(histories, respHistory{
// 			ID: magang.ID,
// 		})
// 	}

// 	h.SuccessResponse(ctx, http.StatusOK, "ini magang", histories, nil)
// }
