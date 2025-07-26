package service

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-rat/chix"
	"github.com/leonelquinteros/gotext"
	"github.com/spf13/cast"

	"github.com/tnborg/panel/internal/app"
	"github.com/tnborg/panel/internal/http/request"
	"github.com/tnborg/panel/pkg/io"
	"github.com/tnborg/panel/pkg/ntp"
	"github.com/tnborg/panel/pkg/shell"
	"github.com/tnborg/panel/pkg/tools"
	"github.com/tnborg/panel/pkg/types"
)

type ToolboxSystemService struct {
	t *gotext.Locale
}

func NewToolboxSystemService(t *gotext.Locale) *ToolboxSystemService {
	return &ToolboxSystemService{
		t: t,
	}
}

// GetDNS 获取 DNS 信息
func (s *ToolboxSystemService) GetDNS(w http.ResponseWriter, r *http.Request) {
	raw, err := io.Read("/etc/resolv.conf")
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	match := regexp.MustCompile(`nameserver\s+(\S+)`).FindAllStringSubmatch(raw, -1)
	dns := make([]string, 0)
	for _, m := range match {
		dns = append(dns, m[1])
	}

	Success(w, dns)
}

// UpdateDNS 设置 DNS 信息
func (s *ToolboxSystemService) UpdateDNS(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ToolboxSystemDNS](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	var dns string
	dns += "nameserver " + req.DNS1 + "\n"
	dns += "nameserver " + req.DNS2 + "\n"

	if err := io.Write("/etc/resolv.conf", dns, 0644); err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("failed to update DNS: %v", err))
		return
	}

	Success(w, nil)
}

// GetSWAP 获取 SWAP 信息
func (s *ToolboxSystemService) GetSWAP(w http.ResponseWriter, r *http.Request) {
	var total, used, free string
	var size int64
	if io.Exists(filepath.Join(app.Root, "swap")) {
		file, err := os.Stat(filepath.Join(app.Root, "swap"))
		if err != nil {
			Error(w, http.StatusInternalServerError, s.t.Get("failed to get SWAP: %v", err))
			return
		}

		size = file.Size() / 1024 / 1024
		total = tools.FormatBytes(float64(file.Size()))
	} else {
		size = 0
		total = "0.00 B"
	}

	raw, err := shell.Execf("free | grep Swap")
	if err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("failed to get SWAP: %v", err))
		return
	}

	match := regexp.MustCompile(`Swap:\s+(\d+)\s+(\d+)\s+(\d+)`).FindStringSubmatch(raw)
	if len(match) >= 4 {
		used = tools.FormatBytes(cast.ToFloat64(match[2]) * 1024)
		free = tools.FormatBytes(cast.ToFloat64(match[3]) * 1024)
	}

	Success(w, chix.M{
		"total": total,
		"size":  size,
		"used":  used,
		"free":  free,
	})
}

