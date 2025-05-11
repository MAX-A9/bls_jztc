package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Package 套餐表
type Package struct {
	Id           int         `json:"id"           description:"套餐ID"`
	Title        string      `json:"title"        description:"套餐名称"`
	Description  string      `json:"description"  description:"套餐简介"`
	Price        float64     `json:"price"        description:"价格(元)"`
	Type         string      `json:"type"         description:"套餐类型: top-置顶套餐, publish-发布套餐"`
	Duration     int         `json:"duration"     description:"时长值"`
	DurationType string      `json:"durationType" description:"时长单位: hour-小时, day-天, month-月"`
	CreatedAt    *gtime.Time `json:"createdAt"    description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    description:"更新时间"`
	DeletedAt    *gtime.Time `json:"deletedAt"    description:"删除时间"`
}
