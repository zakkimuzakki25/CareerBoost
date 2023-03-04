package middleware

import (
	"CareerBoost/entity"
	"CareerBoost/sdk/config"
	"net/http"
	"os"

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
		tokenJwt, err := c.Cookie("token")
		if err != nil {
			ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

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
