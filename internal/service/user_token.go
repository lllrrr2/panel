package service

import (
	"net/http"
	"time"

	"github.com/go-rat/chix"
	"github.com/leonelquinteros/gotext"

	"github.com/tnborg/panel/internal/biz"
	"github.com/tnborg/panel/internal/http/request"
)

type UserTokenService struct {
	t             *gotext.Locale
	userTokenRepo biz.UserTokenRepo
}

func NewUserTokenService(t *gotext.Locale, userToken biz.UserTokenRepo) *UserTokenService {
	return &UserTokenService{
		t:             t,
		userTokenRepo: userToken,
	}
}

func (s *UserTokenService) List(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.UserTokenList](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	userTokens, total, err := s.userTokenRepo.List(req.UserID, req.Page, req.Limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, chix.M{
		"total": total,
		"items": userTokens,
	})
}

func (s *UserTokenService) Create(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.UserTokenCreate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	expiredAt := time.Unix(0, req.ExpiredAt*int64(time.Millisecond))
	if expiredAt.Before(time.Now()) {
		Error(w, http.StatusUnprocessableEntity, s.t.Get("expiration time must be greater than current time"))
		return
	}
	if expiredAt.After(time.Now().AddDate(10, 0, 0)) {
		Error(w, http.StatusUnprocessableEntity, s.t.Get("expiration time must be less than 10 years"))
		return
	}

	userToken, err := s.userTokenRepo.Create(req.UserID, req.IPs, expiredAt)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	// 手动组装响应，因为 Token 设置了 json:"-"
	Success(w, chix.M{
		"id":         userToken.ID,
		"user_id":    userToken.UserID,
		"token":      userToken.Token,
		"ips":        userToken.IPs,
		"expired_at": userToken.ExpiredAt,
		"created_at": userToken.CreatedAt,
		"updated_at": userToken.UpdatedAt,
	})
}

func (s *UserTokenService) Update(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.UserTokenUpdate](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	expiredAt := time.Unix(0, req.ExpiredAt*int64(time.Millisecond))
	if expiredAt.Before(time.Now()) {
		Error(w, http.StatusUnprocessableEntity, s.t.Get("expiration time must be greater than current time"))
		return
	}
	if expiredAt.After(time.Now().AddDate(10, 0, 0)) {
		Error(w, http.StatusUnprocessableEntity, s.t.Get("expiration time must be less than 10 years"))
		return
	}

	userToken, err := s.userTokenRepo.Update(req.ID, req.IPs, expiredAt)
	if err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, userToken)
}

func (s *UserTokenService) Delete(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.ID](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err = s.userTokenRepo.Delete(req.ID); err != nil {
		Error(w, http.StatusInternalServerError, "%v", err)
		return
	}

	Success(w, nil)
}
