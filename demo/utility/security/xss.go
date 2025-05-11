package security

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

var (
	// 危险URL正则表达式
	dangerousUrlRegex = regexp.MustCompile(`(?i)(javascript|data|vbscript|file):`)
	// XSS攻击特征正则表达式
	xssRegex = regexp.MustCompile(`(?i)(<\s*script|on\w+\s*=|\w+script:|<\s*iframe|<\s*object|<\s*embed|alert\s*\()`)
)

// 内容过滤类型常量
const (
	// 严格模式 - 不允许任何HTML标签
	SanitizeTypeStrict = "strict"
	// 评论模式 - 允许有限的标签，如b, i, code等
	SanitizeTypeComment = "comment"
	// 内容模式 - 允许更多格式化标签，适用于富文本内容
	SanitizeTypeContent = "content"
	// 基本UGC模式 - 用户生成内容的基本过滤
	SanitizeTypeUGC = "ugc"
)

// 全局策略实例缓存
var (
	strictPolicy  *bluemonday.Policy
	commentPolicy *bluemonday.Policy
	contentPolicy *bluemonday.Policy
	ugcPolicy     *bluemonday.Policy
)

// 初始化全局策略实例
func init() {
	strictPolicy = bluemonday.StrictPolicy()
	commentPolicy = newCommentPolicy()
	contentPolicy = newContentPolicy()
	ugcPolicy = newUGCPolicy()
}

// SanitizeHTML 过滤HTML内容，防止XSS攻击
func SanitizeHTML(content string, policyType string) string {
	if content == "" {
		return ""
	}

	switch policyType {
	case SanitizeTypeStrict:
		return strictPolicy.Sanitize(content)
	case SanitizeTypeComment:
		return commentPolicy.Sanitize(content)
	case SanitizeTypeContent:
		return contentPolicy.Sanitize(content)
	case SanitizeTypeUGC:
		return ugcPolicy.Sanitize(content)
	default:
		return strictPolicy.Sanitize(content)
	}
}

// StripAllTags 去除所有HTML标签，只保留纯文本
func StripAllTags(content string) string {
	return strictPolicy.Sanitize(content)
}

// SanitizeComment 过滤评论内容
// 适用于用户评论，只允许基本格式化
func SanitizeComment(content string) string {
	return commentPolicy.Sanitize(content)
}

// SanitizeRichText 过滤富文本内容
// 适用于CMS内容、文章等富文本
func SanitizeRichText(content string) string {
	return contentPolicy.Sanitize(content)
}

// SanitizeUserContent 过滤用户生成内容
// 适用于用户博客、帖子等内容
func SanitizeUserContent(content string) string {
	return ugcPolicy.Sanitize(content)
}

// ContainsSuspiciousXSS 检测内容是否包含可疑的XSS攻击特征
func ContainsSuspiciousXSS(content string) bool {
	return xssRegex.MatchString(content)
}

// newCommentPolicy 创建评论安全策略
// 只允许基本的文本格式化和链接
func newCommentPolicy() *bluemonday.Policy {
	p := bluemonday.NewPolicy()

	// 允许的HTML元素
	p.AllowElements("b", "i", "strong", "em", "code", "pre", "a")

	// 为<a>标签允许href属性，但需要过滤危险URL
	p.AllowAttrs("href").OnElements("a")
	p.RequireNoReferrerOnLinks(true)
	p.RequireParseableURLs(true)
	p.AllowURLSchemes("http", "https")

	// 允许换行
	p.AllowElements("br")

	return p
}

// newContentPolicy 创建内容安全策略
// 适用于富文本编辑器内容，允许更多格式化元素
func newContentPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()

	// 增加额外允许的元素
	p.AllowElements("div", "span", "h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowElements("blockquote", "hr")
	p.AllowElements("table", "thead", "tbody", "tr", "th", "td")

	// 允许一些常见样式
	p.AllowAttrs("class").OnElements("div", "span", "p", "h1", "h2", "h3", "h4", "h5", "h6", "table", "tr", "td")
	p.AllowAttrs("style").OnElements("div", "span", "p", "h1", "h2", "h3", "h4", "h5", "h6", "table", "tr", "td")

	// 限制style属性中的内容
	p.AllowStyles("color", "background-color", "text-align", "font-size", "font-weight", "font-style", "text-decoration").OnElements("div", "span", "p", "h1", "h2", "h3", "h4", "h5", "h6")

	return p
}

// newUGCPolicy 创建用户生成内容安全策略
// 基于bluemonday的UGCPolicy，适用于博客、文章等
func newUGCPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()

	// 允许一些额外的属性
	p.AllowAttrs("class").Globally()
	p.AllowAttrs("target").OnElements("a")

	return p
}

// ValidateURL 验证URL是否安全
// 返回true表示URL安全，false表示不安全
func ValidateURL(url string) bool {
	if url == "" {
		return false
	}

	// 检查是否包含危险协议
	return !dangerousUrlRegex.MatchString(url)
}

// TruncateText 安全地截断文本到指定长度
// 会先移除HTML标签，然后截断
func TruncateText(content string, length int) string {
	if content == "" {
		return ""
	}

	// 先移除HTML标签
	text := StripAllTags(content)

	// 如果文本长度小于指定长度，直接返回
	runeText := []rune(text)
	if len(runeText) <= length {
		return text
	}

	// 截断文本并添加省略号
	return string(runeText[:length]) + "..."
}
