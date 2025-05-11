package client

import (
	"context"
	"fmt"
	"math"
	"time"

	clientapi "demo/api/payment/client"
	v1 "demo/api/payment/v1/wechat"
	"demo/internal/consts"
	"demo/internal/model"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
)

// 元转换为分
func yuan2Fen(yuan float64) int {
	return int(math.Round(yuan * 100))
}

// OrderController 客户端订单控制器
type OrderController struct{}

// List 获取客户订单列表
func (c *OrderController) List(ctx context.Context, req *clientapi.OrderListReq) (res *clientapi.OrderListRes, err error) {
	// 获取当前登录的客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}

	// 初始化响应
	res = &clientapi.OrderListRes{}

	// 构建查询条件
	m := g.DB().Model("order").Safe().Where("client_id", clientId)

	// 根据状态筛选
	switch req.Status {
	case "process": // 进行中
		m = m.Where("status", consts.OrderStatusProcessing)
	case "unpaid": // 待支付
		m = m.Where("status", consts.OrderStatusWaitPay)
	case "completed": // 已完成
		m = m.Where("status", consts.OrderStatusCompleted)
	case "cancelled": // 已取消
		m = m.Where("status", consts.OrderStatusCancelled)
	case "refunded": // 已退款
		m = m.Where("status", consts.OrderStatusRefunded)
	case "all": // 全部
		// 不添加状态过滤
	default:
		// 默认显示全部
	}

	// 分页查询
	res.Page = req.Page
	if res.Page <= 0 {
		res.Page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 获取总数
	res.Total, err = m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "获取订单总数失败")
	}

	// 没有数据，直接返回空列表
	if res.Total == 0 {
		res.List = make([]clientapi.OrderListItem, 0)
		return res, nil
	}

	// 查询数据
	var orders []model.Order
	err = m.Page(req.Page, pageSize).Order("id DESC").Scan(&orders)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单列表失败")
	}

	// 检查列表中的待支付订单是否已过期
	var ordersChanged bool // 标记是否有订单状态被改变
	for i := range orders {
		if orders[i].Status == consts.OrderStatusWaitPay {
			err = c.checkOrderStatus(ctx, &orders[i])
			if err != nil {
				g.Log().Error(ctx, "处理过期订单失败:", err)
			} else {
				// 如果订单状态可能已改变，标记需要重新获取
				if orders[i].ExpireTime != nil && gtime.Now().TimestampMilli() > orders[i].ExpireTime.TimestampMilli() {
					ordersChanged = true
				}
			}
		}
	}

	// 如果有订单状态变化，重新查询数据
	if ordersChanged {
		g.Log().Info(ctx, "检测到有订单状态变化，重新获取订单列表")
		err = m.Page(req.Page, pageSize).Order("id DESC").Scan(&orders)
		if err != nil {
			return nil, gerror.Wrap(err, "重新查询订单列表失败")
		}
	}

	// 转换为API响应格式
	res.List = make([]clientapi.OrderListItem, len(orders))
	for i, order := range orders {
		createdAt := ""
		if order.CreatedAt != nil {
			createdAt = order.CreatedAt.String()
		}

		payTime := ""
		if order.PayTime != nil {
			payTime = order.PayTime.String()
		}

		expireTime := ""
		if order.ExpireTime != nil {
			expireTime = order.ExpireTime.String()
		}

		// 获取状态文本
		statusText := getStatusText(order.Status)

		res.List[i] = clientapi.OrderListItem{
			Id:            order.Id,
			OrderNo:       order.OrderNo,
			ContentId:     order.ContentId,
			ProductName:   order.ProductName,
			Amount:        order.Amount,
			Status:        order.Status,
			StatusText:    statusText,
			PaymentMethod: order.PaymentMethod,
			CreatedAt:     createdAt,
			PayTime:       payTime,
			PackageInfo:   order.PackageInfo,
			ExpireTime:    expireTime,
		}
	}

	return res, nil
}

// Detail 获取客户订单详情
func (c *OrderController) Detail(ctx context.Context, req *clientapi.OrderDetailReq) (res *clientapi.OrderDetailRes, err error) {
	// 获取当前登录的客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}

	// 查询订单
	var order model.Order
	err = g.DB().Model("order").
		Where("order_no = ?", req.OrderNo).
		Where("client_id = ?", clientId). // 确保只能查询自己的订单
		Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单详情失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 检查订单状态，如已过期则自动处理
	err = c.checkOrderStatus(ctx, &order)
	if err != nil {
		g.Log().Error(ctx, "处理过期订单失败:", err)
		// 不中断流程，继续返回订单信息
	}

	// 如果刚才自动处理过，需要重新获取最新订单信息
	if order.Status == consts.OrderStatusWaitPay && order.ExpireTime != nil && gtime.Now().TimestampMilli() > order.ExpireTime.TimestampMilli() {
		err = g.DB().Model("order").
			Where("order_no = ?", req.OrderNo).
			Where("client_id = ?", clientId).
			Scan(&order)
		if err != nil {
			return nil, gerror.Wrap(err, "重新查询订单详情失败")
		}
	}

	// 获取状态文本
	statusText := getStatusText(order.Status)

	// 转换为API响应格式
	res = &clientapi.OrderDetailRes{
		Id:            order.Id,
		OrderNo:       order.OrderNo,
		ProductName:   order.ProductName,
		Amount:        order.Amount,
		Status:        order.Status,
		StatusText:    statusText,
		PaymentMethod: order.PaymentMethod,
		TransactionId: order.TransactionId,
		Remark:        order.Remark,
		PackageInfo:   order.PackageInfo,
	}

	if order.CreatedAt != nil {
		res.CreatedAt = order.CreatedAt.String()
	}

	if order.PayTime != nil {
		res.PayTime = order.PayTime.String()
	}

	if order.ExpireTime != nil {
		res.ExpireTime = order.ExpireTime.String()
	}

	return res, nil
}

