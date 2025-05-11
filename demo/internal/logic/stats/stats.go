package stats

import (
	"context"
	"fmt"
	"time"

	v1 "demo/api/backend/v1"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sStats struct{}

// New 创建统计服务实例
func New() service.StatsService {
	return &sStats{}
}

func init() {
	service.RegisterStats(New())
}

// getTimeRange 获取时间范围
func getTimeRange(periodType string) (current, previous struct{ start, end string }, err error) {
	now := gtime.Now()

	switch periodType {
	case "week":
		// 本周
		weekDay := int(now.Weekday())
		if weekDay == 0 {
			weekDay = 7
		}
		weekStart := now.AddDate(0, 0, -(weekDay - 1))
		current.start = weekStart.Format("Y-m-d 00:00:00")
		current.end = now.Format("Y-m-d 23:59:59")
		// 上周
		lastWeekStart := weekStart.AddDate(0, 0, -7)
		lastWeekEnd := lastWeekStart.AddDate(0, 0, 6)
		previous.start = lastWeekStart.Format("Y-m-d 00:00:00")
		previous.end = lastWeekEnd.Format("Y-m-d 23:59:59")

	case "month":
		// 本月
		monthStart := gtime.NewFromStr(now.Format("Y-m-01"))
		current.start = monthStart.Format("Y-m-d 00:00:00")
		current.end = now.Format("Y-m-d 23:59:59")
		// 上月
		lastMonthStart := monthStart.AddDate(0, -1, 0)
		lastMonthEnd := monthStart.AddDate(0, 0, -1)
		previous.start = lastMonthStart.Format("Y-m-d 00:00:00")
		previous.end = lastMonthEnd.Format("Y-m-d 23:59:59")

	case "year":
		// 本年
		yearStart := gtime.NewFromStr(now.Format("Y-01-01"))
		current.start = yearStart.Format("Y-m-d 00:00:00")
		current.end = now.Format("Y-m-d 23:59:59")
		// 上年
		lastYearStart := yearStart.AddDate(-1, 0, 0)
		lastYearEnd := yearStart.AddDate(0, 0, -1)
		previous.start = lastYearStart.Format("Y-m-d 00:00:00")
		previous.end = lastYearEnd.Format("Y-m-d 23:59:59")

	default:
		err = gerror.NewCode(gcode.New(400, "无效的周期类型", nil))
	}

	return
}

// GetStatsData 获取统计数据
func (s *sStats) GetStatsData(ctx context.Context, req *v1.StatsDataReq) (res *v1.StatsDataRes, err error) {
	// 权限检查
	userId, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}
	if role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限查看统计数据", nil))
	}

	// 初始化响应
	res = &v1.StatsDataRes{}

	// 获取时间范围
	current, _, err := getTimeRange(req.PeriodType)
	if err != nil {
		return nil, err
	}

	// 记录日志
	g.Log().Info(ctx, "统计数据",
		"userId", userId,
		"periodType", req.PeriodType,
		"currentStart", current.start,
		"currentEnd", current.end,
	)

	// 1. 获取注册客户数量
	count, err := g.DB().Model("client").
		WhereBetween("created_at", current.start, current.end).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "获取注册客户数量失败")
	}
	res.ClientCount = count

	// 2. 获取兑换数量
	count, err = g.DB().Model("exchange_record").
		Where("status", "completed").
		WhereBetween("exchange_time", current.start, current.end).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "获取兑换数量失败")
	}
	res.ExchangeCount = count

	// 3. 获取发布数量
	count, err = g.DB().Model("content").
		WhereBetween("created_at", current.start, current.end).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "获取发布数量失败")
	}
	res.PublishCount = count

	// 4. 获取收益金额
	var revenueRecord struct {
		TotalAmount float64 `json:"total_amount"`
	}
	err = g.DB().Model("`order`").
		Fields("IFNULL(SUM(amount), 0) as total_amount").
		Where("status", 1). // 已支付
		WhereBetween("pay_time", current.start, current.end).
		Scan(&revenueRecord)
	if err != nil {
		return nil, gerror.Wrap(err, "获取收益金额失败")
	}
	res.RevenueAmount = revenueRecord.TotalAmount

	return res, nil
}

