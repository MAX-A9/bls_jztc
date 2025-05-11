package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	v1 "demo/api/content/client/v1"
	"demo/internal/consts"
	"demo/internal/dao"
	"demo/internal/model/do"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// contentClientImpl 客户端内容服务实现
type contentClientImpl struct {
	contentDao  dao.ContentDao
	categoryDao dao.ContentDao // 复用内容分类
}

// New 创建客户端内容服务实例
func New() service.ContentClientService {
	return &contentClientImpl{
		contentDao:  dao.NewContentDao(),
		categoryDao: dao.NewContentDao(),
	}
}

// ContentPublicDetail 获取公开内容详情
func (s *contentClientImpl) ContentPublicDetail(ctx context.Context, req *v1.ContentPublicDetailReq) (res *v1.ContentPublicDetailRes, err error) {
	res = &v1.ContentPublicDetailRes{}

	// 查询内容数据
	content, err := s.contentDao.FindOne(ctx, req.Id)
	if err != nil {
		g.Log().Error(ctx, "获取内容详情失败", err)
		return nil, gerror.New("获取内容详情失败: " + err.Error())
	}

	if content == nil {
		return nil, gerror.New("内容不存在")
	}

	// 检查内容状态是否为已发布
	if gconv.String(content.Status) != "已发布" {
		return nil, gerror.New("内容未发布或已下架")
	}

	// 更新浏览量
	err = s.contentDao.IncrementViews(ctx, req.Id)
	if err != nil {
		g.Log().Error(ctx, "更新浏览量失败", err)
		// 不返回错误，继续处理
	}

	// 从content中提取图片
	var images []string
	contentText := gconv.String(content.Content)

	// 提取图片，处理JSON格式图片数组
	if len(contentText) > 0 && strings.Contains(contentText, "[") && strings.Contains(contentText, "]") {
		startIdx := strings.Index(contentText, "[")
		endIdx := strings.LastIndex(contentText, "]")
		if startIdx >= 0 && endIdx > startIdx {
			jsonStr := contentText[startIdx : endIdx+1]
			if err := json.Unmarshal([]byte(jsonStr), &images); err == nil {
				// 成功提取图片，更新内容文本
				contentText = strings.TrimSpace(contentText[:startIdx])
			}
		}
	}

	// 如果没找到JSON格式图片，尝试从HTML内容中提取
	if len(images) == 0 {
		images = extractImagesFromHtml(contentText)
	}

	// 处理扩展字段
	var tradePlace, tradeMethod string
	var price, originalPrice float64
	var isTop bool

	if content.Extend != "" {
		var extendMap map[string]interface{}
		extendStr := gconv.String(content.Extend)
		if err := json.Unmarshal([]byte(extendStr), &extendMap); err == nil {
			// 获取价格信息
			if priceVal, ok := extendMap["price"]; ok {
				price = gconv.Float64(priceVal)
			}
			// 获取原价信息
			if originalPriceVal, ok := extendMap["originalPrice"]; ok {
				originalPrice = gconv.Float64(originalPriceVal)
			}
			// 获取交易地点
			if tradePlaceVal, ok := extendMap["tradePlace"]; ok {
				tradePlace = gconv.String(tradePlaceVal)
			}
			// 获取交易方式
			if tradeMethodVal, ok := extendMap["tradeMethod"]; ok {
				tradeMethod = gconv.String(tradeMethodVal)
			}
		}
	}

	// 判断是否置顶
	if content.TopUntil != nil && content.TopUntil.After(gtime.Now()) {
		isTop = true
	} else {
		isTop = false // 如果没有置顶时间或置顶时间已过期，则不置顶
	}

	// 获取发布者信息
	publisher := gconv.String(content.Author)
	clientId := gconv.Int(content.ClientId)
	if clientId > 0 {
		var clientInfo struct {
			RealName string `json:"real_name"`
		}
		_ = g.DB().Model("client").
			Fields("real_name").
			Where("id", clientId).
			Scan(&clientInfo)

		if clientInfo.RealName != "" {
			publisher = clientInfo.RealName // 优先使用真实姓名
		}
	}

	// 设置返回数据
	res.Id = gconv.Int(content.Id)
	res.Title = gconv.String(content.Title)
	res.Content = contentText
	res.Category = gconv.String(content.Category)
	res.PublishTime = content.PublishedAt.Format("Y-m-d H:i:s")
	res.Publisher = publisher
	res.PublisherId = uint(clientId) // 设置发布者ID
	res.IsTop = isTop
	res.TradePlace = tradePlace
	res.Price = price
	res.OriginalPrice = originalPrice
	res.Views = gconv.Int(content.Views)
	res.Likes = gconv.Int(content.Likes)
	res.TradeMethod = tradeMethod
	res.Images = images

	return res, nil
}

