package browse_history

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"

	v1 "demo/api/browse_history/v1"
	"demo/internal/dao"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// browseHistoryImpl 浏览历史记录服务实现
type browseHistoryImpl struct {
	browseHistoryDao *dao.BrowseHistoryDao
}

// New 创建一个浏览历史记录服务实例
func New() service.BrowseHistoryService {
	return &browseHistoryImpl{
		browseHistoryDao: dao.GetBrowseHistoryDao(),
	}
}

// List 获取浏览历史记录列表
func (s *browseHistoryImpl) List(ctx context.Context, req *v1.BrowseHistoryListReq) (res *v1.BrowseHistoryListRes, err error) {
	res = &v1.BrowseHistoryListRes{
		List:  make([]v1.BrowseHistoryItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 从上下文中获取客户ID
	clientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("获取用户信息失败")
	}

	// 记录请求信息
	g.Log().Debug(ctx, "获取浏览历史记录列表",
		"clientId", clientId,
		"timeType", req.TimeType,
		"page", req.Page,
		"pageSize", req.PageSize)

	// 获取浏览历史记录列表
	list, total, err := s.browseHistoryDao.GetList(ctx, clientId, req.TimeType, req.Page, req.PageSize)
	if err != nil {
		g.Log().Error(ctx, "获取浏览历史记录失败", "error", err.Error())
		return nil, err
	}

	res.Total = total

	// 如果有数据，获取内容详情
	if len(list) > 0 {
		// 遍历处理每条记录
		for _, item := range list {
			// 构建默认值
			historyItem := v1.BrowseHistoryItem{
				Id:            item.Id,
				ContentId:     item.ContentId,
				ContentType:   item.ContentType,
				ContentTitle:  "[内容已删除]",
				ContentCover:  "", // 默认空封面
				ContentStatus: 0,  // 默认状态为0
				BrowseTime:    item.BrowseTime.String(),
				Category:      "", // 默认空分类
				Price:         0,  // 默认价格为0
			}

			// 查询内容详情
			contentInfo, err := s.getContentInfo(ctx, item.ContentId, item.ContentType)
			if err == nil && contentInfo != nil {
				// 设置标题
				if title, ok := contentInfo["title"].(string); ok {
					historyItem.ContentTitle = title
				}

				// 设置分类
				if category, ok := contentInfo["category"].(string); ok {
					historyItem.Category = category
				}

				// 设置封面
				if cover, ok := contentInfo["cover"].(string); ok {
					historyItem.ContentCover = cover
				}

				// 设置价格
				if price, ok := contentInfo["price"].(float64); ok {
					historyItem.Price = price
				}

				// 处理状态 - status可能是string类型
				if status, ok := contentInfo["status"].(string); ok {
					// 尝试将状态字符串转换为数字
					switch status {
					case "已发布":
						historyItem.ContentStatus = 1
					case "待审核":
						historyItem.ContentStatus = 0
					case "已下架":
						historyItem.ContentStatus = 2
					default:
						// 尝试直接转换为数字
						if statusNum, err := strconv.Atoi(status); err == nil {
							historyItem.ContentStatus = statusNum
						}
					}
				}
			}

			// 添加到结果列表
			res.List = append(res.List, historyItem)
		}
	}

	return res, nil
}

// Clear 清空浏览历史记录
func (s *browseHistoryImpl) Clear(ctx context.Context, req *v1.BrowseHistoryClearReq) (res *v1.BrowseHistoryClearRes, err error) {
	res = &v1.BrowseHistoryClearRes{}

	// 从上下文中获取客户ID
	clientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("获取用户信息失败")
	}

	// 清空浏览历史记录
	err = s.browseHistoryDao.Clear(ctx, clientId, req.TimeType)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Add 添加浏览历史记录
func (s *browseHistoryImpl) Add(ctx context.Context, req *v1.BrowseHistoryAddReq) (res *v1.BrowseHistoryAddRes, err error) {
	res = &v1.BrowseHistoryAddRes{}

	// 从上下文中获取客户ID
	clientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("获取用户信息失败")
	}

	// 检查内容是否存在
	contentInfo, err := s.getContentInfo(ctx, req.ContentId, req.ContentType)
	if err != nil {
		return nil, gerror.New("检查内容失败")
	}
	if contentInfo == nil {
		return nil, gerror.New("内容不存在")
	}

	// 添加或更新浏览历史记录
	err = s.browseHistoryDao.AddOrUpdate(ctx, clientId, req.ContentId, req.ContentType)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// 获取内容信息
func (s *browseHistoryImpl) getContentInfo(ctx context.Context, contentId int, contentType string) (map[string]interface{}, error) {
	var tableName string

	// 根据内容类型确定表名
	switch contentType {
	case "article", "info", "idle":
		tableName = "content"
	default:
		g.Log().Warning(ctx, "不支持的内容类型", "contentType", contentType)
		return nil, gerror.New("不支持的内容类型")
	}

	// 记录查询信息
	g.Log().Debug(ctx, "获取内容信息",
		"contentId", contentId,
		"contentType", contentType,
		"tableName", tableName)

	// 定义内容结构体 - 包含需要的字段
	type ContentInfo struct {
		Id       int    `json:"id"`
		Title    string `json:"title"`
		Status   string `json:"status"`
		Content  string `json:"content"`  // 富文本内容
		Extend   string `json:"extend"`   // JSON扩展字段
		Category string `json:"category"` // 分类
	}

	// 查询内容信息，包括更多字段
	var content ContentInfo
	err := g.DB().Model(tableName).
		Fields("id, title, status, content, extend, category"). // 查询需要的字段
		Where("id", contentId).
		Scan(&content)
	if err != nil {
		g.Log().Error(ctx, "获取内容信息失败",
			"error", err.Error(),
			"contentId", contentId,
			"contentType", contentType)
		return nil, gerror.New("获取内容信息失败")
	}

	// 检查是否找到内容
	if content.Id == 0 {
		g.Log().Warning(ctx, "内容不存在或已删除",
			"contentId", contentId,
			"contentType", contentType)
		return nil, nil
	}

	// 初始化结果map
	result := map[string]interface{}{
		"id":       content.Id,
		"title":    content.Title,
		"status":   content.Status,
		"category": content.Category,
		"cover":    "", // 默认空封面
		"price":    0,  // 默认价格0
	}

	// 从内容中提取第一张图片作为封面
	if content.Content != "" {
		// 提取第一个图片标签
		re := regexp.MustCompile(`<img[^>]+src="([^"]+)"`)
		matches := re.FindStringSubmatch(content.Content)
		if len(matches) > 1 {
			result["cover"] = matches[1]
		}
	}

	// 从extend字段中提取价格信息
	if content.Extend != "" {
		var extendData map[string]interface{}
		if err := json.Unmarshal([]byte(content.Extend), &extendData); err == nil {
			// 提取价格
			if price, ok := extendData["price"].(float64); ok {
				result["price"] = price
			} else if priceStr, ok := extendData["price"].(string); ok {
				if priceVal, err := strconv.ParseFloat(priceStr, 64); err == nil {
					result["price"] = priceVal
				}
			}
		}
	}

	return result, nil
}

// Count 获取浏览历史记录数量
func (s *browseHistoryImpl) Count(ctx context.Context, req *v1.BrowseHistoryCountReq) (res *v1.BrowseHistoryCountRes, err error) {
	res = &v1.BrowseHistoryCountRes{}

	// 从上下文中获取客户ID
	clientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("获取用户信息失败")
	}

	// 记录请求信息
	g.Log().Debug(ctx, "获取浏览历史记录数量", "clientId", clientId)

	// 查询总数
	count, err := g.DB().Model("browse_history").
		Where("client_id", clientId).
		Count()
	if err != nil {
		g.Log().Error(ctx, "获取浏览历史记录数量失败", "error", err.Error())
		return nil, gerror.New("获取浏览历史记录数量失败")
	}

	res.Count = count
	return res, nil
}
