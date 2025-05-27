package parser_test

import (
	"errors"
	"testing"

	"github.com/G-Fuchter/home24-assignment/internal/adapters/parser"
)

func TestNewMyDocumentParser(t *testing.T) {
	pr := parser.NewWebPageParser()

	if pr == nil {
		t.Fatal("NewMyDocumentParser() returned nil")
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "valid HTML",
			content:     "<html><head><title>Test</title></head><body><h1>Header</h1></body></html>",
			expectError: false,
		},
		{
			name:        "minimal HTML",
			content:     "<html></html>",
			expectError: false,
		},
		{
			name:        "empty string",
			content:     "",
			expectError: false, // htmlquery.Parse handles empty strings gracefully
		},
		{
			name:        "malformed HTML",
			content:     "<html><head><title>Test</title><body><h1>Header</h1></body>", // missing closing tags
			expectError: false,                                                         // HTML parser is forgiving
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prs := parser.NewWebPageParser()
			err := prs.FromString(tt.content)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGetDocumentVersion(t *testing.T) {

	// Test without loading document
	t.Run("should return error when no document is loaded", func(t *testing.T) {
		prsr := parser.NewWebPageParser()
		_, err := prsr.GetDocumentVersion()
		if err != parser.ErrDocumentNotLoaded {
			t.Fatalf("Expected document not loaded error to be return: %v", err)
		}
	})

	t.Run("should return error when there is no doctype", func(t *testing.T) {
		prsr := parser.NewWebPageParser()
		err := prsr.FromString("<html><head><title></title></head></html>")
		if err != nil {
			t.Fatalf("Failed to load document: %v", err)
		}

		_, err = prsr.GetDocumentVersion()
		if err != parser.ErrNoVersionFound {
			t.Fatalf("Expected error was not returned: %v", err)
		}

	})
	t.Run("should return correct html version", func(t *testing.T) {
		tests := []struct {
			name            string
			html            string
			expectedVersion string
		}{
			{
				name:            "version 5",
				html:            "<!DOCTYPE html><html><head><title>Test</title></head></html>",
				expectedVersion: "5",
			},
			{
				name:            "version 4.01",
				html:            "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\"><html><head><title>Test</title></head></html>",
				expectedVersion: "4.01",
			},
			{
				name:            "version 3.2",
				html:            "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 3.2 Final//EN\"><html><head><title>Test</title></head></html>",
				expectedVersion: "3.2",
			},
			{
				name:            "version 2.0",
				html:            "<!DOCTYPE html PUBLIC \"-//IETF//DTD HTML 2.0//EN\"><html><head><title>Test</title></head></html>",
				expectedVersion: "2.0",
			},
		}

		prsr := parser.NewWebPageParser()
		for _, tcase := range tests {
			t.Run(tcase.name, func(t *testing.T) {
				err := prsr.FromString(tcase.html)
				if err != nil {
					t.Fatalf("Failed to load document: %v", err)
				}

				actual, err := prsr.GetDocumentVersion()
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				if actual != tcase.expectedVersion {
					t.Fatalf("Expected %v version and returned %v", tcase.expectedVersion, actual)
				}
			})
		}
	})
}

func TestGetTitle(t *testing.T) {
	tests := []struct {
		name          string
		html          string
		expectedTitle string
		expectError   bool
		errorType     error
	}{
		{
			name:          "valid title",
			html:          "<html><head><title>Test Page</title></head></html>",
			expectedTitle: "Test Page",
			expectError:   false,
		},
		{
			name:          "empty title",
			html:          "<html><head><title></title></head></html>",
			expectedTitle: "",
			expectError:   false,
		},
		{
			name:        "no title element",
			html:        "<html><head></head></html>",
			expectError: true,
			errorType:   parser.ErrElementNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := parser.NewWebPageParser()
			err := parser.FromString(tt.html)
			if err != nil {
				t.Fatalf("Failed to load document: %v", err)
			}

			title, err := parser.GetTitle()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("Expected error type %v, got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if title != tt.expectedTitle {
					t.Errorf("Expected title %q, got %q", tt.expectedTitle, title)
				}
			}
		})
	}
}

func TestGetTitle_DocumentNotLoaded(t *testing.T) {
	prsr := parser.NewWebPageParser()

	title, err := prsr.GetTitle()

	if err == nil {
		t.Error("Expected error but got none")
	}

	if !errors.Is(err, parser.ErrDocumentNotLoaded) {
		t.Errorf("Expected ErrDocumentNotLoaded, got %v", err)
	}

	if title != "" {
		t.Errorf("Expected empty title, got %q", title)
	}
}

func TestGetExternalLinkCount(t *testing.T) {
	parser := parser.NewWebPageParser()

	count, err := parser.GetExternalLinkCount()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}
}

