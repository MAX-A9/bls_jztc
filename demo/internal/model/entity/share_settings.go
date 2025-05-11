package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ShareSettings 分享设置表
type ShareSettings struct {
	Id                int         `json:"id"               description:"设置ID"`
	DefaultShareText  string      `json:"defaultShareText" description:"默认分享语"`
	DefaultShareImage string      `json:"defaultShareImage" description:"默认分享图片"`
	ContentShareText  string      `json:"contentShareText" description:"内容页分享语"`
	ContentShareImage string      `json:"contentShareImage" description:"内容默认分享图片"`
	HomeShareText     string      `json:"homeShareText"    description:"首页分享语"`
	HomeShareImage    string      `json:"homeShareImage"   description:"首页默认分享图片"`
	CreatedAt         *gtime.Time `json:"createdAt"        description:"创建时间"`
	UpdatedAt         *gtime.Time `json:"updatedAt"        description:"更新时间"`
}
