package client

import (
	"github.com/gogf/gf/v2/frame/g"
)

// WxPayUnifiedOrderReq 微信支付统一下单请求
type WxPayUnifiedOrderReq struct {
	g.Meta   `path:"/unified-order" method:"post" tags:"微信支付" summary:"微信支付统一下单" security:"Bearer" description:"微信支付统一下单接口，用于客户创建微信支付订单"`
	OrderNo  string  `v:"required#订单号不能为空" json:"orderNo" dc:"订单号"`
	TotalFee float64 `v:"required|min:0.01#支付金额不能为空|支付金额不能小于0.01" json:"totalFee" dc:"支付金额（元）"`
	Body     string  `v:"required#商品描述不能为空" json:"body" dc:"商品描述"`
}

// WxPayUnifiedOrderRes 微信支付统一下单响应
type WxPayUnifiedOrderRes struct {
	g.Meta     `mime:"application/json" example:"json"`
	AppId      string `json:"appId" dc:"微信AppID"`
	TimeStamp  string `json:"timeStamp" dc:"时间戳"`
	NonceStr   string `json:"nonceStr" dc:"随机字符串"`
	Package    string `json:"package" dc:"订单详情扩展字符串"`
	SignType   string `json:"signType" dc:"签名类型"`
	PaySign    string `json:"paySign" dc:"签名"`
	PrepayId   string `json:"prepayId" dc:"预支付交易会话标识"`
	CodeUrl    string `json:"codeUrl,omitempty" dc:"二维码链接"`
	MwebUrl    string `json:"mwebUrl,omitempty" dc:"H5支付链接"`
	TradeType  string `json:"tradeType" dc:"交易类型 JSAPI、NATIVE、APP、MWEB"`
	OutTradeNo string `json:"outTradeNo" dc:"商户订单号"`
}

// WxPayNotifyReq 微信支付回调通知请求
type WxPayNotifyReq struct {
	g.Meta `path:"/notify" method:"post" tags:"微信支付" summary:"微信支付回调通知" description:"微信支付结果通知接口，接收微信支付结果"`
}

// WxPayNotifyRes 微信支付回调通知响应
type WxPayNotifyRes struct {
	g.Meta     `mime:"application/xml" example:"xml"`
	ReturnCode string `json:"-" xml:"return_code" dc:"返回状态码"`
	ReturnMsg  string `json:"-" xml:"return_msg" dc:"返回信息"`
}

// WxPayOrderQueryReq 微信支付订单查询请求
type WxPayOrderQueryReq struct {
	g.Meta  `path:"/query-order" method:"post" tags:"微信支付" summary:"微信支付订单查询" security:"Bearer" description:"微信支付订单查询接口，客户查询自己的订单支付状态"`
	OrderNo string `v:"required#订单号不能为空" json:"orderNo" dc:"订单号"`
}

// WxPayOrderQueryRes 微信支付订单查询响应
type WxPayOrderQueryRes struct {
	g.Meta         `mime:"application/json" example:"json"`
	OrderNo        string  `json:"orderNo" dc:"订单号"`
	TransactionId  string  `json:"transactionId" dc:"微信支付订单号"`
	TradeState     string  `json:"tradeState" dc:"交易状态"`
	TradeStateDesc string  `json:"tradeStateDesc" dc:"交易状态描述"`
	PayTime        string  `json:"payTime" dc:"支付完成时间"`
	TotalFee       float64 `json:"totalFee" dc:"订单金额"`
}
