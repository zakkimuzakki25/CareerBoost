package handler

import (
	"CareerBoost/entity"
	"CareerBoost/sdk/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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

	h.SuccessResponse(ctx, http.StatusOK, "Successfully registered", nil, nil)
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
			h.ErrorResponse(ctx, http.StatusBadRequest, "Account not found", nil)
			return
		}

	}

	//cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userBody.Password)); err != nil {
		h.ErrorResponse(ctx, http.StatusUnauthorized, "Wrong password", nil)
		return
	}

	expTime := time.Now().Add(time.Minute * 60)
	claim := &config.JWTClaim{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-gin",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	//deklarasi algonya untuk sign in
	tokenAlg := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	token, err := tokenAlg.SignedString(config.JWT_KEY)
	if err != nil {
		h.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	//set token ke cookie agar lebih aman
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	h.SuccessResponse(ctx, http.StatusOK, "Login Succes", nil, nil)
}

// function logout
func (h *handler) userLogout(ctx *gin.Context) {
	//hapus token yang ada di cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	h.SuccessResponse(ctx, http.StatusOK, "Logout Succes", nil, nil)
}
