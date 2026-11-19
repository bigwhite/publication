package ingest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func Run(cfg Config) (Stats, error) {
	var stats Stats

	if cfg.SourceRoot == "" {
		return stats, fmt.Errorf("source root is required")
	}
	if cfg.KnowledgeRoot == "" {
		return stats, fmt.Errorf("knowledge root is required")
	}
	if cfg.ManifestPath == "" {
		return stats, fmt.Errorf("manifest path is required")
	}

	if err := os.MkdirAll(cfg.KnowledgeRoot, 0755); err != nil {
		return stats, err
	}

	currentManifest, err := loadManifest(cfg.ManifestPath)
	if err != nil {
		return stats, err
	}

	nextManifest := make(manifest)
	seenFingerprints := make(map[string]string)
	for _, entry := range currentManifest {
		seenFingerprints[entry.Fingerprint] = entry.Output
	}

	seenPaths := make(map[string]bool)

	err = filepath.WalkDir(cfg.SourceRoot, func(
		path string,
		entry os.DirEntry,
		walkErr error,
	) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(cfg.SourceRoot, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		parser, err := parserFor(path)
		if err != nil {
			// 暂不支持的格式不进入知识库。
			return nil
		}

		stats.Scanned++
		seenPaths[rel] = true

		raw, err := os.ReadFile(path)
		if err != nil {
			stats.Failed++
			return fmt.Errorf("read source %s: %w", rel, err)
		}

		fingerprint := hashString(string(raw))
		id := hashString(rel)

		if old, ok := currentManifest[rel]; ok &&
			old.Fingerprint == fingerprint {
			nextManifest[rel] = old
			stats.Skipped++
			return nil
		}

		// 相同内容使用同一个规范化文件。
		if output, ok := seenFingerprints[fingerprint]; ok {
			nextManifest[rel] = manifestEntry{
				ID:          id,
				Fingerprint: fingerprint,
				Output:      output,
			}

			if _, exists := currentManifest[rel]; exists {
				stats.Updated++
			} else {
				stats.Added++
			}
			return nil
		}

		parsed, err := parser.Parse(path, raw)
		if err != nil {
			stats.Failed++
			return fmt.Errorf("parse source %s: %w", rel, err)
		}

		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		item := KnowledgeItem{
			ID:          id,
			SourcePath:  rel,
			Format:      strings.TrimPrefix(filepath.Ext(path), "."),
			Title:       parsed.Title,
			Headings:    parsed.Headings,
			Content:     parsed.Content,
			Fingerprint: fingerprint,
			UpdatedAt:   info.ModTime().UTC(),
		}

		output := outputName(id, path)
		outputPath := filepath.Join(cfg.KnowledgeRoot, output)

		if err := os.WriteFile(
			outputPath,
			[]byte(renderKnowledge(item)),
			0644,
		); err != nil {
			return err
		}

		nextManifest[rel] = manifestEntry{
			ID:          id,
			Fingerprint: fingerprint,
			Output:      output,
		}

		seenFingerprints[fingerprint] = output

		if _, exists := currentManifest[rel]; exists {
			stats.Updated++
		} else {
			stats.Added++
		}

		return nil
	})
	if err != nil {
		return stats, err
	}

	// 清理已经从来源目录删除的记录。
	for sourcePath, _ := range currentManifest {
		if seenPaths[sourcePath] {
			continue
		}
		stats.Deleted++
	}

	// 只删除不再被任何来源引用的输出文件。
	referenced := make(map[string]bool)
	for _, entry := range nextManifest {
		referenced[entry.Output] = true
	}

	for _, old := range currentManifest {
		if referenced[old.Output] {
			continue
		}
		_ = os.Remove(filepath.Join(cfg.KnowledgeRoot, old.Output))
	}

	if err := saveManifest(cfg.ManifestPath, nextManifest); err != nil {
		return stats, err
	}

	return stats, nil
}

func loadManifest(path string) (manifest, error) {
	result := make(manifest)

	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return result, nil
	}
	if err != nil {
		return nil, err
	}

	if len(raw) == 0 {
		return result, nil
	}

	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}

	return result, nil
}

func saveManifest(path string, data manifest) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0644); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}

func hashString(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func outputName(id, sourcePath string) string {
	base := filepath.Base(sourcePath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")

	return id[:12] + "-" + name + ".md"
}

func renderKnowledge(item KnowledgeItem) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "---\n")
	fmt.Fprintf(&builder, "id: %s\n", item.ID)
	fmt.Fprintf(&builder, "source: %s\n", item.SourcePath)
	fmt.Fprintf(&builder, "format: %s\n", item.Format)
	fmt.Fprintf(
		&builder,
		"updated_at: %s\n",
		item.UpdatedAt.Format(time.RFC3339),
	)
	fmt.Fprintf(&builder, "---\n\n")

	fmt.Fprintf(&builder, "# %s\n\n", item.Title)

	if len(item.Headings) > 0 {
		builder.WriteString("## Document Headings\n\n")
		for _, heading := range item.Headings {
			fmt.Fprintf(&builder, "- %s\n", heading)
		}
		builder.WriteString("\n")
	}

	builder.WriteString(item.Content)

	return builder.String()
}

// 保证 Manifest 序列化前的输出顺序稳定，方便检查和调试。
func sortedManifestKeys(data manifest) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
