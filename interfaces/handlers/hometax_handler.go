package handlers

import (
	"fmt"
	"go-gin-postgre/domain/usecases"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type HometaxHandler interface {
	HandleRequest(c *gin.Context)
}

type hometaxHandler struct {
	usecase usecases.HomtaxMacroUseCase
}

func NewHometaxHandler(u usecases.HomtaxMacroUseCase) HometaxHandler {
	return &hometaxHandler{
		usecase: u,
	}
}

func (h *hometaxHandler) fetchDataFromAPI() (string, error) {
	props := map[string]interface{}{
		"actionId":     "someActionId",
		"screenId":     "someScreenId",
		"realScreenId": "someRealScreenId",
	}

	_, err := h.usecase.SendXmlRequest(props)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	resp, err := http.Get("https://example.com/api/data")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 응답 오류: %v", resp.Status)
	}

	builder := &strings.Builder{}
	_, err = io.Copy(builder, resp.Body)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func (h *hometaxHandler) HandleRequest(c *gin.Context) {
	var data string
	var err error

	var wg sync.WaitGroup

	// 외부 API 호출 시작
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, err = h.fetchDataFromAPI()
	}()

	// goroutine이 완료될 때까지 대기
	wg.Wait()

	// 결과 처리
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("API 호출 중 에러 발생: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
