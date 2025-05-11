package order

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "demo/api/payment/v1"
	"demo/internal/consts"
	"demo/internal/model"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sOrder struct{}

// New 创建订单服务实例
func New() service.OrderService {
	return &sOrder{}
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
func (s *sOrder) List(ctx context.Context, req *v1.OrderListReq) (res *v1.OrderListRes, err error) {
	res = &v1.OrderListRes{}

	// 验证权限
	userId, username, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}
	if role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限查看订单列表", nil))
	}

	// 记录日志
	defer func() {
		if err != nil {
			service.Log().Record(ctx, userId, username, "订单管理", "查询订单列表", 0, err.Error())
		} else {
			service.Log().Record(ctx, userId, username, "订单管理", "查询订单列表", 1, "")
		}
	}()

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

	// 转换为API响应格式
	res.List = make([]v1.OrderListItem, len(orders))
	for i, order := range orders {
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
			CreatedAt:     order.CreatedAt.String(),
			PayTime:       order.PayTime.String(),
		}
	}

	return res, nil
}

// Detail 获取订单详情
func (s *sOrder) Detail(ctx context.Context, req *v1.OrderDetailReq) (res *v1.OrderDetailRes, err error) {
	// 验证权限
	userId, username, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}
	if role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限查看订单详情", nil))
	}

	// 记录日志
	defer func() {
		if err != nil {
			service.Log().Record(ctx, userId, username, "订单管理", "查询订单详情", 0, err.Error())
		} else {
			service.Log().Record(ctx, userId, username, "订单管理", "查询订单详情", 1, "")
		}
	}()

	// 查询订单
	var order model.Order
	err = g.DB().Model("order").Where("order_no = ?", req.OrderNo).Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单详情失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
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
	}

	if order.CreatedAt != nil {
		res.CreatedAt = order.CreatedAt.String()
	}

	if order.PayTime != nil {
		res.PayTime = order.PayTime.String()
	}

	return res, nil
}

// Cancel 取消订单
func (s *sOrder) Cancel(ctx context.Context, req *v1.OrderCancelReq) (res *v1.OrderCancelRes, err error) {
	// 验证权限
	userId, username, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}

	// 检查订单是否存在
	var order model.Order
	err = g.DB().Model("order").Where("order_no = ?", req.OrderNo).Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单信息失败")
	}

	if order.Id == 0 {
		return nil, gerror.NewCode(gcode.New(404, "订单不存在", nil))
	}

	// 检查权限，只有管理员或订单所属客户可以取消
	if role != "admin" && userId != order.ClientId {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限取消此订单", nil))
	}

	// 检查订单状态
	if order.Status != 0 {
		return nil, gerror.NewCode(gcode.New(400, "只有待支付的订单可以取消", nil))
	}

	// 记录日志
	defer func() {
		if err != nil {
			service.Log().Record(ctx, userId, username, "订单管理", "取消订单", 0, err.Error())
		} else {
			service.Log().Record(ctx, userId, username, "订单管理", "取消订单", 1, "")
		}
	}()

	// 更新订单状态为已取消
	_, err = g.DB().Model("order").Where("order_no = ?", req.OrderNo).Update(g.Map{
		"status":     2, // 已取消
		"updated_at": gtime.Now(),
	})

	if err != nil {
		return nil, gerror.Wrap(err, "取消订单失败")
	}

	res = &v1.OrderCancelRes{}
	return res, nil
}

// UpdateOrderStatus 更新订单状态
func (s *sOrder) UpdateOrderStatus(ctx context.Context, orderNo string, status int, transactionId string) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查询订单
		var order model.Order
		err := tx.Model("order").Where("order_no = ?", orderNo).Scan(&order)
		if err != nil {
			return gerror.Wrap(err, "查询订单信息失败")
		}

		if order.Id == 0 {
			return gerror.NewCode(gcode.New(404, "订单不存在", nil))
		}

		// 如果订单已非待支付状态，不再处理
		if order.Status != consts.OrderStatusWaitPay {
			return nil
		}

		// 支付成功后直接将状态设为进行中(4)，而不是已支付(1)
		updateData := g.Map{
			"status":     consts.OrderStatusProcessing, // 直接设置为进行中状态
			"updated_at": gtime.Now(),
			"pay_time":   gtime.Now(),
		}

		// 更新交易信息
		if transactionId != "" {
			updateData["transaction_id"] = transactionId
			updateData["payment_method"] = "wechat"
		}

		// 设置默认的有效期（从全局配置获取过期天数）
		if order.ExpireTime == nil {
			// 从配置中获取默认过期天数，如果配置不存在则记录警告日志
			defaultExpireDaysValue := g.Cfg().MustGet(ctx, consts.ConfigOrderDefaultExpireDays)
			if defaultExpireDaysValue.IsEmpty() {
				g.Log().Warning(ctx, fmt.Sprintf("订单过期天数配置 %s 未设置", consts.ConfigOrderDefaultExpireDays))
				// 配置不存在时，不设置过期时间，正常处理其他逻辑
			} else {
				defaultExpireDays := defaultExpireDaysValue.Int()
				updateData["expire_time"] = gtime.Now().AddDate(0, 0, defaultExpireDays)

				g.Log().Info(ctx,
					fmt.Sprintf("订单 %s 支付成功，自动设置过期时间为 %d 天后",
						orderNo, defaultExpireDays))
			}
		}

		// 更新订单状态
		_, err = tx.Model("order").Where("order_no = ?", orderNo).Update(updateData)
		if err != nil {
			return gerror.Wrap(err, "更新订单状态失败")
		}

		// 处理关联内容
		if order.ContentId > 0 {
			// 更新内容状态为已发布
			_, err = tx.Model("content").Where("id = ?", order.ContentId).Update(g.Map{
				"status":     "已发布",
				"updated_at": gtime.Now(),
			})
			if err != nil {
				return gerror.Wrap(err, "更新内容状态失败")
			}

			g.Log().Info(ctx, fmt.Sprintf("订单 %s 支付成功，已更新关联内容 %d 状态为已发布", orderNo, order.ContentId))
		}

		g.Log().Info(ctx, fmt.Sprintf("订单 %s 支付成功，状态已设置为进行中", orderNo))
		return nil
	})
}