func TestGetInternalLinkCount(t *testing.T) {
	parser := parser.NewWebPageParser()

	count, err := parser.GetInternalLinkCount()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}
}

func TestGetContainsLogin(t *testing.T) {
	parser := parser.NewWebPageParser()

	contains, err := parser.GetContainsLogin()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if contains != false {
		t.Errorf("Expected false, got %t", contains)
	}
}

func TestHeaderCounts(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		h1Count int
		h2Count int
		h3Count int
		h4Count int
		h5Count int
		h6Count int
	}{
		{
			name:    "no headers",
			html:    "<html><body><p>No headers here</p></body></html>",
			h1Count: 0, h2Count: 0, h3Count: 0, h4Count: 0, h5Count: 0, h6Count: 0,
		},
		{
			name:    "single header of each type",
			html:    "<html><body><h1>H1</h1><h2>H2</h2><h3>H3</h3><h4>H4</h4><h5>H5</h5><h6>H6</h6></body></html>",
			h1Count: 1, h2Count: 1, h3Count: 1, h4Count: 1, h5Count: 1, h6Count: 1,
		},
		{
			name:    "multiple headers of same type",
			html:    "<html><body><h1>First</h1><h1>Second</h1><h2>H2</h2></body></html>",
			h1Count: 2, h2Count: 1, h3Count: 0, h4Count: 0, h5Count: 0, h6Count: 0,
		},
		{
			name:    "nested headers",
			html:    "<html><body><div><h1>Title</h1><section><h2>Section</h2></section></div></body></html>",
			h1Count: 1, h2Count: 1, h3Count: 0, h4Count: 0, h5Count: 0, h6Count: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := parser.NewWebPageParser()
			err := parser.FromString(tt.html)
			if err != nil {
				t.Fatalf("Failed to load document: %v", err)
			}

			// Test H1 count
			count, err := parser.GetHeaderOneCount()
			if err != nil {
				t.Errorf("GetHeaderOneCount() error: %v", err)
			}
			if count != tt.h1Count {
				t.Errorf("Expected H1 count %d, got %d", tt.h1Count, count)
			}

			// Test H2 count
			count, err = parser.GetHeaderTwoCount()
			if err != nil {
				t.Errorf("GetHeaderTwoCount() error: %v", err)
			}
			if count != tt.h2Count {
				t.Errorf("Expected H2 count %d, got %d", tt.h2Count, count)
			}

			// Test H3 count
			count, err = parser.GetHeaderThreeCount()
			if err != nil {
				t.Errorf("GetHeaderThreeCount() error: %v", err)
			}
			if count != tt.h3Count {
				t.Errorf("Expected H3 count %d, got %d", tt.h3Count, count)
			}

			// Test H4 count
			count, err = parser.GetHeaderFourCount()
			if err != nil {
				t.Errorf("GetHeaderFourCount() error: %v", err)
			}
			if count != tt.h4Count {
				t.Errorf("Expected H4 count %d, got %d", tt.h4Count, count)
			}

			// Test H5 count
			count, err = parser.GetHeaderFiveCount()
			if err != nil {
				t.Errorf("GetHeaderFiveCount() error: %v", err)
			}
			if count != tt.h5Count {
				t.Errorf("Expected H5 count %d, got %d", tt.h5Count, count)
			}

			// Test H6 count
			count, err = parser.GetHeaderSixCount()
			if err != nil {
				t.Errorf("GetHeaderSixCount() error: %v", err)
			}
			if count != tt.h6Count {
				t.Errorf("Expected H6 count %d, got %d", tt.h6Count, count)
			}
		})
	}
}

func TestHeaderCounts_DocumentNotLoaded(t *testing.T) {
	prsr := parser.NewWebPageParser()

	tests := []struct {
		name     string
		testFunc func() (int, error)
	}{
		{"GetHeaderOneCount", prsr.GetHeaderOneCount},
		{"GetHeaderTwoCount", prsr.GetHeaderTwoCount},
		{"GetHeaderThreeCount", prsr.GetHeaderThreeCount},
		{"GetHeaderFourCount", prsr.GetHeaderFourCount},
		{"GetHeaderFiveCount", prsr.GetHeaderFiveCount},
		{"GetHeaderSixCount", prsr.GetHeaderSixCount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := tt.testFunc()

			if err == nil {
				t.Error("Expected error but got none")
			}

			if !errors.Is(err, parser.ErrDocumentNotLoaded) {
				t.Errorf("Expected ErrDocumentNotLoaded, got %v", err)
			}

			if count != 0 {
				t.Errorf("Expected count 0, got %d", count)
			}
		})
	}
}
