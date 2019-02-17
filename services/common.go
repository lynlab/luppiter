package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrInternal     = errors.New("unknown error, please retry later")
	ErrUnauthorized = errors.New("Unauthorized")
)

type userInfo struct {
	UUID     string
	Username string
	Email    string
}

var httpClient = &http.Client{}

func getUserInfo(accessToken string) (*userInfo, error) {
	req, _ := http.NewRequest("GET", "http://auth.lynlab.co.kr/apis/v1/me", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, _ := httpClient.Do(req)
	if res.StatusCode != 200 {
		return nil, ErrUnauthorized
	}

	var user userInfo
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	json.Unmarshal(buf.Bytes(), &user)
	return &user, nil
}
