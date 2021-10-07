package service

import (
	"goclass/internal/client"
	"goclass/internal/repository"
	"log"
)

type Service struct {
	client     client.Client
	repository repository.Repository
}

func New(c client.Client, r repository.Repository) Service {
	return Service{
		client:     c,
		repository: r,
	}
}

func (s Service) Refresh() {
	posts, err := s.client.Get()
	if err != nil {
		log.Fatalf("cannot get the posts: %v", err)
	}

	err = s.repository.Store(posts)
	if err != nil {
		return
	}
}
