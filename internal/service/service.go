package service

import (
	"github.com/dngrudin/notes-service/internal/api"
	"github.com/dngrudin/notes-service/internal/api/routes"
	"github.com/dngrudin/notes-service/internal/config"
	"github.com/dngrudin/notes-service/internal/storage"
)

type Service struct {
	config  *config.Config
	storage storage.Storage
}

func New(config *config.Config, storage storage.Storage) *Service {
	return &Service{config: config, storage: storage}
}

func (s *Service) Run() {
	rb := routes.NewBuilder(NewExecuter(s.storage))
	srv := api.NewServer(s.config.HTTPServer, rb.Build())
	srv.Start()
}
