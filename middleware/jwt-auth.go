package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

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

type Claims struct {
	Data map[string]interface{} `json:"data"`
	jwt.StandardClaims
}

func CreateToken(data map[string]interface{}, duration time.Duration) (string, error) {
	claims := &Claims{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// "Bearer" 접두어 처리
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(401, gin.H{"error": "Authorization header format is invalid"})
			c.Abort()
			return
		}
		tokenString := authHeader[7:]

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			// 에러 메시지 로깅 및 반환
			log.Printf("토큰 검증 에러: %v", err)
			c.JSON(401, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 만료 검사
		if claims.ExpiresAt < time.Now().Unix() {
			c.JSON(401, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		// 토큰 클레임 gin.Context에 저장
		c.Set("claims", claims)

		c.Next()
	}
}
