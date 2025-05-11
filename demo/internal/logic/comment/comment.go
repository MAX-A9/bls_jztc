package comment

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	v1 "demo/api/comment/v1"
	"demo/internal/dao"
	"demo/internal/model/do"
	"demo/internal/service"
)

// sComment 评论服务实现
type sComment struct{}

// New 创建评论服务实例
func New() service.CommentService {
	return &sComment{}
}

// List 评论列表
func (s *sComment) List(ctx context.Context, req *v1.CommentListReq) (res *v1.CommentListRes, err error) {
	res = &v1.CommentListRes{
		List:  make([]v1.CommentListItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 构建过滤条件
	filter := make(map[string]interface{})
	if req.ContentId > 0 {
		filter["content_id"] = req.ContentId
	}
	if req.Status != "" {
		filter["status"] = req.Status
	}
	if req.RealName != "" {
		filter["real_name"] = req.RealName
	}
	if req.Comment != "" {
		filter["comment"] = req.Comment
	}

	// 查询数据
	commentDao := dao.NewContentCommentDao()
	list, total, err := commentDao.FindList(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		return nil, gerror.New("获取评论列表失败: " + err.Error())
	}

	// 设置响应
	res.Total = total
	res.Page = req.Page

	// 转换数据
	if len(list) > 0 {
		// 创建内容ID到评论项的映射，用于后续查询内容标题
		contentIds := make([]interface{}, 0, len(list))
		idToIndex := make(map[int]int)

		for i, item := range list {
			contentId := gconv.Int(item.ContentId)
			if contentId > 0 {
				found := false
				for _, id := range contentIds {
					if gconv.Int(id) == contentId {
						found = true
						break
					}
				}

				if !found {
					contentIds = append(contentIds, contentId)
				}

				idToIndex[contentId] = i
			}
		}

		// 查询内容标题
		contentDao := dao.NewContentDao()
		contentTitles := make(map[int]string)

		if len(contentIds) > 0 {
			contentList, err := contentDao.FindByIds(ctx, contentIds)
			if err == nil && len(contentList) > 0 {
				for _, content := range contentList {
					contentId := gconv.Int(content.Id)
					contentTitles[contentId] = gconv.String(content.Title)
				}
			}
		}

		// 转换评论数据
		for _, item := range list {
			statusText := ""
			switch gconv.String(item.Status) {
			case "已审核":
				statusText = "已审核"
			case "待审核":
				statusText = "待审核"
			case "已拒绝":
				statusText = "已拒绝"
			default:
				statusText = gconv.String(item.Status)
			}

			contentId := gconv.Int(item.ContentId)
			contentTitle := contentTitles[contentId]

			listItem := v1.CommentListItem{
				Id:           gconv.Int(item.Id),
				ContentId:    contentId,
				ContentTitle: contentTitle,
				ClientId:     gconv.Int(item.ClientId),
				RealName:     gconv.String(item.RealName),
				Comment:      gconv.String(item.Comment),
				Status:       gconv.String(item.Status),
				StatusText:   statusText,
				CreatedAt:    item.CreatedAt.String(),
				UpdatedAt:    item.UpdatedAt.String(),
			}
			res.List = append(res.List, listItem)
		}
	}

	return res, nil
}

// Detail 评论详情
func (s *sComment) Detail(ctx context.Context, req *v1.CommentDetailReq) (res *v1.CommentDetailRes, err error) {
	res = &v1.CommentDetailRes{}

	// 查询数据
	commentDao := dao.NewContentCommentDao()
	comment, err := commentDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("获取评论详情失败: " + err.Error())
	}
	if comment == nil {
		return nil, gerror.New("评论不存在")
	}

	// 设置基本字段
	res.Id = gconv.Int(comment.Id)
	res.ContentId = gconv.Int(comment.ContentId)
	res.ClientId = gconv.Int(comment.ClientId)
	res.RealName = gconv.String(comment.RealName)
	res.Comment = gconv.String(comment.Comment)
	res.Status = gconv.String(comment.Status)
	res.CreatedAt = comment.CreatedAt.String()
	res.UpdatedAt = comment.UpdatedAt.String()

	// 查询内容标题
	contentDao := dao.NewContentDao()
	content, err := contentDao.FindOne(ctx, gconv.Int(comment.ContentId))
	if err == nil && content != nil {
		res.ContentTitle = gconv.String(content.Title)
	}

	return res, nil
}

// ContentComments 内容评论列表
func (s *sComment) ContentComments(ctx context.Context, req *v1.ContentCommentsReq) (res *v1.ContentCommentsRes, err error) {
	res = &v1.ContentCommentsRes{
		ContentId: req.ContentId,
		List:      make([]v1.CommentListItem, 0),
		Total:     0,
		Page:      req.Page,
	}

	// 查询内容是否存在
	contentDao := dao.NewContentDao()
	content, err := contentDao.FindOne(ctx, req.ContentId)
	if err != nil {
		return nil, gerror.New("获取内容信息失败: " + err.Error())
	}
	if content == nil {
		return nil, gerror.New("内容不存在")
	}

	// 设置内容标题
	res.ContentTitle = gconv.String(content.Title)

	// 查询评论数据
	commentDao := dao.NewContentCommentDao()
	list, total, err := commentDao.FindByContentId(ctx, req.ContentId, req.Page, req.PageSize)
	if err != nil {
		return nil, gerror.New("获取评论列表失败: " + err.Error())
	}

	// 设置响应数据
	res.Total = total
	res.Page = req.Page

	// 转换评论数据
	for _, item := range list {
		statusText := ""
		switch gconv.String(item.Status) {
		case "已审核":
			statusText = "已审核"
		case "待审核":
			statusText = "待审核"
		case "已拒绝":
			statusText = "已拒绝"
		default:
			statusText = gconv.String(item.Status)
		}

		listItem := v1.CommentListItem{
			Id:           gconv.Int(item.Id),
			ContentId:    req.ContentId,
			ContentTitle: res.ContentTitle,
			ClientId:     gconv.Int(item.ClientId),
			RealName:     gconv.String(item.RealName),
			Comment:      gconv.String(item.Comment),
			Status:       gconv.String(item.Status),
			StatusText:   statusText,
			CreatedAt:    item.CreatedAt.String(),
			UpdatedAt:    item.UpdatedAt.String(),
		}
		res.List = append(res.List, listItem)
	}

	return res, nil
}

// Create 创建评论
func (s *sComment) Create(ctx context.Context, req *v1.CommentCreateReq) (res *v1.CommentCreateRes, err error) {
	res = &v1.CommentCreateRes{}

	// 查询内容是否存在
	contentDao := dao.NewContentDao()
	content, err := contentDao.FindOne(ctx, req.ContentId)
	if err != nil {
		return nil, gerror.New("查询内容失败: " + err.Error())
	}
	if content == nil {
		return nil, gerror.New("内容不存在")
	}

	// 构建数据
	data := &do.ContentCommentDO{
		ContentId: req.ContentId,
		ClientId:  req.ClientId,
		RealName:  req.RealName,
		Comment:   req.Comment,
		Status:    req.Status,
		CreatedAt: gtime.Now(),
		UpdatedAt: gtime.Now(),
	}

	// 插入数据
	commentDao := dao.NewContentCommentDao()
	lastInsertId, err := commentDao.Insert(ctx, data)
	if err != nil {
		return nil, gerror.New("创建评论失败: " + err.Error())
	}

	// 如果评论状态为已审核，则增加内容的评论数
	if req.Status == "已审核" {
		if err := commentDao.IncrementCommentCount(ctx, req.ContentId); err != nil {
			g.Log().Warning(ctx, "增加内容评论数失败:", err.Error())
		}
	}

	res.Id = int(lastInsertId)
	return res, nil
}

// Update 更新评论
func (s *sComment) Update(ctx context.Context, req *v1.CommentUpdateReq) (res *v1.CommentUpdateRes, err error) {
	res = &v1.CommentUpdateRes{}

	// 查询原有数据
	commentDao := dao.NewContentCommentDao()
	exists, err := commentDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询评论失败: " + err.Error())
	}
	if exists == nil {
		return nil, gerror.New("评论不存在")
	}

	// 判断评论状态是否变更
	oldStatus := gconv.String(exists.Status)
	newStatus := req.Status
	statusChanged := oldStatus != newStatus

	// 构建更新数据
	data := &do.ContentCommentDO{
		RealName:  req.RealName,
		Comment:   req.Comment,
		Status:    req.Status,
		UpdatedAt: gtime.Now(),
	}

	// 更新数据
	_, err = commentDao.Update(ctx, data, req.Id)
	if err != nil {
		return nil, gerror.New("更新评论失败: " + err.Error())
	}

	// 如果状态变更，需要更新内容的评论数
	if statusChanged {
		contentId := gconv.Int(exists.ContentId)
		if oldStatus != "已审核" && newStatus == "已审核" {
			// 状态从非审核变为已审核，评论数+1
			if err := commentDao.IncrementCommentCount(ctx, contentId); err != nil {
				g.Log().Warning(ctx, "增加内容评论数失败:", err.Error())
			}
		} else if oldStatus == "已审核" && newStatus != "已审核" {
			// 状态从已审核变为非审核，评论数-1
			if err := commentDao.DecrementCommentCount(ctx, contentId); err != nil {
				g.Log().Warning(ctx, "减少内容评论数失败:", err.Error())
			}
		}
	}

	return res, nil
}

