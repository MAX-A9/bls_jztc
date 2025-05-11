package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// 地区列表请求
type RegionListReq struct {
	g.Meta   `path:"/region/list" method:"get" tags:"地区管理" summary:"获取地区列表" security:"Bearer" description:"获取地区列表，分页查询"`
	Page     int    `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize int    `v:"max:100#每页最大100条" json:"pageSize" dc:"每页数量"`
	Name     string `json:"name" dc:"地区名称，模糊搜索"`
	Level    string `json:"level" dc:"地区级别: 省,县,乡"`
	Location string `json:"location" dc:"所在地区，模糊搜索"`
	Status   int    `json:"status" default:"-1" dc:"状态 0:启用 1:禁用 -1:全部"`
}

// 地区列表响应
type RegionListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []RegionListItem `json:"list" dc:"地区列表"`
	Total  int              `json:"total" dc:"总数量"`
	Page   int              `json:"page" dc:"当前页码"`
}

// 地区列表项
type RegionListItem struct {
	Id        int         `json:"id" dc:"地区ID"`
	Location  string      `json:"location" dc:"所在地区"`
	Name      string      `json:"name" dc:"地区名称"`
	Level     string      `json:"level" dc:"级别: 省,县,乡"`
	Status    int         `json:"status" dc:"状态 0:启用 1:禁用"`
	CreatedAt *gtime.Time `json:"createdAt" dc:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" dc:"更新时间"`
}

// 地区详情请求
type RegionDetailReq struct {
	g.Meta `path:"/region/detail" method:"get" tags:"地区管理" summary:"获取地区详情" security:"Bearer" description:"根据ID获取地区详情"`
	Id     int `v:"required#地区ID不能为空" json:"id" dc:"地区ID"`
}

// 地区详情响应
type RegionDetailRes struct {
	g.Meta    `mime:"application/json" example:"json"`
	Id        int         `json:"id" dc:"地区ID"`
	Location  string      `json:"location" dc:"所在地区"`
	Name      string      `json:"name" dc:"地区名称"`
	Level     string      `json:"level" dc:"级别: 省,县,乡"`
	Status    int         `json:"status" dc:"状态 0:启用 1:禁用"`
	CreatedAt *gtime.Time `json:"createdAt" dc:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" dc:"更新时间"`
}

// 创建地区请求
type RegionCreateReq struct {
	g.Meta   `path:"/region/create" method:"post" tags:"地区管理" summary:"创建地区" security:"Bearer" description:"创建新地区，可通过级联选择器选择省市县，系统会自动确定当前地区的层级"`
	Location string `v:"required#所在地区不能为空" json:"location" dc:"所在地区，通过级联选择器选择，格式如：北京市/朝阳区"`
	Name     string `json:"name" dc:"地区名称，如级联最后一级或补充的乡镇名称"`
	Level    string `v:"required|in:省,县,乡#级别不能为空|级别必须是'省','县','乡'之一" json:"level" dc:"级别: 省,县,乡"`
	Status   int    `v:"required|in:0,1#状态不能为空|状态必须是0或1" json:"status" dc:"状态 0:启用 1:禁用"`
}

// 创建地区响应
type RegionCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"新创建的地区ID"`
}

// 更新地区请求
type RegionUpdateReq struct {
	g.Meta   `path:"/region/update" method:"put" tags:"地区管理" summary:"更新地区" security:"Bearer" description:"更新地区信息，可通过级联选择器选择省市县"`
	Id       int    `v:"required#地区ID不能为空" json:"id" dc:"地区ID"`
	Location string `v:"required#所在地区不能为空" json:"location" dc:"所在地区，通过级联选择器选择，格式如：北京市/朝阳区"`
	Name     string `json:"name" dc:"地区名称，如级联最后一级或补充的乡镇名称"`
	Level    string `v:"required|in:省,县,乡#级别不能为空|级别必须是'省','县','乡'之一" json:"level" dc:"级别: 省,县,乡"`
	Status   int    `v:"required|in:0,1#状态不能为空|状态必须是0或1" json:"status" dc:"状态 0:启用 1:禁用"`
}

// 更新地区响应
type RegionUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// 删除地区请求
type RegionDeleteReq struct {
	g.Meta `path:"/region/delete" method:"delete" tags:"地区管理" summary:"删除地区" security:"Bearer" description:"删除指定地区"`
	Id     int `v:"required#地区ID不能为空" json:"id" dc:"地区ID"`
}

// 删除地区响应
type RegionDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// 客户端地区列表请求
type WxClientRegionListReq struct {
	g.Meta `path:"/wx/client/region/list" method:"get" tags:"客户端地区" summary:"获取地区列表" description:"微信小程序获取地区列表，根据状态筛选，默认只返回开通的地区"`
	Status int `json:"status" default:"0" dc:"状态 0:启用 1:禁用"`
}

// 客户端地区列表响应
type WxClientRegionListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []RegionListItem `json:"list" dc:"地区列表"`
}