// UpdateStatus 管理员更新订单状态
func (s *sOrder) UpdateStatus(ctx context.Context, req *v1.UpdateOrderStatusReq) (res *v1.UpdateOrderStatusRes, err error) {
	// 验证权限
	userId, username, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}
	if role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限更新订单状态", nil))
	}

	// 记录日志
	defer func() {
		if err != nil {
			service.Log().Record(ctx, userId, username, "订单管理", "更新订单状态", 0, err.Error())
		} else {
			service.Log().Record(ctx, userId, username, "订单管理", "更新订单状态", 1, "")
		}
	}()

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
				// 更新内容状态
				_, err = tx.Model("content").Where("id = ?", order.ContentId).Update(g.Map{
					"status":     "已发布", // 设置为已发布状态
					"updated_at": gtime.Now(),
				})
				if err != nil {
					return gerror.Wrap(err, "更新内容状态失败")
				}
			}

			g.Log().Info(ctx, fmt.Sprintf("管理员将订单 %s 设置为已支付，状态已更新为进行中", req.OrderNo))
		} else {
			// 其他状态直接设置
			updateData["status"] = req.Status

			// 处理套餐信息
			if req.PackageInfo != "" {
				// 验证套餐信息是否为有效的JSON格式
				var packageInfoCheck interface{}
				if err := json.Unmarshal([]byte(req.PackageInfo), &packageInfoCheck); err != nil {
					return gerror.Wrap(err, "套餐信息格式错误，请提供有效的JSON格式")
				}
				updateData["package_info"] = req.PackageInfo
			}

			// 处理到期时间
			if req.ExpireTime != "" {
				expireTime, err := gtime.StrToTime(req.ExpireTime)
				if err != nil {
					return gerror.Wrap(err, "到期时间格式错误，请使用yyyy-MM-dd HH:mm:ss格式")
				}
				updateData["expire_time"] = expireTime
			} else if req.Status == consts.OrderStatusProcessing && order.ExpireTime == nil && req.PackageInfo == "" {
				// 如果是设置为进行中状态，但没有提供到期时间和套餐信息，提示错误
				return gerror.New("设置为进行中状态时必须提供到期时间或套餐信息")
			}

			// 根据状态处理内容关联
			if order.ContentId > 0 {
				var contentUpdateData g.Map

				switch req.Status {
				case consts.OrderStatusProcessing:
					// 进行中状态，内容设为已发布
					contentUpdateData = g.Map{
						"status":     "已发布",
						"updated_at": gtime.Now(),
					}
				case consts.OrderStatusCompleted:
					// 已完成状态，内容状态设为已完成
					contentUpdateData = g.Map{
						"status":     "已完成",
						"updated_at": gtime.Now(),
					}
				case consts.OrderStatusCancelled, consts.OrderStatusRefunded:
					// 已取消或已退款状态，内容状态设为已关闭
					contentUpdateData = g.Map{
						"status":     "已关闭",
						"updated_at": gtime.Now(),
					}
				}

				// 更新内容状态
				if len(contentUpdateData) > 0 {
					_, err = tx.Model("content").Where("id = ?", order.ContentId).Update(contentUpdateData)
					if err != nil {
						return gerror.Wrap(err, "更新内容状态失败")
					}
				}
			}
		}

		// 更新订单状态
		_, err = tx.Model("order").Where("order_no = ?", req.OrderNo).Update(updateData)
		if err != nil {
			return gerror.Wrap(err, "更新订单状态失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 构建响应
	res = &v1.UpdateOrderStatusRes{
		Success: true,
		Message: getStatusChangeMessage(req.Status),
	}
	return res, nil
}

// getStatusChangeMessage 获取状态变更信息
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
func (s *sOrder) Delete(ctx context.Context, req *v1.OrderDeleteReq) (res *v1.OrderDeleteRes, err error) {
	// 验证权限
	userId, username, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}
	if role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限删除订单", nil))
	}

	// 记录日志
	defer func() {
		if err != nil {
			service.Log().Record(ctx, userId, username, "订单管理", "删除订单", 0, err.Error())
		} else {
			service.Log().Record(ctx, userId, username, "订单管理", "删除订单", 1, "")
		}
	}()

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
			// 仅记录日志，不删除内容
			g.Log().Info(ctx, fmt.Sprintf("订单 %s 关联内容ID: %d，请注意该内容可能需要手动处理", req.OrderNo, order.ContentId))
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
