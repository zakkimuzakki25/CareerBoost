package handler

// import (
// 	"CareerBoost/entity"
// 	"CareerBoost/sdk/config"
// 	"net/http"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	"golang.org/x/crypto/bcrypt"
// )

// // function login
// func (h *handler) adminLogin(ctx *gin.Context) {
// 	var adminBody entity.AdminLogin

// 	if err := h.BindBody(ctx, &adminBody); err != nil {
// 		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request login", nil)
// 		return
// 	}

// 	var admin entity.AdminLogin

// 	adminUname := os.Getenv("ADMIN_USERNAME")
// 	if adminBody.Username != adminUname {
// 		h.ErrorResponse(ctx, http.StatusBadRequest, "invalid request login an admin", nil)
// 		return
// 	}

// 	hashPW, _ := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)
// 	adminPW := string(hashPW)

// 	//cek password
// 	if err := bcrypt.CompareHashAndPassword([]byte(adminPW), []byte(adminBody.Password)); err != nil {
// 		h.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), nil)
// 		return
// 	}

// 	tokenJwt, err := config.GenerateTokenAdmin(admin)
// 	if err != nil {
// 		h.ErrorResponse(ctx, http.StatusInternalServerError, "create token failed", nil)
// 		return
// 	}

// 	http.SetCookie(ctx.Writer, &http.Cookie{
// 		Name:     "token",
// 		Path:     "/admin",
// 		Value:    tokenJwt,
// 		HttpOnly: true,
// 	})

// 	h.SuccessResponse(ctx, http.StatusOK, "Login Berhasil", nil, nil)
// }

// // function logout
// func (h *handler) adminLogout(ctx *gin.Context) {
// 	http.SetCookie(ctx.Writer, &http.Cookie{
// 		Name:     "token",
// 		Path:     "/",
// 		Value:    "",
// 		HttpOnly: true,
// 		MaxAge:   -1,
// 	})

// 	h.SuccessResponse(ctx, http.StatusOK, "Logout Berhasil", nil, nil)
// }
