package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// InnerPageBanner 内页轮播图实体
type InnerPageBanner struct {
	Id         int         `json:"id"          description:"内页轮播图ID"`
	Image      string      `json:"image"       description:"轮播图片地址"`
	LinkType   string      `json:"link_type"   description:"跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页"`
	LinkUrl    string      `json:"link_url"    description:"跳转地址"`
	IsEnabled  int         `json:"is_enabled"  description:"是否启用：0-禁用，1-启用"`
	Order      int         `json:"order"       description:"排序值，越小越靠前"`
	BannerType string      `json:"banner_type" description:"轮播图类型：home-首页轮播，idle-闲置轮播"`
	CreatedAt  *gtime.Time `json:"created_at"  description:"创建时间"`
	UpdatedAt  *gtime.Time `json:"updated_at"  description:"更新时间"`
}
