package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User 用户表结构
type User struct {
	Id            int         `json:"id"            description:"用户ID"`
	Username      string      `json:"username"      description:"用户名"`
	Password      string      `json:"-"             description:"密码"`
	Nickname      string      `json:"nickname"      description:"昵称"`
	Role          string      `json:"role"          description:"角色 admin:管理员"`
	Status        int         `json:"status"        description:"状态 0:禁用 1:正常"`
	LastLoginIp   string      `json:"lastLoginIp"   description:"最后登录IP"`
	LastLoginTime *gtime.Time `json:"lastLoginTime" description:"最后登录时间"`
	CreatedAt     *gtime.Time `json:"createdAt"     description:"创建时间"`
	UpdatedAt     *gtime.Time `json:"updatedAt"     description:"更新时间"`
	DeletedAt     *gtime.Time `json:"-"             description:"删除时间"`
}
