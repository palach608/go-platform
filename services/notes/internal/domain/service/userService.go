package service

import (
	apiModel "github.com/SilverName608/go-notes/internal/api/model"
)

type UserService interface {
	Register(req *apiModel.RegisterRequest) (*apiModel.AuthResponse, error)
	Login(req *apiModel.LoginRequest) (*apiModel.AuthResponse, error)
}
