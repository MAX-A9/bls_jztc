package mall

import (
	"context"
	v1 "demo/api/mall/v1"
	"demo/internal/dao"
	"demo/internal/model/entity"
	"demo/internal/service"
	"demo/utility/auth"
	"math"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// 确保 sMall 实现了 MallService 接口
var _ service.MallService = (*sMall)(nil)

func init() {
	service.RegisterMall(New())

	// 初始同步所有商品分类的商品数量
	go func() {
		// 延迟5秒等待系统完全启动
		time.Sleep(5 * time.Second)
		ctx := context.Background()
		shopCategoryDao := &dao.ShopCategoryDao{}
		err := shopCategoryDao.UpdateAllCategoriesProductCount(ctx)
		if err != nil {
			g.Log().Error(ctx, "初始同步商品分类数量失败:", err)
		} else {
			g.Log().Info(ctx, "初始同步商品分类数量成功")
		}
	}()
}

// sMall 商城服务实现
type sMall struct {
	shopCategoryDao *dao.ShopCategoryDao
}

// New 创建商城服务实例
func New() *sMall {
	return &sMall{
		shopCategoryDao: &dao.ShopCategoryDao{},
	}
}

// GetCategory 获取商城分类详情
func (s *sMall) GetCategory(ctx context.Context, req *v1.ShopCategoryGetReq) (res *v1.ShopCategoryGetRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看此商城分类")
	}

	// 获取商城分类详情
	category, err := s.shopCategoryDao.Get(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("获取商城分类失败: " + err.Error())
	}
	if category == nil {
		return nil, gerror.New("商城分类不存在")
	}

	// 构建响应
	res = &v1.ShopCategoryGetRes{
		Id:           category.Id,
		Name:         category.Name,
		SortOrder:    category.SortOrder,
		ProductCount: category.ProductCount,
		Status:       category.Status,
		Image:        category.Image,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}

	return res, nil
}

// CreateCategory 创建商城分类
func (s *sMall) CreateCategory(ctx context.Context, req *v1.ShopCategoryCreateReq) (res *v1.ShopCategoryCreateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限创建商城分类")
	}

	// 创建商城分类
	category := &entity.ShopCategory{
		Name:         req.Name,
		SortOrder:    req.SortOrder,
		ProductCount: req.ProductCount,
		Status:       req.Status,
		Image:        req.Image,
	}

	id, err := s.shopCategoryDao.Create(ctx, category)
	if err != nil {
		return nil, gerror.New("创建商城分类失败: " + err.Error())
	}

	// 构建响应
	res = &v1.ShopCategoryCreateRes{
		Id: int(id),
	}

	return res, nil
}

// UpdateCategory 更新商城分类
func (s *sMall) UpdateCategory(ctx context.Context, req *v1.ShopCategoryUpdateReq) (res *v1.ShopCategoryUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新商城分类")
	}

	// 检查分类是否存在
	_, err = s.shopCategoryDao.Get(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("商城分类不存在")
	}

	// 更新商城分类
	category := &entity.ShopCategory{
		Id:           req.Id,
		Name:         req.Name,
		SortOrder:    req.SortOrder,
		ProductCount: req.ProductCount,
		Status:       req.Status,
		Image:        req.Image,
	}

	err = s.shopCategoryDao.Update(ctx, category)
	if err != nil {
		return nil, gerror.New("更新商城分类失败: " + err.Error())
	}

	// 构建响应
	res = &v1.ShopCategoryUpdateRes{}

	return res, nil
}

// DeleteCategory 删除商城分类
func (s *sMall) DeleteCategory(ctx context.Context, req *v1.ShopCategoryDeleteReq) (res *v1.ShopCategoryDeleteRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限删除商城分类")
	}

	// 检查分类是否存在
	_, err = s.shopCategoryDao.Get(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("商城分类不存在")
	}

	// 检查分类下是否有商品
	count, err := s.shopCategoryDao.GetProductCount(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("检查分类商品失败: " + err.Error())
	}
	if count > 0 {
		return nil, gerror.New("该分类下还有商品，不能删除")
	}

	// 删除商城分类
	err = s.shopCategoryDao.Delete(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("删除商城分类失败: " + err.Error())
	}

	// 构建响应
	res = &v1.ShopCategoryDeleteRes{}

	return res, nil
}

// UpdateCategoryStatus 更新商城分类状态
func (s *sMall) UpdateCategoryStatus(ctx context.Context, req *v1.ShopCategoryStatusUpdateReq) (res *v1.ShopCategoryStatusUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新商城分类状态")
	}

	// 检查分类是否存在
	_, err = s.shopCategoryDao.Get(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("商城分类不存在")
	}

	// 更新商城分类状态
	err = s.shopCategoryDao.UpdateStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, gerror.New("更新商城分类状态失败: " + err.Error())
	}

	// 构建响应
	res = &v1.ShopCategoryStatusUpdateRes{}

	return res, nil
}

