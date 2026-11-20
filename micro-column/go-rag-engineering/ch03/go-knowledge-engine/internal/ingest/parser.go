package ingest

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
)

type MarkdownParser struct{}

func (MarkdownParser) Parse(
	path string,
	raw []byte,
) (ParsedDocument, error) {
	content := normalizeText(string(raw))

	var headings []string
	var title string

	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !strings.HasPrefix(line, "#") {
			continue
		}

		heading := strings.TrimSpace(strings.TrimLeft(line, "#"))
		if heading == "" {
			continue
		}

		headings = append(headings, heading)

		if title == "" {
			title = heading
		}
	}

	if err := scanner.Err(); err != nil {
		return ParsedDocument{}, err
	}

	if title == "" {
		title = filenameTitle(path)
	}

	return ParsedDocument{
		Title:    title,
		Headings: headings,
		Content:  content,
	}, nil
}

type TextParser struct{}

func (TextParser) Parse(
	path string,
	raw []byte,
) (ParsedDocument, error) {
	content := normalizeText(string(raw))

	return ParsedDocument{
		Title:   filenameTitle(path),
		Content: content,
	}, nil
}

type CodeParser struct{}

func (CodeParser) Parse(
	path string,
	raw []byte,
) (ParsedDocument, error) {
	language := strings.TrimPrefix(filepath.Ext(path), ".")

	content := "```" + language + "\n" +
		string(raw) +
		"\n```"

	return ParsedDocument{
		Title:   filenameTitle(path),
		Content: content,
	}, nil
}

func parserFor(path string) (Parser, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".markdown":
		return MarkdownParser{}, nil
	case ".txt":
		return TextParser{}, nil
	case ".go", ".java", ".py", ".js", ".ts", ".rs":
		return CodeParser{}, nil
	default:
		return nil, fmt.Errorf(
			"unsupported format: %s",
			filepath.Ext(path),
		)
	}
}

func normalizeText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	return strings.TrimSpace(text) + "\n"
}

func filenameTitle(path string) string {
	name := filepath.Base(path)
	ext := filepath.Ext(name)
	return strings.TrimSuffix(name, ext)
}
