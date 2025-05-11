package exchange

import (
	"context"
	v1 "demo/api/exchange/v1"
	"demo/internal/dao"
	"demo/internal/model/entity"
	"demo/internal/service"
	"demo/utility/auth"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// 确保 sExchange 实现了 ExchangeService 接口
var _ service.ExchangeService = (*sExchange)(nil)

func init() {
	service.RegisterExchange(New())
}

// sExchange 兑换记录服务实现
type sExchange struct {
	exchangeRecordDao *dao.ExchangeRecordDao
}

// New 创建兑换记录服务实例
func New() *sExchange {
	return &sExchange{
		exchangeRecordDao: &dao.ExchangeRecordDao{},
	}
}

// Get 获取兑换记录详情
func (s *sExchange) Get(ctx context.Context, req *v1.ExchangeRecordGetReq) (res *v1.ExchangeRecordGetRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看此兑换记录")
	}

	// 获取兑换记录详情
	record, err := s.exchangeRecordDao.Get(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("获取兑换记录失败: " + err.Error())
	}
	if record == nil {
		return nil, gerror.New("兑换记录不存在")
	}

	// 构建响应
	res = &v1.ExchangeRecordGetRes{
		Id:              record.Id,
		ClientId:        record.ClientId,
		ClientName:      record.ClientName,
		RechargeAccount: record.RechargeAccount,
		ProductName:     record.ProductName,
		Duration:        record.Duration,
		ExchangeTime:    record.ExchangeTime,
		Status:          record.Status,
		Remark:          record.Remark,
		CreatedAt:       record.CreatedAt,
		UpdatedAt:       record.UpdatedAt,
	}

	return res, nil
}

// Create 创建兑换记录
func (s *sExchange) Create(ctx context.Context, req *v1.ExchangeRecordCreateReq) (res *v1.ExchangeRecordCreateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限创建兑换记录")
	}

	// 设置默认值
	exchangeTime := req.ExchangeTime
	if exchangeTime == nil {
		exchangeTime = gtime.Now()
	}

	status := req.Status
	if status == "" {
		status = "processing"
	}

	// 创建兑换记录
	record := &entity.ExchangeRecord{
		ClientId:        req.ClientId,
		ClientName:      req.ClientName,
		RechargeAccount: req.RechargeAccount,
		ProductName:     req.ProductName,
		Duration:        req.Duration,
		ExchangeTime:    exchangeTime,
		Status:          status,
		Remark:          req.Remark,
	}

	id, err := s.exchangeRecordDao.Create(ctx, record)
	if err != nil {
		return nil, gerror.New("创建兑换记录失败: " + err.Error())
	}

	// 构建响应
	res = &v1.ExchangeRecordCreateRes{
		Id: int(id),
	}

	return res, nil
}

// Update 更新兑换记录
func (s *sExchange) Update(ctx context.Context, req *v1.ExchangeRecordUpdateReq) (res *v1.ExchangeRecordUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新兑换记录")
	}

	// 检查记录是否存在
	_, err = s.exchangeRecordDao.Get(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("兑换记录不存在")
	}

	// 设置默认值
	exchangeTime := req.ExchangeTime
	if exchangeTime == nil {
		exchangeTime = gtime.Now()
	}

	// 更新兑换记录
	record := &entity.ExchangeRecord{
		Id:              req.Id,
		ClientId:        req.ClientId,
		ClientName:      req.ClientName,
		RechargeAccount: req.RechargeAccount,
		ProductName:     req.ProductName,
		Duration:        req.Duration,
		ExchangeTime:    exchangeTime,
		Status:          req.Status,
		Remark:          req.Remark,
	}

	err = s.exchangeRecordDao.Update(ctx, record)
	if err != nil {
		return nil, gerror.New("更新兑换记录失败: " + err.Error())
	}

	// 构建响应
	res = &v1.ExchangeRecordUpdateRes{}

	return res, nil
}

// Delete 删除兑换记录
func (s *sExchange) Delete(ctx context.Context, req *v1.ExchangeRecordDeleteReq) (res *v1.ExchangeRecordDeleteRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限删除兑换记录")
	}

	// 检查记录是否存在
	_, err = s.exchangeRecordDao.Get(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("兑换记录不存在")
	}

	// 删除兑换记录
	err = s.exchangeRecordDao.Delete(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("删除兑换记录失败: " + err.Error())
	}

	// 构建响应
	res = &v1.ExchangeRecordDeleteRes{}

	return res, nil
}

