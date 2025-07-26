package data

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/tnborg/panel/internal/biz"
	"github.com/tnborg/panel/internal/http/request"
	"github.com/tnborg/panel/pkg/shell"
	"github.com/tnborg/panel/pkg/tools"
	"github.com/tnborg/panel/pkg/types"
	"github.com/tnborg/panel/pkg/types/docker/volume"
)

type containerVolumeRepo struct {
	client *resty.Client
}

func NewContainerVolumeRepo() biz.ContainerVolumeRepo {
	return &containerVolumeRepo{
		client: getDockerClient("/var/run/docker.sock"),
	}
}

// List 列出存储卷
func (r *containerVolumeRepo) List() ([]types.ContainerVolume, error) {
	var resp volume.ListResponse
	_, err := r.client.R().SetResult(&resp).Get("/volumes")
	if err != nil {
		return nil, err
	}

	var volumes []types.ContainerVolume
	for _, item := range resp.Volumes {
		volumes = append(volumes, types.ContainerVolume{
			Name:       item.Name,
			Driver:     item.Driver,
			Scope:      item.Scope,
			MountPoint: item.Mountpoint,
			CreatedAt:  item.CreatedAt,
			Labels:     types.MapToKV(item.Labels),
			Options:    types.MapToKV(item.Options),
			RefCount:   item.UsageData.RefCount,
			Size:       tools.FormatBytes(float64(item.UsageData.Size)),
		})
	}

	slices.SortFunc(volumes, func(a types.ContainerVolume, b types.ContainerVolume) int {
		return strings.Compare(a.Name, b.Name)
	})

	return volumes, nil
}

// Create 创建存储卷
func (r *containerVolumeRepo) Create(req *request.ContainerVolumeCreate) (string, error) {
	var sb strings.Builder
	sb.WriteString("docker volume create")
	sb.WriteString(fmt.Sprintf(" %s", req.Name))

	if req.Driver != "" {
		sb.WriteString(fmt.Sprintf(" --driver %s", req.Driver))
	}
	for _, label := range req.Labels {
		sb.WriteString(fmt.Sprintf(" --label %s=%s", label.Key, label.Value))
	}

	for _, option := range req.Options {
		sb.WriteString(fmt.Sprintf(" --opt %s=%s", option.Key, option.Value))
	}

	return shell.Exec(sb.String())
}

// Remove 删除存储卷
func (r *containerVolumeRepo) Remove(id string) error {
	_, err := shell.ExecfWithTimeout(2*time.Minute, "docker volume rm -f %s", id)
	return err
}

// Prune 清理未使用的存储卷
func (r *containerVolumeRepo) Prune() error {
	_, err := shell.ExecfWithTimeout(2*time.Minute, "docker volume prune -f")
	return err
}
