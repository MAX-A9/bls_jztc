// Package license 提供授权验证相关工具
package license

import (
	"context"
	"net"
	"strings"

	"github.com/gogf/gf/v2/os/glog"
)

// ValidateDomain 验证域名格式是否正确
func ValidateDomain(domain string) string {
	// 标准化处理域名
	domain = strings.ToLower(strings.TrimSpace(domain))

	// 去除可能存在的协议前缀
	if strings.HasPrefix(domain, "http://") {
		domain = domain[7:]
	} else if strings.HasPrefix(domain, "https://") {
		domain = domain[8:]
	}

	// 去除可能存在的路径和参数
	if idx := strings.Index(domain, "/"); idx > 0 {
		domain = domain[:idx]
	}

	// 去除可能存在的端口
	if idx := strings.Index(domain, ":"); idx > 0 {
		domain = domain[:idx]
	}

	return domain
}

// ExtractDomainFromHost 从主机地址中提取域名
func ExtractDomainFromHost(host string) string {
	// 去除端口号
	if idx := strings.Index(host, ":"); idx > 0 {
		host = host[:idx]
	}
	return host
}

// IsLocalhost 检查是否是本地环境
func IsLocalhost(domain string) bool {
	domain = strings.ToLower(strings.TrimSpace(domain))
	return domain == "localhost" || domain == "127.0.0.1" || domain == "::1"
}

// GetServerIP 获取当前服务器的IP地址
func GetServerIP() string {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		glog.Warning(context.Background(), "获取网络接口失败:", err)
		return "127.0.0.1"
	}

	// 检查所有网络接口
	for _, iface := range interfaces {
		// 跳过禁用的接口
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// 跳过回环接口
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 获取该接口的所有地址
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// 跳过IPv6地址和回环地址
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			// 如果不是私有IP（如10.x.x.x、192.168.x.x等），则很可能是外网IP
			if !isPrivateIP(ip.String()) {
				return ip.String()
			}
		}
	}

	// 如果没有找到非私有IP，则尝试获取第一个非回环IPv4地址
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			// 返回第一个有效的IPv4地址
			return ip.String()
		}
	}

	// 如果都没找到，返回本地回环地址
	return "127.0.0.1"
}

// isPrivateIP 检查IP是否为内网IP
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// 检查是否为私有IP范围
	privateIPBlocks := []string{
		"10.0.0.0/8",     // 10.0.0.0 - 10.255.255.255
		"172.16.0.0/12",  // 172.16.0.0 - 172.31.255.255
		"192.168.0.0/16", // 192.168.0.0 - 192.168.255.255
		"169.254.0.0/16", // 链路本地地址
	}

	for _, block := range privateIPBlocks {
		_, ipNet, err := net.ParseCIDR(block)
		if err != nil {
			continue
		}
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}
