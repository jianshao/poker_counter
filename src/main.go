package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"private/backend/gamesRoom/src/room"
	"private/backend/gamesRoom/src/user"
	"private/backend/gamesRoom/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Init() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	file, _ := os.OpenFile("gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	router.Use(gin.LoggerWithWriter(file))

	utils.Init()
	room.Init(router)
	user.Init(router)

	router.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func main() {
	Init()
}

// AccessTokenResponse 用于接收 access_token 的响应结构
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

func getAccessToken(appId, appSecret string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appId, appSecret)
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

// QrCodeResponse 用于接收二维码的响应结构
type QrCodeResponse struct {
	QrCodeURL string `json:"url"`
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}

func createMiniProgramQRCode(accessToken, path string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/wxaapp/createwxaqrcode?access_token=%s", accessToken)

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

	appId := "wx7cf078db1b055aef"
	appSecret := "a897857ca6eae5a80f01202cdbeffe64"
	// path := url // 指定跳转的小程序页面路径
	fmt.Printf("path: %s\n", path)

	accessToken, err := getAccessToken(appId, appSecret)
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
