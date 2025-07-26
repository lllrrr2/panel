//go:build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/tnborg/panel/internal/app"
	"github.com/tnborg/panel/internal/apps"
	"github.com/tnborg/panel/internal/bootstrap"
	"github.com/tnborg/panel/internal/data"
	"github.com/tnborg/panel/internal/route"
	"github.com/tnborg/panel/internal/service"
)

// initCli init command line.
func initCli() (*app.Cli, error) {
	panic(wire.Build(bootstrap.ProviderSet, route.ProviderSet, service.ProviderSet, data.ProviderSet, apps.ProviderSet, app.NewCli))
}
