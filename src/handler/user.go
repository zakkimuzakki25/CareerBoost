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
		if len(userBody.FullName) > 50 {
			h.ErrorResponse(ctx, http.StatusBadRequest, "Username minimal 5 karakter", nil)
			return
		}
		if len(userBody.Username) < 5 {
			h.ErrorResponse(ctx, http.StatusBadRequest, "Username minimal 5 karakter", nil)
			return
		}
		if len(userBody.Username) > 20 {
			h.ErrorResponse(ctx, http.StatusBadRequest, "Username maximal 20 karakter", nil)
			return
		}
		if len(userBody.Email) < 1 {
			h.ErrorResponse(ctx, http.StatusBadRequest, "Masukkan email dengan benar", nil)
			return
		}
		if len(userBody.Password) < 8 {
			h.ErrorResponse(ctx, http.StatusBadRequest, "Password minimal 8 karakter", nil)
			return
		}
		if len(userBody.InterestID) < 1 {
			h.ErrorResponse(ctx, http.StatusBadRequest, "Minimal pilih 1 ketertarikan", nil)
			return
		}
		h.ErrorResponse(ctx, http.StatusBadRequest, "Masukkan data dengan benar", nil)
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
			h.ErrorResponse(ctx, http.StatusBadRequest, "akun atau password salah", nil)
			return
		}

	}

	//cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userBody.Password)); err != nil {
		h.ErrorResponse(ctx, http.StatusUnauthorized, "akun atau password salah", nil)
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
	var userBody entity.UserUpdateProfile

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

	var interests []entity.Interest
	if err := h.db.Where("id IN ?", userBody.InterestID).Find(&interests).Error; err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, "error sini", nil)
		return
	}

	if err := h.db.Model(&userDB).Association("Interest").Replace(interests); err != nil {
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

	var interests []entity.RespInterestWithID
	for _, s := range interest {
		interests = append(interests, entity.RespInterestWithID{
			Nama: s.Nama,
			ID:   s.ID,
		})
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
	type mentorParam struct {
		ID uint `uri:"mentor_id" gorm:"column:id"`
	}

	var reqBody mentorParam

	if err := h.BindParam(ctx, &reqBody); err != nil {
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

func (h *handler) UserAddMagang(ctx *gin.Context) {
	var reqBody entity.MagangParam

	if err := h.BindParam(ctx, &reqBody); err != nil {
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

func (h *handler) UserAddCourse(ctx *gin.Context) {
	var reqBody entity.CourseParam

	if err := h.BindParam(ctx, &reqBody); err != nil {
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

func (h *handler) userGetRiwayat(ctx *gin.Context) {
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

	type respHistoryMagang struct {
		Logo       string              `json:"logo"`
		Perusahaan string              `json:"perusahaan"`
		Interest   entity.RespInterest `json:"interest"`
	}

	type respHistoryCourse struct {
		Judul    string              `json:"judul"`
		Interest entity.RespInterest `json:"interest"`
	}

	type respHistoryMentor struct {
		Nama     string              `json:"nama"`
		Interest entity.RespInterest `json:"interest"`
	}

	type respHistory struct {
		Magang []respHistoryMagang `json:"magang"`
		Course []respHistoryCourse `json:"course"`
		Mentor []respHistoryMentor `json:"mentor"`
	}

	var userDB entity.User
	err := h.db.Preload("Mentor.Interest").Where("id = ?", userID).Take(&userDB).Error
	if err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	mentors := userDB.Mentor
	var historiesMentor []respHistoryMentor
	for _, mentor := range mentors {
		var interest entity.RespInterest
		if len(mentor.Interest) > 0 {
			interest = entity.RespInterest{
				Nama: mentor.Interest[0].Nama,
			}
		}
		historiesMentor = append(historiesMentor, respHistoryMentor{
			Nama:     mentor.FullName,
			Interest: interest,
		})
	}

	var userDBMagang entity.User
	errrr := h.db.Preload("Magang.Interest").Where("id = ?", userID).Take(&userDBMagang).Error
	if errrr != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	magangs := userDBMagang.Magang
	var historiesMagang []respHistoryMagang
	for _, magang := range magangs {
		var interest entity.RespInterest
		if len(magang.Interest) > 0 {
			interest = entity.RespInterest{
				Nama: magang.Interest[0].Nama,
			}
		}
		historiesMagang = append(historiesMagang, respHistoryMagang{
			Logo:       magang.Logo,
			Perusahaan: magang.Perusahaan,
			Interest:   interest,
		})
	}

	var courseDB []entity.Course
	errr := h.db.
		Joins("JOIN user_courses ON user_courses.course_id = courses.id").
		Where("user_courses.user_id = ?", userID).
		Find(&courseDB).Error
	if errr != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var historiesCourse []respHistoryCourse
	for _, course := range courseDB {
		var interest entity.RespInterest

		if course.InterestID != 0 {
			interests := &entity.Interest{}
			err := h.db.Model(interests).Where("id = ?", course.InterestID).Take(&interests).Error
			if err != nil {
				h.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
				return
			}
			interest = entity.RespInterest{
				Nama: interests.Nama,
			}
		}

		historiesCourse = append(historiesCourse, respHistoryCourse{
			Judul:    course.Judul,
			Interest: interest,
		})
	}

	var resp respHistory

	resp.Mentor = historiesMentor
	resp.Course = historiesCourse
	resp.Magang = historiesMagang

	h.SuccessResponse(ctx, http.StatusOK, "Success", resp, nil)
}

func (h *handler) userGetLangganan(ctx *gin.Context) {
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

	type respLanggananCourse struct {
		ID    uint   `json:"id"`
		Judul string `json:"judul"`
		Foto  string `json:"foto"`
	}

	type respLanggananMentor struct {
		FullName     string `json:"full_name"`
		PhotoProfile string `json:"photo_profile"`
		WA           string `json:"wa"`
		IG           string `json:"ig"`
		Email        string `json:"email"`
	}

	type respLangganan struct {
		Course respLanggananCourse `json:"course"`
		Mentor respLanggananMentor `json:"mentor"`
	}

	var resp respLangganan

	var langgananMentor respLanggananMentor
	var mentorDB entity.Mentor
	err := h.db.
		Joins("JOIN user_mentors ON user_mentors.mentor_id = mentors.id").
		Where("user_mentors.user_id = ?", userID).
		Order("mentors.created_at DESC").
		Limit(1).
		First(&mentorDB).
		Error

	if err != nil {
		langgananMentor.FullName = ""
		langgananMentor.Email = ""
		langgananMentor.PhotoProfile = ""
		langgananMentor.IG = ""
		langgananMentor.WA = ""
	} else {
		langgananMentor.FullName = mentorDB.FullName
		langgananMentor.Email = mentorDB.Email
		langgananMentor.PhotoProfile = mentorDB.ProfilePhoto
		langgananMentor.IG = mentorDB.IG
		langgananMentor.WA = mentorDB.WA
	}

	var langgananCourse respLanggananCourse
	var courseDB entity.Course
	err2 := h.db.
		Joins("JOIN user_courses ON user_courses.course_id = courses.id").
		Where("user_courses.user_id = ?", userID).
		Order("courses.created_at DESC").
		Limit(1).
		First(&courseDB).
		Error

	if err2 != nil {
		langgananCourse.ID = 0
		langgananCourse.Foto = ""
		langgananCourse.Judul = ""
	} else {
		langgananCourse.ID = courseDB.ID
		langgananCourse.Foto = courseDB.Foto
		langgananCourse.Judul = courseDB.Judul
	}

	resp.Mentor = langgananMentor
	resp.Course = langgananCourse

	h.SuccessResponse(ctx, http.StatusOK, "Success", resp, nil)
}
