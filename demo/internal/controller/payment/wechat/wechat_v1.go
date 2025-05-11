package wechat

import (
	"io/ioutil"

	"demo/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ControllerV1 微信支付控制器V1版本
type ControllerV1 struct{}

// Notify 微信支付回调通知处理
func (c *ControllerV1) Notify(req *ghttp.Request) {
	ctx := req.Context()
	g.Log().Info(ctx, "收到微信支付回调请求")

	// 读取请求体
	notifyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		g.Log().Error(ctx, "读取微信支付回调请求体失败:", err)
		req.Response.WriteXmlExit(g.Map{
			"return_code": "FAIL",
			"return_msg":  "读取请求失败",
		})
		return
	}

	// 记录原始回调数据，便于排查问题
	g.Log().Info(ctx, "微信支付回调原始数据:", string(notifyData))

	// 处理回调通知
	res, err := service.WechatPay().HandleNotify(ctx, notifyData)
	if err != nil {
		g.Log().Error(ctx, "处理微信支付回调失败:", err)
	}

	// 返回结果
	g.Log().Info(ctx, "微信支付回调响应:", res.ReturnCode, res.ReturnMsg)
	req.Response.WriteXmlExit(g.Map{
		"return_code": res.ReturnCode,
		"return_msg":  res.ReturnMsg,
	})
}

// NewV1 创建微信支付控制器V1版本实例
func NewV1() *ControllerV1 {
	return &ControllerV1{}
}
