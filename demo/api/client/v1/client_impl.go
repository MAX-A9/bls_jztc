package v1

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	"demo/internal/dao"
	"demo/utility/auth"
	"demo/utility/storage"
)

// ControllerImpl 控制器实现
type ControllerImpl struct{}

// List 获取客户列表
func (c *ControllerImpl) List(ctx context.Context, req *ClientListReq) (res *ClientListRes, err error) {
	res = &ClientListRes{
		List:  make([]ClientListItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 构建查询参数
	m := g.Map{}
	if req.Username != "" {
		m["username"] = req.Username
	}
	if req.RealName != "" {
		m["real_name"] = req.RealName
	}
	if req.Phone != "" {
		m["phone"] = req.Phone
	}
	if req.Status != 0 {
		m["status"] = req.Status
	}

	// 使用 DAO 层获取客户列表
	clientDao := dao.ClientDao()
	list, total, err := clientDao.GetWithPage(ctx, req.Page, req.PageSize, m)
	if err != nil {
		return nil, err
	}
	res.Total = total

	// 处理结果
	var items []ClientListItem
	if err = list.Structs(&items); err != nil {
		return nil, gerror.New("转换客户列表数据失败")
	}

	// 处理结果
	for _, item := range items {
		// 状态文本
		if item.Status == 0 {
			item.StatusText = "禁用"
		} else {
			item.StatusText = "正常"
		}

		// 来源标识文本
		if item.Identifier == "wxapp" {
			item.IdentifierText = "小程序"
		} else {
			item.IdentifierText = "未知"
		}

		res.List = append(res.List, item)
	}

	return
}

// Create 创建客户
func (c *ControllerImpl) Create(ctx context.Context, req *ClientCreateReq) (res *ClientCreateRes, err error) {
	res = &ClientCreateRes{}

	// 使用 DAO 层检查用户名是否存在
	clientDao := dao.ClientDao()
	count, err := clientDao.Model(ctx).Where("username", req.Username).Count()
	if err != nil {
		return nil, gerror.New("检查用户名失败")
	}
	if count > 0 {
		return nil, gerror.New("用户名已存在")
	}

	// 生成密码哈希
	passwordHash, err := gmd5.EncryptString(req.Password)
	if err != nil {
		return nil, gerror.New("密码加密失败")
	}

	// 构建插入数据
	data := g.Map{
		"username":   req.Username,
		"password":   passwordHash,
		"real_name":  req.RealName,
		"phone":      req.Phone,
		"status":     req.Status,
		"identifier": req.Identifier,
		"created_at": gtime.Now().String(),
	}

	// 使用 DAO 层插入数据
	lastId, err := clientDao.Insert(ctx, data)
	if err != nil {
		return nil, gerror.New("创建客户失败")
	}

	res.Id = int(lastId)
	return
}

// Update 更新客户
func (c *ControllerImpl) Update(ctx context.Context, req *ClientUpdateReq) (res *ClientUpdateRes, err error) {
	res = &ClientUpdateRes{}

	// 使用 DAO 层检查客户是否存在
	clientDao := dao.ClientDao()
	count, err := clientDao.Model(ctx).Where("id", req.Id).Count()
	if err != nil {
		return nil, gerror.New("检查客户失败")
	}
	if count == 0 {
		return nil, gerror.New("客户不存在")
	}

	// 检查用户名是否已被其他用户使用
	count, err = clientDao.Model(ctx).Where("username", req.Username).Where("id <> ?", req.Id).Count()
	if err != nil {
		return nil, gerror.New("检查用户名失败")
	}
	if count > 0 {
		return nil, gerror.New("用户名已存在")
	}

	// 获取客户当前信息
	var oldClientInfo struct {
		RealName string `json:"real_name"`
	}
	err = clientDao.Model(ctx).
		Fields("real_name").
		Where("id", req.Id).
		Scan(&oldClientInfo)
	if err != nil {
		// 继续执行，不中断流程
	}
	oldRealName := oldClientInfo.RealName

	// 构建更新数据
	data := g.Map{
		"username":  req.Username,
		"real_name": req.RealName,
		"phone":     req.Phone,
		"status":    req.Status,
	}

	// 使用 DAO 层更新数据
	err = clientDao.UpdateById(ctx, req.Id, data)
	if err != nil {
		return nil, gerror.New("更新客户失败")
	}

	// 如果更新了真实姓名，同步更新内容表中的作者字段
	if req.RealName != oldRealName {
		// 执行更新内容表的作者字段
		_, err = g.DB().Model("content").
			Data(g.Map{"author": req.RealName}).
			Where("client_id", req.Id).
			Update()
		if err != nil {
			g.Log().Warning(ctx, "更新内容作者信息失败: "+err.Error())
			// 不中断主流程，继续执行
		}
	}

	return
}

// Delete 删除客户
func (c *ControllerImpl) Delete(ctx context.Context, req *ClientDeleteReq) (res *ClientDeleteRes, err error) {
	res = &ClientDeleteRes{}

	// 使用 DAO 层检查客户是否存在
	clientDao := dao.ClientDao()
	count, err := clientDao.Model(ctx).Where("id", req.Id).Count()
	if err != nil {
		return nil, gerror.New("检查客户失败")
	}
	if count == 0 {
		return nil, gerror.New("客户不存在")
	}

	// 使用 DAO 层删除客户
	err = clientDao.DeleteById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("删除客户失败")
	}

	return
}

// WxappLogin 微信小程序登录
func (c *ControllerImpl) WxappLogin(ctx context.Context, req *WxappLoginReq) (res *WxappLoginRes, err error) {
	res = &WxappLoginRes{}

	// 获取微信小程序配置
	wxConfig, err := loadWxappConfig()
	if err != nil {
		return nil, gerror.New("获取微信小程序配置失败")
	}

	// 调用微信API获取openid和session_key
	openid, sessionKey, err := getWxOpenId(ctx, req.Code, wxConfig.AppId, wxConfig.AppSecret)
	if err != nil {
		return nil, gerror.New("微信登录失败，请重试")
	}

	// 如果无法获取openid，返回错误
	if openid == "" {
		return nil, gerror.New("微信登录失败，无法获取用户标识")
	}

	// 检查openid是否存在对应的客户
	var client struct {
		Id         int    `json:"id"`
		Username   string `json:"username"`
		Status     int    `json:"status"`
		OpenId     string `json:"open_id"`
		SessionKey string `json:"session_key"`
	}

	// 查询客户信息
	err = g.DB().Model("client").Where("open_id", openid).Scan(&client)

	// 是否为新用户
	isNewUser := err != nil || client.Id == 0

	// 处理用户注册或登录逻辑
	if isNewUser {
		// 生成随机用户名
		username := "wx_" + grand.S(10)

		// 生成随机密码
		password := grand.S(16)
		passwordHash, err := gmd5.EncryptString(password)
		if err != nil {
			return nil, gerror.New("密码加密失败")
		}

		// 设置真实姓名为"微信用户_随机字符串"
		randomSuffix := grand.S(5)
		realName := "微信用户_" + randomSuffix

		// 插入数据库，创建基础账号
		r, err := g.DB().Model("client").Data(g.Map{
			"username":    username,
			"password":    passwordHash,
			"real_name":   realName,
			"phone":       "",
			"status":      1, // 默认正常
			"identifier":  "wxapp",
			"open_id":     openid,
			"session_key": sessionKey,
			"created_at":  gtime.Now().String(),
		}).Insert()
		if err != nil {
			return nil, gerror.New("创建微信客户失败")
		}

		// 获取插入ID
		lastId, err := r.LastInsertId()
		if err != nil {
			return nil, gerror.New("获取新客户ID失败")
		}

		client = struct {
			Id         int    `json:"id"`
			Username   string `json:"username"`
			Status     int    `json:"status"`
			OpenId     string `json:"open_id"`
			SessionKey string `json:"session_key"`
		}{
			Id:       int(lastId),
			Username: username,
			Status:   1,
			OpenId:   openid,
		}
	} else {
		// 更新session_key
		_, _ = g.DB().Model("client").
			Data(g.Map{
				"session_key": sessionKey,
			}).
			Where("id", client.Id).
			Update()
	}

	// 检查客户状态
	if client.Status == 0 {
		return nil, gerror.New("账号已被禁用")
	}

	// 更新登录信息
	_, _ = g.DB().Model("client").
		Data(g.Map{
			"last_login_time": gtime.Now().String(),
			"last_login_ip":   g.RequestFromCtx(ctx).GetClientIp(),
		}).
		Where("id", client.Id).
		Update()

	// 生成JWT令牌
	token, expire, err := createClientToken(ctx, client.Id, client.Username)
	if err != nil {
		return nil, gerror.New("生成令牌失败")
	}

	// 设置上下文信息
	ctx = auth.SetRequestFrom(ctx, "wxapp")
	ctx = auth.SetOpenID(ctx, openid)
	g.RequestFromCtx(ctx).SetCtx(ctx)

	// 组装返回数据
	res.ClientId = client.Id
	res.Token = token
	res.ExpireIn = expire

	return
}

// getWxOpenId 通过code获取微信OpenID和SessionKey
func getWxOpenId(ctx context.Context, code, appId, appSecret string) (openid, sessionKey string, err error) {
	if code == "" || appId == "" || appSecret == "" {
		return "", "", gerror.New("参数不完整")
	}

	// 微信登录凭证校验接口URL
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appId, appSecret, code,
	)

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 创建带有上下文的请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", gerror.New("创建请求失败: " + err.Error())
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", "", gerror.New("请求微信接口失败: " + err.Error())
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", gerror.New("读取微信响应失败: " + err.Error())
	}

	// 解析响应
	var result struct {
		OpenId     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionId    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", gerror.New("解析微信响应失败: " + err.Error())
	}

	// 检查错误
	if result.ErrCode != 0 {
		return "", "", gerror.New("微信登录失败: " + result.ErrMsg)
	}

	return result.OpenId, result.SessionKey, nil
}

