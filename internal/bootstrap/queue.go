package bootstrap

import (
	"github.com/tnborg/panel/pkg/queue"
)

func NewQueue() *queue.Queue {
	return queue.New(100)
}
