package alert

import (
	"bytes"
	"context"
	"tooladmin/server/common/logger"
	"tooladmin/server/scheduler/model"

	"text/template"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var ScheduleFailMsg = `

{
  "schema": "2.0",
  "config": {
    "wide_screen_mode": true
  },
  "header": {
    "title": {
      "tag": "plain_text",
      "content": "[scheduler]任务执行失败通知"
    }
  },
  "body": {
    "elements": [
	 {
                    "tag": "markdown",
                    "content": "调度任务** {{.TaskName}} **执行失败",
                    "text_align": "left",
                    "text_size": "normal"
                   
},
      {
        "tag": "button",
        "element_id": "account_doc_url",
        "margin": "0px 0px 0px 0px",
        "type": "primary",
        "size": "small",
        "width": "default",
        "text": {
          "tag": "plain_text",
          "content": "查看详情"
        },
        "icon": {
          "tag": "standard_icon",
          "token": "link-copy_outlined"
        },
        "hover_tips": {},
        "disabled": false,
        "disabled_tips": {},
        "behaviors": [
          {
            "type": "open_url",
            "default_url": "{{.LogUrl}}"
          }
        ]
      }
    ]
  }
}

`
var scheduleFailMsgTpl *template.Template

func init() {
	tpl, err := template.New("schedule_fail_msg").Parse(ScheduleFailMsg)
	if err != nil {
		panic(err)
	}
	scheduleFailMsgTpl = tpl
}

func SendScheduleFailMsg(ctx context.Context, msgData *model.ScheduleFailMsg) {
	// 检查LarkClient是否已初始化
	if LarkClient == nil {
		logger.Error(ctx, "SendScheduleFailMsg, LarkClient is not initialized")
		return
	}
	// 检查是否提供了接收者ID
	if len(msgData.ReceiveIds) == 0 {
		logger.Error(ctx, "SendScheduleFailMsg, ReceiveId is empty")
		return
	}
	var buffer bytes.Buffer
	err := scheduleFailMsgTpl.Execute(&buffer, msgData)
	if err != nil {
		logger.Error(ctx, "SendScheduleFailMsg,template execute err:", err)
		return
	}
	msgStr := buffer.String()
	logger.Info(ctx, "SendScheduleFailMsg,msg:", msgStr)

	// 发送消息
	for _, receiveId := range msgData.ReceiveIds {
		resp, err := LarkClient.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
			ReceiveIdType(larkim.ReceiveIdTypeOpenId).
			Body(larkim.NewCreateMessageReqBodyBuilder().
				MsgType(larkim.MsgTypeInteractive).
				Content(msgStr).ReceiveId(receiveId).
				Build()).
			Build())

		if err != nil || !resp.Success() {
			logger.Info(ctx, "SendScheduleFailMsg error:", err, "response:", resp)
		} else {
			logger.Info(ctx, "SendScheduleFailMsg success")
		}
	}
}
