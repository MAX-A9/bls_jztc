package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ButlerReq 专属管家请求
type ButlerReq struct {
	g.Meta `path:"/butler" tags:"专属管家" method:"get" summary:"获取专属管家" security:"Bearer"`
}

// ButlerRes 专属管家响应
type ButlerRes struct {
	Id        int64       `json:"id"         description:"主键ID"`
	ImageUrl  string      `json:"imageUrl"   description:"图片地址"`
	Status    int         `json:"status"     description:"状态 1:启用 0:禁用"`
	CreatedAt *gtime.Time `json:"createdAt"  description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt"  description:"更新时间"`
}

// ========== 管理端接口 ==========

// SaveImageReq 保存专属管家图片请求
type SaveImageReq struct {
	g.Meta   `path:"/butler/image/save" method:"post" tags:"专属管家" summary:"保存专属管家图片" security:"Bearer"`
	ImageUrl string `v:"required#图片URL不能为空" json:"imageUrl" dc:"管家图片URL"`
	Status   int    `v:"required|in:0,1#状态不能为空|状态只能是0或1" d:"1" json:"status" dc:"状态 0:禁用 1:启用，默认1"`
}

// SaveImageRes 保存专属管家图片响应
type SaveImageRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"管家ID"`
}

// ========== 客户端接口 ==========

// GetImageReq 客户端获取专属管家图片请求
type GetImageReq struct {
	g.Meta `path:"/wx/client/butler/image" method:"get" tags:"客户端专属管家" summary:"获取专属管家图片"`
}

// GetImageRes 客户端获取专属管家图片响应
type GetImageRes struct {
	g.Meta   `mime:"application/json" example:"json"`
	ImageUrl string `json:"imageUrl" dc:"管家图片URL"`
}