// GetStatsTrend 获取趋势分析数据
func (s *sStats) GetStatsTrend(ctx context.Context, req *v1.StatsTrendReq) (res *v1.StatsTrendRes, err error) {
	// 权限检查
	userId, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}
	if role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限查看趋势分析数据", nil))
	}

	// 初始化响应
	res = &v1.StatsTrendRes{
		DataType: req.DataType,
		Period:   req.PeriodType,
	}

	now := gtime.Now()
	var startTime, endTime *gtime.Time
	var interval, format string

	// 根据周期类型确定时间范围和格式
	switch req.PeriodType {
	case "week":
		// 本周
		weekDay := int(now.Weekday())
		if weekDay == 0 {
			weekDay = 7
		}
		weekStart := now.AddDate(0, 0, -(weekDay - 1))
		startTime = weekStart
		endTime = now
		interval = "day"
		format = "m-d"

	case "month":
		// 本月
		startTime = gtime.NewFromStr(now.Format("Y-m-01"))
		endTime = now
		interval = "day"
		format = "m-d"

	case "year":
		// 本年
		startTime = gtime.NewFromStr(now.Format("Y-01-01"))
		endTime = now
		interval = "month"
		format = "Y-m"

	default:
		return nil, gerror.NewCode(gcode.New(400, "无效的周期类型", nil))
	}

	// 记录日志
	g.Log().Info(ctx, "趋势分析",
		"userId", userId,
		"periodType", req.PeriodType,
		"dataType", req.DataType,
		"startTime", startTime.String(),
		"endTime", endTime.String(),
	)

	// 构建完整的日期范围
	fullDateRange := make(map[string]bool)

	if interval == "day" {
		// 按天构建日期范围
		for d := startTime.Time; !d.After(endTime.Time); d = d.AddDate(0, 0, 1) {
			dateStr := gtime.NewFromTime(d).Format("Y-m-d")
			fullDateRange[dateStr] = true
		}
	} else {
		// 按月构建日期范围
		for d := startTime.Time; d.Year() < endTime.Time.Year() || (d.Year() == endTime.Time.Year() && d.Month() <= endTime.Time.Month()); d = d.AddDate(0, 1, 0) {
			dateStr := gtime.NewFromTime(d).Format("Y-m")
			fullDateRange[dateStr] = true
		}
	}

	// 转换为有序的日期
	dates := make([]string, 0, len(fullDateRange))
	for date := range fullDateRange {
		dates = append(dates, date)
	}

	// 按日期排序
	for i := 0; i < len(dates)-1; i++ {
		for j := i + 1; j < len(dates); j++ {
			if dates[i] > dates[j] {
				dates[i], dates[j] = dates[j], dates[i]
			}
		}
	}

	// 格式化日期标签
	res.Labels = make([]string, len(dates))
	for i, date := range dates {
		if interval == "day" {
			t, _ := time.Parse("2006-01-02", date)
			res.Labels[i] = gtime.NewFromTime(t).Format(format)
		} else {
			t, _ := time.Parse("2006-01", date)
			res.Labels[i] = gtime.NewFromTime(t).Format(format)
		}
	}

	// 查询时间范围
	startDateTime := startTime.Format("Y-m-d 00:00:00")
	endDateTime := endTime.Format("Y-m-d 23:59:59")

	// 处理"all"类型，获取所有数据
	if req.DataType == "all" {
		res.AllValues = make(map[string][]interface{})
		dataTypes := []string{"clients", "exchanges", "publishes", "revenue"}

		for _, dataType := range dataTypes {
			// 为每种数据类型准备查询
			var records []struct {
				Date  string      `json:"date"`
				Value interface{} `json:"value"`
			}

			switch dataType {
			case "clients":
				// 客户数量
				err = clientTrendQuery(interval, startDateTime, endDateTime, &records)
			case "exchanges":
				// 兑换数量
				err = exchangeTrendQuery(interval, startDateTime, endDateTime, &records)
			case "publishes":
				// 发布数量
				err = publishTrendQuery(interval, startDateTime, endDateTime, &records)
			case "revenue":
				// 收益金额
				err = revenueTrendQuery(interval, startDateTime, endDateTime, &records)
			}

			if err != nil {
				g.Log().Error(ctx, "获取趋势数据失败", err, "dataType", dataType)
				continue // 错误时跳过当前类型，继续处理其他类型
			}

			// 构建结果数据
			dateValues := make(map[string]interface{})
			for _, record := range records {
				dateValues[record.Date] = record.Value
			}

			// 创建与日期一一对应的数据数组
			values := make([]interface{}, len(dates))

			for i, date := range dates {
				// 获取数值，如果不存在则设置为0
				if value, exists := dateValues[date]; exists {
					values[i] = value
				} else {
					if dataType == "revenue" {
						values[i] = 0.0
					} else {
						values[i] = 0
					}
				}
			}

			// 保存该数据类型的结果
			res.AllValues[dataType] = values

			// 如果是客户数据，同时设置为主要显示数据
			if dataType == "clients" {
				res.Values = values
			}
		}

		return res, nil
	}

	// 处理普通的单一数据类型查询
	var records []struct {
		Date  string      `json:"date"`
		Value interface{} `json:"value"`
	}

	switch req.DataType {
	case "clients":
		// 客户数量
		err = clientTrendQuery(interval, startDateTime, endDateTime, &records)
	case "exchanges":
		// 兑换数量
		err = exchangeTrendQuery(interval, startDateTime, endDateTime, &records)
	case "publishes":
		// 发布数量
		err = publishTrendQuery(interval, startDateTime, endDateTime, &records)
	case "revenue":
		// 收益金额
		err = revenueTrendQuery(interval, startDateTime, endDateTime, &records)
	default:
		return nil, gerror.NewCode(gcode.New(400, "无效的数据类型", nil))
	}

	if err != nil {
		return nil, gerror.Wrap(err, "获取趋势数据失败")
	}

	// 构建结果数据
	dateValues := make(map[string]interface{})
	for _, record := range records {
		dateValues[record.Date] = record.Value
	}

	res.Values = make([]interface{}, len(dates))

	for i, date := range dates {
		// 获取数值，如果不存在则设置为0
		if value, exists := dateValues[date]; exists {
			res.Values[i] = value
		} else {
			if req.DataType == "revenue" {
				res.Values[i] = 0.0
			} else {
				res.Values[i] = 0
			}
		}
	}

	return res, nil
}