// HomeContentList 获取首页内容列表
func (s *contentClientImpl) HomeContentList(ctx context.Context, req *v1.HomeContentListReq) (res *v1.HomeContentListRes, err error) {
	res = &v1.HomeContentListRes{
		List:  make([]v1.ContentItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 构建过滤条件
	filter := make(map[string]interface{})

	// 只查询已发布状态的内容
	filter["status"] = "已发布"

	// 根据内容类型过滤
	if req.Type == 1 {
		// 推荐内容
		filter["is_recommended"] = true
	}

	// 查询数据
	list, total, err := s.contentDao.FindList(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		g.Log().Error(ctx, "获取首页内容列表失败", err)
		return nil, gerror.New("获取首页内容列表失败: " + err.Error())
	}

	// 设置响应
	res.Total = total
	res.Page = req.Page

	// 转换数据
	res.List = s.convertToContentItems(list)

	return res, nil
}

// IdleContentList 获取闲置内容列表
func (s *contentClientImpl) IdleContentList(ctx context.Context, req *v1.IdleContentListReq) (res *v1.IdleContentListRes, err error) {
	res = &v1.IdleContentListRes{
		List:  make([]v1.ContentItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 构建过滤条件
	filter := make(map[string]interface{})

	// 只查询已发布状态的内容
	filter["status"] = "已发布"

	// 闲置内容（假设状态为"闲置"的内容被标记在分类中）
	filter["category"] = "闲置"

	// 关键词搜索
	if req.Keyword != "" {
		filter["title"] = req.Keyword
	}

	// 分类过滤
	if req.Category > 0 {
		filter["sub_category"] = req.Category
	}

	// 查询数据，暂时不处理排序，默认按创建时间排序
	list, total, err := s.contentDao.FindList(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		g.Log().Error(ctx, "获取闲置内容列表失败", err)
		return nil, gerror.New("获取闲置内容列表失败: " + err.Error())
	}

	// 设置响应
	res.Total = total
	res.Page = req.Page

	// 转换数据
	res.List = s.convertToContentItems(list)

	return res, nil
}

// CategoryList 获取分类列表
func (s *contentClientImpl) CategoryList(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	res = &v1.CategoryListRes{
		List: make([]v1.CategoryItem, 0),
	}

	// 从数据库获取分类列表
	var categoryItems []v1.CategoryItem

	// 根据分类类型获取对应的分类数据
	if req.Type == 1 {
		// 获取首页分类列表
		homeCategoryDao := dao.NewHomeCategoryDao()
		homeCategories, err := homeCategoryDao.FindList(ctx)
		if err != nil {
			g.Log().Error(ctx, "获取首页分类列表失败", err)
			return nil, gerror.New("获取首页分类列表失败: " + err.Error())
		}

		// 转换首页分类数据
		for _, item := range homeCategories {
			if gconv.Bool(item.IsActive) {
				icon := gconv.String(item.Icon)
				if icon == "" {
					icon = "icon-category.png" // 使用默认图标
				}

				categoryItem := v1.CategoryItem{
					Id:   gconv.Int(item.Id),
					Name: gconv.String(item.Name),
					Icon: icon,
					Type: 1,
				}

				// 获取该分类下内容数量
				filter := map[string]interface{}{
					"category": item.Name,
					"status":   "已发布",
				}
				_, count, _ := s.contentDao.FindList(ctx, filter, 1, 1)
				categoryItem.Count = count

				categoryItems = append(categoryItems, categoryItem)
			}
		}
	} else if req.Type == 2 {
		// 获取闲置分类列表
		idleCategoryDao := dao.NewIdleCategoryDao()
		idleCategories, err := idleCategoryDao.FindList(ctx)
		if err != nil {
			g.Log().Error(ctx, "获取闲置分类列表失败", err)
			return nil, gerror.New("获取闲置分类列表失败: " + err.Error())
		}

		// 转换闲置分类数据
		for _, item := range idleCategories {
			if gconv.Bool(item.IsActive) {
				icon := gconv.String(item.Icon)
				if icon == "" {
					icon = "icon-idle.png" // 使用默认图标
				}

				categoryItem := v1.CategoryItem{
					Id:   gconv.Int(item.Id),
					Name: gconv.String(item.Name),
					Icon: icon,
					Type: 2,
				}

				// 获取该分类下内容数量
				filter := map[string]interface{}{
					"category": item.Name,
					"status":   "已发布",
				}
				_, count, _ := s.contentDao.FindList(ctx, filter, 1, 1)
				categoryItem.Count = count

				categoryItems = append(categoryItems, categoryItem)
			}
		}
	} else {
		// 获取所有分类
		// 先获取首页分类
		homeCategoryDao := dao.NewHomeCategoryDao()
		homeCategories, err := homeCategoryDao.FindList(ctx)
		if err != nil {
			g.Log().Error(ctx, "获取首页分类列表失败", err)
			return nil, gerror.New("获取首页分类列表失败: " + err.Error())
		}

		// 转换首页分类数据
		for _, item := range homeCategories {
			if gconv.Bool(item.IsActive) {
				icon := gconv.String(item.Icon)
				if icon == "" {
					icon = "icon-category.png" // 使用默认图标
				}

				categoryItem := v1.CategoryItem{
					Id:   gconv.Int(item.Id),
					Name: gconv.String(item.Name),
					Icon: icon,
					Type: 1,
				}
				categoryItems = append(categoryItems, categoryItem)
			}
		}

		// 再获取闲置分类
		idleCategoryDao := dao.NewIdleCategoryDao()
		idleCategories, err := idleCategoryDao.FindList(ctx)
		if err != nil {
			g.Log().Error(ctx, "获取闲置分类列表失败", err)
			return nil, gerror.New("获取闲置分类列表失败: " + err.Error())
		}

		// 转换闲置分类数据
		for _, item := range idleCategories {
			if gconv.Bool(item.IsActive) {
				icon := gconv.String(item.Icon)
				if icon == "" {
					icon = "icon-idle.png" // 使用默认图标
				}

				categoryItem := v1.CategoryItem{
					Id:   gconv.Int(item.Id),
					Name: gconv.String(item.Name),
					Icon: icon,
					Type: 2,
				}
				categoryItems = append(categoryItems, categoryItem)
			}
		}
	}

	// 设置返回数据
	res.List = categoryItems

	// 添加性能优化: 缓存分类列表
	g.Log().Debug(ctx, "获取分类列表成功, 共", len(res.List), "项")

	return res, nil
}

// WxIdleCreate 微信客户端-闲置发布
func (s *contentClientImpl) WxIdleCreate(ctx context.Context, req *v1.WxIdleCreateReq) (res *v1.WxIdleCreateRes, err error) {
	res = &v1.WxIdleCreateRes{}

	// 获取当前登录的客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("获取客户信息失败: " + err.Error())
	}

	// 查询客户信息
	type ClientInfo struct {
		Username string `json:"username"`
		RealName string `json:"real_name"`
	}
	var clientInfo *ClientInfo
	err = g.DB().Model("client").
		Fields("username, real_name").
		Where("id", clientId).
		Scan(&clientInfo)
	if err != nil {
		return nil, gerror.New("查询客户信息失败: " + err.Error())
	}

	if clientInfo == nil {
		return nil, gerror.New("客户不存在")
	}

	// 获取作者信息
	authorName := ""
	if clientInfo.RealName != "" {
		authorName = clientInfo.RealName
	} else if clientInfo.Username != "" {
		authorName = clientInfo.Username
	} else {
		authorName = "微信用户"
	}

	// 查询分类信息
	idleCategoryDao := dao.NewIdleCategoryDao()
	idleCategory, err := idleCategoryDao.FindOne(ctx, req.CategoryId)
	if err != nil {
		return nil, gerror.New("查询分类信息失败: " + err.Error())
	}
	if idleCategory == nil {
		return nil, gerror.New("分类不存在")
	}

	// 检查区域是否存在且处于启用状态
	regionExists, err := g.DB().Model("region").
		Where("id", req.RegionId).
		Where("status", 0). // 0表示启用状态
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		return nil, gerror.New("查询地区信息失败: " + err.Error())
	}
	if regionExists == 0 {
		return nil, gerror.New("所选地区不存在或未启用")
	}

	// 处理图片，与WxInfoCreate一致
	var contentStr string
	if len(req.Images) > 0 {
		imagesJson, err := json.Marshal(req.Images)
		if err != nil {
			return nil, gerror.New("处理图片数据失败: " + err.Error())
		}
		contentStr = req.Content + "\n" + string(imagesJson)
	} else {
		contentStr = req.Content
	}

	// 构建内容数据
	contentData := &do.ContentDO{
		Title:         req.Title,
		Category:      gconv.String(idleCategory.Name),
		Author:        authorName,
		RegionId:      req.RegionId,
		ClientId:      clientId,   // 直接存储客户ID
		Content:       contentStr, // 使用处理后的内容，包含图片JSON
		Status:        "已发布",      // 客户端发布默认为已发布状态
		Views:         0,
		Likes:         0,
		Comments:      0,
		IsRecommended: false,
		PublishedAt:   gtime.Now(),
		CreatedAt:     gtime.Now(),
		UpdatedAt:     gtime.Now(),
	}

	// 扩展字段：添加价格、交易地点等信息，但不包含images
	extendData := g.Map{
		"type":          2, // 2表示闲置物品
		"price":         req.Price,
		"originalPrice": req.OriginalPrice,
		"tradePlace":    req.TradePlace,
		"tradeMethod":   req.TradeMethod,
	}

	// 将扩展字段序列化为JSON
	extendJson, err := json.Marshal(extendData)
	if err != nil {
		return nil, gerror.New("处理扩展数据失败: " + err.Error())
	}

	// 设置扩展数据
	contentData.Extend = string(extendJson)

	// 默认设置3天展示期
	contentData.ExpiresAt = gtime.Now().AddDate(0, 0, 3)

	// 验证置顶时长不超过展示时长
	if contentData.TopUntil != nil && contentData.ExpiresAt != nil {
		if contentData.TopUntil.After(contentData.ExpiresAt) {
			g.Log().Warning(ctx, "置顶时长超过展示时长", g.Map{
				"topUntil":  contentData.TopUntil,
				"expiresAt": contentData.ExpiresAt,
			})
			return nil, gerror.New("置顶时长不能超过展示时长，请调整置顶套餐或增加展示时长")
		}
	}

	// 插入数据
	lastInsertId, err := s.contentDao.Insert(ctx, contentData)
	if err != nil {
		return nil, gerror.New("发布闲置信息失败: " + err.Error())
	}

	// 设置返回ID
	res.Id = int(lastInsertId)

	return res, nil
}

// WxInfoCreate 微信客户端-信息发布
func (s *contentClientImpl) WxInfoCreate(ctx context.Context, req *v1.WxInfoCreateReq) (res *v1.WxInfoCreateRes, err error) {
	res = &v1.WxInfoCreateRes{}

	// 记录请求信息，用于调试
	g.Log().Info(ctx, "接收到内容发布请求", req)

	// 获取当前登录的客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取客户信息失败", err)
		return nil, gerror.New("获取客户信息失败: " + err.Error())
	}

	g.Log().Info(ctx, "客户ID验证成功", clientId)

	// 查询客户信息
	type ClientInfo struct {
		Username string `json:"username"`
		RealName string `json:"real_name"`
	}
	var clientInfo *ClientInfo
	err = g.DB().Model("client").
		Fields("username, real_name").
		Where("id", clientId).
		Scan(&clientInfo)
	if err != nil {
		return nil, gerror.New("查询客户信息失败: " + err.Error())
	}

	if clientInfo == nil {
		return nil, gerror.New("客户不存在")
	}

	// 获取作者信息
	authorName := ""
	if clientInfo.RealName != "" {
		authorName = clientInfo.RealName
	} else if clientInfo.Username != "" {
		authorName = clientInfo.Username
	} else {
		authorName = "微信用户"
	}

	// 查询分类信息
	homeCategoryDao := dao.NewHomeCategoryDao()
	homeCategory, err := homeCategoryDao.FindOne(ctx, req.CategoryId)
	if err != nil {
		return nil, gerror.New("查询分类信息失败: " + err.Error())
	}
	if homeCategory == nil {
		return nil, gerror.New("分类不存在")
	}

	// 检查区域是否存在且处于启用状态
	regionExists, err := g.DB().Model("region").
		Where("id", req.RegionId).
		Where("status", 0). // 0表示启用状态
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		return nil, gerror.New("查询地区信息失败: " + err.Error())
	}
	if regionExists == 0 {
		return nil, gerror.New("所选地区不存在或未启用")
	}

	// 处理图片
	var contentStr string
	if len(req.Images) > 0 {
		imagesJson, err := json.Marshal(req.Images)
		if err != nil {
			return nil, gerror.New("处理图片数据失败: " + err.Error())
		}
		contentStr = req.Content + "\n" + string(imagesJson)
	} else {
		contentStr = req.Content
	}

	// 构建内容数据
	contentData := &do.ContentDO{
		Title:         req.Title,
		Category:      gconv.String(homeCategory.Name),
		Author:        authorName,
		RegionId:      req.RegionId, // 设置地区ID
		ClientId:      clientId,     // 直接存储客户ID
		Content:       contentStr,
		Status:        "已发布", // 默认为已发布状态，如需支付则改为待支付
		Views:         0,
		Likes:         0,
		Comments:      0,
		IsRecommended: false,
		PublishedAt:   gtime.Now(),
		CreatedAt:     gtime.Now(),
		UpdatedAt:     gtime.Now(),
	}

	// 扩展字段：添加信息类型
	extendData := g.Map{
		"type": 1, // 1表示普通信息
	}

	// 将扩展字段序列化为JSON
	extendJson, err := json.Marshal(extendData)
	if err != nil {
		return nil, gerror.New("处理扩展数据失败: " + err.Error())
	}

	// 设置扩展数据
	contentData.Extend = string(extendJson)

	// 计算订单金额
	var totalAmount float64 = 0
	var topPackageName, publishPackageName string

	// 检查用户是否选择了套餐
	needCreateOrder := false
	orderRemark := ""

	// 处理置顶申请
	if req.IsTopRequest {
		if req.TopPackageId > 0 {
			// 获取置顶套餐详情
			packageDao := dao.NewPackageDao()
			topPackage, err := packageDao.FindOne(ctx, req.TopPackageId)
			if err != nil {
				g.Log().Error(ctx, "获取置顶套餐信息失败", err, g.Map{"topPackageId": req.TopPackageId})
				return nil, gerror.New("获取置顶套餐信息失败: " + err.Error())
			}
			if topPackage == nil {
				g.Log().Warning(ctx, "置顶套餐不存在", g.Map{"topPackageId": req.TopPackageId})
				return nil, gerror.New("置顶套餐不存在")
			}

			// 验证套餐类型
			packageType := gconv.String(topPackage.Type)
			g.Log().Debug(ctx, "验证置顶套餐类型", g.Map{
				"topPackageId": req.TopPackageId,
				"type":         packageType,
				"expected":     do.PackageTypeTop,
			})

			if packageType != do.PackageTypeTop {
				g.Log().Warning(ctx, "套餐类型错误", g.Map{
					"topPackageId": req.TopPackageId,
					"type":         packageType,
					"expected":     do.PackageTypeTop,
				})
				return nil, gerror.New("套餐类型错误，必须使用置顶套餐")
			}

			// 添加置顶套餐金额
			totalAmount += gconv.Float64(topPackage.Price)
			topPackageName = gconv.String(topPackage.Title)
			orderRemark += "置顶套餐: " + topPackageName + ", "
			needCreateOrder = true

			// 设置置顶截止时间
			durationType := gconv.String(topPackage.DurationType)
			duration := gconv.Int(topPackage.Duration)

			switch durationType {
			case "hour":
				contentData.TopUntil = gtime.Now().Add(time.Hour * time.Duration(duration))
			case "day":
				contentData.TopUntil = gtime.Now().AddDate(0, 0, duration)
			case "month":
				contentData.TopUntil = gtime.Now().AddDate(0, duration, 0)
			default:
				// 默认按天计算
				contentData.TopUntil = gtime.Now().AddDate(0, 0, duration)
			}
			contentData.IsRecommended = true
		} else if req.TopDays > 0 {
			// 根据天数计算置顶时间
			contentData.TopUntil = gtime.Now().AddDate(0, 0, req.TopDays)
			contentData.IsRecommended = true
		}
	}

	// 处理展示套餐
	if req.PublishPackageId > 0 {
		// 获取展示套餐详情
		packageDao := dao.NewPackageDao()
		publishPackage, err := packageDao.FindOne(ctx, req.PublishPackageId)
		if err != nil {
			g.Log().Error(ctx, "获取展示套餐信息失败", err, g.Map{"publishPackageId": req.PublishPackageId})
			return nil, gerror.New("获取展示套餐信息失败: " + err.Error())
		}
		if publishPackage == nil {
			g.Log().Warning(ctx, "展示套餐不存在", g.Map{"publishPackageId": req.PublishPackageId})
			return nil, gerror.New("展示套餐不存在")
		}

		// 验证套餐类型
		packageType := gconv.String(publishPackage.Type)
		g.Log().Debug(ctx, "验证展示套餐类型", g.Map{
			"publishPackageId": req.PublishPackageId,
			"type":             packageType,
			"expected":         do.PackageTypePublish,
		})

		if packageType != do.PackageTypePublish {
			g.Log().Warning(ctx, "套餐类型错误", g.Map{
				"publishPackageId": req.PublishPackageId,
				"type":             packageType,
				"expected":         do.PackageTypePublish,
			})
			return nil, gerror.New("套餐类型错误，必须使用展示套餐")
		}

		// 添加展示套餐金额
		totalAmount += gconv.Float64(publishPackage.Price)
		publishPackageName = gconv.String(publishPackage.Title)
		orderRemark += "展示套餐: " + publishPackageName
		needCreateOrder = true

		// 设置内容到期时间
		durationType := gconv.String(publishPackage.DurationType)
		duration := gconv.Int(publishPackage.Duration)

		switch durationType {
		case "hour":
			contentData.ExpiresAt = gtime.Now().Add(time.Hour * time.Duration(duration))
		case "day":
			contentData.ExpiresAt = gtime.Now().AddDate(0, 0, duration)
		case "month":
			contentData.ExpiresAt = gtime.Now().AddDate(0, duration, 0)
		default:
			// 默认按天计算
			contentData.ExpiresAt = gtime.Now().AddDate(0, 0, duration)
		}
	} else {
		// API验证已确保必须选择展示套餐，这个分支应该不会执行
		return nil, gerror.New("必须选择展示套餐")
	}

	// 验证置顶时长不超过展示时长
	if contentData.TopUntil != nil && contentData.ExpiresAt != nil {
		if contentData.TopUntil.After(contentData.ExpiresAt) {
			g.Log().Warning(ctx, "置顶时长超过展示时长", g.Map{
				"topUntil":  contentData.TopUntil,
				"expiresAt": contentData.ExpiresAt,
			})
			return nil, gerror.New("置顶时长不能超过展示时长，请调整置顶套餐或增加展示时长")
		}
	}

	// 如果需要创建订单，则修改内容状态为待支付
	if needCreateOrder && totalAmount > 0 {
		contentData.Status = "待支付" // 内容状态设为待支付，等待支付完成后再改为已发布
	}

	// 插入数据
	lastInsertId, err := s.contentDao.Insert(ctx, contentData)
	if err != nil {
		return nil, gerror.New("发布信息失败: " + err.Error())
	}

	// 设置返回ID
	res.Id = int(lastInsertId)

	// 创建订单（如果需要）
	if needCreateOrder && totalAmount > 0 {
		contentId := int(lastInsertId)

		// 生成订单号
		orderNo := fmt.Sprintf("ORD%s%d", time.Now().Format("20060102150405"), clientId)

		// 获取订单过期时间配置 - 使用常量代替字符串
		orderExpireMinutes := g.Cfg().MustGet(ctx, consts.ConfigOrderExpireTime, 1).Int()

		// 计算过期时间
		expireTime := gtime.Now().Add(time.Duration(orderExpireMinutes) * time.Minute)

		// 构建订单数据
		orderData := g.Map{
			"order_no":       orderNo,
			"client_id":      clientId,
			"client_name":    authorName,
			"content_id":     contentId,
			"product_name":   "内容发布套餐",
			"amount":         totalAmount,
			"status":         0,                                                                               // 待支付
			"payment_method": g.Cfg().MustGet(ctx, consts.ConfigOrderDefaultPaymentMethod, "wechat").String(), // 从配置读取默认支付方式
			"remark":         orderRemark,
			"created_at":     gtime.Now(),
			"expire_time":    expireTime, // 添加过期时间
		}

		// 插入订单
		_, err = g.DB().Model("order").Data(orderData).Insert()
		if err != nil {
			g.Log().Error(ctx, "创建订单失败:", err)
			// 不需要中断流程，继续返回内容ID
		} else {
			// 将订单号和订单信息关联到内容扩展数据中
			extendData["order_no"] = orderNo
			extendJson, _ = json.Marshal(extendData)

			// 更新内容扩展数据
			g.DB().Model(do.TableContent).Data(g.Map{
				"extend": string(extendJson),
			}).Where("id", contentId).Update()

			// 添加订单号到返回结果
			res.OrderNo = orderNo

			// 添加订单信息到返回结果
			g.Log().Info(ctx, "创建订单成功，订单号:", orderNo)
		}
	} else {
		g.Log().Info(ctx, "无需创建订单或订单金额为0，使用免费套餐")
	}

	return res, nil
}

