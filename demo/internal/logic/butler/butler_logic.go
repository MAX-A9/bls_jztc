package butler

import (
	"demo/internal/service"
)

func init() {
	service.SetButler(New())
}
