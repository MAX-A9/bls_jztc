package client

import (
	"context"
	v1 "demo/api/client/v1"
)

// Reward 奖励记录控制器
type Reward struct {
	v1 *RewardV1
}

// RewardV1 奖励记录V1控制器
type RewardV1 struct {
	iController v1.IRewardController
}

// NewReward 创建奖励记录控制器
func NewReward() *Reward {
	return &Reward{
		v1: &RewardV1{
			iController: &v1.ControllerImpl{},
		},
	}
}

// V1 获取V1版本控制器
func (c *Reward) V1() *RewardV1 {
	return c.v1
}

// RewardRecordList 获取奖励记录列表
func (c *RewardV1) RewardRecordList(ctx context.Context, req *v1.RewardRecordListReq) (res *v1.RewardRecordListRes, err error) {
	return c.iController.RewardRecordList(ctx, req)
}

// RewardRecordDetail 获取奖励记录详情
func (c *RewardV1) RewardRecordDetail(ctx context.Context, req *v1.RewardRecordDetailReq) (res *v1.RewardRecordDetailRes, err error) {
	return c.iController.RewardRecordDetail(ctx, req)
}

// RewardRecordStat 获取奖励记录统计
func (c *RewardV1) RewardRecordStat(ctx context.Context, req *v1.RewardRecordStatReq) (res *v1.RewardRecordStatRes, err error) {
	return c.iController.RewardRecordStat(ctx, req)
}
