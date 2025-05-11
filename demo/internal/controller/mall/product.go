package mall

import (
	"context"
	v1 "demo/api/mall/v1"
	"demo/internal/service"
)

// 商品控制器
type ProductController struct{}

// 创建ProductController实例
func New() *ProductController {
	return &ProductController{}
}

// 获取商品列表
func (c *ProductController) List(ctx context.Context, req *v1.ProductListReq) (res *v1.ProductListRes, err error) {
	return service.Product().GetProductList(ctx, req)
}

// 获取商品详情
func (c *ProductController) Detail(ctx context.Context, req *v1.ProductDetailReq) (res *v1.ProductDetailRes, err error) {
	return service.Product().GetProductDetail(ctx, req)
}

// 创建商品
func (c *ProductController) Create(ctx context.Context, req *v1.ProductCreateReq) (res *v1.ProductCreateRes, err error) {
	return service.Product().CreateProduct(ctx, req)
}

// 更新商品
func (c *ProductController) Update(ctx context.Context, req *v1.ProductUpdateReq) (res *v1.ProductUpdateRes, err error) {
	return service.Product().UpdateProduct(ctx, req)
}

// 删除商品
func (c *ProductController) Delete(ctx context.Context, req *v1.ProductDeleteReq) (res *v1.ProductDeleteRes, err error) {
	return service.Product().DeleteProduct(ctx, req)
}

// 更新商品状态
func (c *ProductController) UpdateStatus(ctx context.Context, req *v1.ProductStatusUpdateReq) (res *v1.ProductStatusUpdateRes, err error) {
	return service.Product().UpdateProductStatus(ctx, req)
}
