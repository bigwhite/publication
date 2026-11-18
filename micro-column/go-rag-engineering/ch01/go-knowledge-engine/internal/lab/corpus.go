package lab

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Document 表示知识库中的一份 Markdown 文档。
type Document struct {
	Name    string
	Content string
}

// LoadCorpus 加载目录下的 Markdown 文件。
func LoadCorpus(root string) ([]Document, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read knowledge directory: %w", err)
	}

	var docs []Document

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".md") || name == "INDEX.md" {
			continue
		}

		path := filepath.Join(root, name)
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}

		docs = append(docs, Document{
			Name:    name,
			Content: string(content),
		})
	}

	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Name < docs[j].Name
	})

	return docs, nil
}

// ReadIndex 读取知识库入口索引。
func ReadIndex(root string) (string, error) {
	path := filepath.Join(root, "INDEX.md")

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read index: %w", err)
	}

	return string(content), nil
}
