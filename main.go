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
		wecom.GET("", func(c *gin.Context) {
			VerifyURL(c)
		})
		wecom.POST("", func(c *gin.Context) {
			ResponseString(c, 200, "wecom post data")
		})
	}

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8000")
}
