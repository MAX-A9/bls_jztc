package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AgreementSettings 协议设置表
type AgreementSettings struct {
	Id            int         `json:"id"            description:"设置ID"`
	PrivacyPolicy string      `json:"privacyPolicy" description:"隐私政策内容"`
	UserAgreement string      `json:"userAgreement" description:"用户协议内容"`
	CreatedAt     *gtime.Time `json:"createdAt"     description:"创建时间"`
	UpdatedAt     *gtime.Time `json:"updatedAt"     description:"更新时间"`
}
