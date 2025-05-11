package content

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	v1 "demo/api/content/v1"
	packagev1 "demo/api/package/v1"
	"demo/internal/dao"
	"demo/internal/model/do"
	"demo/internal/service"
)

type sContent struct{}

// New 创建内容服务实例
func New() service.ContentService {
	return &sContent{}
}

// List 内容列表
func (s *sContent) List(ctx context.Context, req *v1.ContentListReq) (res *v1.ContentListRes, err error) {
	res = &v1.ContentListRes{
		List:  make([]v1.ContentListItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 构建过滤条件
	filter := make(map[string]interface{})
	if req.Title != "" {
		filter["title"] = req.Title
	}
	if req.Category != "" {
		filter["category"] = req.Category
	}
	if req.Status != "" {
		filter["status"] = req.Status
	}
	if req.Author != "" {
		filter["author"] = req.Author
	}

	// 查询数据
	contentDao := dao.NewContentDao()
	list, total, err := contentDao.FindList(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		return nil, gerror.New("获取内容列表失败: " + err.Error())
	}

	// 设置响应
	res.Total = total
	res.Page = req.Page

	// 转换数据
	for _, item := range list {
		statusText := ""
		switch gconv.String(item.Status) {
		case "已发布":
			statusText = "已发布"
		case "待审核":
			statusText = "待审核"
		case "已下架":
			statusText = "已下架"
		default:
			statusText = gconv.String(item.Status)
		}

		// 转换日期格式
		var publishedAt, expiresAt, topUntil *string
		if item.PublishedAt != nil {
			tmp := item.PublishedAt.String()
			publishedAt = &tmp
		}
		if item.ExpiresAt != nil {
			tmp := item.ExpiresAt.String()
			expiresAt = &tmp
		}
		if item.TopUntil != nil {
			tmp := item.TopUntil.String()
			topUntil = &tmp
		}

		// 判断是否推荐（置顶）
		isRecommended := gconv.Bool(item.IsRecommended)
		// 如果置顶时间已过期，则不再显示为推荐状态
		if item.TopUntil != nil && !item.TopUntil.After(gtime.Now()) {
			isRecommended = false
		}

		listItem := v1.ContentListItem{
			Id:            gconv.Int(item.Id),
			Title:         gconv.String(item.Title),
			Category:      gconv.String(item.Category),
			Author:        gconv.String(item.Author),
			Status:        gconv.String(item.Status),
			StatusText:    statusText,
			Views:         gconv.Int(item.Views),
			Likes:         gconv.Int(item.Likes),
			Comments:      gconv.Int(item.Comments),
			IsRecommended: isRecommended,
			PublishedAt:   publishedAt,
			ExpiresAt:     expiresAt,
			TopUntil:      topUntil,
			CreatedAt:     item.CreatedAt.String(),
			UpdatedAt:     item.UpdatedAt.String(),
		}
		res.List = append(res.List, listItem)
	}

	return res, nil
}

// Detail 内容详情
func (s *sContent) Detail(ctx context.Context, req *v1.ContentDetailReq) (res *v1.ContentDetailRes, err error) {
	res = &v1.ContentDetailRes{}

	// 查询数据
	contentDao := dao.NewContentDao()
	content, err := contentDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("获取内容详情失败: " + err.Error())
	}
	if content == nil {
		return nil, gerror.New("内容不存在")
	}

	// 转换日期格式
	var publishedAt, expiresAt, topUntil *string
	if content.PublishedAt != nil {
		tmp := content.PublishedAt.String()
		publishedAt = &tmp
	}
	if content.ExpiresAt != nil {
		tmp := content.ExpiresAt.String()
		expiresAt = &tmp
	}
	if content.TopUntil != nil {
		tmp := content.TopUntil.String()
		topUntil = &tmp
	}

	// 判断是否推荐（置顶）
	isRecommended := gconv.Bool(content.IsRecommended)
	// 如果置顶时间已过期，则不再显示为推荐状态
	if content.TopUntil != nil && !content.TopUntil.After(gtime.Now()) {
		isRecommended = false
	}

	// 设置响应
	res.Id = gconv.Int(content.Id)
	res.Title = gconv.String(content.Title)
	res.Category = gconv.String(content.Category)
	res.Author = gconv.String(content.Author)
	res.Content = gconv.String(content.Content)
	res.Status = gconv.String(content.Status)
	res.Views = gconv.Int(content.Views)
	res.Likes = gconv.Int(content.Likes)
	res.Comments = gconv.Int(content.Comments)
	res.IsRecommended = isRecommended
	res.PublishedAt = publishedAt
	res.ExpiresAt = expiresAt
	res.TopUntil = topUntil
	res.CreatedAt = content.CreatedAt.String()
	res.UpdatedAt = content.UpdatedAt.String()

	// 设置默认内容类型为普通信息
	res.Type = 1

	// 解析扩展字段，提取闲置物品特有信息
	if content.Extend != nil && gconv.String(content.Extend) != "" {
		var extendMap map[string]interface{}
		if err := json.Unmarshal([]byte(gconv.String(content.Extend)), &extendMap); err == nil {
			// 提取内容类型
			if contentType, ok := extendMap["type"]; ok {
				res.Type = gconv.Int(contentType)
			}

			// 如果是闲置物品，提取特有字段
			if res.Type == 2 {
				// 提取价格
				if price, ok := extendMap["price"]; ok {
					res.Price = gconv.Float64(price)
				}
				// 提取原价
				if originalPrice, ok := extendMap["originalPrice"]; ok {
					res.OriginalPrice = gconv.Float64(originalPrice)
				}
				// 提取交易地点
				if tradePlace, ok := extendMap["tradePlace"]; ok {
					res.TradePlace = gconv.String(tradePlace)
				}
				// 提取交易方式
				if tradeMethod, ok := extendMap["tradeMethod"]; ok {
					res.TradeMethod = gconv.String(tradeMethod)
				}
			}
		}
	}

	return res, nil
}

// Create 创建内容
func (s *sContent) Create(ctx context.Context, req *v1.ContentCreateReq) (res *v1.ContentCreateRes, err error) {
	res = &v1.ContentCreateRes{
		HomeCategories: make([]v1.CategoryItem, 0),
		IdleCategories: make([]v1.CategoryItem, 0),
	}

	// 构建数据
	data := &do.ContentDO{
		Title:         req.Title,
		Category:      req.Category,
		Author:        req.Author,
		Content:       req.Content,
		Status:        req.Status,
		Views:         0,
		Likes:         0,
		Comments:      0,
		IsRecommended: req.IsRecommended,
		CreatedAt:     gtime.Now(),
		UpdatedAt:     gtime.Now(),
	}

	// 处理日期字段
	if req.Status == "已发布" {
		data.PublishedAt = gtime.Now()
	}

	// 处理套餐逻辑 - 发布套餐
	if req.PublishPackageId > 0 {
		// 获取发布套餐详情
		packageDetail, err := service.Package().Detail(ctx, &packagev1.PackageDetailReq{Id: req.PublishPackageId})
		if err != nil {
			g.Log().Warning(ctx, "获取发布套餐详情失败:", err.Error())
		} else if packageDetail != nil && packageDetail.Package != nil && packageDetail.Package.Type == "publish" {
			// 设置到期时间，根据时长单位计算
			var expiresAt *gtime.Time
			switch packageDetail.Package.DurationType {
			case "hour":
				expiresAt = gtime.Now().Add(time.Hour * time.Duration(packageDetail.Package.Duration))
			case "day":
				expiresAt = gtime.Now().AddDate(0, 0, packageDetail.Package.Duration)
			case "month":
				expiresAt = gtime.Now().AddDate(0, packageDetail.Package.Duration, 0)
			default:
				// 默认按天计算
				expiresAt = gtime.Now().AddDate(0, 0, packageDetail.Package.Duration)
			}
			data.ExpiresAt = expiresAt
			g.Log().Info(ctx, "使用发布套餐设置展示期限:", expiresAt.String())
		}
	} else if req.ExpiresAt != "" {
		// 用户手动指定了到期时间
		expiresAt, err := gtime.StrToTime(req.ExpiresAt)
		if err != nil {
			return nil, gerror.New("到期时间格式不正确: " + err.Error())
		}
		data.ExpiresAt = expiresAt
	} else {
		// 默认设置3天展示期
		data.ExpiresAt = gtime.Now().AddDate(0, 0, 3)
		g.Log().Info(ctx, "设置默认3天展示期限:", data.ExpiresAt.String())
	}

	// 处理套餐逻辑 - 置顶套餐
	if req.TopPackageId > 0 && req.IsRecommended {
		// 获取置顶套餐详情
		packageDetail, err := service.Package().Detail(ctx, &packagev1.PackageDetailReq{Id: req.TopPackageId})
		if err != nil {
			g.Log().Warning(ctx, "获取置顶套餐详情失败:", err.Error())
		} else if packageDetail != nil && packageDetail.Package != nil && packageDetail.Package.Type == "top" {
			// 设置置顶截止时间，根据时长单位计算
			var topUntil *gtime.Time
			switch packageDetail.Package.DurationType {
			case "hour":
				topUntil = gtime.Now().Add(time.Hour * time.Duration(packageDetail.Package.Duration))
			case "day":
				topUntil = gtime.Now().AddDate(0, 0, packageDetail.Package.Duration)
			case "month":
				topUntil = gtime.Now().AddDate(0, packageDetail.Package.Duration, 0)
			default:
				// 默认按天计算
				topUntil = gtime.Now().AddDate(0, 0, packageDetail.Package.Duration)
			}

			// 校验置顶时间不能超过到期时间
			if data.ExpiresAt != nil && topUntil.After(data.ExpiresAt) {
				topUntil = data.ExpiresAt
				g.Log().Info(ctx, "置顶时间超过到期时间，已调整为到期时间:", topUntil.String())
			}

			data.TopUntil = topUntil
			g.Log().Info(ctx, "使用置顶套餐设置置顶期限:", topUntil.String())
		}
	} else if req.TopUntil != "" && req.IsRecommended {
		// 用户手动指定了置顶截止时间
		topUntil, err := gtime.StrToTime(req.TopUntil)
		if err != nil {
			return nil, gerror.New("置顶截止时间格式不正确: " + err.Error())
		}

		// 校验置顶时间不能超过到期时间
		if data.ExpiresAt != nil && topUntil.After(data.ExpiresAt) {
			return nil, gerror.New("置顶时间不能超过到期时间")
		}

		data.TopUntil = topUntil
	}

	// 插入数据
	contentDao := dao.NewContentDao()
	lastInsertId, err := contentDao.Insert(ctx, data)
	if err != nil {
		return nil, gerror.New("创建内容失败: " + err.Error())
	}

	res.Id = int(lastInsertId)

	// 获取首页分类列表
	homeCategoryDao := dao.NewHomeCategoryDao()
	homeCategories, err := homeCategoryDao.FindList(ctx)
	if err != nil {
		return res, nil // 返回创建的ID，忽略分类获取错误
	}

	// 转换首页分类数据
	for _, item := range homeCategories {
		categoryItem := v1.CategoryItem{
			Id:        gconv.Int(item.Id),
			Name:      gconv.String(item.Name),
			SortOrder: gconv.Int(item.SortOrder),
			IsActive:  gconv.Bool(item.IsActive),
			CreatedAt: item.CreatedAt.String(),
			UpdatedAt: item.UpdatedAt.String(),
		}
		res.HomeCategories = append(res.HomeCategories, categoryItem)
	}

	// 获取闲置分类列表
	idleCategoryDao := dao.NewIdleCategoryDao()
	idleCategories, err := idleCategoryDao.FindList(ctx)
	if err != nil {
		return res, nil // 返回创建的ID和首页分类，忽略闲置分类获取错误
	}

	// 转换闲置分类数据
	for _, item := range idleCategories {
		categoryItem := v1.CategoryItem{
			Id:        gconv.Int(item.Id),
			Name:      gconv.String(item.Name),
			SortOrder: gconv.Int(item.SortOrder),
			IsActive:  gconv.Bool(item.IsActive),
			CreatedAt: item.CreatedAt.String(),
			UpdatedAt: item.UpdatedAt.String(),
		}
		res.IdleCategories = append(res.IdleCategories, categoryItem)
	}

	return res, nil
}

// Update 更新内容
func (s *sContent) Update(ctx context.Context, req *v1.ContentUpdateReq) (res *v1.ContentUpdateRes, err error) {
	res = &v1.ContentUpdateRes{}

	// 查询原有数据
	contentDao := dao.NewContentDao()
	exists, err := contentDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询内容失败: " + err.Error())
	}
	if exists == nil {
		return nil, gerror.New("内容不存在")
	}

	// 构建更新数据
	data := &do.ContentDO{
		Title:         req.Title,
		Category:      req.Category,
		Author:        req.Author,
		Content:       req.Content,
		Status:        req.Status,
		IsRecommended: req.IsRecommended,
		UpdatedAt:     gtime.Now(),
	}

	// 处理日期字段
	if req.Status == "已发布" && (exists.PublishedAt == nil || gconv.String(exists.Status) != "已发布") {
		data.PublishedAt = gtime.Now()
	}

	// 处理套餐逻辑 - 发布套餐
	if req.PublishPackageId > 0 {
		// 获取发布套餐详情
		packageDetail, err := service.Package().Detail(ctx, &packagev1.PackageDetailReq{Id: req.PublishPackageId})
		if err != nil {
			g.Log().Warning(ctx, "获取发布套餐详情失败:", err.Error())
		} else if packageDetail != nil && packageDetail.Package != nil && packageDetail.Package.Type == "publish" {
			// 设置到期时间，根据时长单位计算
			var expiresAt *gtime.Time
			switch packageDetail.Package.DurationType {
			case "hour":
				expiresAt = gtime.Now().Add(time.Hour * time.Duration(packageDetail.Package.Duration))
			case "day":
				expiresAt = gtime.Now().AddDate(0, 0, packageDetail.Package.Duration)
			case "month":
				expiresAt = gtime.Now().AddDate(0, packageDetail.Package.Duration, 0)
			default:
				// 默认按天计算
				expiresAt = gtime.Now().AddDate(0, 0, packageDetail.Package.Duration)
			}
			data.ExpiresAt = expiresAt
			g.Log().Info(ctx, "使用发布套餐设置展示期限:", expiresAt.String())
		}
	} else if req.ExpiresAt != "" {
		// 用户手动指定了到期时间
		expiresAt, err := gtime.StrToTime(req.ExpiresAt)
		if err != nil {
			return nil, gerror.New("到期时间格式不正确: " + err.Error())
		}
		data.ExpiresAt = expiresAt
	}

	// 处理套餐逻辑 - 置顶套餐
	if req.TopPackageId > 0 && req.IsRecommended {
		// 获取置顶套餐详情
		packageDetail, err := service.Package().Detail(ctx, &packagev1.PackageDetailReq{Id: req.TopPackageId})
		if err != nil {
			g.Log().Warning(ctx, "获取置顶套餐详情失败:", err.Error())
		} else if packageDetail != nil && packageDetail.Package != nil && packageDetail.Package.Type == "top" {
			// 设置置顶截止时间为当前时间加上套餐天数
			topUntil := gtime.Now().AddDate(0, 0, packageDetail.Package.Duration)

			// 校验置顶时间不能超过到期时间
			if data.ExpiresAt != nil && topUntil.After(data.ExpiresAt) {
				topUntil = data.ExpiresAt
				g.Log().Info(ctx, "置顶时间超过到期时间，已调整为到期时间:", topUntil.String())
			}

			data.TopUntil = topUntil
			g.Log().Info(ctx, "使用置顶套餐设置置顶期限:", topUntil.String())
		}
	} else if req.TopUntil != "" && req.IsRecommended {
		// 用户手动指定了置顶截止时间
		topUntil, err := gtime.StrToTime(req.TopUntil)
		if err != nil {
			return nil, gerror.New("置顶截止时间格式不正确: " + err.Error())
		}

		// 校验置顶时间不能超过到期时间
		if data.ExpiresAt != nil && topUntil.After(data.ExpiresAt) {
			return nil, gerror.New("置顶时间不能超过到期时间")
		}

		data.TopUntil = topUntil
	} else if !req.IsRecommended {
		// 如果取消推荐，则置顶时间置为空
		data.TopUntil = nil
	}

	// 更新数据
	_, err = contentDao.Update(ctx, data, req.Id)
	if err != nil {
		return nil, gerror.New("更新内容失败: " + err.Error())
	}

	return res, nil
}

