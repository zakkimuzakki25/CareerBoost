package handler

import (
	"CareerBoost/middleware"
	"CareerBoost/sdk/config"
	"fmt"
	"net/http"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	http      *gin.Engine
	config    config.Interface
	db        *gorm.DB
	supClient supabasestorageuploader.SupabaseClientService
}

func Init(config config.Interface, db *gorm.DB, supClient supabasestorageuploader.SupabaseClientService) *handler {
	rest := handler{
		http:      gin.New(),
		config:    config,
		db:        db,
		supClient: supClient,
	}

	rest.registerRoutes()

	return &rest
}

func (h *handler) registerRoutes() {

	// h.http.Use(middleware.CORS())
	api := h.http.Group("api")

	api.Use(middleware.JwtMiddleware())
	api.Use(h.userGetHome)

	admin := h.http.Group("admin")
	admin.Use(middleware.JwtMiddlewareAdmin())

	// Post

	h.http.POST("/admin/login", h.adminLogin)
	admin.GET("/logout", h.adminLogout)
	admin.POST("/mentor/post", h.addNewMentor)
	admin.GET("/", h.ping)

	api.GET("/", h.ping)
	h.http.POST("/user/register", h.userRegister)
	h.http.POST("/user/login", h.userLogin)
	h.http.GET("/user/logout", h.userLogout)
	api.POST("/profile", h.userUpdateProfile)
	api.POST("/profile/photo/update", h.userUploadPhotoProfile)
	api.GET("/profile", h.userGetProfile)
	api.GET("/mentorsinfo", h.getMentorRekomendation, h.getAllMentor)
	api.POST("/mentorsinfo/filter", h.getMentorFilter)
	api.POST("/mentorsinfo/search", h.getMentorSearch)
	api.POST("/mentorinfo/data", h.getMentorData)
	api.POST("/mentorinfo/pengalaman", h.getMentorExp)

}

func (h *handler) ping(ctx *gin.Context) {
	h.SuccessResponse(ctx, http.StatusOK, "ping", nil, nil)
}

func (h *handler) Run() {
	h.http.Run(fmt.Sprintf(":%s", h.config.Get("PORT")))
}
