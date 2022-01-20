package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
)

// 企业微信请求参数
type WecomVerifyURLReq struct {
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

// GetTokenFromWecom
// @Description: 从企业微信获取应用调用接口token
func GetTokenFromWecom(corpid, corpsecret string) (access_token string, err error) {
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

// SendMsgToWecom
// @Description: 发送消息到企业微信接口
func SendMsgToWecom(access_token string, body io.Reader) (sendMsgResp *SendMsgResp, err error) {
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

// GetIPFromWecom
// @Description: 获取企业微信服务器的ip段
func GetIPFromWecom(access_token string) (ipList []string, err error) {
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
	req := new(WecomVerifyURLReq)

	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("verify url err: ", err)
		ResponseString(c, http.StatusBadRequest, err.Error())
		return
	}

	token := ""
	if Token == "" {
		access_token, err := GetTokenFromWecom(CorpID, AgentSecret)
		if err != nil {
			fmt.Println("verify url err: ", err)
			ResponseString(c, http.StatusInternalServerError, err.Error())
			return
		}
		token = access_token
	} else {

		token = Token
	}
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