// Info 获取客户信息
func (c *ControllerImpl) Info(ctx context.Context, req *ClientInfoReq) (res *ClientInfoRes, err error) {
	res = &ClientInfoRes{}

	// 从JWT中获取客户ID
	clientId := ctx.Value("client_id")
	if clientId == nil {
		return nil, gerror.New("未登录")
	}

	// 使用 DAO 层查询客户信息
	clientDao := dao.ClientDao()
	record, err := clientDao.GetById(ctx, clientId.(int))
	if err != nil {
		return nil, err
	}

	// 将查询结果赋值给返回对象
	err = record.Struct(res)
	if err != nil {
		return nil, gerror.New("转换客户数据失败")
	}

	return
}

// UpdateProfile 更新客户个人信息
func (c *ControllerImpl) UpdateProfile(ctx context.Context, req *ClientUpdateProfileReq) (res *ClientUpdateProfileRes, err error) {
	res = &ClientUpdateProfileRes{}

	// 从JWT中获取客户ID
	clientId := ctx.Value("client_id")
	if clientId == nil {
		return nil, gerror.New("未登录")
	}

	// 先获取客户当前信息
	var oldClientInfo struct {
		AvatarUrl string `json:"avatar_url"`
		RealName  string `json:"real_name"`
	}
	err = g.DB().Model("client").
		Fields("avatar_url, real_name").
		Where("id", clientId).
		Scan(&oldClientInfo)
	if err != nil {
		// 继续执行，不中断流程
	}
	oldAvatarUrl := oldClientInfo.AvatarUrl
	oldRealName := oldClientInfo.RealName

	// 更新数据库
	data := g.Map{}
	var newRealName string
	if req.RealName != "" {
		data["real_name"] = req.RealName
		newRealName = req.RealName
	}
	if req.Phone != "" {
		data["phone"] = req.Phone
	}

	// 保存新头像URL，用于后续比较
	var newAvatarUrl string

	// 处理头像上传
	if req.AvatarUrl != "" {
		// 判断是否为base64编码的图片数据
		if len(req.AvatarUrl) > 23 && req.AvatarUrl[:23] == "data:image/jpeg;base64," ||
			len(req.AvatarUrl) > 22 && req.AvatarUrl[:22] == "data:image/png;base64," ||
			len(req.AvatarUrl) > 22 && req.AvatarUrl[:22] == "data:image/gif;base64," {
			// 这是base64编码的图片，需要解码并上传到OSS
			avatarUrl, err := processBase64Avatar(ctx, req.AvatarUrl)
			if err != nil {
				return nil, gerror.New("头像上传失败")
			}
			data["avatar_url"] = avatarUrl
			newAvatarUrl = avatarUrl
		} else {
			// 直接使用URL
			data["avatar_url"] = req.AvatarUrl
			newAvatarUrl = req.AvatarUrl
		}
	}

	// 如果没有更新字段，则返回
	if len(data) == 0 {
		return
	}

	// 执行更新
	_, err = g.DB().Model("client").
		Data(data).
		Where("id", clientId).
		Update()
	if err != nil {
		return nil, gerror.New("更新个人信息失败")
	}

	// 如果更新了真实姓名，同步更新内容表中的作者字段
	if newRealName != "" && newRealName != oldRealName {
		// 执行更新内容表的作者字段
		_, err = g.DB().Model("content").
			Data(g.Map{"author": newRealName}).
			Where("client_id", clientId).
			Update()
		if err != nil {
			g.Log().Warning(ctx, "更新内容作者信息失败: "+err.Error())
			// 不中断主流程，继续执行
		}
	}

	// 如果更新了头像且旧头像存在，尝试删除旧头像
	if newAvatarUrl != "" && oldAvatarUrl != "" && newAvatarUrl != oldAvatarUrl {
		// 异步删除旧头像，不阻塞主流程
		go func(avatarUrl string) {
			defer func() {
				recover()
			}()

			// 创建新的上下文，避免使用已完成的请求上下文
			deleteCtx := context.Background()

			// 获取OSS客户端
			ossClient, err := storage.GetOSSClient(deleteCtx)
			if err != nil {
				return
			}

			// 分析URL获取正确的对象键
			var objectKey string

			// 检查URL格式并提取对象键
			if strings.HasPrefix(avatarUrl, "http://") || strings.HasPrefix(avatarUrl, "https://") {
				// 从完整URL中提取对象键
				uri, err := url.Parse(avatarUrl)
				if err != nil {
					return
				}

				// 获取路径部分作为初始对象键
				objectKey = strings.TrimPrefix(uri.Path, "/")

				// 如果使用了公共访问URL，去除目录前缀
				// 避免出现重复的目录前缀，例如 test/test/xx
				if ossClient.PublicURL != "" {
					pubUri, _ := url.Parse(ossClient.PublicURL)
					if pubUri != nil && uri.Host == pubUri.Host {
						// 这是使用公共URL访问的情况，检查是否存在目录前缀
						if ossClient.Config.Directory != "" && strings.HasPrefix(objectKey, ossClient.Config.Directory) {
							// 如果对象键已经包含了目录前缀，那么不要再次添加前缀
							_ = ossClient.Bucket.DeleteObject(objectKey)
							return
						}
					}
				}
			} else {
				// 如果不是完整URL，直接使用作为对象键
				objectKey = strings.TrimPrefix(avatarUrl, "/")
			}

			// 尝试直接使用对象键删除
			_ = ossClient.Bucket.DeleteObject(objectKey)

			// 尝试使用完整对象键
			fullObjectKey := ossClient.GetFullObjectKey(objectKey)

			_ = ossClient.Bucket.DeleteObject(fullObjectKey)

			// 尝试其他可能的对象键
			if ossClient.Config.Directory != "" {
				// 从URL中提取不包含目录前缀的部分
				parts := strings.Split(objectKey, "/")
				if len(parts) > 1 {
					alternativeKey := strings.Join(parts[1:], "/")

					_ = ossClient.Bucket.DeleteObject(alternativeKey)
				}
			}

			// 如果对象键包含URL参数，去除参数部分
			if strings.Contains(objectKey, "?") {
				cleanKey := strings.Split(objectKey, "?")[0]

				_ = ossClient.Bucket.DeleteObject(cleanKey)
			}
		}(oldAvatarUrl)
	}

	return
}

