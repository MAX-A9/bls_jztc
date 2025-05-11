package backend

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	v1 "demo/api/backend/v1"
	"demo/internal/dao"
	"demo/internal/model"
	"demo/internal/service"
	"demo/utility/auth"
	"demo/utility/storage"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

// 后端服务实现
type backendImpl struct {
	fileDao *dao.FileDao
}

// New 创建后端服务实例
func New() service.BackendService {
	return &backendImpl{
		fileDao: new(dao.FileDao),
	}
}

// FileUpload 文件上传
func (s *backendImpl) FileUpload(ctx context.Context, req *v1.FileUploadReq) (res *v1.FileUploadRes, err error) {
	res = &v1.FileUploadRes{}

	// 获取当前登录用户信息
	userId, username, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if role != "admin" {
		return nil, gerror.New("您没有权限上传文件")
	}

	// 获取OSS客户端
	ossClient, err := storage.GetOSSClient(ctx)
	if err != nil {
		return nil, err
	}

	// 获取上传文件
	uploadFile := req.File
	if uploadFile == nil {
		return nil, gerror.New("请选择需要上传的文件")
	}

	// 获取原始文件名和扩展名
	originalName := uploadFile.Filename
	extName := strings.ToLower(path.Ext(originalName))
	// 处理文件类型
	fileType := s.getFileTypeByExt(extName)

	// 构建文件存储路径
	storagePath := fmt.Sprintf("uploads/%s/%s/%s%s",
		fileType,
		time.Now().Format("20060102"),
		grand.S(16), // 16位随机字符串
		extName,
	)

	// 如果有自定义目录，添加到路径前缀
	if req.Directory != "" {
		// 确保目录以/结尾
		if !strings.HasSuffix(req.Directory, "/") {
			req.Directory += "/"
		}
		// 确保目录不以/开头
		if strings.HasPrefix(req.Directory, "/") {
			req.Directory = req.Directory[1:]
		}
		storagePath = req.Directory + storagePath
	}

	// 处理是否公开访问
	isPublic := req.IsPublic
	if !isPublic && ossClient.Config.PublicAccess {
		isPublic = true // 如果配置中设置了公开访问，则所有文件都公开
	}

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
		Type:        fileType,
		ContentType: uploadFile.Header.Get("Content-Type"),
		Extension:   extName[1:], // 去掉扩展名前面的点
		IsPublic:    isPublic,
		UserId:      userId,
		Username:    username,
		CreatedAt:   gtime.Now(),
		UpdatedAt:   gtime.Now(),
	}

	// 插入数据库
	fileId, err := s.fileDao.Insert(ctx, fileData)
	if err != nil {
		return nil, gerror.New("保存文件信息失败: " + err.Error())
	}

	// 设置响应
	res.Id = int(fileId)
	res.Name = originalName
	res.Url = fileUrl
	res.Size = uploadFile.Size
	res.Type = fileType
	res.Extension = extName[1:] // 去掉扩展名前面的点
	res.IsPublic = isPublic
	res.UploadTime = gtime.Now().String()

	// 记录操作日志
	service.Log().Record(ctx, userId, username, "文件管理", "上传文件", 1, "")

	return res, nil
}

