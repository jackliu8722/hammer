package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/daoone/hammer/logic"
)

type RegisterRequest struct {
	Name 		string `form:"name" json:"name" binding:"required"`
	Password 	string `form:"password" json:"password" binding:"required"`
	Phone 		string `form:"phone" json:"phone" binding:"required"`
	Code		string `form:"code" json:"code" binding:"required"`
}

type LoginRequest struct {
	Password 	string `form:"password" json:"password" binding:"required"`
	Phone 		string `form:"phone" json:"phone" binding:"required"`
	VerifyCode 	string `form:"verify_code" json:"verify_code"`
}

type SendCodeRequest struct {
	Phone 	string `form:"phone" json:"phone" binding:"required"`
}

type SendCodeResponse struct {

}

func Register(c *gin.Context) {
	var rr RegisterRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, logic.UserRegistration(rr.Name, rr.Phone, rr.Password))
	}
}

func Login(c *gin.Context) {
	var lr LoginRequest
	if err := c.Bind(&lr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		//c.JSON(http.StatusOK, logic.UserLogin(rr.Name, rr.Phone, rr.Password))
	}
}

func SendCode(c *gin.Context) {
	var sdr SendCodeRequest
	if err := c.Bind(&sdr); err != nil {
		// 非正常请求不提示详细信息，不用为他人的非正常请求买单
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, logic.SendUserPhoneCode(sdr.Phone))
	}
}
