package biz

import (
	"github.com/tnborg/panel/internal/http/request"
	"github.com/tnborg/panel/pkg/types"
)

type ContainerImageRepo interface {
	List() ([]types.ContainerImage, error)
	Pull(req *request.ContainerImagePull) error
	Remove(id string) error
	Prune() error
}
