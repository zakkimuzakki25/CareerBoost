package handler

import (
	"CareerBoost/middleware"
	"CareerBoost/sdk/config"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	http   *gin.Engine
	config config.Interface
	db     *gorm.DB
}

func Init(config config.Interface, db *gorm.DB) *handler {
	rest := handler{
		http:   gin.Default(),
		config: config,
		db:     db,
	}

	rest.registerRoutes()

	return &rest
}

func (h *handler) registerRoutes() {

	api := h.http.Group("api")

	api.Use(middleware.JwtMiddleware())

	// Post
	api.GET("/", h.ping)
	h.http.POST("/user/register", h.userRegister)
	h.http.POST("/user/login", h.userLogin)
	h.http.GET("/user/logout", h.userLogout)
	api.POST("/profile/update", h.userUpdateProfile)

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
