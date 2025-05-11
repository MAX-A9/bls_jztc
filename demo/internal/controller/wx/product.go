package wx

import (
	"context"
	v1 "demo/api/mall/v1"
	"demo/internal/service"
)

// 微信小程序商品控制器
type ProductController struct{}

// 创建ProductController实例
func New() *ProductController {
	return &ProductController{}
}

// 获取商品列表
func (c *ProductController) List(ctx context.Context, req *v1.WxProductListReq) (res *v1.WxProductListRes, err error) {
	return service.Product().GetWxProductList(ctx, req)
}

// 获取商品详情
func (c *ProductController) Detail(ctx context.Context, req *v1.WxProductDetailReq) (res *v1.WxProductDetailRes, err error) {
	return service.Product().GetWxProductDetail(ctx, req)
}
