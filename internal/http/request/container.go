package request

import "github.com/tnborg/panel/pkg/types"

type ContainerID struct {
	ID string `json:"id" form:"id" validate:"required"`
}

type ContainerRename struct {
	ID   string `form:"id" json:"id" validate:"required"`
	Name string `form:"name" json:"name" validate:"required"`
}

type ContainerCreate struct {
	Name            string                           `form:"name" json:"name" validate:"required"`
	Image           string                           `form:"image" json:"image" validate:"required"`
	Ports           []types.ContainerPort            `form:"ports" json:"ports"`
	Network         string                           `form:"network" json:"network"`
	Volumes         []types.ContainerContainerVolume `form:"volumes" json:"volumes"`
	Labels          []types.KV                       `form:"labels" json:"labels"`
	Env             []types.KV                       `form:"env" json:"env"`
	Entrypoint      []string                         `form:"entrypoint" json:"entrypoint"`
	Command         []string                         `form:"command" json:"command"`
	RestartPolicy   string                           `form:"restart_policy" json:"restart_policy"`
	AutoRemove      bool                             `form:"auto_remove" json:"auto_remove"`
	Privileged      bool                             `form:"privileged" json:"privileged"`
	OpenStdin       bool                             `form:"openStdin" json:"open_stdin"`
	PublishAllPorts bool                             `form:"publish_all_ports" json:"publish_all_ports"`
	Tty             bool                             `form:"tty" json:"tty"`
	CPUShares       int64                            `form:"cpu_shares" json:"cpu_shares"`
	CPUs            int64                            `form:"cpus" json:"cpus"`
	Memory          int64                            `form:"memory" json:"memory"`
}