// UpdateStatus 更新兑换记录状态
func (s *sExchange) UpdateStatus(ctx context.Context, req *v1.ExchangeRecordStatusUpdateReq) (res *v1.ExchangeRecordStatusUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新兑换记录状态")
	}

	// 检查状态值是否有效
	if req.Status != "processing" && req.Status != "completed" && req.Status != "failed" {
		return nil, gerror.New("无效的状态值，应为：processing、completed、failed")
	}

	// 检查记录是否存在
	record, err := s.exchangeRecordDao.Get(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("兑换记录不存在")
	}

	// 如果更新状态为失败，且当前状态不是失败，则需要返还时长
	if req.Status == "failed" && record.Status != "failed" {
		// 获取客户当前的时长信息
		var clientDuration struct {
			RemainingDuration string `json:"remaining_duration"`
			UsedDuration      string `json:"used_duration"`
		}
		err = g.DB().Model("client_duration").
			Fields("remaining_duration", "used_duration").
			Where("client_id", record.ClientId).
			Scan(&clientDuration)

		if err == nil { // 时长记录存在
			// 计算时长（分钟）
			remainingMinutes := parseDurationToMinutes(clientDuration.RemainingDuration)
			usedMinutes := parseDurationToMinutes(clientDuration.UsedDuration)

			// 计算需要返还的时长（分钟）
			returnMinutes := record.Duration * 24 * 60 // 假设Duration字段是天数，转换为分钟

			// 更新客户剩余时长和已使用时长
			newRemainingMinutes := remainingMinutes + returnMinutes
			newUsedMinutes := usedMinutes - returnMinutes
			if newUsedMinutes < 0 {
				newUsedMinutes = 0 // 确保已使用时长不为负数
			}

			newRemainingDuration := formatDuration(newRemainingMinutes)
			newUsedDuration := formatDuration(newUsedMinutes)

			// 更新客户时长记录
			_, updateErr := g.DB().Model("client_duration").
				Data(g.Map{
					"remaining_duration": newRemainingDuration,
					"used_duration":      newUsedDuration,
					"updated_at":         gtime.Now().String(),
				}).
				Where("client_id", record.ClientId).
				Update()

			if updateErr != nil {
				g.Log().Error(ctx, "状态更新为失败，返还时长失败", updateErr)
				// 记录日志但继续执行，不中断流程
			} else {
				g.Log().Info(ctx, fmt.Sprintf("兑换记录ID %d 状态更新为失败，已返还时长 %d天 (%d分钟)",
					req.Id, record.Duration, returnMinutes))
			}
		} else {
			g.Log().Error(ctx, "获取客户时长信息失败，无法返还时长", err)
			// 记录日志但继续执行，不中断流程
		}
	}

	// 更新兑换记录状态
	err = s.exchangeRecordDao.UpdateStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, gerror.New("更新兑换记录状态失败: " + err.Error())
	}

	// 构建响应
	res = &v1.ExchangeRecordStatusUpdateRes{}

	return res, nil
}

// GetList 获取兑换记录列表
func (s *sExchange) GetList(ctx context.Context, req *v1.ExchangeRecordListReq) (res *v1.ExchangeRecordListRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看兑换记录列表")
	}

	// 获取数据
	records, total, err := s.exchangeRecordDao.GetList(ctx, req.Page, req.Size, req.Id, req.ClientId, req.Status)
	if err != nil {
		return nil, gerror.New("获取兑换记录列表失败: " + err.Error())
	}

	// 计算总页数
	pages := 0
	if total > 0 {
		pages = int(math.Ceil(float64(total) / float64(req.Size)))
	}

	// 构建响应
	res = &v1.ExchangeRecordListRes{
		List:  make([]v1.ExchangeRecordItem, len(records)),
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
		Pages: pages,
	}

	// 转换记录列表
	for i, record := range records {
		res.List[i] = v1.ExchangeRecordItem{
			Id:              record.Id,
			ClientId:        record.ClientId,
			ClientName:      record.ClientName,
			RechargeAccount: record.RechargeAccount,
			ProductName:     record.ProductName,
			Duration:        record.Duration,
			ExchangeTime:    record.ExchangeTime,
			Status:          record.Status,
			Remark:          record.Remark,
			CreatedAt:       record.CreatedAt,
		}
	}

	return res, nil
}

