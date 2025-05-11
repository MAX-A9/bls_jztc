package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ContentComment 评论表相关常量
const (
	// TableContentComment 表名
	TableContentComment = "content_comment"
	// ContentCommentColumns 所有字段
	ContentCommentColumns = "id,content_id,client_id,real_name,comment,status,created_at,updated_at"
)

// ContentCommentDO DO结构体
type ContentCommentDO struct {
	g.Meta    `orm:"table:content_comment, do:true"`
	Id        interface{} `orm:"id,primary" json:"id"`        // 评论ID
	ContentId interface{} `orm:"content_id" json:"contentId"` // 内容ID
	ClientId  interface{} `orm:"client_id" json:"clientId"`   // 客户ID
	RealName  interface{} `orm:"real_name" json:"realName"`   // 真实姓名
	Comment   interface{} `orm:"comment" json:"comment"`      // 评论内容
	Status    interface{} `orm:"status" json:"status"`        // 状态：已审核、待审核、已拒绝
	CreatedAt *gtime.Time `orm:"created_at" json:"createdAt"` // 创建时间
	UpdatedAt *gtime.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
}
