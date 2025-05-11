package favorite

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"

	v1 "demo/api/content/client/v1"
	"demo/internal/dao"
	"demo/internal/service"
	"demo/utility/auth"
)

// sFavorite 收藏服务实现
type sFavorite struct{}

// New 创建收藏服务实例
func New() service.FavoriteService {
	return &sFavorite{}
}

// Add 添加收藏
func (s *sFavorite) Add(ctx context.Context, req *v1.FavoriteAddReq) (res *v1.FavoriteAddRes, err error) {
	res = &v1.FavoriteAddRes{
		Success: false,
		Message: "",
	}

	// 获取当前客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("获取客户信息失败: " + err.Error())
	}

	// 检查内容是否存在
	contentDao := dao.NewContentDao()
	content, err := contentDao.FindOne(ctx, req.ContentId)
	if err != nil {
		return nil, gerror.New("获取内容失败: " + err.Error())
	}
	if content == nil {
		return nil, gerror.New("内容不存在")
	}

	// 检查内容状态是否为已发布
	if gconv.String(content.Status) != "已发布" {
		return nil, gerror.New("内容不可收藏，状态不是已发布")
	}

	// 添加收藏
	favoriteDao := dao.NewFavoriteDao()
	err = favoriteDao.Add(ctx, clientId, req.ContentId)
	if err != nil {
		return nil, gerror.New("添加收藏失败: " + err.Error())
	}

	// 设置响应
	res.Success = true
	res.Message = "收藏成功"

	return res, nil
}

// Cancel 取消收藏
func (s *sFavorite) Cancel(ctx context.Context, req *v1.FavoriteCancelReq) (res *v1.FavoriteCancelRes, err error) {
	res = &v1.FavoriteCancelRes{
		Success: false,
		Message: "",
	}

	// 获取当前客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("获取客户信息失败: " + err.Error())
	}

	// 检查内容是否存在
	contentDao := dao.NewContentDao()
	content, err := contentDao.FindOne(ctx, req.ContentId)
	if err != nil {
		return nil, gerror.New("获取内容失败: " + err.Error())
	}
	if content == nil {
		return nil, gerror.New("内容不存在")
	}

	// 取消收藏
	favoriteDao := dao.NewFavoriteDao()
	err = favoriteDao.Cancel(ctx, clientId, req.ContentId)
	if err != nil {
		return nil, gerror.New("取消收藏失败: " + err.Error())
	}

	// 设置响应
	res.Success = true
	res.Message = "取消收藏成功"

	return res, nil
}

// GetStatus 获取收藏状态
func (s *sFavorite) GetStatus(ctx context.Context, req *v1.FavoriteStatusReq) (res *v1.FavoriteStatusRes, err error) {
	res = &v1.FavoriteStatusRes{
		IsFavorite: false,
	}

	// 获取当前客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("获取客户信息失败: " + err.Error())
	}

	// 检查内容是否存在
	contentDao := dao.NewContentDao()
	content, err := contentDao.FindOne(ctx, req.ContentId)
	if err != nil {
		return nil, gerror.New("获取内容失败: " + err.Error())
	}
	if content == nil {
		return nil, gerror.New("内容不存在")
	}

	// 获取收藏状态
	favoriteDao := dao.NewFavoriteDao()
	isFavorite, err := favoriteDao.IsFavorite(ctx, clientId, req.ContentId)
	if err != nil {
		return nil, gerror.New("获取收藏状态失败: " + err.Error())
	}

	// 设置响应
	res.IsFavorite = isFavorite

	return res, nil
}

// GetList 获取收藏列表
func (s *sFavorite) GetList(ctx context.Context, req *v1.FavoriteListReq) (res *v1.FavoriteListRes, err error) {
	res = &v1.FavoriteListRes{
		List:  make([]v1.FavoriteItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 获取当前客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("获取客户信息失败: " + err.Error())
	}

	// 获取收藏列表
	favoriteDao := dao.NewFavoriteDao()
	list, total, err := favoriteDao.GetFavoriteList(ctx, clientId, req.Page, req.PageSize, req.Type, req.Category)
	if err != nil {
		return nil, gerror.New("获取收藏列表失败: " + err.Error())
	}

	// 设置响应
	res.Total = total
	res.Page = req.Page

	// 转换数据
	for _, item := range list {
		// 创建列表项
		listItem := v1.FavoriteItem{
			Id:       gconv.Int(item["id"]),
			Title:    gconv.String(item["title"]),
			Category: gconv.String(item["category"]),
			// 使用Author字段作为Publisher
			Publisher: gconv.String(item["author"]),
		}

		// 处理图片
		// 尝试从内容中提取第一张图片
		coverImage := ""
		contentStr := gconv.String(item["content"])

		// 简单判断内容是否包含JSON格式的图片数组
		if len(contentStr) > 0 && strings.Contains(contentStr, "[") && strings.Contains(contentStr, "]") {
			var images []string
			startIdx := strings.Index(contentStr, "[")
			endIdx := strings.LastIndex(contentStr, "]")
			if startIdx >= 0 && endIdx > startIdx {
				jsonStr := contentStr[startIdx : endIdx+1]
				if err := json.Unmarshal([]byte(jsonStr), &images); err == nil && len(images) > 0 {
					coverImage = images[0]
				}
			}
		}

		listItem.Image = coverImage

		// 解析扩展字段，提取内容类型和价格
		if item["extend"] != nil && gconv.String(item["extend"]) != "" {
			var extendMap map[string]interface{}
			if err := json.Unmarshal([]byte(gconv.String(item["extend"])), &extendMap); err == nil {
				// 提取内容类型
				if contentType, ok := extendMap["type"]; ok {
					listItem.Type = gconv.Int(contentType)
				} else {
					// 默认为普通信息(首页分类)
					listItem.Type = 1
				}

				// 如果是闲置物品，提取价格
				if listItem.Type == 2 {
					// 提取价格
					if price, ok := extendMap["price"]; ok {
						listItem.Price = gconv.Float64(price)
					}
				}
			}
		} else {
			// 如果没有extend字段或为空，默认为普通信息(首页分类)
			listItem.Type = 1
		}

		res.List = append(res.List, listItem)
	}

	return res, nil
}

// GetCount 获取收藏总数
func (s *sFavorite) GetCount(ctx context.Context, req *v1.FavoriteCountReq) (res *v1.FavoriteCountRes, err error) {
	res = &v1.FavoriteCountRes{
		Total: 0,
	}

	// 获取当前客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("获取客户信息失败: " + err.Error())
	}

	// 获取收藏总数
	favoriteDao := dao.NewFavoriteDao()
	total, err := favoriteDao.GetFavoriteCount(ctx, clientId)
	if err != nil {
		return nil, gerror.New("获取收藏总数失败: " + err.Error())
	}

	// 设置响应
	res.Total = total

	return res, nil
}
