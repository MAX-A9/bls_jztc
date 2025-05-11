package service

import (
	"context"
	v1 "demo/api/mall/v1"
)

// 商品服务接口
type IProductService interface {
	// 获取商品列表
	GetProductList(ctx context.Context, req *v1.ProductListReq) (res *v1.ProductListRes, err error)

	// 获取商品详情
	GetProductDetail(ctx context.Context, req *v1.ProductDetailReq) (res *v1.ProductDetailRes, err error)

	// 创建商品
	CreateProduct(ctx context.Context, req *v1.ProductCreateReq) (res *v1.ProductCreateRes, err error)

	// 更新商品
	UpdateProduct(ctx context.Context, req *v1.ProductUpdateReq) (res *v1.ProductUpdateRes, err error)

	// 删除商品
	DeleteProduct(ctx context.Context, req *v1.ProductDeleteReq) (res *v1.ProductDeleteRes, err error)

	// 更新商品状态
	UpdateProductStatus(ctx context.Context, req *v1.ProductStatusUpdateReq) (res *v1.ProductStatusUpdateRes, err error)

	// 微信小程序端接口 - 获取商品列表
	GetWxProductList(ctx context.Context, req *v1.WxProductListReq) (res *v1.WxProductListRes, err error)

	// 微信小程序端接口 - 获取商品详情
	GetWxProductDetail(ctx context.Context, req *v1.WxProductDetailReq) (res *v1.WxProductDetailRes, err error)
}
