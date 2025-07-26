package service

import (
	"net/http"

	"github.com/go-rat/chix"

	"github.com/tnborg/panel/internal/biz"
	"github.com/tnborg/panel/internal/http/request"
)

type ContainerImageService struct {
	containerImageRepo biz.ContainerImageRepo
}

func NewContainerImageService(containerImage biz.ContainerImageRepo) *ContainerImageService {
	return &ContainerImageService{
		containerImageRepo: containerImage,
	}
}

func (s *ContainerImageService) List(w http.ResponseWriter, r *http.Request) {
	images, err := s.containerImageRepo.List()
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	paged, total := Paginate(r, images)

	Success(w, chix.M{
		"total": total,
		"items": paged,
	})
}

func (s *ContainerImageService) Pull(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerImagePull](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.containerImageRepo.Pull(req); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *ContainerImageService) Remove(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ContainerImageID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.containerImageRepo.Remove(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}

func (s *ContainerImageService) Prune(w http.ResponseWriter, r *http.Request) {
	if err := s.containerImageRepo.Prune(); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}
