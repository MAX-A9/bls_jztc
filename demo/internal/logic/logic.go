package logic

import (
	_ "demo/internal/logic/backend"
	_ "demo/internal/logic/client"
	_ "demo/internal/logic/content"
	_ "demo/internal/logic/dbtool"
	_ "demo/internal/logic/log"
	_ "demo/internal/logic/order"
	_ "demo/internal/logic/package"
	_ "demo/internal/logic/payment"
	_ "demo/internal/logic/region"
	_ "demo/internal/logic/storage"
	_ "demo/internal/logic/user"
	_ "demo/internal/logic/wechat_pay"
)

func init() {
	// 逻辑层初始化代码
}
