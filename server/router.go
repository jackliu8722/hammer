package main

import(
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/daoone/hammer/controller"
	"github.com/daoone/hammer/middleware"
)

func (hs *HammerServer)NewRouter() {
	r := hs.router
	v1 := r.Group("/v1")
	{
		// @todo 增加Ratelimit middleware
		v1.POST("/login", controller.Login)
		v1.POST("/register", controller.Register)
		v1.POST("/send_code", controller.SendCode)

		authorize := v1.Group("/")
		authorize.Use(middleware.Auth)

		authorize.GET("/ready", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "Let DAOONE rule the world :)",
			})
		})
	}

	// 微信授权登录
	wx := r.Group("/wx")
	{
		wx.POST("/login", controller.WxLogin)

		oauth := wx.Group("/")
		oauth.Use(middleware.WxCheckLogin)
		oauth.POST("/get_dot", controller.WxGetDot)
		oauth.POST("/get_rp", controller.WxGetRp)
		oauth.POST("/get_dot_rp", controller.WxGetDotAndRp)
		oauth.POST("/get_wallet", controller.WxGetWallet)
		oauth.POST("/get_dotlog", controller.WxGetDotLog)
		oauth.POST("/get_dotlog_detail", controller.WxGetDotLogDetail)
		oauth.POST("/unlock_wallet", controller.WxUnlockWallet)
		oauth.POST("/lock_wallet", controller.WxLockWallet)
		// 第一次创建钱包
		oauth.POST("/create_wallet", controller.WxCreateWallet)
		//oauth.POST("/create_account", controller.WxCreateAccount)
		//oauth.POST("/import_account", controller.WxImportAccount)
		// 交易
		oauth.POST("/do_transfer", controller.WxDoTransfer)

		// 任务接口
		//oauth.POST()

	}

}