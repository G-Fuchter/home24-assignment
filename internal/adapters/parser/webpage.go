package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

var ErrCouldNotLoadDocument error = errors.New("could not load document")
var ErrDocumentNotLoaded error = errors.New("document has not been loaded")
var ErrFailedQuerying error = errors.New("failed to query document")
var ErrElementNotFound error = errors.New("could not find element")
var ErrNoVersionFound error = errors.New("could not find document version")

var regexHostnameURL = regexp.MustCompile(`^(https?:\/\/)?([^/?#:]+)`)

type WebPageParser struct {
	document    *html.Node
	documentURL string
}

func NewWebPageParser() *WebPageParser {
	return &WebPageParser{
		documentURL: "",
	}
}

// DownloadPage implements the DocumentParser interface.
func (p *WebPageParser) DownloadDocument(location string) error {
	var err error
	p.documentURL = location
	p.document, err = htmlquery.LoadURL(location)
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrCouldNotLoadDocument, err)
	}
	return err
}

func (p *WebPageParser) FromString(content string, url string) error {
	p.documentURL = url
	var err error
	p.document, err = htmlquery.Parse(strings.NewReader(content))
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrCouldNotLoadDocument, err)
	}
	return err

}

// GetDocumentVersion implements the DocumentParser interface.
func (p *WebPageParser) GetDocumentVersion() (string, error) {
	if p.document == nil {
		return "", ErrDocumentNotLoaded
	}
	doctypeNode := p.document.FirstChild
	if doctypeNode.Type != html.DoctypeNode {
		return "", ErrNoVersionFound
	}

	numOfAttributes := len(doctypeNode.Attr)
	isVersionFive := numOfAttributes == 0
	if isVersionFive {
		return "5", nil
	}

	rx, err := regexp.Compile(`[Hh][Tt][Mm][Ll] ([0-9](\.[0-9][0-9]*)*)`)
	if err != nil {
		return "", err
	}
	rxResult := rx.FindAllStringSubmatch(doctypeNode.Attr[0].Val, -1)
	if len(rxResult) == 0 {
		return "", ErrNoVersionFound
	}

	version := rxResult[0][1]

	return version, nil
}

// GetTitle implements the DocumentParser interface.
func (p *WebPageParser) GetTitle() (string, error) {
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
	if p.document == nil {
		return 0, ErrDocumentNotLoaded
	}
	count := 0
	links, err := getAllLinks(p.document)
	if err != nil {
		return 0, err
	}
	for _, link := range links {
		isInternal, err := isLinkURLInternal(p.documentURL, link)
		if err != nil {
			return 0, err
		}
		if !isInternal {
			count++
		}
	}
	return count, nil
}

// GetInternalLinkCount implements the DocumentParser interface.
func (p *WebPageParser) GetInternalLinkCount() (int, error) {
	if p.document == nil {
		return 0, ErrDocumentNotLoaded
	}
	count := 0
	links, err := getAllLinks(p.document)
	if err != nil {
		return 0, err
	}
	for _, link := range links {
		isInternal, err := isLinkURLInternal(p.documentURL, link)
		if err != nil {
			return 0, err
		}
		if isInternal {
			count++
		}
	}
	return count, nil
}

func getAllLinks(doc *html.Node) ([]string, error) {
	links := []string{}
	linkNodes, err := htmlquery.QueryAll(doc, "//a")
	if err != nil {
		return nil, err
	}
	for _, node := range linkNodes {
		attributes := node.Attr
		for _, attr := range attributes {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}
	return links, nil
}

func isLinkURLInternal(pageURL string, linkURL string) (bool, error) {
	if len(linkURL) > 0 && linkURL[0] == '/' {
		return true, nil
	}
	linkHostname, err := getHostname(linkURL)
	if err != nil {
		return false, nil
	}
	pageHostname, err := getHostname(pageURL)
	if err != nil {
		return false, fmt.Errorf("the website's URL is not valid: %w", err)
	}
	return linkHostname == pageHostname, nil
}

// GetContainsLogin implements the DocumentParser interface.
func (p *WebPageParser) GetContainsLogin() (bool, error) {
	doc := p.document
	if doc == nil {
		return false, ErrDocumentNotLoaded
	}
	forms, err := htmlquery.QueryAll(doc, "//form")
	if err != nil {
		return false, err
	}

	for _, f := range forms {
		inputs, err := htmlquery.QueryAll(f, "//input")
		hasUserInput, hasPasswordInput := false, false
		if err != nil {
			return false, fmt.Errorf("could not retrieve login: %w", err)
		}
		for _, i := range inputs {
			for _, attr := range i.Attr {
				if attr.Key == "type" {
					if attr.Val == "text" || attr.Val == "email" {
						hasUserInput = true
					} else if attr.Val == "password" {
						hasPasswordInput = true
					}
				}
			}
		}
		if hasUserInput && hasPasswordInput {
			return true, nil
		}
	}
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

func getHostname(url string) (string, error) {
	match := regexHostnameURL.FindStringSubmatch(url)

	if len(match) >= 3 {
		// The base URL is in the first capturing group (index 1)
		baseURL := match[2]
		return baseURL, nil
	} else {
		return "", fmt.Errorf("URL: %s -> No hostname found", url)
	}
}
