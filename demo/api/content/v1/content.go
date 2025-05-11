package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ContentListReq 内容列表请求
type ContentListReq struct {
	g.Meta   `path:"/list" method:"get" tags:"内容管理" summary:"获取内容列表" security:"Bearer" description:"获取内容列表，支持分页和搜索，需要管理员权限"`
	Page     int    `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize int    `v:"max:100#每页最大100条" json:"pageSize" dc:"每页数量"`
	Title    string `json:"title" dc:"标题关键词"`
	Category string `json:"category" dc:"分类"`
	Status   string `json:"status" dc:"状态"`
	Author   string `json:"author" dc:"作者"`
}

// ContentListRes 内容列表响应
type ContentListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []ContentListItem `json:"list" dc:"内容列表"`
	Total  int               `json:"total" dc:"总数量"`
	Page   int               `json:"page" dc:"当前页码"`
}

// ContentListItem 内容列表项
type ContentListItem struct {
	Id            int     `json:"id" dc:"内容ID"`
	Title         string  `json:"title" dc:"标题"`
	Category      string  `json:"category" dc:"分类"`
	Author        string  `json:"author" dc:"作者"`
	Status        string  `json:"status" dc:"状态"`
	StatusText    string  `json:"statusText" dc:"状态文本"`
	Views         int     `json:"views" dc:"浏览量"`
	Likes         int     `json:"likes" dc:"想要数量"`
	Comments      int     `json:"comments" dc:"评论数"`
	IsRecommended bool    `json:"isRecommended" dc:"是否推荐"`
	PublishedAt   *string `json:"publishedAt" dc:"发布时间"`
	ExpiresAt     *string `json:"expiresAt" dc:"到期时间"`
	TopUntil      *string `json:"topUntil" dc:"置顶截止时间"`
	CreatedAt     string  `json:"createdAt" dc:"创建时间"`
	UpdatedAt     string  `json:"updatedAt" dc:"更新时间"`
}

// ContentDetailReq 内容详情请求
type ContentDetailReq struct {
	g.Meta `path:"/detail" method:"get" tags:"内容管理" summary:"获取内容详情" security:"Bearer" description:"获取内容详情，需要管理员权限"`
	Id     int `v:"required#内容ID不能为空" json:"id" dc:"内容ID"`
}

// ContentDetailRes 内容详情响应
type ContentDetailRes struct {
	g.Meta        `mime:"application/json" example:"json"`
	Id            int     `json:"id" dc:"内容ID"`
	Title         string  `json:"title" dc:"标题"`
	Category      string  `json:"category" dc:"分类"`
	Author        string  `json:"author" dc:"作者"`
	Content       string  `json:"content" dc:"内容详情"`
	Status        string  `json:"status" dc:"状态"`
	Views         int     `json:"views" dc:"浏览量"`
	Likes         int     `json:"likes" dc:"想要数量"`
	Comments      int     `json:"comments" dc:"评论数"`
	IsRecommended bool    `json:"isRecommended" dc:"是否推荐"`
	PublishedAt   *string `json:"publishedAt" dc:"发布时间"`
	ExpiresAt     *string `json:"expiresAt" dc:"到期时间"`
	TopUntil      *string `json:"topUntil" dc:"置顶截止时间"`
	CreatedAt     string  `json:"createdAt" dc:"创建时间"`
	UpdatedAt     string  `json:"updatedAt" dc:"更新时间"`
	// 闲置物品特有字段
	Type          int     `json:"type" dc:"内容类型 1:普通信息 2:闲置物品"`
	Price         float64 `json:"price" dc:"价格，闲置物品有效"`
	OriginalPrice float64 `json:"originalPrice" dc:"原价，闲置物品有效"`
	TradePlace    string  `json:"tradePlace" dc:"交易地点，闲置物品有效"`
	TradeMethod   string  `json:"tradeMethod" dc:"交易方式，闲置物品有效"`
}

// ContentCreateReq 创建内容请求
type ContentCreateReq struct {
	g.Meta           `path:"/create" method:"post" tags:"内容管理" summary:"创建内容" security:"Bearer" description:"创建新内容，需要管理员权限"`
	Title            string `v:"required#标题不能为空" json:"title" dc:"标题"`
	Category         string `v:"required#分类不能为空" json:"category" dc:"分类"`
	Author           string `v:"required#作者不能为空" json:"author" dc:"作者"`
	Content          string `v:"required#内容不能为空" json:"content" dc:"内容详情"`
	Status           string `v:"required|in:已发布,待审核,已下架#状态不能为空|状态值不正确" json:"status" dc:"状态 已发布/待审核/已下架"`
	IsRecommended    bool   `json:"isRecommended" dc:"是否推荐"`
	ExpiresAt        string `json:"expiresAt" dc:"到期时间"`
	TopUntil         string `json:"topUntil" dc:"置顶截止时间"`
	TopPackageId     int    `json:"topPackageId" dc:"置顶套餐ID，0表示不使用套餐"`
	PublishPackageId int    `json:"publishPackageId" dc:"发布套餐ID，0表示使用默认3天展示期"`
}

// ContentCreateRes 创建内容响应
type ContentCreateRes struct {
	g.Meta         `mime:"application/json" example:"json"`
	Id             int            `json:"id" dc:"内容ID"`
	HomeCategories []CategoryItem `json:"homeCategories" dc:"首页分类列表"`
	IdleCategories []CategoryItem `json:"idleCategories" dc:"闲置分类列表"`
}

// ContentUpdateReq 更新内容请求
type ContentUpdateReq struct {
	g.Meta           `path:"/update" method:"put" tags:"内容管理" summary:"更新内容" security:"Bearer" description:"更新内容信息，需要管理员权限"`
	Id               int    `v:"required#内容ID不能为空" json:"id" dc:"内容ID"`
	Title            string `v:"required#标题不能为空" json:"title" dc:"标题"`
	Category         string `v:"required#分类不能为空" json:"category" dc:"分类"`
	Author           string `v:"required#作者不能为空" json:"author" dc:"作者"`
	Content          string `v:"required#内容不能为空" json:"content" dc:"内容详情"`
	Status           string `v:"required|in:已发布,待审核,已下架#状态不能为空|状态值不正确" json:"status" dc:"状态 已发布/待审核/已下架"`
	IsRecommended    bool   `json:"isRecommended" dc:"是否推荐"`
	ExpiresAt        string `json:"expiresAt" dc:"到期时间"`
	TopUntil         string `json:"topUntil" dc:"置顶截止时间"`
	TopPackageId     int    `json:"topPackageId" dc:"置顶套餐ID，0表示不使用套餐"`
	PublishPackageId int    `json:"publishPackageId" dc:"发布套餐ID，0表示使用默认3天展示期"`
}

// ContentUpdateRes 更新内容响应
type ContentUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// ContentDeleteReq 删除内容请求
type ContentDeleteReq struct {
	g.Meta `path:"/delete" method:"delete" tags:"内容管理" summary:"删除内容" security:"Bearer" description:"删除内容，需要管理员权限"`
	Id     int `v:"required#内容ID不能为空" json:"id" dc:"内容ID"`
}

// ContentDeleteRes 删除内容响应
type ContentDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// ContentStatusUpdateReq 更新内容状态请求
type ContentStatusUpdateReq struct {
	g.Meta `path:"/status/update" method:"put" tags:"内容管理" summary:"更新内容状态" security:"Bearer" description:"更新内容状态，需要管理员权限"`
	Id     int    `v:"required#内容ID不能为空" json:"id" dc:"内容ID"`
	Status string `v:"required|in:已发布,待审核,已下架#状态不能为空|状态值不正确" json:"status" dc:"状态 已发布/待审核/已下架"`
}

// ContentStatusUpdateRes 更新内容状态响应
type ContentStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// ContentRecommendUpdateReq 更新内容推荐状态请求
type ContentRecommendUpdateReq struct {
	g.Meta        `path:"/recommend/update" method:"put" tags:"内容管理" summary:"更新内容推荐状态" security:"Bearer" description:"更新内容推荐状态，需要管理员权限"`
	Id            int    `v:"required#内容ID不能为空" json:"id" dc:"内容ID"`
	IsRecommended bool   `json:"isRecommended" dc:"是否推荐"`
	TopUntil      string `json:"topUntil" dc:"置顶截止时间"`
	TopPackageId  int    `json:"topPackageId" dc:"置顶套餐ID，0表示不使用套餐"`
}

// ContentRecommendUpdateRes 更新内容推荐状态响应
type ContentRecommendUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// GetCategoriesReq 获取所有分类请求
type GetCategoriesReq struct {
	g.Meta `path:"/categories" method:"get" tags:"内容管理" summary:"获取所有分类" security:"Bearer" description:"获取所有分类，包括首页分类和闲置分类，需要管理员权限"`
}

// GetCategoriesRes 获取所有分类响应
type GetCategoriesRes struct {
	g.Meta         `mime:"application/json" example:"json"`
	HomeCategories []CategoryItem `json:"homeCategories" dc:"首页分类列表"`
	IdleCategories []CategoryItem `json:"idleCategories" dc:"闲置分类列表"`
}
