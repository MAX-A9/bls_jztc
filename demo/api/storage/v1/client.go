package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// WxUploadImageReq 微信客户端-上传图片请求
type WxUploadImageReq struct {
	g.Meta `path:"/wx/upload/image" method:"post" mime:"multipart/form-data" tags:"微信客户端" summary:"上传图片" security:"Bearer" description:"微信客户端上传图片"`
	File   *ghttp.UploadFile `p:"file" v:"required#请选择需要上传的图片" dc:"上传的图片文件"`
}

// WxUploadImageRes 微信客户端-上传图片响应
type WxUploadImageRes struct {
	g.Meta   `mime:"application/json" example:"json"`
	Url      string `json:"url" dc:"图片URL"`
	Filename string `json:"filename" dc:"文件名"`
	Size     int64  `json:"size" dc:"文件大小(字节)"`
}
