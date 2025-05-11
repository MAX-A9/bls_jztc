package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 客户列表请求
type ClientListReq struct {
	g.Meta   `path:"/list" method:"get" tags:"客户管理" summary:"获取客户列表" security:"Bearer" description:"获取客户列表，需要管理员权限"`
	Page     int    `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize int    `v:"max:100#每页最大100条" json:"pageSize" dc:"每页数量"`
	Username string `json:"username" dc:"用户名"`
	RealName string `json:"realName" dc:"真实姓名"`
	Phone    string `json:"phone" dc:"手机号"`
	Status   int    `json:"status" dc:"状态 0:禁用 1:正常"`
}

// 客户列表响应
type ClientListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []ClientListItem `json:"list" dc:"客户列表"`
	Total  int              `json:"total" dc:"总数量"`
	Page   int              `json:"page" dc:"当前页码"`
}

// 客户列表项
type ClientListItem struct {
	Id             int    `json:"id" dc:"客户ID"`
	Username       string `json:"username" dc:"用户名"`
	RealName       string `json:"realName" dc:"真实姓名"`
	Phone          string `json:"phone" dc:"手机号"`
	Status         int    `json:"status" dc:"状态 0:禁用 1:正常"`
	StatusText     string `json:"statusText" dc:"状态文本"`
	Identifier     string `json:"identifier" dc:"来源标识 wxapp:小程序 unknown:未知"`
	IdentifierText string `json:"identifierText" dc:"来源标识文本"`
	AvatarUrl      string `json:"avatarUrl" dc:"头像地址"`
	CreatedAt      string `json:"createdAt" dc:"创建时间"`
	LastLoginAt    string `json:"lastLoginAt" dc:"最后登录时间"`
	LastLoginIp    string `json:"lastLoginIp" dc:"最后登录IP"`
}

// 创建客户请求
type ClientCreateReq struct {
	g.Meta     `path:"/create" method:"post" tags:"客户管理" summary:"创建客户" security:"Bearer" description:"创建新客户，需要管理员权限"`
	Username   string `v:"required|length:3,20|regex:^[a-zA-Z0-9_-]+$#用户名不能为空|用户名长度应为3-20个字符|用户名只能包含字母、数字、下划线和连字符" json:"username" dc:"用户名"`
	Password   string `v:"required|length:6,20|password#密码不能为空|密码长度不能少于6个字符|密码必须包含大小写字母和数字" json:"password" dc:"密码"`
	RealName   string `v:"required#姓名不能为空" json:"realName" dc:"真实姓名"`
	Phone      string `v:"required|phone#手机号不能为空|手机号格式不正确" json:"phone" dc:"手机号"`
	Status     int    `v:"required|in:0,1#状态不能为空|状态值不正确" json:"status" dc:"状态 0:禁用 1:正常"`
	Identifier string `v:"required|in:wxapp,unknown#来源标识不能为空|来源标识值不正确" json:"identifier" dc:"来源标识 wxapp:小程序 unknown:未知"`
}

// 创建客户响应
type ClientCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"客户ID"`
}

// 更新客户请求
type ClientUpdateReq struct {
	g.Meta   `path:"/update" method:"put" tags:"客户管理" summary:"更新客户" security:"Bearer" description:"更新客户信息，需要管理员权限"`
	Id       int    `v:"required#客户ID不能为空" json:"id" dc:"客户ID"`
	Username string `v:"required|length:3,20|regex:^[a-zA-Z0-9_-]+$#用户名不能为空|用户名长度应为3-20个字符|用户名只能包含字母、数字、下划线和连字符" json:"username" dc:"用户名"`
	RealName string `v:"required#姓名不能为空" json:"realName" dc:"真实姓名"`
	Phone    string `v:"required|phone#手机号不能为空|手机号格式不正确" json:"phone" dc:"手机号"`
	Status   int    `v:"required|in:0,1#状态不能为空|状态值不正确" json:"status" dc:"状态 0:禁用 1:正常"`
}

// 更新客户响应
type ClientUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// 删除客户请求
type ClientDeleteReq struct {
	g.Meta `path:"/delete" method:"delete" tags:"客户管理" summary:"删除客户" security:"Bearer" description:"删除客户，需要管理员权限"`
	Id     int `v:"required#客户ID不能为空" json:"id" dc:"客户ID"`
}