// Delete 删除内容
func (s *sContent) Delete(ctx context.Context, req *v1.ContentDeleteReq) (res *v1.ContentDeleteRes, err error) {
	res = &v1.ContentDeleteRes{}

	// 查询原有数据
	contentDao := dao.NewContentDao()
	exists, err := contentDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询内容失败: " + err.Error())
	}
	if exists == nil {
		return nil, gerror.New("内容不存在")
	}

	// 删除数据（硬删除）
	_, err = contentDao.Delete(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("删除内容失败: " + err.Error())
	}

	return res, nil
}

// UpdateStatus 更新内容状态
func (s *sContent) UpdateStatus(ctx context.Context, req *v1.ContentStatusUpdateReq) (res *v1.ContentStatusUpdateRes, err error) {
	res = &v1.ContentStatusUpdateRes{}

	// 查询原有数据
	contentDao := dao.NewContentDao()
	exists, err := contentDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询内容失败: " + err.Error())
	}
	if exists == nil {
		return nil, gerror.New("内容不存在")
	}

	// 更新状态
	_, err = contentDao.UpdateStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, gerror.New("更新内容状态失败: " + err.Error())
	}

	return res, nil
}

// UpdateRecommend 更新内容推荐状态
func (s *sContent) UpdateRecommend(ctx context.Context, req *v1.ContentRecommendUpdateReq) (res *v1.ContentRecommendUpdateRes, err error) {
	res = &v1.ContentRecommendUpdateRes{}

	// 查询原有数据
	contentDao := dao.NewContentDao()
	exists, err := contentDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询内容失败: " + err.Error())
	}
	if exists == nil {
		return nil, gerror.New("内容不存在")
	}

	// 处理置顶时间
	var topUntil *gtime.Time = nil
	if req.IsRecommended {
		// 如果使用置顶套餐
		if req.TopPackageId > 0 {
			// 获取套餐信息
			packageDao := dao.NewPackageDao()
			packageInfo, err := packageDao.FindOne(ctx, req.TopPackageId)
			if err != nil {
				g.Log().Error(ctx, "获取置顶套餐信息失败", err)
				return nil, gerror.New("获取置顶套餐信息失败: " + err.Error())
			}
			if packageInfo == nil {
				return nil, gerror.New("置顶套餐不存在")
			}
			if packageInfo.Type != "top" {
				return nil, gerror.New("套餐类型错误，必须使用置顶套餐")
			}

			// 根据时长单位设置置顶截止时间
			durationType := gconv.String(packageInfo.DurationType)
			duration := gconv.Int(packageInfo.Duration)

			switch durationType {
			case "hour":
				topUntil = gtime.Now().Add(time.Hour * time.Duration(duration))
			case "day":
				topUntil = gtime.Now().AddDate(0, 0, duration)
			case "month":
				topUntil = gtime.Now().AddDate(0, duration, 0)
			default:
				// 默认按天计算
				topUntil = gtime.Now().AddDate(0, 0, duration)
			}

			// 校验置顶时间不能超过到期时间
			if exists.ExpiresAt != nil && topUntil.After(exists.ExpiresAt) {
				topUntil = exists.ExpiresAt
				g.Log().Info(ctx, "置顶时间超过到期时间，已调整为到期时间:", topUntil.String())
			}
		} else if req.TopUntil != "" {
			// 手动设置置顶时间
			tmp, err := gtime.StrToTime(req.TopUntil)
			if err != nil {
				return nil, gerror.New("置顶截止时间格式不正确: " + err.Error())
			}

			// 校验置顶时间不能超过到期时间
			if exists.ExpiresAt != nil && tmp.After(exists.ExpiresAt) {
				return nil, gerror.New("置顶时间不能超过到期时间")
			}

			topUntil = tmp
		}
	}

	// 更新推荐状态
	_, err = contentDao.UpdateRecommend(ctx, req.Id, req.IsRecommended, topUntil)
	if err != nil {
		return nil, gerror.New("更新内容推荐状态失败: " + err.Error())
	}

	return res, nil
}

