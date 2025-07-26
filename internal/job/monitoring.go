package job

import (
	"log/slog"
	"time"

	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/tnborg/panel/internal/app"
	"github.com/tnborg/panel/internal/biz"
	"github.com/tnborg/panel/pkg/tools"
)

// Monitoring 系统监控
type Monitoring struct {
	db          *gorm.DB
	log         *slog.Logger
	settingRepo biz.SettingRepo
}

func NewMonitoring(db *gorm.DB, log *slog.Logger, setting biz.SettingRepo) *Monitoring {
	return &Monitoring{
		db:          db,
		log:         log,
		settingRepo: setting,
	}
}

func (r *Monitoring) Run() {
	if app.Status != app.StatusNormal {
		return
	}

	monitor, err := r.settingRepo.Get(biz.SettingKeyMonitor)
	if err != nil || !cast.ToBool(monitor) {
		return
	}

	info := tools.CurrentInfo(nil, nil)

	// 去除部分数据以减少数据库存储
	info.Disk = nil
	info.Cpus = nil

	if app.Status != app.StatusNormal {
		return
	}

	if err = r.db.Create(&biz.Monitor{Info: info}).Error; err != nil {
		r.log.Warn("[Monitor] failed to create monitor record", slog.Any("err", err))
		return
	}

	// 删除过期数据
	dayStr, err := r.settingRepo.Get(biz.SettingKeyMonitorDays)
	if err != nil {
		return
	}
	day := cast.ToInt(dayStr)
	if day <= 0 || app.Status != app.StatusNormal {
		return
	}
	if err = r.db.Where("created_at < ?", time.Now().AddDate(0, 0, -day).Format(time.DateTime)).Delete(&biz.Monitor{}).Error; err != nil {
		r.log.Warn("[Monitor] failed to delete monitor record", slog.Any("err", err))
		return
	}
}
