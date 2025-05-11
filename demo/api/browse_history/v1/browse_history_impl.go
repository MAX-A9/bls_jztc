package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"demo/utility/auth"
)

// ControllerImpl 控制器实现
type ControllerImpl struct{}

// 获取浏览历史记录列表
func (c *ControllerImpl) List(ctx context.Context, req *BrowseHistoryListReq) (res *BrowseHistoryListRes, err error) {
	res = &BrowseHistoryListRes{
		List:  make([]BrowseHistoryItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 从上下文中获取客户ID
	clientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("获取用户信息失败")
	}

	// 构建查询条件
	model := g.DB().Model("browse_history").
		Where("client_id", clientId).
		Order("browse_time DESC")

	// 根据时间类型添加过滤条件
	model = c.applyTimeFilter(ctx, model, req.TimeType)

	// 查询总数
	total, err := model.Count()
	if err != nil {
		return nil, gerror.New("获取浏览历史记录总数失败")
	}
	res.Total = total

	// 分页查询数据
	var records []struct {
		Id          int         `json:"id"`
		ContentId   int         `json:"content_id"`
		ContentType string      `json:"content_type"`
		BrowseTime  *gtime.Time `json:"browse_time"`
		// 用于连接内容表字段
		Title    string `json:"title"`
		Content  string `json:"content"` // 添加content字段用于提取封面
		Status   int    `json:"status"`
		Category string `json:"category"`
		Extend   string `json:"extend"` // JSON扩展字段，可能包含价格
	}

	// 修改关联查询，获取分类和扩展字段
	err = model.As("h").
		LeftJoin("content c", "h.content_id=c.id").
		Fields("h.id, h.content_id, h.content_type, h.browse_time, c.title, c.content, c.category, c.status, c.extend").
		Page(req.Page, req.PageSize).
		Scan(&records)
	if err != nil {
		return nil, gerror.New("获取浏览历史记录失败")
	}

	// 处理结果
	for _, record := range records {
		// 提取价格信息
		var price float64 = 0
		if record.Extend != "" {
			var extendData map[string]interface{}
			if err := json.Unmarshal([]byte(record.Extend), &extendData); err == nil {
				if priceValue, ok := extendData["price"]; ok {
					switch v := priceValue.(type) {
					case float64:
						price = v
					case string:
						if p, err := strconv.ParseFloat(v, 64); err == nil {
							price = p
						}
					}
				}
			}
		}

		// 从内容中提取第一张图片作为封面
		cover := ""
		if record.Content != "" {
			re := regexp.MustCompile(`<img[^>]+src="([^"]+)"`)
			matches := re.FindStringSubmatch(record.Content)
			if len(matches) > 1 {
				cover = matches[1]
			}
		}

		// 构建浏览历史记录项
		item := BrowseHistoryItem{
			Id:            record.Id,
			ContentId:     record.ContentId,
			ContentType:   record.ContentType,
			ContentTitle:  record.Title,
			ContentCover:  cover, // 设置提取的封面
			ContentStatus: record.Status,
			BrowseTime:    record.BrowseTime.String(),
			Category:      record.Category,
			Price:         price,
		}

		// 处理内容为null的情况（可能内容已删除）
		if item.ContentTitle == "" {
			item.ContentTitle = "[内容已删除]"
		}

		res.List = append(res.List, item)
	}

	return res, nil
}

// 清空浏览历史记录
func (c *ControllerImpl) Clear(ctx context.Context, req *BrowseHistoryClearReq) (res *BrowseHistoryClearRes, err error) {
	res = &BrowseHistoryClearRes{}

	// 从上下文中获取客户ID
	clientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("获取用户信息失败")
	}

	// 构建删除条件
	model := g.DB().Model("browse_history").
		Where("client_id", clientId)

	// 根据时间类型添加过滤条件
	model = c.applyTimeFilter(ctx, model, req.TimeType)

	// 执行删除
	_, err = model.Delete()
	if err != nil {
		return nil, gerror.New("清空浏览历史记录失败")
	}

	return res, nil
}

