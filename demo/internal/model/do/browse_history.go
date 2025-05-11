package do

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// BrowseHistoryDO 浏览历史记录数据对象
type BrowseHistoryDO struct {
	Id          int         `json:"id"           description:"记录ID"`
	ClientId    int         `json:"client_id"    description:"客户ID"`
	ContentId   int         `json:"content_id"   description:"内容ID"`
	ContentType string      `json:"content_type" description:"内容类型"`
	BrowseTime  *gtime.Time `json:"browse_time"  description:"浏览时间"`
}
