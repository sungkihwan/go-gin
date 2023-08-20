package usecases

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"go-gin-postgre/domain"
	"io"
	"net/http"
)

// 인터페이스 정의
type HomtaxMacroUseCase interface {
	SendXmlRequest(props map[string]interface{}) (*http.Response, error)
	LoginAtHometax(props map[string]interface{}) (domain.XMLMap, error)
}

// 구현체 정의
type HomtaxMacroUseCaseImpl struct {
	httpClient *http.Client
}

// 생성자 함수
func NewHomtaxMacroService(client *http.Client) HomtaxMacroUseCase {
	return &HomtaxMacroUseCaseImpl{
		httpClient: client,
	}
}

// SendXmlRequest 메서드 구현
func (ms *HomtaxMacroUseCaseImpl) SendXmlRequest(props map[string]interface{}) (*http.Response, error) {
	prefix := "www"
	if v, ok := props["prefix"]; ok {
		prefix = v.(string)
	}
	actionId := props["actionId"].(string)
	screenId := props["screenId"].(string)
	popupYn := false
	if v, ok := props["popupYn"]; ok {
		popupYn = v.(bool)
	}
	realScreenId := props["realScreenId"].(string)

	query := fmt.Sprintf("actionId=%s&screenId=%s&popupYn=%v&realScreenId=%s", actionId, screenId, popupYn, realScreenId)

	// XML 작성
	xmlData, xmlParseErr := xml.Marshal(props["map"])
	if xmlParseErr != nil {
		return nil, xmlParseErr
	}

	url := fmt.Sprintf("https://%s.hometax.go.kr/wqAction.do?%s", prefix, query)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(xmlData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/xml; charset=UTF-8")

	res, err := ms.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(data))
	return res, nil
}

func (ms *HomtaxMacroUseCaseImpl) LoginAtHometax(props map[string]interface{}) (domain.XMLMap, error) {
	response, err := ms.SendXmlRequest(props)
	if err != nil {
		return domain.XMLMap{}, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.XMLMap{}, err
	}

	var result domain.XMLRoot
	err = extractDataFromXML(data, &result)
	if err != nil {
		return domain.XMLMap{}, err
	}

	// result 값 안에서 원하는 데이터를 찾습니다.
	for _, m := range result.Maps {
		if m.ID == "원하는ID" {
			return m, nil
		}
	}

	return domain.XMLMap{}, fmt.Errorf("ID %s에 대한 데이터를 찾을 수 없습니다", "원하는ID")
}

func extractDataFromXML(data []byte, result interface{}) error {
	return xml.Unmarshal(data, result)
}
