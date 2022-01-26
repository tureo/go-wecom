package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
)

// 企业微信验证url请求参数
type VerifyURLReq struct {
	MsgSignature string `form:"msg_signature" json:"msg_signature" example:"5c45ff5e21c57e6ad56bac8758b79b1d9ac89fd3"`                                     // 企业微信加密签名
	Timestamp    int    `form:"timestamp" json:"timestamp" example:"1409659589"`                                                                           // 时间戳
	Nonce        string `form:"nonce" json:"nonce" example:"263014780"`                                                                                    // 随机数
	EchoStr      string `form:"echostr" json:"echostr" example:"P9nAzCzyDtyTWESHep1vC5X9xho/qYX3Zpb4yKa9SKld1DsH3Iyt3tP3zNdtp+4RPcs8TgAE7OaBO+FZXvnaqQ=="` // 加密的字符串
}

// 获取token响应字段
type GetTokenResp struct {
	ErrCode     int    `json:"errcode"`      // 出错返回码
	ErrMsg      string `json:"errmsg"`       // 返回码提示语
	AccessToken string `json:"access_token"` // 获取到的凭证
	ExpiresIn   int    `json:"expires_in"`   // 凭证的有效时间（秒）
}

// 发送消息企业微信响应字段
type SendMsgResp struct {
	ErrCode      int    `json:"errcode"`       // 返回码
	ErrMsg       string `json:"errmsg"`        // 对返回码的文本描述内容
	InvalidUser  string `json:"invaliduser"`   // 不合法的userid，不区分大小写，统一转为小写
	InvalidParty string `json:"invalidparty"`  // 不合法的partyid
	InvalidTag   string `json:"invalidtag"`    // 不合法的标签id
	MsgID        string `json:"msgid"`         // 消息id
	ResponseCode string `json:"response_code"` // 仅消息类型为“按钮交互型”，“投票选择型”和“多项选择型”的模板卡片消息返回
}

// 获取服务器ip响应字段
type GetIPResp struct {
	ErrCode int      `json:"errcode"` // 错误码
	ErrMsg  string   `json:"errmsg"`  // 错误信息
	IPList  []string `json:"ip_list"` // 企业微信回调的IP段
}

// ***callback start***//
// 企业微信回调url请求参数
type CallbackReq struct {
	MsgSignature string `form:"msg_signature" json:"msg_signature" example:"5c45ff5e21c57e6ad56bac8758b79b1d9ac89fd3"` // 企业微信加密签名
	Timestamp    int    `form:"timestamp" json:"timestamp" example:"1409659589"`                                       // 时间戳
	Nonce        string `form:"nonce" json:"nonce" example:"263014780"`                                                // 随机数
}

// 企业微信回调消息体解密后的公共字段
type CallbackMsgContentCommon struct {
	ToUserName   string `xml:"ToUserName"`   // 企业微信的CorpID
	FromUserName string `xml:"FromUserName"` // 成员UserID
	CreateTime   int    `xml:"CreateTime"`   // 消息创建时间戳
	MsgType      string `xml:"MsgType"`      // 消息类型
	AgentID      int    `xml:"AgentID"`      // 企业应用的id
}

// ***callback end***//

// GetToken
// @Description: 从企业微信获取应用调用接口token
func GetToken(corpid, corpsecret string) (access_token string, err error) {
	fmt.Println("get token from wecom...")
	getTokenResp := new(GetTokenResp)
	content, err := HttpGet(Host + "/cgi-bin/gettoken?corpid=" + corpid + "&corpsecret=" + corpsecret)
	if err != nil {
		fmt.Println("get token from wecom err: ", err)
		return
	}
	err = json.Unmarshal(content, &getTokenResp)
	if err != nil {
		fmt.Println("get token from wecom err: ", err)
		return
	}
	if getTokenResp.ErrCode != 0 {
		errStr := strconv.Itoa(getTokenResp.ErrCode) + getTokenResp.ErrMsg
		fmt.Println("get token from wecom err: ", errStr)
		return "", errors.New(errStr)
	}
	access_token = getTokenResp.AccessToken
	fmt.Println("get token from wecom success: ", access_token)
	return access_token, nil
}

// SendMsg
// @Description: 发送消息到企业微信接口
func SendMsg(access_token string, body io.Reader) (sendMsgResp *SendMsgResp, err error) {
	fmt.Println("send message to wecom...")

	content, err := HttpPost(Host+"/cgi-bin/message/send?access_token="+access_token, "application/json", body)
	if err != nil {
		fmt.Println("send message to wecom err: ", err)
		return
	}
	sendMsgResp = new(SendMsgResp)
	err = json.Unmarshal(content, &sendMsgResp)
	if err != nil {
		fmt.Println("send message to wecom err: ", err)
		return
	}
	if sendMsgResp.ErrCode != 0 {
		errStr := strconv.Itoa(sendMsgResp.ErrCode) + sendMsgResp.ErrMsg
		fmt.Println("send message to wecom err: ", errStr)
		return sendMsgResp, errors.New(errStr)
	}
	fmt.Println("send message to wecom success: ", sendMsgResp)
	return sendMsgResp, nil
}

