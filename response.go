package main

import (
	"github.com/gin-gonic/gin"
)

type ResponseJSONEntry struct {
	Code int         `json:"code"` // 0代表成功;其他代表失败
	Msg  string      `json:"msg"`  // 错误信息,成功时为"OK"
	Data interface{} `json:"data"` // 返回数据
}

// 返回JSON格式
func ResponseJSON(ctx *gin.Context, statusCode int, code int, msg string, data interface{}) {
	ctx.JSON(statusCode, ResponseJSONEntry{
		Code: code,
		Msg:  msg,
		Data: data,
	})
	return
}

// 返回字符串格式
func ResponseString(ctx *gin.Context, statusCode int, msg string) {
	ctx.JSON(statusCode, msg)
	return
}
