package middleware

import (
	"demo/utility/security"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

// XssMiddleware XSS防护中间件
// 对请求参数进行XSS防护过滤
func XssMiddleware(r *ghttp.Request) {
	// 检查路径是否在排除列表中
	excludePaths := g.Cfg().MustGet(r.GetCtx(), "security.xss.excludePaths").Strings()
	path := r.URL.Path
	for _, excludePath := range excludePaths {
		if strings.HasPrefix(path, excludePath) {
			r.Middleware.Next()
			return
		}
	}

	// 处理GET请求参数
	if len(r.GetMap()) > 0 {
		getMap := r.GetMap()
		for k, v := range getMap {
			if vStr, ok := v.(string); ok {
				getMap[k] = sanitizeByParamName(vStr, k)
			}
		}
		// 重新设置GET参数
		// 注意：这里不修改原有的URL查询参数，仅修改解析后的Map
	}

	// 处理POST请求参数
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		// 处理表单数据
		if r.Form != nil {
			sanitizeUrlValues(r.Form)
		}

		// 处理POST表单数据
		if r.PostForm != nil {
			sanitizeUrlValues(r.PostForm)
		}

		// 处理JSON数据
		if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			jsonMap := make(map[string]interface{})
			if err := r.Parse(&jsonMap); err == nil && len(jsonMap) > 0 {
				sanitizeJsonData(jsonMap)
				// 将处理后的数据写回请求体
				// GoFrame v2不支持直接修改已解析的请求体
				// 这里只对数据进行了清洗，依赖后续的业务逻辑获取过滤后的值
				g.Log().Debug(r.GetCtx(), "XSS中间件已过滤JSON数据")
			}
		}
	}

	r.Middleware.Next()
}

// sanitizeUrlValues 过滤URL参数中的XSS内容
func sanitizeUrlValues(values url.Values) {
	for key, vals := range values {
		for i, val := range vals {
			values[key][i] = sanitizeByParamName(val, key)
		}
	}
}

// sanitizeJsonData 递归处理JSON数据，过滤XSS内容
func sanitizeJsonData(data map[string]interface{}) {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			data[key] = sanitizeByParamName(v, key)
		case map[string]interface{}:
			sanitizeJsonData(v)
		case []interface{}:
			for i, item := range v {
				if itemStr, ok := item.(string); ok {
					v[i] = sanitizeByParamName(itemStr, key)
				} else if itemMap, ok := item.(map[string]interface{}); ok {
					sanitizeJsonData(itemMap)
				}
			}
		}
	}
}

// sanitizeByParamName 根据参数名选择合适的过滤策略
func sanitizeByParamName(content string, paramName string) string {
	if content == "" {
		return ""
	}

	paramName = strings.ToLower(paramName)

	// 选择过滤策略
	policyType := security.SanitizeTypeStrict

	// 根据参数名特征选择策略
	if strings.Contains(paramName, "content") ||
		strings.Contains(paramName, "detail") ||
		strings.Contains(paramName, "description") ||
		strings.Contains(paramName, "html") {
		policyType = security.SanitizeTypeContent
	} else if strings.Contains(paramName, "comment") ||
		strings.Contains(paramName, "reply") {
		policyType = security.SanitizeTypeComment
	} else if strings.Contains(paramName, "blog") ||
		strings.Contains(paramName, "article") ||
		strings.Contains(paramName, "post") {
		policyType = security.SanitizeTypeUGC
	}

	// 进行内容过滤
	return security.SanitizeHTML(content, policyType)
}

// XSSFilterHandler API请求XSS过滤处理器
// 用于控制器中手动调用，对请求参数进行XSS过滤
// 用法示例:
//
//	func (c *Controller) Create(ctx context.Context, req *api.CreateReq) (res *api.CreateRes, err error) {
//	   // 对请求参数进行XSS过滤
//	   middleware.XSSFilterHandler(req, map[string]string{
//	       "Content": security.SanitizeTypeContent,
//	       "Comment": security.SanitizeTypeComment,
//	   })
//	   // 业务处理...
//	}
func XSSFilterHandler(req interface{}, policyMap map[string]string) {
	reqMap := gconv.Map(req)
	if reqMap == nil {
		return
	}

	// 处理请求数据
	for key, value := range reqMap {
		if valueStr, ok := value.(string); ok {
			policy := security.SanitizeTypeStrict
			if p, ok := policyMap[key]; ok {
				policy = p
			}
			reqMap[key] = security.SanitizeHTML(valueStr, policy)
		} else if valueMap, ok := value.(map[string]interface{}); ok {
			sanitizeJsonData(valueMap)
		} else if valueArray, ok := value.([]interface{}); ok {
			for i, item := range valueArray {
				if itemStr, ok := item.(string); ok {
					policy := security.SanitizeTypeStrict
					if p, ok := policyMap[key]; ok {
						policy = p
					}
					valueArray[i] = security.SanitizeHTML(itemStr, policy)
				} else if itemMap, ok := item.(map[string]interface{}); ok {
					sanitizeJsonData(itemMap)
				}
			}
		}
	}

	// 将处理后的数据转换回原结构
	gconv.Struct(reqMap, req)
}
