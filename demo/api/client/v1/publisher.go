package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 发布人信息请求
type PublisherInfoReq struct {
	g.Meta      `path:"/info" method:"get" tags:"发布人管理" summary:"获取发布人信息" description:"获取内容发布人的信息，包含头像、姓名、发布数量"`
	PublisherId uint `v:"min:1#发布人ID不能为空" json:"publisher_id" dc:"发布人ID"`
}

// 发布人信息响应
type PublisherInfoRes struct {
	g.Meta       `mime:"application/json" example:"json"`
	RealName     string `json:"real_name" dc:"发布人姓名"`
	AvatarUrl    string `json:"avatar_url" dc:"发布人头像"`
	FollowNum    int    `json:"follow_num" dc:"关注数"`
	FansNum      int    `json:"fans_num" dc:"粉丝数"`
	PublishCount int    `json:"publish_count" dc:"已发布内容数量"`
}

// 关注发布人请求
type FollowPublisherReq struct {
	g.Meta `path:"/follow" method:"post" tags:"发布人管理" summary:"关注发布人" security:"Bearer" description:"关注指定的发布人"`
	// 被关注的发布人ID
	PublisherId uint `v:"required#发布人ID不能为空" json:"publisher_id" dc:"发布人ID"`
}

// 关注发布人响应
type FollowPublisherRes struct {
	g.Meta  `mime:"application/json" example:"json"`
	Success bool `json:"success" dc:"操作是否成功"`
}

// 取消关注发布人请求
type UnfollowPublisherReq struct {
	g.Meta `path:"/unfollow" method:"post" tags:"发布人管理" summary:"取消关注发布人" security:"Bearer" description:"取消关注指定的发布人"`
	// 取消关注的发布人ID
	PublisherId uint `v:"required#发布人ID不能为空" json:"publisher_id" dc:"发布人ID"`
}

// 取消关注发布人响应
type UnfollowPublisherRes struct {
	g.Meta  `mime:"application/json" example:"json"`
	Success bool `json:"success" dc:"操作是否成功"`
}

// 获取关注状态请求
type FollowStatusReq struct {
	g.Meta `path:"/follow/status" method:"get" tags:"发布人管理" summary:"获取关注状态" security:"Bearer" description:"获取当前用户是否关注了指定发布人"`
	// 发布人ID
	PublisherId uint `v:"required#发布人ID不能为空" json:"publisher_id" dc:"发布人ID"`
}

// 获取关注状态响应
type FollowStatusRes struct {
	g.Meta     `mime:"application/json" example:"json"`
	IsFollowed bool `json:"is_followed" dc:"是否已关注"`
}

// 获取关注人列表请求
type FollowingListReq struct {
	g.Meta `path:"/following/list" method:"get" tags:"发布人管理" summary:"获取我关注的发布人列表" security:"Bearer" description:"获取当前登录用户关注的发布人列表"`
	Page   int `v:"min:1#页码最小值为1" d:"1" json:"page" dc:"页码"`
	Size   int `v:"max:50#每页最大50条" d:"10" json:"size" dc:"每页数量"`
}

// 获取关注人列表响应
type FollowingListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []FollowingItem `json:"list" dc:"关注人列表"`
	Total  int             `json:"total" dc:"总数"`
	Page   int             `json:"page" dc:"当前页码"`
	Size   int             `json:"size" dc:"每页数量"`
}

// 关注人列表项
type FollowingItem struct {
	ClientId     uint   `json:"client_id" dc:"发布人ID"`
	RealName     string `json:"real_name" dc:"发布人姓名"`
	AvatarUrl    string `json:"avatar_url" dc:"发布人头像"`
	PublishCount int    `json:"publish_count" dc:"已发布内容数量"`
	FollowTime   string `json:"follow_time" dc:"关注时间"`
}

// 获取关注人总数请求
type FollowingCountReq struct {
	g.Meta `path:"/following/count" method:"get" tags:"发布人管理" summary:"获取我关注的发布人总数" security:"Bearer" description:"获取当前登录用户关注的发布人总数"`
}

// 获取关注人总数响应
type FollowingCountRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Count  int `json:"count" dc:"关注人总数"`
}
