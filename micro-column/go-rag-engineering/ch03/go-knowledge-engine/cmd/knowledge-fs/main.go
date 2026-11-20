package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"example.com/go-knowledge-engine/internal/knowledgefs"
)

func main() {
	root := flag.String(
		"root",
		"./knowledge",
		"knowledge directory",
	)

	action := flag.String(
		"action",
		"list",
		"list, index, read, or search",
	)

	path := flag.String(
		"path",
		"",
		"relative markdown path",
	)

	query := flag.String(
		"query",
		"",
		"keyword query",
	)

	topK := flag.Int(
		"topk",
		5,
		"maximum search results",
	)

	flag.Parse()

	store, err := knowledgefs.Open(*root)
	if err != nil {
		fail(err)
	}

	switch strings.ToLower(*action) {
	case "list":
		docs, err := store.List()
		if err != nil {
			fail(err)
		}

		for _, doc := range docs {
			fmt.Println(doc.Path)
		}

	case "index":
		content, err := store.ReadIndex()
		if err != nil {
			fail(err)
		}

		fmt.Println(content)

	case "read":
		if strings.TrimSpace(*path) == "" {
			fail(fmt.Errorf("-path is required"))
		}

		doc, err := store.Read(*path)
		if err != nil {
			fail(err)
		}

		fmt.Printf("File: %s\n\n", doc.Path)
		fmt.Println(doc.Content)

	case "search":
		if strings.TrimSpace(*query) == "" {
			fail(fmt.Errorf("-query is required"))
		}

		results, err := store.Search(*query, *topK)
		if err != nil {
			fail(err)
		}

		for _, result := range results {
			fmt.Printf(
				"[%d] %s\n",
				result.Score,
				result.Document.Path,
			)
		}

	default:
		fail(fmt.Errorf("unsupported action: %s", *action))
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	log.Fatal()
}
