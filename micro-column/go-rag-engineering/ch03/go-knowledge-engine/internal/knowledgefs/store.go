package knowledgefs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Store 表示一个基于文件系统的知识库。
type Store struct {
	root string
}

// Document 表示一份知识文档。
type Document struct {
	Path    string
	Content string
}

// SearchResult 表示一次关键词查找结果。
type SearchResult struct {
	Document Document
	Score    int
}

// Open 打开一个知识库目录。
func Open(root string) (*Store, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve knowledge root: %w", err)
	}

	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return nil, fmt.Errorf("resolve knowledge root symlink: %w", err)
	}

	info, err := os.Stat(real)
	if err != nil {
		return nil, fmt.Errorf("stat knowledge root: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("knowledge root is not a directory: %s", real)
	}

	return &Store{root: real}, nil
}

// Root 返回知识库的绝对路径。
func (s *Store) Root() string {
	return s.root
}

// resolve 将相对路径解析为知识库内部的真实文件路径。
// 它同时检查路径穿越和符号链接逃逸。
func (s *Store) resolve(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("empty document path")
	}

	if filepath.IsAbs(name) {
		return "", fmt.Errorf("absolute path is not allowed")
	}

	clean := filepath.Clean(name)
	candidate := filepath.Join(s.root, clean)

	rel, err := filepath.Rel(s.root, candidate)
	if err != nil {
		return "", fmt.Errorf("resolve relative path: %w", err)
	}

	if rel == ".." || strings.HasPrefix(
		rel,
		".."+string(filepath.Separator),
	) {
		return "", fmt.Errorf("path outside knowledge root")
	}

	real, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return "", fmt.Errorf("resolve document path: %w", err)
	}

	realRel, err := filepath.Rel(s.root, real)
	if err != nil {
		return "", fmt.Errorf("check real document path: %w", err)
	}

	if realRel == ".." || strings.HasPrefix(
		realRel,
		".."+string(filepath.Separator),
	) {
		return "", fmt.Errorf("symlink outside knowledge root")
	}

	return real, nil
}

// Read 读取知识库中的一份 Markdown 文档。
func (s *Store) Read(name string) (Document, error) {
	if !strings.EqualFold(filepath.Ext(name), ".md") {
		return Document{}, fmt.Errorf("only markdown files are supported")
	}

	fullPath, err := s.resolve(name)
	if err != nil {
		return Document{}, err
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return Document{}, fmt.Errorf("stat document: %w", err)
	}

	if !info.Mode().IsRegular() {
		return Document{}, fmt.Errorf("not a regular file: %s", name)
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return Document{}, fmt.Errorf("read document: %w", err)
	}

	return Document{
		Path:    filepath.ToSlash(filepath.Clean(name)),
		Content: string(content),
	}, nil
}

// List 递归列出知识库中的 Markdown 文档。
// 本方法只返回路径，不读取文档正文。
func (s *Store) List() ([]Document, error) {
	var docs []Document

	err := filepath.WalkDir(s.root, func(
		path string,
		entry fs.DirEntry,
		walkErr error,
	) error {
		if walkErr != nil {
			return walkErr
		}

		// 不跟随符号链接，也不将符号链接作为知识文件。
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}

		if entry.IsDir() {
			return nil
		}

		if !strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			return nil
		}

		rel, err := filepath.Rel(s.root, path)
		if err != nil {
			return err
		}

		docs = append(docs, Document{
			Path: filepath.ToSlash(rel),
		})

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk knowledge directory: %w", err)
	}

	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Path < docs[j].Path
	})

	return docs, nil
}

// ReadIndex 读取知识库根目录中的 INDEX.md。
func (s *Store) ReadIndex() (string, error) {
	doc, err := s.Read("INDEX.md")
	if err != nil {
		return "", err
	}

	return doc.Content, nil
}

// Search 使用命中的不同关键词数量进行简单检索。
// 每个关键词最多贡献 1 分，避免单个关键词的重复出现
// 掩盖文档对其他查询关键词的匹配情况。
// 这是教学用实现，不是生产级全文检索。
func (s *Store) Search(
	query string,
	topK int,
) ([]SearchResult, error) {
	if topK <= 0 {
		return nil, nil
	}

	// 将查询拆分为小写词项，并去除重复词项。
	rawTerms := strings.Fields(strings.ToLower(query))
	if len(rawTerms) == 0 {
		return nil, nil
	}

	termSet := make(map[string]struct{}, len(rawTerms))
	terms := make([]string, 0, len(rawTerms))

	for _, term := range rawTerms {
		if _, exists := termSet[term]; exists {
			continue
		}

		termSet[term] = struct{}{}
		terms = append(terms, term)
	}

	docs, err := s.List()
	if err != nil {
		return nil, err
	}

	var results []SearchResult

	for _, item := range docs {
		// INDEX.md 只负责导航，不参与正文搜索。
		if strings.EqualFold(filepath.Base(item.Path), "INDEX.md") {
			continue
		}

		doc, err := s.Read(item.Path)
		if err != nil {
			return nil, err
		}

		content := strings.ToLower(doc.Content)
		score := 0

		// 每个不同的查询词最多贡献 1 分。
		for _, term := range terms {
			if strings.Contains(content, term) {
				score++
			}
		}

		if score == 0 {
			continue
		}

		results = append(results, SearchResult{
			Document: doc,
			Score:    score,
		})
	}

	// 得分越高，说明命中的不同关键词越多。
	// 得分相同时，按照路径字典序排序，保证结果稳定。
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Document.Path <
				results[j].Document.Path
		}

		return results[i].Score > results[j].Score
	})

	if len(results) > topK {
		results = results[:topK]
	}

	return results, nil
}
