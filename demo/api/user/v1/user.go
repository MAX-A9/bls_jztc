package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 用户列表请求
type UserListReq struct {
	g.Meta   `path:"/list" method:"get" tags:"用户管理" summary:"获取用户列表" security:"Bearer" description:"获取用户列表，需要管理员权限"`
	Page     int    `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize int    `v:"max:100#每页最大100条" json:"pageSize" dc:"每页数量"`
	Username string `json:"username" dc:"用户名"`
	Nickname string `json:"nickname" dc:"昵称"`
	Status   int    `json:"status" dc:"状态"`
}

// 用户列表响应
type UserListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []UserListItem `json:"list" dc:"用户列表"`
	Total  int            `json:"total" dc:"总数量"`
	Page   int            `json:"page" dc:"当前页码"`
}

// 用户列表项
type UserListItem struct {
	Id            int    `json:"id" dc:"用户ID"`
	Username      string `json:"username" dc:"用户名"`
	Nickname      string `json:"nickname" dc:"昵称"`
	Status        int    `json:"status" dc:"状态"`
	StatusText    string `json:"statusText" dc:"状态文本"`
	LastLoginIp   string `json:"lastLoginIp" dc:"最后登录IP"`
	LastLoginTime string `json:"lastLoginTime" dc:"最后登录时间"`
}

// 创建用户请求
type UserCreateReq struct {
	g.Meta   `path:"/create" method:"post" tags:"用户管理" summary:"创建用户" security:"Bearer" description:"创建新用户，需要管理员权限"`
	Username string `v:"required|length:3,20|regex:^[a-zA-Z0-9_-]+$#用户名不能为空|用户名长度应为3-20个字符|用户名只能包含字母、数字、下划线和连字符" json:"username" dc:"用户名"`
	Password string `v:"required|length:6,20|password#密码不能为空|密码长度不能少于6个字符|密码必须包含大小写字母和数字" json:"password" dc:"密码"`
	Nickname string `v:"required#昵称不能为空" json:"nickname" dc:"昵称"`
}

// 创建用户响应
type UserCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"用户ID"`
}

// 更新用户请求
type UserUpdateReq struct {
	g.Meta   `path:"/update" method:"put" tags:"用户管理" summary:"更新用户" security:"Bearer" description:"更新用户信息，需要管理员权限"`
	Id       int    `v:"required#用户ID不能为空" json:"id" dc:"用户ID"`
	Username string `v:"required|length:3,20|regex:^[a-zA-Z0-9_-]+$#用户名不能为空|用户名长度应为3-20个字符|用户名只能包含字母、数字、下划线和连字符" json:"username" dc:"用户名"`
	Nickname string `v:"required#昵称不能为空" json:"nickname" dc:"昵称"`
	Status   int    `v:"required|in:0,1#状态不能为空|状态值不正确" json:"status" dc:"状态"`
}

// 更新用户响应
type UserUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// 删除用户请求
type UserDeleteReq struct {
	g.Meta `path:"/delete" method:"delete" tags:"用户管理" summary:"删除用户" security:"Bearer" description:"删除用户，需要管理员权限"`
	Id     int `v:"required#用户ID不能为空" json:"id" dc:"用户ID"`
}

// 删除用户响应
type UserDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// 用户个人信息请求
type UserInfoReq struct {
	g.Meta `path:"/info" method:"get" tags:"用户" summary:"获取用户个人信息" security:"Bearer" description:"获取当前登录用户的个人信息"`
}

// 用户个人信息响应
type UserInfoRes struct {
	g.Meta    `mime:"application/json" example:"json"`
	Id        int    `json:"id" dc:"用户ID"`
	Username  string `json:"username" dc:"用户名"`
	Nickname  string `json:"nickname" dc:"昵称"`
	Role      string `json:"role" dc:"角色 admin:管理员"`
	Status    int    `json:"status" dc:"状态 0:禁用 1:正常"`
	LastLogin string `json:"lastLogin" dc:"最后登录时间"`
}
