package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// FileUploadReq 文件上传请求
type FileUploadReq struct {
	g.Meta    `path:"/file/upload" method:"post" mime:"multipart/form-data" tags:"文件管理" summary:"上传文件" security:"Bearer" description:"上传文件到阿里云OSS，需要管理员权限"`
	File      *ghttp.UploadFile `type:"file" v:"required#请选择需要上传的文件" dc:"上传的文件"`
	Directory string            `dc:"存储目录，选填，默认为配置中的目录"`
	IsPublic  bool              `dc:"是否公开访问，选填，默认为配置中的设置"`
}

// FileUploadRes 文件上传响应
type FileUploadRes struct {
	g.Meta     `mime:"application/json" example:"json"`
	Id         int    `json:"id" dc:"文件ID"`
	Name       string `json:"name" dc:"文件名"`
	Url        string `json:"url" dc:"文件访问链接"`
	Size       int64  `json:"size" dc:"文件大小(字节)"`
	Type       string `json:"type" dc:"文件类型"`
	Extension  string `json:"extension" dc:"文件扩展名"`
	IsPublic   bool   `json:"isPublic" dc:"是否公开访问"`
	UploadTime string `json:"uploadTime" dc:"上传时间"`
}

// FileListReq 文件列表请求
type FileListReq struct {
	g.Meta   `path:"/file/list" method:"get" tags:"文件管理" summary:"获取文件列表" security:"Bearer" description:"获取文件列表，支持分页和搜索，需要管理员权限"`
	Page     int    `d:"1" v:"min:1#页码最小为1" dc:"页码，默认1"`
	PageSize int    `d:"10" v:"min:1#每页大小最小为1" dc:"每页数量，默认10"`
	Keyword  string `dc:"搜索关键词，支持文件名搜索"`
	Type     string `dc:"文件类型筛选，例如：image、document等"`
	IsPublic *bool  `dc:"是否公开访问筛选，可选值：true/false"`
}

// FileListRes 文件列表响应
type FileListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []FileInfo `json:"list" dc:"文件列表"`
	Total  int        `json:"total" dc:"总数量"`
	Page   int        `json:"page" dc:"当前页码"`
}

// FileInfo 文件信息结构
type FileInfo struct {
	Id          int    `json:"id" dc:"文件ID"`
	Name        string `json:"name" dc:"文件名"`
	Url         string `json:"url" dc:"文件访问链接"`
	Path        string `json:"path" dc:"文件存储路径"`
	Size        int64  `json:"size" dc:"文件大小(字节)"`
	SizeFormat  string `json:"sizeFormat" dc:"文件大小格式化"`
	Type        string `json:"type" dc:"文件类型"`
	ContentType string `json:"contentType" dc:"内容类型"`
	Extension   string `json:"extension" dc:"文件扩展名"`
	IsPublic    bool   `json:"isPublic" dc:"是否公开访问"`
	UserId      int    `json:"userId" dc:"上传用户ID"`
	Username    string `json:"username" dc:"上传用户名"`
	CreatedAt   string `json:"createdAt" dc:"创建时间"`
}

// FileDeleteReq 删除文件请求
type FileDeleteReq struct {
	g.Meta `path:"/file/delete" method:"delete" tags:"文件管理" summary:"删除文件" security:"Bearer" description:"从阿里云OSS删除文件，需要管理员权限"`
	Id     int `v:"required#文件ID不能为空" dc:"文件ID"`
}

// FileDeleteRes 删除文件响应
type FileDeleteRes struct {
	g.Meta  `mime:"application/json" example:"json"`
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"提示信息"`
}

// FileBatchDeleteReq 批量删除文件请求
type FileBatchDeleteReq struct {
	g.Meta `path:"/file/batch-delete" method:"delete" tags:"文件管理" summary:"批量删除文件" security:"Bearer" description:"从阿里云OSS批量删除文件，需要管理员权限"`
	Ids    []int `v:"required#文件ID列表不能为空" dc:"文件ID列表"`
}

// FileBatchDeleteRes 批量删除文件响应
type FileBatchDeleteRes struct {
	g.Meta  `mime:"application/json" example:"json"`
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"提示信息"`
	Count   int    `json:"count" dc:"成功删除数量"`
}

// FileDetailReq 文件详情请求
type FileDetailReq struct {
	g.Meta `path:"/file/detail" method:"get" tags:"文件管理" summary:"获取文件详情" security:"Bearer" description:"获取文件详情信息，需要管理员权限"`
	Id     int `v:"required#文件ID不能为空" dc:"文件ID"`
}

// FileDetailRes 文件详情响应
type FileDetailRes struct {
	g.Meta `mime:"application/json" example:"json"`
	FileInfo
}

// FileUpdatePublicReq 更新文件公开状态请求
type FileUpdatePublicReq struct {
	g.Meta   `path:"/file/update-public" method:"put" tags:"文件管理" summary:"更新文件公开状态" security:"Bearer" description:"更新文件的公开访问状态，需要管理员权限"`
	Id       int  `v:"required#文件ID不能为空" dc:"文件ID"`
	IsPublic bool `dc:"是否公开访问"`
}

// FileUpdatePublicRes 更新文件公开状态响应
type FileUpdatePublicRes struct {
	g.Meta  `mime:"application/json" example:"json"`
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"提示信息"`
	Url     string `json:"url" dc:"文件新的访问链接"`
}