// WxClientPackageList 微信客户端-获取套餐列表
func (s *contentClientImpl) WxClientPackageList(ctx context.Context, req *v1.WxClientPackageListReq) (res *v1.WxClientPackageListRes, err error) {
	// 记录请求信息
	g.Log().Info(ctx, "接收到获取套餐列表请求")

	res = &v1.WxClientPackageListRes{
		TopPackages:     make([]*v1.Package, 0),
		PublishPackages: make([]*v1.Package, 0),
	}

	// 获取套餐DAO
	packageDao := dao.NewPackageDao()

	// 获取系统配置DAO，用于查询总开关状态
	systemConfigDao := &dao.SystemConfigDao{}

	// 查询置顶套餐总开关状态
	topEnabled, err := systemConfigDao.GetTopPackageEnabled(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取置顶套餐总开关状态失败", err)
		return nil, gerror.New("获取置顶套餐总开关状态失败: " + err.Error())
	}
	res.TopEnabled = topEnabled

	// 查询发布套餐总开关状态
	publishEnabled, err := systemConfigDao.GetPublishPackageEnabled(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取发布套餐总开关状态失败", err)
		return nil, gerror.New("获取发布套餐总开关状态失败: " + err.Error())
	}
	res.PublishEnabled = publishEnabled

	// 获取所有套餐并分类
	g.Log().Debug(ctx, "正在查询所有套餐")
	packages, err := packageDao.FindAll(ctx, req.Sort, req.Order)
	if err != nil {
		g.Log().Error(ctx, "获取套餐列表失败", err)
		return nil, gerror.New("获取套餐列表失败: " + err.Error())
	}
	g.Log().Debug(ctx, "查询到套餐总数量", len(packages))

	// 分类处理
	for _, p := range packages {
		if p != nil {
			pkg := &v1.Package{
				Id:           gconv.Int(p.Id),
				Title:        gconv.String(p.Title),
				Description:  gconv.String(p.Description),
				Price:        gconv.Float64(p.Price),
				Type:         v1.PackageType(gconv.String(p.Type)),
				Duration:     gconv.Int(p.Duration),
				DurationType: v1.DurationType(gconv.String(p.DurationType)),
				SortOrder:    gconv.Int(p.SortOrder),
			}

			// 根据类型分配到不同列表
			packageType := gconv.String(p.Type)
			if packageType == string(v1.PackageTypeTop) {
				res.TopPackages = append(res.TopPackages, pkg)
			} else if packageType == string(v1.PackageTypePublish) {
				res.PublishPackages = append(res.PublishPackages, pkg)
			} else {
				g.Log().Warning(ctx, "发现未知类型的套餐", g.Map{
					"id":   gconv.Int(p.Id),
					"type": packageType,
				})
			}
		}
	}

	g.Log().Info(ctx, "返回套餐列表结果", g.Map{
		"topPackagesCount":     len(res.TopPackages),
		"publishPackagesCount": len(res.PublishPackages),
		"topEnabled":           res.TopEnabled,
		"publishEnabled":       res.PublishEnabled,
	})
	return res, nil
}

