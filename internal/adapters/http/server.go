package http

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type Config struct {
	hostname string
	port     string
}

type Server struct {
	srv *echo.Echo
	cfg Config
}

func NewServer(srv *echo.Echo) *Server {
	return &Server{
		srv: srv,
	}
}

func (s *Server) AddHandlers(h []Handler) error {
	for _, value := range h {
		switch value.GetMethod() {
		case Get:
			s.srv.GET(
				value.GetEndpoint(),
				func(c echo.Context) error { return value.Handle(c) },
			)
			return nil
		case Post:
			s.srv.POST(
				value.GetEndpoint(),
				func(c echo.Context) error { return value.Handle(c) },
			)
			return nil
		case Put:
			s.srv.PUT(
				value.GetEndpoint(),
				func(c echo.Context) error { return value.Handle(c) },
			)
			return nil
		case Delete:
			s.srv.DELETE(
				value.GetEndpoint(),
				func(c echo.Context) error { return value.Handle(c) },
			)
		case Patch:
			s.srv.PATCH(
				value.GetEndpoint(),
				func(c echo.Context) error { return value.Handle(c) },
			)
		default:
			return fmt.Errorf(
				"Handler for endpoint: %v has an invalid Method %v",
				value.GetEndpoint(),
				value.GetMethod(),
			)
		}
	}
	return nil
}

func (s *Server) Start(h []Handler) error {
	port := s.cfg.port
	hostname := s.cfg.hostname
	return s.srv.Start(fmt.Sprintf("%v:%v", port, hostname))
}
