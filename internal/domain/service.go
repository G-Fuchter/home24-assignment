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

	if err := s.parser.DownloadDocument(location); err != nil {
		return model.WebPageReport{}, fmt.Errorf("%w: %w", ErrInvlidPage, err)
	}

	documentVersion, err := s.parser.GetDocumentVersion()
	if err != nil {
		return model.WebPageReport{}, fmt.Errorf("failed to get document version: %w", err)
	}

	title, err := s.parser.GetTitle()
	if err != nil {
		return model.WebPageReport{}, fmt.Errorf("failed to get title: %w", err)
	}

	externalLinkCount, err := s.parser.GetExternalLinkCount()
	if err != nil {
		return model.WebPageReport{}, fmt.Errorf("failed to get external link count: %w", err)
	}

	internalLinkCount, err := s.parser.GetInternalLinkCount()
	if err != nil {
		return model.WebPageReport{}, fmt.Errorf("failed to get internal link count: %w", err)
	}

	containsLogin, err := s.parser.GetContainsLogin()
	if err != nil {
		return model.WebPageReport{}, fmt.Errorf("failed to check login form: %w", err)
	}

	headerOneCount, err := s.parser.GetHeaderOneCount()
	if err != nil {
		return model.WebPageReport{}, fmt.Errorf("failed to get H1 count: %w", err)
	}

	headerTwoCount, err := s.parser.GetHeaderTwoCount()
	if err != nil {
		return model.WebPageReport{}, fmt.Errorf("failed to get H2 count: %w", err)
	}

	headerThreeCount, err := s.parser.GetHeaderThreeCount()
	if err != nil {
		return model.WebPageReport{}, fmt.Errorf("failed to get H3 count: %w", err)
	}

	headerFourCount, err := s.parser.GetHeaderFourCount()
	if err != nil {
		return model.WebPageReport{}, fmt.Errorf("failed to get H4 count: %w", err)
	}

	headerFiveCount, err := s.parser.GetHeaderFiveCount()
	if err != nil {
		return model.WebPageReport{}, fmt.Errorf("failed to get H5 count: %w", err)
	}

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
