package handler

import (
	"CareerBoost/sdk/config"
	"CareerBoost/sdk/middleware"
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
	// api.Use(h.userGetHome)

	admin := h.http.Group("admin")
	admin.Use(middleware.JwtMiddlewareAdmin())

	// Post

	h.http.POST("/admin/login", h.adminLogin)
	api.GET("/logout", h.adminLogout)
	admin.POST("/mentor/post", h.addNewMentor)
	admin.POST("/magang/post", h.addNewMagang)
	admin.POST("/course/post", h.addNewCourse)
	admin.GET("/", h.ping)

	h.http.GET("/", h.userGetHome)
	h.http.POST("/user/register", h.userRegister)
	h.http.POST("/user/login", h.userLogin)
	api.PUT("/profile", h.userUpdateProfile)
	api.PUT("/profile/photo/update", h.userUploadPhotoProfile)
	api.GET("/profile", h.userGetProfile)
	api.GET("/profile/history", h.userGetRiwayat)
	// api.GET("/profile/history/course", h.userGetCourses)
	// api.GET("/profile/history/mentor", h.userGetMentors)
	// api.GET("/profile/langganan", h.userGetMentors)

	api.GET("/mentorinfo/all", h.getAllMentor)
	api.GET("/mentorinfo/rekomendasi", h.getMagangRecomendation)
	api.GET("/mentorinfo", h.getMentorFilter)
	api.GET("/mentorinfo/data", h.getMentorData)
	api.GET("/mentorinfo/pengalaman", h.getMentorExp)
	api.POST("/mentorinfo/checkout", h.UserAddMentor)

	// api.GET("/maganginfo", h.getAllMagang)
	api.GET("/maganginfo", h.getMagangFilter)
	api.GET("/maganginfo/data", h.getMagangData)
	api.POST("/maganginfo/checkout", h.UserAddMagang)

	api.GET("/courseinfo", h.getAllCourse)
	// api.GET("/courseinfo/recomendation", h.getAllCourseRecomendation)
	api.GET("/courseinfo/all", h.getAllCourse)
	api.GET("/courseinfo/data", h.getCourseData)
	api.POST("/courseinfo/checkout", h.UserAddCourse)

}

func (h *handler) ping(ctx *gin.Context) {
	h.SuccessResponse(ctx, http.StatusOK, "ping", nil, nil)
}

func (h *handler) Run() {
	h.http.Run(fmt.Sprintf(":%s", h.config.Get("PORT")))
}
