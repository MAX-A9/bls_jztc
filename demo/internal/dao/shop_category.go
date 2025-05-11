package dao

import (
	"context"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ShopCategoryDao 商城分类数据访问对象
type ShopCategoryDao struct{}

// ShopCategoryFilter 商城分类过滤条件
type ShopCategoryFilter struct {
	Name      string `json:"name"`      // 分类名称，模糊搜索
	Status    int    `json:"status"`    // 状态 -1:全部 0:禁用 1:启用
	SortField string `json:"sortField"` // 排序字段
	SortOrder string `json:"sortOrder"` // 排序方式: asc, desc
}

// Get 根据ID获取商城分类
func (d *ShopCategoryDao) Get(ctx context.Context, id int) (*entity.ShopCategory, error) {
	var category entity.ShopCategory
	err := g.DB().Model("shop_category").Where("id=?", id).Scan(&category)
	if err != nil || category.Id == 0 {
		return nil, err
	}
	return &category, nil
}

// Create 创建商城分类
func (d *ShopCategoryDao) Create(ctx context.Context, category *entity.ShopCategory) (int64, error) {
	now := gtime.Now()
	result, err := g.DB().Model("shop_category").
		Data(g.Map{
			"name":          category.Name,
			"sort_order":    category.SortOrder,
			"product_count": category.ProductCount,
			"status":        category.Status,
			"image":         category.Image,
			"created_at":    now,
			"updated_at":    now,
		}).
		Insert()
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

// Update 更新商城分类
func (d *ShopCategoryDao) Update(ctx context.Context, category *entity.ShopCategory) error {
	_, err := g.DB().Model("shop_category").
		Data(g.Map{
			"name":          category.Name,
			"sort_order":    category.SortOrder,
			"product_count": category.ProductCount,
			"status":        category.Status,
			"image":         category.Image,
			"updated_at":    gtime.Now(),
		}).
		Where("id=?", category.Id).
		Update()
	return err
}

// Delete 删除商城分类
func (d *ShopCategoryDao) Delete(ctx context.Context, id int) error {
	_, err := g.DB().Model("shop_category").
		Where("id=?", id).
		Delete()
	return err
}

// UpdateStatus 更新商城分类状态
func (d *ShopCategoryDao) UpdateStatus(ctx context.Context, id int, status int) error {
	_, err := g.DB().Model("shop_category").
		Data(g.Map{
			"status":     status,
			"updated_at": gtime.Now(),
		}).
		Where("id=?", id).
		Update()
	return err
}

// GetList 获取商城分类列表
func (d *ShopCategoryDao) GetList(ctx context.Context, filter *ShopCategoryFilter) ([]*entity.ShopCategory, error) {
	model := g.DB().Model("shop_category")

	// 过滤条件
	if filter != nil {
		// 分类名称模糊搜索
		if filter.Name != "" {
			model = model.WhereLike("name", "%"+filter.Name+"%")
		}
		// 状态过滤
		if filter.Status >= 0 && filter.Status <= 1 {
			model = model.Where("status=?", filter.Status)
		}
	}

	// 排序处理
	orderBy := "sort_order ASC"
	if filter != nil && filter.SortField != "" {
		sortField := filter.SortField
		sortOrder := "ASC"
		if filter.SortOrder == "desc" {
			sortOrder = "DESC"
		}
		orderBy = sortField + " " + sortOrder
	}

	var categories []*entity.ShopCategory
	err := model.Order(orderBy).Scan(&categories)
	return categories, err
}

// GetPage 分页获取商城分类列表
func (d *ShopCategoryDao) GetPage(ctx context.Context, page, size int, filter *ShopCategoryFilter) ([]entity.ShopCategory, int, error) {
	model := g.DB().Model("shop_category")

	// 过滤条件
	if filter != nil {
		// 分类名称模糊搜索
		if filter.Name != "" {
			model = model.WhereLike("name", "%"+filter.Name+"%")
		}
		// 状态过滤
		if filter.Status >= 0 && filter.Status <= 1 {
			model = model.Where("status=?", filter.Status)
		}
	}

	// 计算总数
	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}

	// 排序处理
	orderBy := "sort_order ASC"
	if filter != nil && filter.SortField != "" {
		sortField := filter.SortField
		sortOrder := "ASC"
		if filter.SortOrder == "desc" {
			sortOrder = "DESC"
		}
		orderBy = sortField + " " + sortOrder
	}

	// 获取分页数据
	var categories []entity.ShopCategory
	err = model.Page(page, size).
		Order(orderBy).
		Scan(&categories)

	return categories, total, err
}

// UpdateProductCount 更新分类的商品数量
func (d *ShopCategoryDao) UpdateProductCount(ctx context.Context, id int, count int) error {
	_, err := g.DB().Model("shop_category").
		Data(g.Map{
			"product_count": count,
			"updated_at":    gtime.Now(),
		}).
		Where("id=?", id).
		Update()
	return err
}

// GetProductCount 获取分类的商品数量
func (d *ShopCategoryDao) GetProductCount(ctx context.Context, id int) (int, error) {
	// 查询分类下商品数量（由于已移除deleted_at字段，不需要考虑软删除）
	count, err := g.DB().Model("product").Where("category_id=?", id).Count()
	return count, err
}

// UpdateAllCategoriesProductCount 更新所有分类的商品数量
func (d *ShopCategoryDao) UpdateAllCategoriesProductCount(ctx context.Context) error {
	// 获取所有分类
	categories, err := d.GetList(ctx, nil)
	if err != nil {
		return err
	}

	// 逐个更新分类的商品数量
	for _, category := range categories {
		count, err := d.GetProductCount(ctx, category.Id)
		if err != nil {
			return err
		}

		err = d.UpdateProductCount(ctx, category.Id, count)
		if err != nil {
			return err
		}
	}

	return nil
}

// SyncCategoryProductCount 同步指定分类的商品数量
func (d *ShopCategoryDao) SyncCategoryProductCount(ctx context.Context, categoryId int) error {
	count, err := d.GetProductCount(ctx, categoryId)
	if err != nil {
		return err
	}

	return d.UpdateProductCount(ctx, categoryId, count)
}
