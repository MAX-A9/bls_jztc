package storage

import (
	"context"
	"fmt"
	"path"
	"strings"

	v1 "demo/api/storage/v1"
	"demo/internal/dao"
	"demo/internal/model"
	"demo/utility/auth"
	"demo/utility/storage"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// WxUploadImage 微信客户端-上传图片
func (s *storageImpl) WxUploadImage(ctx context.Context, req *v1.WxUploadImageReq) (res *v1.WxUploadImageRes, err error) {
	res = &v1.WxUploadImageRes{}

	// 获取当前登录的客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("获取客户信息失败: " + err.Error())
	}

	// 获取OSS客户端
	ossClient, err := storage.GetOSSClient(ctx)
	if err != nil {
		return nil, err
	}

	// 获取上传文件
	uploadFile := req.File
	if uploadFile == nil {
		return nil, gerror.New("请选择需要上传的图片")
	}

	// 验证文件类型
	contentType := uploadFile.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, gerror.New("只能上传图片文件")
	}

	// 验证文件大小
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if uploadFile.Size > maxSize {
		return nil, gerror.New("图片大小不能超过5MB")
	}

	// 获取原始文件名和扩展名
	originalName := uploadFile.Filename
	extName := strings.ToLower(path.Ext(originalName))

	// 允许的图片扩展名
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	extAllowed := false
	for _, ext := range allowedExts {
		if extName == ext {
			extAllowed = true
			break
		}
	}

	if !extAllowed {
		return nil, gerror.New("只支持JPG、PNG、GIF、BMP和WEBP格式的图片")
	}

	// 使用通用函数生成存储路径
	storagePath := ossClient.GenerateStoragePath("wx_image", extName)

	// 处理是否公开访问 - 客户端上传的图片默认公开访问
	isPublic := true

	// 读取上传文件
	tempFile, err := uploadFile.Open()
	if err != nil {
		return nil, gerror.New("打开上传文件失败: " + err.Error())
	}
	defer tempFile.Close()

	// 上传到OSS
	err = ossClient.Bucket.PutObject(ossClient.GetFullObjectKey(storagePath), tempFile)
	if err != nil {
		return nil, gerror.New("上传到OSS失败: " + err.Error())
	}

	// 获取文件URL
	fileUrl, err := ossClient.GetObjectURL(ossClient.GetFullObjectKey(storagePath), !isPublic, 0)
	if err != nil {
		return nil, gerror.New("获取文件URL失败: " + err.Error())
	}

	// 保存文件信息到数据库
	fileData := &model.FileDO{
		Name:        originalName,
		Path:        storagePath,
		Size:        uploadFile.Size,
		Type:        "image",
		ContentType: contentType,
		Extension:   extName[1:], // 去掉扩展名前面的点
		IsPublic:    isPublic,
		UserId:      0,                                 // 客户端上传，不关联管理员用户
		Username:    fmt.Sprintf("客户端用户_%d", clientId), // 记录客户ID
		CreatedAt:   gtime.Now(),
		UpdatedAt:   gtime.Now(),
	}

	// 插入数据库
	fileDao := new(dao.FileDao)
	_, err = fileDao.Insert(ctx, fileData)
	if err != nil {
		g.Log().Warning(ctx, "保存客户端上传图片记录失败", err)
		// 不阻止返回URL，即使数据库记录失败也返回成功
	}

	// 设置响应
	res.Url = fileUrl
	res.Filename = originalName
	res.Size = uploadFile.Size

	// 记录日志
	g.Log().Debug(ctx, "客户端上传图片成功", "clientId", clientId, "url", fileUrl)

	return res, nil
}
