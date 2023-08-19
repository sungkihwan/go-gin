package middleware

import (
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey string

func InitJWTSecretKey() {
	secretKey = os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatalf("JWT_SECRET_KEY 환경 변수가 설정되지 않았습니다.")
	}
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 성공적으로 인증된 경우 다음 핸들러로 넘어감
		c.Next()
	}
}
