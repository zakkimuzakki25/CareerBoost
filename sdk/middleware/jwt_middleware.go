package middleware

import (
	"CareerBoost/sdk/config"
	"CareerBoost/src/entity"
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
		if !strings.HasPrefix(authorization, "Bearer ") {
			ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		tokenJwt := authorization[7:]
		claims := entity.UserClaims{}
		jwtKey := os.Getenv("SECRET_KEY")

		if err := config.DecodeToken(tokenJwt, &claims, jwtKey); err != nil {
			ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		c.Set("user", claims)
	}
}

func JwtMiddlewareAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenJwt, err := c.Cookie("tokenAdmin")
		if err != nil {
			ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		claims := entity.AdminClaims{}
		jwtKey := os.Getenv("SECRET_KEY")

		if err := config.DecodeToken(tokenJwt, &claims, jwtKey); err != nil {
			ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		c.Set("admin", claims)
	}
}