// 删除客户响应
type ClientDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// 微信小程序登录请求
type WxappLoginReq struct {
	g.Meta `path:"/wx/wxapp-login" method:"post" tags:"客户接口" summary:"小程序登录" description:"微信小程序登录接口，通过code获取openid并登录，没有账号则自动创建基础账号"`
	Code   string `v:"required#小程序code不能为空" json:"code" dc:"微信小程序登录code"`
}

// 微信用户信息（该结构已废弃，仅用于兼容旧代码）
// 根据微信最新规范，应使用头像昵称填写组件让用户主动填写
// <button open-type="chooseAvatar">和<input type="nickname">
type WxUserInfo struct {
	NickName  string `json:"nickName" dc:"昵称"`
	AvatarUrl string `json:"avatarUrl" dc:"头像地址"`
	Gender    int    `json:"gender" dc:"性别 0:未知 1:男 2:女"`
	Province  string `json:"province" dc:"省份"`
	City      string `json:"city" dc:"城市"`
	Country   string `json:"country" dc:"国家"`
}

// 微信小程序登录响应
type WxappLoginRes struct {
	g.Meta   `mime:"application/json" example:"json"`
	ClientId int    `json:"clientId" dc:"客户ID"`
	Token    string `json:"token" dc:"JWT令牌，后续请求需要在Authorization头中携带此令牌"`
	ExpireIn int    `json:"expireIn" dc:"令牌过期时间(秒)"`
}

// 获取客户信息请求
type ClientInfoReq struct {
	g.Meta `path:"/client/info" method:"get" tags:"客户接口" summary:"获取客户信息" security:"Bearer" description:"获取当前登录客户的信息"`
}

// 获取客户信息响应
type ClientInfoRes struct {
	g.Meta     `mime:"application/json" example:"json"`
	Id         int    `json:"id" dc:"客户ID"`
	Username   string `json:"username" dc:"用户名"`
	RealName   string `json:"realName" dc:"真实姓名"`
	Phone      string `json:"phone" dc:"手机号"`
	Status     int    `json:"status" dc:"状态 0:禁用 1:正常"`
	Identifier string `json:"identifier" dc:"来源标识 wxapp:小程序 unknown:未知"`
	AvatarUrl  string `json:"avatarUrl" dc:"头像地址"`
	CreatedAt  string `json:"createdAt" dc:"创建时间"`
}

// 更新客户个人信息请求
type ClientUpdateProfileReq struct {
	g.Meta    `path:"/client/update-profile" method:"put" tags:"客户接口" summary:"更新个人信息" security:"Bearer" description:"客户更新自己的个人信息"`
	RealName  string `json:"realName" dc:"真实姓名"`
	Phone     string `v:"phone#手机号格式不正确" json:"phone" dc:"手机号"`
	AvatarUrl string `json:"avatarUrl" dc:"头像地址，支持base64编码的图片数据"`
}