// FileList 获取文件列表
func (s *backendImpl) FileList(ctx context.Context, req *v1.FileListReq) (res *v1.FileListRes, err error) {
	res = &v1.FileListRes{
		List:  make([]v1.FileInfo, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 获取当前登录用户信息
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if role != "admin" {
		return nil, gerror.New("您没有权限查看文件列表")
	}

	// 获取OSS客户端
	ossClient, err := storage.GetOSSClient(ctx)
	if err != nil {
		return nil, err
	}

	// 构建查询条件
	filter := make(map[string]interface{})
	if req.Keyword != "" {
		filter["keyword"] = req.Keyword
	}
	if req.Type != "" {
		filter["type"] = req.Type
	}
	if req.IsPublic != nil {
		filter["isPublic"] = *req.IsPublic
	}

	// 查询列表
	list, total, err := s.fileDao.FindList(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		return nil, gerror.New("查询文件列表失败: " + err.Error())
	}

	// 设置响应
	res.Total = total
	res.Page = req.Page

	// 转换数据
	for _, item := range list {
		// 获取文件URL
		fileUrl, _ := ossClient.GetObjectURL(ossClient.GetFullObjectKey(gconv.String(item.Path)), !gconv.Bool(item.IsPublic), 0)

		// 格式化文件大小
		sizeFormat := s.formatFileSize(gconv.Int64(item.Size))

		// 添加到结果列表
		res.List = append(res.List, v1.FileInfo{
			Id:          gconv.Int(item.Id),
			Name:        gconv.String(item.Name),
			Url:         fileUrl,
			Path:        gconv.String(item.Path),
			Size:        gconv.Int64(item.Size),
			SizeFormat:  sizeFormat,
			Type:        gconv.String(item.Type),
			ContentType: gconv.String(item.ContentType),
			Extension:   gconv.String(item.Extension),
			IsPublic:    gconv.Bool(item.IsPublic),
			UserId:      gconv.Int(item.UserId),
			Username:    gconv.String(item.Username),
			CreatedAt:   item.CreatedAt.String(),
		})
	}

	return res, nil
}

// FileDetail 获取文件详情
func (s *backendImpl) FileDetail(ctx context.Context, req *v1.FileDetailReq) (res *v1.FileDetailRes, err error) {
	res = &v1.FileDetailRes{}

	// 获取当前登录用户信息
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if role != "admin" {
		return nil, gerror.New("您没有权限查看文件详情")
	}

	// 获取OSS客户端
	ossClient, err := storage.GetOSSClient(ctx)
	if err != nil {
		return nil, err
	}

	// 查询文件信息
	fileInfo, err := s.fileDao.FindById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询文件信息失败: " + err.Error())
	}
	if fileInfo == nil {
		return nil, gerror.New("文件不存在")
	}

	// 获取文件URL
	fileUrl, _ := ossClient.GetObjectURL(ossClient.GetFullObjectKey(gconv.String(fileInfo.Path)), !gconv.Bool(fileInfo.IsPublic), 0)

	// 格式化文件大小
	sizeFormat := s.formatFileSize(gconv.Int64(fileInfo.Size))

	// 设置响应
	res.Id = gconv.Int(fileInfo.Id)
	res.Name = gconv.String(fileInfo.Name)
	res.Url = fileUrl
	res.Path = gconv.String(fileInfo.Path)
	res.Size = gconv.Int64(fileInfo.Size)
	res.SizeFormat = sizeFormat
	res.Type = gconv.String(fileInfo.Type)
	res.ContentType = gconv.String(fileInfo.ContentType)
	res.Extension = gconv.String(fileInfo.Extension)
	res.IsPublic = gconv.Bool(fileInfo.IsPublic)
	res.UserId = gconv.Int(fileInfo.UserId)
	res.Username = gconv.String(fileInfo.Username)
	res.CreatedAt = fileInfo.CreatedAt.String()

	return res, nil
}

// FileDelete 删除文件
func (s *backendImpl) FileDelete(ctx context.Context, req *v1.FileDeleteReq) (res *v1.FileDeleteRes, err error) {
	res = &v1.FileDeleteRes{}

	// 获取当前登录用户信息
	userId, username, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if role != "admin" {
		return nil, gerror.New("您没有权限删除文件")
	}

	// 获取OSS客户端
	ossClient, err := storage.GetOSSClient(ctx)
	if err != nil {
		return nil, err
	}

	// 查询文件信息
	fileInfo, err := s.fileDao.FindById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询文件信息失败: " + err.Error())
	}
	if fileInfo == nil {
		return nil, gerror.New("文件不存在")
	}

	// 从OSS删除文件
	err = ossClient.Bucket.DeleteObject(ossClient.GetFullObjectKey(gconv.String(fileInfo.Path)))
	if err != nil {
		return nil, gerror.New("从OSS删除文件失败: " + err.Error())
	}

	// 从数据库删除记录
	_, err = s.fileDao.Delete(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("从数据库删除文件记录失败: " + err.Error())
	}

	// 设置响应
	res.Success = true
	res.Message = "删除成功"

	// 记录操作日志
	service.Log().Record(ctx, userId, username, "文件管理", "删除文件", 1, "")

	return res, nil
}

