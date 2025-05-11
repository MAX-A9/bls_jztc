package dao

import (
	"context"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RewardRecordDao 奖励记录数据访问对象
type RewardRecordDao struct{}

// Get 根据ID获取奖励记录
func (d *RewardRecordDao) Get(ctx context.Context, id int) (*entity.RewardRecord, error) {
	var record entity.RewardRecord
	err := g.DB().Model("reward_record").Where("id=?", id).Scan(&record)
	if err != nil || record.Id == 0 {
		return nil, err
	}
	return &record, nil
}

// GetListByClientId 获取指定客户的有效奖励记录列表
func (d *RewardRecordDao) GetListByClientId(ctx context.Context, clientId int) ([]*entity.RewardRecord, error) {
	var records []*entity.RewardRecord
	err := g.DB().Model("reward_record").
		Where("client_id=?", clientId).
		Where("status=?", 1). // 有效状态
		Where("expire_at>?", gtime.Now()).
		Order("expire_at ASC"). // 按过期时间升序排序，先过期的记录先使用
		Scan(&records)
	return records, err
}

// GetDailyRewardCount 获取指定客户当天的奖励记录数量
func (d *RewardRecordDao) GetDailyRewardCount(ctx context.Context, clientId int) (int, error) {
	now := gtime.Now()
	startTime := now.Format("Y-m-d 00:00:00")
	endTime := now.Format("Y-m-d 23:59:59")

	count, err := g.DB().Model("reward_record").
		Where("client_id=?", clientId).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Count()
	return count, err
}

// GetTotalRewardMinutes 获取指定客户的有效奖励总分钟数
func (d *RewardRecordDao) GetTotalRewardMinutes(ctx context.Context, clientId int) (int, error) {
	var result struct {
		Total int `json:"total"`
	}
	err := g.DB().Model("reward_record").
		Where("client_id=?", clientId).
		Where("status=?", 1). // 有效状态
		Where("expire_at>?", gtime.Now()).
		Fields("SUM(remaining_minutes) as total").
		Scan(&result)
	if err != nil {
		return 0, err
	}
	return result.Total, nil
}

// Save 保存奖励记录
func (d *RewardRecordDao) Save(ctx context.Context, record *entity.RewardRecord) error {
	now := gtime.Now()

	if record.Id > 0 {
		// 更新已有记录
		_, err := g.DB().Model("reward_record").
			Data(g.Map{
				"remaining_minutes": record.RemainingMinutes,
				"used_minutes":      record.UsedMinutes,
				"status":            record.Status,
				"updated_at":        now,
			}).
			Where("id=?", record.Id).
			Update()
		return err
	} else {
		// 创建新记录
		_, err := g.DB().Model("reward_record").
			Data(g.Map{
				"client_id":            record.ClientId,
				"reward_minutes":       record.RewardMinutes,
				"reward_days":          record.RewardDays,
				"is_first_view":        record.IsFirstView,
				"remaining_minutes":    record.RemainingMinutes,
				"total_reward_minutes": record.TotalRewardMinutes,
				"used_minutes":         record.UsedMinutes,
				"status":               record.Status,
				"expire_at":            record.ExpireAt,
				"created_at":           now,
				"updated_at":           now,
			}).
			Insert()
		return err
	}
}

// DeleteExpiredRecords 删除过期的奖励记录
func (d *RewardRecordDao) DeleteExpiredRecords(ctx context.Context) (int64, error) {
	now := gtime.Now()

	// 将过期记录状态更新为已过期
	result, err := g.DB().Model("reward_record").
		Data(g.Map{
			"status":     0, // 已过期状态
			"updated_at": now,
		}).
		Where("expire_at<?", now).
		Where("status=?", 1). // 只更新当前为有效状态的记录
		Update()

	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	return affected, err
}

// CheckFirstView 检查是否是客户的首次观看
func (d *RewardRecordDao) CheckFirstView(ctx context.Context, clientId int) (bool, error) {
	count, err := g.DB().Model("reward_record").
		Where("client_id=?", clientId).
		Count()
	if err != nil {
		return false, err
	}
	return count == 0, nil // 如果没有记录则是首次观看
}
