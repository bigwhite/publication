package parser

import (
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if !opts.IncludeComments {
		t.Error("Expected IncludeComments to be true")
	}

	if !opts.IncludeMetadata {
		t.Error("Expected IncludeMetadata to be true")
	}

	if !opts.IncludeTimestamps {
		t.Error("Expected IncludeTimestamps to be true")
	}

	if !opts.IncludeUserLinks {
		t.Error("Expected IncludeUserLinks to be true")
	}

	if !opts.EmojisEnabled {
		t.Error("Expected EmojisEnabled to be true")
	}

	if !opts.PreserveLineBreaks {
		t.Error("Expected PreserveLineBreaks to be true")
	}
}

func TestNewParser(t *testing.T) {
	tests := []struct {
		name string
		opts *Options
	}{
		{
			name: "with nil options",
			opts: nil,
		},
		{
			name: "with custom options",
			opts: &Options{
				IncludeComments: false,
				IncludeMetadata: true,
			},
		},
		{
			name: "with default options",
			opts: DefaultOptions(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.opts)
			if parser == nil {
				t.Error("NewParser() returned nil")
			}
			if parser.options == nil {
				t.Error("NewParser().options is nil")
			}
		})
	}
}

func TestMarkdownDocument(t *testing.T) {
	doc := &MarkdownDocument{
		Title:   "Test Document",
		Content: "# Test\nThis is a test document.",
		Metadata: map[string]string{
			"author": "testuser",
			"date":   "2023-01-01",
		},
	}

	if doc.Title != "Test Document" {
		t.Errorf("Expected title 'Test Document', got %s", doc.Title)
	}

	if doc.Content != "# Test\nThis is a test document." {
		t.Errorf("Expected content mismatch")
	}

	if doc.Metadata["author"] != "testuser" {
		t.Errorf("Expected author 'testuser', got %s", doc.Metadata["author"])
	}
}

func TestProcessingError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		code    string
		details string
		want    string
	}{
		{
			name:    "basic error",
			message: "test error",
			code:    "TEST_ERR",
			details: "",
			want:    "test error",
		},
		{
			name:    "error with details",
			message: "test error",
			code:    "TEST_ERR",
			details: "additional details",
			want:    "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewProcessingError(tt.message, tt.code, tt.details)
			if err == nil {
				t.Error("NewProcessingError() returned nil")
			}
			if err.Error() != tt.want {
				t.Errorf("ProcessingError.Error() = %v, want %v", err.Error(), tt.want)
			}
			if err.Code != tt.code {
				t.Errorf("ProcessingError.Code = %v, want %v", err.Code, tt.code)
			}
		})
	}
}