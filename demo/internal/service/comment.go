package service

import (
	"context"

	v1 "demo/api/comment/v1"
)

// CommentService 评论服务接口
type CommentService interface {
	// List 评论列表
	List(ctx context.Context, req *v1.CommentListReq) (res *v1.CommentListRes, err error)
	// Detail 评论详情
	Detail(ctx context.Context, req *v1.CommentDetailReq) (res *v1.CommentDetailRes, err error)
	// ContentComments 内容评论列表
	ContentComments(ctx context.Context, req *v1.ContentCommentsReq) (res *v1.ContentCommentsRes, err error)
	// Create 创建评论
	Create(ctx context.Context, req *v1.CommentCreateReq) (res *v1.CommentCreateRes, err error)
	// Update 更新评论
	Update(ctx context.Context, req *v1.CommentUpdateReq) (res *v1.CommentUpdateRes, err error)
	// Delete 删除评论
	Delete(ctx context.Context, req *v1.CommentDeleteReq) (res *v1.CommentDeleteRes, err error)
	// UpdateStatus 更新评论状态
	UpdateStatus(ctx context.Context, req *v1.CommentStatusUpdateReq) (res *v1.CommentStatusUpdateRes, err error)
	// WxClientList 微信客户端评论列表
	WxClientList(ctx context.Context, req *v1.WxClientCommentListReq) (res *v1.WxClientCommentListRes, err error)
	// WxClientCreate 微信客户端创建评论
	WxClientCreate(ctx context.Context, req *v1.WxClientCommentCreateReq) (res *v1.WxClientCommentCreateRes, err error)
}

// 声明本地服务变量
var localComment CommentService

// Comment 获取评论服务
func Comment() CommentService {
	if localComment == nil {
		panic("implement not found for interface CommentService, forgot register?")
	}
	return localComment
}

// SetComment 设置评论服务实现
func SetComment(s CommentService) {
	localComment = s
}
