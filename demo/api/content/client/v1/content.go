package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// 首页内容列表请求
type HomeContentListReq struct {
	g.Meta   `path:"/home/list" method:"get" tags:"客户端内容" summary:"获取首页内容列表" security:"Bearer" description:"获取首页推荐内容列表"`
	Page     int `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize int `v:"max:50#每页最大50条" json:"pageSize" dc:"每页数量"`
	Type     int `json:"type" dc:"内容类型 0:全部 1:推荐 2:最新"`
}

// 首页内容列表响应
type HomeContentListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []ContentItem `json:"list" dc:"内容列表"`
	Total  int           `json:"total" dc:"总数量"`
	Page   int           `json:"page" dc:"当前页码"`
}

// 闲置内容列表请求
type IdleContentListReq struct {
	g.Meta    `path:"/idle/list" method:"get" tags:"客户端内容" summary:"获取闲置内容列表" security:"Bearer" description:"获取闲置物品内容列表"`
	Page      int    `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize  int    `v:"max:50#每页最大50条" json:"pageSize" dc:"每页数量"`
	Keyword   string `json:"keyword" dc:"搜索关键词"`
	Category  int    `json:"category" dc:"分类ID"`
	SortBy    string `json:"sortBy" dc:"排序方式 price:价格 time:时间"`
	SortOrder string `json:"sortOrder" dc:"排序顺序 asc:升序 desc:降序"`
}

// 闲置内容列表响应
type IdleContentListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []ContentItem `json:"list" dc:"内容列表"`
	Total  int           `json:"total" dc:"总数量"`
	Page   int           `json:"page" dc:"当前页码"`
}

// ContentItem 内容列表项
type ContentItem struct {
	Id             int         `json:"id" dc:"内容ID"`
	Title          string      `json:"title" dc:"标题"`
	Category       string      `json:"category" dc:"分类"`
	Author         string      `json:"author" dc:"作者"`
	Status         string      `json:"status" dc:"状态"`
	Views          int         `json:"views" dc:"浏览量"`
	Likes          int         `json:"likes" dc:"想要数量"`
	Comments       int         `json:"comments" dc:"评论数"`
	IsRecommended  bool        `json:"isRecommended" dc:"是否推荐"`
	Type           int         `json:"type" dc:"内容类型 1:普通信息 2:闲置物品"`
	Price          float64     `json:"price" dc:"价格，闲置物品有效"`
	PublishedAt    *gtime.Time `json:"publishedAt" dc:"发布时间"`
	PublishedAtStr string      `json:"publishedAtStr" dc:"发布时间(格式化)"`
	Publisher      string      `json:"publisher" dc:"发布者，对应author字段"`
	PublishTime    string      `json:"publishTime" dc:"发布时间，前端显示格式，对应publishedAtStr字段"`
	IsTop          bool        `json:"isTop" dc:"是否置顶，对应isRecommended字段"`
	Image          string      `json:"image" dc:"封面图片，对应coverImage字段"`
	Summary        string      `json:"summary" dc:"内容摘要"`
	CoverImage     string      `json:"coverImage" dc:"封面图片"`
	Content        string      `json:"content" dc:"内容详情，详情页返回"`
}

// 作者信息
type Author struct {
	Id       int    `json:"id" dc:"用户ID"`
	Nickname string `json:"nickname" dc:"昵称"`
	Avatar   string `json:"avatar" dc:"头像"`
}

// 分类列表请求
type CategoryListReq struct {
	g.Meta `path:"/wx/client/content/categories" method:"get" tags:"客户端内容" summary:"获取内容分类列表" description:"获取内容分类列表"`
	Type   int `json:"type" dc:"分类类型 1:首页分类 2:闲置分类"`
}

// 分类列表响应
type CategoryListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []CategoryItem `json:"list" dc:"分类列表"`
}

// 分类项
type CategoryItem struct {
	Id    int    `json:"id" dc:"分类ID"`
	Name  string `json:"name" dc:"分类名称"`
	Icon  string `json:"icon,omitempty" dc:"分类图标"`
	Count int    `json:"count,omitempty" dc:"该分类下内容数量"`
	Type  int    `json:"type" dc:"分类类型 1:首页分类 2:闲置分类"`
}

// 微信客户端-闲置发布请求
type WxIdleCreateReq struct {
	g.Meta        `path:"/wx/client/content/idle/create" method:"post" tags:"客户端内容" summary:"发布闲置信息" security:"Bearer" description:"客户端发布闲置物品信息"`
	CategoryId    int      `v:"required#分类ID不能为空" json:"categoryId" dc:"闲置分类ID"`
	RegionId      int      `v:"required#地区ID不能为空" json:"regionId" dc:"所属地区ID"`
	Title         string   `v:"required#标题不能为空" json:"title" dc:"标题"`
	Content       string   `v:"required#内容不能为空" json:"content" dc:"内容描述"`
	Images        []string `v:"required#图片不能为空" json:"images" dc:"图片列表"`
	Price         float64  `v:"required#价格不能为空" json:"price" dc:"价格"`
	OriginalPrice float64  `json:"originalPrice" dc:"原价，选填"`
	TradePlace    string   `v:"required#交易地点不能为空" json:"tradePlace" dc:"交易地点"`
	TradeMethod   string   `v:"required#交易方式不能为空" json:"tradeMethod" dc:"交易方式，如：自提、快递等"`
}

