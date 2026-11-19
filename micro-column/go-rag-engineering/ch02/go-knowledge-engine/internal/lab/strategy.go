package lab

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

type Result struct {
	Mode    string
	Files   []string
	Context string
}

func FullContext(docs []Document) Result {
	var builder strings.Builder
	var files []string

	for _, doc := range docs {
		files = append(files, doc.Name)

		fmt.Fprintf(
			&builder,
			"\n\n--- FILE: %s ---\n%s",
			doc.Name,
			doc.Content,
		)
	}

	return Result{
		Mode:    "full",
		Files:   files,
		Context: strings.TrimSpace(builder.String()),
	}
}

func Navigate(root string, selected []string) (Result, error) {
	index, err := ReadIndex(root)
	if err != nil {
		return Result{}, err
	}

	var builder strings.Builder
	builder.WriteString("# Knowledge Index\n")
	builder.WriteString(index)

	var files []string

	for _, name := range selected {
		if filepath.Base(name) != name {
			return Result{}, fmt.Errorf("invalid document name: %s", name)
		}

		if !strings.HasSuffix(name, ".md") ||
			name == "INDEX.md" {
			return Result{}, fmt.Errorf("invalid knowledge file: %s", name)
		}

		path := filepath.Join(root, name)
		content, err := os.ReadFile(path)
		if err != nil {
			return Result{}, fmt.Errorf("read %s: %w", name, err)
		}

		files = append(files, name)

		fmt.Fprintf(
			&builder,
			"\n\n--- FILE: %s ---\n%s",
			name,
			string(content),
		)
	}

	return Result{
		Mode:    "navigate",
		Files:   files,
		Context: strings.TrimSpace(builder.String()),
	}, nil
}

// Search 使用简单的词项匹配，选择得分最高的文档。
// 中文暂时按连续汉字片段匹配，不是真正的中文分词。
func Search(docs []Document, query string, topK int) Result {
	if topK <= 0 {
		return Result{Mode: "search"}
	}

	queryTerms := tokenize(query)
	if len(queryTerms) == 0 {
		return Result{Mode: "search"}
	}

	type scoredDoc struct {
		doc   Document
		score int
	}

	var candidates []scoredDoc

	for _, doc := range docs {
		content := strings.ToLower(doc.Content)
		score := 0

		for _, term := range queryTerms {
			score += strings.Count(content, term)
		}

		if score > 0 {
			candidates = append(candidates, scoredDoc{
				doc:   doc,
				score: score,
			})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].score == candidates[j].score {
			return candidates[i].doc.Name < candidates[j].doc.Name
		}
		return candidates[i].score > candidates[j].score
	})

	if len(candidates) > topK {
		candidates = candidates[:topK]
	}

	var builder strings.Builder
	var files []string

	for _, candidate := range candidates {
		files = append(files, candidate.doc.Name)

		fmt.Fprintf(
			&builder,
			"\n\n--- FILE: %s (score=%d) ---\n%s",
			candidate.doc.Name,
			candidate.score,
			candidate.doc.Content,
		)
	}

	return Result{
		Mode:    "search",
		Files:   files,
		Context: strings.TrimSpace(builder.String()),
	}
}

func tokenize(text string) []string {
	text = strings.ToLower(text)

	var terms []string
	var current strings.Builder
	var currentIsChinese bool

	flush := func() {
		if current.Len() == 0 {
			return
		}
		terms = append(terms, current.String())
		current.Reset()
	}

	for _, r := range text {
		isChinese := unicode.In(r, unicode.Han)
		isASCII := r <= unicode.MaxASCII &&
			(unicode.IsLetter(r) || unicode.IsDigit(r))

		if isChinese {
			if current.Len() > 0 && !currentIsChinese {
				flush()
			}
			currentIsChinese = true
			current.WriteRune(r)
		} else if isASCII {
			if current.Len() > 0 && currentIsChinese {
				flush()
			}
			currentIsChinese = false
			current.WriteRune(r)
		} else {
			flush()
		}
	}

	flush()
	return terms
}
