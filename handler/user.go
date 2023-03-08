package handler

import (
	"CareerBoost/entity"
	"CareerBoost/sdk/config"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// function register
func (h *handler) userRegister(ctx *gin.Context) {
	var userBody entity.UserRegister

	if err := h.BindBody(ctx, &userBody); err != nil {
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request register", nil)
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
		h.ErrorResponse(ctx, http.StatusInternalServerError, "Email atau Username tidak tersedia", nil)
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
		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request register", nil)
		return
	}

	var user entity.User

	if err := h.db.Where("email = ?", userBody.Email).First(&user).Error; err != nil {

		if err := h.db.Where("username = ?", userBody.Email).First(&user).Error; err != nil {
			h.ErrorResponse(ctx, http.StatusBadRequest, "Email atau Password salah", nil)
			return
		}

	}

	//cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userBody.Password)); err != nil {
		h.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	tokenJwt, err := config.GenerateTokenUser(user)
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

	h.SuccessResponse(ctx, http.StatusOK, "Login Berhasil", nil, nil)
}

// function logout
func (h *handler) userLogout(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	h.SuccessResponse(ctx, http.StatusOK, "Logout Berhasil", nil, nil)
}

func (h *handler) userUpdateProfile(ctx *gin.Context) {
	var userBody entity.UserProfilePage

	if err := h.BindBody(ctx, &userBody); err != nil {
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

	h.SuccessResponse(ctx, http.StatusOK, "Succesfully Update", nil, nil)
}

// upload foto
func (h *handler) userUploadPhotoProfile(ctx *gin.Context) {
	file, err := ctx.FormFile("profile")
	if err != nil {
		h.ErrorResponse(ctx, 400, err.Error(), nil)
		return
	}

	link, err := h.supClient.Upload(file)
	if err != nil {
		h.ErrorResponse(ctx, 400, err.Error(), nil)
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

	h.SuccessResponse(ctx, http.StatusOK, "Succesfully Upload", link, nil)
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

	userResp := entity.UserProfilePage{
		Email:        userDB.Email,
		FullName:     userDB.FullName,
		Lokasi:       userDB.Lokasi,
		ProfilePhoto: userDB.ProfilePhoto,
		Deskripsi:    userDB.Deskripsi,
		TanggalLahir: userDB.TanggalLahir,
		TempatLahir:  userDB.TempatLahir,
		InterestID:   userDB.Interest,
	}

	fmt.Println(userDB.TanggalLahir)

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
