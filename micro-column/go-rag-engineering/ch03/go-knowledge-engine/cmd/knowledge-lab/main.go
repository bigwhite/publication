package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"example.com/go-knowledge-engine/internal/ingest"
	"example.com/go-knowledge-engine/internal/lab"
)

func main() {
	mode := flag.String("mode", "full", "full, navigate, or search")
	knowledgeDir := flag.String("knowledge", "./knowledge", "knowledge directory")
	query := flag.String("query", "", "search query")
	files := flag.String("files", "", "comma-separated files for navigation")
	topK := flag.Int("topk", 2, "number of search results")

	ingestMode := flag.Bool("ingest", false, "run knowledge ingestion")
	sourceDir := flag.String("source", "./sources", "raw knowledge source directory")
	manifestPath := flag.String(
		"manifest",
		"./knowledge/.manifest.json",
		"ingestion manifest path",
	)

	flag.Parse()

	if *ingestMode {
		stats, err := ingest.Run(ingest.Config{
			SourceRoot:    *sourceDir,
			KnowledgeRoot: *knowledgeDir,
			ManifestPath:  *manifestPath,
		})
		if err != nil {
			fail(err)
		}

		fmt.Printf("Scanned: %d\n", stats.Scanned)
		fmt.Printf("Added:   %d\n", stats.Added)
		fmt.Printf("Updated: %d\n", stats.Updated)
		fmt.Printf("Skipped: %d\n", stats.Skipped)
		fmt.Printf("Deleted: %d\n", stats.Deleted)
		fmt.Printf("Failed:  %d\n", stats.Failed)
		return
	}

	docs, err := lab.LoadCorpus(*knowledgeDir)
	if err != nil {
		fail(err)
	}

	var result lab.Result

	switch *mode {
	case "full":
		result = lab.FullContext(docs)

	case "navigate":
		var selected []string

		if strings.TrimSpace(*files) != "" {
			for _, name := range strings.Split(*files, ",") {
				name = strings.TrimSpace(name)
				if name != "" {
					selected = append(selected, name)
				}
			}
		}

		result, err = lab.Navigate(*knowledgeDir, selected)
		if err != nil {
			fail(err)
		}

	case "search":
		result = lab.Search(docs, *query, *topK)

	default:
		fail(fmt.Errorf("unsupported mode: %s", *mode))
	}

	fmt.Printf("Mode: %s\n", result.Mode)
	fmt.Printf("Files: %s\n", strings.Join(result.Files, ", "))
	fmt.Printf("Context chars: %d\n", utf8.RuneCountInString(result.Context))
	fmt.Printf("Approx tokens: %d\n", approximateTokens(result.Context))
	fmt.Println("\n========== CONTEXT ==========")
	fmt.Println(result.Context)
}

func approximateTokens(text string) int {
	return (utf8.RuneCountInString(text) + 3) / 4
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
