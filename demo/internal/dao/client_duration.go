package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// ClientDuration DAO 客户时长表数据访问对象
type ClientDuration struct{}

// Model 获取客户时长表模型
func (d *ClientDuration) Model(ctx context.Context) *gdb.Model {
	return g.DB().Model("client_duration").Ctx(ctx)
}

// GetWithPage 分页获取客户时长列表
func (d *ClientDuration) GetWithPage(ctx context.Context, page, pageSize int, condition g.Map) (list gdb.Result, total int, err error) {
	model := d.Model(ctx).Where(condition)
	total, err = model.Count()
	if err != nil {
		return nil, 0, gerror.New("获取客户时长总数失败")
	}

	list, err = model.Fields(
		"id",
		"client_id",
		"client_name",
		"remaining_duration",
		"total_duration",
		"used_duration",
		"created_at",
		"updated_at",
	).Order("id DESC").Page(page, pageSize).All()
	if err != nil {
		return nil, 0, gerror.New("获取客户时长列表失败")
	}

	return list, total, nil
}

// GetById 根据ID获取客户时长信息
func (d *ClientDuration) GetById(ctx context.Context, id int) (gdb.Record, error) {
	record, err := d.Model(ctx).Where("id", id).One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, gerror.New("客户时长记录不存在")
	}
	return record, nil
}

// GetByClientId 根据客户ID获取客户时长信息
func (d *ClientDuration) GetByClientId(ctx context.Context, clientId int) (gdb.Record, error) {
	record, err := d.Model(ctx).Where("client_id", clientId).One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, gerror.New("客户时长记录不存在")
	}
	return record, nil
}

// Insert 插入客户时长记录
func (d *ClientDuration) Insert(ctx context.Context, data g.Map) (int64, error) {
	res, err := d.Model(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateById 更新客户时长记录
func (d *ClientDuration) UpdateById(ctx context.Context, id int, data g.Map) error {
	_, err := d.Model(ctx).Where("id", id).Data(data).Update()
	return err
}

// DeleteById 删除客户时长记录
func (d *ClientDuration) DeleteById(ctx context.Context, id int) error {
	_, err := d.Model(ctx).Where("id", id).Delete()
	return err
}

// New 创建并返回客户时长表数据访问对象
func NewClientDuration() *ClientDuration {
	return &ClientDuration{}
}

// 客户时长表单例
var clientDurationDao = NewClientDuration()

// ClientDurationDao 获取客户时长表数据访问对象单例
func ClientDurationDao() *ClientDuration {
	return clientDurationDao
}