// 微信客户端-闲置发布响应
type WxIdleCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"内容ID"`
}

// 微信客户端-信息发布请求
type WxInfoCreateReq struct {
	g.Meta           `path:"/wx/client/content/info/create" method:"post" tags:"客户端内容" summary:"发布信息" security:"Bearer" description:"客户端发布普通信息"`
	CategoryId       int      `v:"required#分类ID不能为空" json:"categoryId" dc:"信息分类ID"`
	RegionId         int      `v:"required#地区ID不能为空" json:"regionId" dc:"所属地区ID"`
	Title            string   `v:"required#标题不能为空" json:"title" dc:"标题"`
	Content          string   `v:"required#内容不能为空" json:"content" dc:"内容"`
	Images           []string `json:"images" dc:"图片列表"`
	IsTopRequest     bool     `json:"isTopRequest" dc:"是否申请置顶"`
	TopDays          int      `json:"topDays" dc:"置顶天数，申请置顶时必填"`
	TopPackageId     int      `json:"topPackageId" dc:"置顶套餐ID，应用置顶时必填"`
	PublishPackageId int      `v:"required#展示套餐不能为空" json:"publishPackageId" dc:"展示套餐ID，必选"`
}

// 微信客户端-信息发布响应
type WxInfoCreateRes struct {
	g.Meta  `mime:"application/json" example:"json"`
	Id      int    `json:"id" dc:"内容ID"`
	OrderNo string `json:"orderNo" dc:"订单号，有订单时返回"`
}

// 按地区获取内容列表请求
type RegionContentListReq struct {
	g.Meta   `path:"/wx/client/content/region/list" method:"get" tags:"客户端内容" summary:"按地区获取内容列表" description:"按地区获取内容列表，只返回普通信息内容，公开接口"`
	RegionId int    `v:"required#地区ID不能为空" json:"regionId" dc:"地区ID"`
	Page     int    `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize int    `v:"max:50#每页最大50条" json:"pageSize" dc:"每页数量"`
	Category int    `json:"category" dc:"分类ID，0表示全部分类"`
	Keyword  string `json:"keyword" dc:"搜索关键词"`
}

// 按地区获取内容列表响应
type RegionContentListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []RegionContentItem `json:"list" dc:"内容列表"`
	Total  int                 `json:"total" dc:"总数量"`
	Page   int                 `json:"page" dc:"当前页码"`
}

// 按地区内容列表项
type RegionContentItem struct {
	Id          int    `json:"id" dc:"内容ID"`
	Title       string `json:"title" dc:"标题"`
	Category    string `json:"category" dc:"分类"`
	Publisher   string `json:"publisher" dc:"发布者"`
	PublishTime string `json:"publishTime" dc:"发布时间"`
	IsTop       bool   `json:"isTop" dc:"是否置顶"`
	Image       string `json:"image" dc:"封面图片"`
}

// 按地区获取闲置物品列表请求
type RegionIdleListReq struct {
	g.Meta   `path:"/wx/client/content/region/idle/list" method:"get" tags:"客户端内容" summary:"按地区获取闲置物品列表" description:"按地区获取闲置物品列表，只返回闲置物品内容，公开接口"`
	RegionId int    `v:"required#地区ID不能为空" json:"regionId" dc:"地区ID"`
	Page     int    `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize int    `v:"max:50#每页最大50条" json:"pageSize" dc:"每页数量"`
	Category int    `json:"category" dc:"分类ID，0表示全部分类"`
	Keyword  string `json:"keyword" dc:"搜索关键词"`
}

// 按地区获取闲置物品列表响应
type RegionIdleListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []RegionIdleItem `json:"list" dc:"闲置物品列表"`
	Total  int              `json:"total" dc:"总数量"`
	Page   int              `json:"page" dc:"当前页码"`
}

// 按地区闲置物品列表项
type RegionIdleItem struct {
	Id         int     `json:"id" dc:"内容ID"`
	Title      string  `json:"title" dc:"标题"`
	Summary    string  `json:"summary" dc:"摘要，截取content字段中的10个字，过滤图片"`
	TradePlace string  `json:"tradePlace" dc:"交易地点"`
	Price      float64 `json:"price" dc:"价格"`
	Likes      int     `json:"likes" dc:"想要数量"`
	Image      string  `json:"image" dc:"封面图片"`
}

// 公开内容详情请求
type ContentPublicDetailReq struct {
	g.Meta `path:"/wx/client/content/public/detail" method:"get" tags:"客户端内容" summary:"获取内容公开详情" description:"获取内容详情信息，公开接口无需登录"`
	Id     int `v:"required#内容ID不能为空" json:"id" dc:"内容ID"`
}

// 公开内容详情响应
type ContentPublicDetailRes struct {
	g.Meta        `mime:"application/json" example:"json"`
	Id            int      `json:"id" dc:"内容ID"`
	Title         string   `json:"title" dc:"标题"`
	Content       string   `json:"content" dc:"内容详情"`
	Category      string   `json:"category" dc:"分类名称"`
	PublishTime   string   `json:"publishTime" dc:"发布时间"`
	Publisher     string   `json:"publisher" dc:"发布者"`
	PublisherId   uint     `json:"publisher_id" dc:"发布者ID，对应client表的id字段"`
	IsTop         bool     `json:"isTop" dc:"是否置顶"`
	TradePlace    string   `json:"tradePlace,omitempty" dc:"交易地点(闲置物品)"`
	Price         float64  `json:"price,omitempty" dc:"价格(闲置物品)"`
	OriginalPrice float64  `json:"originalPrice,omitempty" dc:"原价(闲置物品)"`
	Views         int      `json:"views" dc:"查看次数"`
	Likes         int      `json:"likes" dc:"想要/收藏数量"`
	TradeMethod   string   `json:"tradeMethod,omitempty" dc:"交易方式(闲置物品)"`
	Images        []string `json:"images" dc:"图片列表"`
}
