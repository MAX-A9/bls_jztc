package main

import (
	_ "demo/internal/packed"
	// 确保聊天逻辑包被初始化

	"github.com/gogf/gf/v2/os/gctx"

	"demo/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
