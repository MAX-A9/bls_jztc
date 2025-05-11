package region

import (
	"demo/internal/service"
)

func init() {
	// 注册地区服务
	service.SetRegion(New())
}
