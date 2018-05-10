package middleware

import(
	"github.com/gin-gonic/gin"
	. "github.com/daoone/hammer/util"
	"net/http"
	"net/url"
)

type WxSessionRequest struct {
	Token string `json:"token" form:"token"`
}

func WxCheckLogin(c *gin.Context) {
	if token, ok := c.GetPostForm("token"); ok {
		token, _ = url.QueryUnescape(token)
		redisClient := GetRedis()
		defer redisClient.Close()
		if  val2, _ := redisClient.WithPrefix("wx").HGET(token, "status");val2 == "1" {
			unionId, _ := redisClient.WithPrefix("wx").HGET(token, "union_id")
			c.Set("unionId", unionId)
			c.Set("token", token)
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired", "message": "请重新登录"})
		c.Abort()
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "token missing", "message": "请重新登录"})
}

// 发送code，获取sessionkey
func WxCheckLoginWithCode(c *gin.Context) {
	// code 换 sessionKey
	if val, ok := c.Get("code"); ok {
		sessionKey, _ := WxGetSessionKey(val.(string))
		redisClient := GetRedis()
		defer redisClient.Close()
		if val, _ := redisClient.WithPrefix("wx:").HGET(sessionKey, "status");val == "1" {
			c.Next()
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session Expired."})
		c.Abort()
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Code missing."})
}
