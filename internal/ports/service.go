package ports

import "github.com/G-Fuchter/home24-assignment/internal/domain/model"

type Service interface {
	GenerateWebPageReport(location string) (model.WebPageReport, error)
}
