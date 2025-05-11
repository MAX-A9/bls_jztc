package message

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	v1 "github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// SendMessage 发送消息
func SendMessage(ctx context.Context, senderId int, senderName string, receiverId int, content string) (id int64, err error) {
	// 获取接收者信息
	receiverInfo, err := g.Model("client").
		Where("id", receiverId).
		Where("status", 1).
		Where("deleted_at IS NULL").
		One()
	if err != nil {
		return 0, err
	}
	if receiverInfo.IsEmpty() {
		return 0, gerror.New("接收者不存在或已被禁用")
	}

	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 插入消息
		messageId, err := tx.Model("client_message").
			Data(g.Map{
				"sender_id":     senderId,
				"sender_name":   senderName,
				"receiver_id":   receiverId,
				"receiver_name": receiverInfo["real_name"],
				"content":       content,
				"is_read":       0,
				"status":        1,
			}).
			InsertAndGetId()
		if err != nil {
			return err
		}

		// 2. 更新或创建发送者与接收者的会话
		// 发送者会话
		senderConversation, err := tx.Model("client_conversation").
			Where("client_id", senderId).
			Where("target_id", receiverId).
			One()
		if err != nil {
			return err
		}

		nowTime := gtime.Now()
		if senderConversation.IsEmpty() {
			// 创建发送者会话
			_, err = tx.Model("client_conversation").
				Data(g.Map{
					"client_id":     senderId,
					"target_id":     receiverId,
					"target_name":   receiverInfo["real_name"],
					"target_avatar": receiverInfo["avatar_url"],
					"last_message":  content,
					"unread_count":  0,
					"last_time":     nowTime,
					"status":        1,
				}).
				Insert()
		} else {
			// 更新发送者会话
			_, err = tx.Model("client_conversation").
				Where("client_id", senderId).
				Where("target_id", receiverId).
				Data(g.Map{
					"target_name":   receiverInfo["real_name"],
					"target_avatar": receiverInfo["avatar_url"],
					"last_message":  content,
					"last_time":     nowTime,
				}).
				Update()
		}
		if err != nil {
			return err
		}

		// 3. 获取发送者信息，用于更新接收者会话
		senderInfo, err := tx.Model("client").
			Where("id", senderId).
			One()
		if err != nil {
			return err
		}

		// 接收者会话
		receiverConversation, err := tx.Model("client_conversation").
			Where("client_id", receiverId).
			Where("target_id", senderId).
			One()
		if err != nil {
			return err
		}

		if receiverConversation.IsEmpty() {
			// 创建接收者会话
			_, err = tx.Model("client_conversation").
				Data(g.Map{
					"client_id":     receiverId,
					"target_id":     senderId,
					"target_name":   senderInfo["real_name"],
					"target_avatar": senderInfo["avatar_url"],
					"last_message":  content,
					"unread_count":  1,
					"last_time":     nowTime,
					"status":        1,
				}).
				Insert()
		} else {
			// 更新接收者会话，并增加未读数
			_, err = tx.Model("client_conversation").
				Where("client_id", receiverId).
				Where("target_id", senderId).
				Data(g.Map{
					"target_name":   senderInfo["real_name"],
					"target_avatar": senderInfo["avatar_url"],
					"last_message":  content,
					"last_time":     nowTime,
					"unread_count":  gdb.Raw("unread_count + 1"),
				}).
				Update()
		}
		if err != nil {
			return err
		}

		id = messageId
		return nil
	})

	return
}

// GetMessages 获取消息列表
func GetMessages(ctx context.Context, clientId, targetId, page, size int) (list []gdb.Record, totalCount int, err error) {
	// 检查对方是否存在
	targetExists, err := g.Model("client").
		Where("id", targetId).
		Where("status", 1).
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		return nil, 0, err
	}
	if targetExists == 0 {
		return nil, 0, gerror.New("对方不存在或已被禁用")
	}

	// 查询消息列表
	model := g.Model("client_message").
		Where("(sender_id = ? AND receiver_id = ? AND status = 1) OR (sender_id = ? AND receiver_id = ? AND status = 1)",
			clientId, targetId, targetId, clientId).
		Order("created_at DESC")

	// 获取总数
	totalCount, err = model.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	list, err = model.Page(page, size).All()
	if err != nil {
		return nil, 0, err
	}

	return
}

// GetConversations 获取会话列表
func GetConversations(ctx context.Context, clientId, page, size int) (list []gdb.Record, totalCount int, err error) {
	// 查询会话列表
	model := g.Model("client_conversation").
		Where("client_id = ? AND status = 1", clientId).
		Order("last_time DESC")

	// 获取总数
	totalCount, err = model.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	list, err = model.Page(page, size).All()
	if err != nil {
		return nil, 0, err
	}

	return
}

