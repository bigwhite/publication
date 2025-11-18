package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// WeatherArgs 定义了天气查询工具的输入参数
type WeatherArgs struct {
	City string `json:"city" jsonschema:"city name"`
}

// WeatherOutput 定义了天气查询工具的返回值
type WeatherOutput struct {
	Summary   string    `json:"summary"`
	Timestamp time.Time `json:"timestamp"`
}

// GetWeather 是我们工具的核心实现。
// 它接收一个 tool.Context 和我们定义的输入参数结构体，返回输出结构体。
func GetWeather(ctx tool.Context, args WeatherArgs) (WeatherOutput, error) {
	// 在真实的场景中，这里会调用一个外部的天气 API。
	// 为了简化，我们在这里模拟一下。
	city := strings.ToLower(args.City)
	var summary string

	switch city {
	case "beijing":
		summary = "Sunny, 25°C"
	case "shanghai":
		summary = "Cloudy, 22°C"
	case "shenzhen":
		summary = "Rainy, 28°C"
	default:
		summary = fmt.Sprintf("Sorry, I don't have the weather information for %s.", args.City)
	}

	return WeatherOutput{
		Summary:   summary,
		Timestamp: time.Now(),
	},nil
}

func main() {
	ctx := context.Background()

	// 1. 初始化模型 (与上一讲相同)
	model, err := gemini.NewModel(ctx, "gemini-2.0-flash-lite", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	// 2. 将 Go 函数包装成 Tool
	weatherTool, err := functiontool.New[WeatherArgs, WeatherOutput](functiontool.Config{
		Name:        "get_weather",
		Description: "Retrieves the current weather report for a specified city.",
	//}, functiontool.Func[WeatherArgs, WeatherOutput](GetWeather))
	}, GetWeather)
	if err != nil {
		log.Fatalf("Failed to create weather tool: %v", err)
	}

	// 3. 创建 LLMAgent 并注入 Tool
	weatherAgent, err := llmagent.New(llmagent.Config{
		Name:        "weather_agent",
		Description: "An agent that can provide real-time weather information.",
		Model:       model,
		Tools:       []tool.Tool{weatherTool}, // 在这里注入我们的工具！
		Instruction: "You are a friendly weather assistant. Your job is to answer user's questions about the weather. Use the available tools to get the weather information. If the information for a city is not available, inform the user politely.",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// 4. 配置并启动 Launcher (与上一讲相同)
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(weatherAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
