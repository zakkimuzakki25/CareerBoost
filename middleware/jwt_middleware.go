package middleware

import (
	"CareerBoost/entity"
	"CareerBoost/sdk/config"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, code int64, message string, data interface{}) {
	c.JSON(int(code), entity.HTTPResponse{
		Message:    message,
		IsSuccess:  false,
		Data:       data,
		Pagination: nil,
	})
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(authorization, "Bearer") {
			ErrorResponse(c, http.StatusInternalServerError, "Unauthorized1", nil)
			return
		}
		tokenJwt := authorization
		claims := entity.UserClaims{}
		jwtKey := os.Getenv("SECRET_KEY")
		if err := config.DecodeToken(tokenJwt, &claims, jwtKey); err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Unauthorized2", nil)
			return
		}
		c.Set("user", claims)
	}
}
