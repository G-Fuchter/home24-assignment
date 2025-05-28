package http

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Hostname string
	Port     string
}

type Server struct {
	srv *echo.Echo
	cfg Config
}

func NewServer(srv *echo.Echo, config Config) *Server {
	return &Server{
		srv: srv,
		cfg: config,
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

func (s *Server) Start() error {
	s.srv.Use(middleware.CORS())
	port := s.cfg.Port
	hostname := s.cfg.Hostname
	return s.srv.Start(fmt.Sprintf("%v:%v", hostname, port))
}
