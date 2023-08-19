package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v3"
)

type IceServer struct {
	URLs       []string `json:"urls"`
	Username   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}

func GetIceServers(c *gin.Context) {
	servers := []IceServer{
		{
			URLs: []string{"stun:localhost:3478"},
		},
		{
			URLs:       []string{"turn:localhost:3478"},
			Username:   "your_turn_username",
			Credential: "test_key",
		},
	}

	c.JSON(http.StatusOK, servers)
}

func Offer(c *gin.Context) {
	// Pion WebRTC API 사용하여 PeerConnection 생성
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Fatalf("PeerConnection 생성 중 오류: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Offer 생성
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		log.Fatalf("Offer 생성 중 오류: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Offer 설정
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		log.Fatalf("Offer 설정 중 오류: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Offer 반환
	c.JSON(http.StatusOK, gin.H{"offer": offer})
}

var storedAnswer map[string]interface{}

func Answer(c *gin.Context) {
	var data map[string]interface{}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storedAnswer = data

	c.JSON(http.StatusOK, gin.H{"message": "Answer received"})
}

func GetAnswer(c *gin.Context) {
	if storedAnswer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}

	c.JSON(http.StatusOK, storedAnswer)
}
