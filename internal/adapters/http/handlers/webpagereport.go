package handlers

import (
	httpgo "net/http"

	"github.com/G-Fuchter/home24-assignment/internal/adapters/http"
	"github.com/G-Fuchter/home24-assignment/internal/ports"
)

type PostWebPageReportRequestBody struct {
	URL string `json:"url"`
}

type PostWebPageReportResponseBody struct {
	DocumentVersion   string `json:"documentVersion"`
	Title             string `json:"title"`
	ExternalLinkCount int    `json:"externalLinkCount"`
	InternalLinkCount int    `json:"internalLinkCount"`
	ContainsLogin     bool   `json:"containsLogin"`
	HeaderOneCount    int    `json:"headerOneCount"`
	HeaderTwoCount    int    `json:"headerTwoCount"`
	HeaderThreeCount  int    `json:"headerThreeCount"`
	HeaderFourCount   int    `json:"headerFourCount"`
	HeaderFiveCount   int    `json:"headerFiveCount"`
	HeaderSixCount    int    `json:"headerSixCount"`
}

type CreateWebPageReport struct {
	webpageReportService ports.Service
}

func NewCreateWebPageReport(webpageReportService ports.Service) *CreateWebPageReport {
	return &CreateWebPageReport{
		webpageReportService: webpageReportService,
	}
}

func (h *CreateWebPageReport) GetMethod() http.Method {
	return http.Post
}

func (h *CreateWebPageReport) GetEndpoint() string {
	return "/reports/webpage"
}

func (h *CreateWebPageReport) Handle(c http.Context) error {
	var body PostWebPageReportRequestBody
	err := c.Bind(&body)
	if err != nil {
		return c.NoContent(httpgo.StatusBadRequest)
	}
	model, err := h.webpageReportService.GenerateWebPageReport(body.URL)
	if err != nil {
		return c.NoContent(httpgo.StatusInternalServerError)
	}
	resBody := &PostWebPageReportResponseBody{
		DocumentVersion:   model.DocumentVersion,
		Title:             model.Title,
		ExternalLinkCount: model.ExternalLinkCount,
		InternalLinkCount: model.InternalLinkCount,
		ContainsLogin:     model.ContainsLogin,
		HeaderOneCount:    model.HeaderOneCount,
		HeaderTwoCount:    model.HeaderTwoCount,
		HeaderThreeCount:  model.HeaderThreeCount,
		HeaderFourCount:   model.HeaderFourCount,
		HeaderFiveCount:   model.HeaderFiveCount,
		HeaderSixCount:    model.HeaderSixCount,
	}
	c.JSON(httpgo.StatusCreated, resBody)
	return nil
}
