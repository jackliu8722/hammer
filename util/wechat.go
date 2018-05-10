package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
)

var errorCode = map[string]int{
	"IllegalAesKey":     -41001,
	"IllegalIv":         -41002,
	"IllegalBuffer":     -41003,
	"DecodeBase64Error": -41004,
}

type WxSessionData struct {
	SessionKey string `json:"session_key"`
	Openid	string `json:"openid"`
	Unionid string `json:"unionid"`
}

func WxGetSessionKey(code string) (sessionKey, openId string){
	c := url.Values{}
	c.Add("appid", GetConfig("wechat", "appId"))
	c.Add("secret", GetConfig("wechat", "appSecret"))
	c.Add("js_code", code)
	c.Add("grant_type", "authorization_code")

	res, err := http.Get(GetConfig("wechat", "wxUrl") + "?" + c.Encode())
	if err != nil {
		DoLog(err.Error(), "wx_error")
		return
	}
	defer res.Body.Close()
	retByte, _ := ioutil.ReadAll(res.Body)
	var wxsd WxSessionData
	json.Unmarshal(retByte, &wxsd)
	sessionKey = wxsd.SessionKey
	openId = wxsd.Openid

	return
}

// WxBizDataCrypt represents an active WxBizDataCrypt object
type WxBizDataCrypt struct {
	AppID      string
	SessionKey string
}

type showError struct {
	errorCode int
	errorMsg  error
}

func (e showError) Error() string {
	return fmt.Sprintf("{code: %v, error: \"%v\"}", e.errorCode, e.errorMsg)
}

// Decrypt Weixin APP's AES Data
// If isJSON is true, Decrypt return JSON type.
// If isJSON is false, Decrypt return map type.
func (wxCrypt *WxBizDataCrypt) Decrypt(encryptedData string, iv string, isJSON bool) (interface{}, error) {
	if len(wxCrypt.SessionKey) != 24 {
		return nil, showError{errorCode["IllegalAesKey"], errors.New("sessionKey length is error")}
	}
	aesKey, err := base64.StdEncoding.DecodeString(wxCrypt.SessionKey)
	if err != nil {
		return nil, showError{errorCode["DecodeBase64Error"], err}
	}

	if len(iv) != 24 {
		return nil, showError{errorCode["IllegalIv"], errors.New("iv length is error")}
	}
	aesIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, showError{errorCode["DecodeBase64Error"], err}
	}

	aesCipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, showError{errorCode["DecodeBase64Error"], err}
	}
	aesPlantText := make([]byte, len(aesCipherText))

	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, showError{errorCode["IllegalBuffer"], err}
	}

	mode := cipher.NewCBCDecrypter(aesBlock, aesIV)
	mode.CryptBlocks(aesPlantText, aesCipherText)
	aesPlantText = WxPKCS7UnPadding(aesPlantText)

	var decrypted map[string]interface{}
	aesPlantText = []byte(strings.Replace(string(aesPlantText), "\a", "", -1))
	err = json.Unmarshal([]byte(aesPlantText), &decrypted)
	if err != nil {
		return nil, showError{errorCode["IllegalBuffer"], err}
	}

	if decrypted["watermark"].(map[string]interface{})["appid"] != wxCrypt.AppID {
		return nil, showError{errorCode["IllegalBuffer"], errors.New("appId is not match")}
	}

	if isJSON == true {
		return aesPlantText, nil
	}

	return decrypted, nil
}

// WxPKCS7UnPadding return unpadding []Byte plantText
func WxPKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	if unPadding < 1 || unPadding > 32 {
		unPadding = 0
	}
	return plantText[:(length - unPadding)]
}
