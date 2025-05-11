package service

import (
	"context"

	v1 "demo/api/content/v1"
)

// ContentService 内容管理服务接口
type ContentService interface {
	// List 内容列表
	List(ctx context.Context, req *v1.ContentListReq) (res *v1.ContentListRes, err error)
	// Detail 内容详情
	Detail(ctx context.Context, req *v1.ContentDetailReq) (res *v1.ContentDetailRes, err error)
	// Create 创建内容
	Create(ctx context.Context, req *v1.ContentCreateReq) (res *v1.ContentCreateRes, err error)
	// Update 更新内容
	Update(ctx context.Context, req *v1.ContentUpdateReq) (res *v1.ContentUpdateRes, err error)
	// Delete 删除内容
	Delete(ctx context.Context, req *v1.ContentDeleteReq) (res *v1.ContentDeleteRes, err error)
	// UpdateStatus 更新内容状态
	UpdateStatus(ctx context.Context, req *v1.ContentStatusUpdateReq) (res *v1.ContentStatusUpdateRes, err error)
	// UpdateRecommend 更新内容推荐状态
	UpdateRecommend(ctx context.Context, req *v1.ContentRecommendUpdateReq) (res *v1.ContentRecommendUpdateRes, err error)
	// GetCategories 获取所有分类
	GetCategories(ctx context.Context, req *v1.GetCategoriesReq) (res *v1.GetCategoriesRes, err error)
}

// CategoryService 分类管理服务接口
type CategoryService interface {
	// HomeList 获取首页分类列表
	HomeList(ctx context.Context, req *v1.HomeCategoryListReq) (res *v1.HomeCategoryListRes, err error)
	// HomeCreate 创建首页分类
	HomeCreate(ctx context.Context, req *v1.HomeCategoryCreateReq) (res *v1.HomeCategoryCreateRes, err error)
	// HomeUpdate 更新首页分类
	HomeUpdate(ctx context.Context, req *v1.HomeCategoryUpdateReq) (res *v1.HomeCategoryUpdateRes, err error)
	// HomeDelete 删除首页分类
	HomeDelete(ctx context.Context, req *v1.HomeCategoryDeleteReq) (res *v1.HomeCategoryDeleteRes, err error)
	// IdleList 获取闲置分类列表
	IdleList(ctx context.Context, req *v1.IdleCategoryListReq) (res *v1.IdleCategoryListRes, err error)
	// IdleCreate 创建闲置分类
	IdleCreate(ctx context.Context, req *v1.IdleCategoryCreateReq) (res *v1.IdleCategoryCreateRes, err error)
	// IdleUpdate 更新闲置分类
	IdleUpdate(ctx context.Context, req *v1.IdleCategoryUpdateReq) (res *v1.IdleCategoryUpdateRes, err error)
	// IdleDelete 删除闲置分类
	IdleDelete(ctx context.Context, req *v1.IdleCategoryDeleteReq) (res *v1.IdleCategoryDeleteRes, err error)
}

// 声明服务变量
var (
	localContent  ContentService
	localCategory CategoryService
)

// Content 获取内容服务
func Content() ContentService {
	if localContent == nil {
		panic("implement not found for interface ContentService, forgot register?")
	}
	return localContent
}

// Category 获取分类服务
func Category() CategoryService {
	if localCategory == nil {
		panic("implement not found for interface CategoryService, forgot register?")
	}
	return localCategory
}

// SetContent 设置内容服务实现
func SetContent(s ContentService) {
	localContent = s
}

// SetCategory 设置分类服务实现
func SetCategory(s CategoryService) {
	localCategory = s
}
