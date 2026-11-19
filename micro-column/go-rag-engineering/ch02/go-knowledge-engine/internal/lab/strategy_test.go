package lab

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestFullContext(t *testing.T) {
	docs := []Document{
		{Name: "a.md", Content: "architecture"},
		{Name: "b.md", Content: "timeout"},
	}

	result := FullContext(docs)

	if result.Mode != "full" {
		t.Fatalf("unexpected mode: %s", result.Mode)
	}

	if len(result.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(result.Files))
	}

	if !strings.Contains(result.Context, "architecture") ||
		!strings.Contains(result.Context, "timeout") {
		t.Fatal("full context should contain all documents")
	}
}

func TestSearch(t *testing.T) {
	docs := []Document{
		{
			Name:    "architecture.md",
			Content: "service architecture and service boundary",
		},
		{
			Name:    "timeout.md",
			Content: "http timeout retry timeout",
		},
		{
			Name:    "incident.md",
			Content: "timeout retry connection pool",
		},
	}

	// 增加 connection 和 pool 两个关键词，
	// 使历史故障文档的匹配得分最高。
	result := Search(docs, "timeout retry connection pool", 1)

	if result.Mode != "search" {
		t.Fatalf("unexpected mode: %s", result.Mode)
	}

	expected := []string{"incident.md"}
	if !reflect.DeepEqual(result.Files, expected) {
		t.Fatalf("expected %v, got %v", expected, result.Files)
	}
}

func TestNavigate(t *testing.T) {
	root := t.TempDir()

	index := "# Knowledge Index\n- timeout.md"
	if err := os.WriteFile(
		filepath.Join(root, "INDEX.md"),
		[]byte(index),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	content := "# Timeout\nHTTP timeout is 2 seconds."
	if err := os.WriteFile(
		filepath.Join(root, "timeout.md"),
		[]byte(content),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	result, err := Navigate(root, []string{"timeout.md"})
	if err != nil {
		t.Fatal(err)
	}

	if result.Mode != "navigate" {
		t.Fatalf("unexpected mode: %s", result.Mode)
	}

	if !strings.Contains(result.Context, "Knowledge Index") {
		t.Fatal("context should contain index")
	}

	if !strings.Contains(result.Context, "HTTP timeout") {
		t.Fatal("context should contain selected document")
	}
}

func TestNavigateRejectsPathTraversal(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(
		filepath.Join(root, "INDEX.md"),
		[]byte("# Index"),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	_, err := Navigate(root, []string{"../secret.md"})
	if err == nil {
		t.Fatal("expected path traversal to be rejected")
	}
}