// Delete 删除评论
func (s *sComment) Delete(ctx context.Context, req *v1.CommentDeleteReq) (res *v1.CommentDeleteRes, err error) {
	res = &v1.CommentDeleteRes{}

	// 查询原有数据
	commentDao := dao.NewContentCommentDao()
	exists, err := commentDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询评论失败: " + err.Error())
	}
	if exists == nil {
		return nil, gerror.New("评论不存在")
	}

	// 记录评论状态和内容ID
	status := gconv.String(exists.Status)
	contentId := gconv.Int(exists.ContentId)

	// 删除数据（软删除）
	_, err = commentDao.Delete(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("删除评论失败: " + err.Error())
	}

	// 如果删除的是已审核状态的评论，减少内容的评论数
	if status == "已审核" {
		if err := commentDao.DecrementCommentCount(ctx, contentId); err != nil {
			g.Log().Warning(ctx, "减少内容评论数失败:", err.Error())
		}
	}

	return res, nil
}

// UpdateStatus 更新评论状态
func (s *sComment) UpdateStatus(ctx context.Context, req *v1.CommentStatusUpdateReq) (res *v1.CommentStatusUpdateRes, err error) {
	res = &v1.CommentStatusUpdateRes{}

	// 查询原有数据
	commentDao := dao.NewContentCommentDao()
	exists, err := commentDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询评论失败: " + err.Error())
	}
	if exists == nil {
		return nil, gerror.New("评论不存在")
	}

	// 记录原状态和内容ID
	oldStatus := gconv.String(exists.Status)
	contentId := gconv.Int(exists.ContentId)

	// 如果状态未变更，直接返回
	if oldStatus == req.Status {
		return res, nil
	}

	// 更新状态
	_, err = commentDao.UpdateStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, gerror.New("更新评论状态失败: " + err.Error())
	}

	// 更新内容评论数
	if oldStatus != "已审核" && req.Status == "已审核" {
		// 从非审核变为已审核，评论数+1
		if err := commentDao.IncrementCommentCount(ctx, contentId); err != nil {
			g.Log().Warning(ctx, "增加内容评论数失败:", err.Error())
		}
	} else if oldStatus == "已审核" && req.Status != "已审核" {
		// 从已审核变为非审核，评论数-1
		if err := commentDao.DecrementCommentCount(ctx, contentId); err != nil {
			g.Log().Warning(ctx, "减少内容评论数失败:", err.Error())
		}
	}

	return res, nil
}

