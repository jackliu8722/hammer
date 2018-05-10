package logic

import (
	. "github.com/daoone/hammer/util"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
	"strconv"
	"github.com/daoone/hammer/model"
)

const (
	errPhoneDuplicated = iota + 1000
	errPhoneNumberWrong
	errUsernameDuplicated
	errUsernameWrong
	errSendRateLimit
	errSendCode
	errCodeExpired
)

// 注册
func UserRegistration(name, phoneNum, passwd string) gin.H {
	return gin.H{}
}

// 登录
func UserLogin(phone, passwd string) gin.H {
	return gin.H{}
}

// 发送注册码
func SendUserPhoneCode(phoneNum string) gin.H {
	if !IsPhone(phoneNum) {
		return Error("errPhoneNumberWrong", errPhoneNumberWrong, "Phone number is invalid.")
	}
	var redisPrefix = "phone_code"
	if PhoneNumberExists(phoneNum) {
		return Error("errPhoneDuplicated",
			errPhoneDuplicated, "This phone has been used.")
	}

	rand.Seed(time.Now().UnixNano())
	redisClient := GetRedis()
	defer redisClient.Close()
	// sms has been sent to this phone number
	if redisClient.WithPrefix(redisPrefix).GET(phoneNum) != "" {
		return Error("errSendRateLimit",
			errSendRateLimit, "Code has already be sent.")
	}
	if err := redisClient.WithPrefix(redisPrefix).SETEX(
		phoneNum, 180,
		string([]byte(strconv.Itoa(rand.Int()))[:4])); err != nil {
		return Error("errSendCode", errSendCode, "Error when send code.")
	}

	return Success()
}

// 校验注册码
func CheckPhoneCode(code, phoneNum string) bool {
	redisClient := GetRedis()
	defer redisClient.Close()

	if redisClient.WithPrefix("phone_code").GET(phoneNum) != code {
		return false
	}

	return true
}

// 检测手机号已存在
func PhoneNumberExists(phoneNum string) bool {
	user := &model.User{}
	_, err := model.GetEngine().Where("phone=?", phoneNum).Get(user)
	if err != nil {
		DoLog(err.Error(),"logic_user")
		panic(err.Error())
	}
	if user.Uid == 0 {
		return false
	}

	return true
}
