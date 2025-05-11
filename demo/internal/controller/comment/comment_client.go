package comment

import (
	"context"

	v1 "demo/api/comment/v1"
	"demo/internal/service"
)

// ClientController 微信客户端评论控制器
type ClientController struct{}

// WxClientList 获取评论列表
func (c *ClientController) WxClientList(ctx context.Context, req *v1.WxClientCommentListReq) (res *v1.WxClientCommentListRes, err error) {
	return service.Comment().WxClientList(ctx, req)
}

// WxClientCreate 创建评论
func (c *ClientController) WxClientCreate(ctx context.Context, req *v1.WxClientCommentCreateReq) (res *v1.WxClientCommentCreateRes, err error) {
	return service.Comment().WxClientCreate(ctx, req)
}
