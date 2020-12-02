package handler

import (
	service "github.com/mtanzim/event-driven-bookstore/bookstore-server/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s}
}