// stripHtmlTags 简单实现的HTML标签移除函数
func (s *contentClientImpl) stripHtmlTags(html string) string {
	// 处理空输入
	if html == "" {
		return ""
	}

	// 处理无效的UTF-8序列
	if !utf8.ValidString(html) {
		// 尝试将可能存在的非UTF-8编码转换为UTF-8
		htmlBytes := []byte(html)
		validBytes := make([]byte, 0, len(htmlBytes))

		for i := 0; i < len(htmlBytes); {
			r, size := utf8.DecodeRune(htmlBytes[i:])
			if r == utf8.RuneError && size == 1 {
				// 跳过无效字节
				i++
			} else {
				validBytes = append(validBytes, htmlBytes[i:i+size]...)
				i += size
			}
		}
		html = string(validBytes)
	}

	// 移除HTML标签
	var result strings.Builder
	inTag := false

	for _, r := range html {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}

	// 替换HTML实体
	output := result.String()
	output = strings.ReplaceAll(output, "&nbsp;", " ")
	output = strings.ReplaceAll(output, "&lt;", "<")
	output = strings.ReplaceAll(output, "&gt;", ">")
	output = strings.ReplaceAll(output, "&amp;", "&")
	output = strings.ReplaceAll(output, "&quot;", "\"")
	output = strings.ReplaceAll(output, "&#39;", "'")

	// 移除多余空白
	output = strings.TrimSpace(output)

	return output
}

