package v1

import (
	"context"
	"demo/internal/dao"
	"demo/internal/model/entity"
	"demo/utility/auth"
	"math"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RewardRecordList 获取奖励记录列表
func (c *ControllerImpl) RewardRecordList(ctx context.Context, req *RewardRecordListReq) (res *RewardRecordListRes, err error) {
	// 验证管理员权限
	_, _, _, err = auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("需要管理员权限")
	}

	// 查询条件
	whereCondition := g.Map{}

	// 如果指定了客户ID，则按客户ID查询
	if req.ClientId > 0 {
		whereCondition["client_id"] = req.ClientId
	}

	// 如果指定了状态，则按状态查询
	if req.Status > -1 {
		whereCondition["status"] = req.Status
	}

	// 分页参数处理
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 查询记录总数
	total, err := g.DB().Model("reward_record").
		Where(whereCondition).
		Count()
	if err != nil {
		return nil, err
	}

	// 分页查询记录
	var records []*entity.RewardRecord
	err = g.DB().Model("reward_record").
		Where(whereCondition).
		Order("created_at DESC"). // 按创建时间降序排序
		Limit((req.Page-1)*req.PageSize, req.PageSize).
		Scan(&records)
	if err != nil {
		return nil, err
	}

	// 获取客户信息（用于展示客户名称）
	clientIds := make([]int, 0, len(records))
	for _, record := range records {
		clientIds = append(clientIds, record.ClientId)
	}

	// 查询客户信息
	clientMap := make(map[int]string)
	if len(clientIds) > 0 {
		var clients []struct {
			Id       int    `json:"id"`
			Username string `json:"username"`
		}
		err = g.DB().Model("client").
			Fields("id, username").
			WhereIn("id", clientIds).
			Scan(&clients)
		if err != nil {
			return nil, err
		}

		for _, client := range clients {
			clientMap[client.Id] = client.Username
		}
	}

	// 组装返回数据
	res = &RewardRecordListRes{
		List:  make([]RewardRecordListItem, 0),
		Total: total,
		Page:  req.Page,
	}

	for _, record := range records {
		item := RewardRecordListItem{
			Id:                 record.Id,
			ClientId:           record.ClientId,
			ClientName:         clientMap[record.ClientId],
			RewardMinutes:      record.RewardMinutes,
			RewardDays:         record.RewardDays,
			IsFirstView:        record.IsFirstView,
			RemainingMinutes:   record.RemainingMinutes,
			TotalRewardMinutes: record.TotalRewardMinutes,
			UsedMinutes:        record.UsedMinutes,
			Status:             record.Status,
			StatusText:         getStatusText(record.Status),
			ExpireAt:           record.ExpireAt,
			CreatedAt:          record.CreatedAt,
		}
		res.List = append(res.List, item)
	}

	return res, nil
}

// RewardRecordDetail 获取奖励记录详情
func (c *ControllerImpl) RewardRecordDetail(ctx context.Context, req *RewardRecordDetailReq) (res *RewardRecordDetailRes, err error) {
	// 验证管理员权限
	_, _, _, err = auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("需要管理员权限")
	}

	// 查询记录
	var record entity.RewardRecord
	err = g.DB().Model("reward_record").
		Where("id = ?", req.Id).
		Scan(&record)
	if err != nil {
		return nil, err
	}

	// 如果记录不存在
	if record.Id == 0 {
		return nil, gerror.New("奖励记录不存在")
	}

	// 查询客户信息
	var clientName string
	if record.ClientId > 0 {
		var client struct {
			Username string `json:"username"`
		}
		err = g.DB().Model("client").
			Fields("username").
			Where("id = ?", record.ClientId).
			Scan(&client)
		if err == nil && client.Username != "" {
			clientName = client.Username
		}
	}

	// 组装返回数据
	res = &RewardRecordDetailRes{
		Id:                 record.Id,
		ClientId:           record.ClientId,
		ClientName:         clientName,
		RewardMinutes:      record.RewardMinutes,
		RewardDays:         record.RewardDays,
		IsFirstView:        record.IsFirstView,
		RemainingMinutes:   record.RemainingMinutes,
		TotalRewardMinutes: record.TotalRewardMinutes,
		UsedMinutes:        record.UsedMinutes,
		Status:             record.Status,
		StatusText:         getStatusText(record.Status),
		ExpireAt:           record.ExpireAt,
		CreatedAt:          record.CreatedAt,
		UpdatedAt:          record.UpdatedAt,
	}

	return res, nil
}

