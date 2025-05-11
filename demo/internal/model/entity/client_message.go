package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ClientMessage 用户消息表结构
type ClientMessage struct {
	Id           int         `json:"id"           description:"消息ID"`
	SenderId     int         `json:"senderId"     description:"发送者ID"`
	SenderName   string      `json:"senderName"   description:"发送者名称"`
	ReceiverId   int         `json:"receiverId"   description:"接收者ID"`
	ReceiverName string      `json:"receiverName" description:"接收者名称"`
	Content      string      `json:"content"      description:"消息内容"`
	IsRead       int         `json:"isRead"       description:"是否已读：0未读，1已读"`
	Status       int         `json:"status"       description:"状态：0删除，1正常"`
	CreatedAt    *gtime.Time `json:"createdAt"    description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    description:"更新时间"`
}
