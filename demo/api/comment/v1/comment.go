package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CommentListReq 评论列表请求
type CommentListReq struct {
	g.Meta    `path:"/list" method:"get" tags:"评论管理" summary:"获取评论列表" security:"Bearer" description:"获取评论列表，支持分页和搜索，需要管理员权限"`
	Page      int    `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize  int    `v:"max:100#每页最大100条" json:"pageSize" dc:"每页数量"`
	ContentId int    `json:"contentId" dc:"内容ID，不填则查询所有内容的评论"`
	Status    string `json:"status" dc:"状态：已审核、待审核、已拒绝，不填则查询所有状态"`
	RealName  string `json:"realName" dc:"真实姓名关键词"`
	Comment   string `json:"comment" dc:"评论内容关键词"`
}

// CommentListRes 评论列表响应
type CommentListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []CommentListItem `json:"list" dc:"评论列表"`
	Total  int               `json:"total" dc:"总数量"`
	Page   int               `json:"page" dc:"当前页码"`
}

// CommentListItem 评论列表项
type CommentListItem struct {
	Id           int    `json:"id" dc:"评论ID"`
	ContentId    int    `json:"contentId" dc:"内容ID"`
	ContentTitle string `json:"contentTitle" dc:"内容标题"`
	ClientId     int    `json:"clientId" dc:"客户ID"`
	RealName     string `json:"realName" dc:"真实姓名"`
	Comment      string `json:"comment" dc:"评论内容"`
	Status       string `json:"status" dc:"状态"`
	StatusText   string `json:"statusText" dc:"状态文本"`
	CreatedAt    string `json:"createdAt" dc:"创建时间"`
	UpdatedAt    string `json:"updatedAt" dc:"更新时间"`
}

// CommentDetailReq 评论详情请求
type CommentDetailReq struct {
	g.Meta `path:"/detail" method:"get" tags:"评论管理" summary:"获取评论详情" security:"Bearer" description:"获取评论详情，需要管理员权限"`
	Id     int `v:"required#评论ID不能为空" json:"id" dc:"评论ID"`
}

// CommentDetailRes 评论详情响应
type CommentDetailRes struct {
	g.Meta       `mime:"application/json" example:"json"`
	Id           int    `json:"id" dc:"评论ID"`
	ContentId    int    `json:"contentId" dc:"内容ID"`
	ContentTitle string `json:"contentTitle" dc:"内容标题"`
	ClientId     int    `json:"clientId" dc:"客户ID"`
	RealName     string `json:"realName" dc:"真实姓名"`
	Comment      string `json:"comment" dc:"评论内容"`
	Status       string `json:"status" dc:"状态"`
	CreatedAt    string `json:"createdAt" dc:"创建时间"`
	UpdatedAt    string `json:"updatedAt" dc:"更新时间"`
}

// ContentCommentsReq 内容评论列表请求
type ContentCommentsReq struct {
	g.Meta    `path:"/content-comments" method:"get" tags:"评论管理" summary:"获取指定内容的评论列表" security:"Bearer" description:"获取指定内容的评论列表，支持分页，需要管理员权限"`
	ContentId int `v:"required#内容ID不能为空" json:"contentId" dc:"内容ID"`
	Page      int `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize  int `v:"max:100#每页最大100条" json:"pageSize" dc:"每页数量"`
}

// ContentCommentsRes 内容评论列表响应
type ContentCommentsRes struct {
	g.Meta       `mime:"application/json" example:"json"`
	ContentId    int               `json:"contentId" dc:"内容ID"`
	ContentTitle string            `json:"contentTitle" dc:"内容标题"`
	List         []CommentListItem `json:"list" dc:"评论列表"`
	Total        int               `json:"total" dc:"总数量"`
	Page         int               `json:"page" dc:"当前页码"`
}

// CommentCreateReq 创建评论请求
type CommentCreateReq struct {
	g.Meta    `path:"/create" method:"post" tags:"评论管理" summary:"创建评论" security:"Bearer" description:"创建新评论，需要管理员权限"`
	ContentId int    `v:"required#内容ID不能为空" json:"contentId" dc:"内容ID"`
	ClientId  int    `v:"required#客户ID不能为空" json:"clientId" dc:"客户ID"`
	RealName  string `v:"required#真实姓名不能为空" json:"realName" dc:"真实姓名"`
	Comment   string `v:"required#评论内容不能为空" json:"comment" dc:"评论内容"`
	Status    string `v:"required|in:已审核,待审核,已拒绝#状态不能为空|状态值不正确" json:"status" dc:"状态：已审核、待审核、已拒绝"`
}

// CommentCreateRes 创建评论响应
type CommentCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"评论ID"`
}

// CommentUpdateReq 更新评论请求
type CommentUpdateReq struct {
	g.Meta   `path:"/update" method:"put" tags:"评论管理" summary:"更新评论" security:"Bearer" description:"更新评论信息，需要管理员权限"`
	Id       int    `v:"required#评论ID不能为空" json:"id" dc:"评论ID"`
	RealName string `v:"required#真实姓名不能为空" json:"realName" dc:"真实姓名"`
	Comment  string `v:"required#评论内容不能为空" json:"comment" dc:"评论内容"`
	Status   string `v:"required|in:已审核,待审核,已拒绝#状态不能为空|状态值不正确" json:"status" dc:"状态：已审核、待审核、已拒绝"`
}

// CommentUpdateRes 更新评论响应
type CommentUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// CommentDeleteReq 删除评论请求
type CommentDeleteReq struct {
	g.Meta `path:"/delete" method:"delete" tags:"评论管理" summary:"删除评论" security:"Bearer" description:"删除评论，需要管理员权限"`
	Id     int `v:"required#评论ID不能为空" json:"id" dc:"评论ID"`
}

// CommentDeleteRes 删除评论响应
type CommentDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// CommentStatusUpdateReq 更新评论状态请求
type CommentStatusUpdateReq struct {
	g.Meta `path:"/status/update" method:"put" tags:"评论管理" summary:"更新评论状态" security:"Bearer" description:"更新评论状态，需要管理员权限"`
	Id     int    `v:"required#评论ID不能为空" json:"id" dc:"评论ID"`
	Status string `v:"required|in:已审核,待审核,已拒绝#状态不能为空|状态值不正确" json:"status" dc:"状态：已审核、待审核、已拒绝"`
}

// CommentStatusUpdateRes 更新评论状态响应
type CommentStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}
