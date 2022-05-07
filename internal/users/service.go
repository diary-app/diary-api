package users

import (
	"diary-api/internal/auth"
)

type Service interface {
	GetSharingKey(id string) string
	Login(login, password string) auth.LoginResult
	Register(login, password string) auth.LoginResult
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo Repo
}

func (s *service) GetSharingKey(id string) string {
	//TODO implement me
	panic("implement me")
}

func (s *service) Login(login, password string) auth.LoginResult {
	//TODO implement me
	panic("implement me")
}

func (s *service) Register(login, password string) auth.LoginResult {
	//TODO implement me
	panic("implement me")
}
