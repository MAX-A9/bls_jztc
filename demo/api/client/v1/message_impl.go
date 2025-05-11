package v1

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"demo/internal/logic/message"
	"demo/utility/auth"
)

// SendMessage 发送消息
func (SendMessageReq) SendMessage(ctx context.Context, req *SendMessageReq) (res *SendMessageRes, err error) {
	// 获取当前登录用户信息
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录")
	}

	// 检查是否自己给自己发消息
	if clientId == req.ReceiverId {
		return nil, gerror.New("不能给自己发送消息")
	}

	// 定义消息ID变量
	var messageId int64

	// 获取接收者信息
	receiverInfo, err := g.Model("client").
		Where("id", req.ReceiverId).
		Where("status", 1).
		Where("deleted_at IS NULL").
		One()
	if err != nil {
		return nil, err
	}
	if receiverInfo.IsEmpty() {
		return nil, gerror.New("接收者不存在或已被禁用")
	}

	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 获取当前登录用户信息
		senderInfo, err := tx.Model("client").
			Where("id", clientId).
			One()
		if err != nil {
			return err
		}

		// 记录当前时间
		nowTime := gtime.Now()

		// 1. 插入消息记录
		result, err := tx.Model("client_message").
			Data(g.Map{
				"sender_id":     clientId,
				"sender_name":   senderInfo["real_name"],
				"receiver_id":   req.ReceiverId,
				"receiver_name": receiverInfo["real_name"],
				"content":       req.Content,
				"is_read":       0, // 0表示未读
				"status":        1, // 1表示正常
				"created_at":    nowTime,
			}).
			Insert()
		if err != nil {
			return err
		}

		// 获取插入的消息ID
		messageId, err = result.LastInsertId()
		if err != nil {
			return err
		}

		// 2. 更新或创建发送者会话
		senderConversation, err := tx.Model("client_conversation").
			Where("client_id", clientId).
			Where("target_id", req.ReceiverId).
			One()
		if err != nil {
			return err
		}

		if senderConversation.IsEmpty() {
			// 创建发送者会话
			_, err = tx.Model("client_conversation").
				Data(g.Map{
					"client_id":     clientId,
					"target_id":     req.ReceiverId,
					"target_name":   receiverInfo["real_name"],
					"target_avatar": receiverInfo["avatar_url"],
					"last_message":  req.Content,
					"unread_count":  0,
					"last_time":     nowTime,
					"status":        1,
					"created_at":    nowTime,
				}).
				Insert()
		} else {
			// 更新发送者会话
			updateData := g.Map{
				"target_name":   receiverInfo["real_name"],
				"target_avatar": receiverInfo["avatar_url"],
				"last_message":  req.Content,
				"last_time":     nowTime,
				"updated_at":    nowTime,
			}

			// 如果会话已被删除，则恢复状态
			if gconv.Int(senderConversation["status"]) == 0 {
				updateData["status"] = 1
			}

			_, err = tx.Model("client_conversation").
				Where("client_id", clientId).
				Where("target_id", req.ReceiverId).
				Data(updateData).
				Update()
		}
		if err != nil {
			return err
		}

		// 3. 更新或创建接收者会话
		// 首先强制获取接收者的会话（即使状态为0）
		receiverConversation, err := tx.Model("client_conversation").
			Where("client_id = ? AND target_id = ?", req.ReceiverId, clientId).
			One()
		if err != nil {
			return err
		}

		if receiverConversation.IsEmpty() {
			// 如果接收者会话不存在，则创建新会话
			_, err = tx.Model("client_conversation").
				Data(g.Map{
					"client_id":     req.ReceiverId,
					"target_id":     clientId,
					"target_name":   senderInfo["real_name"],
					"target_avatar": senderInfo["avatar_url"],
					"last_message":  req.Content,
					"unread_count":  1,
					"last_time":     nowTime,
					"status":        1,
					"created_at":    nowTime,
				}).
				Insert()
		} else {
			// 对于所有存在的接收者会话，无论状态如何，都执行以下更新
			// 确保将删除的会话重新激活
			_, err = tx.Model("client_conversation").
				Where("client_id = ? AND target_id = ?", req.ReceiverId, clientId).
				Data(g.Map{
					"target_name":   senderInfo["real_name"],
					"target_avatar": senderInfo["avatar_url"],
					"last_message":  req.Content,
					"last_time":     nowTime,
					"unread_count":  gconv.Int(receiverConversation["unread_count"]) + 1,
					"status":        1, // 强制设置状态为1（正常）
					"updated_at":    nowTime,
				}).
				Update()
		}
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	res = &SendMessageRes{
		Id:        int(messageId),
		SenderId:  clientId,
		CreatedAt: gtime.Now(),
	}
	return
}