// WxGetPage 微信客户端分页获取兑换记录列表
func (s *sExchange) WxGetPage(ctx context.Context, req *v1.WxExchangeRecordPageReq) (res *v1.WxExchangeRecordPageRes, err error) {
	// 获取当前登录客户信息
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 获取分页数据
	records, total, err := s.exchangeRecordDao.GetWxPage(ctx, req.Page, req.Size, clientId)
	if err != nil {
		return nil, gerror.New("获取兑换记录列表失败: " + err.Error())
	}

	// 计算总页数
	pages := 0
	if total > 0 {
		pages = int(math.Ceil(float64(total) / float64(req.Size)))
	}

	// 构建响应
	res = &v1.WxExchangeRecordPageRes{
		List:  make([]v1.ExchangeRecordItem, len(records)),
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
		Pages: pages,
	}

	// 转换记录列表
	for i, record := range records {
		res.List[i] = v1.ExchangeRecordItem{
			Id:              record.Id,
			ClientId:        record.ClientId,
			ClientName:      record.ClientName,
			RechargeAccount: record.RechargeAccount,
			ProductName:     record.ProductName,
			Duration:        record.Duration,
			ExchangeTime:    record.ExchangeTime,
			Status:          record.Status,
			Remark:          record.Remark,
			CreatedAt:       record.CreatedAt,
		}
	}

	return res, nil
}

// WxCreate 微信客户端创建兑换记录
func (s *sExchange) WxCreate(ctx context.Context, req *v1.WxExchangeRecordCreateReq) (res *v1.WxExchangeRecordCreateRes, err error) {
	// 获取当前登录客户信息
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 查询客户信息获取客户名称
	var clientInfo struct {
		RealName string `json:"real_name"`
	}
	err = g.DB().Model("client").
		Fields("real_name").
		Where("id", clientId).
		Scan(&clientInfo)
	if err != nil {
		return nil, gerror.New("获取客户信息失败: " + err.Error())
	}

	// 检查客户剩余时长是否足够
	var clientDuration struct {
		RemainingDuration string `json:"remaining_duration"`
		UsedDuration      string `json:"used_duration"`
	}
	err = g.DB().Model("client_duration").
		Fields("remaining_duration", "used_duration").
		Where("client_id", clientId).
		Scan(&clientDuration)

	// 如果查询出错但不是因为记录不存在，则返回错误
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return nil, gerror.New("获取客户时长信息失败: " + err.Error())
	}

	// 计算客户剩余时长（分钟）
	remainingMinutes := 0
	usedMinutes := 0
	if err == nil { // 时长记录存在
		// 解析剩余时长字符串为分钟数
		remainingMinutes = parseDurationToMinutes(clientDuration.RemainingDuration)
		// 解析已使用时长字符串为分钟数
		usedMinutes = parseDurationToMinutes(clientDuration.UsedDuration)
	}

	// 兑换所需时长（分钟）- req.Duration可能是天数，需要转换为分钟
	requiredMinutes := req.Duration * 24 * 60 // 假设Duration字段是天数，转换为分钟

	// 检查剩余时长是否足够
	if remainingMinutes < requiredMinutes {
		g.Log().Info(ctx, fmt.Sprintf("兑换失败：客户ID %d, 剩余时长 %s (%d分钟), 需要时长 %d天 (%d分钟)",
			clientId, clientDuration.RemainingDuration, remainingMinutes, req.Duration, requiredMinutes))
		return nil, gerror.New("剩余时长不足，无法完成兑换")
	}

	// 创建兑换记录
	record := &entity.ExchangeRecord{
		ClientId:        clientId,
		ClientName:      clientInfo.RealName,
		RechargeAccount: req.RechargeAccount,
		ProductName:     req.ProductName,
		Duration:        req.Duration,
		ExchangeTime:    gtime.Now(),
		Status:          "processing",
		Remark:          "",
	}

	id, err := s.exchangeRecordDao.Create(ctx, record)
	if err != nil {
		return nil, gerror.New("创建兑换记录失败: " + err.Error())
	}

	// 更新客户剩余时长和已使用时长
	newRemainingMinutes := remainingMinutes - requiredMinutes
	newUsedMinutes := usedMinutes + requiredMinutes
	newRemainingDuration := formatDuration(newRemainingMinutes)
	newUsedDuration := formatDuration(newUsedMinutes)

	// 更新客户时长记录
	_, err = g.DB().Model("client_duration").
		Data(g.Map{
			"remaining_duration": newRemainingDuration,
			"used_duration":      newUsedDuration,
			"updated_at":         gtime.Now().String(),
		}).
		Where("client_id", clientId).
		Update()

	if err != nil {
		g.Log().Error(ctx, "更新客户剩余时长失败", err)
		// 即使更新剩余时长失败，我们仍返回成功创建兑换记录的响应
		// 这个问题需要后续人工处理
	}

	// 构建响应
	res = &v1.WxExchangeRecordCreateRes{
		Id: int(id),
	}

	return res, nil
}

