package settings

import (
	"context"
	"demo/internal/dao"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
)

// InitRewardCleanTask 初始化奖励过期记录清理定时任务
func InitRewardCleanTask() {
	// 从配置中读取定时任务执行频率，默认每天0点执行一次
	cronExpr := g.Cfg().MustGet(gctx.New(), "reward.cronCleanExpired", "0 0 0 * * *").String()

	// 初始化定时任务
	_, err := gcron.Add(gctx.New(), cronExpr, func(ctx context.Context) {
		// 清理过期奖励记录
		cleanExpiredRewards(ctx)
	})

	if err != nil {
		g.Log().Error(gctx.New(), "初始化奖励过期记录清理任务失败:", err)
	} else {
		g.Log().Info(gctx.New(), "奖励过期记录清理任务已初始化，定时表达式：", cronExpr)
	}
}

// cleanExpiredRewards 清理过期奖励记录
func cleanExpiredRewards(ctx context.Context) {
	rewardRecordDao := &dao.RewardRecordDao{}

	// 将过期记录状态更新为已过期
	count, err := rewardRecordDao.DeleteExpiredRecords(ctx)
	if err != nil {
		g.Log().Error(ctx, "清理过期奖励记录失败:", err)
		return
	}

	if count > 0 {
		g.Log().Info(ctx, "成功清理", count, "条过期奖励记录")
	} else {
		g.Log().Debug(ctx, "没有需要清理的过期奖励记录")
	}
}
