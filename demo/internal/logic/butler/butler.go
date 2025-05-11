package butler

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	v1 "demo/api/butler/v1"
	"demo/internal/dao"
	"demo/internal/model/do"
	"demo/internal/service"
)

// butlerImpl 专属管家服务实现
type butlerImpl struct {
	butlerDao dao.ButlerDao
}

// New 创建专属管家服务实例
func New() service.Butler {
	return &butlerImpl{
		butlerDao: dao.NewButlerDao(),
	}
}

// SaveImage 保存专属管家图片
func (s *butlerImpl) SaveImage(ctx context.Context, req *v1.SaveImageReq) (res *v1.SaveImageRes, err error) {
	g.Log().Debug(ctx, "butler.SaveImage", req)
	res = &v1.SaveImageRes{}

	// 构建数据
	data := &do.ButlerDO{
		ImageUrl:  req.ImageUrl,
		Status:    req.Status,
		CreatedAt: gtime.Now(),
	}

	// 保存图片
	id, err := s.butlerDao.SaveImage(ctx, data)
	if err != nil {
		g.Log().Error(ctx, "保存专属管家图片失败", err)
		return nil, gerror.New("保存专属管家图片失败")
	}

	res.Id = int(id)

	return res, nil
}

// GetImage 获取专属管家图片
func (s *butlerImpl) GetImage(ctx context.Context, req *v1.GetImageReq) (res *v1.GetImageRes, err error) {
	g.Log().Debug(ctx, "butler.GetImage", req)
	res = &v1.GetImageRes{}

	// 获取最新图片
	butler, err := s.butlerDao.GetLatestImage(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取专属管家图片失败", err)
		return nil, gerror.New("获取专属管家图片失败")
	}

	// 如果没有启用的图片，返回空
	if butler == nil {
		return res, nil
	}

	// 设置返回值
	res.ImageUrl = gconv.String(butler.ImageUrl)

	return res, nil
}
