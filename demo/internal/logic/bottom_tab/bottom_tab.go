package bottom_tab

import (
	"context"
	v1 "demo/api/content/v1"
	"demo/internal/dao"
	"demo/internal/model/entity"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gerror"
)

// sBottomTab 底部导航栏服务实现
type sBottomTab struct {
	bottomTabDao *dao.BottomTabDao
}

// New 创建底部导航栏服务
func New() service.IBottomTab {
	return &sBottomTab{
		bottomTabDao: &dao.BottomTabDao{},
	}
}

// GetBottomTabList 获取底部导航栏列表
func (s *sBottomTab) GetBottomTabList(ctx context.Context, req *v1.BottomTabListReq) (res *v1.BottomTabListRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看底部导航栏列表")
	}

	// 初始化响应对象
	res = &v1.BottomTabListRes{
		List: make([]v1.BottomTabItem, 0),
	}

	// 获取底部导航栏列表
	tabs, err := s.bottomTabDao.GetList(ctx, false)
	if err != nil {
		return nil, gerror.New("获取底部导航栏列表失败: " + err.Error())
	}

	// 转换为API响应格式
	for _, tab := range tabs {
		res.List = append(res.List, v1.BottomTabItem{
			Id:           tab.Id,
			Name:         tab.Name,
			Icon:         tab.Icon,
			SelectedIcon: tab.SelectedIcon,
			Path:         tab.Path,
			Order:        tab.Order,
			IsEnabled:    tab.IsEnabled == 1,
		})
	}

	return res, nil
}

// CreateBottomTab 创建底部导航项
func (s *sBottomTab) CreateBottomTab(ctx context.Context, req *v1.BottomTabCreateReq) (res *v1.BottomTabCreateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限创建底部导航项")
	}

	// 创建功能已禁用
	return nil, gerror.New("创建底部导航项功能已禁用，请使用已有的导航项")

	/*
		// 以下代码暂时不启用
		// 检查底部导航项数量，限制最多不超过5个
		tabs, err := s.bottomTabDao.GetList(ctx, false)
		if err != nil {
			return nil, gerror.New("获取底部导航栏列表失败: " + err.Error())
		}
		if len(tabs) >= 5 {
			return nil, gerror.New("底部导航项最多为5个，无法继续添加")
		}

		// 准备保存的数据
		tab := &entity.BottomTab{
			Name:         req.Name,
			Icon:         req.Icon,
			SelectedIcon: req.SelectedIcon,
			Path:         req.Path,
			Order:        req.Order,
			IsEnabled:    0,
		}

		if req.IsEnabled {
			tab.IsEnabled = 1
		}

		// 创建底部导航项
		id, err := s.bottomTabDao.Create(ctx, tab)
		if err != nil {
			return nil, gerror.New("创建底部导航项失败: " + err.Error())
		}

		return &v1.BottomTabCreateRes{
			Id: int(id),
		}, nil
	*/
}

// UpdateBottomTab 更新底部导航项
func (s *sBottomTab) UpdateBottomTab(ctx context.Context, req *v1.BottomTabUpdateReq) (res *v1.BottomTabUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新底部导航项")
	}

	// 检查底部导航项是否存在
	_, err = s.bottomTabDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("底部导航项不存在或已被删除")
	}

	// 准备更新的数据
	tab := &entity.BottomTab{
		Id:           req.Id,
		Name:         req.Name,
		Icon:         req.Icon,
		SelectedIcon: req.SelectedIcon,
		Path:         req.Path,
		Order:        req.Order,
		IsEnabled:    0,
	}

	if req.IsEnabled {
		tab.IsEnabled = 1
	}

	// 更新底部导航项
	if err := s.bottomTabDao.Update(ctx, tab); err != nil {
		return nil, gerror.New("更新底部导航项失败: " + err.Error())
	}

	return &v1.BottomTabUpdateRes{}, nil
}

// DeleteBottomTab 删除底部导航项
func (s *sBottomTab) DeleteBottomTab(ctx context.Context, req *v1.BottomTabDeleteReq) (res *v1.BottomTabDeleteRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限删除底部导航项")
	}

	// 删除功能已禁用
	return nil, gerror.New("删除底部导航项功能已禁用，请使用状态开关来控制导航项的显示")

	/*
		// 以下代码暂时不启用
		// 检查底部导航项是否存在
		_, err = s.bottomTabDao.GetById(ctx, req.Id)
		if err != nil {
			return nil, gerror.New("底部导航项不存在或已被删除")
		}

		// 删除底部导航项
		if err := s.bottomTabDao.Delete(ctx, req.Id); err != nil {
			return nil, gerror.New("删除底部导航项失败: " + err.Error())
		}

		return &v1.BottomTabDeleteRes{}, nil
	*/
}

// UpdateBottomTabStatus 更新底部导航项状态
func (s *sBottomTab) UpdateBottomTabStatus(ctx context.Context, req *v1.BottomTabStatusUpdateReq) (res *v1.BottomTabStatusUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新底部导航项状态")
	}

	// 检查底部导航项是否存在
	_, err = s.bottomTabDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("底部导航项不存在或已被删除")
	}

	// 更新状态
	isEnabled := 0
	if req.IsEnabled {
		isEnabled = 1
	}

	if err := s.bottomTabDao.UpdateStatus(ctx, req.Id, isEnabled); err != nil {
		return nil, gerror.New("更新底部导航项状态失败: " + err.Error())
	}

	return &v1.BottomTabStatusUpdateRes{}, nil
}

// WxGetBottomTabList 微信客户端获取底部导航栏列表
func (s *sBottomTab) WxGetBottomTabList(ctx context.Context, req *v1.WxBottomTabListReq) (res *v1.WxBottomTabListRes, err error) {
	// 初始化响应对象
	res = &v1.WxBottomTabListRes{
		List: make([]v1.BottomTabItem, 0),
	}

	// 获取底部导航栏列表，只获取已启用的项目
	tabs, err := s.bottomTabDao.GetList(ctx, true)
	if err != nil {
		return nil, gerror.New("获取底部导航栏列表失败: " + err.Error())
	}

	// 转换为API响应格式
	for _, tab := range tabs {
		res.List = append(res.List, v1.BottomTabItem{
			Id:           tab.Id,
			Name:         tab.Name,
			Icon:         tab.Icon,
			SelectedIcon: tab.SelectedIcon,
			Path:         tab.Path,
			Order:        tab.Order,
			IsEnabled:    tab.IsEnabled == 1,
		})
	}

	return res, nil
}
