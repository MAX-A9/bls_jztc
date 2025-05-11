package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// WxClientCommentListReq 微信客户端获取评论列表请求
type WxClientCommentListReq struct {
	g.Meta    `path:"/wx/client/comment/list" method:"get" tags:"客户端评论" summary:"获取评论列表" description:"获取内容的评论列表，仅获取已审核的评论"`
	ContentId int `v:"required#内容ID不能为空" json:"contentId" dc:"内容ID"`
	Page      int `v:"min:1#页码最小值为1" json:"page" dc:"页码，默认1" d:"1"`
	PageSize  int `v:"max:50#每页最大50条" json:"pageSize" dc:"每页数量，默认10" d:"10"`
}

// WxClientCommentListRes 微信客户端获取评论列表响应
type WxClientCommentListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []WxClientCommentItem `json:"list" dc:"评论列表"`
	Total  int                   `json:"total" dc:"总数量"`
	Page   int                   `json:"page" dc:"当前页码"`
}

// WxClientCommentItem 微信客户端评论列表项
type WxClientCommentItem struct {
	Id        int    `json:"id" dc:"评论ID"`
	RealName  string `json:"realName" dc:"真实姓名"`
	AvatarUrl string `json:"avatarUrl" dc:"头像URL"`
	Comment   string `json:"comment" dc:"评论内容"`
	CreatedAt string `json:"createdAt" dc:"创建时间"`
}

// WxClientCommentCreateReq 微信客户端创建评论请求
type WxClientCommentCreateReq struct {
	g.Meta    `path:"/wx/client/comment/create" method:"post" tags:"客户端评论" summary:"发表评论" security:"WxToken" description:"客户发表评论，需要客户端登录"`
	ContentId int    `v:"required#内容ID不能为空" json:"contentId" dc:"内容ID"`
	Comment   string `v:"required#评论内容不能为空" json:"comment" dc:"评论内容"`
}

// WxClientCommentCreateRes 微信客户端创建评论响应
type WxClientCommentCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"评论ID"`
}
