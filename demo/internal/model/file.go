package model

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FileInfo 文件信息模型
type FileInfo struct {
	Id          int         `json:"id"`          // 文件ID
	Name        string      `json:"name"`        // 文件名
	Path        string      `json:"path"`        // 文件路径
	Size        int64       `json:"size"`        // 文件大小(字节)
	Type        string      `json:"type"`        // 文件类型
	ContentType string      `json:"contentType"` // 内容类型
	Extension   string      `json:"extension"`   // 扩展名
	IsPublic    bool        `json:"isPublic"`    // 是否公开
	Url         string      `json:"url"`         // 访问URL
	UserId      int         `json:"userId"`      // 上传用户ID
	Username    string      `json:"username"`    // 上传用户名
	CreatedAt   *gtime.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"`   // 更新时间
}

// FileDO 文件数据对象
type FileDO struct {
	Id          int64       `orm:"id,primary,auto_increment" json:"id"` // 文件ID
	Name        interface{} `orm:"name" json:"name"`                    // 文件名
	Path        interface{} `orm:"path" json:"path"`                    // 文件路径
	Size        interface{} `orm:"size" json:"size"`                    // 文件大小(字节)
	Type        interface{} `orm:"type" json:"type"`                    // 文件类型
	ContentType interface{} `orm:"content_type" json:"contentType"`     // 内容类型
	Extension   interface{} `orm:"extension" json:"extension"`          // 扩展名
	IsPublic    interface{} `orm:"is_public" json:"isPublic"`           // 是否公开
	UserId      interface{} `orm:"user_id" json:"userId"`               // 上传用户ID
	Username    interface{} `orm:"username" json:"username"`            // 上传用户名
	CreatedAt   *gtime.Time `orm:"created_at" json:"createdAt"`         // 创建时间
	UpdatedAt   *gtime.Time `orm:"updated_at" json:"updatedAt"`         // 更新时间
}

// FileUploadInfo 文件上传信息
type FileUploadInfo struct {
	OriginalName string // 原始文件名
	Size         int64  // 文件大小(字节)
	ContentType  string // 内容类型
	Path         string // 存储路径
	IsPublic     bool   // 是否公开
	Extension    string // 扩展名
	Type         string // 文件类型(image/document/video等)
	URL          string // 访问URL
}

// 文件表相关常量
const (
	// TableFile 表名
	TableFile = "file"
	// FileColumns 所有字段
	FileColumns = "id,name,path,size,type,content_type,extension,is_public,user_id,username,created_at,updated_at"
)