// ReadMessages 标记消息已读
func ReadMessages(ctx context.Context, clientId, targetId int) (err error) {
	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 更新消息已读状态
		_, err = tx.Model("client_message").
			Where("receiver_id = ? AND sender_id = ? AND is_read = 0",
				clientId, targetId).
			Data(g.Map{"is_read": 1, "updated_at": v1.Now()}).
			Update()
		if err != nil {
			return err
		}

		// 2. 更新会话未读数
		_, err = tx.Model("client_conversation").
			Where("client_id = ? AND target_id = ?",
				clientId, targetId).
			Data(g.Map{"unread_count": 0, "updated_at": v1.Now()}).
			Update()
		if err != nil {
			return err
		}

		return nil
	})

	return
}

// DeleteConversation 删除会话
func DeleteConversation(ctx context.Context, conversationId, clientId int) (err error) {
	// 查询会话是否存在
	conversation, err := g.Model("client_conversation").
		Where("id = ? AND client_id = ?", conversationId, clientId).
		One()
	if err != nil {
		return err
	}
	if conversation.IsEmpty() {
		return gerror.New("会话不存在或已被删除")
	}

	targetId := gconv.Int(conversation["target_id"])

	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 标记会话为删除状态
		_, err = tx.Model("client_conversation").
			Where("id = ? AND client_id = ?", conversationId, clientId).
			Data(g.Map{"status": 0, "updated_at": v1.Now()}).
			Update()
		if err != nil {
			return err
		}

		// 2. 标记该用户相关的消息为不可见（不影响对方可见性）
		// 如果clientId是发送者，则设置sender_visible为0
		_, err = tx.Model("client_message").
			Where("sender_id = ? AND receiver_id = ?", clientId, targetId).
			Data(g.Map{"sender_visible": 0, "updated_at": v1.Now()}).
			Update()
		if err != nil {
			return err
		}

		// 如果clientId是接收者，则设置receiver_visible为0
		_, err = tx.Model("client_message").
			Where("sender_id = ? AND receiver_id = ?", targetId, clientId).
			Data(g.Map{"receiver_visible": 0, "updated_at": v1.Now()}).
			Update()
		if err != nil {
			return err
		}

		// 3. 将同时满足sender_visible=0和receiver_visible=0的消息状态设置为0（删除状态）
		_, err = tx.Model("client_message").
			Where("((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND sender_visible = 0 AND receiver_visible = 0",
				clientId, targetId, targetId, clientId).
			Data(g.Map{"status": 0, "updated_at": v1.Now()}).
			Update()
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// ClearReadConversations 清除已读会话
func ClearReadConversations(ctx context.Context, clientId int) (clearCount int, err error) {
	// 查询所有未读消息数为0的会话
	readConversations, err := g.Model("client_conversation").
		Where("client_id = ? AND status = 1 AND unread_count = 0", clientId).
		All()
	if err != nil {
		return 0, err
	}

	if len(readConversations) == 0 {
		return 0, nil
	}

	// 收集所有需要标记删除的会话ID和目标用户ID
	var conversationIds []int
	targetUserMap := make(map[int]bool)

	for _, conv := range readConversations {
		conversationIds = append(conversationIds, gconv.Int(conv["id"]))
		targetUserMap[gconv.Int(conv["target_id"])] = true
	}

	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 标记会话为删除状态
		result, err := tx.Model("client_conversation").
			Where("client_id = ? AND id IN(?)", clientId, conversationIds).
			Data(g.Map{"status": 0, "updated_at": v1.Now()}).
			Update()
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		clearCount = int(rowsAffected)

		// 2. 为每个目标用户处理消息可见性
		for targetId := range targetUserMap {
			// 如果clientId是发送者，则设置sender_visible为0
			_, err = tx.Model("client_message").
				Where("sender_id = ? AND receiver_id = ?", clientId, targetId).
				Data(g.Map{"sender_visible": 0, "updated_at": v1.Now()}).
				Update()
			if err != nil {
				return err
			}

			// 如果clientId是接收者，则设置receiver_visible为0
			_, err = tx.Model("client_message").
				Where("sender_id = ? AND receiver_id = ?", targetId, clientId).
				Data(g.Map{"receiver_visible": 0, "updated_at": v1.Now()}).
				Update()
			if err != nil {
				return err
			}

			// 3. 将同时满足sender_visible=0和receiver_visible=0的消息状态设置为0（删除状态）
			_, err = tx.Model("client_message").
				Where("((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND sender_visible = 0 AND receiver_visible = 0",
					clientId, targetId, targetId, clientId).
				Data(g.Map{"status": 0, "updated_at": v1.Now()}).
				Update()
			if err != nil {
				return err
			}
		}

		return nil
	})

	return
}
