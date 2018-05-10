package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/daoone/hammer/logic"
	"fmt"
)

type WxLoginRequest struct {
	Code string `json:"code" form:"code" binding:"required"`
	EncryptedData string `json:"encrypted_data" form:"encrypted_data" binding:"required"`
	Iv string `json:"iv" form:"iv" binding:"required"`
}

type WxTransferRequest struct {
	ToAccount string `json:"to_account" form:"to_account" binding:"required"`
	Amount string `json:"amount" form:"amount" binding:"required"`
	Memo string `json:"memo" form:"memo" binding:"required"`
	Passcode string `json:"passcode" form:"passcode" binding:"required"`
}


func WxLogin(c *gin.Context) {
	var wxl WxLoginRequest
	if err := c.Bind(&wxl); err != nil {
		// 非正常请求不提示详细信息，不用为他人的非正常请求买单
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
	} else {
		token, err := logic.WxCheckLogin(wxl.Code, wxl.EncryptedData, wxl.Iv)
		fmt.Println(token, err)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"success":true, "token": token})
	}
}

func WxGetDot(c *gin.Context) {
	unionId, err := c.Get("unionId")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": true})
	}
	c.JSON(http.StatusOK, logic.WxGetUserDot(unionId.(string)))
}

func WxGetRp(c *gin.Context) {
	unionId, err := c.Get("unionId")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": true})
	}
	c.JSON(http.StatusOK, logic.WxGetUserRp(unionId.(string)))
}

func WxGetDotAndRp(c *gin.Context) {
	unionId, err := c.Get("unionId")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": true})
	}
	c.JSON(http.StatusOK, logic.WxGetUserDotAndRp(unionId.(string)))
}

func WxGetWallet(c *gin.Context) {
	unionId, err := c.Get("unionId")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": true})
	}
	c.JSON(http.StatusOK, logic.WxGetUserWallet(unionId.(string)))
}

func WxCreateWallet(c *gin.Context) {
	unionId, err := c.Get("unionId")
	passcode, err2 := c.GetPostForm("passcode")
	if !err || !err2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true})
	}

	if len(passcode) != 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Passcode new to be 8 characters."})
	}
	c.JSON(http.StatusOK, logic.WxCreateUserWallet(unionId.(string), passcode))
}

func WxGetDotLog(c *gin.Context) {
	unionId, err := c.Get("unionId")
	page, err2 := c.GetPostForm("page")
	if !err || !err2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true})
	}
	c.JSON(http.StatusOK, logic.WxGetUserAccountLog(unionId.(string), page))
}

func WxGetDotLogDetail(c *gin.Context) {

}

func WxUnlockWallet(c *gin.Context) {
	unionId, err := c.Get("unionId")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": true})
	}
	c.JSON(http.StatusOK, logic.WxChangeWalletStatus(unionId.(string), 1))
}

func WxLockWallet(c *gin.Context) {
	unionId, err := c.Get("unionId")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": true})
	}
	c.JSON(http.StatusOK, logic.WxChangeWalletStatus(unionId.(string), 0))
}

func WxDoTransfer(c *gin.Context) {
	unionId, err := c.Get("unionId")
	var wdr WxTransferRequest
	// @todo 添加自定义validator
	if err2 := c.Bind(&wdr); err2 != nil || !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err2.Error()})
	}
	c.JSON(http.StatusOK, logic.WxDoTransfer(unionId.(string), wdr.ToAccount, wdr.Amount, wdr.Memo, wdr.Passcode))
}