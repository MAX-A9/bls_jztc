package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SendMessageReq 发送消息请求
type SendMessageReq struct {
	g.Meta     `path:"/client/message/send" method:"post" tags:"消息" summary:"发送消息" security:"Bearer" description:"发送消息，需要客户端身份验证"`
	ReceiverId int    `json:"receiverId" v:"required" dc:"接收者ID"`
	Content    string `json:"content" v:"required" dc:"消息内容"`
}

// SendMessageRes 发送消息响应
type SendMessageRes struct {
	Id        int         `json:"id" dc:"消息ID"`
	SenderId  int         `json:"senderId" dc:"发送者ID"`
	CreatedAt *gtime.Time `json:"createdAt" dc:"发送时间"`
}

// GetMessagesReq 获取消息列表请求
type GetMessagesReq struct {
	g.Meta   `path:"/client/message/list" method:"get" tags:"消息" summary:"获取消息列表" security:"Bearer" description:"获取消息列表，需要客户端身份验证"`
	TargetId int `json:"targetId" v:"required" dc:"会话对象ID"`
	Page     int `json:"page" v:"min:1" dc:"页码，默认1"`
	Size     int `json:"size" v:"max:50" dc:"每页数量，默认20"`
}

// GetMessagesRes 获取消息列表响应
type GetMessagesRes struct {
	List        []MessageItem `json:"list" dc:"消息列表"`
	TotalCount  int           `json:"totalCount" dc:"总数"`
	TotalPage   int           `json:"totalPage" dc:"总页数"`
	CurrentPage int           `json:"currentPage" dc:"当前页码"`
	Size        int           `json:"size" dc:"每页数量"`
}

// MessageItem 消息项
type MessageItem struct {
	Id             int         `json:"id" dc:"消息ID"`
	SenderId       int         `json:"senderId" dc:"发送者ID"`
	SenderName     string      `json:"senderName" dc:"发送者名称"`
	SenderAvatar   string      `json:"senderAvatar" dc:"发送者头像"`
	ReceiverId     int         `json:"receiverId" dc:"接收者ID"`
	ReceiverName   string      `json:"receiverName" dc:"接收者名称"`
	ReceiverAvatar string      `json:"receiverAvatar" dc:"接收者头像"`
	Content        string      `json:"content" dc:"消息内容"`
	IsRead         int         `json:"isRead" dc:"是否已读：0未读，1已读"`
	CreatedAt      *gtime.Time `json:"createdAt" dc:"创建时间"`
	IsSelf         bool        `json:"isSelf" dc:"是否自己发送的消息"`
}

// GetConversationsReq 获取会话列表请求
type GetConversationsReq struct {
	g.Meta `path:"/client/conversation/list" method:"get" tags:"消息" summary:"获取会话列表" security:"Bearer" description:"获取会话列表，需要客户端身份验证"`
	Page   int `json:"page" v:"min:1" dc:"页码，默认1"`
	Size   int `json:"size" v:"max:50" dc:"每页数量，默认20"`
}

// GetConversationsRes 获取会话列表响应
type GetConversationsRes struct {
	List        []ConversationItem `json:"list" dc:"会话列表"`
	TotalCount  int                `json:"totalCount" dc:"总数"`
	TotalPage   int                `json:"totalPage" dc:"总页数"`
	CurrentPage int                `json:"currentPage" dc:"当前页码"`
	Size        int                `json:"size" dc:"每页数量"`
}

// ConversationItem 会话项
type ConversationItem struct {
	Id           int         `json:"id" dc:"会话ID"`
	TargetId     int         `json:"targetId" dc:"对方ID"`
	TargetName   string      `json:"targetName" dc:"对方名称"`
	TargetAvatar string      `json:"targetAvatar" dc:"对方头像"`
	LastMessage  string      `json:"lastMessage" dc:"最后一条消息内容"`
	UnreadCount  int         `json:"unreadCount" dc:"未读消息数"`
	LastTime     *gtime.Time `json:"lastTime" dc:"最后消息时间"`
}

// ReadMessageReq 标记消息已读请求
type ReadMessageReq struct {
	g.Meta   `path:"/client/message/read" method:"post" tags:"消息" summary:"标记消息已读" security:"Bearer" description:"标记消息已读，需要客户端身份验证"`
	TargetId int `json:"targetId" v:"required" dc:"会话对象ID"`
}

// ReadMessageRes 标记消息已读响应
type ReadMessageRes struct {
	Success bool `json:"success" dc:"是否成功"`
}

// CreateConversationReq 创建会话请求
type CreateConversationReq struct {
	g.Meta   `path:"/client/conversation/create" method:"post" tags:"消息" summary:"创建会话" security:"Bearer" description:"创建会话，需要客户端身份验证"`
	TargetId int `json:"targetId" v:"required" dc:"目标用户ID"`
}

// CreateConversationRes 创建会话响应
type CreateConversationRes struct {
	ConversationId int         `json:"conversationId" dc:"会话ID"`
	TargetId       int         `json:"targetId" dc:"目标用户ID"`
	TargetName     string      `json:"targetName" dc:"目标用户名称"`
	TargetAvatar   string      `json:"targetAvatar" dc:"目标用户头像"`
	CreatedAt      *gtime.Time `json:"createdAt" dc:"创建时间"`
}

// GetUnreadCountReq 获取未读消息数量请求
type GetUnreadCountReq struct {
	g.Meta `path:"/client/message/unread/count" method:"get" tags:"消息" summary:"获取未读消息数量" security:"Bearer" description:"获取未读消息总数量，需要客户端身份验证"`
}

// GetUnreadCountRes 获取未读消息数量响应
type GetUnreadCountRes struct {
	UnreadCount int `json:"unreadCount" dc:"未读消息总数"`
}

// DeleteConversationReq 删除会话请求
type DeleteConversationReq struct {
	g.Meta `path:"/client/conversation/delete" method:"post" tags:"消息" summary:"删除会话" security:"Bearer" description:"删除会话，需要客户端身份验证"`
	Id     int `json:"id" v:"required" dc:"会话ID"`
}

// DeleteConversationRes 删除会话响应
type DeleteConversationRes struct {
	Success bool `json:"success" dc:"是否成功"`
}

// ClearReadConversationsReq 清除已读会话请求
type ClearReadConversationsReq struct {
	g.Meta `path:"/client/conversation/clear-read" method:"post" tags:"消息" summary:"清除已读会话" security:"Bearer" description:"清除已读会话，需要客户端身份验证"`
}

// ClearReadConversationsRes 清除已读会话响应
type ClearReadConversationsRes struct {
	Success bool `json:"success" dc:"是否成功"`
	Count   int  `json:"count" dc:"清除的会话数量"`
}