// createClientToken 创建客户JWT令牌
func createClientToken(ctx context.Context, id int, username string) (token string, expire int, err error) {
	return auth.CreateClientToken(ctx, id, username)
}

// processBase64Avatar 处理base64编码的头像图片
func processBase64Avatar(ctx context.Context, base64Data string) (string, error) {
	// 获取图片格式和数据
	var contentType, base64Image string
	if len(base64Data) > 23 && base64Data[:23] == "data:image/jpeg;base64," {
		contentType = "image/jpeg"
		base64Image = base64Data[23:]
	} else if len(base64Data) > 22 && base64Data[:22] == "data:image/png;base64," {
		contentType = "image/png"
		base64Image = base64Data[22:]
	} else if len(base64Data) > 22 && base64Data[:22] == "data:image/gif;base64," {
		contentType = "image/gif"
		base64Image = base64Data[22:]
	} else {
		return "", gerror.New("不支持的图片格式")
	}

	// 解码base64
	imgData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", gerror.New("解码图片数据失败: " + err.Error())
	}

	// 获取OSS客户端
	ossClient, err := storage.GetOSSClient(ctx)
	if err != nil {
		return "", gerror.New("获取OSS客户端失败: " + err.Error())
	}

	// 确定文件扩展名
	var ext string
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	default:
		ext = ".jpg"
	}

	// 使用通用函数生成存储路径
	objectKey := ossClient.GenerateStoragePath("wx_avatar", ext)

	// 获取完整的对象键（含目录前缀）
	fullObjectKey := ossClient.GetFullObjectKey(objectKey)

	// 上传头像到OSS
	err = ossClient.Bucket.PutObject(fullObjectKey, bytes.NewReader(imgData))
	if err != nil {
		return "", gerror.New("上传头像到OSS失败: " + err.Error())
	}

	// 判断是否为公开访问
	if ossClient.Config.PublicAccess {
		// 公开访问，直接返回URL
		return ossClient.PublicURL + "/" + fullObjectKey, nil
	} else {
		// 私有访问，返回临时URL（有效期30天）
		signedURL, err := ossClient.GetObjectURL(fullObjectKey, true, 30*24*3600)
		if err != nil {
			return "", gerror.New("生成头像访问URL失败: " + err.Error())
		}
		return signedURL, nil
	}
}