// 更新客户个人信息响应
type ClientUpdateProfileRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// 客户时长列表请求
type ClientDurationListReq struct {
	g.Meta   `path:"/duration/list" method:"get" tags:"客户时长" summary:"获取客户时长列表" security:"Bearer" description:"获取客户时长列表，需要管理员权限"`
	Page     int `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize int `v:"max:100#每页最大100条" json:"pageSize" dc:"每页数量"`
	ClientId int `json:"clientId" dc:"客户ID"`
}

// 客户时长列表响应
type ClientDurationListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []ClientDurationItem `json:"list" dc:"客户时长列表"`
	Total  int                  `json:"total" dc:"总数量"`
	Page   int                  `json:"page" dc:"当前页码"`
}

// 客户时长项
type ClientDurationItem struct {
	Id                int    `json:"id" dc:"记录ID"`
	ClientId          int    `json:"clientId" dc:"客户ID"`
	ClientName        string `json:"clientName" dc:"客户名称"`
	RemainingDuration string `json:"remainingDuration" dc:"剩余时长"`
	TotalDuration     string `json:"totalDuration" dc:"累计获得"`
	UsedDuration      string `json:"usedDuration" dc:"已使用"`
	CreatedAt         string `json:"createdAt" dc:"创建时间"`
	UpdatedAt         string `json:"updatedAt" dc:"更新时间"`
}

// 获取客户时长详情请求
type ClientDurationDetailReq struct {
	g.Meta `path:"/duration/detail" method:"get" tags:"客户时长" summary:"获取客户时长详情" security:"Bearer" description:"获取客户时长详情，需要管理员权限"`
	Id     int `v:"required#记录ID不能为空" json:"id" dc:"记录ID"`
}

// 获取客户时长详情响应
type ClientDurationDetailRes struct {
	g.Meta            `mime:"application/json" example:"json"`
	Id                int    `json:"id" dc:"记录ID"`
	ClientId          int    `json:"clientId" dc:"客户ID"`
	ClientName        string `json:"clientName" dc:"客户名称"`
	RemainingDuration string `json:"remainingDuration" dc:"剩余时长"`
	TotalDuration     string `json:"totalDuration" dc:"累计获得"`
	UsedDuration      string `json:"usedDuration" dc:"已使用"`
	CreatedAt         string `json:"createdAt" dc:"创建时间"`
	UpdatedAt         string `json:"updatedAt" dc:"更新时间"`
}

// 获取客户端用户剩余时长请求
type WxClientRemainingDurationReq struct {
	g.Meta `path:"/client/duration/remaining" method:"get" tags:"客户接口" summary:"获取剩余时长" security:"Bearer" description:"获取当前登录客户的剩余时长"`
}

// 获取客户端用户剩余时长响应
type WxClientRemainingDurationRes struct {
	g.Meta            `mime:"application/json" example:"json"`
	RemainingDuration string `json:"remainingDuration" dc:"剩余时长"`
}

// 创建客户时长请求
type ClientDurationCreateReq struct {
	g.Meta            `path:"/duration/create" method:"post" tags:"客户时长" summary:"创建客户时长" security:"Bearer" description:"创建客户时长记录，需要管理员权限"`
	ClientId          int    `v:"required#客户ID不能为空" json:"clientId" dc:"客户ID"`
	ClientName        string `v:"required#客户名称不能为空" json:"clientName" dc:"客户名称"`
	RemainingDuration string `v:"required#剩余时长不能为空" json:"remainingDuration" dc:"剩余时长，如3天18小时42分钟"`
	TotalDuration     string `v:"required#累计获得不能为空" json:"totalDuration" dc:"累计获得，如3天18小时42分钟"`
	UsedDuration      string `v:"required#已使用不能为空" json:"usedDuration" dc:"已使用，如4天8小时59分钟"`
}

// 创建客户时长响应
type ClientDurationCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"记录ID"`
}

// 更新客户时长请求
type ClientDurationUpdateReq struct {
	g.Meta            `path:"/duration/update" method:"put" tags:"客户时长" summary:"更新客户时长" security:"Bearer" description:"更新客户时长记录，需要管理员权限"`
	Id                int    `v:"required#记录ID不能为空" json:"id" dc:"记录ID"`
	ClientId          int    `v:"required#客户ID不能为空" json:"clientId" dc:"客户ID"`
	ClientName        string `v:"required#客户名称不能为空" json:"clientName" dc:"客户名称"`
	RemainingDuration string `v:"required#剩余时长不能为空" json:"remainingDuration" dc:"剩余时长，如3天18小时42分钟"`
	TotalDuration     string `v:"required#累计获得不能为空" json:"totalDuration" dc:"累计获得，如3天18小时42分钟"`
	UsedDuration      string `v:"required#已使用不能为空" json:"usedDuration" dc:"已使用，如4天8小时59分钟"`
}

// 更新客户时长响应
type ClientDurationUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// 删除客户时长请求
type ClientDurationDeleteReq struct {
	g.Meta `path:"/duration/delete" method:"delete" tags:"客户时长" summary:"删除客户时长" security:"Bearer" description:"删除客户时长记录，需要管理员权限"`
	Id     int `v:"required#记录ID不能为空" json:"id" dc:"记录ID"`
}

// 删除客户时长响应
type ClientDurationDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}
