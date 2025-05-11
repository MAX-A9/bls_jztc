package middleware

import (
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

// SwaggerMiddleware 自定义Swagger中间件，为OpenAPI响应添加安全定义
func SwaggerMiddleware(r *ghttp.Request) {
	// 对.map文件的请求直接返回空内容，避免大量404日志
	if strings.HasSuffix(r.URL.Path, ".map") {
		r.Response.WriteStatus(404)
		r.Exit()
		return
	}

	// 仅处理/api.json请求
	if r.URL.Path == "/api.json" {
		// 让请求继续处理
		r.Middleware.Next()

		// 如果响应内容为JSON，添加安全定义
		contentType := r.Response.Header().Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			// 获取原始响应内容
			responseContent := r.Response.BufferString()

			// 解析为JSON
			var openApiSpec map[string]interface{}
			if err := json.Unmarshal([]byte(responseContent), &openApiSpec); err == nil {
				// 添加安全定义
				if _, ok := openApiSpec["components"]; !ok {
					openApiSpec["components"] = make(map[string]interface{})
				}

				components := openApiSpec["components"].(map[string]interface{})

				// 添加安全方案
				securitySchemes := map[string]interface{}{
					"Bearer": map[string]interface{}{
						"type":         "http",
						"scheme":       "bearer",
						"bearerFormat": "JWT",
						"description":  "在下方输入您的JWT令牌(不需要Bearer前缀)",
					},
				}

				components["securitySchemes"] = securitySchemes

				// 将更新后的spec转回为JSON
				if modifiedJson, err := json.Marshal(openApiSpec); err == nil {
					// 替换响应内容
					r.Response.ClearBuffer()
					r.Response.Write(modifiedJson)
				}
			}
		}
	} else {
		// 其他请求正常处理
		r.Middleware.Next()
	}
}
