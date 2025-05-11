package do

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// 专属管家数据表名
const TableButler = "butler"

// ButlerDO 专属管家数据对象
type ButlerDO struct {
	Id        int64       `json:"id"         description:"主键ID"`
	ImageUrl  string      `json:"imageUrl"   description:"图片地址"`
	Status    int         `json:"status"     description:"状态 1:启用 0:禁用"`
	CreatedAt *gtime.Time `json:"createdAt"  description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt"  description:"更新时间"`
}
