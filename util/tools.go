package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"errors"
	"regexp"
	"encoding/gob"
	"time"
	"fmt"
	"math/rand"
	"crypto/md5"
	"io"
	"crypto/aes"
	"crypto/cipher"
)

// 公用请求
func DoPost(url string, payload interface{}) (a []byte, e error) {
	jsonValue, _ := json.Marshal(payload)
	res, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonValue),
	)
	if err != nil {
		DoLog(err.Error(), "tool_doPost")
		return a, err
	}
	defer res.Body.Close()
	retByte, _ := ioutil.ReadAll(res.Body)
	if (res.StatusCode < http.StatusOK) || (res.StatusCode > http.StatusIMUsed) {
		return a, errors.New(string(retByte))
	}

	return retByte, nil
}

func DoGet(url string) (a []byte, e error) {
	res, err := http.Get(url)
	if err != nil {
		DoLog(err.Error(), "tool_doPost")
		return a, err
	}
	defer res.Body.Close()
	retByte, _ := ioutil.ReadAll(res.Body)
	if (res.StatusCode < http.StatusOK) || (res.StatusCode > http.StatusIMUsed) {
		return a, errors.New(string(retByte))
	}

	return retByte, nil
}

// 正则验证手机号
func IsPhone(phone string) bool {
	reg := `^1[3|4|5|7|8][0-9]{9}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

// interface转byte
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MakeRandomStr(l int) string {
	chars := "abcdefghijkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789"
	clen := float64(len(chars))
	res := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < l; i++ {
		rfi := int(clen * rand.Float64())
		res += fmt.Sprintf("%c", chars[rfi])
	}

	return res
}

func MakeRandomStrLower(l int) string {
	chars := "12345abcdefghijklmnopqrstuvwxyz"
	clen := float64(len(chars))
	res := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < l; i++ {
		rfi := int(clen * rand.Float64())
		res += fmt.Sprintf("%c", chars[rfi])
	}

	return res
}

func DotUnitToInt(am float64) int64 {
	am = am * 10000
	i := int(am)
	return int64(i)
}

func DotUnitToFloat(am int64) float64 {

	f := float64(am)
	f = f/10000
	return f
}

func ByteToString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

func B2S(bs []int8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}

func Md5(text string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, text)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}


// AES normal
func Encrypt(plantText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}
	plantText = PKCS7Padding(plantText, block.BlockSize())
	blockModel := cipher.NewCBCEncrypter(block, key)
	ciphertext := make([]byte, len(plantText))
	blockModel.CryptBlocks(ciphertext, plantText)
	return ciphertext, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}


func Decrypt(ciphertext, key []byte) ([]byte, error) {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes) //选择加密算法
	if err != nil {
		return nil, err
	}
	blockModel := cipher.NewCBCDecrypter(block, keyBytes)
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)
	plantText = PKCS7UnPadding(plantText, block.BlockSize())
	return plantText, nil
}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}
