package service

import (
	"context"

	v1 "demo/api/region/v1"
)

// RegionService 地区服务接口
type RegionService interface {
	// List 获取地区列表
	List(ctx context.Context, req *v1.RegionListReq) (res *v1.RegionListRes, err error)

	// Detail 获取地区详情
	Detail(ctx context.Context, req *v1.RegionDetailReq) (res *v1.RegionDetailRes, err error)

	// Create 创建地区
	Create(ctx context.Context, req *v1.RegionCreateReq) (res *v1.RegionCreateRes, err error)

	// Update 更新地区
	Update(ctx context.Context, req *v1.RegionUpdateReq) (res *v1.RegionUpdateRes, err error)

	// Delete 删除地区
	Delete(ctx context.Context, req *v1.RegionDeleteReq) (res *v1.RegionDeleteRes, err error)

	// WxClientList 客户端获取地区列表
	WxClientList(ctx context.Context, req *v1.WxClientRegionListReq) (res *v1.WxClientRegionListRes, err error)
}

var (
	localRegion RegionService
)

// Region 获取地区服务实例
func Region() RegionService {
	if localRegion == nil {
		panic("Region service not initialized")
	}
	return localRegion
}

// SetRegion 设置地区服务实例
func SetRegion(s RegionService) {
	localRegion = s
}
