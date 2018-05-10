package util

import "github.com/gin-gonic/gin"


func Success() gin.H {
	return gin.H{
		"success": true,
	}
}

func Error(name string, code int, message string) gin.H {
	return gin.H{
		"name": name,
		"code": code,
		"message": message,
	}
}

func WxError(msg string) gin.H {
	return gin.H{
		"error": true,
		"message": msg,
	}
}