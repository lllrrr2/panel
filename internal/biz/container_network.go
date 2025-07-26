package biz

import (
	"github.com/tnborg/panel/internal/http/request"
	"github.com/tnborg/panel/pkg/types"
)

type ContainerNetworkRepo interface {
	List() ([]types.ContainerNetwork, error)
	Create(req *request.ContainerNetworkCreate) (string, error)
	Remove(id string) error
	Prune() error
}
