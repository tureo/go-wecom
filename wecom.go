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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
)

// 企业微信请求参数
type VerifyURLReq struct {
	MsgSignature string `form:"msg_signature" json:"msg_signature" example:"5c45ff5e21c57e6ad56bac8758b79b1d9ac89fd3"`                                     // 企业微信加密签名
	Timestamp    int    `form:"timestamp" json:"timestamp" example:"1409659589"`                                                                           // 时间戳
	Nonce        string `form:"nonce" json:"nonce" example:"263014780"`                                                                                    // 随机数
	EchoStr      string `form:"echostr" json:"echostr" example:"P9nAzCzyDtyTWESHep1vC5X9xho/qYX3Zpb4yKa9SKld1DsH3Iyt3tP3zNdtp+4RPcs8TgAE7OaBO+FZXvnaqQ=="` // 加密的字符串
}

// 获取token响应结构体
type GetTokenResp struct {
	ErrCode     int    `json:"errcode"`      // 出错返回码
	ErrMsg      string `json:"errmsg"`       // 返回码提示语
	AccessToken string `json:"access_token"` // 获取到的凭证
	ExpiresIn   int    `json:"expires_in"`   // 凭证的有效时间（秒）
}

// 发送消息响应结构体
type SendMsgResp struct {
	ErrCode      int    `json:"errcode"`       // 返回码
	ErrMsg       string `json:"errmsg"`        // 对返回码的文本描述内容
	InvalidUser  string `json:"invaliduser"`   // 不合法的userid，不区分大小写，统一转为小写
	InvalidParty string `json:"invalidparty"`  // 不合法的partyid
	InvalidTag   string `json:"invalidtag"`    // 不合法的标签id
	MsgID        string `json:"msgid"`         // 消息id
	ResponseCode string `json:"response_code"` // 仅消息类型为“按钮交互型”，“投票选择型”和“多项选择型”的模板卡片消息返回
}

// 发送文本消息
type SendMsgText struct {
	Touser                 string `json:"touser"`                   // 指定接收消息的成员，成员ID列表
	Toparty                string `json:"toparty"`                  // 指定接收消息的部门，部门ID列表
	ToTag                  string `json:"totag"`                    // 指定接收消息的标签，标签ID列表
	MsgType                string `json:"msgtype"`                  // 消息类型
	AgentID                int    `json:"agentid"`                  // 企业应用的id
	Text                   Text   `json:"text"`                     // 消息内容
	Safe                   int    `json:"safe"`                     // 是否是保密消息
	EnableIDTrans          int    `json:"enable_id_trans"`          // 是否开启id转译
	EnableDuplicateCheck   int    `json:"enable_duplicate_check"`   // 是否开启重复消息检查
	DuplicateCheckInterval int    `json:"duplicate_check_interval"` // 是否重复消息检查的时间间隔
}

// 文本消息内容
type Text struct {
	Content string `json:"content"` // 消息内容
}

// 获取服务器ip响应结构体
type GetIPResp struct {
	ErrCode int      `json:"errcode"` // 错误码
	ErrMsg  string   `json:"errmsg"`  // 错误信息
	IPList  []string `json:"ip_list"` // 企业微信回调的IP段
}

// 企业微信回调url请求参数
type CallbackReq struct {
	MsgSignature string `form:"msg_signature" json:"msg_signature" example:"5c45ff5e21c57e6ad56bac8758b79b1d9ac89fd3"` // 企业微信加密签名
	Timestamp    int    `form:"timestamp" json:"timestamp" example:"1409659589"`                                       // 时间戳
	Nonce        string `form:"nonce" json:"nonce" example:"263014780"`                                                // 随机数
}

// 企业微信回调消息体解密后的数据
type ReqMsgContent struct {
	ToUserName   string `xml:"ToUserName"`   // 企业微信的CorpID
	AgentID      string `xml:"AgentID"`      // 接收的应用id
	FromUserName string `xml:"FromUserName"` // 发送者的userid
	EventKey     string `xml:"EventKey"`     // 按钮key值
	TaskID       string `xml:"TaskId"`       // 任务id
}

// // 企业微信回调被动响应包文本数据
// type RespText struct {
// 	XMLName      xml.Name `xml:"xml"`
// 	ToUserName   string   `xml:"ToUserName"`   // 成员UserID
// 	FromUserName string   `xml:"FromUserName"` // 企业微信CorpID
// 	CreateTime   int      `xml:"CreateTime"`   // 消息创建时间
// 	MsgType      string   `xml:"MsgType"`      // 消息类型
// 	Content      string   `xml:"Content"`      // 文本消息内容
// }

// 企业微信回调被动响应包更新按钮文案数据
type RespUpdateButton struct {
	XMLName      xml.Name            `xml:"xml"`
	ToUserName   string              `xml:"ToUserName"`   // 成员UserID
	FromUserName string              `xml:"FromUserName"` // 企业微信CorpID
	CreateTime   int                 `xml:"CreateTime"`   // 消息创建时间
	MsgType      string              `xml:"MsgType"`      // 消息类型
	Button       UpdateButtonReplace `xml:"Button"`       // 按钮
}

// 更新按钮文案
type UpdateButtonReplace struct {
	ReplaceName string `xml:"ReplaceName"` // 点击卡片按钮后显示的按钮名称
}

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
		fmt.Println("verify url fail: ", cryptErr.ErrMsg)
		ResponseString(c, http.StatusBadRequest, cryptErr.ErrMsg)
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
		fmt.Println("callback decrypt msg err: ", cryptErr.ErrMsg)
		ResponseString(c, http.StatusBadRequest, cryptErr.ErrMsg)
		return
	}
	fmt.Println("callback decrypt msg: ", string(msg))

	// 解析xml
	var reqMsgContent ReqMsgContent
	err = xml.Unmarshal(msg, &reqMsgContent)
	if err != nil {
		fmt.Println("callback unmarshal xml err: ", err)
		ResponseString(c, http.StatusBadRequest, err.Error())
	}
	fmt.Println("callback unmarshal xml: ", reqMsgContent)

	// ********************************* //
	// 业务处理
	// TODO
	var buttonReplaceText = ""
	if reqMsgContent.EventKey == "approve" {
		buttonReplaceText = "审核已通过"
	}
	if reqMsgContent.EventKey == "reject" {
		buttonReplaceText = "审核已驳回"
	}

	// ********************************* //
	// 构造被动响应消息
	respMsgContent := &RespUpdateButton{
		ToUserName:   reqMsgContent.FromUserName,
		FromUserName: reqMsgContent.ToUserName,
		CreateTime:   int(time.Now().Unix()),
		MsgType:      "update_button",
		Button: UpdateButtonReplace{
			ReplaceName: buttonReplaceText,
		},
	}

	// 消息内容转成xml文本
	respMsgContentXML, err := xml.Marshal(respMsgContent)
	if err != nil {
		fmt.Println("callback marshal xml err: ", err)
		ResponseString(c, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("callback marshal xml : ", string(respMsgContentXML))

	// 加密签名
	encryptMsg, cryptErr := wxcpt.EncryptMsg(string(respMsgContentXML), strconv.Itoa(req.Timestamp), req.Nonce)
	if cryptErr != nil {
		fmt.Println("callback encrypt msg eror: ", cryptErr.ErrMsg)
		ResponseString(c, http.StatusInternalServerError, cryptErr.ErrMsg)
	}
	fmt.Println("callback encrypt msg : ", string(encryptMsg))

	ResponseString(c, http.StatusOK, string(encryptMsg))
}
