package handlers

import (
	"context"
	"go-gin-postgre/middleware"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

var googleOauthConfig *oauth2.Config

func InitGoogleOauthConfig() {
	scopesStr := os.Getenv("GOOGLE_OAUTH_SCOPES")
	scopes := strings.Split(scopesStr, ",")

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       scopes,
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

	log.Println(token)

	person, err := getUserInfoFromGoogle(token.AccessToken)

	if err != nil {
		log.Println(err)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Google 정보를 가져오는 중 오류가 발생했습니다."})
		// return
	}

	data := make(map[string]interface{})

	if person != nil && len(person.Names) > 0 && person.Names[0].GivenName != "" { // Name이 있으면
		data["name"] = person.Names[0].GivenName
	}

	if person != nil && len(person.EmailAddresses) > 0 && person.EmailAddresses[0].Value != "" { // EmailAddresses의 첫 번째 이메일이 있으면
		data["email"] = person.EmailAddresses[0].Value
	}

	// JWT 토큰 생성
	// jwtToken, err := middleware.CreateToken(data, time.Hour*24)  // 24시간 만료 시간 설정
	jwtToken, err := middleware.CreateToken(data, time.Minute*30) // 30분 만료 시간 설정

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwtToken})
}

func getUserInfoFromGoogle(accessToken string) (*people.Person, error) {
	ctx := context.Background()

	ts := googleOauthConfig.TokenSource(ctx, &oauth2.Token{AccessToken: accessToken})
	service, err := people.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, err
	}

	// people/me는 현재 인증된 사용자를 의미합니다.
	person, err := service.People.Get("people/me").PersonFields("names,emailAddresses").Do()
	if err != nil {
		return nil, err
	}

	return person, nil
}
