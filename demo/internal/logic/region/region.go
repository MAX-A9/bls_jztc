package region

import (
	"context"
	"strings"

	v1 "demo/api/region/v1"
	"demo/internal/consts"
	"demo/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// regionImpl 地区服务实现
type regionImpl struct{}

// New 创建地区服务实例
func New() service.RegionService {
	return &regionImpl{}
}

// List 获取地区列表
func (s *regionImpl) List(ctx context.Context, req *v1.RegionListReq) (res *v1.RegionListRes, err error) {
	// 初始化响应
	res = &v1.RegionListRes{
		List:  make([]v1.RegionListItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 构建查询模型
	model := g.DB().Model("region").Safe()

	// 应用过滤条件
	if req.Name != "" {
		model = model.WhereLike("name", "%"+req.Name+"%")
	}

	if req.Level != "" {
		model = model.Where("level", req.Level)
	}

	if req.Location != "" {
		model = model.WhereLike("location", "%"+req.Location+"%")
	}

	// 状态过滤，只有当Status明确指定为0或1时才过滤
	// Status默认应为-1，表示查询所有状态
	if req.Status == 0 || req.Status == 1 {
		model = model.Where("status", req.Status)
	}

	// 统计总记录数
	count, err := model.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeRegionListFailed, err.Error())
	}
	res.Total = count

	// 没有记录，直接返回空数据
	if count == 0 {
		return res, nil
	}

	// 设置分页参数
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10 // 默认每页10条
	}

	page := req.Page
	if page <= 0 {
		page = 1 // 默认第1页
	}

	model = model.Page(page, pageSize).Order("id ASC")

	// 获取记录列表
	records, err := model.All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeRegionListFailed, err.Error())
	}

	// 转换记录为响应对象
	for _, record := range records {
		item := v1.RegionListItem{
			Id:       record["id"].Int(),
			Location: record["location"].String(),
			Name:     record["name"].String(),
			Level:    record["level"].String(),
			Status:   record["status"].Int(),
		}

		// 设置时间字段
		if !record["created_at"].IsEmpty() {
			createdAt := record["created_at"].GTime()
			item.CreatedAt = createdAt
		}

		if !record["updated_at"].IsEmpty() {
			updatedAt := record["updated_at"].GTime()
			item.UpdatedAt = updatedAt
		}

		res.List = append(res.List, item)
	}

	return res, nil
}

// Detail 获取地区详情
func (s *regionImpl) Detail(ctx context.Context, req *v1.RegionDetailReq) (res *v1.RegionDetailRes, err error) {
	// 获取地区详情
	record, err := g.DB().Model("region").Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeRegionDetailFailed, err.Error())
	}

	if record.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeRegionNotExists)
	}

	// 构建详情响应
	res = &v1.RegionDetailRes{
		Id:       record["id"].Int(),
		Location: record["location"].String(),
		Name:     record["name"].String(),
		Level:    record["level"].String(),
		Status:   record["status"].Int(),
	}

	// 设置时间字段
	if !record["created_at"].IsEmpty() {
		createdAt := record["created_at"].GTime()
		res.CreatedAt = createdAt
	}

	if !record["updated_at"].IsEmpty() {
		updatedAt := record["updated_at"].GTime()
		res.UpdatedAt = updatedAt
	}

	return res, nil
}

// Create 创建地区
func (s *regionImpl) Create(ctx context.Context, req *v1.RegionCreateReq) (res *v1.RegionCreateRes, err error) {
	res = &v1.RegionCreateRes{}

	// 处理地区名称，如果为空则使用Location的最后一级
	name := req.Name
	if name == "" {
		locationParts := strings.Split(req.Location, "/")
		if len(locationParts) > 0 {
			name = locationParts[len(locationParts)-1]
		}
	}

	// 检查该Location下是否已存在同名地区
	count, err := g.DB().Model("region").
		Where("location", req.Location).
		Where("name", name).
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeRegionCreateFailed, err.Error())
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeRegionNameExists)
	}

	// 检查级别与所在地区的格式匹配
	locationParts := strings.Split(req.Location, "/")
	if req.Level == "省" && len(locationParts) != 1 {
		return nil, gerror.NewCode(consts.CodeInvalidRegionLevel, "省级地区的所在地区格式应为'省名'")
	} else if req.Level == "县" && len(locationParts) != 2 {
		return nil, gerror.NewCode(consts.CodeInvalidRegionLevel, "县级地区的所在地区格式应为'省名/县名'")
	} else if req.Level == "乡" && len(locationParts) != 3 {
		return nil, gerror.NewCode(consts.CodeInvalidRegionLevel, "乡级地区的所在地区格式应为'省名/县名/乡名'")
	}

	// 准备插入数据
	data := g.Map{
		"location":   req.Location,
		"name":       name,
		"level":      req.Level,
		"status":     req.Status,
		"created_at": gtime.Now(),
	}

	// 插入数据
	result, err := g.DB().Model("region").Data(data).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeRegionCreateFailed, err.Error())
	}

	// 获取新插入记录的ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeRegionCreateFailed, err.Error())
	}

	res.Id = int(id)
	return res, nil
}