// getCategoryById 根据ID获取分类
func (s *contentClientImpl) getCategoryById(ctx context.Context, id int) *v1.CategoryItem {
	g.Log().Debug(ctx, "开始查询分类ID:", id)

	// 如果是首页分类，则查询首页分类
	homeCategoryDao := dao.NewHomeCategoryDao()
	homeCategory, err := homeCategoryDao.FindOne(ctx, id)
	if err == nil && homeCategory != nil {
		// 找到首页分类
		icon := gconv.String(homeCategory.Icon)
		if icon == "" {
			icon = "icon-category.png" // 默认图标
		}

		g.Log().Debug(ctx, "查询到首页分类:", gconv.String(homeCategory.Name))
		return &v1.CategoryItem{
			Id:   gconv.Int(homeCategory.Id),
			Name: gconv.String(homeCategory.Name),
			Icon: icon,
			Type: 1, // 首页分类
		}
	}

	// 查询闲置分类
	idleCategoryDao := dao.NewIdleCategoryDao()
	idleCategory, err := idleCategoryDao.FindOne(ctx, id)
	if err == nil && idleCategory != nil {
		// 找到闲置分类
		icon := gconv.String(idleCategory.Icon)
		if icon == "" {
			icon = "icon-idle.png" // 默认图标
		}

		g.Log().Debug(ctx, "查询到闲置分类:", gconv.String(idleCategory.Name))
		return &v1.CategoryItem{
			Id:   gconv.Int(idleCategory.Id),
			Name: gconv.String(idleCategory.Name),
			Icon: icon,
			Type: 2, // 闲置分类
		}
	}

	// 如果都没找到，返回nil
	g.Log().Warning(ctx, "未找到分类ID:", id)
	return nil
}

