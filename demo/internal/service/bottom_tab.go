package service

import (
	"context"
	v1 "demo/api/content/v1"
)

// IBottomTab 底部导航栏服务接口
type IBottomTab interface {
	// GetBottomTabList 获取底部导航栏列表
	GetBottomTabList(ctx context.Context, req *v1.BottomTabListReq) (res *v1.BottomTabListRes, err error)

	// CreateBottomTab 创建底部导航项
	CreateBottomTab(ctx context.Context, req *v1.BottomTabCreateReq) (res *v1.BottomTabCreateRes, err error)

	// UpdateBottomTab 更新底部导航项
	UpdateBottomTab(ctx context.Context, req *v1.BottomTabUpdateReq) (res *v1.BottomTabUpdateRes, err error)

	// DeleteBottomTab 删除底部导航项
	DeleteBottomTab(ctx context.Context, req *v1.BottomTabDeleteReq) (res *v1.BottomTabDeleteRes, err error)

	// UpdateBottomTabStatus 更新底部导航项状态
	UpdateBottomTabStatus(ctx context.Context, req *v1.BottomTabStatusUpdateReq) (res *v1.BottomTabStatusUpdateRes, err error)

	// WxGetBottomTabList 微信客户端获取底部导航栏列表
	WxGetBottomTabList(ctx context.Context, req *v1.WxBottomTabListReq) (res *v1.WxBottomTabListRes, err error)
}

var localBottomTab IBottomTab

// BottomTab 获取底部导航栏服务
func BottomTab() IBottomTab {
	if localBottomTab == nil {
		panic("请实现 IBottomTab 接口并通过 SetBottomTab 设置服务实现")
	}
	return localBottomTab
}

// SetBottomTab 设置底部导航栏服务
func SetBottomTab(service IBottomTab) {
	localBottomTab = service
}