// UpdateSWAP 设置 SWAP 信息
func (s *ToolboxSystemService) UpdateSWAP(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ToolboxSystemSWAP](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if io.Exists(filepath.Join(app.Root, "swap")) {
		if _, err = shell.Execf("swapoff '%s'", filepath.Join(app.Root, "swap")); err != nil {
			Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
		if _, err = shell.Execf("rm -f '%s'", filepath.Join(app.Root, "swap")); err != nil {
			Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
		if _, err = shell.Execf(`sed -i "\|^%s|d" /etc/fstab`, filepath.Join(app.Root, "swap")); err != nil {
			Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
	}

	if req.Size > 1 {
		var free string
		free, err = shell.Execf("df -k %s | awk '{print $4}' | tail -n 1", app.Root)
		if err != nil {
			Error(w, http.StatusInternalServerError, s.t.Get("failed to get disk space: %v", err))
			return
		}
		if cast.ToInt64(free)*1024 < req.Size*1024*1024 {
			Error(w, http.StatusInternalServerError, s.t.Get("disk space is insufficient, current free %s", tools.FormatBytes(cast.ToFloat64(free))))
			return
		}

		btrfsCheck, _ := shell.Execf("df -T %s | awk '{print $2}' | tail -n 1", app.Root)
		if strings.Contains(btrfsCheck, "btrfs") {
			if _, err = shell.Execf("btrfs filesystem mkswapfile --size %dM --uuid clear %s", req.Size, filepath.Join(app.Root, "swap")); err != nil {
				Error(w, http.StatusInternalServerError, "%v", err)
				return
			}
		} else {
			if _, err = shell.Execf("dd if=/dev/zero of=%s bs=1M count=%d", filepath.Join(app.Root, "swap"), req.Size); err != nil {
				Error(w, http.StatusInternalServerError, "%v", err)
				return
			}
			if _, err = shell.Execf("mkswap -f '%s'", filepath.Join(app.Root, "swap")); err != nil {
				Error(w, http.StatusInternalServerError, "%v", err)
				return
			}
			if err = io.Chmod(filepath.Join(app.Root, "swap"), 0600); err != nil {
				Error(w, http.StatusInternalServerError, s.t.Get("failed to set SWAP permission: %v", err))
				return
			}
		}
		if _, err = shell.Execf("swapon '%s'", filepath.Join(app.Root, "swap")); err != nil {
			Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
		if _, err = shell.Execf("echo '%s    swap    swap    defaults    0 0' >> /etc/fstab", filepath.Join(app.Root, "swap")); err != nil {
			Error(w, http.StatusInternalServerError, "%v", err)
			return
		}
	}

	Success(w, nil)
}

// GetTimezone 获取时区
func (s *ToolboxSystemService) GetTimezone(w http.ResponseWriter, r *http.Request) {
	raw, err := shell.Execf("timedatectl | grep zone")
	if err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("failed to get timezone: %v", err))
		return
	}

	match := regexp.MustCompile(`zone:\s+(\S+)`).FindStringSubmatch(raw)
	if len(match) == 0 {
		match = append(match, "")
	}

	zonesRaw, err := shell.Execf("timedatectl list-timezones")
	if err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("failed to get available timezones: %v", err))
		return
	}
	zones := strings.Split(zonesRaw, "\n")

	var zonesList []types.LV
	for _, z := range zones {
		zonesList = append(zonesList, types.LV{
			Label: z,
			Value: z,
		})
	}

	Success(w, chix.M{
		"timezone":  match[1],
		"timezones": zonesList,
	})
}

// UpdateTimezone 设置时区
func (s *ToolboxSystemService) UpdateTimezone(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ToolboxSystemTimezone](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if _, err = shell.Execf("timedatectl set-timezone '%s'", req.Timezone); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

// UpdateTime 设置时间
func (s *ToolboxSystemService) UpdateTime(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ToolboxSystemTime](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = ntp.UpdateSystemTime(req.Time); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)

}

// SyncTime 同步时间
func (s *ToolboxSystemService) SyncTime(w http.ResponseWriter, r *http.Request) {
	now, err := ntp.Now()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	if err = ntp.UpdateSystemTime(now); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

// GetHostname 获取主机名
func (s *ToolboxSystemService) GetHostname(w http.ResponseWriter, r *http.Request) {
	hostname, _ := io.Read("/etc/hostname")
	Success(w, strings.TrimSpace(hostname))
}

// UpdateHostname 设置主机名
func (s *ToolboxSystemService) UpdateHostname(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ToolboxSystemHostname](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if _, err = shell.Execf("hostnamectl set-hostname '%s'", req.Hostname); err != nil {
		// 直接写 /etc/hostname
		if err = io.Write("/etc/hostname", req.Hostname, 0644); err != nil {
			Error(w, http.StatusInternalServerError, s.t.Get("failed to set hostname: %v", err))
			return
		}
	}

	Success(w, nil)
}

// GetHosts 获取 hosts 信息
func (s *ToolboxSystemService) GetHosts(w http.ResponseWriter, r *http.Request) {
	hosts, _ := io.Read("/etc/hosts")
	Success(w, hosts)
}

// UpdateHosts 设置 hosts 信息
func (s *ToolboxSystemService) UpdateHosts(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ToolboxSystemHosts](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = io.Write("/etc/hosts", req.Hosts, 0644); err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("failed to set hosts: %v", err))
		return
	}

	Success(w, nil)
}

// UpdateRootPassword 设置 root 密码
func (s *ToolboxSystemService) UpdateRootPassword(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ToolboxSystemPassword](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	req.Password = strings.ReplaceAll(req.Password, `'`, `\'`)
	if _, err = shell.Execf(`yes '%s' | passwd root`, req.Password); err != nil {
		Error(w, http.StatusInternalServerError, "%v", s.t.Get("failed to set root password: %v", err))
		return
	}

	Success(w, nil)
}