// RewardRecordStat 获取奖励记录统计
func (c *ControllerImpl) RewardRecordStat(ctx context.Context, req *RewardRecordStatReq) (res *RewardRecordStatRes, err error) {
	// 验证管理员权限
	_, _, _, err = auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("需要管理员权限")
	}

	// 查询条件
	whereCondition := g.Map{
		"status": 1, // 默认只统计有效的奖励记录
	}

	// 如果指定了客户ID，则按客户ID查询
	if req.ClientId > 0 {
		whereCondition["client_id"] = req.ClientId
	}

	// 获取DAO实例
	rewardRecordDao := &dao.RewardRecordDao{}

	// 获取有效奖励总分钟数
	totalMinutes := 0
	var totalMinutesResult struct {
		Total int `json:"total"`
	}
	err = g.DB().Model("reward_record").
		Where(whereCondition).
		Fields("SUM(remaining_minutes) as total").
		Scan(&totalMinutesResult)
	if err != nil {
		return nil, err
	}
	totalMinutes = totalMinutesResult.Total

	// 获取今日领取总次数
	now := gtime.Now()
	startTime := now.Format("Y-m-d 00:00:00")
	endTime := now.Format("Y-m-d 23:59:59")

	todayCount, err := g.DB().Model("reward_record").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Count()
	if err != nil {
		return nil, err
	}

	// 查询奖励设置
	var rewardSettings struct {
		DailyRewardLimit int `json:"daily_reward_limit"`
	}
	err = g.DB().Model("reward_settings").
		Where("id = 1"). // 假设只有一条记录，ID为1
		Fields("daily_reward_limit").
		Scan(&rewardSettings)
	if err != nil {
		return nil, err
	}

	// 计算剩余可领取次数（只有在指定客户ID时才有意义）
	remainingCount := rewardSettings.DailyRewardLimit
	if req.ClientId > 0 {
		clientTodayCount, err := rewardRecordDao.GetDailyRewardCount(ctx, req.ClientId)
		if err != nil {
			return nil, err
		}
		remainingCount = rewardSettings.DailyRewardLimit - clientTodayCount
		if remainingCount < 0 {
			remainingCount = 0
		}
	}

	// 计算总天数（分钟转换为天）
	totalDays := float64(totalMinutes) / (24 * 60) // 24小时*60分钟=1天的分钟数
	// 保留两位小数
	totalDays = math.Round(totalDays*100) / 100

	// 获取有奖励记录的客户数量
	clientCount, err := g.DB().Model("reward_record").
		Fields("COUNT(DISTINCT client_id) as count").
		Value()
	if err != nil {
		return nil, err
	}

	// 获取已过期的奖励记录数量
	expiredCount, err := g.DB().Model("reward_record").
		Where("status = 0").
		Count()
	if err != nil {
		return nil, err
	}

	// 获取有效的奖励记录数量
	validCount, err := g.DB().Model("reward_record").
		Where("status = 1").
		Count()
	if err != nil {
		return nil, err
	}

	// 组装返回数据
	res = &RewardRecordStatRes{
		TotalMinutes:   totalMinutes,
		TotalDays:      totalDays,
		TodayCount:     todayCount,
		RemainingCount: remainingCount,
		ClientCount:    clientCount.Int(),
		ExpiredCount:   expiredCount,
		ValidCount:     validCount,
	}

	return res, nil
}

// 获取状态文本
func getStatusText(status int) string {
	if status == 0 {
		return "已过期"
	} else {
		return "有效"
	}
}
