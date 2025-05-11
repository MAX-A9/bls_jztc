package model

import "github.com/gogf/gf/v2/os/gtime"

// Region 地区模型
type Region struct {
	Id        int         `json:"id"         description:"地区ID"`
	Location  string      `json:"location"   description:"所在地区，如：北京市/市辖区/东城区"`
	Name      string      `json:"name"       description:"地区名称"`
	Level     string      `json:"level"      description:"级别: 省,县,乡"`
	Status    int         `json:"status"     description:"状态 0:启用 1:禁用"`
	CreatedAt *gtime.Time `json:"createdAt"  description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt"  description:"更新时间"`
}
