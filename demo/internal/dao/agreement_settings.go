package dao

import (
	"context"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AgreementSettingsDao 协议设置数据访问对象
type AgreementSettingsDao struct{}

// Get 获取协议设置
func (d *AgreementSettingsDao) Get(ctx context.Context) (*entity.AgreementSettings, error) {
	var settings entity.AgreementSettings
	err := g.DB().Model("agreement_settings").Where("id=?", 1).Scan(&settings)
	// 如果没有设置数据，初始化一个空对象
	if err != nil || settings.Id == 0 {
		return &entity.AgreementSettings{
			Id:            1,
			PrivacyPolicy: "隐私政策内容",
			UserAgreement: "用户协议内容",
		}, nil
	}
	return &settings, nil
}

// Save 保存协议设置
func (d *AgreementSettingsDao) Save(ctx context.Context, settings *entity.AgreementSettings) error {
	now := gtime.Now()
	// 检查是否已有记录
	var count int
	count, err := g.DB().Model("agreement_settings").Where("id=?", 1).Count()
	if err != nil {
		return err
	}

	// 如果已存在记录则更新，否则创建
	if count > 0 {
		_, err = g.DB().Model("agreement_settings").
			Data(g.Map{
				"privacy_policy": settings.PrivacyPolicy,
				"user_agreement": settings.UserAgreement,
				"updated_at":     now,
			}).
			Where("id=?", 1).
			Update()
		return err
	} else {
		_, err = g.DB().Model("agreement_settings").
			Data(g.Map{
				"id":             1,
				"privacy_policy": settings.PrivacyPolicy,
				"user_agreement": settings.UserAgreement,
				"created_at":     now,
				"updated_at":     now,
			}).
			Insert()
		return err
	}
}
