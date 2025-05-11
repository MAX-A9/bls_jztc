package middleware

import (
	"context"
	"demo/internal/service"
	"demo/utility/auth"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

// OperationLogContextKey 操作日志中间件的上下文键类型
type OperationLogContextKey string

// 定义上下文键常量
const (
	RequestObjectKey OperationLogContextKey = "RequestObject"
)

// OperationLog 操作日志中间件
func OperationLog(r *ghttp.Request) {
	// 跳过不需要记录日志的请求
	if shouldSkipLog(r) {
		r.Middleware.Next()
		return
	}

	// 先执行后续中间件和处理程序，确保Auth中间件有机会先设置用户信息
	r.Middleware.Next()

	// 获取用户名和请求上下文
	var (
		username string
		userId   int
		ctx      = r.GetCtx()
	)

	// 从auth包获取用户信息
	userId, username, _, err := auth.GetUserInfo(ctx)
	if err != nil || username == "" {
		username = "未登录用户"
	}

	// 优化操作类型显示，全部使用中文
	var action string
	pathParts := strings.Split(r.URL.Path, "/")

	// 根据路径和HTTP方法确定操作类型，使用中文模块名
	if len(pathParts) >= 3 {
		module := pathParts[1]
		operation := pathParts[2]

		// 将模块名转换为中文
		moduleCN := getModuleNameCN(module)

		// 特定操作类型的优化处理
		if operation == "list" {
			action = "查询" + moduleCN
		} else if operation == "create" {
			action = "创建" + moduleCN
		} else if operation == "update" {
			action = "更新" + moduleCN
		} else if operation == "delete" {
			action = "删除" + moduleCN
		} else if operation == "export" {
			action = "导出" + moduleCN
		} else {
			// 兜底处理
			action = getActionFromMethod(r.Method) + moduleCN
		}
	} else if len(pathParts) == 2 && pathParts[1] != "" {
		// 处理只有一级路径的情况，如 /login
		module := pathParts[1]
		if module == "login" {
			action = "用户登录"
		} else {
			// 将模块名转换为中文
			moduleCN := getModuleNameCN(module)
			action = getActionFromMethod(r.Method) + moduleCN
		}
	} else {
		// 处理根路径
		action = getActionFromMethod(r.Method) + "根路径"
	}

	// 获取响应状态
	status := r.Response.Status
	operationStatus := 1 // 默认成功
	if status < 200 || status >= 400 {
		operationStatus = 0 // 失败
	}

	// 不记录查询操作，包括查询日志
	if r.Method == "GET" {
		return
	}

	// 创建一个包含请求对象的上下文，确保IP地址能够被正确提取
	reqCtx := context.WithValue(ctx, RequestObjectKey, r)

	// 异步记录操作日志
	go service.Log().Record(reqCtx, userId, username, "system", action, operationStatus, "")
}

// 将模块名转换为中文
func getModuleNameCN(module string) string {
	switch module {
	case "user":
		return "用户"
	case "log":
		return "日志"
	case "role":
		return "角色"
	case "menu":
		return "菜单"
	case "dept":
		return "部门"
	case "post":
		return "岗位"
	case "dict":
		return "字典"
	case "config":
		return "配置"
	case "notice":
		return "通知"
	case "file":
		return "文件"
	case "job":
		return "任务"
	case "monitor":
		return "监控"
	default:
		return module // 对于未知模块，保留英文名
	}
}

// shouldSkipLog 判断是否应该跳过日志记录
func shouldSkipLog(r *ghttp.Request) bool {
	// 跳过静态资源请求
	if strings.HasPrefix(r.URL.Path, "/static") ||
		strings.HasPrefix(r.URL.Path, "/swagger") ||
		strings.HasPrefix(r.URL.Path, "/export") ||
		strings.HasSuffix(r.URL.Path, ".map") ||
		strings.HasSuffix(r.URL.Path, ".ico") {
		return true
	}

	// 跳过健康检查等不需要记录的接口
	if r.URL.Path == "/health" || r.URL.Path == "/metrics" {
		return true
	}

	// 跳过API文档相关
	if r.URL.Path == "/api.json" {
		return true
	}

	// 跳过客户端相关接口（所有/wx前缀的接口）
	if strings.HasPrefix(r.URL.Path, "/wx") {
		return true
	}

	// 所有GET请求都跳过，包括日志查询
	if r.Method == "GET" {
		return true
	}

	return false
}

// getActionFromMethod 根据HTTP方法获取操作类型
func getActionFromMethod(method string) string {
	switch method {
	case "GET":
		return "查询"
	case "POST":
		return "创建"
	case "PUT":
		return "更新"
	case "DELETE":
		return "删除"
	default:
		return "其他"
	}
}
