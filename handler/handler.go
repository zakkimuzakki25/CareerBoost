package handler

import (
	"CareerBoost/middleware"
	"CareerBoost/sdk/config"
	"fmt"
	"net/http"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-contrib/cors"
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
		http:      gin.Default(),
		config:    config,
		db:        db,
		supClient: supClient,
	}

	rest.registerRoutes()

	return &rest
}

func (h *handler) registerRoutes() {

	h.http.Use(cors.Default())

	api := h.http.Group("api")

	api.Use(middleware.JwtMiddleware())

	// Post
	api.GET("/", h.ping)
	h.http.POST("/user/register", h.userRegister)
	h.http.POST("/user/login", h.userLogin)
	h.http.GET("/user/logout", h.userLogout)
	api.POST("/profile", h.userUpdateProfile)
	api.POST("/profile/photo/upload", h.userUploadPhotoProfile)
	api.GET("/profile", h.userGetProfile)
	api.GET("/", h.userGetHome)

	// v1.GET("/post/:post_id", h.getPost)
	// v1.PUT("/post/:post_id", h.updatePost) // 1 -> aku mau update post yang id nya 1
	// v1.DELETE("/post/:post_id", h.deletePost)
}

func (h *handler) ping(ctx *gin.Context) {
	h.SuccessResponse(ctx, http.StatusOK, "pong", nil, nil)
}

func (h *handler) Run() {
	h.http.Run(fmt.Sprintf(":%s", h.config.Get("PORT")))
}
