package service

import (
	"context"
	v1 "demo/api/client/v1"
)

// ClientService 客户服务接口
type ClientService interface {
	// List 获取客户列表
	List(ctx context.Context, req *v1.ClientListReq) (res *v1.ClientListRes, err error)

	// Create 创建客户
	Create(ctx context.Context, req *v1.ClientCreateReq) (res *v1.ClientCreateRes, err error)

	// Update 更新客户
	Update(ctx context.Context, req *v1.ClientUpdateReq) (res *v1.ClientUpdateRes, err error)

	// Delete 删除客户
	Delete(ctx context.Context, req *v1.ClientDeleteReq) (res *v1.ClientDeleteRes, err error)

	// WxappLogin 微信小程序登录
	WxappLogin(ctx context.Context, req *v1.WxappLoginReq) (res *v1.WxappLoginRes, err error)

	// Info 获取客户信息
	Info(ctx context.Context, req *v1.ClientInfoReq) (res *v1.ClientInfoRes, err error)

	// UpdateProfile 更新客户个人信息
	UpdateProfile(ctx context.Context, req *v1.ClientUpdateProfileReq) (res *v1.ClientUpdateProfileRes, err error)

	// GetWxappConfig 获取微信小程序配置
	GetWxappConfig(ctx context.Context, req *v1.WxappConfigGetReq) (res *v1.WxappConfigGetRes, err error)

	// SaveWxappConfig 保存微信小程序配置
	SaveWxappConfig(ctx context.Context, req *v1.WxappConfigSaveReq) (res *v1.WxappConfigSaveRes, err error)

	// PublisherInfo 获取发布人信息
	PublisherInfo(ctx context.Context, req *v1.PublisherInfoReq) (res *v1.PublisherInfoRes, err error)

	// FollowPublisher 关注发布人
	FollowPublisher(ctx context.Context, req *v1.FollowPublisherReq) (res *v1.FollowPublisherRes, err error)

	// UnfollowPublisher 取消关注发布人
	UnfollowPublisher(ctx context.Context, req *v1.UnfollowPublisherReq) (res *v1.UnfollowPublisherRes, err error)

	// FollowStatus 获取关注状态
	FollowStatus(ctx context.Context, req *v1.FollowStatusReq) (res *v1.FollowStatusRes, err error)

	// FollowingList 获取关注人列表
	FollowingList(ctx context.Context, req *v1.FollowingListReq) (res *v1.FollowingListRes, err error)

	// FollowingCount 获取关注人总数
	FollowingCount(ctx context.Context, req *v1.FollowingCountReq) (res *v1.FollowingCountRes, err error)

	// DurationList 获取客户时长列表
	DurationList(ctx context.Context, req *v1.ClientDurationListReq) (res *v1.ClientDurationListRes, err error)

	// DurationDetail 获取客户时长详情
	DurationDetail(ctx context.Context, req *v1.ClientDurationDetailReq) (res *v1.ClientDurationDetailRes, err error)

	// DurationCreate 创建客户时长
	DurationCreate(ctx context.Context, req *v1.ClientDurationCreateReq) (res *v1.ClientDurationCreateRes, err error)

	// DurationUpdate 更新客户时长
	DurationUpdate(ctx context.Context, req *v1.ClientDurationUpdateReq) (res *v1.ClientDurationUpdateRes, err error)

	// DurationDelete 删除客户时长
	DurationDelete(ctx context.Context, req *v1.ClientDurationDeleteReq) (res *v1.ClientDurationDeleteRes, err error)

	// WxGetRemainingDuration 获取客户端用户剩余时长
	WxGetRemainingDuration(ctx context.Context, req *v1.WxClientRemainingDurationReq) (res *v1.WxClientRemainingDurationRes, err error)
}

var (
	clientService ClientService
)

// SetClient 设置客户服务
func SetClient(s ClientService) {
	clientService = s
}

// Client 获取客户服务
func Client() ClientService {
	return clientService
}