// FileBatchDelete 批量删除文件
func (s *backendImpl) FileBatchDelete(ctx context.Context, req *v1.FileBatchDeleteReq) (res *v1.FileBatchDeleteRes, err error) {
	res = &v1.FileBatchDeleteRes{}

	// 获取当前登录用户信息
	userId, username, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if role != "admin" {
		return nil, gerror.New("您没有权限批量删除文件")
	}

	// 获取OSS客户端
	ossClient, err := storage.GetOSSClient(ctx)
	if err != nil {
		return nil, err
	}

	// 遍历ID列表，逐个查询文件信息并删除
	deleteCount := 0
	for _, id := range req.Ids {
		// 查询文件信息
		fileInfo, err := s.fileDao.FindById(ctx, id)
		if err != nil {
			g.Log().Warning(ctx, "查询文件信息失败", err, "fileId", id)
			continue
		}
		if fileInfo == nil {
			g.Log().Warning(ctx, "文件不存在", "fileId", id)
			continue
		}

		// 从OSS删除文件
		err = ossClient.Bucket.DeleteObject(ossClient.GetFullObjectKey(gconv.String(fileInfo.Path)))
		if err != nil {
			g.Log().Warning(ctx, "从OSS删除文件失败", err, "fileId", id, "path", fileInfo.Path)
			continue
		}

		deleteCount++
	}

	// 从数据库批量删除记录
	if len(req.Ids) > 0 {
		_, err = s.fileDao.BatchDelete(ctx, req.Ids)
		if err != nil {
			g.Log().Warning(ctx, "从数据库批量删除文件记录失败", err)
		}
	}

	// 设置响应
	res.Success = true
	res.Message = "批量删除成功"
	res.Count = deleteCount

	// 记录操作日志
	service.Log().Record(ctx, userId, username, "文件管理", "批量删除文件", 1, fmt.Sprintf("删除%d个文件", deleteCount))

	return res, nil
}

// FileUpdatePublic 更新文件公开状态
func (s *backendImpl) FileUpdatePublic(ctx context.Context, req *v1.FileUpdatePublicReq) (res *v1.FileUpdatePublicRes, err error) {
	res = &v1.FileUpdatePublicRes{}

	// 获取当前登录用户信息
	userId, username, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	// 权限检查
	if role != "admin" {
		return nil, gerror.New("您没有权限修改文件状态")
	}

	// 获取OSS客户端
	ossClient, err := storage.GetOSSClient(ctx)
	if err != nil {
		return nil, err
	}

	// 查询文件信息
	fileInfo, err := s.fileDao.FindById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询文件信息失败: " + err.Error())
	}
	if fileInfo == nil {
		return nil, gerror.New("文件不存在")
	}

	// 更新数据库中的文件状态
	_, err = s.fileDao.UpdatePublicStatus(ctx, req.Id, req.IsPublic)
	if err != nil {
		return nil, gerror.New("更新文件状态失败: " + err.Error())
	}

	// 获取新的文件URL
	fileUrl, _ := ossClient.GetObjectURL(ossClient.GetFullObjectKey(gconv.String(fileInfo.Path)), !req.IsPublic, 0)

	// 设置响应
	res.Success = true
	res.Message = "更新成功"
	res.Url = fileUrl

	// 记录操作日志
	statusText := "公开"
	if !req.IsPublic {
		statusText = "私有"
	}
	service.Log().Record(ctx, userId, username, "文件管理", "修改文件状态", 1, fmt.Sprintf("将文件状态修改为%s", statusText))

	return res, nil
}

// formatFileSize 格式化文件大小
func (s *backendImpl) formatFileSize(sizeBytes int64) string {
	if sizeBytes < 1024 {
		return fmt.Sprintf("%d B", sizeBytes)
	} else if sizeBytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(sizeBytes)/1024)
	} else if sizeBytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(sizeBytes)/1024/1024)
	} else {
		return fmt.Sprintf("%.2f GB", float64(sizeBytes)/1024/1024/1024)
	}
}

// getFileTypeByExt 根据扩展名获取文件类型
func (s *backendImpl) getFileTypeByExt(ext string) string {
	ext = strings.ToLower(ext)

	// 图片类型
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}
	if gstr.InArray(imageExts, ext) {
		return "image"
	}

	// 文档类型
	docExts := []string{".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".csv", ".rtf"}
	if gstr.InArray(docExts, ext) {
		return "document"
	}

	// 视频类型
	videoExts := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv"}
	if gstr.InArray(videoExts, ext) {
		return "video"
	}

	// 音频类型
	audioExts := []string{".mp3", ".wav", ".ogg", ".m4a", ".flac"}
	if gstr.InArray(audioExts, ext) {
		return "audio"
	}

	// 压缩文件类型
	archiveExts := []string{".zip", ".rar", ".7z", ".tar", ".gz"}
	if gstr.InArray(archiveExts, ext) {
		return "archive"
	}

	// 默认类型
	return "other"
}