// convertToContentItems 将DO转换为API响应结构
func (s *contentClientImpl) convertToContentItems(list []*do.ContentDO) []v1.ContentItem {
	items := make([]v1.ContentItem, 0)
	if len(list) == 0 {
		return items
	}

	for _, content := range list {
		// 处理作者信息
		authorName := gconv.String(content.Author)

		// 从ClientId获取发布者信息 (可以优化作者名称显示)
		clientId := gconv.Int(content.ClientId)
		if clientId > 0 {
			var clientInfo struct {
				RealName string `json:"real_name"`
				Avatar   string `json:"avatar_url"`
			}
			_ = g.DB().Model("client").
				Fields("real_name, avatar_url").
				Where("id", clientId).
				Scan(&clientInfo)

			if clientInfo.RealName != "" {
				authorName = clientInfo.RealName // 优先使用真实姓名
			}
		}

		// 处理图片
		var coverImage string
		contentStr := gconv.String(content.Content)
		// 简单判断内容是否包含JSON格式的图片数组
		if len(contentStr) > 0 && strings.Contains(contentStr, "[") && strings.Contains(contentStr, "]") {
			var images []string
			startIdx := strings.Index(contentStr, "[")
			endIdx := strings.LastIndex(contentStr, "]")
			if startIdx >= 0 && endIdx > startIdx {
				jsonStr := contentStr[startIdx : endIdx+1]
				if err := json.Unmarshal([]byte(jsonStr), &images); err == nil && len(images) > 0 {
					coverImage = images[0]
				}
			}
		}

		// 处理内容摘要
		summary := s.stripHtmlTags(contentStr)
		if len(summary) > 100 {
			summary = summary[:100] + "..."
		}

		// 处理扩展字段
		contentType := 1 // 默认为普通信息
		price := 0.0
		if content.Extend != nil {
			extendStr := gconv.String(content.Extend)
			var extendMap map[string]interface{}
			if err := json.Unmarshal([]byte(extendStr), &extendMap); err == nil {
				// 获取内容类型
				if typeVal, ok := extendMap["type"]; ok {
					contentType = gconv.Int(typeVal)
				}
				// 获取价格信息
				if contentType == 2 && extendMap["price"] != nil {
					price = gconv.Float64(extendMap["price"])
				}
			}
		}

		// 构建列表项
		item := v1.ContentItem{
			Id:         gconv.Int(content.Id),
			Title:      gconv.String(content.Title),
			Category:   gconv.String(content.Category),
			Author:     authorName,
			Status:     gconv.String(content.Status),
			Views:      gconv.Int(content.Views),
			Likes:      gconv.Int(content.Likes),
			Comments:   gconv.Int(content.Comments),
			Type:       contentType,
			Price:      price,
			Summary:    summary,
			CoverImage: coverImage,
		}

		// 判断是否推荐和置顶
		// 1. 判断推荐状态：必须同时满足is_recommended=1和置顶时间未过期
		if gconv.Bool(content.IsRecommended) && content.TopUntil != nil && content.TopUntil.After(gtime.Now()) {
			item.IsRecommended = true
		} else {
			item.IsRecommended = false
		}

		// 2. 判断置顶状态：必须同时满足is_recommended=1和置顶时间未过期
		if gconv.Bool(content.IsRecommended) && content.TopUntil != nil && content.TopUntil.After(gtime.Now()) {
			item.IsTop = true
		} else {
			item.IsTop = false
		}

		// 处理时间
		if content.PublishedAt != nil {
			item.PublishedAt = content.PublishedAt
			item.PublishedAtStr = content.PublishedAt.Format("Y-m-d H:i:s")
		} else if content.CreatedAt != nil {
			item.PublishedAt = content.CreatedAt
			item.PublishedAtStr = content.CreatedAt.Format("Y-m-d H:i:s")
		}

		items = append(items, item)
	}
	return items
}

