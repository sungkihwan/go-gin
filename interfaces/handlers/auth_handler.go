package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

func InitGoogleOauthConfig() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleLoginHandler(c *gin.Context) {
	log.Printf("googleRedirectURL: %v", googleOauthConfig.RedirectURL)

	url := googleOauthConfig.AuthCodeURL("randomState")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallbackHandler(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	token, err := googleOauthConfig.Exchange(c, code)
	if err != nil {
		log.Printf("Error during token exchange: %v", err)
		c.JSON(http.StatusBadRequest, "Invalid google auth code.")
		return
	}
	// token을 사용하여 추가 정보 요청 및 JWT 토큰 생성 등 추가 작업 진행
	c.JSON(http.StatusOK, gin.H{"token": token.AccessToken})
}
