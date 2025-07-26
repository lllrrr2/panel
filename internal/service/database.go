package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/tnborg/panel/internal/biz"
	"github.com/tnborg/panel/internal/http/request"
)

type DatabaseService struct {
	databaseRepo biz.DatabaseRepo
}

func NewDatabaseService(database biz.DatabaseRepo) *DatabaseService {
	return &DatabaseService{
		databaseRepo: database,
	}
}

func (s *DatabaseService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.Paginate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	databases, total, err := s.databaseRepo.List(req.Page, req.Limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, chix.M{
		"total": total,
		"items": databases,
	})
}

func (s *DatabaseService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.DatabaseCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.databaseRepo.Create(req); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *DatabaseService) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.DatabaseDelete](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.databaseRepo.Delete(req.ServerID, req.Name); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *DatabaseService) Comment(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.DatabaseComment](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.databaseRepo.Comment(req); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}
