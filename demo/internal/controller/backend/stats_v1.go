package backend

import (
	"context"
	v1 "demo/api/backend/v1"
	"demo/internal/service"
)

type StatsController struct{}

// GetStatsData 获取统计数据
func (c *StatsController) GetStatsData(ctx context.Context, req *v1.StatsDataReq) (res *v1.StatsDataRes, err error) {
	return service.Stats().GetStatsData(ctx, req)
}

// GetStatsTrend 获取趋势分析数据
func (c *StatsController) GetStatsTrend(ctx context.Context, req *v1.StatsTrendReq) (res *v1.StatsTrendRes, err error) {
	return service.Stats().GetStatsTrend(ctx, req)
}
