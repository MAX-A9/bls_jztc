package client

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "demo/api/client/v1"
	"demo/utility/auth"
)

// MessageController 消息控制器
type MessageController struct{}

// 自定义API接口，需要符合GoFrame框架要求的接口定义
type IMessage interface {
	SendMessage(ctx context.Context, req *v1.SendMessageReq) (res *v1.SendMessageRes, err error)
	GetMessages(ctx context.Context, req *v1.GetMessagesReq) (res *v1.GetMessagesRes, err error)
	GetConversations(ctx context.Context, req *v1.GetConversationsReq) (res *v1.GetConversationsRes, err error)
	ReadMessage(ctx context.Context, req *v1.ReadMessageReq) (res *v1.ReadMessageRes, err error)
	CreateConversation(ctx context.Context, req *v1.CreateConversationReq) (res *v1.CreateConversationRes, err error)
	GetUnreadCount(ctx context.Context, req *v1.GetUnreadCountReq) (res *v1.GetUnreadCountRes, err error)
	DeleteConversation(ctx context.Context, req *v1.DeleteConversationReq) (res *v1.DeleteConversationRes, err error)
	ClearReadConversations(ctx context.Context, req *v1.ClearReadConversationsReq) (res *v1.ClearReadConversationsRes, err error)
}

// 定义API处理结构体
type messageApi struct{}

// 实现消息发送接口
func (a *messageApi) SendMessage(ctx context.Context, req *v1.SendMessageReq) (res *v1.SendMessageRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	return v1.SendMessageReq{}.SendMessage(ctx, req)
}

// 实现获取消息列表接口
func (a *messageApi) GetMessages(ctx context.Context, req *v1.GetMessagesReq) (res *v1.GetMessagesRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	return v1.GetMessagesReq{}.GetMessages(ctx, req)
}

// 实现获取会话列表接口
func (a *messageApi) GetConversations(ctx context.Context, req *v1.GetConversationsReq) (res *v1.GetConversationsRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	return v1.GetConversationsReq{}.GetConversations(ctx, req)
}

// 实现标记消息已读接口
func (a *messageApi) ReadMessage(ctx context.Context, req *v1.ReadMessageReq) (res *v1.ReadMessageRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	return v1.ReadMessageReq{}.ReadMessage(ctx, req)
}

// 实现创建会话接口
func (a *messageApi) CreateConversation(ctx context.Context, req *v1.CreateConversationReq) (res *v1.CreateConversationRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	return v1.CreateConversationReq{}.CreateConversation(ctx, req)
}

// 实现获取未读消息数量接口
func (a *messageApi) GetUnreadCount(ctx context.Context, req *v1.GetUnreadCountReq) (res *v1.GetUnreadCountRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	return v1.GetUnreadCountReq{}.GetUnreadCount(ctx, req)
}

// 实现删除会话接口
func (a *messageApi) DeleteConversation(ctx context.Context, req *v1.DeleteConversationReq) (res *v1.DeleteConversationRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	return v1.DeleteConversationReq{}.DeleteConversation(ctx, req)
}

// 实现清除已读会话接口
func (a *messageApi) ClearReadConversations(ctx context.Context, req *v1.ClearReadConversationsReq) (res *v1.ClearReadConversationsRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	return v1.ClearReadConversationsReq{}.ClearReadConversations(ctx, req)
}

// Register 注册路由
func (c *MessageController) Register(group *ghttp.RouterGroup) {
	api := &messageApi{}

	// 发送消息
	group.Bind(api.SendMessage)

	// 获取消息列表
	group.Bind(api.GetMessages)

	// 获取会话列表
	group.Bind(api.GetConversations)

	// 标记消息已读
	group.Bind(api.ReadMessage)

	// 创建会话
	group.Bind(api.CreateConversation)

	// 获取未读消息数量
	group.Bind(api.GetUnreadCount)

	// 删除会话
	group.Bind(api.DeleteConversation)

	// 清除已读会话
	group.Bind(api.ClearReadConversations)
}
