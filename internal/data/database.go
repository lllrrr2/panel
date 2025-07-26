package data

import (
	"errors"
	"fmt"
	"slices"

	"github.com/leonelquinteros/gotext"
	"gorm.io/gorm"

	"github.com/tnborg/panel/internal/biz"
	"github.com/tnborg/panel/internal/http/request"
	"github.com/tnborg/panel/pkg/db"
)

type databaseRepo struct {
	t      *gotext.Locale
	db     *gorm.DB
	server biz.DatabaseServerRepo
	user   biz.DatabaseUserRepo
}

func NewDatabaseRepo(t *gotext.Locale, db *gorm.DB, server biz.DatabaseServerRepo, user biz.DatabaseUserRepo) biz.DatabaseRepo {
	return &databaseRepo{
		t:      t,
		db:     db,
		server: server,
		user:   user,
	}
}

func (r databaseRepo) List(page, limit uint) ([]*biz.Database, int64, error) {
	var databaseServer []*biz.DatabaseServer
	if err := r.db.Model(&biz.DatabaseServer{}).Order("id desc").Find(&databaseServer).Error; err != nil {
		return nil, 0, err
	}

	database := make([]*biz.Database, 0)
	for _, server := range databaseServer {
		switch server.Type {
		case biz.DatabaseTypeMysql:
			mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
			if err == nil {
				if databases, err := mysql.Databases(); err == nil {
					for item := range slices.Values(databases) {
						database = append(database, &biz.Database{
							Type:     biz.DatabaseTypeMysql,
							Name:     item.Name,
							Server:   server.Name,
							ServerID: server.ID,
							Encoding: item.CharSet,
						})
					}
				}
				_ = mysql.Close()
			}
		case biz.DatabaseTypePostgresql:
			postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
			if err == nil {
				if databases, err := postgres.Databases(); err == nil {
					for item := range slices.Values(databases) {
						database = append(database, &biz.Database{
							Type:     biz.DatabaseTypePostgresql,
							Name:     item.Name,
							Server:   server.Name,
							ServerID: server.ID,
							Encoding: item.Encoding,
							Comment:  item.Comment,
						})
					}
				}
				_ = postgres.Close()
			}
		}
	}

	return database[(page-1)*limit:], int64(len(database)), nil
}

func (r databaseRepo) Create(req *request.DatabaseCreate) error {
	server, err := r.server.Get(req.ServerID)
	if err != nil {
		return err
	}

	switch server.Type {
	case biz.DatabaseTypeMysql:
		mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err != nil {
			return err
		}
		defer func(mysql *db.MySQL) {
			_ = mysql.Close()
		}(mysql)
		if req.CreateUser {
			if err = r.user.Create(&request.DatabaseUserCreate{
				ServerID: req.ServerID,
				Username: req.Username,
				Password: req.Password,
				Host:     req.Host,
			}); err != nil {
				return err
			}
		}
		if err = mysql.DatabaseCreate(req.Name); err != nil {
			return err
		}
		if req.Username != "" {
			if err = mysql.PrivilegesGrant(req.Username, req.Name, req.Host); err != nil {
				return err
			}
		}
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err != nil {
			return err
		}
		defer func(postgres *db.Postgres) {
			_ = postgres.Close()
		}(postgres)
		if req.CreateUser {
			if err = r.user.Create(&request.DatabaseUserCreate{
				ServerID: req.ServerID,
				Username: req.Username,
				Password: req.Password,
				Host:     req.Host,
			}); err != nil {
				return err
			}
		}
		if err = postgres.DatabaseCreate(req.Name); err != nil {
			return err
		}
		if req.Username != "" {
			if err = postgres.PrivilegesGrant(req.Username, req.Name); err != nil {
				return err
			}
		}
		if err = postgres.DatabaseComment(req.Name, req.Comment); err != nil {
			return err
		}
	}

	return nil
}

func (r databaseRepo) Delete(serverID uint, name string) error {
	server, err := r.server.Get(serverID)
	if err != nil {
		return err
	}

	switch server.Type {
	case biz.DatabaseTypeMysql:
		mysql, err := db.NewMySQL(server.Username, server.Password, fmt.Sprintf("%s:%d", server.Host, server.Port))
		if err != nil {
			return err
		}
		defer func(mysql *db.MySQL) {
			_ = mysql.Close()
		}(mysql)
		return mysql.DatabaseDrop(name)
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err != nil {
			return err
		}
		defer func(postgres *db.Postgres) {
			_ = postgres.Close()
		}(postgres)
		return postgres.DatabaseDrop(name)
	}

	return nil
}

func (r databaseRepo) Comment(req *request.DatabaseComment) error {
	server, err := r.server.Get(req.ServerID)
	if err != nil {
		return err
	}

	switch server.Type {
	case biz.DatabaseTypeMysql:
		return errors.New(r.t.Get("mysql not support database comment"))
	case biz.DatabaseTypePostgresql:
		postgres, err := db.NewPostgres(server.Username, server.Password, server.Host, server.Port)
		if err != nil {
			return err
		}
		defer func(postgres *db.Postgres) {
			_ = postgres.Close()
		}(postgres)
		return postgres.DatabaseComment(req.Name, req.Comment)
	}

	return nil
}
