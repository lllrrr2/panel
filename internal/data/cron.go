package data

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-rat/utils/str"
	"github.com/leonelquinteros/gotext"
	"gorm.io/gorm"

	"github.com/tnborg/panel/internal/app"
	"github.com/tnborg/panel/internal/biz"
	"github.com/tnborg/panel/internal/http/request"
	"github.com/tnborg/panel/pkg/io"
	"github.com/tnborg/panel/pkg/os"
	"github.com/tnborg/panel/pkg/shell"
	"github.com/tnborg/panel/pkg/systemctl"
)

type cronRepo struct {
	t  *gotext.Locale
	db *gorm.DB
}

func NewCronRepo(t *gotext.Locale, db *gorm.DB) biz.CronRepo {
	return &cronRepo{
		t:  t,
		db: db,
	}
}

func (r *cronRepo) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&biz.Cron{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *cronRepo) List(page, limit uint) ([]*biz.Cron, int64, error) {
	cron := make([]*biz.Cron, 0)
	var total int64
	err := r.db.Model(&biz.Cron{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&cron).Error
	return cron, total, err
}

func (r *cronRepo) Get(id uint) (*biz.Cron, error) {
	cron := new(biz.Cron)
	if err := r.db.Where("id = ?", id).First(cron).Error; err != nil {
		return nil, err
	}

	return cron, nil
}

func (r *cronRepo) Create(req *request.CronCreate) error {
	var script string
	if req.Type == "backup" {
		if req.BackupType == "website" {
			script = fmt.Sprintf(`#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

panel-cli backup website -n '%s' -p '%s'
panel-cli backup clear -t website -f '%s' -s '%d' -p '%s'
`, req.Target, req.BackupPath, req.Target, req.Save, req.BackupPath)
		}
		if req.BackupType == "mysql" || req.BackupType == "postgres" {
			script = fmt.Sprintf(`#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

panel-cli backup database -t '%s' -n '%s' -p '%s'
panel-cli backup clear -t '%s' -f '%s' -s '%d' -p '%s'
`, req.BackupType, req.Target, req.BackupPath, req.BackupType, req.Target, req.Save, req.BackupPath)
		}
	}
	if req.Type == "cutoff" {
		script = fmt.Sprintf(`#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

panel-cli cutoff website -n '%s' -p '%s'
panel-cli cutoff clear -t website -f '%s' -s '%d' -p '%s'
`, req.Target, req.BackupPath, req.Target, req.Save, req.BackupPath)
	}
	if req.Type == "shell" {
		script = req.Script
	}

	shellDir := fmt.Sprintf("%s/server/cron/", app.Root)
	shellLogDir := fmt.Sprintf("%s/server/cron/logs/", app.Root)
	if !io.Exists(shellDir) {
		return errors.New(r.t.Get("cron directory %s not exists", shellDir))
	}
	if !io.Exists(shellLogDir) {
		return errors.New(r.t.Get("cron log directory %s not exists", shellLogDir))
	}
	shellFile := strconv.Itoa(int(time.Now().Unix())) + str.Random(16)
	if err := io.Write(filepath.Join(shellDir, shellFile+".sh"), script, 0700); err != nil {
		return errors.New(err.Error())
	}
	// 编码转换
	_, _ = shell.Execf("dos2unix %s%s.sh", shellDir, shellFile)

	cron := new(biz.Cron)
	cron.Name = req.Name
	cron.Type = req.Type
	cron.Status = true
	cron.Time = req.Time
	cron.Shell = shellDir + shellFile + ".sh"
	cron.Log = shellLogDir + shellFile + ".log"

	if err := r.db.Create(cron).Error; err != nil {
		return err
	}
	if err := r.addToSystem(cron); err != nil {
		return err
	}

	return nil
}

func (r *cronRepo) Update(req *request.CronUpdate) error {
	cron, err := r.Get(req.ID)
	if err != nil {
		return err
	}

	cron.Time = req.Time
	cron.Name = req.Name
	if err = r.db.Save(cron).Error; err != nil {
		return err
	}

	if err = io.Write(cron.Shell, req.Script, 0700); err != nil {
		return err
	}
	if out, err := shell.Execf("dos2unix %s", cron.Shell); err != nil {
		return errors.New(out)
	}

	if err = r.deleteFromSystem(cron); err != nil {
		return err
	}
	if cron.Status {
		if err = r.addToSystem(cron); err != nil {
			return err
		}
	}

	return nil
}

func (r *cronRepo) Delete(id uint) error {
	cron, err := r.Get(id)
	if err != nil {
		return err
	}

	if err = r.deleteFromSystem(cron); err != nil {
		return err
	}
	if err = io.Remove(cron.Shell); err != nil {
		return err
	}

	return r.db.Delete(cron).Error
}

func (r *cronRepo) Status(id uint, status bool) error {
	cron, err := r.Get(id)
	if err != nil {
		return err
	}

	if err = r.deleteFromSystem(cron); err != nil {
		return err
	}
	if status {
		if err = r.addToSystem(cron); err != nil {
			return err
		}
	}

	cron.Status = status

	return r.db.Save(cron).Error
}

// addToSystem 添加到系统
func (r *cronRepo) addToSystem(cron *biz.Cron) error {
	if _, err := shell.Execf(`( crontab -l; echo "%s %s >> %s 2>&1" ) | sort - | uniq - | crontab -`, cron.Time, cron.Shell, cron.Log); err != nil {
		return err
	}

	return r.restartCron()
}

// deleteFromSystem 从系统中删除
func (r *cronRepo) deleteFromSystem(cron *biz.Cron) error {
	if _, err := shell.Execf(`( crontab -l | grep -v -F "%s %s >> %s 2>&1" ) | crontab -`, cron.Time, cron.Shell, cron.Log); err != nil {
		return err
	}

	return r.restartCron()
}

// restartCron 重启 cron 服务
func (r *cronRepo) restartCron() error {
	if os.IsRHEL() {
		return systemctl.Restart("crond")
	}

	if os.IsDebian() || os.IsUbuntu() {
		return systemctl.Restart("cron")
	}

	return errors.New(r.t.Get("unsupported system"))
}
