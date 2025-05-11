package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// FavoriteAddReq 添加收藏请求
type FavoriteAddReq struct {
	g.Meta    `path:"/favorite/add" method:"post" tags:"客户端内容" summary:"添加收藏" security:"Bearer" description:"添加收藏，需要客户端身份验证"`
	ContentId int `v:"required#内容ID不能为空" json:"contentId" dc:"内容ID"`
}

// FavoriteAddRes 添加收藏响应
type FavoriteAddRes struct {
	g.Meta  `mime:"application/json" example:"json"`
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"提示信息"`
}

// FavoriteCancelReq 取消收藏请求
type FavoriteCancelReq struct {
	g.Meta    `path:"/favorite/cancel" method:"post" tags:"客户端内容" summary:"取消收藏" security:"Bearer" description:"取消收藏，需要客户端身份验证"`
	ContentId int `v:"required#内容ID不能为空" json:"contentId" dc:"内容ID"`
}

// FavoriteCancelRes 取消收藏响应
type FavoriteCancelRes struct {
	g.Meta  `mime:"application/json" example:"json"`
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"提示信息"`
}

// FavoriteStatusReq 获取收藏状态请求
type FavoriteStatusReq struct {
	g.Meta    `path:"/favorite/status" method:"get" tags:"客户端内容" summary:"获取收藏状态" security:"Bearer" description:"获取收藏状态，需要客户端身份验证"`
	ContentId int `v:"required#内容ID不能为空" json:"contentId" dc:"内容ID"`
}

// FavoriteStatusRes 获取收藏状态响应
type FavoriteStatusRes struct {
	g.Meta     `mime:"application/json" example:"json"`
	IsFavorite bool `json:"isFavorite" dc:"是否已收藏"`
}

// FavoriteListReq 获取收藏列表请求
type FavoriteListReq struct {
	g.Meta   `path:"/favorite/list" method:"get" tags:"客户端内容" summary:"获取收藏列表" security:"Bearer" description:"获取收藏列表，需要客户端身份验证"`
	Page     int `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize int `v:"max:50#每页最大50条" json:"pageSize" dc:"每页数量"`
	Type     int `json:"type" dc:"内容分类类型 0:全部 1:首页分类 2:闲置分类"`
	Category int `json:"category" dc:"具体分类ID，0表示全部分类"`
}

// FavoriteListRes 获取收藏列表响应
type FavoriteListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []FavoriteItem `json:"list" dc:"收藏列表"`
	Total  int            `json:"total" dc:"总数量"`
	Page   int            `json:"page" dc:"当前页码"`
}

// FavoriteItem 收藏列表项
type FavoriteItem struct {
	Id        int     `json:"id" dc:"内容ID"`
	Title     string  `json:"title" dc:"标题"`
	Category  string  `json:"category" dc:"分类"`
	Publisher string  `json:"publisher" dc:"发布者"`
	Price     float64 `json:"price" dc:"价格，闲置物品有效"`
	Image     string  `json:"image" dc:"封面图片"`
	Type      int     `json:"type" dc:"内容类型 1:普通信息 2:闲置物品"`
}

// FavoriteCountReq 获取收藏总数请求
type FavoriteCountReq struct {
	g.Meta `path:"/favorite/count" method:"get" tags:"客户端内容" summary:"获取收藏总数量" security:"Bearer" description:"获取收藏总数量，需要客户端身份验证"`
}

// FavoriteCountRes 获取收藏总数响应
type FavoriteCountRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Total  int `json:"total" dc:"收藏总数量"`
}
