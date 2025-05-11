package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Favorite 收藏表相关常量
const (
	// TableFavorite 表名
	TableFavorite = "favorite"
	// FavoriteColumns 所有字段
	FavoriteColumns = "id,client_id,content_id,created_at"
)

// FavoriteDO DO结构体
type FavoriteDO struct {
	g.Meta    `orm:"table:favorite, do:true"`
	Id        interface{} `orm:"id,primary" json:"id"`        // 收藏ID
	ClientId  interface{} `orm:"client_id" json:"clientId"`   // 客户ID
	ContentId interface{} `orm:"content_id" json:"contentId"` // 内容ID
	CreatedAt *gtime.Time `orm:"created_at" json:"createdAt"` // 创建时间
}
