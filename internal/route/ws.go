package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/tnborg/panel/internal/service"
)

type Ws struct {
	ws *service.WsService
}

func NewWs(ws *service.WsService) *Ws {
	return &Ws{
		ws: ws,
	}
}

func (route *Ws) Register(r *chi.Mux) {
	r.Route("/api/ws", func(r chi.Router) {
		r.Get("/ssh", route.ws.Session)
		r.Get("/exec", route.ws.Exec)
	})
}
