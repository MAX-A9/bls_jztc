package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User 用户表相关常量
const (
	// TableUser 表名
	TableUser = "user"
	// UserColumns 所有字段
	UserColumns = "id,username,password,nickname,role,status,last_login_ip,last_login_time,created_at,updated_at,deleted_at"
)

// UserDO DO结构体
type UserDO struct {
	g.Meta        `orm:"table:user, do:true"`
	Id            interface{} `orm:"id,primary" json:"id"`                 // 用户ID
	Username      interface{} `orm:"username" json:"username"`             // 用户名
	Password      interface{} `orm:"password" json:"password"`             // 密码
	Nickname      interface{} `orm:"nickname" json:"nickname"`             // 昵称
	Role          interface{} `orm:"role" json:"role"`                     // 角色 admin:管理员
	Status        interface{} `orm:"status" json:"status"`                 // 状态 0:禁用 1:正常
	LastLoginIp   interface{} `orm:"last_login_ip" json:"lastLoginIp"`     // 最后登录IP
	LastLoginTime *gtime.Time `orm:"last_login_time" json:"lastLoginTime"` // 最后登录时间
	CreatedAt     *gtime.Time `orm:"created_at" json:"createdAt"`          // 创建时间
	UpdatedAt     *gtime.Time `orm:"updated_at" json:"updatedAt"`          // 更新时间
	DeletedAt     *gtime.Time `orm:"deleted_at" json:"deletedAt"`          // 删除时间
}
