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
	api.GET("/logout", h.adminLogout)
	admin.POST("/mentor/post", h.addNewMentor)
	admin.POST("/magang/post", h.addNewMagang)
	admin.GET("/", h.ping)

	api.GET("/", h.ping)
	h.http.POST("/user/register", h.userRegister)
	h.http.POST("/user/login", h.userLogin)
	api.PUT("/profile", h.userUpdateProfile)
	api.PUT("/profile/photo/update", h.userUploadPhotoProfile)
	api.GET("/profile", h.userGetProfile)
	api.GET("/profile/history", h.userGetMagang, h.userGetMentor)
	api.GET("/profile/langganan", h.userGetMagang, h.userGetMentor)

	api.GET("/mentorsinfo", h.getMentorRekomendation, h.getAllMentor)
	api.POST("/mentorsinfo", h.getMentorRekomendation, h.getMentorFilter)
	api.POST("/mentorinfo/data", h.getMentorData)
	api.POST("/mentorinfo/pengalaman", h.getMentorExp)
	api.POST("/mentorinfo/checkout", h.UserAddMentor)

	api.GET("/magangsinfo", h.getMagangRekomendation, h.getAllMagang)
	api.POST("/magangsinfo", h.getMagangRekomendation, h.getMagangFilter)
	api.POST("/maganginfo/data", h.getMagangRekomendation, h.getMagangData)
	api.POST("/maganginfo/checkout", h.UserAddMagang)

}

func (h *handler) ping(ctx *gin.Context) {
	h.SuccessResponse(ctx, http.StatusOK, "ping", nil, nil)
}

func (h *handler) Run() {
	h.http.Run(fmt.Sprintf(":%s", h.config.Get("PORT")))
}