// Update 更新地区
func (s *regionImpl) Update(ctx context.Context, req *v1.RegionUpdateReq) (res *v1.RegionUpdateRes, err error) {
	res = &v1.RegionUpdateRes{}

	// 检查地区是否存在
	record, err := g.DB().Model("region").Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeRegionUpdateFailed, err.Error())
	}
	if record.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeRegionNotExists)
	}

	// 处理地区名称，如果为空则使用Location的最后一级
	name := req.Name
	if name == "" {
		locationParts := strings.Split(req.Location, "/")
		if len(locationParts) > 0 {
			name = locationParts[len(locationParts)-1]
		}
	}

	// 检查该Location下是否已存在同名地区(排除自身)
	count, err := g.DB().Model("region").
		Where("location", req.Location).
		Where("name", name).
		Where("id != ?", req.Id).
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeRegionUpdateFailed, err.Error())
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeRegionNameExists)
	}

	// 检查级别与所在地区的格式匹配
	locationParts := strings.Split(req.Location, "/")
	if req.Level == "省" && len(locationParts) != 1 {
		return nil, gerror.NewCode(consts.CodeInvalidRegionLevel, "省级地区的所在地区格式应为'省名'")
	} else if req.Level == "县" && len(locationParts) != 2 {
		return nil, gerror.NewCode(consts.CodeInvalidRegionLevel, "县级地区的所在地区格式应为'省名/县名'")
	} else if req.Level == "乡" && len(locationParts) != 3 {
		return nil, gerror.NewCode(consts.CodeInvalidRegionLevel, "乡级地区的所在地区格式应为'省名/县名/乡名'")
	}

	// 准备更新数据
	data := g.Map{
		"location":   req.Location,
		"name":       name,
		"level":      req.Level,
		"status":     req.Status,
		"updated_at": gtime.Now(),
	}

	// 更新地区
	_, err = g.DB().Model("region").Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeRegionUpdateFailed, err.Error())
	}

	return res, nil
}

// Delete 删除地区
func (s *regionImpl) Delete(ctx context.Context, req *v1.RegionDeleteReq) (res *v1.RegionDeleteRes, err error) {
	res = &v1.RegionDeleteRes{}

	// 检查地区是否存在
	record, err := g.DB().Model("region").Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeRegionDeleteFailed, err.Error())
	}

	if record.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeRegionNotExists)
	}

	// 直接删除地区
	_, err = g.DB().Model("region").Where("id", req.Id).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeRegionDeleteFailed, err.Error())
	}

	return res, nil
}

// WxClientList 客户端获取地区列表
func (s *regionImpl) WxClientList(ctx context.Context, req *v1.WxClientRegionListReq) (res *v1.WxClientRegionListRes, err error) {
	// 初始化响应
	res = &v1.WxClientRegionListRes{
		List: make([]v1.RegionListItem, 0),
	}

	// 构建查询模型，只查询未删除的记录
	model := g.DB().Model("region").Where("deleted_at IS NULL")

	// 应用状态过滤条件，默认只查询启用状态(status=0)的地区
	status := req.Status
	if status != 0 && status != 1 {
		status = 0 // 默认只查询启用状态
	}
	model = model.Where("status", status)

	// 按ID正序排序
	model = model.Order("id ASC")

	// 获取记录列表
	records, err := model.All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeRegionListFailed, err.Error())
	}

	// 转换记录为响应对象
	for _, record := range records {
		item := v1.RegionListItem{
			Id:       record["id"].Int(),
			Location: record["location"].String(),
			Name:     record["name"].String(),
			Level:    record["level"].String(),
			Status:   record["status"].Int(),
		}

		// 设置时间字段
		if !record["created_at"].IsEmpty() {
			createdAt := record["created_at"].GTime()
			item.CreatedAt = createdAt
		}

		if !record["updated_at"].IsEmpty() {
			updatedAt := record["updated_at"].GTime()
			item.UpdatedAt = updatedAt
		}

		res.List = append(res.List, item)
	}

	return res, nil
}
