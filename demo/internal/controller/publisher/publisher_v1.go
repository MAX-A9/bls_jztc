package publisher

import (
	"context"

	v1 "demo/api/client/v1"
	"demo/internal/service"
)

type ControllerV1 struct {
	service service.ClientService
}

// PublisherInfo 获取发布人信息
func (c *ControllerV1) PublisherInfo(ctx context.Context, req *v1.PublisherInfoReq) (res *v1.PublisherInfoRes, err error) {
	return c.service.PublisherInfo(ctx, req)
}

// FollowPublisher 关注发布人
func (c *ControllerV1) FollowPublisher(ctx context.Context, req *v1.FollowPublisherReq) (res *v1.FollowPublisherRes, err error) {
	return c.service.FollowPublisher(ctx, req)
}

// UnfollowPublisher 取消关注发布人
func (c *ControllerV1) UnfollowPublisher(ctx context.Context, req *v1.UnfollowPublisherReq) (res *v1.UnfollowPublisherRes, err error) {
	return c.service.UnfollowPublisher(ctx, req)
}

// FollowStatus 获取关注状态
func (c *ControllerV1) FollowStatus(ctx context.Context, req *v1.FollowStatusReq) (res *v1.FollowStatusRes, err error) {
	return c.service.FollowStatus(ctx, req)
}

// FollowingList 获取关注人列表
func (c *ControllerV1) FollowingList(ctx context.Context, req *v1.FollowingListReq) (res *v1.FollowingListRes, err error) {
	return c.service.FollowingList(ctx, req)
}

// FollowingCount 获取关注人总数
func (c *ControllerV1) FollowingCount(ctx context.Context, req *v1.FollowingCountReq) (res *v1.FollowingCountRes, err error) {
	return c.service.FollowingCount(ctx, req)
}
