package payment

import (
	"context"
	"os"

	v1 "demo/api/payment/v1"
	"demo/internal/model"
	"demo/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

// 配置文件路径
const (
	configFilePath = "manifest/config/payment.json"
)

// 支付服务实现
type paymentImpl struct{}

// New 创建支付服务实例
func New() service.PaymentService {
	return &paymentImpl{}
}

// GetConfig 获取支付配置
func (s *paymentImpl) GetConfig(ctx context.Context, req *v1.PaymentConfigReq) (res *v1.PaymentConfigRes, err error) {
	config, err := s.readConfig(ctx)
	if err != nil {
		return nil, err
	}

	// 创建响应对象
	res = &v1.PaymentConfigRes{
		AppId:     config.AppId,
		MchId:     config.MchId,
		ApiKey:    s.maskSecret(config.ApiKey), // 掩码处理密钥
		NotifyUrl: config.NotifyUrl,
		IsEnabled: config.IsEnabled,
	}

	return res, nil
}

// SaveConfig 保存支付配置
func (s *paymentImpl) SaveConfig(ctx context.Context, req *v1.SavePaymentConfigReq) (res *v1.SavePaymentConfigRes, err error) {
	// 获取现有配置
	oldConfig, err := s.readConfig(ctx)
	if err != nil {
		// 如果读取失败但是不是因为文件不存在，则返回错误
		if !os.IsNotExist(err) {
			g.Log().Error(ctx, "读取支付配置失败", err)
			return nil, gerror.New("读取现有配置失败")
		}
	}

	// 创建新的配置对象
	config := &model.PaymentConfig{
		AppId:     req.AppId,
		MchId:     req.MchId,
		NotifyUrl: req.NotifyUrl,
		IsEnabled: req.IsEnabled,
	}

	// 处理API密钥
	// 如果是掩码（没有修改），则使用旧配置中的值
	if req.ApiKey == "********" && oldConfig != nil {
		config.ApiKey = oldConfig.ApiKey
	} else {
		config.ApiKey = req.ApiKey
	}

	// 创建配置目录（如果不存在）
	if err = gfile.Mkdir(gfile.Dir(configFilePath)); err != nil {
		g.Log().Error(ctx, "创建配置目录失败", err)
		return nil, gerror.New("创建配置目录失败")
	}

	// 将配置对象转换为JSON并保存到文件
	if err = gfile.PutContents(configFilePath, gjson.MustEncodeString(config)); err != nil {
		g.Log().Error(ctx, "保存支付配置失败", err)
		return nil, gerror.New("保存配置失败")
	}

	// 创建响应对象
	res = &v1.SavePaymentConfigRes{
		Success: true,
		Message: "配置保存成功",
	}

	// 记录日志
	g.Log().Info(ctx, "支付配置已更新", gtime.Now())

	return res, nil
}

// readConfig 读取支付配置
func (s *paymentImpl) readConfig(ctx context.Context) (config *model.PaymentConfig, err error) {
	// 检查配置文件是否存在
	if !gfile.Exists(configFilePath) {
		// 文件不存在时返回默认配置
		return &model.PaymentConfig{
			AppId:     "",
			MchId:     "",
			ApiKey:    "",
			NotifyUrl: "",
			IsEnabled: false,
		}, nil
	}

	// 读取配置文件内容
	content := gfile.GetContents(configFilePath)
	if content == "" {
		g.Log().Error(ctx, "读取配置文件失败")
		return nil, gerror.New("读取配置文件失败")
	}

	// 解析JSON内容
	config = &model.PaymentConfig{}
	if err = gjson.DecodeTo(content, config); err != nil {
		g.Log().Error(ctx, "解析配置文件失败", err)
		return nil, gerror.New("解析配置文件失败")
	}

	return config, nil
}

// maskSecret 对密钥进行掩码处理
func (s *paymentImpl) maskSecret(secret string) string {
	if secret == "" {
		return ""
	}
	return "********"
}
