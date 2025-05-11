package dao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// BrowseHistoryDao 浏览历史记录数据访问对象
type BrowseHistoryDao struct{}

// GetBrowseHistoryDao 获取浏览历史记录DAO实例
func GetBrowseHistoryDao() *BrowseHistoryDao {
	return &BrowseHistoryDao{}
}

// 浏览历史记录结构
type BrowseHistoryDO struct {
	Id          int         `json:"id"`
	ClientId    int         `json:"client_id"`
	ContentId   int         `json:"content_id"`
	ContentType string      `json:"content_type"`
	BrowseTime  *gtime.Time `json:"browse_time"`
}

// 获取浏览历史记录列表
func (d *BrowseHistoryDao) GetList(ctx context.Context, clientId int, timeType string, page, pageSize int) (list []*BrowseHistoryDO, total int, err error) {
	model := g.DB().Model("browse_history").
		Where("client_id", clientId).
		Order("browse_time DESC")

	// 根据时间类型添加过滤条件
	model = d.ApplyTimeFilter(ctx, model, timeType)

	// 查询总数
	total, err = model.Count()
	if err != nil {
		return nil, 0, gerror.New("获取浏览历史记录总数失败")
	}

	// 分页查询数据
	err = model.Page(page, pageSize).Scan(&list)
	if err != nil {
		return nil, 0, gerror.New("获取浏览历史记录失败")
	}

	return list, total, nil
}

// ApplyTimeFilter 根据时间类型应用过滤条件
func (d *BrowseHistoryDao) ApplyTimeFilter(_ context.Context, model *gdb.Model, timeType string) *gdb.Model {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch timeType {
	case "today":
		// 今天的记录
		model = model.WhereGTE("browse_time", today.Format("2006-01-02 15:04:05"))
	case "yesterday":
		// 昨天的记录
		yesterdayStart := today.AddDate(0, 0, -1)
		yesterdayEnd := today.Add(-time.Nanosecond)
		model = model.WhereBetween("browse_time", yesterdayStart.Format("2006-01-02 15:04:05"), yesterdayEnd.Format("2006-01-02 15:04:05"))
	case "this_week":
		// 本周记录（周一到现在）
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // 周日视为一周的第7天
		}
		weekStart := today.AddDate(0, 0, -(weekday - 1))
		model = model.WhereGTE("browse_time", weekStart.Format("2006-01-02 15:04:05"))
	case "this_month":
		// 本月记录
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		model = model.WhereGTE("browse_time", monthStart.Format("2006-01-02 15:04:05"))
	}

	return model
}

// 清空浏览历史记录
func (d *BrowseHistoryDao) Clear(ctx context.Context, clientId int, timeType string) error {
	model := g.DB().Model("browse_history").
		Where("client_id", clientId)

	// 根据时间类型添加过滤条件
	model = d.ApplyTimeFilter(ctx, model, timeType)

	// 执行删除
	_, err := model.Delete()
	if err != nil {
		return gerror.New("清空浏览历史记录失败")
	}

	return nil
}

// 添加或更新浏览历史记录
func (d *BrowseHistoryDao) AddOrUpdate(ctx context.Context, clientId, contentId int, contentType string) error {
	// 检查是否已经存在该浏览记录
	count, err := g.DB().Model("browse_history").
		Where("client_id", clientId).
		Where("content_id", contentId).
		Where("content_type", contentType).
		Count()
	if err != nil {
		return gerror.New("检查浏览记录失败")
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
			Where("content_id", contentId).
			Where("content_type", contentType).
			Update()
		if err != nil {
			return gerror.New("更新浏览时间失败")
		}
	} else {
		// 否则插入新记录
		_, err = g.DB().Model("browse_history").
			Data(g.Map{
				"client_id":    clientId,
				"content_id":   contentId,
				"content_type": contentType,
				"browse_time":  now,
			}).
			Insert()
		if err != nil {
			return gerror.New("添加浏览记录失败")
		}
	}

	return nil
}
