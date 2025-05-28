package main

import (
	"fmt"

	"github.com/G-Fuchter/home24-assignment/internal/adapters/http"
	"github.com/G-Fuchter/home24-assignment/internal/adapters/http/handlers"
	"github.com/G-Fuchter/home24-assignment/internal/adapters/parser"
	"github.com/G-Fuchter/home24-assignment/internal/domain"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := http.Config{
		Hostname: "",
		Port:     "8080",
	}
	srv := http.NewServer(e, cfg)
	handlers := getHandlers()
	srv.AddHandlers(handlers)
	srv.EnableCORS()
	srv.EnableStaticWebsite()
	err := srv.Start()
	fmt.Print(err.Error())
}

func getHandlers() []http.Handler {
	webParser := parser.NewWebPageParser()
	service := domain.NewService(webParser)
	return []http.Handler{
		handlers.NewCreateWebPageReport(service),
	}
}