// DurationList 获取客户时长列表
func (c *ControllerImpl) DurationList(ctx context.Context, req *ClientDurationListReq) (res *ClientDurationListRes, err error) {
	res = &ClientDurationListRes{
		List:  make([]ClientDurationItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 构建查询参数
	m := g.Map{}
	if req.ClientId != 0 {
		m["client_id"] = req.ClientId
	}

	// 使用 DAO 层获取客户时长列表
	clientDurationDao := dao.ClientDurationDao()
	list, total, err := clientDurationDao.GetWithPage(ctx, req.Page, req.PageSize, m)
	if err != nil {
		return nil, err
	}
	res.Total = total

	// 处理结果
	var items []ClientDurationItem
	if err = list.Structs(&items); err != nil {
		return nil, gerror.New("转换客户时长列表数据失败")
	}

	res.List = items
	return
}

// DurationDetail 获取客户时长详情
func (c *ControllerImpl) DurationDetail(ctx context.Context, req *ClientDurationDetailReq) (res *ClientDurationDetailRes, err error) {
	res = &ClientDurationDetailRes{}

	// 使用 DAO 层获取客户时长详情
	clientDurationDao := dao.ClientDurationDao()
	record, err := clientDurationDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 将查询结果赋值给返回对象
	err = record.Struct(res)
	if err != nil {
		return nil, gerror.New("转换客户时长数据失败")
	}

	return
}

// DurationCreate 创建客户时长
func (c *ControllerImpl) DurationCreate(ctx context.Context, req *ClientDurationCreateReq) (res *ClientDurationCreateRes, err error) {
	res = &ClientDurationCreateRes{}

	// 使用 DAO 层检查客户是否存在
	clientDao := dao.ClientDao()
	count, err := clientDao.Model(ctx).Where("id", req.ClientId).Count()
	if err != nil {
		return nil, gerror.New("检查客户失败")
	}
	if count == 0 {
		return nil, gerror.New("客户不存在")
	}

	// 构建插入数据
	data := g.Map{
		"client_id":          req.ClientId,
		"client_name":        req.ClientName,
		"remaining_duration": req.RemainingDuration,
		"total_duration":     req.TotalDuration,
		"used_duration":      req.UsedDuration,
		"created_at":         gtime.Now().String(),
	}

	// 使用 DAO 层插入数据
	clientDurationDao := dao.ClientDurationDao()
	lastId, err := clientDurationDao.Insert(ctx, data)
	if err != nil {
		return nil, gerror.New("创建客户时长记录失败")
	}

	res.Id = int(lastId)
	return
}

// DurationUpdate 更新客户时长
func (c *ControllerImpl) DurationUpdate(ctx context.Context, req *ClientDurationUpdateReq) (res *ClientDurationUpdateRes, err error) {
	res = &ClientDurationUpdateRes{}

	// 使用 DAO 层检查记录是否存在
	clientDurationDao := dao.ClientDurationDao()
	_, err = clientDurationDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("客户时长记录不存在")
	}

	// 使用 DAO 层检查客户是否存在
	clientDao := dao.ClientDao()
	count, err := clientDao.Model(ctx).Where("id", req.ClientId).Count()
	if err != nil {
		return nil, gerror.New("检查客户失败")
	}
	if count == 0 {
		return nil, gerror.New("客户不存在")
	}

	// 构建更新数据
	data := g.Map{
		"client_id":          req.ClientId,
		"client_name":        req.ClientName,
		"remaining_duration": req.RemainingDuration,
		"total_duration":     req.TotalDuration,
		"used_duration":      req.UsedDuration,
		"updated_at":         gtime.Now().String(),
	}

	// 使用 DAO 层更新数据
	err = clientDurationDao.UpdateById(ctx, req.Id, data)
	if err != nil {
		return nil, gerror.New("更新客户时长记录失败")
	}

	return
}

// DurationDelete 删除客户时长
func (c *ControllerImpl) DurationDelete(ctx context.Context, req *ClientDurationDeleteReq) (res *ClientDurationDeleteRes, err error) {
	res = &ClientDurationDeleteRes{}

	// 使用 DAO 层检查记录是否存在
	clientDurationDao := dao.ClientDurationDao()
	_, err = clientDurationDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("客户时长记录不存在")
	}

	// 使用 DAO 层删除记录
	err = clientDurationDao.DeleteById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("删除客户时长记录失败")
	}

	return
}

// WxGetRemainingDuration 获取客户端用户剩余时长
func (c *ControllerImpl) WxGetRemainingDuration(ctx context.Context, req *WxClientRemainingDurationReq) (res *WxClientRemainingDurationRes, err error) {
	res = &WxClientRemainingDurationRes{}

	// 从上下文获取客户ID
	clientId, err := auth.GetClientInfo(ctx)
	if err != nil || clientId == 0 {
		return nil, gerror.New("未登录或登录已过期")
	}

	// 使用 DAO 层获取客户时长信息
	clientDurationDao := dao.ClientDurationDao()
	record, err := clientDurationDao.GetByClientId(ctx, clientId)
	if err != nil {
		// 如果是记录不存在的错误，返回空时长信息而不是报错
		if gerror.Cause(err).Error() == "客户时长记录不存在" {
			res.RemainingDuration = "0分钟"
			return res, nil
		}
		return nil, gerror.New("获取客户时长信息失败")
	}

	// 从记录中获取剩余时长信息
	res.RemainingDuration = record["remaining_duration"].String()

	return
}
