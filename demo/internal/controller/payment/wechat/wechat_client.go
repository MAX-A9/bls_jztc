package wechat

import (
	"context"

	clientapi "demo/api/payment/client"
	v1 "demo/api/payment/v1/wechat"
	"demo/internal/model"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// ClientControllerV1 微信支付客户端控制器V1版本
type ClientControllerV1 struct{}

// UnifiedOrder 微信支付统一下单
func (c *ClientControllerV1) UnifiedOrder(ctx context.Context, req *clientapi.WxPayUnifiedOrderReq) (res *clientapi.WxPayUnifiedOrderRes, err error) {
	// 检查客户权限
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(403, "未登录或登录已过期", nil))
	}

	// 验证订单所属权限
	var order model.Order
	err = g.DB().Model("order").Where("order_no = ?", req.OrderNo).Scan(&order)
	if err != nil {
		g.Log().Error(ctx, "查询订单失败:", err)
		return nil, gerror.Wrap(err, "查询订单失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 验证订单所属客户
	if order.ClientId != clientId {
		g.Log().Error(ctx, "客户尝试支付非自己的订单:", req.OrderNo)
		return nil, gerror.NewCode(gcode.New(403, "无权操作此订单", nil))
	}

	// 调用微信支付服务
	serviceReq := &v1.WxPayUnifiedOrderReq{
		OrderNo:  req.OrderNo,
		TotalFee: req.TotalFee,
		Body:     req.Body,
	}
	serviceRes, err := service.WechatPay().UnifiedOrder(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	res = &clientapi.WxPayUnifiedOrderRes{
		AppId:      serviceRes.AppId,
		TimeStamp:  serviceRes.TimeStamp,
		NonceStr:   serviceRes.NonceStr,
		Package:    serviceRes.Package,
		SignType:   serviceRes.SignType,
		PaySign:    serviceRes.PaySign,
		PrepayId:   serviceRes.PrepayId,
		CodeUrl:    serviceRes.CodeUrl,
		MwebUrl:    serviceRes.MwebUrl,
		TradeType:  serviceRes.TradeType,
		OutTradeNo: serviceRes.OutTradeNo,
	}

	return res, nil
}

// QueryOrder 微信支付订单查询
func (c *ClientControllerV1) QueryOrder(ctx context.Context, req *clientapi.WxPayOrderQueryReq) (res *clientapi.WxPayOrderQueryRes, err error) {
	// 检查客户权限
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(403, "未登录或登录已过期", nil))
	}

	// 验证订单所属权限
	var order model.Order
	err = g.DB().Model("order").Where("order_no = ?", req.OrderNo).Scan(&order)
	if err != nil {
		g.Log().Error(ctx, "查询订单失败:", err)
		return nil, gerror.Wrap(err, "查询订单失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 验证订单所属客户
	if order.ClientId != clientId {
		g.Log().Error(ctx, "客户尝试查询非自己的订单:", req.OrderNo)
		return nil, gerror.NewCode(gcode.New(403, "无权查询此订单", nil))
	}

	// 调用微信支付服务
	serviceReq := &v1.WxPayOrderQueryReq{
		OrderNo: req.OrderNo,
	}
	serviceRes, err := service.WechatPay().QueryOrder(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	res = &clientapi.WxPayOrderQueryRes{
		OrderNo:        serviceRes.OrderNo,
		TransactionId:  serviceRes.TransactionId,
		TradeState:     serviceRes.TradeState,
		TradeStateDesc: serviceRes.TradeStateDesc,
		PayTime:        serviceRes.PayTime,
		TotalFee:       serviceRes.TotalFee,
	}

	return res, nil
}

// NewClientV1 创建微信支付客户端控制器V1版本实例
func NewClientV1() *ClientControllerV1 {
	return &ClientControllerV1{}
}