// GetMessages 获取消息列表
func (GetMessagesReq) GetMessages(ctx context.Context, req *GetMessagesReq) (res *GetMessagesRes, err error) {
	// 获取当前登录用户信息
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录")
	}

	// 检查是否查询自己和自己的消息
	if clientId == req.TargetId {
		return nil, gerror.New("不能查询与自己的聊天记录")
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	// 检查对方是否存在并获取对方信息
	targetUser, err := g.Model("client").
		Where("id", req.TargetId).
		Where("status", 1).
		Where("deleted_at IS NULL").
		One()
	if err != nil {
		return nil, err
	}
	if targetUser.IsEmpty() {
		return nil, gerror.New("对方不存在或已被禁用")
	}

	// 获取自己的信息
	selfUser, err := g.Model("client").
		Where("id", clientId).
		One()
	if err != nil {
		return nil, err
	}

	// 检查会话是否存在或已删除
	conversationExists, err := g.Model("client_conversation").
		Where("client_id = ? AND target_id = ? AND status = 1", clientId, req.TargetId).
		Count()
	if err != nil {
		return nil, err
	}
	if conversationExists == 0 {
		// 如果会话不存在或已删除，则返回空列表
		res = &GetMessagesRes{
			List:        []MessageItem{},
			TotalCount:  0,
			TotalPage:   0,
			CurrentPage: req.Page,
			Size:        req.Size,
		}
		return
	}

	// 构建基础查询 - 考虑消息可见性
	condition := "(sender_id = ? AND receiver_id = ? AND sender_visible = 1 AND status = 1) OR " +
		"(sender_id = ? AND receiver_id = ? AND receiver_visible = 1 AND status = 1)"
	params := []interface{}{clientId, req.TargetId, req.TargetId, clientId}

	// 获取总数
	totalCount, err := g.Model("client_message").
		Where(condition, params...).
		Count()
	if err != nil {
		return nil, err
	}

	// 查询消息列表 - 只包括对当前用户可见的消息
	records, err := g.Model("client_message").
		Fields("id, sender_id, sender_name, receiver_id, receiver_name, content, is_read, created_at").
		Where(condition, params...).
		Order("created_at DESC"). // 按时间倒序
		Page(req.Page, req.Size).
		All()
	if err != nil {
		return nil, err
	}

	// 构建消息列表
	var messageList []MessageItem
	if err = gconv.Scan(records, &messageList); err != nil {
		return nil, err
	}

	// 标记消息是否是自己发送的并添加头像
	for i := range messageList {
		messageList[i].IsSelf = messageList[i].SenderId == clientId

		// 根据发送者ID设置头像
		if messageList[i].SenderId == clientId {
			messageList[i].SenderAvatar = gconv.String(selfUser["avatar_url"])
		} else {
			messageList[i].SenderAvatar = gconv.String(targetUser["avatar_url"])
		}

		// 根据接收者ID设置头像
		if messageList[i].ReceiverId == clientId {
			messageList[i].ReceiverAvatar = gconv.String(selfUser["avatar_url"])
		} else {
			messageList[i].ReceiverAvatar = gconv.String(targetUser["avatar_url"])
		}
	}

	// 构建响应
	res = &GetMessagesRes{
		List:        messageList,
		TotalCount:  totalCount,
		TotalPage:   (totalCount + req.Size - 1) / req.Size,
		CurrentPage: req.Page,
		Size:        req.Size,
	}
	return
}

