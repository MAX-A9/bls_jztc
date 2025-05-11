package client

import (
	"context"

	v1 "demo/api/client/v1"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PublisherInfo 获取发布人信息
func (s *sClient) PublisherInfo(ctx context.Context, req *v1.PublisherInfoReq) (res *v1.PublisherInfoRes, err error) {
	// 初始化响应对象
	res = &v1.PublisherInfoRes{}

	// 获取发布人基本信息
	var publisherInfo struct {
		RealName  string `json:"real_name"`
		AvatarUrl string `json:"avatar_url"`
	}
	err = g.DB().Model("client").
		Fields("real_name, avatar_url").
		Where("id", req.PublisherId).
		Where("status", 1). // 仅查询状态为正常的客户
		Scan(&publisherInfo)
	if err != nil {
		return nil, gerror.New("获取发布人信息失败")
	}

	// 如果没有找到发布人信息
	if publisherInfo.RealName == "" {
		return nil, gerror.New("发布人不存在或已被禁用")
	}

	// 填充基本信息
	res.RealName = publisherInfo.RealName
	res.AvatarUrl = publisherInfo.AvatarUrl

	// 查询该发布人关注的人数
	followNum, err := g.DB().Model("publisher_follow").
		Where("client_id", req.PublisherId).
		Count()
	if err != nil {
		// 查询失败不阻断流程，给默认值0
		followNum = 0
	}
	res.FollowNum = followNum

	// 查询关注该发布人的粉丝数
	fansNum, err := g.DB().Model("publisher_follow").
		Where("publisher_id", req.PublisherId).
		Count()
	if err != nil {
		// 查询失败不阻断流程，给默认值0
		fansNum = 0
	}
	res.FansNum = fansNum

	// 查询发布内容数量
	publishCount, err := g.DB().Model("content").
		Where("client_id", req.PublisherId).
		Where("status", "已发布"). // 只统计已发布的内容
		Count()
	if err != nil {
		// 查询失败不阻断流程，给默认值0
		publishCount = 0
	}
	res.PublishCount = publishCount

	return res, nil
}

// FollowPublisher 关注发布人
func (s *sClient) FollowPublisher(ctx context.Context, req *v1.FollowPublisherReq) (res *v1.FollowPublisherRes, err error) {
	// 初始化响应对象
	res = &v1.FollowPublisherRes{Success: false}

	// 获取当前登录用户ID
	currentClientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 不能关注自己
	if uint(currentClientId) == req.PublisherId {
		return nil, gerror.New("不能关注自己")
	}

	// 检查目标发布人是否存在
	publisherCount, err := g.DB().Model("client").
		Where("id", req.PublisherId).
		Where("status", 1). // 仅关注状态为正常的客户
		Count()
	if err != nil {
		return nil, gerror.New("查询发布人信息失败")
	}
	if publisherCount == 0 {
		return nil, gerror.New("发布人不存在或已被禁用")
	}

	// 检查是否已关注
	favoriteCount, err := g.DB().Model("publisher_follow").
		Where("client_id", currentClientId).
		Where("publisher_id", req.PublisherId).
		Count()
	if err != nil {
		return nil, gerror.New("查询关注状态失败")
	}
	if favoriteCount > 0 {
		// 已关注，返回成功
		res.Success = true
		return res, nil
	}

	// 添加关注关系
	_, err = g.DB().Model("publisher_follow").Insert(g.Map{
		"client_id":    currentClientId,
		"publisher_id": req.PublisherId,
		"created_at":   gtime.Now(),
	})
	if err != nil {
		return nil, gerror.New("关注失败")
	}

	res.Success = true
	return res, nil
}

// UnfollowPublisher 取消关注发布人
func (s *sClient) UnfollowPublisher(ctx context.Context, req *v1.UnfollowPublisherReq) (res *v1.UnfollowPublisherRes, err error) {
	// 初始化响应对象
	res = &v1.UnfollowPublisherRes{Success: false}

	// 获取当前登录用户ID
	currentClientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 删除关注关系
	_, err = g.DB().Model("publisher_follow").
		Where("client_id", currentClientId).
		Where("publisher_id", req.PublisherId).
		Delete()
	if err != nil {
		return nil, gerror.New("取消关注失败")
	}

	res.Success = true
	return res, nil
}

// FollowingList 获取关注人列表
func (s *sClient) FollowingList(ctx context.Context, req *v1.FollowingListReq) (res *v1.FollowingListRes, err error) {
	// 初始化响应对象
	res = &v1.FollowingListRes{
		List:  make([]v1.FollowingItem, 0),
		Page:  req.Page,
		Size:  req.Size,
		Total: 0,
	}

	// 获取当前登录用户ID
	currentClientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 计算分页参数
	offset := (req.Page - 1) * req.Size

	// 查询我关注的发布人总数
	total, err := g.DB().Model("publisher_follow").
		Where("client_id", currentClientId).
		Count()
	if err != nil {
		return nil, gerror.New("查询关注列表失败")
	}
	res.Total = total

	// 如果没有关注任何发布人，直接返回空列表
	if total == 0 {
		return res, nil
	}

	// 查询发布人列表
	var followingList []struct {
		PublisherId uint        `json:"publisher_id"`
		CreatedAt   *gtime.Time `json:"created_at"`
	}

	err = g.DB().Model("publisher_follow").
		Fields("publisher_id, created_at").
		Where("client_id", currentClientId).
		Order("created_at DESC"). // 按关注时间倒序排列，最新关注的排在前面
		Limit(req.Size).
		Offset(offset).
		Scan(&followingList)
	if err != nil {
		return nil, gerror.New("查询关注列表失败")
	}

	// 如果没有数据，返回空列表
	if len(followingList) == 0 {
		return res, nil
	}

	// 提取所有发布人ID
	var publisherIds []uint
	for _, item := range followingList {
		publisherIds = append(publisherIds, item.PublisherId)
	}

	// 查询这些发布人的基本信息
	var publisherInfoList []struct {
		Id        uint   `json:"id"`
		RealName  string `json:"real_name"`
		AvatarUrl string `json:"avatar_url"`
	}

	err = g.DB().Model("client").
		Fields("id, real_name, avatar_url").
		WhereIn("id", publisherIds).
		Where("status", 1). // 只查询状态正常的发布人
		Scan(&publisherInfoList)
	if err != nil {
		return nil, gerror.New("查询发布人信息失败")
	}

	// 构建发布人基本信息映射，方便快速查找
	publisherMap := make(map[uint]struct {
		RealName  string
		AvatarUrl string
	})

	for _, info := range publisherInfoList {
		publisherMap[info.Id] = struct {
			RealName  string
			AvatarUrl string
		}{
			RealName:  info.RealName,
			AvatarUrl: info.AvatarUrl,
		}
	}

	// 查询每个发布人的发布内容数量
	type PublishCount struct {
		ClientId uint `json:"client_id"`
		Count    int  `json:"count"`
	}

	var publishCountList []PublishCount
	err = g.DB().Model("content").
		Fields("client_id, COUNT(1) as count").
		WhereIn("client_id", publisherIds).
		Where("status", "已发布").
		Group("client_id").
		Scan(&publishCountList)

	// 如果查询出错，记录错误但继续处理（将使用空列表）
	if err != nil {
		g.Log().Warning(ctx, "查询发布人内容数量失败:", err)
	}

	// 构建发布数量映射
	publishCountMap := make(map[uint]int)
	for _, item := range publishCountList {
		publishCountMap[item.ClientId] = item.Count
	}

	// 组装最终的响应列表
	for _, follow := range followingList {
		// 查找发布人信息
		publisherInfo, exists := publisherMap[follow.PublisherId]
		if !exists {
			// 如果找不到发布人信息（可能已被删除），则跳过
			continue
		}

		// 获取发布内容数量，如果没有则默认为0
		publishCount := publishCountMap[follow.PublisherId]

		// 添加到结果列表
		item := v1.FollowingItem{
			ClientId:     follow.PublisherId,
			RealName:     publisherInfo.RealName,
			AvatarUrl:    publisherInfo.AvatarUrl,
			PublishCount: publishCount,
			FollowTime:   follow.CreatedAt.Format("Y-m-d H:i:s"),
		}

		res.List = append(res.List, item)
	}

	return res, nil
}

// FollowingCount 获取关注人总数
func (s *sClient) FollowingCount(ctx context.Context, req *v1.FollowingCountReq) (res *v1.FollowingCountRes, err error) {
	// 初始化响应对象
	res = &v1.FollowingCountRes{Count: 0}

	// 获取当前登录用户ID
	currentClientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 查询关注人总数
	count, err := g.DB().Model("publisher_follow").
		Where("client_id", currentClientId).
		Count()
	if err != nil {
		return nil, gerror.New("查询关注人数量失败")
	}

	res.Count = count
	return res, nil
}

// FollowStatus 获取关注状态
func (s *sClient) FollowStatus(ctx context.Context, req *v1.FollowStatusReq) (res *v1.FollowStatusRes, err error) {
	// 初始化响应对象
	res = &v1.FollowStatusRes{IsFollowed: false}

	// 获取当前登录用户ID
	currentClientId, err := auth.GetClientId(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 检查目标发布人是否存在
	publisherCount, err := g.DB().Model("client").
		Where("id", req.PublisherId).
		Where("status", 1). // 仅查询状态为正常的客户
		Count()
	if err != nil {
		return nil, gerror.New("查询发布人信息失败")
	}
	if publisherCount == 0 {
		return nil, gerror.New("发布人不存在或已被禁用")
	}

	// 查询是否已关注
	followCount, err := g.DB().Model("publisher_follow").
		Where("client_id", currentClientId).
		Where("publisher_id", req.PublisherId).
		Count()
	if err != nil {
		return nil, gerror.New("查询关注状态失败")
	}

	// 设置关注状态
	res.IsFollowed = followCount > 0

	return res, nil
}
