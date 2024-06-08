package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"private/backend/poker_counter/src/config"
)

var (
	URL_GET_OPENID    = "https://api.weixin.qq.com/sns/jscode2session?grant_type=authorization_code"
	URL_GET_TOKEN     = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential"
	URL_CREATE_QRCODE = "https://api.weixin.qq.com/cgi-bin/wxaapp/createwxaqrcode?"
)

// QrCodeResponse 用于接收二维码的响应结构
type QrCodeResponse struct {
	QrCodeURL string `json:"url"`
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}

var result struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// AccessTokenResponse 用于接收 access_token 的响应结构
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

func GetWechatOpenidAndSessionKey(code string) (openid, sessionKey string, err error) {
	url := fmt.Sprintf("%s&appid=%s&secret=%s&js_code=%s", URL_GET_OPENID, config.APP_ID, config.APP_SECRET, code)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", "", err
	}

	if result.ErrCode != 0 {
		return "", "", fmt.Errorf("error from wechat: %d - %s", result.ErrCode, result.ErrMsg)
	}

	return result.Openid, result.SessionKey, nil
}

func getAccessToken(appId, appSecret string) (string, error) {
	url := fmt.Sprintf("%s&appid=%s&secret=%s", URL_GET_TOKEN, appId, appSecret)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp AccessTokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", err
	}

	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("error getting access token: %s", tokenResp.ErrMsg)
	}

	return tokenResp.AccessToken, nil
}

func createMiniProgramQRCode(accessToken, path string) (string, error) {
	url := fmt.Sprintf("%saccess_token=%s", URL_CREATE_QRCODE, accessToken)

	data := map[string]string{
		"path": path,
	}

	j, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(j))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println(string(body))

	var qrCodeResp QrCodeResponse
	err = json.Unmarshal(body, &qrCodeResp)
	if err != nil {
		return "", err
	}

	if qrCodeResp.ErrCode != 0 {
		return "", fmt.Errorf("error creating QR code: %s", qrCodeResp.ErrMsg)
	}

	return qrCodeResp.QrCodeURL, nil
}

func generateQrCode(path string) string {

	// path := url // 指定跳转的小程序页面路径
	fmt.Printf("path: %s\n", path)

	accessToken, err := getAccessToken(config.APP_ID, config.APP_SECRET)
	if err != nil {
		fmt.Println("Error getting access token:", err)
		return ""
	}

	fmt.Printf("Access Token: %s, path: %s", accessToken, path)
	qrCodeURL, err := createMiniProgramQRCode(accessToken, path)
	if err != nil {
		fmt.Println("Error creating QR code:", err)
		return ""
	}

	fmt.Println("QR Code URL:", qrCodeURL)
	return qrCodeURL
}