// GetCategories 获取所有分类
func (s *sContent) GetCategories(ctx context.Context, req *v1.GetCategoriesReq) (res *v1.GetCategoriesRes, err error) {
	res = &v1.GetCategoriesRes{
		HomeCategories: make([]v1.CategoryItem, 0),
		IdleCategories: make([]v1.CategoryItem, 0),
	}

	// 获取首页分类列表
	homeCategoryDao := dao.NewHomeCategoryDao()
	homeCategories, err := homeCategoryDao.FindList(ctx)
	if err != nil {
		return nil, gerror.New("获取首页分类列表失败: " + err.Error())
	}

	// 转换首页分类数据
	for _, item := range homeCategories {
		categoryItem := v1.CategoryItem{
			Id:        gconv.Int(item.Id),
			Name:      gconv.String(item.Name),
			SortOrder: gconv.Int(item.SortOrder),
			IsActive:  gconv.Bool(item.IsActive),
			CreatedAt: item.CreatedAt.String(),
			UpdatedAt: item.UpdatedAt.String(),
		}
		res.HomeCategories = append(res.HomeCategories, categoryItem)
	}

	// 获取闲置分类列表
	idleCategoryDao := dao.NewIdleCategoryDao()
	idleCategories, err := idleCategoryDao.FindList(ctx)
	if err != nil {
		return nil, gerror.New("获取闲置分类列表失败: " + err.Error())
	}

	// 转换闲置分类数据
	for _, item := range idleCategories {
		categoryItem := v1.CategoryItem{
			Id:        gconv.Int(item.Id),
			Name:      gconv.String(item.Name),
			SortOrder: gconv.Int(item.SortOrder),
			IsActive:  gconv.Bool(item.IsActive),
			CreatedAt: item.CreatedAt.String(),
			UpdatedAt: item.UpdatedAt.String(),
		}
		res.IdleCategories = append(res.IdleCategories, categoryItem)
	}

	return res, nil
}
