package comment

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "demo/api/comment/v1"
	"demo/internal/service"
)

// Controller 评论控制器
type Controller struct{}

// List 获取评论列表
func (c *Controller) List(ctx context.Context, req *v1.CommentListReq) (res *v1.CommentListRes, err error) {
	return service.Comment().List(ctx, req)
}

// Detail 获取评论详情
func (c *Controller) Detail(ctx context.Context, req *v1.CommentDetailReq) (res *v1.CommentDetailRes, err error) {
	return service.Comment().Detail(ctx, req)
}

// ContentComments 获取指定内容的评论列表
func (c *Controller) ContentComments(ctx context.Context, req *v1.ContentCommentsReq) (res *v1.ContentCommentsRes, err error) {
	return service.Comment().ContentComments(ctx, req)
}

// Create 创建评论
func (c *Controller) Create(ctx context.Context, req *v1.CommentCreateReq) (res *v1.CommentCreateRes, err error) {
	// 检查必要参数
	if req.ContentId <= 0 {
		return nil, gerror.New("内容ID必须大于0")
	}
	if req.ClientId <= 0 {
		return nil, gerror.New("客户ID必须大于0")
	}
	if req.RealName == "" {
		return nil, gerror.New("真实姓名不能为空")
	}
	if req.Comment == "" {
		return nil, gerror.New("评论内容不能为空")
	}

	return service.Comment().Create(ctx, req)
}

// Update 更新评论
func (c *Controller) Update(ctx context.Context, req *v1.CommentUpdateReq) (res *v1.CommentUpdateRes, err error) {
	// 检查必要参数
	if req.Id <= 0 {
		return nil, gerror.New("评论ID必须大于0")
	}
	if req.RealName == "" {
		return nil, gerror.New("真实姓名不能为空")
	}
	if req.Comment == "" {
		return nil, gerror.New("评论内容不能为空")
	}

	return service.Comment().Update(ctx, req)
}

// Delete 删除评论
func (c *Controller) Delete(ctx context.Context, req *v1.CommentDeleteReq) (res *v1.CommentDeleteRes, err error) {
	// 检查必要参数
	if req.Id <= 0 {
		return nil, gerror.New("评论ID必须大于0")
	}

	return service.Comment().Delete(ctx, req)
}

// UpdateStatus 更新评论状态
func (c *Controller) UpdateStatus(ctx context.Context, req *v1.CommentStatusUpdateReq) (res *v1.CommentStatusUpdateRes, err error) {
	// 检查必要参数
	if req.Id <= 0 {
		return nil, gerror.New("评论ID必须大于0")
	}
	if req.Status == "" {
		return nil, gerror.New("状态不能为空")
	}

	return service.Comment().UpdateStatus(ctx, req)
}
