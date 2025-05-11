package user

import (
	"context"
	v1 "demo/api/user/v1"
	"demo/internal/consts"
	"demo/internal/dao"
	"demo/internal/model/do"
	"demo/internal/model/entity"
	"demo/internal/service"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

// userImpl 用户服务实现
type userImpl struct {
	userDao *dao.UserDao
}

// New 创建一个用户服务实例
func New() service.UserService {
	return &userImpl{
		userDao: new(dao.UserDao),
	}
}

// List 获取用户列表
func (s *userImpl) List(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error) {
	res = &v1.UserListRes{
		List:  make([]v1.UserListItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 构建查询条件
	filter := make(map[string]interface{})
	if req.Username != "" {
		filter["username"] = req.Username
	}
	if req.Nickname != "" {
		filter["nickname"] = req.Nickname
	}
	if req.Status != 0 {
		filter["status"] = req.Status
	}

	// 执行查询
	list, total, err := s.userDao.FindList(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeUserListFailed, err.Error())
	}

	// 转换为API响应格式
	res.Total = total
	if len(list) > 0 {
		for _, item := range list {
			// 状态转换为文本
			statusText := "正常"
			if gconv.Int(item.Status) == 0 {
				statusText = "禁用"
			}

			// 添加到结果列表
			res.List = append(res.List, v1.UserListItem{
				Id:            gconv.Int(item.Id),
				Username:      gconv.String(item.Username),
				Nickname:      gconv.String(item.Nickname),
				Status:        gconv.Int(item.Status),
				StatusText:    statusText,
				LastLoginIp:   gconv.String(item.LastLoginIp),
				LastLoginTime: item.LastLoginTime.String(),
			})
		}
	}

	return res, nil
}

// Create 创建用户
func (s *userImpl) Create(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error) {
	res = &v1.UserCreateRes{}

	// 检查用户名是否存在
	exists, err := s.userDao.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeUserCreateFailed, err.Error())
	}
	if exists != nil {
		return nil, gerror.NewCode(consts.CodeUserExists)
	}

	// 密码加密
	password, err := gmd5.EncryptString(req.Password)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeUserCreateFailed, "密码加密失败")
	}

	// 准备数据
	data := &do.UserDO{
		Username: req.Username,
		Password: password,
		Nickname: req.Nickname,
		Role:     "admin", // 设置所有用户为管理员角色
		Status:   1,       // 默认正常状态
	}

	// 插入数据
	insertId, err := s.userDao.Insert(ctx, data)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeUserCreateFailed, err.Error())
	}

	res.Id = int(insertId)
	return res, nil
}

// Update 更新用户
func (s *userImpl) Update(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error) {
	res = &v1.UserUpdateRes{}

	// 检查用户是否存在
	user, err := s.userDao.FindById(ctx, req.Id)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeUserUpdateFailed, err.Error())
	}
	if user == nil {
		return nil, gerror.NewCode(consts.CodeUserNotExists)
	}

	// 如果修改了用户名，则检查新用户名是否已存在
	if gconv.String(user.Username) != req.Username {
		exists, err := s.userDao.FindByUsername(ctx, req.Username)
		if err != nil {
			return nil, gerror.NewCode(consts.CodeUserUpdateFailed, err.Error())
		}
		if exists != nil && gconv.Int(exists.Id) != req.Id {
			return nil, gerror.NewCode(consts.CodeUserExists)
		}
	}

	// 准备更新数据
	data := &do.UserDO{
		Username: req.Username,
		Nickname: req.Nickname,
		Role:     "admin", // 确保更新时保持管理员角色
		Status:   req.Status,
	}

	// 执行更新
	_, err = s.userDao.Update(ctx, data, req.Id)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeUserUpdateFailed, err.Error())
	}

	return res, nil
}

// Delete 删除用户
func (s *userImpl) Delete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error) {
	res = &v1.UserDeleteRes{}

	// 检查用户是否存在
	user, err := s.userDao.FindById(ctx, req.Id)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeUserDeleteFailed, err.Error())
	}
	if user == nil {
		return nil, gerror.NewCode(consts.CodeUserNotExists)
	}

	// 执行删除
	_, err = s.userDao.Delete(ctx, req.Id)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeUserDeleteFailed, err.Error())
	}

	return res, nil
}

// GetUserById 通过ID获取用户
func (s *userImpl) GetUserById(ctx context.Context, id int) (user *entity.User, err error) {
	// 查询用户数据
	data, err := s.userDao.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	// 数据转换
	user = &entity.User{}
	err = gconv.Struct(data, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername 通过用户名获取用户
func (s *userImpl) GetUserByUsername(ctx context.Context, username string) (user *entity.User, err error) {
	// 查询用户数据
	data, err := s.userDao.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	// 数据转换
	user = &entity.User{}
	err = gconv.Struct(data, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