// GetCategoryList 获取商城分类列表
func (s *sMall) GetCategoryList(ctx context.Context, req *v1.ShopCategoryListReq) (res *v1.ShopCategoryListRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看商城分类列表")
	}

	// 初始化响应
	res = &v1.ShopCategoryListRes{
		List: make([]v1.ShopCategoryItem, 0),
	}

	// 构建过滤条件
	filter := &dao.ShopCategoryFilter{
		Name:      req.Name,
		Status:    -1, // 默认获取所有状态
		SortField: req.SortField,
		SortOrder: req.SortOrder,
	}

	// 如果指定了有效状态值，使用指定值
	if req.Status == 0 || req.Status == 1 {
		filter.Status = req.Status
	}

	// 判断是否需要分页
	if req.Page > 0 && req.Size > 0 {
		// 使用分页查询
		categories, total, err := s.shopCategoryDao.GetPage(ctx, req.Page, req.Size, filter)
		if err != nil {
			return nil, gerror.New("获取商城分类列表失败: " + err.Error())
		}

		// 计算总页数
		pages := 0
		if total > 0 {
			pages = int(math.Ceil(float64(total) / float64(req.Size)))
		}

		// 填充分页信息
		res.Total = total
		res.Page = req.Page
		res.Size = req.Size
		res.Pages = pages

		// 转换分类列表
		res.List = make([]v1.ShopCategoryItem, len(categories))
		for i, category := range categories {
			res.List[i] = v1.ShopCategoryItem{
				Id:           category.Id,
				Name:         category.Name,
				SortOrder:    category.SortOrder,
				ProductCount: category.ProductCount,
				Status:       category.Status,
				Image:        category.Image,
				CreatedAt:    category.CreatedAt,
			}
		}
	} else {
		// 不分页，获取全部列表
		categories, err := s.shopCategoryDao.GetList(ctx, filter)
		if err != nil {
			return nil, gerror.New("获取商城分类列表失败: " + err.Error())
		}

		// 转换分类列表
		res.List = make([]v1.ShopCategoryItem, len(categories))
		for i, category := range categories {
			res.List[i] = v1.ShopCategoryItem{
				Id:           category.Id,
				Name:         category.Name,
				SortOrder:    category.SortOrder,
				ProductCount: category.ProductCount,
				Status:       category.Status,
				Image:        category.Image,
				CreatedAt:    category.CreatedAt,
			}
		}
	}

	return res, nil
}

// WxGetCategoryList 微信客户端获取商城分类列表
func (s *sMall) WxGetCategoryList(ctx context.Context, req *v1.WxShopCategoryListReq) (res *v1.WxShopCategoryListRes, err error) {
	// 初始化响应
	res = &v1.WxShopCategoryListRes{
		List: make([]v1.ShopCategoryItem, 0),
	}

	// 创建过滤条件 (微信端只获取启用的分类)
	filter := &dao.ShopCategoryFilter{
		Status:    1, // 固定为启用状态
		SortField: req.SortField,
		SortOrder: req.SortOrder,
	}

	// 判断是否需要分页
	if req.Page > 0 && req.Size > 0 {
		// 使用分页查询
		categories, total, err := s.shopCategoryDao.GetPage(ctx, req.Page, req.Size, filter)
		if err != nil {
			return nil, gerror.New("获取商城分类列表失败: " + err.Error())
		}

		// 计算总页数
		pages := 0
		if total > 0 {
			pages = int(math.Ceil(float64(total) / float64(req.Size)))
		}

		// 填充分页信息
		res.Total = total
		res.Page = req.Page
		res.Size = req.Size
		res.Pages = pages

		// 转换分类列表
		res.List = make([]v1.ShopCategoryItem, len(categories))
		for i, category := range categories {
			res.List[i] = v1.ShopCategoryItem{
				Id:           category.Id,
				Name:         category.Name,
				SortOrder:    category.SortOrder,
				ProductCount: category.ProductCount,
				Status:       category.Status,
				Image:        category.Image,
				CreatedAt:    category.CreatedAt,
			}
		}
	} else {
		// 不分页，获取全部列表
		categories, err := s.shopCategoryDao.GetList(ctx, filter)
		if err != nil {
			return nil, gerror.New("获取商城分类列表失败: " + err.Error())
		}

		// 转换分类列表
		res.List = make([]v1.ShopCategoryItem, len(categories))
		for i, category := range categories {
			res.List[i] = v1.ShopCategoryItem{
				Id:           category.Id,
				Name:         category.Name,
				SortOrder:    category.SortOrder,
				ProductCount: category.ProductCount,
				Status:       category.Status,
				Image:        category.Image,
				CreatedAt:    category.CreatedAt,
			}
		}
	}

	return res, nil
}

// SyncCategoriesProductCount 同步商品分类数量
func (s *sMall) SyncCategoriesProductCount(ctx context.Context, req *v1.ShopCategorySyncProductCountReq) (res *v1.ShopCategorySyncProductCountRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限执行此操作")
	}

	// 执行同步操作
	err = s.shopCategoryDao.UpdateAllCategoriesProductCount(ctx)
	if err != nil {
		return nil, gerror.New("同步商品分类数量失败: " + err.Error())
	}

	// 构建响应
	res = &v1.ShopCategorySyncProductCountRes{
		Status: true,
	}

	return res, nil
}
