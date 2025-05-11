package dao

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"demo/internal/model/do"
)

// ContentDao 内容表DAO接口
type ContentDao interface {
	// Insert 插入一条数据
	Insert(ctx context.Context, data *do.ContentDO) (lastInsertId int64, err error)
	// Update 更新数据
	Update(ctx context.Context, data *do.ContentDO, id interface{}) (rowsAffected int64, err error)
	// Delete 删除数据
	Delete(ctx context.Context, id interface{}) (rowsAffected int64, err error)
	// FindOne 查询单条数据
	FindOne(ctx context.Context, id interface{}) (*do.ContentDO, error)
	// FindList 查询列表数据
	FindList(ctx context.Context, filter map[string]interface{}, page, size int) (list []*do.ContentDO, total int, err error)
	// UpdateStatus 更新状态
	UpdateStatus(ctx context.Context, id interface{}, status string) (rowsAffected int64, err error)
	// UpdateRecommend 更新推荐状态
	UpdateRecommend(ctx context.Context, id interface{}, isRecommended bool, topUntil *gtime.Time) (rowsAffected int64, err error)
	// IncrementViews 增加浏览量
	IncrementViews(ctx context.Context, id interface{}) error
	// FindByIds 根据ID列表查询多条数据
	FindByIds(ctx context.Context, ids []interface{}) ([]*do.ContentDO, error)
	// GetContentInfo 获取内容信息
	GetContentInfo(ctx context.Context, tableName string, contentId int, result interface{}) error
}

// contentDao 内容表DAO实现
type contentDao struct{}

// NewContentDao 创建内容表DAO
func NewContentDao() ContentDao {
	return &contentDao{}
}

// Insert 插入一条数据
func (d *contentDao) Insert(ctx context.Context, data *do.ContentDO) (lastInsertId int64, err error) {
	result, err := g.DB().Model(do.TableContent).Ctx(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 更新数据
func (d *contentDao) Update(ctx context.Context, data *do.ContentDO, id interface{}) (rowsAffected int64, err error) {
	result, err := g.DB().Model(do.TableContent).Ctx(ctx).Data(data).Where("id", id).Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete 删除数据
func (d *contentDao) Delete(ctx context.Context, id interface{}) (rowsAffected int64, err error) {
	// 硬删除
	result, err := g.DB().Model(do.TableContent).Ctx(ctx).
		Where("id", id).
		Delete()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// FindOne 查询单条数据
func (d *contentDao) FindOne(ctx context.Context, id interface{}) (*do.ContentDO, error) {
	var content *do.ContentDO
	err := g.DB().Model(do.TableContent).Ctx(ctx).
		Where("id", id).
		Scan(&content)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, nil
	}
	return content, nil
}

// FindList 查询列表数据
func (d *contentDao) FindList(ctx context.Context, filter map[string]interface{}, page, size int) (list []*do.ContentDO, total int, err error) {
	model := g.DB().Model(do.TableContent).Ctx(ctx).Safe()

	// 处理内容过期逻辑：将已过期但状态不是"已下架"的内容更新为"已下架"
	now := gtime.Now()
	_, err = g.DB().Model(do.TableContent).
		Where("expires_at IS NOT NULL").
		Where("expires_at < ?", now).
		Where("status != ?", "已下架").
		Data(g.Map{
			"status":     "已下架",
			"updated_at": now,
		}).
		Update()

	if err != nil {
		g.Log().Error(ctx, "更新过期内容状态失败:", err)
		// 继续执行查询，不直接返回错误
	} else {
		g.Log().Debug(ctx, "已更新过期内容状态为已下架")
	}

	// 添加过滤条件
	if filter != nil {
		// 标题模糊搜索
		if title, ok := filter["title"]; ok && title != "" {
			model = model.WhereLike("title", "%"+title.(string)+"%")
		}

		// 分类查询
		if category, ok := filter["category"]; ok && category != "" {
			model = model.Where("category", category)
		}

		// 状态查询
		if status, ok := filter["status"]; ok && status != "" {
			model = model.Where("status", status)
		}

		// 作者查询
		if author, ok := filter["author"]; ok && author != "" {
			model = model.WhereLike("author", "%"+author.(string)+"%")
		}

		// 地区ID查询
		if regionId, ok := filter["region_id"]; ok && regionId != nil {
			model = model.Where("region_id", regionId)
		}

		// 内容类型查询，从extend字段中解析type值
		if contentType, ok := filter["content_type"]; ok && contentType != nil {
			model = model.Where("JSON_EXTRACT(extend, '$.type') = ?", contentType)
		}

		// 是否推荐查询
		if isRecommended, ok := filter["is_recommended"]; ok {
			model = model.Where("is_recommended", isRecommended)
		}

		// 客户ID查询
		if clientId, ok := filter["client_id"]; ok && clientId != nil {
			model = model.Where("client_id", clientId)
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

	// 排序规则：同时满足is_recommended=1和置顶时间未过期才优先显示，再按创建时间倒序
	model = model.Order("CASE WHEN is_recommended = 1 AND top_until IS NOT NULL AND top_until > NOW() THEN 0 ELSE 1 END") // 有效置顶内容优先
	model = model.Order("created_at DESC")                                                                                // 最后按创建时间倒序

	// 查询数据
	var contents []*do.ContentDO
	err = model.Scan(&contents)
	if err != nil {
		return nil, 0, err
	}

	return contents, count, nil
}

// UpdateStatus 更新状态
func (d *contentDao) UpdateStatus(ctx context.Context, id interface{}, status string) (rowsAffected int64, err error) {
	data := do.ContentDO{
		Status:    status,
		UpdatedAt: gtime.Now(),
	}

	// 如果状态是已发布，则设置发布时间
	if status == "已发布" {
		data.PublishedAt = gtime.Now()
	}

	result, err := g.DB().Model(do.TableContent).Ctx(ctx).
		Data(data).
		Where("id", id).
		Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateRecommend 更新推荐状态
func (d *contentDao) UpdateRecommend(ctx context.Context, id interface{}, isRecommended bool, topUntil *gtime.Time) (rowsAffected int64, err error) {
	result, err := g.DB().Model(do.TableContent).Ctx(ctx).
		Data(g.Map{
			"is_recommended": isRecommended,
			"top_until":      topUntil,
			"updated_at":     gtime.Now(),
		}).
		Where("id", id).
		Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// IncrementViews 增加浏览量
func (d *contentDao) IncrementViews(ctx context.Context, id interface{}) error {
	_, err := g.DB().Model(do.TableContent).Ctx(ctx).
		Data("views=views+1").
		Where("id", id).
		Update()
	return err
}

// FindByIds 根据ID列表查询多条数据
func (d *contentDao) FindByIds(ctx context.Context, ids []interface{}) ([]*do.ContentDO, error) {
	if len(ids) == 0 {
		return []*do.ContentDO{}, nil
	}

	var contents []*do.ContentDO
	err := g.DB().Model(do.TableContent).Ctx(ctx).
		WhereIn("id", ids).
		Scan(&contents)

	if err != nil {
		return nil, err
	}

	return contents, nil
}

// GetContentInfo 获取内容信息
func (d *contentDao) GetContentInfo(ctx context.Context, tableName string, contentId int, result interface{}) error {
	err := g.DB().Model(tableName).
		Fields("id, title, cover, status").
		Where("id", contentId).
		Scan(result)
	if err != nil {
		return gerror.New("获取内容信息失败")
	}
	return nil
}
