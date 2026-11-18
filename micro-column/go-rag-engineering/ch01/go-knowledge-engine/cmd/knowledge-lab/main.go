package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"example.com/go-knowledge-engine/internal/lab"
)

func main() {
	mode := flag.String("mode", "full", "full, navigate, or search")
	knowledgeDir := flag.String("knowledge", "./knowledge", "knowledge directory")
	query := flag.String("query", "", "search query")
	files := flag.String("files", "", "comma-separated files for navigation")
	topK := flag.Int("topk", 2, "number of search results")
	flag.Parse()

	docs, err := lab.LoadCorpus(*knowledgeDir)
	if err != nil {
		fail(err)
	}

	if len(docs) == 0 {
		log.Fatalf("no knowledge documents found in %q", *knowledgeDir)
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
	// 仅用于实验观察，不代表任何模型的真实 Tokenizer。
	return (utf8.RuneCountInString(text) + 3) / 4
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
