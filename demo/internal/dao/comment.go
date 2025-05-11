package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"demo/internal/model/do"
)

// ContentCommentDao 评论表DAO接口
type ContentCommentDao interface {
	// Insert 插入一条数据
	Insert(ctx context.Context, data *do.ContentCommentDO) (lastInsertId int64, err error)
	// Update 更新数据
	Update(ctx context.Context, data *do.ContentCommentDO, id interface{}) (rowsAffected int64, err error)
	// Delete 删除数据
	Delete(ctx context.Context, id interface{}) (rowsAffected int64, err error)
	// FindOne 查询单条数据
	FindOne(ctx context.Context, id interface{}) (*do.ContentCommentDO, error)
	// FindList 查询列表数据
	FindList(ctx context.Context, filter map[string]interface{}, page, size int) (list []*do.ContentCommentDO, total int, err error)
	// FindByContentId 根据内容ID查询评论列表
	FindByContentId(ctx context.Context, contentId interface{}, page, size int) (list []*do.ContentCommentDO, total int, err error)
	// UpdateStatus 更新状态
	UpdateStatus(ctx context.Context, id interface{}, status string) (rowsAffected int64, err error)
	// IncrementCommentCount 增加内容的评论数
	IncrementCommentCount(ctx context.Context, contentId interface{}) error
	// DecrementCommentCount 减少内容的评论数
	DecrementCommentCount(ctx context.Context, contentId interface{}) error
}

// contentCommentDao 评论表DAO实现
type contentCommentDao struct{}

// NewContentCommentDao 创建评论表DAO
func NewContentCommentDao() ContentCommentDao {
	return &contentCommentDao{}
}

// Insert 插入一条数据
func (d *contentCommentDao) Insert(ctx context.Context, data *do.ContentCommentDO) (lastInsertId int64, err error) {
	result, err := g.DB().Model(do.TableContentComment).Ctx(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 更新数据
func (d *contentCommentDao) Update(ctx context.Context, data *do.ContentCommentDO, id interface{}) (rowsAffected int64, err error) {
	result, err := g.DB().Model(do.TableContentComment).Ctx(ctx).Data(data).Where("id", id).Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete 删除数据
func (d *contentCommentDao) Delete(ctx context.Context, id interface{}) (rowsAffected int64, err error) {
	// 硬删除
	result, err := g.DB().Model(do.TableContentComment).Ctx(ctx).
		Where("id", id).
		Delete()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// FindOne 查询单条数据
func (d *contentCommentDao) FindOne(ctx context.Context, id interface{}) (*do.ContentCommentDO, error) {
	var comment *do.ContentCommentDO
	err := g.DB().Model(do.TableContentComment).Ctx(ctx).
		Where("id", id).
		Scan(&comment)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, nil
	}
	return comment, nil
}

// FindList 查询列表数据
func (d *contentCommentDao) FindList(ctx context.Context, filter map[string]interface{}, page, size int) (list []*do.ContentCommentDO, total int, err error) {
	model := g.DB().Model(do.TableContentComment).Ctx(ctx).Safe()

	// 添加过滤条件
	if filter != nil {
		// 内容ID过滤
		if contentId, ok := filter["content_id"]; ok && contentId != nil {
			model = model.Where("content_id", contentId)
		}

		// 客户ID过滤
		if clientId, ok := filter["client_id"]; ok && clientId != nil {
			model = model.Where("client_id", clientId)
		}

		// 状态过滤
		if status, ok := filter["status"]; ok && status != "" {
			model = model.Where("status", status)
		}

		// 真实姓名关键字搜索
		if realName, ok := filter["real_name"]; ok && realName != "" {
			model = model.WhereLike("real_name", "%"+realName.(string)+"%")
		}

		// 评论内容关键字搜索
		if comment, ok := filter["comment"]; ok && comment != "" {
			model = model.WhereLike("comment", "%"+comment.(string)+"%")
		}
	}

	// 计算总数
	count, err := model.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && size > 0 {
		model = model.Page(page, size)
	}

	// 按创建时间倒序
	model = model.Order("created_at DESC")

	// 查询数据
	var comments []*do.ContentCommentDO
	err = model.Scan(&comments)
	if err != nil {
		return nil, 0, err
	}

	return comments, count, nil
}

// FindByContentId 根据内容ID查询评论列表
func (d *contentCommentDao) FindByContentId(ctx context.Context, contentId interface{}, page, size int) (list []*do.ContentCommentDO, total int, err error) {
	filter := map[string]interface{}{
		"content_id": contentId,
		"status":     "已审核", // 只获取已审核的评论
	}
	return d.FindList(ctx, filter, page, size)
}

// UpdateStatus 更新状态
func (d *contentCommentDao) UpdateStatus(ctx context.Context, id interface{}, status string) (rowsAffected int64, err error) {
	result, err := g.DB().Model(do.TableContentComment).Ctx(ctx).
		Data(g.Map{
			"status":     status,
			"updated_at": gtime.Now(),
		}).
		Where("id", id).
		Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// IncrementCommentCount 增加内容的评论数
func (d *contentCommentDao) IncrementCommentCount(ctx context.Context, contentId interface{}) error {
	_, err := g.DB().Model(do.TableContent).Ctx(ctx).
		Data("comments=comments+1").
		Where("id", contentId).
		Update()
	return err
}

// DecrementCommentCount 减少内容的评论数
func (d *contentCommentDao) DecrementCommentCount(ctx context.Context, contentId interface{}) error {
	_, err := g.DB().Model(do.TableContent).Ctx(ctx).
		Data("comments=IF(comments>0, comments-1, 0)"). // 确保不会减到负数
		Where("id", contentId).
		Update()
	return err
}
