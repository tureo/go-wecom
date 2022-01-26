package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
)

// 企业微信回调模板卡片按钮消息体解密后的数据
type ReqMsgContentTemplateCardButton struct {
	CallbackMsgContentCommon
	EventKey string `xml:"EventKey"` // 按钮key值
	TaskID   string `xml:"TaskId"`   // 任务id
}

// 企业微信回调文本消息体解密后的数据
type ReqMsgContentText struct {
	CallbackMsgContentCommon
	Content string `xml:"Content"` // 文本消息内容
	MsgID   int64  `xml:"MsgId"`   // 消息id
}

// 企业微信回调被动响应包文本数据
type RespText struct {
	CallbackMsgContentCommon
	XMLName xml.Name `xml:"xml"`
	Content string   `xml:"Content"` // 文本消息内容
}

// 企业微信回调被动响应包更新按钮文案数据
type RespUpdateButton struct {
	CallbackMsgContentCommon
	XMLName xml.Name            `xml:"xml"`
	Button  UpdateButtonReplace `xml:"Button"` // 按钮
}

// 更新按钮文案
type UpdateButtonReplace struct {
	ReplaceName string `xml:"ReplaceName"` // 点击卡片按钮后显示的按钮名称
}

// CallbackTemplateCardButtonTest
// @Description: 测试企业微信模板卡片按钮消息回调处理
func CallbackTemplateCardButtonTest(c *gin.Context, wxcpt *wxbizmsgcrypt.WXBizMsgCrypt, req *CallbackReq, msg []byte) (httpStatus int, encryptMsg []byte, err error) {
	// 解析消息内容xml
	reqMsgContent := new(ReqMsgContentTemplateCardButton)
	err = xml.Unmarshal(msg, &reqMsgContent)
	if err != nil {
		fmt.Println("callback unmarshal xml err: ", err)
		return http.StatusBadRequest, nil, err
	}
	fmt.Println("callback unmarshal xml: ", reqMsgContent)

	// 业务逻辑处理
	buttonReplaceText := ""
	if reqMsgContent.EventKey == "approve" {
		buttonReplaceText = "审核已通过"
	}
	if reqMsgContent.EventKey == "reject" {
		buttonReplaceText = "审核已驳回"
	}

	// 构造被动响应消息
	respMsgContent := new(RespUpdateButton)
	respMsgContent.ToUserName = reqMsgContent.FromUserName
	respMsgContent.FromUserName = reqMsgContent.ToUserName
	respMsgContent.CreateTime = int(time.Now().Unix())
	respMsgContent.MsgType = "update_button"
	respMsgContent.Button = UpdateButtonReplace{
		ReplaceName: buttonReplaceText,
	}

	// 消息内容转成xml文本
	respMsgContentXML, err := xml.Marshal(respMsgContent)
	if err != nil {
		fmt.Println("callback marshal xml err: ", err)
		return http.StatusInternalServerError, nil, err
	}
	fmt.Println("callback marshal xml : ", string(respMsgContentXML))

	// 加密签名
	encryptMsg, cryptErr := wxcpt.EncryptMsg(string(respMsgContentXML), strconv.Itoa(req.Timestamp), req.Nonce)
	if cryptErr != nil {
		errStr := strconv.Itoa(cryptErr.ErrCode) + cryptErr.ErrMsg
		fmt.Println("callback encrypt msg err: ", errStr)
		return http.StatusInternalServerError, nil, errors.New(errStr)
	}
	fmt.Println("callback encrypt msg : ", string(encryptMsg))

	return http.StatusOK, encryptMsg, nil
}

// CallbackTextTest
// @Description: 测试企业微信文本消息回调处理
func CallbackTextTest(c *gin.Context, wxcpt *wxbizmsgcrypt.WXBizMsgCrypt, req *CallbackReq, msg []byte) (httpStatus int, encryptMsg []byte, err error) {
	// 解析消息内容xml
	reqMsgContent := new(ReqMsgContentText)
	err = xml.Unmarshal(msg, &reqMsgContent)
	if err != nil {
		fmt.Println("callback unmarshal xml err: ", err)
		return http.StatusBadRequest, nil, err
	}
	fmt.Println("callback unmarshal xml: ", reqMsgContent)

	// 业务逻辑处理
	content := reqMsgContent.Content + ", get it!"

	// 构造被动响应消息
	respMsgContent := new(RespText)
	respMsgContent.ToUserName = reqMsgContent.FromUserName
	respMsgContent.FromUserName = reqMsgContent.ToUserName
	respMsgContent.CreateTime = int(time.Now().Unix())
	respMsgContent.MsgType = "text"
	respMsgContent.Content = content

	// 消息内容转成xml文本
	respMsgContentXML, err := xml.Marshal(respMsgContent)
	if err != nil {
		fmt.Println("callback marshal xml err: ", err)
		return http.StatusInternalServerError, nil, err
	}
	fmt.Println("callback marshal xml : ", string(respMsgContentXML))

	// 加密签名
	encryptMsg, cryptErr := wxcpt.EncryptMsg(string(respMsgContentXML), strconv.Itoa(req.Timestamp), req.Nonce)
	if cryptErr != nil {
		errStr := strconv.Itoa(cryptErr.ErrCode) + cryptErr.ErrMsg
		fmt.Println("callback encrypt msg err: ", errStr)
		return http.StatusInternalServerError, nil, errors.New(errStr)
	}
	fmt.Println("callback encrypt msg : ", string(encryptMsg))

	return http.StatusOK, encryptMsg, nil
}