// Pay 支付订单
func (c *OrderController) Pay(ctx context.Context, req *clientapi.OrderPayReq) (res *clientapi.OrderPayRes, err error) {
	// 获取当前登录的客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}

	// 查询订单
	var order model.Order
	err = g.DB().Model("order").
		Where("order_no = ?", req.OrderNo).
		Where("client_id = ?", clientId). // 确保只能支付自己的订单
		Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 检查订单状态
	if order.Status != 0 {
		return nil, gerror.NewCode(gcode.New(400, "只有待支付的订单才能发起支付", nil))
	}

	// 检查订单是否已过期并自动处理
	err = c.checkOrderStatus(ctx, &order)
	if err != nil {
		g.Log().Error(ctx, "处理过期订单失败:", err)
	}

	// 如果订单在上面的检查中被标记为过期，则直接返回错误
	if order.Status == consts.OrderStatusWaitPay && order.ExpireTime != nil && gtime.Now().TimestampMilli() > order.ExpireTime.TimestampMilli() {
		return nil, gerror.NewCode(gcode.New(400, "订单已过期，请重新创建订单", nil))
	}

	// 为避免订单号重复的问题，创建新的临时订单号
	// 原订单号格式: ORD + 时间 + 客户ID
	// 新订单号格式: ORD + 时间 + 客户ID + 6位随机数
	newOrderNo := fmt.Sprintf("%s%06d", order.OrderNo, grand.N(100000, 999999))

	// 记录订单和临时订单号的关系，用于支付成功后的回调处理
	// 注意: 真实环境应该在数据库中添加一个临时订单号字段，或创建订单支付记录表
	_, err = g.DB().Model("order_pay_map").Insert(g.Map{
		"original_order_no": order.OrderNo,
		"temp_order_no":     newOrderNo,
		"client_id":         clientId,
		"created_at":        time.Now(),
	})
	if err != nil {
		return nil, gerror.Wrap(err, "创建支付映射失败")
	}

	// 构建微信支付请求
	prepayReq := &v1.WxPayUnifiedOrderReq{
		OrderNo:  newOrderNo, // 使用临时订单号
		TotalFee: order.Amount,
		Body:     order.ProductName,
	}

	// 调用微信支付统一下单接口
	prepayRes, err := service.WechatPay().UnifiedOrder(ctx, prepayReq)
	if err != nil {
		return nil, gerror.Wrap(err, "创建支付订单失败")
	}

	// 转换响应
	res = &clientapi.OrderPayRes{
		OrderNo:    order.OrderNo,
		PrepayId:   prepayRes.PrepayId,
		NonceStr:   prepayRes.NonceStr,
		TimeStamp:  prepayRes.TimeStamp,
		Package:    prepayRes.Package,
		SignType:   prepayRes.SignType,
		PaySign:    prepayRes.PaySign,
		TotalFee:   yuan2Fen(order.Amount), // 将元转换为分
		OutTradeNo: prepayRes.OutTradeNo,
	}

	return res, nil
}

// handleExpiredOrder 处理过期订单
func (c *OrderController) handleExpiredOrder(ctx context.Context, order *model.Order) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 更新订单状态为已取消
		_, err := tx.Model("order").Where("id", order.Id).Update(g.Map{
			"status":     2, // 已取消
			"updated_at": gtime.Now(),
			"remark":     "订单已过期自动取消",
		})
		if err != nil {
			return gerror.Wrap(err, "更新订单状态失败")
		}

		// 检查是否有关联内容 - 先查询是否有content_id
		var contentInfo struct {
			ContentId int `json:"content_id"`
		}
		err = tx.Model("order").Where("id", order.Id).Fields("content_id").Scan(&contentInfo)
		if err != nil {
			return gerror.Wrap(err, "查询订单内容关联失败")
		}

		// 检查是否有关联内容
		if contentInfo.ContentId > 0 {
			// 查询内容状态
			var content struct {
				Id     int    `json:"id"`
				Status string `json:"status"`
			}
			err = tx.Model("content").Where("id", contentInfo.ContentId).Scan(&content)
			if err != nil {
				return gerror.Wrap(err, "查询内容信息失败")
			}

			// 只处理待审核状态的内容，删除而不是更新状态
			if content.Id > 0 && content.Status == "待审核" {
				// 直接删除关联内容
				_, err = tx.Model("content").Where("id", contentInfo.ContentId).Delete()
				if err != nil {
					return gerror.Wrap(err, "删除关联内容失败")
				}
				g.Log().Info(ctx, fmt.Sprintf("已删除过期订单 %s 关联的内容 ID: %d", order.OrderNo, contentInfo.ContentId))
			}
		}

		g.Log().Info(ctx, fmt.Sprintf("订单 %s 已过期自动取消", order.OrderNo))
		return nil
	})
}

