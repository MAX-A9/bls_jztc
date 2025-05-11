package content

import (
	"context"

	v1 "demo/api/content/v1"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gerror"
)

// Controller 内容管理控制器
type Controller struct{}

// V1 创建V1版本控制器
func (c *Controller) V1() *ControllerV1 {
	return &ControllerV1{}
}

// ControllerV1 V1版本内容管理控制器
type ControllerV1 struct{}

// 内容管理相关接口
// ================================================================================

// List 获取内容列表
func (c *ControllerV1) List(ctx context.Context, req *v1.ContentListReq) (res *v1.ContentListRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看内容列表")
	}

	return service.Content().List(ctx, req)
}

// Detail 获取内容详情
func (c *ControllerV1) Detail(ctx context.Context, req *v1.ContentDetailReq) (res *v1.ContentDetailRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看内容详情")
	}

	return service.Content().Detail(ctx, req)
}

// Create 创建内容
func (c *ControllerV1) Create(ctx context.Context, req *v1.ContentCreateReq) (res *v1.ContentCreateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限创建内容")
	}

	return service.Content().Create(ctx, req)
}

// Update 更新内容
func (c *ControllerV1) Update(ctx context.Context, req *v1.ContentUpdateReq) (res *v1.ContentUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新内容")
	}

	return service.Content().Update(ctx, req)
}

// Delete 删除内容
func (c *ControllerV1) Delete(ctx context.Context, req *v1.ContentDeleteReq) (res *v1.ContentDeleteRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限删除内容")
	}

	return service.Content().Delete(ctx, req)
}

// UpdateStatus 更新内容状态
func (c *ControllerV1) UpdateStatus(ctx context.Context, req *v1.ContentStatusUpdateReq) (res *v1.ContentStatusUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新内容状态")
	}

	return service.Content().UpdateStatus(ctx, req)
}

// UpdateRecommend 更新内容推荐状态
func (c *ControllerV1) UpdateRecommend(ctx context.Context, req *v1.ContentRecommendUpdateReq) (res *v1.ContentRecommendUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新内容推荐状态")
	}

	return service.Content().UpdateRecommend(ctx, req)
}

// 首页分类相关接口
// ================================================================================

// HomeList 获取首页分类列表
func (c *ControllerV1) HomeList(ctx context.Context, req *v1.HomeCategoryListReq) (res *v1.HomeCategoryListRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看首页分类")
	}

	return service.Category().HomeList(ctx, req)
}

// HomeCreate 创建首页分类
func (c *ControllerV1) HomeCreate(ctx context.Context, req *v1.HomeCategoryCreateReq) (res *v1.HomeCategoryCreateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限创建首页分类")
	}

	return service.Category().HomeCreate(ctx, req)
}

// HomeUpdate 更新首页分类
func (c *ControllerV1) HomeUpdate(ctx context.Context, req *v1.HomeCategoryUpdateReq) (res *v1.HomeCategoryUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新首页分类")
	}

	return service.Category().HomeUpdate(ctx, req)
}

// HomeDelete 删除首页分类
func (c *ControllerV1) HomeDelete(ctx context.Context, req *v1.HomeCategoryDeleteReq) (res *v1.HomeCategoryDeleteRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限删除首页分类")
	}

	return service.Category().HomeDelete(ctx, req)
}

// 闲置分类相关接口
// ================================================================================

// IdleList 获取闲置分类列表
func (c *ControllerV1) IdleList(ctx context.Context, req *v1.IdleCategoryListReq) (res *v1.IdleCategoryListRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看闲置分类")
	}

	return service.Category().IdleList(ctx, req)
}

// IdleCreate 创建闲置分类
func (c *ControllerV1) IdleCreate(ctx context.Context, req *v1.IdleCategoryCreateReq) (res *v1.IdleCategoryCreateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限创建闲置分类")
	}

	return service.Category().IdleCreate(ctx, req)
}

// IdleUpdate 更新闲置分类
func (c *ControllerV1) IdleUpdate(ctx context.Context, req *v1.IdleCategoryUpdateReq) (res *v1.IdleCategoryUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新闲置分类")
	}

	return service.Category().IdleUpdate(ctx, req)
}

// IdleDelete 删除闲置分类
func (c *ControllerV1) IdleDelete(ctx context.Context, req *v1.IdleCategoryDeleteReq) (res *v1.IdleCategoryDeleteRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限删除闲置分类")
	}

	return service.Category().IdleDelete(ctx, req)
}

// GetCategories 获取所有分类
func (c *ControllerV1) GetCategories(ctx context.Context, req *v1.GetCategoriesReq) (res *v1.GetCategoriesRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限获取分类")
	}

	return service.Content().GetCategories(ctx, req)
}

// New 创建内容管理控制器
func New() *Controller {
	return &Controller{}
}
