package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// OrderListReq 订单列表请求
type OrderListReq struct {
	g.Meta     `path:"/list" method:"get" tags:"订单管理" summary:"获取订单列表" security:"Bearer" description:"获取订单列表，需要管理员权限"`
	Page       int    `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize   int    `v:"max:100#每页最大100条" json:"pageSize" dc:"每页数量"`
	OrderNo    string `json:"orderNo" dc:"订单号"`
	ClientName string `json:"clientName" dc:"客户名称"`
	Status     string `json:"status" dc:"订单状态：0-待支付 1-已支付 2-已取消 3-已退款 4-进行中 5-已完成"`
	StartTime  string `json:"startTime" dc:"开始时间"`
	EndTime    string `json:"endTime" dc:"结束时间"`
	Product    string `json:"product" dc:"商品名称"`
}

// OrderListRes 订单列表响应
type OrderListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []OrderListItem `json:"list" dc:"订单列表"`
	Total  int             `json:"total" dc:"总数量"`
	Page   int             `json:"page" dc:"当前页码"`
}

// OrderListItem 订单列表项
type OrderListItem struct {
	Id            int     `json:"id" dc:"订单ID"`
	OrderNo       string  `json:"orderNo" dc:"订单号"`
	ClientName    string  `json:"clientName" dc:"客户名称"`
	ContentId     int     `json:"contentId" dc:"内容ID"`
	ProductName   string  `json:"productName" dc:"商品名称"`
	Amount        float64 `json:"amount" dc:"订单金额"`
	Status        int     `json:"status" dc:"状态：0-待支付 1-已支付 2-已取消 3-已退款"`
	StatusText    string  `json:"statusText" dc:"状态文本"`
	PaymentMethod string  `json:"paymentMethod" dc:"支付方式"`
	CreatedAt     string  `json:"createdAt" dc:"创建时间"`
	PayTime       string  `json:"payTime" dc:"支付时间"`
	ExpireTime    string  `json:"expireTime" dc:"过期时间"`
}

// OrderDetailReq 订单详情请求
type OrderDetailReq struct {
	g.Meta  `path:"/detail" method:"get" tags:"订单管理" summary:"获取订单详情" security:"Bearer" description:"获取订单详情，需要管理员权限"`
	OrderNo string `v:"required#订单号不能为空" json:"orderNo" dc:"订单号"`
}

// OrderDetailRes 订单详情响应
type OrderDetailRes struct {
	g.Meta        `mime:"application/json" example:"json"`
	Id            int     `json:"id" dc:"订单ID"`
	OrderNo       string  `json:"orderNo" dc:"订单号"`
	ClientId      int     `json:"clientId" dc:"客户ID"`
	ClientName    string  `json:"clientName" dc:"客户名称"`
	ContentId     int     `json:"contentId" dc:"内容ID"`
	ProductName   string  `json:"productName" dc:"商品名称"`
	Amount        float64 `json:"amount" dc:"订单金额"`
	Status        int     `json:"status" dc:"状态：0-待支付 1-已支付 2-已取消 3-已退款 4-进行中 5-已完成"`
	StatusText    string  `json:"statusText" dc:"状态文本"`
	PaymentMethod string  `json:"paymentMethod" dc:"支付方式"`
	TransactionId string  `json:"transactionId" dc:"交易流水号"`
	CreatedAt     string  `json:"createdAt" dc:"创建时间"`
	PayTime       string  `json:"payTime" dc:"支付时间"`
	Remark        string  `json:"remark" dc:"备注"`
	PackageInfo   string  `json:"packageInfo" dc:"套餐信息"`
	ExpireTime    string  `json:"expireTime" dc:"到期时间"`
}

// OrderCancelReq 取消订单请求
type OrderCancelReq struct {
	g.Meta  `path:"/cancel" method:"post" tags:"订单管理" summary:"取消订单" security:"Bearer" description:"取消订单，需要管理员权限或订单所属用户权限"`
	OrderNo string `v:"required#订单号不能为空" json:"orderNo" dc:"订单号"`
}

// OrderCancelRes 取消订单响应
type OrderCancelRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// OrderDeleteReq 删除订单请求
type OrderDeleteReq struct {
	g.Meta  `path:"/delete" method:"post" tags:"订单管理" summary:"删除订单" security:"Bearer" description:"删除订单，需要管理员权限"`
	OrderNo string `v:"required#订单号不能为空" json:"orderNo" dc:"订单号"`
}

// OrderDeleteRes 删除订单响应
type OrderDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// UpdateOrderStatusReq 更新订单状态请求
type UpdateOrderStatusReq struct {
	g.Meta        `path:"/update-status" method:"post" tags:"订单管理" summary:"更新订单支付状态" security:"Bearer" description:"更新订单支付状态，需要管理员权限"`
	OrderNo       string `v:"required#订单号不能为空" json:"orderNo" dc:"订单号"`
	Status        int    `v:"required|in:0,1,2,3,4,5#状态不能为空|状态值必须是合法的订单状态值" json:"status" dc:"订单状态：0-待支付 1-已支付 2-已取消 3-已退款 4-进行中 5-已完成"`
	TransactionId string `json:"transactionId" dc:"交易流水号，当状态为已支付时需要提供"`
	PaymentMethod string `json:"paymentMethod" dc:"支付方式，当状态为已支付时需要提供"`
	Remark        string `json:"remark" dc:"备注信息"`
	PackageInfo   string `json:"packageInfo" dc:"套餐信息"`
	ExpireTime    string `json:"expireTime" dc:"到期时间，格式：yyyy-MM-dd HH:mm:ss"`
}

// UpdateOrderStatusRes 更新订单状态响应
type UpdateOrderStatusRes struct {
	g.Meta  `mime:"application/json" example:"json"`
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"处理结果信息"`
}
