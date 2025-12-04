package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bigwhite/my-issue2md/internal/cli"
	"github.com/bigwhite/my-issue2md/internal/config"
	"github.com/bigwhite/my-issue2md/internal/converter"
	"github.com/bigwhite/my-issue2md/internal/github"
	"github.com/bigwhite/my-issue2md/internal/parser"
)

const (
	name    = "issue2md"
	version = "1.0.0"
	usage   = `issue2md - Convert GitHub issues to Markdown

Usage:
  issue2md [owner/repo] [issue-number] [flags]

Examples:
  issue2md facebook/react 12345
  issue2md facebook/react 12345 --output=issue.md
  issue2md facebook/react 12345 --format=html --no-comments

Flags:
  -h, --help              Show help information
  -v, --version           Show version information
  -o, --output string     Output file path (default: "output.md")
  -f, --format string     Output format: markdown, html, json (default: "markdown")
  -t, --token string      GitHub token (or set GITHUB_TOKEN env var)
  --no-comments          Exclude comments from output
  --no-metadata          Exclude metadata from output
  --no-timestamps        Exclude timestamps from output
  --overwrite            Overwrite existing output file
  --debug                Enable debug logging`
)

func main() {
	// 创建CLI实例
	app := cli.NewCLI(name, os.Args[1:])

	// 设置信号处理
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	// 加载配置
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()

	// 解析命令行参数
	if err := runCLI(ctx, app, cfg); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// runCLI 执行CLI逻辑
func runCLI(ctx context.Context, app *cli.CLI, cfg *config.Config) error {
	// TODO: 实现命令行参数解析和业务逻辑
	fmt.Printf("%s v%s\n", name, version)
	fmt.Println("Phase 1 Foundation completed successfully!")
	fmt.Println("Core data structures have been initialized.")

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}

// initializeServices 初始化服务
func initializeServices(cfg *config.Config) (*github.GitHubClient, *parser.MarkdownParser, converter.Converter, error) {
	// 初始化GitHub客户端
	githubClient := github.NewClient(cfg.GitHubToken)

	// 初始化解析器
	parserOptions := &parser.Options{
		IncludeComments:    cfg.Parser.IncludeComments,
		IncludeMetadata:    cfg.Parser.IncludeMetadata,
		IncludeTimestamps:  cfg.Parser.IncludeTimestamps,
		IncludeUserLinks:   cfg.Parser.IncludeUserLinks,
		EmojisEnabled:      cfg.Parser.EmojisEnabled,
		PreserveLineBreaks: cfg.Parser.PreserveLineBreaks,
	}
	markdownParser := parser.NewParser(parserOptions)

	// 初始化转换器
	converterOptions := converter.DefaultConverterOptions()
	var conv converter.Converter

	switch cfg.Output.Format {
	case "html":
		conv = converter.NewHTMLConverter(converterOptions)
	case "json":
		conv = converter.NewJSONConverter(converterOptions)
	default:
		conv = converter.NewMarkdownConverter(converterOptions)
	}

	return githubClient, markdownParser, conv, nil
}