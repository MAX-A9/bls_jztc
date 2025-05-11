package order

import (
	"context"
	"fmt"

	v1 "demo/api/payment/v1"
	"demo/internal/consts"
	"demo/internal/model"
	"demo/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// 订单服务实现
type orderImpl struct{}

// New 创建订单服务实例
func New() service.OrderService {
	return &orderImpl{}
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

// List 获取订单列表
func (s *orderImpl) List(ctx context.Context, req *v1.OrderListReq) (res *v1.OrderListRes, err error) {
	res = &v1.OrderListRes{}

	// 构建查询条件
	m := g.DB().Model("order").Safe()

	if req.OrderNo != "" {
		m = m.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}

	if req.ClientName != "" {
		m = m.Where("client_name LIKE ?", "%"+req.ClientName+"%")
	}

	if req.Status != "" {
		m = m.Where("status = ?", req.Status)
	}

	if req.StartTime != "" {
		m = m.Where("created_at >= ?", req.StartTime)
	}

	if req.EndTime != "" {
		m = m.Where("created_at <= ?", req.EndTime)
	}

	if req.Product != "" {
		m = m.Where("product_name LIKE ?", "%"+req.Product+"%")
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
		res.List = make([]v1.OrderListItem, 0)
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
			// 检查订单是否已过期
			if err := s.checkOrderExpired(ctx, &orders[i]); err != nil {
				g.Log().Error(ctx, "处理过期订单失败:", err)
			} else if orders[i].Status == consts.OrderStatusCancelled {
				// 如果订单状态已改变为已取消，标记需要重新获取列表
				ordersChanged = true
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
	res.List = make([]v1.OrderListItem, len(orders))
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

		res.List[i] = v1.OrderListItem{
			Id:            order.Id,
			OrderNo:       order.OrderNo,
			ClientName:    order.ClientName,
			ContentId:     order.ContentId,
			ProductName:   order.ProductName,
			Amount:        order.Amount,
			Status:        order.Status,
			StatusText:    getStatusText(order.Status),
			PaymentMethod: order.PaymentMethod,
			CreatedAt:     createdAt,
			PayTime:       payTime,
			ExpireTime:    expireTime,
		}
	}

	return res, nil
}

// checkOrderExpired 检查订单是否已过期并处理
func (s *orderImpl) checkOrderExpired(ctx context.Context, order *model.Order) error {
	// 只检查待支付状态的订单
	if order.Status != consts.OrderStatusWaitPay {
		return nil
	}

	// 检查订单是否已过期
	now := gtime.Now()
	if order.ExpireTime != nil && now.TimestampMilli() > order.ExpireTime.TimestampMilli() {
		// 订单已过期，处理过期订单
		return s.handleExpiredOrder(ctx, order)
	}

	return nil
}

// handleExpiredOrder 处理过期订单
func (s *orderImpl) handleExpiredOrder(ctx context.Context, order *model.Order) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 更新订单状态为已取消
		_, err := tx.Model("order").Where("id", order.Id).Update(g.Map{
			"status":     consts.OrderStatusCancelled,
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

		// 修改订单对象的状态（用于返回前端显示）
		order.Status = consts.OrderStatusCancelled
		order.Remark = "订单已过期自动取消"

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

// Detail 获取订单详情
func (s *orderImpl) Detail(ctx context.Context, req *v1.OrderDetailReq) (res *v1.OrderDetailRes, err error) {
	// 查询订单
	var order model.Order
	err = g.DB().Model("order").Where("order_no = ?", req.OrderNo).Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单详情失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 检查订单是否已过期并自动处理
	if order.Status == consts.OrderStatusWaitPay {
		if err := s.checkOrderExpired(ctx, &order); err != nil {
			g.Log().Error(ctx, "处理过期订单失败:", err)
		}

		// 如果订单被标记为过期，需要重新获取最新状态
		if order.Status == consts.OrderStatusCancelled {
			err = g.DB().Model("order").Where("order_no = ?", req.OrderNo).Scan(&order)
			if err != nil {
				return nil, gerror.Wrap(err, "重新查询订单详情失败")
			}
		}
	}

	// 转换为API响应格式
	res = &v1.OrderDetailRes{
		Id:            order.Id,
		OrderNo:       order.OrderNo,
		ClientId:      order.ClientId,
		ClientName:    order.ClientName,
		ContentId:     order.ContentId,
		ProductName:   order.ProductName,
		Amount:        order.Amount,
		Status:        order.Status,
		StatusText:    getStatusText(order.Status),
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

// Cancel 取消订单
func (s *orderImpl) Cancel(ctx context.Context, req *v1.OrderCancelReq) (res *v1.OrderCancelRes, err error) {
	// 检查订单是否存在
	var order model.Order
	err = g.DB().Model("order").Where("order_no = ?", req.OrderNo).Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单信息失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 检查订单状态
	if order.Status != 0 {
		return nil, gerror.NewCode(gcode.New(400, "只有待支付的订单可以取消", nil))
	}

	// 使用事务同时取消订单并处理关联内容
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 更新订单状态为已取消
		_, err := tx.Model("order").Where("order_no = ?", req.OrderNo).Update(g.Map{
			"status":     2, // 已取消
			"updated_at": gtime.Now(),
			"remark":     "用户手动取消订单",
		})
		if err != nil {
			return gerror.Wrap(err, "取消订单失败")
		}

		// 2. 查询订单是否有关联内容
		var contentInfo struct {
			ContentId int `json:"content_id"`
		}
		err = tx.Model("order").Where("order_no = ?", req.OrderNo).Fields("content_id").Scan(&contentInfo)
		if err != nil {
			return gerror.Wrap(err, "查询订单内容关联失败")
		}

		// 3. 如果存在关联内容，检查内容状态并删除
		if contentInfo.ContentId > 0 {
			var content struct {
				Id     int    `json:"id"`
				Status string `json:"status"`
			}
			err = tx.Model("content").Where("id", contentInfo.ContentId).Scan(&content)
			if err != nil {
				return gerror.Wrap(err, "查询内容信息失败")
			}

			// 只删除待支付状态的内容，保留待审核状态的内容
			if content.Id > 0 && content.Status == "待支付" {
				// 删除关联内容
				_, err = tx.Model("content").Where("id", contentInfo.ContentId).Delete()
				if err != nil {
					return gerror.Wrap(err, "删除关联内容失败")
				}
				g.Log().Info(ctx, fmt.Sprintf("已删除订单 %s 关联的内容 ID: %d", req.OrderNo, contentInfo.ContentId))
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	res = &v1.OrderCancelRes{}
	return res, nil
}

// UpdateOrderStatus 更新订单状态
func (s *orderImpl) UpdateOrderStatus(ctx context.Context, orderNo string, status int, transactionId string) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询订单
		var order struct {
			Id        int    `json:"id"`
			OrderNo   string `json:"orderNo"`
			Status    int    `json:"status"`
			ContentId int    `json:"contentId"`
			Remark    string `json:"remark"`
		}

		err := tx.Model("order").Where("order_no = ?", orderNo).Fields("id,order_no,status,content_id,remark").Scan(&order)
		if err != nil {
			return gerror.Wrap(err, "查询订单信息失败")
		}

		if order.Id == 0 {
			return gerror.NewCode(gcode.New(404, "订单不存在", nil))
		}

		// 如果订单已支付，不再处理
		if order.Status == 1 {
			return nil
		}

		updateData := g.Map{
			"updated_at": gtime.Now(),
		}

		// 将已支付状态(1)自动转换为进行中状态(4)
		if status == consts.OrderStatusPaid {
			updateData["status"] = consts.OrderStatusProcessing // 设置为进行中
			updateData["pay_time"] = gtime.Now()
			updateData["transaction_id"] = transactionId
			updateData["payment_method"] = "wechat"

			// 设置订单过期时间
			// 从配置中获取默认过期天数，如果配置不存在则记录警告日志
			defaultExpireDaysValue := g.Cfg().MustGet(ctx, consts.ConfigOrderDefaultExpireDays)
			if defaultExpireDaysValue.IsEmpty() {
				g.Log().Warning(ctx, fmt.Sprintf("订单过期天数配置 %s 未设置", consts.ConfigOrderDefaultExpireDays))
				// 配置不存在时，不设置过期时间
			} else {
				defaultExpireDays := defaultExpireDaysValue.Int()
				updateData["expire_time"] = gtime.Now().AddDate(0, 0, defaultExpireDays)
				g.Log().Info(ctx, fmt.Sprintf("订单 %s 支付成功，自动设置为进行中状态，过期时间为 %d 天后",
					orderNo, defaultExpireDays))
			}
		} else {
			updateData["status"] = status

			if status == 1 {
				updateData["pay_time"] = gtime.Now()
				updateData["transaction_id"] = transactionId
				updateData["payment_method"] = "wechat"
			}
		}

		// 处理支付成功的情况，需要更新关联的内容状态
		if (status == consts.OrderStatusPaid || status == consts.OrderStatusProcessing) && order.ContentId > 0 {
			// 更新内容状态
			_, err = tx.Model("content").Where("id = ?", order.ContentId).Update(g.Map{
				"status":     "已发布", // 设置为已发布状态，而不是数值1
				"updated_at": gtime.Now(),
			})
			if err != nil {
				return gerror.Wrap(err, "更新内容状态失败")
			}
		}

		// 如果是取消状态或退款状态，且原来是待支付状态，需要处理关联内容
		if (status == consts.OrderStatusCancelled || status == consts.OrderStatusRefunded) &&
			order.Status == consts.OrderStatusWaitPay && order.ContentId > 0 {

			// 查询内容状态
			var content struct {
				Id     int    `json:"id"`
				Status string `json:"status"`
			}
			err = tx.Model("content").Where("id", order.ContentId).Scan(&content)
			if err != nil {
				return gerror.Wrap(err, "查询内容信息失败")
			}

			// 只处理待支付状态的内容
			if content.Id > 0 && content.Status == "待支付" {
				// 删除关联内容
				_, err = tx.Model("content").Where("id", order.ContentId).Delete()
				if err != nil {
					return gerror.Wrap(err, "删除关联内容失败")
				}
				g.Log().Info(ctx, fmt.Sprintf("已删除订单 %s 关联的待支付内容 ID: %d", orderNo, order.ContentId))
			}
		}

		// 更新订单状态
		_, err = tx.Model("order").Where("order_no = ?", orderNo).Update(updateData)
		if err != nil {
			return gerror.Wrap(err, "更新订单状态失败")
		}

		return nil
	})
}

// UpdateStatus 管理员更新订单状态
func (s *orderImpl) UpdateStatus(ctx context.Context, req *v1.UpdateOrderStatusReq) (res *v1.UpdateOrderStatusRes, err error) {
	// 初始化响应
	res = &v1.UpdateOrderStatusRes{
		Success: false,
		Message: "",
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询订单
		var order model.Order
		err := tx.Model("order").Where("order_no = ?", req.OrderNo).Scan(&order)
		if err != nil {
			return gerror.Wrap(err, "查询订单信息失败")
		}

		if order.Id == 0 {
			return gerror.NewCode(gcode.New(404, "订单不存在", nil))
		}

		// 先检查订单是否已过期，如果是待支付状态的订单
		if order.Status == consts.OrderStatusWaitPay {
			// 检查是否过期
			now := gtime.Now()
			if order.ExpireTime != nil && now.TimestampMilli() > order.ExpireTime.TimestampMilli() {
				// 订单已过期，但管理员仍要更新状态
				g.Log().Warning(ctx, fmt.Sprintf("订单 %s 已过期，管理员尝试更新其状态为 %d",
					req.OrderNo, req.Status))

				// 如果管理员尝试将过期订单设置为已取消状态，则允许操作继续
				// 否则提醒管理员订单已过期
				if req.Status != consts.OrderStatusCancelled {
					return gerror.New("该订单已过期，请先取消该订单或确认您的操作")
				}
			}
		}

		// 构建更新数据
		updateData := g.Map{
			"updated_at": gtime.Now(),
		}

		// 处理备注信息
		if req.Remark != "" {
			updateData["remark"] = req.Remark
		}

		// 如果是设置为已支付状态，直接改为进行中状态
		if req.Status == consts.OrderStatusPaid {
			// 将状态设置为进行中，而不是已支付
			updateData["status"] = consts.OrderStatusProcessing
			updateData["pay_time"] = gtime.Now()

			// 设置交易流水号
			if req.TransactionId != "" {
				updateData["transaction_id"] = req.TransactionId
			}

			// 设置支付方式
			if req.PaymentMethod != "" {
				updateData["payment_method"] = req.PaymentMethod
			} else {
				updateData["payment_method"] = "admin" // 如果未指定支付方式，默认为管理员手动设置
			}

			// 设置默认的有效期（如果没有设置）
			if order.ExpireTime == nil && req.ExpireTime == "" {
				// 从配置中获取默认过期天数，如果配置不存在则记录警告日志
				defaultExpireDaysValue := g.Cfg().MustGet(ctx, consts.ConfigOrderDefaultExpireDays)
				if defaultExpireDaysValue.IsEmpty() {
					g.Log().Warning(ctx, fmt.Sprintf("订单过期天数配置 %s 未设置", consts.ConfigOrderDefaultExpireDays))
					return gerror.New("订单过期天数配置未设置，无法设置订单过期时间")
				}

				defaultExpireDays := defaultExpireDaysValue.Int()
				updateData["expire_time"] = gtime.Now().AddDate(0, 0, defaultExpireDays)

				g.Log().Info(ctx,
					fmt.Sprintf("管理员将订单 %s 设置为已支付，自动设置过期时间为 %d 天后",
						req.OrderNo, defaultExpireDays))
			} else if req.ExpireTime != "" {
				// 使用请求中的过期时间
				expireTime, err := gtime.StrToTime(req.ExpireTime)
				if err != nil {
					return gerror.Wrap(err, "过期时间格式不正确")
				}
				updateData["expire_time"] = expireTime
			}

			// 如果订单有关联内容，更新内容状态
			if order.ContentId > 0 {
				var content struct {
					Id     int    `json:"id"`
					Status string `json:"status"`
				}
				err = tx.Model("content").Where("id = ?", order.ContentId).Scan(&content)
				if err != nil {
					return gerror.Wrap(err, "查询内容信息失败")
				}

				if content.Id > 0 && content.Status == "待审核" {
					_, err = tx.Model("content").Where("id = ?", order.ContentId).Data(g.Map{
						"status":     "已发布",
						"updated_at": gtime.Now(),
					}).Update()
					if err != nil {
						return gerror.Wrap(err, "更新内容状态失败")
					}
					g.Log().Info(ctx, fmt.Sprintf("已将订单 %s 关联的内容 ID: %d 状态更新为已发布", req.OrderNo, order.ContentId))
				}
			}
		} else if req.Status == consts.OrderStatusProcessing && order.ExpireTime == nil && req.PackageInfo == "" {
			// 如果是从其他状态改为进行中，且未设置过期时间，则设置默认过期时间
			defaultExpireDaysValue := g.Cfg().MustGet(ctx, consts.ConfigOrderDefaultExpireDays)
			if defaultExpireDaysValue.IsEmpty() {
				g.Log().Warning(ctx, fmt.Sprintf("订单过期天数配置 %s 未设置", consts.ConfigOrderDefaultExpireDays))
				return gerror.New("订单过期天数配置未设置，无法设置订单过期时间")
			}

			defaultExpireDays := defaultExpireDaysValue.Int()
			updateData["expire_time"] = gtime.Now().AddDate(0, 0, defaultExpireDays)
			g.Log().Info(ctx,
				fmt.Sprintf("管理员将订单 %s 设置为进行中，自动设置过期时间为 %d 天后",
					req.OrderNo, defaultExpireDays))
		} else {
			// 设置其他状态
			updateData["status"] = req.Status
		}

		// 如果提供了套餐信息
		if req.PackageInfo != "" {
			updateData["package_info"] = req.PackageInfo
		}

		// 更新订单
		_, err = tx.Model("order").Where("order_no = ?", req.OrderNo).Data(updateData).Update()
		if err != nil {
			return gerror.Wrap(err, "更新订单状态失败")
		}

		// 记录状态变更日志
		g.Log().Info(ctx, fmt.Sprintf("管理员更新订单 %s 状态为 %s", req.OrderNo, getStatusChangeMessage(req.Status)))

		return nil
	})

	if err != nil {
		res.Message = err.Error()
		return res, nil
	}

	res = &v1.UpdateOrderStatusRes{
		Success: true,
		Message: fmt.Sprintf("订单状态已更新为%s", getStatusChangeMessage(req.Status)),
	}

	return res, nil
}

// getStatusChangeMessage 获取状态变更消息
func getStatusChangeMessage(status int) string {
	switch status {
	case consts.OrderStatusWaitPay:
		return "订单已设置为待支付状态"
	case consts.OrderStatusPaid:
		return "订单已设置为已支付状态"
	case consts.OrderStatusCancelled:
		return "订单已设置为已取消状态"
	case consts.OrderStatusRefunded:
		return "订单已设置为已退款状态"
	case consts.OrderStatusProcessing:
		return "订单已设置为进行中状态"
	case consts.OrderStatusCompleted:
		return "订单已设置为已完成状态"
	default:
		return "订单状态已更新"
	}
}

// Delete 删除订单
func (s *orderImpl) Delete(ctx context.Context, req *v1.OrderDeleteReq) (res *v1.OrderDeleteRes, err error) {
	// 查询订单
	var order model.Order
	err = g.DB().Model("order").Where("order_no = ?", req.OrderNo).Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单信息失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 使用事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 处理关联内容
		if order.ContentId > 0 {
			// 查询关联内容的状态
			var content struct {
				Id     int    `json:"id"`
				Status string `json:"status"`
			}
			err := tx.Model("content").Where("id", order.ContentId).Scan(&content)
			if err != nil {
				return gerror.Wrap(err, "查询内容信息失败")
			}

			// 如果是待支付状态的内容，则删除
			if content.Id > 0 && content.Status == "待支付" {
				_, err = tx.Model("content").Where("id", order.ContentId).Delete()
				if err != nil {
					return gerror.Wrap(err, "删除关联内容失败")
				}
				g.Log().Info(ctx, fmt.Sprintf("已删除订单 %s 关联的待支付内容 ID: %d", req.OrderNo, order.ContentId))
			} else {
				// 仅记录日志，不删除其他状态的内容
				g.Log().Info(ctx, fmt.Sprintf("订单 %s 关联内容ID: %d，状态为 %s，不自动删除", req.OrderNo, order.ContentId, content.Status))
			}
		}

		// 删除订单
		_, err := tx.Model("order").Where("order_no = ?", req.OrderNo).Delete()
		if err != nil {
			return gerror.Wrap(err, "删除订单失败")
		}

		// 记录日志
		g.Log().Info(ctx, fmt.Sprintf("订单 %s 已成功删除", req.OrderNo))

		return nil
	})

	if err != nil {
		return nil, err
	}

	res = &v1.OrderDeleteRes{}
	return res, nil
}