// RegionContentList 按地区获取内容列表
func (s *contentClientImpl) RegionContentList(ctx context.Context, req *v1.RegionContentListReq) (res *v1.RegionContentListRes, err error) {
	res = &v1.RegionContentListRes{
		List:  make([]v1.RegionContentItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 构建过滤条件
	filter := make(map[string]interface{})

	// 只查询已发布状态的内容
	filter["status"] = "已发布"

	// 根据地区ID过滤
	filter["region_id"] = req.RegionId

	// 关键词搜索
	if req.Keyword != "" {
		filter["title"] = req.Keyword
	}

	// 分类过滤
	if req.Category > 0 {
		// 查询分类信息
		categoryItem := s.getCategoryById(ctx, req.Category)
		if categoryItem != nil {
			filter["category"] = categoryItem.Name
		}
	}

	// 设置为普通信息类型(type=1)，不返回闲置物品
	filter["content_type"] = 1

	// 查询数据
	list, total, err := s.contentDao.FindList(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		g.Log().Error(ctx, "按地区获取内容列表失败", err)
		return nil, gerror.New("按地区获取内容列表失败: " + err.Error())
	}

	// 设置响应
	res.Total = total
	res.Page = req.Page

	// 转换数据
	for _, content := range list {
		// 处理图片
		coverImage := ""
		contentStr := gconv.String(content.Content)

		// 简单判断内容是否包含JSON格式的图片数组
		if len(contentStr) > 0 && strings.Contains(contentStr, "[") && strings.Contains(contentStr, "]") {
			var images []string
			startIdx := strings.Index(contentStr, "[")
			endIdx := strings.LastIndex(contentStr, "]")
			if startIdx >= 0 && endIdx > startIdx {
				jsonStr := contentStr[startIdx : endIdx+1]
				if err := json.Unmarshal([]byte(jsonStr), &images); err == nil && len(images) > 0 {
					coverImage = images[0]
				}
			}
		}

		// 从HTML内容中提取第一张图片作为封面
		if coverImage == "" {
			images := extractImagesFromHtml(contentStr)
			if len(images) > 0 {
				coverImage = images[0]
			}
		}

		// 处理发布时间格式
		publishTime := ""
		if content.PublishedAt != nil {
			publishTime = content.PublishedAt.Format("Y-m-d H:i:s")
		} else if content.CreatedAt != nil {
			publishTime = content.CreatedAt.Format("Y-m-d H:i:s")
		}

		// 获取发布者信息
		publisher := gconv.String(content.Author)
		clientId := gconv.Int(content.ClientId)
		if clientId > 0 {
			var clientInfo struct {
				RealName string `json:"real_name"`
			}
			_ = g.DB().Model("client").
				Fields("real_name").
				Where("id", clientId).
				Scan(&clientInfo)

			if clientInfo.RealName != "" {
				publisher = clientInfo.RealName // 优先使用真实姓名
			}
		}

		// 构建列表项
		item := v1.RegionContentItem{
			Id:          gconv.Int(content.Id),
			Title:       gconv.String(content.Title),
			Category:    gconv.String(content.Category),
			Publisher:   publisher,
			PublishTime: publishTime,
			Image:       coverImage,
		}

		// 判断是否置顶 (同时检查IsRecommended和置顶时间是否过期)
		if gconv.Bool(content.IsRecommended) && content.TopUntil != nil && content.TopUntil.After(gtime.Now()) {
			item.IsTop = true
		} else {
			item.IsTop = false // 如果没有满足条件，则不置顶
		}

		res.List = append(res.List, item)
	}

	return res, nil
}

// RegionIdleList 按地区获取闲置物品列表
func (s *contentClientImpl) RegionIdleList(ctx context.Context, req *v1.RegionIdleListReq) (res *v1.RegionIdleListRes, err error) {
	res = &v1.RegionIdleListRes{
		List:  make([]v1.RegionIdleItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 构建过滤条件
	filter := make(map[string]interface{})

	// 只查询已发布状态的内容
	filter["status"] = "已发布"

	// 根据地区ID过滤
	filter["region_id"] = req.RegionId

	// 关键词搜索
	if req.Keyword != "" {
		filter["title"] = req.Keyword
	}

	// 分类过滤
	if req.Category > 0 {
		// 直接查询闲置分类信息，不使用getCategoryById方法（避免获取到首页分类）
		idleCategoryDao := dao.NewIdleCategoryDao()
		idleCategory, err := idleCategoryDao.FindOne(ctx, req.Category)
		if err == nil && idleCategory != nil {
			// 找到闲置分类
			g.Log().Debug(ctx, "按闲置分类ID过滤:", req.Category, ", 分类名称:", gconv.String(idleCategory.Name))
			filter["category"] = gconv.String(idleCategory.Name)
		} else {
			g.Log().Error(ctx, "未找到闲置分类ID:", req.Category)
		}
	}

	// 设置为闲置物品类型(type=2)
	filter["content_type"] = 2

	g.Log().Debug(ctx, "闲置物品查询条件:", filter)

	// 查询数据
	list, total, err := s.contentDao.FindList(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		g.Log().Error(ctx, "按地区获取闲置物品列表失败", err)
		return nil, gerror.New("按地区获取闲置物品列表失败: " + err.Error())
	}

	// 日志记录查询结果
	g.Log().Debug(ctx, "闲置物品查询结果数量:", total)

	// 设置响应
	res.Total = total
	res.Page = req.Page

	// 转换数据
	for _, content := range list {
		// 处理内容摘要
		contentStr := gconv.String(content.Content)

		// 过滤HTML标签，获取纯文本
		plainText := s.stripHtmlTags(contentStr)

		// 去除URL链接
		plainText = regexp.MustCompile(`https?://\S+`).ReplaceAllString(plainText, "")
		plainText = regexp.MustCompile(`www\.\S+`).ReplaceAllString(plainText, "")

		// 去除JSON数组和其他特殊字符
		plainText = regexp.MustCompile(`\[["']`).ReplaceAllString(plainText, "")
		plainText = regexp.MustCompile(`\n\[`).ReplaceAllString(plainText, "")
		plainText = regexp.MustCompile(`\n\["?`).ReplaceAllString(plainText, "")

		// 去除引号和其他特殊符号
		plainText = strings.ReplaceAll(plainText, "\"", "")
		plainText = strings.ReplaceAll(plainText, "[", "")
		plainText = strings.ReplaceAll(plainText, "]", "")
		plainText = strings.ReplaceAll(plainText, "\n", " ")

		// 多个空格替换为单个空格
		plainText = regexp.MustCompile(`\s+`).ReplaceAllString(plainText, " ")

		// 截取前10个字符作为摘要（使用rune处理中文）
		summary := ""
		runes := []rune(plainText)
		if len(runes) > 0 {
			// 先去除前后空白
			plainText = strings.TrimSpace(plainText)
			runes = []rune(plainText)

			// 截取前10个字符（按照中文字符计算）
			if len(runes) > 10 {
				summary = string(runes[:10])
			} else {
				summary = string(runes)
			}
		}

		// 从extend中获取交易地点和价格
		tradePlace := ""
		price := 0.0

		// 处理图片
		coverImage := ""

		if content.Extend != nil {
			extendStr := gconv.String(content.Extend)
			var extendMap map[string]interface{}
			if err := json.Unmarshal([]byte(extendStr), &extendMap); err == nil {
				// 获取交易地点
				if tradePlaceVal, ok := extendMap["tradePlace"]; ok {
					tradePlace = gconv.String(tradePlaceVal)
				}
				// 获取价格
				if priceVal, ok := extendMap["price"]; ok {
					price = gconv.Float64(priceVal)
				}
				// 获取图片
				if imagesVal, ok := extendMap["images"]; ok {
					if images, ok := imagesVal.([]interface{}); ok && len(images) > 0 {
						coverImage = gconv.String(images[0])
					}
				}
			}
		}

		// 如果在扩展字段中没有找到图片，尝试从内容中提取
		if coverImage == "" {
			// 简单判断内容是否包含JSON格式的图片数组
			if len(contentStr) > 0 && strings.Contains(contentStr, "[") && strings.Contains(contentStr, "]") {
				var images []string
				startIdx := strings.Index(contentStr, "[")
				endIdx := strings.LastIndex(contentStr, "]")
				if startIdx >= 0 && endIdx > startIdx {
					jsonStr := contentStr[startIdx : endIdx+1]
					if err := json.Unmarshal([]byte(jsonStr), &images); err == nil && len(images) > 0 {
						coverImage = images[0]
					}
				}
			}

			// 如果仍未找到图片，尝试从HTML内容中提取
			if coverImage == "" {
				images := extractImagesFromHtml(contentStr)
				if len(images) > 0 {
					coverImage = images[0]
				}
			}
		}

		// 构建列表项
		item := v1.RegionIdleItem{
			Id:         gconv.Int(content.Id),
			Title:      gconv.String(content.Title),
			Summary:    summary,
			TradePlace: tradePlace,
			Price:      price,
			Likes:      gconv.Int(content.Likes),
			Image:      coverImage,
		}

		res.List = append(res.List, item)
	}

	return res, nil
}

// 从HTML内容中提取图片URL
func extractImagesFromHtml(html string) []string {
	var images []string
	imgStart := "<img src=\""
	for len(html) > 0 {
		startIdx := strings.Index(html, imgStart)
		if startIdx < 0 {
			break
		}

		startIdx += len(imgStart)
		endIdx := strings.Index(html[startIdx:], "\"")
		if endIdx < 0 {
			break
		}

		imgUrl := html[startIdx : startIdx+endIdx]
		images = append(images, imgUrl)
		html = html[startIdx+endIdx:]
	}
	return images
}

// WxMyPublishList 微信客户端-获取我的发布列表
func (s *contentClientImpl) WxMyPublishList(ctx context.Context, req *v1.WxMyPublishListReq) (res *v1.WxMyPublishListRes, err error) {
	res = &v1.WxMyPublishListRes{}

	// 获取当前登录的客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 构建查询条件
	filter := map[string]interface{}{
		"client_id": clientId,
	}

	// 根据类型过滤
	if req.Type > 0 {
		filter["content_type"] = req.Type
	}

	// 根据状态过滤
	if req.Status != "" {
		filter["status"] = req.Status
	}

	// 查询数据
	contentDao := dao.NewContentDao()
	contents, total, err := contentDao.FindList(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		return nil, gerror.New("获取我的发布列表失败: " + err.Error())
	}

	// 计算总页数
	pages := 0
	if total > 0 {
		pages = int(math.Ceil(float64(total) / float64(req.PageSize)))
	}

	// 构建响应
	res.List = make([]v1.MyPublishItem, 0, len(contents))
	res.Total = total
	res.Page = req.Page
	res.Pages = pages

	// 处理结果
	for _, content := range contents {
		var item v1.MyPublishItem
		contentId := gconv.Int(content.Id)
		item.Id = contentId
		item.Title = gconv.String(content.Title)
		item.Category = gconv.String(content.Category)
		item.Status = gconv.String(content.Status)
		item.PublishedAt = content.PublishedAt

		res.List = append(res.List, item)
	}

	return res, nil
}

// WxMyPublishCount 微信客户端-获取我的发布数量
func (s *contentClientImpl) WxMyPublishCount(ctx context.Context, req *v1.WxMyPublishCountReq) (res *v1.WxMyPublishCountRes, err error) {
	res = &v1.WxMyPublishCountRes{}

	// 获取当前登录的客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 计算总发布数量
	total, err := g.DB().Model("content").
		Where("client_id", clientId).
		Count()
	if err != nil {
		return nil, gerror.New("获取发布总数量失败: " + err.Error())
	}
	res.Total = total

	return res, nil
}
