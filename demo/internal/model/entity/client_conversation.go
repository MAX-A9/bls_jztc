package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ClientConversation 会话表结构
type ClientConversation struct {
	Id           int         `json:"id"          description:"会话ID"`
	ClientId     int         `json:"clientId"    description:"用户ID"`
	TargetId     int         `json:"targetId"    description:"对方ID"`
	TargetName   string      `json:"targetName"  description:"对方名称"`
	TargetAvatar string      `json:"targetAvatar" description:"对方头像"`
	LastMessage  string      `json:"lastMessage" description:"最后一条消息内容"`
	UnreadCount  int         `json:"unreadCount" description:"未读消息数"`
	LastTime     *gtime.Time `json:"lastTime"    description:"最后消息时间"`
	Status       int         `json:"status"      description:"状态：0删除，1正常"`
	CreatedAt    *gtime.Time `json:"createdAt"   description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt"   description:"更新时间"`
}
