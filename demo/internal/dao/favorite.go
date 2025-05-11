package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"demo/internal/model/do"
)

// FavoriteDao 收藏数据访问对象
type FavoriteDao struct{}

// NewFavoriteDao 创建收藏数据访问对象
func NewFavoriteDao() *FavoriteDao {
	return &FavoriteDao{}
}

// Add 添加收藏
func (d *FavoriteDao) Add(ctx context.Context, clientId int, contentId int) error {
	// 开启事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 检查是否已经收藏
		count, err := tx.Model(do.TableFavorite).
			Where("client_id=? AND content_id=?", clientId, contentId).
			Count()
		if err != nil {
			return err
		}
		if count > 0 {
			// 已经收藏过了，不需要添加
			return nil
		}

		// 2. 添加收藏记录
		_, err = tx.Model(do.TableFavorite).
			Data(do.FavoriteDO{
				ClientId:  clientId,
				ContentId: contentId,
				CreatedAt: gtime.Now(),
			}).
			Insert()
		if err != nil {
			return err
		}

		// 3. 更新内容的likes数量+1
		_, err = tx.Model(do.TableContent).
			Where("id=?", contentId).
			Increment("likes", 1)

		return err
	})
}

// Cancel 取消收藏
func (d *FavoriteDao) Cancel(ctx context.Context, clientId int, contentId int) error {
	// 开启事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 检查是否已经收藏
		count, err := tx.Model(do.TableFavorite).
			Where("client_id=? AND content_id=?", clientId, contentId).
			Count()
		if err != nil {
			return err
		}
		if count == 0 {
			// 未收藏过，不需要取消
			return nil
		}

		// 2. 直接删除收藏记录
		_, err = tx.Model(do.TableFavorite).
			Where("client_id=? AND content_id=?", clientId, contentId).
			Delete()
		if err != nil {
			return err
		}

		// 3. 更新内容的likes数量-1
		_, err = tx.Model(do.TableContent).
			Where("id=?", contentId).
			Decrement("likes", 1)

		return err
	})
}

// IsFavorite 检查用户是否已收藏内容
func (d *FavoriteDao) IsFavorite(ctx context.Context, clientId int, contentId int) (bool, error) {
	count, err := g.DB().Model(do.TableFavorite).
		Where("client_id=? AND content_id=?", clientId, contentId).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetFavoriteList 获取用户收藏列表
func (d *FavoriteDao) GetFavoriteList(ctx context.Context, clientId int, page, pageSize int, contentType int, categoryId int) ([]map[string]interface{}, int, error) {
	db := g.DB()

	// 构建基础查询
	query := db.Model(do.TableFavorite+" AS f").
		LeftJoin(do.TableContent+" AS c", "f.content_id=c.id").
		Where("f.client_id=?", clientId)

	// 按内容类型筛选（首页分类或闲置分类）
	if contentType > 0 {
		// 首页分类为普通内容，闲置分类为闲置物品
		// 从extend字段中提取type字段
		if contentType == 1 {
			// 普通信息
			query = query.Where("JSON_EXTRACT(c.extend, '$.type') = 1 OR JSON_EXTRACT(c.extend, '$.type') IS NULL")
		} else if contentType == 2 {
			// 闲置物品
			query = query.Where("JSON_EXTRACT(c.extend, '$.type') = 2")
		}
	}

	// 按具体分类ID筛选
	if categoryId > 0 {
		// 首页分类在content.category字段中
		// 闲置分类在extend字段的categoryId中
		if contentType == 1 {
			// 首页分类
			// 需要查询首页分类的名称，然后按名称筛选
			var categoryName string
			err := db.Model("home_category").
				Where("id", categoryId).
				Fields("name").
				Scan(&categoryName)
			if err == nil && categoryName != "" {
				query = query.Where("c.category=?", categoryName)
			}
		} else if contentType == 2 {
			// 闲置分类
			query = query.Where("JSON_EXTRACT(c.extend, '$.categoryId') = ?", categoryId)
		} else {
			// 全部类型，查询首页分类名称
			var categoryName string
			err := db.Model("home_category").
				Where("id", categoryId).
				Fields("name").
				Scan(&categoryName)

			if err == nil && categoryName != "" {
				// 首页分类查询
				query = query.Where("(c.category=? OR JSON_EXTRACT(c.extend, '$.categoryId') = ?)",
					categoryName, categoryId)
			} else {
				// 仅按闲置分类ID查询
				query = query.Where("JSON_EXTRACT(c.extend, '$.categoryId') = ?", categoryId)
			}
		}
	}

	// 计算总数 - 使用简单计数的方式
	total, err := query.Clone().Count()
	if err != nil {
		return nil, 0, err
	}

	// 查询列表数据
	list, err := query.Clone().
		Fields("c.*").
		Order("f.created_at DESC").
		Page(page, pageSize).
		All()

	if err != nil {
		return nil, 0, err
	}

	return list.List(), total, nil
}

// GetFavoriteCount 获取用户收藏总数
func (d *FavoriteDao) GetFavoriteCount(ctx context.Context, clientId int) (int, error) {
	db := g.DB()

	// 构建基础查询
	query := db.Model(do.TableFavorite+" AS f").
		LeftJoin(do.TableContent+" AS c", "f.content_id=c.id").
		Where("f.client_id=?", clientId)

	// 计算总数
	total, err := query.Count()
	if err != nil {
		return 0, err
	}

	return total, nil
}