// clientTrendQuery 客户趋势查询
func clientTrendQuery(interval, startTime, endTime string, records *[]struct {
	Date  string      `json:"date"`
	Value interface{} `json:"value"`
}) error {
	m := g.DB().Model("client").
		Fields(fmt.Sprintf("DATE_FORMAT(created_at, '%%Y-%%m-%%d') as date, COUNT(id) as value")).
		WhereBetween("created_at", startTime, endTime).
		Group("date").
		Order("date")

	if interval != "day" {
		m = g.DB().Model("client").
			Fields(fmt.Sprintf("DATE_FORMAT(created_at, '%%Y-%%m') as date, COUNT(id) as value")).
			WhereBetween("created_at", startTime, endTime).
			Group("date").
			Order("date")
	}

	return m.Scan(records)
}

// exchangeTrendQuery 兑换趋势查询
func exchangeTrendQuery(interval, startTime, endTime string, records *[]struct {
	Date  string      `json:"date"`
	Value interface{} `json:"value"`
}) error {
	m := g.DB().Model("exchange_record").
		Fields(fmt.Sprintf("DATE_FORMAT(exchange_time, '%%Y-%%m-%%d') as date, COUNT(id) as value")).
		Where("status", "completed").
		WhereBetween("exchange_time", startTime, endTime).
		Group("date").
		Order("date")

	if interval != "day" {
		m = g.DB().Model("exchange_record").
			Fields(fmt.Sprintf("DATE_FORMAT(exchange_time, '%%Y-%%m') as date, COUNT(id) as value")).
			Where("status", "completed").
			WhereBetween("exchange_time", startTime, endTime).
			Group("date").
			Order("date")
	}

	return m.Scan(records)
}

// publishTrendQuery 发布趋势查询
func publishTrendQuery(interval, startTime, endTime string, records *[]struct {
	Date  string      `json:"date"`
	Value interface{} `json:"value"`
}) error {
	m := g.DB().Model("content").
		Fields(fmt.Sprintf("DATE_FORMAT(created_at, '%%Y-%%m-%%d') as date, COUNT(id) as value")).
		WhereBetween("created_at", startTime, endTime).
		Group("date").
		Order("date")

	if interval != "day" {
		m = g.DB().Model("content").
			Fields(fmt.Sprintf("DATE_FORMAT(created_at, '%%Y-%%m') as date, COUNT(id) as value")).
			WhereBetween("created_at", startTime, endTime).
			Group("date").
			Order("date")
	}

	return m.Scan(records)
}

// revenueTrendQuery 收益趋势查询
func revenueTrendQuery(interval, startTime, endTime string, records *[]struct {
	Date  string      `json:"date"`
	Value interface{} `json:"value"`
}) error {
	m := g.DB().Model("`order`").
		Fields(fmt.Sprintf("DATE_FORMAT(pay_time, '%%Y-%%m-%%d') as date, IFNULL(SUM(amount), 0) as value")).
		Where("status", 1).
		WhereBetween("pay_time", startTime, endTime).
		Group("date").
		Order("date")

	if interval != "day" {
		m = g.DB().Model("`order`").
			Fields(fmt.Sprintf("DATE_FORMAT(pay_time, '%%Y-%%m') as date, IFNULL(SUM(amount), 0) as value")).
			Where("status", 1).
			WhereBetween("pay_time", startTime, endTime).
			Group("date").
			Order("date")
	}

	return m.Scan(records)
}
