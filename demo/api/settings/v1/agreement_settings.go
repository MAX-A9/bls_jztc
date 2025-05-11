package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// AgreementSettingsReq 获取协议设置请求
type AgreementSettingsReq struct {
	g.Meta `path:"/agreement/settings" method:"get" tags:"协议设置" summary:"获取协议设置" security:"Bearer" description:"获取隐私政策和用户协议设置，需要管理员权限"`
}

// AgreementSettingsRes 获取协议设置响应
type AgreementSettingsRes struct {
	g.Meta        `mime:"application/json" example:"json"`
	PrivacyPolicy string `json:"privacyPolicy" dc:"隐私政策内容"`
	UserAgreement string `json:"userAgreement" dc:"用户协议内容"`
}

// AgreementSettingsSaveReq 保存协议设置请求
type AgreementSettingsSaveReq struct {
	g.Meta        `path:"/agreement/settings/save" method:"post" tags:"协议设置" summary:"保存协议设置" security:"Bearer" description:"保存隐私政策和用户协议设置，需要管理员权限"`
	PrivacyPolicy string `json:"privacyPolicy" v:"required#隐私政策内容不能为空" dc:"隐私政策内容"`
	UserAgreement string `json:"userAgreement" v:"required#用户协议内容不能为空" dc:"用户协议内容"`
}

// AgreementSettingsSaveRes 保存协议设置响应
type AgreementSettingsSaveRes struct {
	g.Meta    `mime:"application/json" example:"json"`
	IsSuccess bool `json:"isSuccess" dc:"是否成功"`
}

// WxAgreementGetReq 微信客户端获取协议设置请求
type WxAgreementGetReq struct {
	g.Meta `path:"/wx/agreement/get" method:"get" tags:"客户端基础设置" summary:"获取协议设置" description:"微信客户端获取隐私政策和用户协议设置"`
	Type   string `json:"type" v:"required|in:privacy,user#类型不能为空|类型只能是privacy或user" in:"query" dc:"协议类型：privacy-隐私政策，user-用户协议"`
}

// WxAgreementGetRes 微信客户端获取协议设置响应
type WxAgreementGetRes struct {
	g.Meta  `mime:"application/json" example:"json"`
	Content string `json:"content" dc:"协议内容"`
}
