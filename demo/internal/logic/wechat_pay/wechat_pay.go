package wechat_pay

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	v1 "demo/api/payment/v1/wechat"
	"demo/internal/consts"
	"demo/internal/model"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
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
	// 从配置文件中获取微信支付配置
	configFilePath := "manifest/config/payment.json"

	// 检查文件是否存在
	if !g.Cfg().Available(ctx) {
		// 尝试从文件系统读取
		if !gfile.Exists(configFilePath) {
			g.Log().Error(ctx, "支付配置文件不存在:", configFilePath)
			return nil, gerror.New("微信支付配置文件不存在")
		}

		// 读取配置文件内容
		content := gfile.GetContents(configFilePath)
		if content == "" {
			g.Log().Error(ctx, "读取支付配置文件失败")
			return nil, gerror.New("读取支付配置文件失败")
		}

		// 解析JSON内容
		var payConfig struct {
			AppId     string `json:"appId"`
			MchId     string `json:"mchId"`
			ApiKey    string `json:"apiKey"`
			NotifyUrl string `json:"notifyUrl"`
			IsEnabled bool   `json:"isEnabled"`
		}

		if err := gjson.DecodeTo(content, &payConfig); err != nil {
			g.Log().Error(ctx, "解析支付配置文件失败:", err)
			return nil, gerror.Wrap(err, "解析支付配置文件失败")
		}

		// 验证配置是否完整
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

	// 如果配置系统可用，尝试从配置系统获取
	configValue, err := g.Cfg().Get(ctx, "payment.wechat")
	if err == nil && !configValue.IsEmpty() {
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

	// 最后尝试直接读取payment.json文件
	if gfile.Exists(configFilePath) {
		content := gfile.GetContents(configFilePath)
		if content == "" {
			g.Log().Error(ctx, "读取支付配置文件失败")
			return nil, gerror.New("读取支付配置文件失败")
		}

		var payConfig struct {
			AppId     string `json:"appId"`
			MchId     string `json:"mchId"`
			ApiKey    string `json:"apiKey"`
			NotifyUrl string `json:"notifyUrl"`
			IsEnabled bool   `json:"isEnabled"`
		}

		if err := gjson.DecodeTo(content, &payConfig); err != nil {
			g.Log().Error(ctx, "解析支付配置文件失败:", err)
			return nil, gerror.Wrap(err, "解析支付配置文件失败")
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

	return nil, gerror.New("无法获取微信支付配置")
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

	// 记录签名字符串（不显示实际ApiKey）
	signStr := buf.String()
	strParts := strings.Split(signStr, "key=")
	if len(strParts) > 0 {
		// 安全记录签名字符串，隐藏真实密钥
		safeStr := strParts[0] + "key=***"
		g.Log().Debug(context.Background(), "签名字符串(安全版):", safeStr)
	}

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

// 获取字符串之间的内容，并处理CDATA标记
func getBetween(str, start, end string) string {
	startIndex := strings.Index(str, start)
	if startIndex == -1 {
		return ""
	}
	startIndex += len(start)
	endIndex := strings.Index(str[startIndex:], end)
	if endIndex == -1 {
		return ""
	}

	// 提取标签之间的内容
	result := str[startIndex : startIndex+endIndex]

	// 处理CDATA标记
	if strings.Contains(result, "![CDATA[") && strings.Contains(result, "]]") {
		cdataStart := strings.Index(result, "![CDATA[")
		cdataEnd := strings.LastIndex(result, "]]")
		if cdataStart != -1 && cdataEnd != -1 && cdataEnd > cdataStart {
			// 提取CDATA内部的实际内容（加8是因为"![CDATA["的长度是8）
			result = result[cdataStart+8 : cdataEnd]
		}
	}

	return result
}

// UnifiedOrder 统一下单
func (s *sWechatPay) UnifiedOrder(ctx context.Context, req *v1.WxPayUnifiedOrderReq) (res *v1.WxPayUnifiedOrderRes, err error) {
	// 验证客户权限（不再验证管理员权限）
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取客户信息失败:", err)
		return nil, gerror.NewCode(gcode.New(403, "未登录或登录已过期", nil))
	}

	// 获取微信支付配置
	config, err := s.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	// 确定订单号 - 处理可能的临时订单号情况
	orderNo := req.OrderNo
	originalOrderNo := ""

	// 检查是否是临时订单号格式（带有额外的随机数）
	if len(orderNo) > 6 && strings.Contains(orderNo, "ORD") {
		// 检查order_pay_map表，看是否是临时订单号
		var orderMap struct {
			OriginalOrderNo string `json:"original_order_no"`
		}
		err = g.DB().Model("order_pay_map").
			Where("temp_order_no = ?", orderNo).
			Scan(&orderMap)

		if err == nil && orderMap.OriginalOrderNo != "" {
			// 找到了对应的原始订单号
			g.Log().Info(ctx, fmt.Sprintf("找到临时订单号 %s 对应的原始订单号 %s", orderNo, orderMap.OriginalOrderNo))
			originalOrderNo = orderMap.OriginalOrderNo
			// 直接用临时订单号与微信支付系统通信，但查询原始订单
			orderNo = originalOrderNo
		}
	}

	// 查询订单是否存在
	var order model.Order
	err = g.DB().Model("order").Where("order_no = ?", orderNo).Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 检查权限，只有订单所有者可以支付
	if order.ClientId != clientId {
		g.Log().Error(ctx, fmt.Sprintf("客户 %d 尝试支付非自己的订单 %s", clientId, orderNo))
		return nil, gerror.NewCode(gcode.New(403, "您没有权限支付此订单", nil))
	}

	// 检查订单状态
	if order.Status != 0 {
		return nil, gerror.NewCode(gcode.New(400, "订单状态不正确，不能发起支付", nil))
	}

	// 检查订单是否已过期
	now := gtime.Now()
	if order.ExpireTime != nil && now.TimestampMilli() > order.ExpireTime.TimestampMilli() {
		g.Log().Warning(ctx, fmt.Sprintf("客户 %d 尝试支付已过期的订单 %s", clientId, orderNo))
		return nil, gerror.NewCode(gcode.New(400, "订单已过期，请重新创建订单", nil))
	}

	// 检查支付金额是否一致
	if math.Abs(order.Amount-req.TotalFee) > 0.01 {
		g.Log().Warning(ctx, fmt.Sprintf("支付金额不一致: 订单金额=%.2f, 请求金额=%.2f", order.Amount, req.TotalFee))
		return nil, gerror.NewCode(gcode.New(400, "支付金额不正确", nil))
	}

	// 记录日志
	defer func() {
		if err != nil {
			service.Log().Record(ctx, clientId, "", "微信支付", "发起支付", 0, err.Error())
		} else {
			service.Log().Record(ctx, clientId, "", "微信支付", "发起支付", 1, "")
		}
	}()

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

	// 判断是否来自小程序的请求
	tradeType := "NATIVE" // 默认生成二维码支付
	openId := ""

	// 从请求上下文中获取请求来源和OpenID
	requestFrom := auth.GetRequestFrom(ctx)
	if requestFrom == "wxapp" {
		// 来自小程序的请求，使用JSAPI支付
		tradeType = "JSAPI"
		openId = auth.GetOpenID(ctx)
		if openId == "" {
			return nil, gerror.NewCode(gcode.New(401, "无法获取OpenID，请重新登录", nil))
		}
	}

	// 确定使用哪个订单号与微信支付系统通信
	outTradeNo := req.OrderNo // 默认使用传入的订单号

	// 构建请求参数
	params := map[string]string{
		"appid":            config.AppId,
		"mch_id":           config.MchId,
		"nonce_str":        nonceStr,
		"body":             req.Body,
		"out_trade_no":     outTradeNo,
		"total_fee":        strconv.Itoa(totalFee),
		"spbill_create_ip": g.RequestFromCtx(ctx).GetClientIp(),
		"notify_url":       notifyUrl,
		"trade_type":       tradeType,
	}

	// 如果是JSAPI支付，添加openid参数
	if tradeType == "JSAPI" {
		params["openid"] = openId
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
	// 使用原始订单号返回给前端
	orderNoForResponse := originalOrderNo
	if orderNoForResponse == "" {
		orderNoForResponse = orderNo // 如果不是临时订单号情况，直接返回原始订单号
	}

	res = &v1.WxPayUnifiedOrderRes{
		AppId:      config.AppId,
		TimeStamp:  timeStamp,
		NonceStr:   nonceStr,
		Package:    "prepay_id=" + wxResponse.PrepayId,
		SignType:   signType,
		TradeType:  wxResponse.TradeType,
		PrepayId:   wxResponse.PrepayId,
		OutTradeNo: orderNoForResponse,
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

// HandleNotify 处理支付回调通知
func (s *sWechatPay) HandleNotify(ctx context.Context, notifyData []byte) (res *v1.WxPayNotifyRes, err error) {
	// 获取微信支付配置
	config, err := s.GetConfig(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取微信支付配置失败:", err)
		return &v1.WxPayNotifyRes{
			ReturnCode: "FAIL",
			ReturnMsg:  "支付配置获取失败",
		}, err
	}

	// 解析通知数据
	var notifyResult model.WxPayNotifyResult
	err = xml.Unmarshal(notifyData, &notifyResult)
	if err != nil {
		g.Log().Error(ctx, "解析微信支付通知数据失败:", err)
		return &v1.WxPayNotifyRes{
			ReturnCode: "FAIL",
			ReturnMsg:  "数据格式错误",
		}, err
	}

	// 记录通知数据
	logData, _ := json.Marshal(notifyResult)
	g.Log().Info(ctx, "收到微信支付通知:", string(logData))

	// 检查XML数据中可能影响签名计算的特殊格式
	notifyXmlStr := string(notifyData)
	if strings.Contains(notifyXmlStr, "![CDATA[") {
		g.Log().Warning(ctx, "微信支付通知XML中包含CDATA标记，可能影响签名验证")
	}
	if strings.Contains(notifyXmlStr, "\r") || strings.Contains(notifyXmlStr, "\n") {
		g.Log().Warning(ctx, "微信支付通知XML中包含换行符，可能影响签名验证")
	}

	// 验证返回状态
	if notifyResult.ReturnCode != "SUCCESS" {
		g.Log().Error(ctx, "微信支付通知返回失败:", notifyResult.ReturnMsg)
		return &v1.WxPayNotifyRes{
			ReturnCode: "FAIL",
			ReturnMsg:  "通信失败",
		}, errors.New(notifyResult.ReturnMsg)
	}

	// 验证业务结果
	if notifyResult.ResultCode != "SUCCESS" {
		g.Log().Error(ctx, "微信支付结果失败:", notifyResult.ErrCode, notifyResult.ErrCodeDes)
		return &v1.WxPayNotifyRes{
			ReturnCode: "FAIL",
			ReturnMsg:  "业务失败",
		}, errors.New(notifyResult.ErrCodeDes)
	}

	// 验证签名（构建签名参数）
	signParams := make(map[string]string)
	notifyXml := string(notifyData)

	// 先提取可能不包含CDATA的简单字段
	fields := []struct {
		key   string
		start string
		end   string
	}{
		{"appid", "<appid>", "</appid>"},
		{"bank_type", "<bank_type>", "</bank_type>"},
		{"cash_fee", "<cash_fee>", "</cash_fee>"},
		{"fee_type", "<fee_type>", "</fee_type>"},
		{"is_subscribe", "<is_subscribe>", "</is_subscribe>"},
		{"mch_id", "<mch_id>", "</mch_id>"},
		{"nonce_str", "<nonce_str>", "</nonce_str>"},
		{"openid", "<openid>", "</openid>"},
		{"out_trade_no", "<out_trade_no>", "</out_trade_no>"},
		{"result_code", "<result_code>", "</result_code>"},
		{"return_code", "<return_code>", "</return_code>"},
		{"time_end", "<time_end>", "</time_end>"},
		{"total_fee", "<total_fee>", "</total_fee>"},
		{"trade_type", "<trade_type>", "</trade_type>"},
		{"transaction_id", "<transaction_id>", "</transaction_id>"},
	}

	for _, field := range fields {
		value := getBetween(notifyXml, field.start, field.end)
		if value != "" {
			signParams[field.key] = value
		}
	}

	// 确保total_fee字段是数字
	if _, ok := signParams["total_fee"]; !ok || signParams["total_fee"] == "" {
		// 尝试直接提取不带CDATA的数值
		start := "<total_fee>"
		end := "</total_fee>"
		totalFeeStart := strings.Index(notifyXml, start)
		if totalFeeStart != -1 {
			totalFeeStart += len(start)
			totalFeeEnd := strings.Index(notifyXml[totalFeeStart:], end)
			if totalFeeEnd != -1 {
				signParams["total_fee"] = notifyXml[totalFeeStart : totalFeeStart+totalFeeEnd]
			}
		}
	}

	// 记录日志，帮助排查签名验证问题
	g.Log().Debug(ctx, "微信支付配置信息 - AppId:", config.AppId, "MchId:", config.MchId)
	paramsJSON, _ := json.Marshal(signParams)
	g.Log().Debug(ctx, "微信支付签名参数(CDATA已处理):", string(paramsJSON))

	// 验证签名
	sign := generateSign(signParams, config.ApiKey)
	g.Log().Debug(ctx, "计算得到的签名:", sign)
	g.Log().Debug(ctx, "通知中的签名:", notifyResult.Sign)

	if sign != notifyResult.Sign {
		g.Log().Error(ctx, "微信支付通知签名验证失败")
		return &v1.WxPayNotifyRes{
			ReturnCode: "FAIL",
			ReturnMsg:  "签名验证失败",
		}, errors.New("签名验证失败")
	}

	// 获取订单号并检查是否是临时订单号
	tempOrderNo := notifyResult.OutTradeNo
	orderNo := tempOrderNo

	// 如果是临时订单号，查询原始订单号
	if len(tempOrderNo) > 6 && strings.Contains(tempOrderNo, "ORD") {
		var orderMap struct {
			OriginalOrderNo string `json:"original_order_no"`
		}
		err = g.DB().Model("order_pay_map").
			Where("temp_order_no = ?", tempOrderNo).
			Scan(&orderMap)

		if err == nil && orderMap.OriginalOrderNo != "" {
			// 找到了对应的原始订单号
			g.Log().Info(ctx, fmt.Sprintf("回调通知: 找到临时订单号 %s 对应的原始订单号 %s", tempOrderNo, orderMap.OriginalOrderNo))
			orderNo = orderMap.OriginalOrderNo
		}
	}

	// 检查订单金额是否正确
	var order model.Order
	err = g.DB().Model("order").Where("order_no = ?", orderNo).Scan(&order)
	if err != nil {
		g.Log().Error(ctx, "查询订单失败:", err)
		return &v1.WxPayNotifyRes{
			ReturnCode: "FAIL",
			ReturnMsg:  "订单查询失败",
		}, err
	}

	if order.Id == 0 {
		g.Log().Error(ctx, "订单不存在:", orderNo)
		return &v1.WxPayNotifyRes{
			ReturnCode: "FAIL",
			ReturnMsg:  "订单不存在",
		}, errors.New("订单不存在")
	}

	// 验证订单金额（将元转为分进行比较）
	orderFee := yuan2Fen(order.Amount)
	if orderFee != notifyResult.TotalFee {
		g.Log().Error(ctx, fmt.Sprintf("订单金额不匹配: 订单金额=%d, 通知金额=%d", orderFee, notifyResult.TotalFee))
		return &v1.WxPayNotifyRes{
			ReturnCode: "FAIL",
			ReturnMsg:  "订单金额不匹配",
		}, errors.New("订单金额不匹配")
	}

	// 更新订单状态
	err = service.Order().UpdateOrderStatus(ctx, orderNo, consts.OrderStatusPaid, notifyResult.TransactionId)
	if err != nil {
		g.Log().Error(ctx, "更新订单状态失败:", err)
		return &v1.WxPayNotifyRes{
			ReturnCode: "FAIL",
			ReturnMsg:  "处理订单失败",
		}, err
	}

	// 返回成功
	return &v1.WxPayNotifyRes{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
	}, nil
}

// QueryOrder 查询订单
func (s *sWechatPay) QueryOrder(ctx context.Context, req *v1.WxPayOrderQueryReq) (res *v1.WxPayOrderQueryRes, err error) {
	// 验证客户权限（不再验证管理员权限）
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取客户信息失败:", err)
		return nil, gerror.NewCode(gcode.New(403, "未登录或登录已过期", nil))
	}

	// 获取微信支付配置
	config, err := s.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	// 获取订单号并检查是否是临时订单号
	tempOrderNo := req.OrderNo
	orderNo := tempOrderNo

	// 如果是临时订单号，查询原始订单号
	if len(tempOrderNo) > 6 && strings.Contains(tempOrderNo, "ORD") {
		var orderMap struct {
			OriginalOrderNo string `json:"original_order_no"`
		}
		err = g.DB().Model("order_pay_map").
			Where("temp_order_no = ?", tempOrderNo).
			Scan(&orderMap)

		if err == nil && orderMap.OriginalOrderNo != "" {
			// 找到了对应的原始订单号
			g.Log().Info(ctx, fmt.Sprintf("订单查询: 找到临时订单号 %s 对应的原始订单号 %s", tempOrderNo, orderMap.OriginalOrderNo))
			orderNo = orderMap.OriginalOrderNo
		}
	}

	// 查询订单是否存在
	var order model.Order
	err = g.DB().Model("order").Where("order_no = ?", orderNo).Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 检查权限，只有订单所有者可以查询
	if order.ClientId != clientId {
		g.Log().Error(ctx, fmt.Sprintf("客户 %d 尝试查询非自己的订单 %s", clientId, orderNo))
		return nil, gerror.NewCode(gcode.New(403, "您没有权限查询此订单", nil))
	}

	// 记录日志
	defer func() {
		if err != nil {
			service.Log().Record(ctx, clientId, "", "微信支付", "查询订单", 0, err.Error())
		} else {
			service.Log().Record(ctx, clientId, "", "微信支付", "查询订单", 1, "")
		}
	}()

	// 生成随机字符串
	nonceStr := generateNonceStr()

	// 订单查询API地址
	apiUrl := "https://api.mch.weixin.qq.com/pay/orderquery"

	// 构建请求参数
	params := map[string]string{
		"appid":        config.AppId,
		"mch_id":       config.MchId,
		"out_trade_no": tempOrderNo, // 使用临时订单号查询微信支付结果
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
	var wxResponse struct {
		ReturnCode     string `xml:"return_code"`
		ReturnMsg      string `xml:"return_msg"`
		ResultCode     string `xml:"result_code"`
		ErrCode        string `xml:"err_code"`
		ErrCodeDes     string `xml:"err_code_des"`
		TradeState     string `xml:"trade_state"`
		TradeStateDesc string `xml:"trade_state_desc"`
		TransactionId  string `xml:"transaction_id"`
		TimeEnd        string `xml:"time_end"`
		TotalFee       string `xml:"total_fee"`
	}

	err = xml.Unmarshal([]byte(xmlResp), &wxResponse)
	if err != nil {
		return nil, gerror.Wrap(err, "解析微信支付订单查询响应失败")
	}

	// 检查返回结果
	if wxResponse.ReturnCode != "SUCCESS" {
		return nil, gerror.Newf("微信支付订单查询失败: %s", wxResponse.ReturnMsg)
	}

	// 准备返回结果
	res = &v1.WxPayOrderQueryRes{
		OrderNo:        orderNo, // 返回原始订单号
		TransactionId:  wxResponse.TransactionId,
		TradeState:     wxResponse.TradeState,
		TradeStateDesc: wxResponse.TradeStateDesc,
		PayTime:        wxResponse.TimeEnd,
	}

	// 解析金额
	if wxResponse.TotalFee != "" {
		totalFee, err := strconv.Atoi(wxResponse.TotalFee)
		if err == nil {
			res.TotalFee = fen2Yuan(totalFee)
		}
	}

	// 检查支付状态
	if wxResponse.TradeState == "SUCCESS" {
		// 查询成功并且支付成功，自动更新订单状态（仅当订单状态为待支付时）
		if order.Status == 0 {
			err = service.Order().UpdateOrderStatus(ctx, orderNo, consts.OrderStatusProcessing, wxResponse.TransactionId)
			if err != nil {
				g.Log().Error(ctx, "自动更新订单状态失败:", err)
				// 不返回错误，仍然返回查询结果
			}
		}
	}

	return res, nil
}