// 添加浏览历史记录
func (c *ControllerImpl) Add(ctx context.Context, req *BrowseHistoryAddReq) (res *BrowseHistoryAddRes, err error) {
	res = &BrowseHistoryAddRes{}

	// 从上下文中获取客户ID
	clientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("获取用户信息失败")
	}

	// 记录请求信息
	g.Log().Debug(ctx, "添加浏览历史记录请求",
		"clientId", clientId,
		"contentId", req.ContentId,
		"contentType", req.ContentType)

	// 检查内容是否存在
	contentExists, err := c.checkContentExists(ctx, req.ContentId, req.ContentType)
	if err != nil {
		g.Log().Error(ctx, "内容检查失败",
			"error", err.Error(),
			"contentId", req.ContentId,
			"contentType", req.ContentType)
		return nil, err
	}
	if !contentExists {
		g.Log().Warning(ctx, "内容不存在",
			"contentId", req.ContentId,
			"contentType", req.ContentType)
		return nil, gerror.New("内容不存在")
	}

	// 检查是否已经存在该浏览记录
	count, err := g.DB().Model("browse_history").
		Where("client_id", clientId).
		Where("content_id", req.ContentId).
		Where("content_type", req.ContentType).
		Count()
	if err != nil {
		g.Log().Error(ctx, "检查浏览记录失败", "error", err.Error())
		return nil, gerror.New("检查浏览记录失败")
	}

	// 当前时间
	now := gtime.Now()

	// 如果已存在记录，则更新浏览时间
	if count > 0 {
		_, err = g.DB().Model("browse_history").
			Data(g.Map{
				"browse_time": now,
			}).
			Where("client_id", clientId).
			Where("content_id", req.ContentId).
			Where("content_type", req.ContentType).
			Update()
		if err != nil {
			g.Log().Error(ctx, "更新浏览时间失败", "error", err.Error())
			return nil, gerror.New("更新浏览时间失败")
		}
		g.Log().Debug(ctx, "成功更新浏览时间")
	} else {
		// 否则插入新记录
		_, err = g.DB().Model("browse_history").
			Data(g.Map{
				"client_id":    clientId,
				"content_id":   req.ContentId,
				"content_type": req.ContentType,
				"browse_time":  now,
			}).
			Insert()
		if err != nil {
			g.Log().Error(ctx, "添加浏览记录失败", "error", err.Error())
			return nil, gerror.New("添加浏览记录失败")
		}
		g.Log().Debug(ctx, "成功添加浏览记录")
	}

	return res, nil
}

// 根据时间类型应用过滤条件
func (c *ControllerImpl) applyTimeFilter(_ context.Context, model *gdb.Model, timeType string) *gdb.Model {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch timeType {
	case BrowseHistoryTimeTypeToday:
		// 今天的记录
		model = model.WhereGTE("browse_time", today.Format("2006-01-02 15:04:05"))
	case BrowseHistoryTimeTypeYesterday:
		// 昨天的记录
		yesterdayStart := today.AddDate(0, 0, -1)
		yesterdayEnd := today.Add(-time.Nanosecond)
		model = model.WhereBetween("browse_time", yesterdayStart.Format("2006-01-02 15:04:05"), yesterdayEnd.Format("2006-01-02 15:04:05"))
	case BrowseHistoryTimeTypeThisWeek:
		// 本周记录（周一到现在）
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // 周日视为一周的第7天
		}
		weekStart := today.AddDate(0, 0, -(weekday - 1))
		model = model.WhereGTE("browse_time", weekStart.Format("2006-01-02 15:04:05"))
	case BrowseHistoryTimeTypeThisMonth:
		// 本月记录
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		model = model.WhereGTE("browse_time", monthStart.Format("2006-01-02 15:04:05"))
	case BrowseHistoryTimeTypeAll:
		// 全部记录，不添加时间过滤
	}

	return model
}

// 检查内容是否存在
func (c *ControllerImpl) checkContentExists(ctx context.Context, contentId int, contentType string) (bool, error) {
	var tableName string

	// 根据内容类型确定表名
	switch contentType {
	case "article", "info", "idle":
		tableName = "content"
	default:
		return false, gerror.New(fmt.Sprintf("不支持的内容类型: %s", contentType))
	}

	// 查询内容是否存在，考虑软删除字段
	model := g.DB().Model(tableName).
		Where("id", contentId).
		Where("deleted_at IS NULL") // 增加软删除过滤条件

	// 添加日志记录
	g.Log().Debug(ctx, "检查内容是否存在",
		"contentId", contentId,
		"contentType", contentType,
		"tableName", tableName)

	// 执行查询
	count, err := model.Count()
	if err != nil {
		g.Log().Error(ctx, "检查内容失败",
			"error", err.Error(),
			"contentId", contentId,
			"contentType", contentType)
		return false, gerror.New("检查内容失败")
	}

	return count > 0, nil
}

// 获取浏览历史记录数量
func (c *ControllerImpl) Count(ctx context.Context, req *BrowseHistoryCountReq) (res *BrowseHistoryCountRes, err error) {
	res = &BrowseHistoryCountRes{}

	// 从上下文中获取客户ID
	clientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("获取用户信息失败")
	}

	// 构建查询条件
	model := g.DB().Model("browse_history").
		Where("client_id", clientId)

	// 查询总数
	count, err := model.Count()
	if err != nil {
		g.Log().Error(ctx, "获取浏览历史记录数量失败", "error", err.Error())
		return nil, gerror.New("获取浏览历史记录数量失败")
	}

	res.Count = count
	return res, nil
}