// WxClientList 微信客户端评论列表
func (s *sComment) WxClientList(ctx context.Context, req *v1.WxClientCommentListReq) (res *v1.WxClientCommentListRes, err error) {
	res = &v1.WxClientCommentListRes{
		List:  make([]v1.WxClientCommentItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 查询评论数据
	commentDao := dao.NewContentCommentDao()
	filter := map[string]interface{}{
		"content_id": req.ContentId,
		"status":     "已审核", // 只获取已审核的评论
	}
	list, total, err := commentDao.FindList(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		return nil, gerror.New("获取评论列表失败: " + err.Error())
	}

	// 设置响应数据
	res.Total = total
	res.Page = req.Page

	// 获取所有客户ID用于批量查询头像
	clientIds := make([]interface{}, 0, len(list))
	clientMap := make(map[int]string) // 客户ID到头像URL的映射

	for _, item := range list {
		clientId := gconv.Int(item.ClientId)
		if clientId > 0 {
			found := false
			for _, id := range clientIds {
				if gconv.Int(id) == clientId {
					found = true
					break
				}
			}
			if !found {
				clientIds = append(clientIds, clientId)
			}
		}
	}

	// 批量查询客户头像
	if len(clientIds) > 0 {
		var clients []struct {
			Id        int    `json:"id"`
			AvatarUrl string `json:"avatar_url"`
		}
		err = g.DB().Model("client").
			Fields("id, avatar_url").
			WhereIn("id", clientIds).
			Scan(&clients)

		if err == nil && len(clients) > 0 {
			for _, client := range clients {
				clientMap[client.Id] = client.AvatarUrl
			}
		}
	}

	// 转换评论数据
	for _, item := range list {
		clientId := gconv.Int(item.ClientId)
		avatarUrl := clientMap[clientId] // 获取头像URL，如果没有则为空字符串

		listItem := v1.WxClientCommentItem{
			Id:        gconv.Int(item.Id),
			RealName:  gconv.String(item.RealName),
			AvatarUrl: avatarUrl,
			Comment:   gconv.String(item.Comment),
			CreatedAt: item.CreatedAt.String(),
		}
		res.List = append(res.List, listItem)
	}

	return res, nil
}

// WxClientCreate 微信客户端创建评论
func (s *sComment) WxClientCreate(ctx context.Context, req *v1.WxClientCommentCreateReq) (res *v1.WxClientCommentCreateRes, err error) {
	res = &v1.WxClientCommentCreateRes{}

	// 从上下文获取客户信息
	clientId := gconv.Int(ctx.Value("client_id"))
	if clientId <= 0 {
		return nil, gerror.New("未获取到客户ID，请先登录")
	}

	// 查询客户信息
	var client struct {
		Username string `json:"username"`
		RealName string `json:"real_name"`
	}
	err = g.DB().Model("client").
		Fields("username, real_name").
		Where("id", clientId).
		Scan(&client)
	if err != nil {
		return nil, gerror.New("获取客户信息失败: " + err.Error())
	}

	// 设置真实姓名，优先使用真实姓名
	realName := client.RealName
	if realName == "" {
		realName = client.Username
	}

	// 查询内容是否存在
	contentDao := dao.NewContentDao()
	content, err := contentDao.FindOne(ctx, req.ContentId)
	if err != nil {
		return nil, gerror.New("查询内容失败: " + err.Error())
	}
	if content == nil {
		return nil, gerror.New("内容不存在")
	}

	// 构建评论数据
	data := &do.ContentCommentDO{
		ContentId: req.ContentId,
		ClientId:  clientId,
		RealName:  realName,
		Comment:   req.Comment,
		Status:    "待审核", // 默认设置为待审核状态
		CreatedAt: gtime.Now(),
		UpdatedAt: gtime.Now(),
	}

	// 插入数据
	commentDao := dao.NewContentCommentDao()
	lastInsertId, err := commentDao.Insert(ctx, data)
	if err != nil {
		return nil, gerror.New("创建评论失败: " + err.Error())
	}

	res.Id = int(lastInsertId)
	return res, nil
}
