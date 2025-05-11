package consts

import (
	"github.com/gogf/gf/v2/errors/gcode"
)

// 用户相关错误码
var (
	// CodeUserNotExists 用户不存在
	CodeUserNotExists = gcode.New(10001, "用户不存在", nil)

	// CodeUserExists 用户已存在
	CodeUserExists = gcode.New(10002, "用户已存在", nil)

	// CodeUserPasswordError 密码错误
	CodeUserPasswordError = gcode.New(10003, "密码错误", nil)

	// CodeUserForbidden 用户被禁用
	CodeUserForbidden = gcode.New(10004, "用户已被禁用", nil)

	// CodeUserCreateFailed 创建用户失败
	CodeUserCreateFailed = gcode.New(10005, "创建用户失败", nil)

	// CodeUserUpdateFailed 更新用户失败
	CodeUserUpdateFailed = gcode.New(10006, "更新用户失败", nil)

	// CodeUserDeleteFailed 删除用户失败
	CodeUserDeleteFailed = gcode.New(10007, "删除用户失败", nil)

	// CodeUserListFailed 获取用户列表失败
	CodeUserListFailed = gcode.New(10008, "获取用户列表失败", nil)

	// 验证码相关错误码
	// CodeCaptchaGenerateFailed 验证码生成失败
	CodeCaptchaGenerateFailed = gcode.New(10101, "验证码生成失败", nil)

	// CodeCaptchaInvalid 验证码无效或已过期
	CodeCaptchaInvalid = gcode.New(10102, "验证码无效或已过期", nil)

	// CodeCaptchaVerifyFailed 验证码验证失败
	CodeCaptchaVerifyFailed = gcode.New(10103, "验证码验证失败", nil)

	// 登录相关错误码
	// CodeLoginFailed 登录失败
	CodeLoginFailed = gcode.New(10201, "登录失败", nil)

	// CodeTokenGenerateFailed 令牌生成失败
	CodeTokenGenerateFailed = gcode.New(10202, "令牌生成失败", nil)

	// CodeTokenInvalid 令牌无效或已过期
	CodeTokenInvalid = gcode.New(10203, "令牌无效或已过期", nil)

	// CodeUnauthorized 未授权访问
	CodeUnauthorized = gcode.New(10204, "未授权访问", nil)

	// CodeForbidden 禁止访问，没有权限
	CodeForbidden = gcode.New(10205, "禁止访问，没有权限", nil)

	// CodeNonAdminToken 非管理员令牌
	CodeNonAdminToken = gcode.New(10206, "非管理员令牌", nil)

	// CodeNonClientToken 非客户令牌
	CodeNonClientToken = gcode.New(10207, "非客户令牌", nil)

	// CodeLogoutFailed 退出登录失败
	CodeLogoutFailed = gcode.New(10208, "退出登录失败", nil)

	// 操作日志相关错误码
	// CodeOperationLogListFailed 获取操作日志列表失败
	CodeOperationLogListFailed = gcode.New(10301, "获取操作日志列表失败", nil)

	// CodeOperationLogDeleteFailed 删除操作日志失败
	CodeOperationLogDeleteFailed = gcode.New(10302, "删除操作日志失败", nil)

	// CodeOperationLogExportFailed 导出操作日志失败
	CodeOperationLogExportFailed = gcode.New(10303, "导出操作日志失败", nil)

	// 客户相关错误码
	// CodeClientNotExists 客户不存在
	CodeClientNotExists = gcode.New(10401, "客户不存在", nil)

	// CodeClientExists 客户已存在
	CodeClientExists = gcode.New(10402, "客户已存在", nil)

	// CodeClientCreateFailed 创建客户失败
	CodeClientCreateFailed = gcode.New(10403, "创建客户失败", nil)

	// CodeClientUpdateFailed 更新客户失败
	CodeClientUpdateFailed = gcode.New(10404, "更新客户失败", nil)

	// CodeClientDeleteFailed 删除客户失败
	CodeClientDeleteFailed = gcode.New(10405, "删除客户失败", nil)

	// CodeClientListFailed 获取客户列表失败
	CodeClientListFailed = gcode.New(10406, "获取客户列表失败", nil)

	// 地区管理相关错误码
	// CodeRegionListFailed 获取地区列表失败
	CodeRegionListFailed = gcode.New(10501, "获取地区列表失败", nil)

	// CodeRegionDetailFailed 获取地区详情失败
	CodeRegionDetailFailed = gcode.New(10502, "获取地区详情失败", nil)

	// CodeRegionCreateFailed 创建地区失败
	CodeRegionCreateFailed = gcode.New(10503, "创建地区失败", nil)

	// CodeRegionUpdateFailed 更新地区失败
	CodeRegionUpdateFailed = gcode.New(10504, "更新地区失败", nil)

	// CodeRegionDeleteFailed 删除地区失败
	CodeRegionDeleteFailed = gcode.New(10505, "删除地区失败", nil)

	// CodeRegionNotExists 地区不存在
	CodeRegionNotExists = gcode.New(10506, "地区不存在", nil)

	// CodeRegionHasChildren 地区下有子地区不能删除
	CodeRegionHasChildren = gcode.New(10507, "该地区下有子地区，不能删除", nil)

	// CodeRegionCascaderFailed 获取级联地区数据失败
	CodeRegionCascaderFailed = gcode.New(10508, "获取级联地区数据失败", nil)

	// CodeParentRegionNotExists 父级地区不存在
	CodeParentRegionNotExists = gcode.New(10509, "父级地区不存在", nil)

	// CodeInvalidRegionLevel 地区级别不匹配
	CodeInvalidRegionLevel = gcode.New(10510, "地区级别不匹配", nil)

	// CodeRegionParentIsItself 不能将地区自身设为父级
	CodeRegionParentIsItself = gcode.New(10511, "不能将地区自身设为父级", nil)

	// CodeRegionNameExists 地区名称已存在
	CodeRegionNameExists = gcode.New(10512, "同一父级下已存在同名地区", nil)
)
