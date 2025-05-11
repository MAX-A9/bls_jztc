package client

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 客户订单列表请求
type OrderListReq struct {
	g.Meta   `path:"/order/list" method:"get" tags:"客户订单" summary:"获取我的订单列表" security:"Bearer" description:"获取当前客户的订单列表"`
	Page     int    `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize int    `v:"max:50#每页最大50条" json:"pageSize" dc:"每页数量"`
	Status   string `json:"status" dc:"订单状态：all-全部 process-进行中 unpaid-待支付 completed-已完成 cancelled-已取消 refunded-已退款"`
}

// 客户订单列表响应
type OrderListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []OrderListItem `json:"list" dc:"订单列表"`
	Total  int             `json:"total" dc:"总数量"`
	Page   int             `json:"page" dc:"当前页码"`
}

// 客户订单列表项
type OrderListItem struct {
	Id            int     `json:"id" dc:"订单ID"`
	OrderNo       string  `json:"orderNo" dc:"订单号"`
	ContentId     int     `json:"contentId" dc:"内容ID"`
	ProductName   string  `json:"productName" dc:"商品名称"`
	Amount        float64 `json:"amount" dc:"订单金额"`
	Status        int     `json:"status" dc:"状态：0-待支付 1-已支付 2-已取消 3-已退款 4-进行中 5-已完成"`
	StatusText    string  `json:"statusText" dc:"状态文本"`
	PaymentMethod string  `json:"paymentMethod" dc:"支付方式"`
	CreatedAt     string  `json:"createdAt" dc:"创建时间"`
	PayTime       string  `json:"payTime" dc:"支付时间"`
	PackageInfo   string  `json:"packageInfo" dc:"套餐信息"`
	ExpireTime    string  `json:"expireTime" dc:"到期时间"`
}

// 客户订单详情请求
type OrderDetailReq struct {
	g.Meta  `path:"/order/detail" method:"get" tags:"客户订单" summary:"获取订单详情" security:"Bearer" description:"获取当前客户的订单详情"`
	OrderNo string `v:"required#订单号不能为空" json:"orderNo" dc:"订单号"`
}

// 客户订单详情响应
type OrderDetailRes struct {
	g.Meta        `mime:"application/json" example:"json"`
	Id            int     `json:"id" dc:"订单ID"`
	OrderNo       string  `json:"orderNo" dc:"订单号"`
	ProductName   string  `json:"productName" dc:"商品名称"`
	Amount        float64 `json:"amount" dc:"订单金额"`
	Status        int     `json:"status" dc:"状态：0-待支付 1-已支付 2-已取消 3-已退款 4-进行中 5-已完成"`
	StatusText    string  `json:"statusText" dc:"状态文本"`
	PaymentMethod string  `json:"paymentMethod" dc:"支付方式"`
	CreatedAt     string  `json:"createdAt" dc:"创建时间"`
	PayTime       string  `json:"payTime" dc:"支付时间"`
	TransactionId string  `json:"transactionId" dc:"交易流水号"`
	Remark        string  `json:"remark" dc:"备注"`
	PackageInfo   string  `json:"packageInfo" dc:"套餐信息"`
	ExpireTime    string  `json:"expireTime" dc:"到期时间"`
}

// 客户订单支付请求
type OrderPayReq struct {
	g.Meta  `path:"/order/pay" method:"post" tags:"客户订单" summary:"支付订单" security:"Bearer" description:"支付待支付状态的订单，避免订单号重复问题"`
	OrderNo string `v:"required#订单号不能为空" json:"orderNo" dc:"订单号"`
}

// 客户订单支付响应
type OrderPayRes struct {
	g.Meta     `mime:"application/json" example:"json"`
	OrderNo    string `json:"orderNo" dc:"原始订单号"`
	PrepayId   string `json:"prepayId" dc:"预支付交易会话标识"`
	NonceStr   string `json:"nonceStr" dc:"随机字符串"`
	TimeStamp  string `json:"timeStamp" dc:"时间戳"`
	Package    string `json:"package" dc:"扩展字段"`
	SignType   string `json:"signType" dc:"签名方式"`
	PaySign    string `json:"paySign" dc:"签名"`
	TotalFee   int    `json:"totalFee" dc:"订单总金额，单位：分"`
	OutTradeNo string `json:"outTradeNo" dc:"商户订单号（临时）"`
}

// 客户取消订单请求
type OrderCancelReq struct {
	g.Meta  `path:"/order/cancel" method:"post" tags:"客户订单" summary:"取消订单" security:"Bearer" description:"取消待支付状态的订单"`
	OrderNo string `v:"required#订单号不能为空" json:"orderNo" dc:"订单号"`
	Reason  string `json:"reason" dc:"取消原因"`
}

// 客户取消订单响应
type OrderCancelRes struct {
	g.Meta    `mime:"application/json" example:"json"`
	OrderNo   string `json:"orderNo" dc:"订单号"`
	Cancelled bool   `json:"cancelled" dc:"是否成功取消"`
}
