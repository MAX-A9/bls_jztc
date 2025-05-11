package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Content 内容表相关常量
const (
	// TableContent 表名
	TableContent = "content"
	// ContentColumns 所有字段
	ContentColumns = "id,title,category,author,region_id,content,status,views,likes,comments,is_recommended,published_at,expires_at,top_until,created_at,updated_at"
)

// ContentDO DO结构体
type ContentDO struct {
	g.Meta        `orm:"table:content, do:true"`
	Id            interface{} `orm:"id,primary" json:"id"`                // 内容ID
	Title         interface{} `orm:"title" json:"title"`                  // 标题
	Category      interface{} `orm:"category" json:"category"`            // 分类
	Author        interface{} `orm:"author" json:"author"`                // 作者
	RegionId      interface{} `orm:"region_id" json:"regionId"`           // 所属地区ID
	Content       interface{} `orm:"content" json:"content"`              // 内容详情(富文本)
	Status        interface{} `orm:"status" json:"status"`                // 状态：已发布、待审核、已下架
	Views         interface{} `orm:"views" json:"views"`                  // 浏览量
	Likes         interface{} `orm:"likes" json:"likes"`                  // 想要数量
	Comments      interface{} `orm:"comments" json:"comments"`            // 评论数
	IsRecommended interface{} `orm:"is_recommended" json:"isRecommended"` // 是否置顶推荐：1是，0否
	PublishedAt   *gtime.Time `orm:"published_at" json:"publishedAt"`     // 发布时间
	ExpiresAt     *gtime.Time `orm:"expires_at" json:"expiresAt"`         // 到期时间
	TopUntil      *gtime.Time `orm:"top_until" json:"topUntil"`           // 置顶截止时间
	CreatedAt     *gtime.Time `orm:"created_at" json:"createdAt"`         // 创建时间
	UpdatedAt     *gtime.Time `orm:"updated_at" json:"updatedAt"`         // 更新时间
	Extend        interface{} `orm:"extend" json:"extend"`                // 扩展字段，存储JSON数据
	ClientId      interface{} `orm:"client_id" json:"clientId"`           // 客户ID
}
