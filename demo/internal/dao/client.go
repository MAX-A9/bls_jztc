package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// Client DAO 客户表数据访问对象
type Client struct{}

// Model 获取客户表模型
func (d *Client) Model(ctx context.Context) *gdb.Model {
	return g.DB().Model("client").Ctx(ctx)
}

// GetWithPage 分页获取客户列表
func (d *Client) GetWithPage(ctx context.Context, page, pageSize int, condition g.Map) (list gdb.Result, total int, err error) {
	model := d.Model(ctx).Where(condition)
	total, err = model.Count()
	if err != nil {
		return nil, 0, gerror.New("获取客户总数失败")
	}

	list, err = model.Fields(
		"id",
		"username",
		"real_name",
		"phone",
		"status",
		"identifier",
		"avatar_url",
		"created_at",
		"last_login_time as last_login_at",
		"last_login_ip",
	).Order("id DESC").Page(page, pageSize).All()
	if err != nil {
		return nil, 0, gerror.New("获取客户列表失败")
	}

	return list, total, nil
}

// GetById 根据ID获取客户信息
func (d *Client) GetById(ctx context.Context, id int) (gdb.Record, error) {
	record, err := d.Model(ctx).Where("id", id).One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, gerror.New("客户不存在")
	}
	return record, nil
}

// Insert 插入客户记录
func (d *Client) Insert(ctx context.Context, data g.Map) (int64, error) {
	res, err := d.Model(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateById 更新客户记录
func (d *Client) UpdateById(ctx context.Context, id int, data g.Map) error {
	_, err := d.Model(ctx).Where("id", id).Data(data).Update()
	return err
}

// DeleteById 删除客户记录
func (d *Client) DeleteById(ctx context.Context, id int) error {
	_, err := d.Model(ctx).Where("id", id).Delete()
	return err
}

// New 创建并返回客户表数据访问对象
func NewClient() *Client {
	return &Client{}
}

// 客户表单例
var clientDao = NewClient()

// ClientDao 获取客户表数据访问对象单例
func ClientDao() *Client {
	return clientDao
}
