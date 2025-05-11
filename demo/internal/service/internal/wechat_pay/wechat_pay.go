package wechat_pay

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	v1 "demo/api/payment/v1/wechat"
	"demo/internal/consts"
	"demo/internal/model"
	"demo/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/guid"
)

type sWechatPay struct{}

// New 创建微信支付服务实例
func New() service.WechatPayService {
	return &sWechatPay{}
}

// GetConfig 获取微信支付配置
func (s *sWechatPay) GetConfig(ctx context.Context) (*model.WxPayConfig, error) {
	// 从数据库或配置中获取微信支付配置
	configValue, err := g.Cfg().Get(ctx, "payment.wechat")
	if err != nil {
		g.Log().Error(ctx, "获取微信支付配置失败:", err)
		// 尝试从数据库获取
		var config struct {
			AppId     string `json:"appId"`
			MchId     string `json:"mchId"`
			ApiKey    string `json:"apiKey"`
			NotifyUrl string `json:"notifyUrl"`
			IsEnabled bool   `json:"isEnabled"`
		}
		err = g.DB().Model("payment_config").Where("1=1").Limit(1).Scan(&config)
		if err != nil {
			return nil, gerror.Wrap(err, "获取微信支付配置失败")
		}

		if config.AppId == "" || config.MchId == "" || config.ApiKey == "" {
			return nil, gerror.New("微信支付配置未设置")
		}

		if !config.IsEnabled {
			return nil, gerror.New("微信支付功能未启用")
		}

		return &model.WxPayConfig{
			AppId:     config.AppId,
			MchId:     config.MchId,
			ApiKey:    config.ApiKey,
			NotifyUrl: config.NotifyUrl,
		}, nil
	}

	var payConfig struct {
		AppId     string `json:"appId"`
		MchId     string `json:"mchId"`
		ApiKey    string `json:"apiKey"`
		NotifyUrl string `json:"notifyUrl"`
		IsEnabled bool   `json:"isEnabled"`
	}

	err = configValue.Scan(&payConfig)
	if err != nil {
		return nil, gerror.Wrap(err, "解析微信支付配置失败")
	}

	if payConfig.AppId == "" || payConfig.MchId == "" || payConfig.ApiKey == "" {
		return nil, gerror.New("微信支付配置未设置")
	}

	if !payConfig.IsEnabled {
		return nil, gerror.New("微信支付功能未启用")
	}

	return &model.WxPayConfig{
		AppId:     payConfig.AppId,
		MchId:     payConfig.MchId,
		ApiKey:    payConfig.ApiKey,
		NotifyUrl: payConfig.NotifyUrl,
	}, nil
}

// 生成随机字符串
func generateNonceStr() string {
	return guid.S()
}

