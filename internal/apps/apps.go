package apps

import (
	"github.com/google/wire"

	"github.com/tnborg/panel/internal/apps/codeserver"
	"github.com/tnborg/panel/internal/apps/docker"
	"github.com/tnborg/panel/internal/apps/fail2ban"
	"github.com/tnborg/panel/internal/apps/frp"
	"github.com/tnborg/panel/internal/apps/gitea"
	"github.com/tnborg/panel/internal/apps/memcached"
	"github.com/tnborg/panel/internal/apps/minio"
	"github.com/tnborg/panel/internal/apps/mysql"
	"github.com/tnborg/panel/internal/apps/nginx"
	"github.com/tnborg/panel/internal/apps/php74"
	"github.com/tnborg/panel/internal/apps/php80"
	"github.com/tnborg/panel/internal/apps/php81"
	"github.com/tnborg/panel/internal/apps/php82"
	"github.com/tnborg/panel/internal/apps/php83"
	"github.com/tnborg/panel/internal/apps/php84"
	"github.com/tnborg/panel/internal/apps/phpmyadmin"
	"github.com/tnborg/panel/internal/apps/podman"
	"github.com/tnborg/panel/internal/apps/postgresql"
	"github.com/tnborg/panel/internal/apps/pureftpd"
	"github.com/tnborg/panel/internal/apps/redis"
	"github.com/tnborg/panel/internal/apps/rsync"
	"github.com/tnborg/panel/internal/apps/s3fs"
	"github.com/tnborg/panel/internal/apps/supervisor"
)

var ProviderSet = wire.NewSet(
	codeserver.NewApp,
	docker.NewApp,
	fail2ban.NewApp,
	frp.NewApp,
	gitea.NewApp,
	memcached.NewApp,
	minio.NewApp,
	mysql.NewApp,
	nginx.NewApp,
	php74.NewApp,
	php80.NewApp,
	php81.NewApp,
	php82.NewApp,
	php83.NewApp,
	php84.NewApp,
	phpmyadmin.NewApp,
	podman.NewApp,
	postgresql.NewApp,
	pureftpd.NewApp,
	redis.NewApp,
	rsync.NewApp,
	s3fs.NewApp,
	supervisor.NewApp,
)
