package model

import "github.com/gogf/gf/v2/os/gtime"

// Order 订单模型
type Order struct {
	Id            int         `json:"id"`            // 订单ID
	OrderNo       string      `json:"orderNo"`       // 订单号
	ClientId      int         `json:"clientId"`      // 客户ID
	ClientName    string      `json:"clientName"`    // 客户名称
	ContentId     int         `json:"contentId"`     // 内容ID
	ProductName   string      `json:"productName"`   // 商品名称
	Amount        float64     `json:"amount"`        // 订单金额
	Status        int         `json:"status"`        // 状态：0-待支付 1-已支付 2-已取消 3-已退款 4-进行中 5-已完成
	PaymentMethod string      `json:"paymentMethod"` // 支付方式
	PayTime       *gtime.Time `json:"payTime"`       // 支付时间
	ExpireTime    *gtime.Time `json:"expireTime"`    // 订单过期时间
	TransactionId string      `json:"transactionId"` // 交易流水号
	CreatedAt     *gtime.Time `json:"createdAt"`     // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"`     // 更新时间
	PackageInfo   string      `json:"packageInfo"`   // 套餐信息，JSON格式
	Remark        string      `json:"remark"`        // 备注
}

// OrderListItem 订单列表项
type OrderListItem struct {
	Id            int         `json:"id"`            // 订单ID
	OrderNo       string      `json:"orderNo"`       // 订单号
	ClientName    string      `json:"clientName"`    // 客户名称
	ContentId     int         `json:"contentId"`     // 内容ID
	ProductName   string      `json:"productName"`   // 商品名称
	Amount        float64     `json:"amount"`        // 订单金额
	Status        int         `json:"status"`        // 状态
	StatusText    string      `json:"statusText"`    // 状态文本
	PaymentMethod string      `json:"paymentMethod"` // 支付方式
	CreatedAt     *gtime.Time `json:"createdAt"`     // 创建时间
	PayTime       *gtime.Time `json:"payTime"`       // 支付时间
}

// 微信支付相关模型

// WxPayConfig 微信支付配置
type WxPayConfig struct {
	AppId     string // 微信支付AppID
	MchId     string // 微信支付商户号
	ApiKey    string // API密钥
	NotifyUrl string // 回调通知地址
}

// WxPayOrder 微信支付订单
type WxPayOrder struct {
	AppId       string  `json:"appId"`       // 微信支付AppID
	MchId       string  `json:"mchId"`       // 微信支付商户号
	OutTradeNo  string  `json:"outTradeNo"`  // 商户订单号
	Body        string  `json:"body"`        // 商品描述
	TotalFee    int     `json:"totalFee"`    // 支付金额(分)
	TotalAmount float64 `json:"totalAmount"` // 支付金额(元)
	NotifyUrl   string  `json:"notifyUrl"`   // 通知地址
	TradeType   string  `json:"tradeType"`   // 交易类型 JSAPI、NATIVE、APP、MWEB
	OpenId      string  `json:"openId"`      // 用户OpenID(JSAPI必填)
}

// WxPayNotifyResult 微信支付回调结果
type WxPayNotifyResult struct {
	ReturnCode    string `xml:"return_code"`    // 返回状态码
	ReturnMsg     string `xml:"return_msg"`     // 返回信息
	AppId         string `xml:"appid"`          // 微信支付AppID
	MchId         string `xml:"mch_id"`         // 微信支付商户号
	DeviceInfo    string `xml:"device_info"`    // 设备号
	NonceStr      string `xml:"nonce_str"`      // 随机字符串
	Sign          string `xml:"sign"`           // 签名
	SignType      string `xml:"sign_type"`      // 签名类型
	ResultCode    string `xml:"result_code"`    // 业务结果
	ErrCode       string `xml:"err_code"`       // 错误代码
	ErrCodeDes    string `xml:"err_code_des"`   // 错误代码描述
	OpenId        string `xml:"openid"`         // 用户标识
	IsSubscribe   string `xml:"is_subscribe"`   // 是否关注公众号
	TradeType     string `xml:"trade_type"`     // 交易类型
	BankType      string `xml:"bank_type"`      // 付款银行
	TotalFee      int    `xml:"total_fee"`      // 订单金额(分)
	FeeType       string `xml:"fee_type"`       // 货币种类
	CashFee       int    `xml:"cash_fee"`       // 现金支付金额(分)
	CashFeeType   string `xml:"cash_fee_type"`  // 现金支付货币类型
	TransactionId string `xml:"transaction_id"` // 微信支付订单号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户订单号
	Attach        string `xml:"attach"`         // 商家数据包
	TimeEnd       string `xml:"time_end"`       // 支付完成时间
}

// WxPayQueryResult 微信支付查询结果
type WxPayQueryResult struct {
	ReturnCode     string `xml:"return_code"`      // 返回状态码
	ReturnMsg      string `xml:"return_msg"`       // 返回信息
	AppId          string `xml:"appid"`            // 微信支付AppID
	MchId          string `xml:"mch_id"`           // 微信支付商户号
	NonceStr       string `xml:"nonce_str"`        // 随机字符串
	Sign           string `xml:"sign"`             // 签名
	ResultCode     string `xml:"result_code"`      // 业务结果
	ErrCode        string `xml:"err_code"`         // 错误代码
	ErrCodeDes     string `xml:"err_code_des"`     // 错误代码描述
	DeviceInfo     string `xml:"device_info"`      // 设备号
	OpenId         string `xml:"openid"`           // 用户标识
	IsSubscribe    string `xml:"is_subscribe"`     // 是否关注公众号
	TradeType      string `xml:"trade_type"`       // 交易类型
	TradeState     string `xml:"trade_state"`      // 交易状态
	TradeStateDesc string `xml:"trade_state_desc"` // 交易状态描述
	BankType       string `xml:"bank_type"`        // 付款银行
	TotalFee       int    `xml:"total_fee"`        // 订单金额(分)
	FeeType        string `xml:"fee_type"`         // 货币种类
	CashFee        int    `xml:"cash_fee"`         // 现金支付金额(分)
	CashFeeType    string `xml:"cash_fee_type"`    // 现金支付货币类型
	TransactionId  string `xml:"transaction_id"`   // 微信支付订单号
	OutTradeNo     string `xml:"out_trade_no"`     // 商户订单号
	Attach         string `xml:"attach"`           // 商家数据包
	TimeEnd        string `xml:"time_end"`         // 支付完成时间
}