// Cancel 客户取消订单
func (c *OrderController) Cancel(ctx context.Context, req *clientapi.OrderCancelReq) (res *clientapi.OrderCancelRes, err error) {
	// 获取当前登录的客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}

	// 查询订单
	var order model.Order
	err = g.DB().Model("order").
		Where("order_no = ?", req.OrderNo).
		Where("client_id = ?", clientId). // 确保只能取消自己的订单
		Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 检查订单状态，只能取消待支付的订单
	if order.Status != consts.OrderStatusWaitPay {
		return nil, gerror.NewCode(gcode.New(400, "只有待支付的订单才能取消", nil))
	}

	// 初始化响应
	res = &clientapi.OrderCancelRes{
		OrderNo:   order.OrderNo,
		Cancelled: false,
	}

	// 执行取消订单操作
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 更新订单状态为已取消
		remark := "客户主动取消订单"
		if req.Reason != "" {
			remark = "客户取消原因: " + req.Reason
		}

		_, err := tx.Model("order").Where("id", order.Id).Update(g.Map{
			"status":     consts.OrderStatusCancelled,
			"updated_at": gtime.Now(),
			"remark":     remark,
		})
		if err != nil {
			return gerror.Wrap(err, "更新订单状态失败")
		}

		// 检查是否有关联内容 - 先查询是否有content_id
		var contentInfo struct {
			ContentId int `json:"content_id"`
		}
		err = tx.Model("order").Where("id", order.Id).Fields("content_id").Scan(&contentInfo)
		if err != nil {
			return gerror.Wrap(err, "查询订单内容关联失败")
		}

		// 检查是否有关联内容
		if contentInfo.ContentId > 0 {
			// 查询内容状态
			var content struct {
				Id     int    `json:"id"`
				Status string `json:"status"`
			}
			err = tx.Model("content").Where("id", contentInfo.ContentId).Scan(&content)
			if err != nil {
				return gerror.Wrap(err, "查询内容信息失败")
			}

			// 只处理待审核状态的内容，删除而不是更新状态
			if content.Id > 0 && content.Status == "待审核" {
				// 直接删除关联内容
				_, err = tx.Model("content").Where("id", contentInfo.ContentId).Delete()
				if err != nil {
					return gerror.Wrap(err, "删除关联内容失败")
				}
				g.Log().Info(ctx, fmt.Sprintf("已删除已取消订单 %s 关联的内容 ID: %d", order.OrderNo, contentInfo.ContentId))
			}
		}

		g.Log().Info(ctx, fmt.Sprintf("订单 %s 已被客户取消", order.OrderNo))
		return nil
	})

	if err != nil {
		return nil, gerror.Wrap(err, "取消订单失败")
	}

	// 设置取消成功标志
	res.Cancelled = true

	return res, nil
}

// getStatusText 获取订单状态文本
func getStatusText(status int) string {
	switch status {
	case consts.OrderStatusWaitPay:
		return "待支付"
	case consts.OrderStatusPaid:
		return "已支付"
	case consts.OrderStatusCancelled:
		return "已取消"
	case consts.OrderStatusRefunded:
		return "已退款"
	case consts.OrderStatusProcessing:
		return "进行中"
	case consts.OrderStatusCompleted:
		return "已完成"
	default:
		return "未知状态"
	}
}

// checkOrderStatus 检查订单状态，如已过期则自动处理
func (c *OrderController) checkOrderStatus(ctx context.Context, order *model.Order) error {
	// 只检查待支付状态的订单
	if order.Status != consts.OrderStatusWaitPay {
		return nil
	}

	// 检查订单是否已过期
	now := gtime.Now()
	if order.ExpireTime != nil && now.TimestampMilli() > order.ExpireTime.TimestampMilli() {
		// 订单已过期，自动取消订单并清理关联内容
		g.Log().Info(ctx, fmt.Sprintf("订单 %s 已过期，自动处理", order.OrderNo))
		return c.handleExpiredOrder(ctx, order)
	}

	return nil
}

// New 创建客户端订单控制器实例
func New() *OrderController {
	return &OrderController{}
}
