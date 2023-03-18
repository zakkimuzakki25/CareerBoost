package handler

import (
	"CareerBoost/sdk/config"
	"CareerBoost/sdk/middleware"
	"fmt"

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

	api := h.http.Group("api")

	api.Use(middleware.JwtMiddleware())

	admin := h.http.Group("admin")
	admin.Use(middleware.JwtMiddlewareAdmin())

	h.http.POST("/admin/login", h.adminLogin)
	admin.GET("/logout", h.adminLogout)
	admin.POST("/mentor/post", h.addNewMentor)
	admin.POST("/magang/post", h.addNewMagang)
	admin.POST("/course/post", h.addNewCourse)

	h.http.GET("/", h.userGetHome)
	h.http.POST("/user/register", h.userRegister)
	h.http.POST("/user/login", h.userLogin)
	api.PUT("/profile", h.userUpdateProfile)
	api.PUT("/profile/photo/update", h.userUploadPhotoProfile)
	api.GET("/profile", h.userGetProfile)
	api.GET("/profile/history", h.userGetRiwayat)
	api.GET("/profile/langganan", h.userGetLangganan)

	api.GET("/mentorinfo/rekomendasi", h.getAllMentor)
	api.GET("/mentorinfo", h.getMentorFilter)
	api.GET("/mentorinfo/data/:mentor_id", h.getMentorData)
	api.GET("/mentorinfo/pengalaman/:mentor_id", h.getMentorExp)
	api.POST("/mentorinfo/checkout/:mentor_id", h.UserAddMentor)

	api.GET("/maganginfo/hiring", h.getMagangFilter)
	api.GET("/maganginfo/rekomendasi", h.getMagangRecomendation)
	api.GET("/maganginfo/:magang_id", h.getMagangData)
	api.POST("/maganginfo/checkout/:magang_id", h.UserAddMagang)

	api.GET("/courseinfo/rekomendasi", h.getAllCourseHome)
	api.GET("/courseinfo", h.getCourseRekomendasi)
	api.GET("/courseinfo/all", h.getAllCourse)
	api.GET("/courseinfo/:course_id", h.getCourseInfo)
	api.GET("/coursedata/:course_id", h.getCourseData)
	api.POST("/courseinfo/checkout/:course_id", h.UserAddCourse)

}

func (h *handler) Run() {
	h.http.Run(fmt.Sprintf(":%s", h.config.Get("PORT")))
}
