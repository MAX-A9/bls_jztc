package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MiniProgramSettings 小程序基础设置表
type MiniProgramSettings struct {
	Id          int         `json:"id"          description:"设置ID"`
	Name        string      `json:"name"        description:"小程序名称"`
	Description string      `json:"description" description:"小程序描述"`
	Logo        string      `json:"logo"        description:"小程序Logo"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:"更新时间"`
}
