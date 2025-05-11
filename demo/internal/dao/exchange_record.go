package dao

import (
	"context"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExchangeRecordDao 兑换记录数据访问对象
type ExchangeRecordDao struct{}

// Get 根据ID获取兑换记录
func (d *ExchangeRecordDao) Get(ctx context.Context, id int) (*entity.ExchangeRecord, error) {
	var record entity.ExchangeRecord
	err := g.DB().Model("exchange_record").Where("id=?", id).Scan(&record)
	if err != nil || record.Id == 0 {
		return nil, err
	}
	return &record, nil
}

// Create 创建兑换记录
func (d *ExchangeRecordDao) Create(ctx context.Context, record *entity.ExchangeRecord) (int64, error) {
	now := gtime.Now()
	result, err := g.DB().Model("exchange_record").
		Data(g.Map{
			"client_id":        record.ClientId,
			"client_name":      record.ClientName,
			"recharge_account": record.RechargeAccount,
			"product_name":     record.ProductName,
			"duration":         record.Duration,
			"exchange_time":    record.ExchangeTime,
			"status":           record.Status,
			"remark":           record.Remark,
			"created_at":       now,
			"updated_at":       now,
		}).
		Insert()
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

// Update 更新兑换记录
func (d *ExchangeRecordDao) Update(ctx context.Context, record *entity.ExchangeRecord) error {
	_, err := g.DB().Model("exchange_record").
		Data(g.Map{
			"client_id":        record.ClientId,
			"client_name":      record.ClientName,
			"recharge_account": record.RechargeAccount,
			"product_name":     record.ProductName,
			"duration":         record.Duration,
			"exchange_time":    record.ExchangeTime,
			"status":           record.Status,
			"remark":           record.Remark,
			"updated_at":       gtime.Now(),
		}).
		Where("id=?", record.Id).
		Update()
	return err
}

// Delete 删除兑换记录
func (d *ExchangeRecordDao) Delete(ctx context.Context, id int) error {
	_, err := g.DB().Model("exchange_record").
		Where("id=?", id).
		Delete()
	return err
}

// UpdateStatus 更新兑换记录状态
func (d *ExchangeRecordDao) UpdateStatus(ctx context.Context, id int, status string) error {
	_, err := g.DB().Model("exchange_record").
		Data(g.Map{
			"status":     status,
			"updated_at": gtime.Now(),
		}).
		Where("id=?", id).
		Update()
	return err
}

// GetPage 分页获取兑换记录列表
func (d *ExchangeRecordDao) GetPage(ctx context.Context, page, size int, clientId int, status string) ([]entity.ExchangeRecord, int, error) {
	model := g.DB().Model("exchange_record")

	// 条件筛选
	if clientId > 0 {
		model = model.Where("client_id=?", clientId)
	}
	if status != "" {
		model = model.Where("status=?", status)
	}

	// 计算总数
	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	var records []entity.ExchangeRecord
	err = model.Page(page, size).
		Order("exchange_time DESC").
		Scan(&records)

	return records, total, err
}

// GetByClient 获取指定客户的兑换记录列表
func (d *ExchangeRecordDao) GetByClient(ctx context.Context, clientId int) ([]*entity.ExchangeRecord, error) {
	var records []*entity.ExchangeRecord
	err := g.DB().Model("exchange_record").
		Where("client_id=?", clientId).
		Order("exchange_time DESC").
		Scan(&records)
	return records, err
}

// GetWxPage 微信客户端分页获取兑换记录
func (d *ExchangeRecordDao) GetWxPage(ctx context.Context, page, size int, clientId int) ([]entity.ExchangeRecord, int, error) {
	model := g.DB().Model("exchange_record").
		Where("client_id=?", clientId)

	// 计算总数
	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	var records []entity.ExchangeRecord
	err = model.Page(page, size).
		Order("exchange_time DESC").
		Scan(&records)

	return records, total, err
}

// GetList 获取兑换记录列表（支持分页和筛选）
func (d *ExchangeRecordDao) GetList(ctx context.Context, page, size int, recordId int, clientId int, status string) ([]entity.ExchangeRecord, int, error) {
	model := g.DB().Model("exchange_record")

	// 条件筛选
	if recordId > 0 {
		model = model.Where("id=?", recordId)
	}
	if clientId > 0 {
		model = model.Where("client_id=?", clientId)
	}
	if status != "" {
		model = model.Where("status=?", status)
	}

	// 计算总数
	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	var records []entity.ExchangeRecord
	err = model.Page(page, size).
		Order("exchange_time DESC").
		Scan(&records)

	return records, total, err
}

// GetLatestRecords 获取最新的兑换记录列表
func (d *ExchangeRecordDao) GetLatestRecords(ctx context.Context, limit int) ([]entity.ExchangeRecord, error) {
	var records []entity.ExchangeRecord

	// 获取已完成的记录，使用英文状态值
	err := g.DB().Model("exchange_record").
		Where("status = ?", "completed").
		Order("exchange_time DESC").
		Limit(limit).
		Scan(&records)

	return records, err
}