// GetConversations 获取会话列表
func (GetConversationsReq) GetConversations(ctx context.Context, req *GetConversationsReq) (res *GetConversationsRes, err error) {
	// 获取当前登录用户信息
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录")
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	// 查询会话列表
	model := g.Model("client_conversation").
		Where("client_id = ? AND status = 1", clientId).
		Order("last_time DESC")

	// 获取总数
	totalCount, err := model.Count()
	if err != nil {
		return nil, err
	}

	// 分页查询
	records, err := model.Page(req.Page, req.Size).All()
	if err != nil {
		return nil, err
	}

	// 构建会话列表
	var conversationList []ConversationItem
	if err = gconv.Scan(records, &conversationList); err != nil {
		return nil, err
	}

	// 构建响应
	res = &GetConversationsRes{
		List:        conversationList,
		TotalCount:  totalCount,
		TotalPage:   (totalCount + req.Size - 1) / req.Size,
		CurrentPage: req.Page,
		Size:        req.Size,
	}
	return
}

// ReadMessage 标记消息已读
func (ReadMessageReq) ReadMessage(ctx context.Context, req *ReadMessageReq) (res *ReadMessageRes, err error) {
	// 获取当前登录用户信息
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录")
	}

	// 检查是否标记自己发给自己的消息
	if clientId == req.TargetId {
		return nil, gerror.New("不能标记与自己的聊天记录")
	}

	// 检查对方是否存在
	targetExists, err := g.Model("client").
		Where("id", req.TargetId).
		Where("status", 1).
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		return nil, err
	}
	if targetExists == 0 {
		return nil, gerror.New("对方不存在或已被禁用")
	}

	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 更新消息已读状态 - 只标记对方发给自己的消息
		_, err = tx.Model("client_message").
			Where("receiver_id = ? AND sender_id = ? AND is_read = 0",
				clientId, req.TargetId).
			Data(g.Map{"is_read": 1, "updated_at": gtime.Now()}).
			Update()
		if err != nil {
			return err
		}

		// 2. 更新会话未读数
		_, err = tx.Model("client_conversation").
			Where("client_id = ? AND target_id = ?",
				clientId, req.TargetId).
			Data(g.Map{"unread_count": 0, "updated_at": gtime.Now()}).
			Update()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	res = &ReadMessageRes{
		Success: true,
	}
	return
}

