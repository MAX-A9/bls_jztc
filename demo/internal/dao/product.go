package dao

import (
	"context"
	"demo/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// 产品数据访问对象
type ProductDao struct{}

// 获取商品列表
func (d *ProductDao) GetList(ctx context.Context, filter *model.ProductFilter, page, size int) (list []*model.Product, total int, err error) {
	// 初始化查询构建器
	m := g.DB().Model("product")

	// 过滤条件
	if filter != nil {
		// 商品名称模糊搜索
		if filter.Name != "" {
			m = m.WhereLike("name", "%"+filter.Name+"%")
		}
		// 根据分类ID筛选
		if filter.CategoryId > 0 {
			m = m.Where("category_id", filter.CategoryId)
		}
		// 根据状态筛选，只有当状态值>=0且<=2时才添加状态过滤
		if filter.Status >= 0 && filter.Status <= 2 {
			m = m.Where("status", filter.Status)
		}
		// 最小时长
		if filter.MinDuration > 0 {
			m = m.WhereGTE("duration", filter.MinDuration)
		}
		// 最大时长
		if filter.MaxDuration > 0 {
			m = m.WhereLTE("duration", filter.MaxDuration)
		}
		// 最小库存
		if filter.MinStock > 0 {
			m = m.WhereGTE("stock", filter.MinStock)
		}
		// 根据标签筛选
		if filter.Tags != "" {
			m = m.WhereLike("tags", "%"+filter.Tags+"%")
		}
	}

	// 获取总记录数
	total, err = m.Count()
	if err != nil {
		return nil, 0, err
	}

	// 排序处理
	orderBy := "sort_order asc" // 默认排序
	if filter != nil && filter.SortField != "" {
		sortField := filter.SortField
		sortOrder := "asc"
		if filter.SortOrder == "desc" {
			sortOrder = "desc"
		}
		orderBy = sortField + " " + sortOrder
	}

	// 分页查询
	err = m.Page(page, size).Order(orderBy).Scan(&list)
	return list, total, err
}

// 根据ID获取商品详情
func (d *ProductDao) GetById(ctx context.Context, id int) (*model.Product, error) {
	var product *model.Product
	err := g.DB().Model("product").Where("id", id).Scan(&product)
	return product, err
}

// 创建商品
func (d *ProductDao) Create(ctx context.Context, data *model.Product) (int64, error) {
	result, err := g.DB().Model("product").Data(data).Insert()
	if err != nil {
		return 0, err
	}

	// 更新商品分类中的商品数量
	if data.CategoryId > 0 {
		shopCategoryDao := &ShopCategoryDao{}
		err = shopCategoryDao.SyncCategoryProductCount(ctx, data.CategoryId)
		if err != nil {
			g.Log().Warning(ctx, "更新商品分类数量失败:", err)
		}
	}

	return result.LastInsertId()
}

// 更新商品信息
func (d *ProductDao) Update(ctx context.Context, id int, data g.Map) error {
	// 查询商品原有分类信息
	var oldCategoryId int
	oldProduct, _ := d.GetById(ctx, id)
	if oldProduct != nil {
		oldCategoryId = oldProduct.CategoryId
	}

	// 更新商品信息
	_, err := g.DB().Model("product").Where("id", id).Data(data).Update()
	if err != nil {
		return err
	}

	// 更新商品分类中的商品数量
	shopCategoryDao := &ShopCategoryDao{}

	// 如果更新了分类ID，则需要更新新旧两个分类的商品数量
	if data["category_id"] != nil && oldCategoryId > 0 && data["category_id"] != oldCategoryId {
		// 更新旧分类
		err = shopCategoryDao.SyncCategoryProductCount(ctx, oldCategoryId)
		if err != nil {
			g.Log().Warning(ctx, "更新旧商品分类数量失败:", err)
		}

		// 更新新分类
		newCategoryId := data["category_id"].(int)
		err = shopCategoryDao.SyncCategoryProductCount(ctx, newCategoryId)
		if err != nil {
			g.Log().Warning(ctx, "更新新商品分类数量失败:", err)
		}
	} else if oldCategoryId > 0 {
		// 只更新当前分类
		err = shopCategoryDao.SyncCategoryProductCount(ctx, oldCategoryId)
		if err != nil {
			g.Log().Warning(ctx, "更新商品分类数量失败:", err)
		}
	}

	return nil
}

// 删除商品
func (d *ProductDao) Delete(ctx context.Context, id int) error {
	// 查询商品分类信息
	var categoryId int
	product, _ := d.GetById(ctx, id)
	if product != nil {
		categoryId = product.CategoryId
	}

	// 删除商品
	_, err := g.DB().Model("product").Where("id", id).Delete()
	if err != nil {
		return err
	}

	// 更新商品分类中的商品数量
	if categoryId > 0 {
		shopCategoryDao := &ShopCategoryDao{}
		err = shopCategoryDao.SyncCategoryProductCount(ctx, categoryId)
		if err != nil {
			g.Log().Warning(ctx, "更新商品分类数量失败:", err)
		}
	}

	return nil
}

// 更新商品状态
func (d *ProductDao) UpdateStatus(ctx context.Context, id int, status int) error {
	_, err := g.DB().Model("product").Where("id", id).Data(g.Map{"status": status}).Update()
	return err
}

// 检查商品是否存在
func (d *ProductDao) Exists(ctx context.Context, id int) (bool, error) {
	count, err := g.DB().Model("product").Where("id", id).Count()
	return count > 0, err
}

// 获取微信小程序端商品列表（只获取已上架的商品）
func (d *ProductDao) GetWxList(ctx context.Context, filter *model.ProductFilter, page, size int) (list []*model.Product, total int, err error) {
	// 初始化查询构建器
	m := g.DB().Model("product").Where("status", model.ProductStatusListed).Where("stock > 0")

	// 过滤条件
	if filter != nil {
		// 商品名称模糊搜索
		if filter.Name != "" {
			m = m.WhereLike("name", "%"+filter.Name+"%")
		}
		// 根据分类ID筛选
		if filter.CategoryId > 0 {
			m = m.Where("category_id", filter.CategoryId)
		}
		// 根据时长筛选
		if filter.MinDuration > 0 {
			m = m.WhereGTE("duration", filter.MinDuration)
		}
		if filter.MaxDuration > 0 {
			m = m.WhereLTE("duration", filter.MaxDuration)
		}
		// 根据标签筛选
		if filter.Tags != "" {
			m = m.WhereLike("tags", "%"+filter.Tags+"%")
		}
	}

	// 获取总记录数
	total, err = m.Count()
	if err != nil {
		return nil, 0, err
	}

	// 排序处理
	orderBy := "sort_order asc" // 默认排序
	if filter != nil && filter.SortField != "" {
		sortField := filter.SortField
		sortOrder := "asc"
		if filter.SortOrder == "desc" {
			sortOrder = "desc"
		}
		orderBy = sortField + " " + sortOrder
	}

	// 分页查询
	err = m.Page(page, size).Order(orderBy).Scan(&list)
	return list, total, err
}
