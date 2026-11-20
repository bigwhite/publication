package knowledgefs

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func newTestStore(t *testing.T) (*Store, string) {
	t.Helper()

	root := t.TempDir()

	files := map[string]string{
		"INDEX.md": "# Knowledge Index\n",
		"engineering/timeout.md": `
# Timeout Policy

HTTP timeout and retry policy.
`,
		"incidents/order-timeout.md": `
# Order Incident

Timeout caused connection pool exhaustion.
`,
		"operations/monitoring.md": `
# Monitoring

Monitor request latency and error rate.
`,
		"notes.txt": "not a markdown file",
	}

	for name, content := range files {
		path := filepath.Join(root, filepath.FromSlash(name))

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	store, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}

	return store, root
}

func TestList(t *testing.T) {
	store, _ := newTestStore(t)

	docs, err := store.List()
	if err != nil {
		t.Fatal(err)
	}

	var got []string
	for _, doc := range docs {
		got = append(got, doc.Path)
	}

	want := []string{
		"INDEX.md",
		"engineering/timeout.md",
		"incidents/order-timeout.md",
		"operations/monitoring.md",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestRead(t *testing.T) {
	store, _ := newTestStore(t)

	doc, err := store.Read("engineering/timeout.md")
	if err != nil {
		t.Fatal(err)
	}

	if doc.Path != "engineering/timeout.md" {
		t.Fatalf("unexpected path: %s", doc.Path)
	}

	if doc.Content == "" {
		t.Fatal("expected non-empty content")
	}
}

func TestReadRejectsPathTraversal(t *testing.T) {
	store, root := newTestStore(t)

	outside := filepath.Join(filepath.Dir(root), "outside.md")
	if err := os.WriteFile(outside, []byte("secret"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := store.Read("../outside.md")
	if err == nil {
		t.Fatal("expected path traversal to be rejected")
	}
}

func TestReadRejectsNonMarkdown(t *testing.T) {
	store, _ := newTestStore(t)

	_, err := store.Read("notes.txt")
	if err == nil {
		t.Fatal("expected non-markdown file to be rejected")
	}
}

func TestReadIndex(t *testing.T) {
	store, _ := newTestStore(t)

	content, err := store.ReadIndex()
	if err != nil {
		t.Fatal(err)
	}

	if content == "" {
		t.Fatal("expected non-empty index")
	}
}

func TestSearch(t *testing.T) {
	store, _ := newTestStore(t)

	results, err := store.Search("timeout connection", 2)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) == 0 {
		t.Fatal("expected search results")
	}

	if results[0].Document.Path != "incidents/order-timeout.md" {
		t.Fatalf(
			"expected incidents/order-timeout.md, got %s",
			results[0].Document.Path,
		)
	}
}

func TestSearchTopK(t *testing.T) {
	store, _ := newTestStore(t)

	results, err := store.Search("timeout", 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) > 1 {
		t.Fatalf("expected at most 1 result, got %d", len(results))
	}
}
