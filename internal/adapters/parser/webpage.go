package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

var ErrCouldNotLoadDocument error = errors.New("could not load document")
var ErrDocumentNotLoaded error = errors.New("document has not been loaded")
var ErrFailedQuerying error = errors.New("failed to query document")
var ErrElementNotFound error = errors.New("could not find element")

type WebPageParser struct {
	document    *html.Node
	documentURL string
}

func NewMyDocumentParser() *WebPageParser {
	return &WebPageParser{
		documentURL: "",
	}
}

// DownloadPage implements the DocumentParser interface.
func (p *WebPageParser) DownloadPage(location string) error {
	var err error
	p.documentURL = location
	p.document, err = htmlquery.LoadURL("http://example.com/")
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrCouldNotLoadDocument, err)
	}
	return err
}

func (p *WebPageParser) FromString(content string) error {
	var err error
	p.document, err = htmlquery.Parse(strings.NewReader(content))
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrCouldNotLoadDocument, err)
	}
	return err

}

// GetDocumentVersion implements the DocumentParser interface.
func (p *WebPageParser) GetDocumentVersion() (string, error) {
	// TODO: Add logic to retrieve the document version.
	return "", nil // Empty string for now
}

// GetTitle implements the DocumentParser interface.
func (p *WebPageParser) GetTitle() (string, error) {
	// TODO: Add logic to retrieve the document title.
	doc := p.document
	if doc == nil {
		return "", ErrDocumentNotLoaded
	}
	title, err := htmlquery.Query(doc, "//title")
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedQuerying, err)
	}
	if title == nil {
		return "", ErrElementNotFound
	}
	txtNode := title.FirstChild
	if txtNode == nil {
		return "", nil
	} else {
		return txtNode.Data, nil
	}
}

// GetExternalLinkCount implements the DocumentParser interface.
func (p *WebPageParser) GetExternalLinkCount() (int, error) {
	// TODO: Add logic to count external links.
	return 0, nil // Zero for now
}

// GetInternalLinkCount implements the DocumentParser interface.
func (p *WebPageParser) GetInternalLinkCount() (int, error) {
	// TODO: Add logic to count internal links.
	return 0, nil // Zero for now
}

// GetContainsLogin implements the DocumentParser interface.
func (p *WebPageParser) GetContainsLogin() (bool, error) {
	// TODO: Add logic to determine if the page contains a login form/button.
	return false, nil // False for now
}

// GetHeaderOneCount implements the DocumentParser interface.
func (p *WebPageParser) GetHeaderOneCount() (int, error) {
	count, err := p.getElementCount("//h1")
	return count, err
}

// GetHeaderTwoCount implements the DocumentParser interface.
func (p *WebPageParser) GetHeaderTwoCount() (int, error) {
	count, err := p.getElementCount("//h2")
	return count, err
}

// GetHeaderThreeCount implements the DocumentParser interface.
func (p *WebPageParser) GetHeaderThreeCount() (int, error) {
	count, err := p.getElementCount("//h3")
	return count, err
}

// GetHeaderFourCount implements the DocumentParser interface.
func (p *WebPageParser) GetHeaderFourCount() (int, error) {
	count, err := p.getElementCount("//h4")
	return count, err
}

func (p *WebPageParser) GetHeaderFiveCount() (int, error) {
	count, err := p.getElementCount("//h5")
	return count, err
}

func (p *WebPageParser) GetHeaderSixCount() (int, error) {
	count, err := p.getElementCount("//h6")
	return count, err
}

func (p *WebPageParser) getElementCount(element string) (int, error) {

	expr := fmt.Sprintf("count(%v)", element)
	comp, err := xpath.Compile(expr)
	if err != nil {
		return 0, err
	}
	doc := p.document
	if doc == nil {
		return 0, ErrDocumentNotLoaded
	}
	count := comp.Evaluate(htmlquery.CreateXPathNavigator(doc)).(float64)
	return int(count), nil

}