// CreateConversation 创建会话
func (CreateConversationReq) CreateConversation(ctx context.Context, req *CreateConversationReq) (res *CreateConversationRes, err error) {
	// 获取当前登录用户信息
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录")
	}

	// 检查是否创建与自己的会话
	if clientId == req.TargetId {
		return nil, gerror.New("不能创建与自己的会话")
	}

	// 检查目标用户是否存在
	targetInfo, err := g.Model("client").
		Fields("id, real_name, avatar_url").
		Where("id", req.TargetId).
		Where("status", 1).
		Where("deleted_at IS NULL").
		One()
	if err != nil {
		return nil, err
	}
	if targetInfo.IsEmpty() {
		return nil, gerror.New("目标用户不存在或已被禁用")
	}

	// 获取当前用户信息
	currentUserInfo, err := g.Model("client").
		Fields("id, real_name, avatar_url").
		Where("id", clientId).
		One()
	if err != nil {
		return nil, err
	}

	// 检查会话是否已存在
	existConversation, err := g.Model("client_conversation").
		Where("client_id", clientId).
		Where("target_id", req.TargetId).
		One()
	if err != nil {
		return nil, err
	}

	var conversationId int64
	nowTime := gtime.Now()

	// 如果会话已存在，直接返回
	if !existConversation.IsEmpty() {
		res = &CreateConversationRes{
			ConversationId: gconv.Int(existConversation["id"]),
			TargetId:       req.TargetId,
			TargetName:     gconv.String(targetInfo["real_name"]),
			TargetAvatar:   gconv.String(targetInfo["avatar_url"]),
			CreatedAt:      gtime.NewFromStr(gconv.String(existConversation["created_at"])),
		}
		return
	}

	// 开启事务创建会话
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 创建会话记录
		result, err := tx.Model("client_conversation").
			Data(g.Map{
				"client_id":     clientId,
				"target_id":     req.TargetId,
				"target_name":   targetInfo["real_name"],
				"target_avatar": targetInfo["avatar_url"],
				"last_message":  "",
				"unread_count":  0,
				"last_time":     nowTime,
				"status":        1,
				"created_at":    nowTime,
			}).
			Insert()
		if err != nil {
			return err
		}

		// 获取会话ID
		conversationId, err = result.LastInsertId()
		if err != nil {
			return err
		}

		// 同时为对方创建会话记录（如果不存在）
		targetConversation, err := tx.Model("client_conversation").
			Where("client_id", req.TargetId).
			Where("target_id", clientId).
			One()
		if err != nil {
			return err
		}

		if targetConversation.IsEmpty() {
			_, err = tx.Model("client_conversation").
				Data(g.Map{
					"client_id":     req.TargetId,
					"target_id":     clientId,
					"target_name":   currentUserInfo["real_name"],
					"target_avatar": currentUserInfo["avatar_url"],
					"last_message":  "",
					"unread_count":  0,
					"last_time":     nowTime,
					"status":        1,
					"created_at":    nowTime,
				}).
				Insert()
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &CreateConversationRes{
		ConversationId: int(conversationId),
		TargetId:       req.TargetId,
		TargetName:     gconv.String(targetInfo["real_name"]),
		TargetAvatar:   gconv.String(targetInfo["avatar_url"]),
		CreatedAt:      nowTime,
	}
	return
}

// GetUnreadCount 获取未读消息数量
func (GetUnreadCountReq) GetUnreadCount(ctx context.Context, req *GetUnreadCountReq) (res *GetUnreadCountRes, err error) {
	// 获取当前登录用户信息
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录")
	}

	// 查询所有会话的未读消息总数
	var count struct {
		Total int `json:"total"`
	}

	err = g.Model("client_conversation").
		Where("client_id = ? AND status = 1", clientId).
		Fields("SUM(unread_count) as total").
		Scan(&count)
	if err != nil {
		return nil, err
	}

	res = &GetUnreadCountRes{
		UnreadCount: count.Total,
	}
	return
}

// DeleteConversation 删除会话
func (DeleteConversationReq) DeleteConversation(ctx context.Context, req *DeleteConversationReq) (res *DeleteConversationRes, err error) {
	// 获取当前登录用户信息
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录")
	}

	// 调用业务逻辑层删除会话
	err = message.DeleteConversation(ctx, req.Id, clientId)
	if err != nil {
		return nil, err
	}

	res = &DeleteConversationRes{
		Success: true,
	}
	return
}

// ClearReadConversations 清除已读会话
func (ClearReadConversationsReq) ClearReadConversations(ctx context.Context, req *ClearReadConversationsReq) (res *ClearReadConversationsRes, err error) {
	// 获取当前登录用户信息
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录")
	}

	// 调用业务逻辑层清除已读会话
	clearCount, err := message.ClearReadConversations(ctx, clientId)
	if err != nil {
		return nil, err
	}

	res = &ClearReadConversationsRes{
		Success: true,
		Count:   clearCount,
	}
	return
}