// WxGetPublicList 微信客户端获取公开兑换记录列表
func (s *sExchange) WxGetPublicList(ctx context.Context, req *v1.WxExchangeRecordPublicListReq) (res *v1.WxExchangeRecordPublicListRes, err error) {
	// 获取最新的兑换记录列表
	limit := req.Limit
	if limit <= 0 {
		limit = 10 // 默认获取10条
	}

	records, err := s.exchangeRecordDao.GetLatestRecords(ctx, limit)
	if err != nil {
		return nil, gerror.New("获取兑换记录列表失败: " + err.Error())
	}

	// 构建响应
	res = &v1.WxExchangeRecordPublicListRes{
		List: make([]v1.ExchangeRecordPublicItem, len(records)),
	}

	// 转换记录列表，并对客户名称进行处理
	for i, record := range records {
		// 处理客户名称，中间4位用*号代替
		maskedName := s.maskClientName(record.ClientName)

		res.List[i] = v1.ExchangeRecordPublicItem{
			ClientName:   maskedName,
			ProductName:  record.ProductName,
			ExchangeTime: record.ExchangeTime,
		}
	}

	return res, nil
}

// maskClientName 对客户名称进行脱敏处理（中间4位用*号屏蔽）
func (s *sExchange) maskClientName(name string) string {
	nameLen := len([]rune(name))

	// 如果名称长度小于等于2，直接返回原名称
	if nameLen <= 2 {
		return name
	}

	// 如果名称长度在3到6之间，保留首尾字符，中间用*替换
	if nameLen >= 3 && nameLen <= 6 {
		firstChar := string([]rune(name)[:1])
		lastChar := string([]rune(name)[nameLen-1:])
		stars := ""
		for i := 0; i < nameLen-2; i++ {
			stars += "*"
		}
		return firstChar + stars + lastChar
	}

	// 如果名称长度大于6，保留前两个和后两个字符，中间4位用*替换
	firstChars := string([]rune(name)[:2])
	lastChars := string([]rune(name)[nameLen-2:])
	return firstChars + "****" + lastChars
}

// 从易读的时长字符串解析为分钟数
func parseDurationToMinutes(durationStr string) int {
	// 默认为0分钟
	if durationStr == "" {
		return 0
	}

	totalMinutes := 0

	// 解析天数
	daysMatch := regexp.MustCompile(`(\d+)天`).FindStringSubmatch(durationStr)
	if len(daysMatch) > 1 {
		days, _ := strconv.Atoi(daysMatch[1])
		totalMinutes += days * 24 * 60
	}

	// 解析小时
	hoursMatch := regexp.MustCompile(`(\d+)小时`).FindStringSubmatch(durationStr)
	if len(hoursMatch) > 1 {
		hours, _ := strconv.Atoi(hoursMatch[1])
		totalMinutes += hours * 60
	}

	// 解析分钟
	minsMatch := regexp.MustCompile(`(\d+)分钟`).FindStringSubmatch(durationStr)
	if len(minsMatch) > 1 {
		mins, _ := strconv.Atoi(minsMatch[1])
		totalMinutes += mins
	}

	return totalMinutes
}

// 格式化时长为易读的字符串，如"3天18小时42分钟"
func formatDuration(minutes int) string {
	days := minutes / (24 * 60)
	hours := (minutes % (24 * 60)) / 60
	mins := minutes % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", days, hours, mins)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟", hours, mins)
	} else {
		return fmt.Sprintf("%d分钟", mins)
	}
}
