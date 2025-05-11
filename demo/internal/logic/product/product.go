package product

import (
	"context"
	v1 "demo/api/mall/v1"
	"demo/internal/dao"
	"demo/internal/model"
	"demo/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// 商品服务实现
type sProduct struct{}

// 初始化服务并将其注册到服务容器中
func init() {
	service.RegisterProduct(New())
}

// 创建新的商品服务实例
func New() service.IProductService {
	return &sProduct{}
}

// 获取商品列表
func (s *sProduct) GetProductList(ctx context.Context, req *v1.ProductListReq) (res *v1.ProductListRes, err error) {
	// 初始化响应
	res = &v1.ProductListRes{
		List:  make([]*model.Product, 0),
		Total: 0,
		Page:  req.Page,
		Size:  req.Size,
	}

	// 构建查询过滤条件
	filter := &model.ProductFilter{
		Name:       req.Name,
		CategoryId: req.CategoryId,
		Status:     -1, // 默认为-1，表示不按状态筛选，获取所有状态的商品
		SortField:  req.SortField,
		SortOrder:  req.SortOrder,
		Tags:       req.Tags,
	}

	// 如果前端传入了有效的状态值，则使用前端传入的值
	if req.Status >= 0 && req.Status <= 2 {
		filter.Status = req.Status
	}

	// 若有Duration参数，设置为最小时长
	if req.Duration > 0 {
		filter.MinDuration = req.Duration
	}

	// 若有Stock参数，设置为最小库存
	if req.Stock > 0 {
		filter.MinStock = req.Stock
	}

	// 调用DAO层获取数据
	productDao := &dao.ProductDao{}
	productList, total, err := productDao.GetList(ctx, filter, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// 设置响应数据
	res.List = productList
	res.Total = total
	return res, nil
}

// 获取商品详情
func (s *sProduct) GetProductDetail(ctx context.Context, req *v1.ProductDetailReq) (res *v1.ProductDetailRes, err error) {
	res = &v1.ProductDetailRes{}

	// 获取商品信息
	productDao := &dao.ProductDao{}
	product, err := productDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, gerror.New("商品不存在")
	}

	res.Product = product
	return res, nil
}

// 创建商品
func (s *sProduct) CreateProduct(ctx context.Context, req *v1.ProductCreateReq) (res *v1.ProductCreateRes, err error) {
	res = &v1.ProductCreateRes{}

	// 构建商品数据
	product := &model.Product{}
	if err = gconv.Struct(req, product); err != nil {
		return nil, err
	}

	// 调用DAO层创建商品
	productDao := &dao.ProductDao{}
	lastInsertId, err := productDao.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	res.Id = int(lastInsertId)
	return res, nil
}

// 更新商品
func (s *sProduct) UpdateProduct(ctx context.Context, req *v1.ProductUpdateReq) (res *v1.ProductUpdateRes, err error) {
	res = &v1.ProductUpdateRes{}

	// 检查商品是否存在
	productDao := &dao.ProductDao{}
	exists, err := productDao.Exists(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, gerror.New("商品不存在")
	}

	// 构建更新数据
	data := g.Map{}
	if req.Name != "" {
		data["name"] = req.Name
	}
	if req.CategoryId > 0 {
		data["category_id"] = req.CategoryId
	}
	if req.CategoryName != "" {
		data["category_name"] = req.CategoryName
	}
	if req.Price > 0 {
		data["price"] = req.Price
	}
	if req.Duration > 0 {
		data["duration"] = req.Duration
	}
	if req.Stock >= 0 { // 允许将库存设为0
		data["stock"] = req.Stock
		// 如果库存为0且状态为已上架，自动修改状态为已售罄
		if req.Stock == 0 && req.Status == model.ProductStatusListed {
			data["status"] = model.ProductStatusSoldOut
		}
	}
	if req.Status >= 0 {
		data["status"] = req.Status
	}
	if req.SortOrder >= 0 {
		data["sort_order"] = req.SortOrder
	}
	if req.Description != "" {
		data["description"] = req.Description
	}
	if req.Thumbnail != "" {
		data["thumbnail"] = req.Thumbnail
	}
	if req.Images != "" {
		data["images"] = req.Images
	}
	// 更新Tags字段
	data["tags"] = req.Tags

	// 若没有更新数据，直接返回
	if len(data) == 0 {
		return res, nil
	}

	// 调用DAO层更新商品
	if err = productDao.Update(ctx, req.Id, data); err != nil {
		return nil, err
	}

	return res, nil
}

// 删除商品
func (s *sProduct) DeleteProduct(ctx context.Context, req *v1.ProductDeleteReq) (res *v1.ProductDeleteRes, err error) {
	res = &v1.ProductDeleteRes{}

	// 检查商品是否存在
	productDao := &dao.ProductDao{}
	exists, err := productDao.Exists(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, gerror.New("商品不存在")
	}

	// 调用DAO层删除商品
	if err = productDao.Delete(ctx, req.Id); err != nil {
		return nil, err
	}

	return res, nil
}

// 更新商品状态
func (s *sProduct) UpdateProductStatus(ctx context.Context, req *v1.ProductStatusUpdateReq) (res *v1.ProductStatusUpdateRes, err error) {
	res = &v1.ProductStatusUpdateRes{}

	// 检查商品是否存在
	productDao := &dao.ProductDao{}
	product, err := productDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, gerror.New("商品不存在")
	}

	// 检查状态值
	if req.Status < 0 || req.Status > 2 {
		return nil, gerror.New("无效的状态值")
	}

	// 检查库存为0时是否允许上架
	if req.Status == model.ProductStatusListed && product.Stock == 0 {
		return nil, gerror.New("库存为0，无法上架")
	}

	// 调用DAO层更新商品状态
	if err = productDao.UpdateStatus(ctx, req.Id, req.Status); err != nil {
		return nil, err
	}

	return res, nil
}

// 微信小程序端 - 获取商品列表
func (s *sProduct) GetWxProductList(ctx context.Context, req *v1.WxProductListReq) (res *v1.WxProductListRes, err error) {
	// 初始化响应
	res = &v1.WxProductListRes{
		List:  make([]*model.Product, 0),
		Total: 0,
		Page:  req.Page,
		Size:  req.Size,
	}

	// 构建查询过滤条件
	filter := &model.ProductFilter{
		Name:       req.Name,
		CategoryId: req.CategoryId,
		SortField:  req.SortField,
		SortOrder:  req.SortOrder,
		Tags:       req.Tags,
	}

	// 若有Duration参数，设置为最小时长
	if req.Duration > 0 {
		filter.MinDuration = req.Duration
	}

	// 调用DAO层获取数据（微信端只查询已上架且有库存的商品）
	productDao := &dao.ProductDao{}
	productList, total, err := productDao.GetWxList(ctx, filter, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// 设置响应数据
	res.List = productList
	res.Total = total
	return res, nil
}

// 微信小程序端 - 获取商品详情
func (s *sProduct) GetWxProductDetail(ctx context.Context, req *v1.WxProductDetailReq) (res *v1.WxProductDetailRes, err error) {
	res = &v1.WxProductDetailRes{}

	// 获取商品信息
	productDao := &dao.ProductDao{}
	product, err := productDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, gerror.New("商品不存在")
	}

	// 微信端只允许查看已上架的商品
	if product.Status != model.ProductStatusListed {
		return nil, gerror.New("该商品未上架或已售罄")
	}

	res.Product = product
	return res, nil
}
