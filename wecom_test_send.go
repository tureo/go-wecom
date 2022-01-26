package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// ***send msg start***//
//发送消息公共字段
type SendMsgCommon struct {
	ToUser  string `json:"touser"`  // 成员ID列表
	ToParty string `json:"toparty"` // 部门ID列表
	ToTag   string `json:"totag"`   // 标签ID列表
	MsgType string `json:"msgtype"` // 消息类型
	AgentID int    `json:"agentid"` // 企业应用的id
}

// ***文本消息 start***//
// 发送文本消息字段
type SendMsgText struct {
	SendMsgCommon
	Text                   Text `json:"text"`                     // 消息内容
	Safe                   int  `json:"safe"`                     // 是否是保密消息
	EnableIDTrans          int  `json:"enable_id_trans"`          // 是否开启id转译
	EnableDuplicateCheck   int  `json:"enable_duplicate_check"`   // 是否开启重复消息检查
	DuplicateCheckInterval int  `json:"duplicate_check_interval"` // 是否重复消息检查的时间间隔
}

// 文本消息内容
type Text struct {
	Content string `json:"content"` // 消息内容
}

// ***文本消息 end***//

// ***模板卡片消息 start***//
// 发送按钮交互型消息字段
type SendMsgTemplateCardButton struct {
	SendMsgCommon
	TemplateCard           TemplateCard `json:"template_card"`            // 模板卡片
	EnableIDTrans          int          `json:"enable_id_trans"`          // 是否开启id转译
	EnableDuplicateCheck   int          `json:"enable_duplicate_check"`   // 是否开启重复消息检查
	DuplicateCheckInterval int          `json:"duplicate_check_interval"` // 是否重复消息检查的时间间隔
}

// 一级标题
type MainTitle struct {
	Title string `json:"title"` // 一级标题
}

// 二级标题+文本
type HorizontalContent struct {
	Keyname string `json:"keyname"` // 二级标题
	Value   string `json:"value"`   // 二级文本
	Type    int    `json:"type"`    // 链接类型
	UserID  string `json:"userid"`  // 成员详情的userid
}

// 按钮
type Button struct {
	Text  string `json:"text"`  // 按钮文案
	Style int    `json:"style"` // 按钮样式
	Key   string `json:"key"`   // 按钮key值
}

// 模板卡片消息内容
type TemplateCard struct {
	CardType              string              `json:"card_type"`               // 模板卡片类型
	MainTitle             MainTitle           `json:"main_title"`              // 一级标题
	SubTitleText          string              `json:"sub_title_text"`          // 二级普通文本
	HorizontalContentList []HorizontalContent `json:"horizontal_content_list"` // 二级标题+文本列表
	TaskID                string              `json:"task_id"`                 // 任务id
	ButtonList            []Button            `json:"button_list"`             // 按钮列表
}

// ***模板卡片消息 end***//

// ***send msg end***//

// SendMsgButtonTest
// @Description: 测试发送模板卡片按钮交互消息到企业微信
func SendMsgButtonTest() {
	access_token, err := GetToken(CorpID, AgentSecret)
	if err != nil {
		fmt.Println("send msg button err: ", err)
	}
	fmt.Println("send msg button access_token: ", access_token)
	sendMsg := SendMsgTemplateCardButton{
		SendMsgCommon: SendMsgCommon{
			ToUser:  UserID,
			ToParty: "",
			ToTag:   "",
			MsgType: "template_card",
			AgentID: AgentID,
		},
		TemplateCard: TemplateCard{
			CardType: "button_interaction",
			MainTitle: MainTitle{
				Title: "abc后端项目",
			},
			SubTitleText: "abc：测试企业微信应用发送模板卡片（按钮交互型）消息" + strconv.FormatInt(time.Now().Unix(), 10),
			HorizontalContentList: []HorizontalContent{
				{
					Keyname: "游戏",
					Value:   "abc1111",
				},
				{
					Keyname: "任务id",
					Value:   "88880001",
				},
				{
					Type:    3,
					Keyname: "提交人",
					Value:   "点击查看",
					UserID:  UserID,
				},
			},
			TaskID: "task_id_multi_user_test" + strconv.FormatInt(time.Now().Unix(), 10),
			ButtonList: []Button{
				{
					Text:  "通过",
					Style: 1,
					Key:   "approve",
				},
				{
					Text:  "驳回",
					Style: 2,
					Key:   "reject",
				},
			},
		},
		EnableIDTrans:          0,
		EnableDuplicateCheck:   0,
		DuplicateCheckInterval: 1800,
	}
	sendMsgJSON, err := json.Marshal(sendMsg)
	if err != nil {
		fmt.Println("send msg button json err: ", err)
	}
	fmt.Println("send msg button json: ", string(sendMsgJSON))
	sendMsgResp, err := SendMsg(access_token, bytes.NewReader(sendMsgJSON))
	if err != nil {
		fmt.Println("send msg text err: ", err)
	}
	fmt.Println("send msg text resp: ", sendMsgResp)
}

// SendMsgTextTest
// @Description: 测试发送文本消息到企业微信
func SendMsgTextTest() {
	access_token, err := GetToken(CorpID, AgentSecret)
	if err != nil {
		fmt.Println("send msg text err: ", err)
	}
	fmt.Println("send msg text access_token: ", access_token)
	sendMsg := SendMsgText{
		SendMsgCommon: SendMsgCommon{
			ToUser:  UserID,
			ToParty: "",
			ToTag:   "",
			MsgType: "text",
			AgentID: AgentID,
		},
		Text: Text{
			Content: "abc：测试企业微信应用发送文本消息" + strconv.FormatInt(time.Now().Unix(), 10),
		},
		Safe:                   0,
		EnableIDTrans:          0,
		EnableDuplicateCheck:   0,
		DuplicateCheckInterval: 1800,
	}

	sendMsgJSON, err := json.Marshal(sendMsg)
	if err != nil {
		fmt.Println("send msg text json err: ", err)
	}
	fmt.Println("send msg text json: ", string(sendMsgJSON))
	sendMsgResp, err := SendMsg(access_token, bytes.NewReader(sendMsgJSON))
	if err != nil {
		fmt.Println("send msg text err: ", err)
	}
	fmt.Println("send msg text resp: ", sendMsgResp)
}
