package service

import (
	"context"
	v1 "demo/api/backend/v1"
)

// StatsService 统计服务接口
type StatsService interface {
	// GetStatsData 获取统计数据
	GetStatsData(ctx context.Context, req *v1.StatsDataReq) (res *v1.StatsDataRes, err error)

	// GetStatsTrend 获取趋势分析数据
	GetStatsTrend(ctx context.Context, req *v1.StatsTrendReq) (res *v1.StatsTrendRes, err error)
}

var localStats StatsService

// Stats 获取统计服务
func Stats() StatsService {
	if localStats == nil {
		panic("implement not found for interface StatsService, forgot register?")
	}
	return localStats
}

// RegisterStats 注册统计服务
func RegisterStats(i StatsService) {
	localStats = i
}
