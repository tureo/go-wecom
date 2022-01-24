package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		ResponseString(c, 200, "welcome to go-wecom")
	})
	r.GET("/ping", func(c *gin.Context) {
		ResponseString(c, 200, "pong")
	})

	wecom := r.Group("/wecom")
	{
		// 验证回调URL
		wecom.GET("", func(c *gin.Context) {
			VerifyURL(c)
		})
		// 处理回调消息
		wecom.POST("", func(c *gin.Context) {
			Callback(c)
		})
	}

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8000")
}