// GetIP
// @Description: 获取企业微信服务器的ip段
func GetIP(access_token string) (ipList []string, err error) {
	fmt.Println("get ip from wecom...")
	content, err := HttpGet(Host + "/cgi-bin/getcallbackip?access_token=" + access_token)
	if err != nil {
		fmt.Println("get ip from wecom err: ", err)
		return
	}
	getIPResp := new(GetIPResp)
	err = json.Unmarshal(content, &getIPResp)
	if err != nil {
		fmt.Println("get ip from wecom err: ", err)
		return
	}
	if getIPResp.ErrCode != 0 {
		errStr := strconv.Itoa(getIPResp.ErrCode) + getIPResp.ErrMsg
		fmt.Println("get ip from wecom err: ", errStr)
		return nil, errors.New(errStr)
	}
	fmt.Println("get ip from wecom success: ", getIPResp.IPList)
	return getIPResp.IPList, nil
}

// VerifyURL
// @Description: 验证回调URL
func VerifyURL(c *gin.Context) {
	req := new(VerifyURLReq)

	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("verify url err: ", err)
		ResponseString(c, http.StatusBadRequest, err.Error())
		return
	}

	token := Token                   // 这里是回调URL的token，不是调用接口的access_token
	receiverId := CorpID             // 这里是corpid
	encodingAeskey := EncodingAeskey // 由英文或数字组成且长度为43位的自定义字符串
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(token, encodingAeskey, receiverId, wxbizmsgcrypt.XmlType)

	// 解析出url上的参数值如下：
	verifyMsgSign := req.MsgSignature
	verifyTimestamp := strconv.Itoa(req.Timestamp)
	verifyNonce := req.Nonce
	verifyEchoStr := req.EchoStr

	// 验证并获取明文
	echoStr, cryptErr := wxcpt.VerifyURL(verifyMsgSign, verifyTimestamp, verifyNonce, verifyEchoStr)
	if cryptErr != nil {
		errStr := strconv.Itoa(cryptErr.ErrCode) + cryptErr.ErrMsg
		fmt.Println("verify url fail: ", errStr)
		ResponseString(c, http.StatusBadRequest, errStr)
		return
	}
	fmt.Println("verify url success echoStr: ", string(echoStr))

	ResponseString(c, http.StatusOK, string(echoStr))
}

// Callback
// @Description: 接收企业微信回调业务数据
func Callback(c *gin.Context) {
	req := new(CallbackReq)
	token := Token                   // 这里是回调URL的token，不是调用接口的access_token
	receiverId := CorpID             // 这里是corpid
	encodingAeskey := EncodingAeskey // 由英文或数字组成且长度为43位的自定义字符串
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(token, encodingAeskey, receiverId, wxbizmsgcrypt.XmlType)

	// 解析url参数
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("callback req params err: ", err)
		ResponseString(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("callback req params: ", req)

	// 读取post body xml 数据
	reqData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("callback req body err: ", err)
		ResponseString(c, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("callback req body: ", string(reqData))

	// 验证并获取明文
	msg, cryptErr := wxcpt.DecryptMsg(req.MsgSignature, strconv.Itoa(req.Timestamp), req.Nonce, reqData)
	if cryptErr != nil {
		errStr := strconv.Itoa(cryptErr.ErrCode) + cryptErr.ErrMsg
		fmt.Println("callback decrypt msg err: ", errStr)
		ResponseString(c, http.StatusBadRequest, errStr)
		return
	}
	fmt.Println("callback decrypt msg: ", string(msg))

	// 解析消息内容xml
	var reqMsgContentCommon CallbackMsgContentCommon
	err = xml.Unmarshal(msg, &reqMsgContentCommon)
	if err != nil {
		fmt.Println("callback unmarshal common xml err: ", err)
		ResponseString(c, http.StatusBadRequest, err.Error())
	}
	fmt.Println("callback unmarshal common xml: ", reqMsgContentCommon)

	httpStatus := http.StatusOK
	respMsg := []byte("")
	// ********************************* //
	// 业务开始处理

	switch reqMsgContentCommon.MsgType {
	// 文本消息
	case "text":
		fmt.Println("callback text msg type")
		// 测试回复消息文本
		httpStatus, respMsg, err = CallbackTextTest(c, wxcpt, req, msg)
	// 事件消息
	case "event":
		fmt.Println("callback event msg type")
		// 测试回复更新模板卡片按钮交互文案
		httpStatus, respMsg, err = CallbackTemplateCardButtonTest(c, wxcpt, req, msg)
	// 默认处理
	default:
		fmt.Println("callback default event msg type")
	}

	// 业务处理结束
	// ********************************* //

	ResponseString(c, httpStatus, string(respMsg))
}
