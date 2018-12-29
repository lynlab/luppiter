package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
		return nil, fmt.Errorf("authorization failed (" + res.Status + ")")
	}

	var user userInfo
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	json.Unmarshal(buf.Bytes(), &user)
	return &user, nil
}

// http://auth.lynlab.co.kr/web/null?access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX3V1aWQiOiI2NDUxMWEyZS0yYTZkLTQyMGEtOTY0OS0yMDFkYmFkNTRkZDEiLCJleHAiOjE1NDY2NzcyMjgsImlhdCI6MTU0NjA3MjQyOCwiaXNzIjoiTFluTGFiL0F1dGgifQ.brEljBT9Gd3Vpeg35mqnyb3sBKJgYXI-YuItqd2ZZNA&refresh_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX3V1aWQiOiI2NDUxMWEyZS0yYTZkLTQyMGEtOTY0OS0yMDFkYmFkNTRkZDEiLCJleHAiOjE1NDg0OTE2MjgsImlhdCI6MTU0NjA3MjQyOCwiaXNzIjoiTFluTGFiL0F1dGgifQ.F3rvWfE4HxU2DH1IuRxiWNYLIbzosEbUzMKq4IFW730
