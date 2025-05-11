package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemConfig 系统配置表
type SystemConfig struct {
	Id          int         `json:"id" description:"配置ID"`
	Module      string      `json:"module" description:"模块名称"`
	Key         string      `json:"key" description:"配置键名"`
	Value       string      `json:"value" description:"配置值"`
	Description string      `json:"description" description:"配置描述"`
	CreatedAt   *gtime.Time `json:"createdAt" description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt" description:"更新时间"`
}
