package application

import (
	"github.com/G-Fuchter/home24-assignment/internal/domain/model"
	"github.com/G-Fuchter/home24-assignment/internal/ports"
)

// Application service seems redundant on such a small project, but it will help with scaling in the future
type Service struct {
	domainService ports.Service
}

func (s *Service) GenerateWebPageReport(location string) (model.WebPageReport, error) {
	return s.domainService.GenerateWebPageReport(location)
}
