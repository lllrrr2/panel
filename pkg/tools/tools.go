// Package tools 存放辅助方法
package tools

import (
	"context"
	"errors"
	"fmt"
	stdnet "net"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"

	"github.com/tnborg/panel/pkg/shell"
	"github.com/tnborg/panel/pkg/types"
)

// CurrentInfo 获取监控数据
func CurrentInfo(nets, disks []string) types.CurrentInfo {
	var res types.CurrentInfo
	res.Cpus, _ = cpu.Info()
	res.Percents, _ = cpu.Percent(100*time.Millisecond, true)
	percent, _ := cpu.Percent(100*time.Millisecond, false)
	if len(percent) > 0 {
		res.Percent = percent[0]
	}
	res.Load, _ = load.Avg()
	res.Host, _ = host.Info()
	res.Mem, _ = mem.VirtualMemory()
	res.Swap, _ = mem.SwapMemory()
	// 硬盘IO
	ioCounters, _ := disk.IOCounters(disks...)
	for _, info := range ioCounters {
		res.DiskIO = append(res.DiskIO, info)
	}
	// 硬盘使用
	var excludes = []string{"/dev", "/boot", "/sys", "/dev", "/run", "/proc", "/usr", "/var", "/snap"}
	excludes = append(excludes, "/mnt/cdrom") // CDROM
	excludes = append(excludes, "/mnt/wsl")   // Windows WSL
	res.Disk, _ = disk.Partitions(false)
	res.Disk = slices.DeleteFunc(res.Disk, func(d disk.PartitionStat) bool {
		for _, exclude := range excludes {
			if strings.HasPrefix(d.Mountpoint, exclude) {
				return true
			}
			// 去除内存盘和overlay容器盘
			if slices.Contains([]string{"tmpfs", "overlay"}, d.Fstype) {
				continue
			}
		}
		return false
	})
	// 分区使用
	for _, partition := range res.Disk {
		usage, _ := disk.Usage(partition.Mountpoint)
		res.DiskUsage = append(res.DiskUsage, *usage)
	}
	// 网络
	if len(nets) == 0 {
		netInfo, _ := net.IOCounters(false)
		res.Net = netInfo
	} else {
		var netStats []net.IOCountersStat
		netInfo, _ := net.IOCounters(true)
		for _, state := range netInfo {
			if slices.Contains(nets, state.Name) {
				netStats = append(netStats, state)
			}
		}
		res.Net = netStats
	}

	res.Time = time.Now()
	return res
}

// RestartPanel 重启面板
func RestartPanel() {
	_ = shell.ExecfAsync("sleep 1 && systemctl restart panel")
}

// IsChina 是否中国大陆
func IsChina() bool {
	client := resty.New()
	client.SetLogger(NoopLogger{})
	client.SetDisableWarn(true)
	client.SetTimeout(3 * time.Second)
	client.SetRetryCount(3)

	resp, err := client.R().Get("https://www.qualcomm.cn/cdn-cgi/trace")
	if err != nil || !resp.IsSuccess() {
		return false
	}

	if strings.Contains(resp.String(), "loc=CN") {
		return true
	}

	return false
}

// GetPublicIPv4 获取公网IPv4
func GetPublicIPv4() (string, error) {
	client := resty.New()
	client.SetLogger(NoopLogger{})
	client.SetDisableWarn(true)
	client.SetTimeout(3 * time.Second)
	client.SetRetryCount(3)
	client.SetTransport(&http.Transport{
		DialContext: func(ctx context.Context, network string, addr string) (stdnet.Conn, error) {
			return (&stdnet.Dialer{}).DialContext(ctx, "tcp4", addr)
		},
	})

	resp, err := client.R().Get("https://www.qualcomm.cn/cdn-cgi/trace")
	if err != nil || !resp.IsSuccess() {
		return "", errors.New("failed to get public ipv4 address")
	}

	return strings.TrimPrefix(strings.Split(resp.String(), "\n")[2], "ip="), nil
}

// GetPublicIPv6 获取公网IPv6
func GetPublicIPv6() (string, error) {
	client := resty.New()
	client.SetLogger(NoopLogger{})
	client.SetDisableWarn(true)
	client.SetTimeout(3 * time.Second)
	client.SetRetryCount(3)
	client.SetTransport(&http.Transport{
		DialContext: func(ctx context.Context, network string, addr string) (stdnet.Conn, error) {
			return (&stdnet.Dialer{}).DialContext(ctx, "tcp6", addr)
		},
	})

	resp, err := client.R().Get("https://www.qualcomm.cn/cdn-cgi/trace")
	if err != nil || !resp.IsSuccess() {
		return "", errors.New("failed to get public ipv6 address")
	}

	return strings.TrimPrefix(strings.Split(resp.String(), "\n")[2], "ip="), nil
}

// GetLocalIPv4 获取本地IPv4
func GetLocalIPv4() (string, error) {
	conn, err := stdnet.Dial("udp", "119.29.29.29:53")
	if err != nil {
		return "", err
	}
	defer func(conn stdnet.Conn) {
		_ = conn.Close()
	}(conn)

	local := conn.LocalAddr().(*stdnet.UDPAddr)
	return local.IP.String(), nil
}

// GetLocalIPv6 获取本地IPv6
func GetLocalIPv6() (string, error) {
	conn, err := stdnet.Dial("udp", "[2402:4e00::]:53")
	if err != nil {
		return "", err
	}
	defer func(conn stdnet.Conn) {
		_ = conn.Close()
	}(conn)

	local := conn.LocalAddr().(*stdnet.UDPAddr)
	return local.IP.String(), nil
}

// FormatBytes 格式化bytes
func FormatBytes(size float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	i := 0
	for ; size >= 1024 && i < len(units); i++ {
		size /= 1024
	}

	return fmt.Sprintf("%.2f %s", size, units[i])
}
