package domain

import (
	"errors"
	"fmt"

	"github.com/G-Fuchter/home24-assignment/internal/domain/model"
	"github.com/G-Fuchter/home24-assignment/internal/ports"
)

var ErrInvlidPage = errors.New("URL is invalid or unreachable")

type Service struct {
	parser ports.DocumentParser
}

func NewService(p ports.DocumentParser) *Service {
	return &Service{
		parser: p,
	}
}

func (s *Service) GenerateWebPageReport(location string) (model.WebPageReport, error) {
	err := s.parser.DownloadPage(location)
	if err == nil {
		return model.WebPageReport{}, fmt.Errorf("%w: %w", ErrInvlidPage, err)
	}
	documentVersion := s.parser.GetDocumentVersion()
	title := s.parser.GetTitle()
	externalLinkCount := s.parser.GetExternalLinkCount()
	internalLinkCount := s.parser.GetInternalLinkCount()
	containsLogin := s.parser.GetContainsLogin()
	headerOneCount := s.parser.GetHeaderOneCount()
	headerTwoCount := s.parser.GetHeaderTwoCount()
	headerThreeCount := s.parser.GetHeaderThreeCount()
	headerFourCount := s.parser.GetHeaderFourCount()
	headerFiveCount := s.parser.GetHeaderFiveCount()

	return model.WebPageReport{
		DocumentVersion:   documentVersion,
		Title:             title,
		ExternalLinkCount: externalLinkCount,
		InternalLinkCount: internalLinkCount,
		ContainsLogin:     containsLogin,
		HeaderOneCount:    headerOneCount,
		HeaderTwoCount:    headerTwoCount,
		HeaderThreeCount:  headerThreeCount,
		HeaderFourCount:   headerFourCount,
		HeaderFiveCount:   headerFiveCount,
	}, err

}