// 生成签名
func generateSign(params map[string]string, apiKey string) string {
	// 按照参数名ASCII码从小到大排序
	var keys []string
	for k := range params {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接字符串
	var buf strings.Builder
	for _, k := range keys {
		if params[k] == "" {
			continue
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(params[k])
		buf.WriteString("&")
	}
	buf.WriteString("key=")
	buf.WriteString(apiKey)

	// MD5签名并转为大写
	h := md5.New()
	h.Write([]byte(buf.String()))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// 元转换为分
func yuan2Fen(yuan float64) int {
	return int(math.Round(yuan * 100))
}

// 分转换为元
func fen2Yuan(fen int) float64 {
	return float64(fen) / 100
}

// UnifiedOrder 统一下单
func (s *sWechatPay) UnifiedOrder(ctx context.Context, req *v1.WxPayUnifiedOrderReq) (res *v1.WxPayUnifiedOrderRes, err error) {
	// 获取微信支付配置
	config, err := s.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	// 查询订单是否存在
	var order model.Order
	err = g.DB().Model("order").Where("order_no = ?", req.OrderNo).Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 检查订单状态
	if order.Status != 0 {
		return nil, gerror.NewCode(gcode.New(400, "订单状态不正确，不能发起支付", nil))
	}

	// 转换金额为分
	totalFee := yuan2Fen(req.TotalFee)
	if totalFee <= 0 {
		return nil, gerror.NewCode(gcode.New(400, "支付金额必须大于0", nil))
	}

	// 生成随机字符串
	nonceStr := generateNonceStr()

	// 统一下单API地址
	apiUrl := "https://api.mch.weixin.qq.com/pay/unifiedorder"

	// 准备请求参数
	notifyUrl := config.NotifyUrl

	// 构建请求参数
	params := map[string]string{
		"appid":            config.AppId,
		"mch_id":           config.MchId,
		"nonce_str":        nonceStr,
		"body":             req.Body,
		"out_trade_no":     req.OrderNo,
		"total_fee":        strconv.Itoa(totalFee),
		"spbill_create_ip": g.RequestFromCtx(ctx).GetClientIp(),
		"notify_url":       notifyUrl,
		"trade_type":       "NATIVE", // 默认生成二维码支付
	}

	// 生成签名
	params["sign"] = generateSign(params, config.ApiKey)

	// 将参数转为XML
	var xmlBuf strings.Builder
	xmlBuf.WriteString("<xml>")
	for k, v := range params {
		xmlBuf.WriteString("<")
		xmlBuf.WriteString(k)
		xmlBuf.WriteString(">")
		xmlBuf.WriteString(v)
		xmlBuf.WriteString("</")
		xmlBuf.WriteString(k)
		xmlBuf.WriteString(">")
	}
	xmlBuf.WriteString("</xml>")

	// 发送请求
	client := g.Client()
	client.SetHeader("Content-Type", "application/xml")
	response, err := client.Post(ctx, apiUrl, xmlBuf.String())
	if err != nil {
		return nil, gerror.Wrap(err, "调用微信支付统一下单接口失败")
	}
	defer response.Close()

	// 解析响应XML
	xmlResp := response.ReadAllString()
	var wxResponse struct {
		ReturnCode string `xml:"return_code"`
		ReturnMsg  string `xml:"return_msg"`
		ResultCode string `xml:"result_code"`
		ErrCode    string `xml:"err_code"`
		ErrCodeDes string `xml:"err_code_des"`
		AppId      string `xml:"appid"`
		MchId      string `xml:"mch_id"`
		NonceStr   string `xml:"nonce_str"`
		Sign       string `xml:"sign"`
		PrepayId   string `xml:"prepay_id"`
		TradeType  string `xml:"trade_type"`
		CodeUrl    string `xml:"code_url"`
		MwebUrl    string `xml:"mweb_url"`
	}

	err = xml.Unmarshal([]byte(xmlResp), &wxResponse)
	if err != nil {
		return nil, gerror.Wrap(err, "解析微信支付统一下单响应失败")
	}

	// 检查返回结果
	if wxResponse.ReturnCode != "SUCCESS" {
		return nil, gerror.Newf("微信支付统一下单失败: %s", wxResponse.ReturnMsg)
	}

	if wxResponse.ResultCode != "SUCCESS" {
		return nil, gerror.Newf("微信支付统一下单业务失败: %s - %s", wxResponse.ErrCode, wxResponse.ErrCodeDes)
	}

	// 生成支付参数
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	signType := "MD5"

	// 准备返回结果
	res = &v1.WxPayUnifiedOrderRes{
		AppId:      config.AppId,
		TimeStamp:  timeStamp,
		NonceStr:   nonceStr,
		Package:    "prepay_id=" + wxResponse.PrepayId,
		SignType:   signType,
		TradeType:  wxResponse.TradeType,
		PrepayId:   wxResponse.PrepayId,
		OutTradeNo: req.OrderNo,
	}

	// 根据交易类型返回不同的结果
	if wxResponse.TradeType == "NATIVE" {
		res.CodeUrl = wxResponse.CodeUrl
	} else if wxResponse.TradeType == "MWEB" {
		res.MwebUrl = wxResponse.MwebUrl
	}

	// 计算支付签名
	paySignParams := map[string]string{
		"appId":     config.AppId,
		"timeStamp": timeStamp,
		"nonceStr":  nonceStr,
		"package":   "prepay_id=" + wxResponse.PrepayId,
		"signType":  signType,
	}
	res.PaySign = generateSign(paySignParams, config.ApiKey)

	// 记录日志
	logData, _ := json.Marshal(res)
	g.Log().Info(ctx, "微信支付统一下单成功:", string(logData))

	return res, nil
}

// HandleNotify 处理微信支付回调通知
func (s *sWechatPay) HandleNotify(ctx context.Context, notifyData []byte) (res *v1.WxPayNotifyRes, err error) {
	// 初始化响应
	res = &v1.WxPayNotifyRes{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
	}

	// 解析通知数据
	var notifyResult model.WxPayNotifyResult
	err = xml.Unmarshal(notifyData, &notifyResult)
	if err != nil {
		g.Log().Error(ctx, "解析微信支付回调通知失败:", err)
		res.ReturnCode = "FAIL"
		res.ReturnMsg = "解析通知数据失败"
		return res, err
	}

	// 验证返回结果
	if notifyResult.ReturnCode != "SUCCESS" {
		g.Log().Error(ctx, "微信支付回调通知失败:", notifyResult.ReturnMsg)
		return res, nil
	}

	if notifyResult.ResultCode != "SUCCESS" {
		g.Log().Error(ctx, "微信支付回调业务结果失败:", notifyResult.ErrCode, notifyResult.ErrCodeDes)
		return res, nil
	}

	// 获取订单号
	orderNo := notifyResult.OutTradeNo
	if orderNo == "" {
		g.Log().Error(ctx, "微信支付回调通知中没有订单号")
		res.ReturnCode = "FAIL"
		res.ReturnMsg = "订单号不能为空"
		return res, errors.New("订单号不能为空")
	}

	// 检查是否是临时订单号
	var orderPayMap struct {
		Id              int    `json:"id"`
		OriginalOrderNo string `json:"original_order_no"`
		TempOrderNo     string `json:"temp_order_no"`
		ClientId        int    `json:"client_id"`
		Status          int    `json:"status"`
	}

	err = g.DB().Model("order_pay_map").Where("temp_order_no", orderNo).Scan(&orderPayMap)
	if err != nil {
		g.Log().Error(ctx, "查询支付映射关系失败:", err)
		// 这里不要直接返回错误，因为可能是正常订单号
	}

	// 如果找到映射关系，使用原始订单号
	if orderPayMap.Id > 0 && orderPayMap.Status == 0 {
		g.Log().Info(ctx, "找到临时订单映射:", orderPayMap)

		// 更新映射状态为已处理
		_, err = g.DB().Model("order_pay_map").
			Where("id", orderPayMap.Id).
			Data(g.Map{
				"status":     1,
				"updated_at": gtime.Now(),
			}).
			Update()
		if err != nil {
			g.Log().Error(ctx, "更新支付映射状态失败:", err)
			// 继续执行，不影响订单状态更新
		}

		// 使用原始订单号
		orderNo = orderPayMap.OriginalOrderNo
	}

	// 获取支付金额（单位：分）
	amount := notifyResult.TotalFee

	// 将金额转换为元
	payAmount := float64(amount) / 100

	// 查询订单
	var order model.Order
	err = g.DB().Model("order").Where("order_no", orderNo).Scan(&order)
	if err != nil {
		g.Log().Error(ctx, "查询订单失败:", err)
		res.ReturnCode = "FAIL"
		res.ReturnMsg = "订单不存在"
		return res, err
	}

	if order.Id == 0 {
		g.Log().Error(ctx, "订单不存在:", orderNo)
		res.ReturnCode = "FAIL"
		res.ReturnMsg = "订单不存在"
		return res, errors.New("订单不存在")
	}

	// 验证支付金额
	if math.Abs(order.Amount-payAmount) > 0.01 {
		g.Log().Warning(ctx, "支付金额不匹配:", g.Map{
			"orderAmount": order.Amount,
			"payAmount":   payAmount,
		})
		// 记录日志，但继续处理订单
	}

	// 更新订单状态
	err = service.Order().UpdateOrderStatus(ctx, orderNo, consts.OrderStatusProcessing, notifyResult.TransactionId)
	if err != nil {
		g.Log().Error(ctx, "更新订单状态失败:", err)
		res.ReturnCode = "FAIL"
		res.ReturnMsg = "更新订单状态失败"
		return res, err
	}

	g.Log().Info(ctx, "订单支付成功处理完成:", orderNo)
	return res, nil
}

// QueryOrder 查询订单
func (s *sWechatPay) QueryOrder(ctx context.Context, req *v1.WxPayOrderQueryReq) (res *v1.WxPayOrderQueryRes, err error) {
	// 获取微信支付配置
	config, err := s.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	// 查询订单是否存在
	var order model.Order
	err = g.DB().Model("order").Where("order_no = ?", req.OrderNo).Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 生成随机字符串
	nonceStr := generateNonceStr()

	// 订单查询API地址
	apiUrl := "https://api.mch.weixin.qq.com/pay/orderquery"

	// 构建请求参数
	params := map[string]string{
		"appid":        config.AppId,
		"mch_id":       config.MchId,
		"out_trade_no": req.OrderNo,
		"nonce_str":    nonceStr,
	}

	// 生成签名
	params["sign"] = generateSign(params, config.ApiKey)

	// 将参数转为XML
	var xmlBuf strings.Builder
	xmlBuf.WriteString("<xml>")
	for k, v := range params {
		xmlBuf.WriteString("<")
		xmlBuf.WriteString(k)
		xmlBuf.WriteString(">")
		xmlBuf.WriteString(v)
		xmlBuf.WriteString("</")
		xmlBuf.WriteString(k)
		xmlBuf.WriteString(">")
	}
	xmlBuf.WriteString("</xml>")

	// 发送请求
	client := g.Client()
	client.SetHeader("Content-Type", "application/xml")
	response, err := client.Post(ctx, apiUrl, xmlBuf.String())
	if err != nil {
		return nil, gerror.Wrap(err, "调用微信支付订单查询接口失败")
	}
	defer response.Close()

	// 解析响应XML
	xmlResp := response.ReadAllString()
	var queryResult model.WxPayQueryResult
	err = xml.Unmarshal([]byte(xmlResp), &queryResult)
	if err != nil {
		return nil, gerror.Wrap(err, "解析微信支付订单查询响应失败")
	}

	// 检查返回结果
	if queryResult.ReturnCode != "SUCCESS" {
		return nil, gerror.Newf("微信支付订单查询失败: %s", queryResult.ReturnMsg)
	}

	if queryResult.ResultCode != "SUCCESS" {
		return nil, gerror.Newf("微信支付订单查询业务失败: %s - %s", queryResult.ErrCode, queryResult.ErrCodeDes)
	}

	// 准备返回结果
	res = &v1.WxPayOrderQueryRes{
		OrderNo:        req.OrderNo,
		TransactionId:  queryResult.TransactionId,
		TradeState:     queryResult.TradeState,
		TradeStateDesc: queryResult.TradeStateDesc,
		TotalFee:       fen2Yuan(queryResult.TotalFee),
	}

	// 格式化支付时间
	if queryResult.TimeEnd != "" {
		timeLayout := "20060102150405"
		payTime, err := time.Parse(timeLayout, queryResult.TimeEnd)
		if err == nil {
			res.PayTime = payTime.Format("2006-01-02 15:04:05")
		}
	}

	// 更新本地订单状态
	if queryResult.TradeState == "SUCCESS" && order.Status == 0 {
		// 支付成功，但本地订单状态未更新，更新订单状态
		err = service.Order().UpdateOrderStatus(ctx, req.OrderNo, consts.OrderStatusProcessing, queryResult.TransactionId)
		if err != nil {
			g.Log().Error(ctx, "更新订单状态失败:", err)
		}
	}

	return res, nil
}
