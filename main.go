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
	// // 测试发送消息到企业微信
	// go func() {
	// 	time.Sleep(time.Second * 5)
	// 	// 发送模板卡片按钮交互消息
	// 	SendMsgButtonTest()
	// 	time.Sleep(time.Second * 5)
	// 	// 发送文本消息
	// 	SendMsgTextTest()
	// }()
	r.Run(":8000")
}
