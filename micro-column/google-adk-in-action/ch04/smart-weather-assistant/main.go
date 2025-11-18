package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// WeatherArgs 定义了天气查询工具的输入参数
type WeatherArgs struct {
	City string `json:"city" jsonschema:"city name"`
}

// WeatherData 用于内部状态存储
type WeatherData struct {
	City      string `json:"city"`
	Summary   string `json:"summary"`
	Timestamp string `json:"timestamp"`
}

func GetWeather(ctx tool.Context, args WeatherArgs) (string, error) {
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
		return fmt.Sprintf("Sorry, I don't have the weather information for %s.", args.City), nil
	}

	// 存储到状态中用于对比
	currentWeather := WeatherData{
		City:      args.City,
		Summary:   summary,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var weatherHistory []WeatherData
	val, err := ctx.State().Get("weather_comparison_list")
	if err != nil {
		if errors.Is(err, session.ErrStateKeyNotExist) {
			weatherHistory = []WeatherData{}
		} else {
			log.Printf("Warning: failed to get state: %v", err)
			weatherHistory = []WeatherData{}
		}
	} else {
		if err := mapstructure.Decode(val, &weatherHistory); err != nil {
			log.Printf("Warning: 'weather_comparison_list' has unexpected structure, resetting. Error: %v", err)
			weatherHistory = []WeatherData{}
		}
	}

	weatherHistory = append(weatherHistory, currentWeather)
	if err := ctx.State().Set("weather_comparison_list", weatherHistory); err != nil {
		log.Printf("Warning: failed to set state 'weather_comparison_list': %v", err)
	}

	// 返回简单的字符串
	return fmt.Sprintf("The weather in %s is %s (updated at %s)", 
		args.City, summary, time.Now().Format("15:04:05")), nil
}

type CompareWeatherArgs struct{}

func CompareWeather(ctx tool.Context, args CompareWeatherArgs) (string, error) {
	var weatherHistory []WeatherData
	val, err := ctx.State().Get("weather_comparison_list")
	if err != nil {
		return "", fmt.Errorf("You haven't queried any weather yet. Please ask me about some cities first")
	}

	if err := mapstructure.Decode(val, &weatherHistory); err != nil || len(weatherHistory) == 0 {
		return "", fmt.Errorf("No weather data available to compare")
	}

	// 构建对比结果字符串
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Weather comparison for %d cities:\n", len(weatherHistory)))
	for i, w := range weatherHistory {
		result.WriteString(fmt.Sprintf("%d. %s: %s (checked at %s)\n", 
			i+1, w.City, w.Summary, w.Timestamp))
	}

	return result.String(), nil
}

type RecallFavoriteCityArgs struct{}

// RecallFavoriteCity 工具从长期记忆中搜索常用城市
func RecallFavoriteCity(ctx tool.Context, args RecallFavoriteCityArgs) (string, error) {
	log.Println("INFO: [Tool] Attempting to recall favorite city from long-term memory...")

	knownCities := []string{"beijing", "shanghai", "shenzhen"}
	cityCounts := make(map[string]int)

	// 遍历所有已知城市，并为每一个城市执行一次关键词搜索
	for _, city := range knownCities {
		searchResp, err := ctx.SearchMemory(ctx, city)
		if err != nil {
			log.Printf("Warning: SearchMemory for city '%s' failed: %v", city, err)
			continue // 继续搜索下一个城市
		}
		if len(searchResp.Memories) > 0 {
			// 简单地用返回的记忆条目数量作为该城市的“提及次数”
			cityCounts[city] = len(searchResp.Memories)
			log.Printf("INFO: [Tool] Found %d memories mentioning '%s'", len(searchResp.Memories), city)
		}
	}

	if len(cityCounts) == 0 {
		return "", fmt.Errorf("I don't have any memory of cities you've asked about before")
	}

	// 找出最高频的城市
	favoriteCity := ""
	maxCount := 0
	for city, count := range cityCounts {
		if count > maxCount {
			maxCount = count
			favoriteCity = city
		}
	}

	if favoriteCity == "" {
		// 这种情况理论上不会发生，除非 cityCounts 不为空但所有 count 都是 0
		return "", fmt.Errorf("I found some memories but couldn't determine a favorite city.")
	}

	log.Printf("INFO: [Tool] Recalled favorite city: %s", favoriteCity)
	return fmt.Sprintf("Based on your history, your favorite city is %s (mentioned %d times)", 
		favoriteCity, maxCount), nil
}

func demonstrateLongTermMemory(ctx context.Context, agentToDemo agent.Agent, sessionService session.Service, memoryService memory.Service) {
	appName := "smart_weather_app"
	userID := "test_user_123"

	r, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          agentToDemo,
		SessionService: sessionService,
		MemoryService:  memoryService,
	})
	if err != nil {
		log.Fatalf("Failed to create runner for demo: %v", err)
	}

	fmt.Println("\n--- 🚀 Starting Session 1: Establishing a Memory ---")

	session1ID := "session-1"
	_, err = sessionService.Create(ctx, &session.CreateRequest{AppName: appName, UserID: userID, SessionID: session1ID})
	if err != nil {
		log.Fatalf("Failed to create session 1: %v", err)
	}

	fmt.Println("User -> what is the weather in beijing?")
	stream1 := r.Run(ctx, userID, session1ID, genai.NewContentFromText("what is the weather in beijing?", genai.RoleUser), agent.RunConfig{})
	for event, err := range stream1 {
		if err != nil {
			log.Fatalf("Session 1 run failed: %v", err)
		}
		if !event.Partial && event.LLMResponse.Content != nil && len(event.LLMResponse.Content.Parts) > 0 {
			fmt.Printf("Agent -> %s\n", event.LLMResponse.Content.Parts[0].Text)
		}
	}

	session1, err := sessionService.Get(ctx, &session.GetRequest{AppName: appName, UserID: userID, SessionID: session1ID})
	if err != nil {
		log.Fatalf("Failed to get session 1 after run: %v", err)
	}
	if err := memoryService.AddSession(ctx, session1.Session); err != nil {
		log.Fatalf("Failed to add session 1 to memory: %v", err)
	}
	fmt.Println("--- ✅ Session 1 Ended. Memory should now contain 'beijing'. ---")

	fmt.Println("\n--- 🚀 Starting Session 2: Recalling from Memory ---")

	session2ID := "session-2"
	_, err = sessionService.Create(ctx, &session.CreateRequest{AppName: appName, UserID: userID, SessionID: session2ID})
	if err != nil {
		log.Fatalf("Failed to create session 2: %v", err)
	}

	fmt.Println("User -> Can you check the weather for me?")
	stream2 := r.Run(ctx, userID, session2ID, genai.NewContentFromText("Can you check the weather for me?", genai.RoleUser), agent.RunConfig{})
	for event, err := range stream2 {
		if err != nil {
			log.Fatalf("Session 2 run failed: %v", err)
		}
		if !event.Partial && event.LLMResponse.Content != nil && len(event.LLMResponse.Content.Parts) > 0 {
			fmt.Printf("Agent -> %s\n", event.LLMResponse.Content.Parts[0].Text)
		}
	}
	fmt.Println("--- ✅ Session 2 Ended. ---")
}

func main() {
	ctx := context.Background()
	model, err := gemini.NewModel(ctx, "gemini-2.0-flash", &genai.ClientConfig{APIKey: os.Getenv("GOOGLE_API_KEY")})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	weatherTool, err := functiontool.New(
		functiontool.Config{
			Name:        "get_weather",
			Description: "Retrieves the current weather report for a specified city.",
		},
		functiontool.Func[WeatherArgs, string](GetWeather),
	)
	if err != nil {
		log.Fatalf("Failed to create weather tool: %v", err)
	}

	compareTool, err := functiontool.New(
		functiontool.Config{
			Name:        "compare_weather_of_queried_cities",
			Description: "Compares the weather of all cities queried in the current conversation and provides a recommendation. Use this when the user asks for a comparison or recommendation.",
		},
		functiontool.Func[CompareWeatherArgs, string](CompareWeather),
	)
	if err != nil {
		log.Fatalf("Failed to create compare tool: %v", err)
	}

	recallFavoriteTool, err := functiontool.New(
		functiontool.Config{
			Name:        "recall_favorite_city",
			Description: "Recalls the user's most frequently queried city from past conversations. Use this ONLY when the user asks for weather but does not specify any city.",
		},
		functiontool.Func[RecallFavoriteCityArgs, string](RecallFavoriteCity),
	)
	if err != nil {
		log.Fatalf("Failed to create recall favorite tool: %v", err)
	}

	weatherAgent, err := llmagent.New(llmagent.Config{
		Name:        "smart_weather_assistant",
		Description: "An agent that provides weather information and helps with travel decisions.",
		Model:       model,
		Tools:       []tool.Tool{weatherTool, compareTool, recallFavoriteTool},
		Instruction: `You are a helpful travel weather assistant.
- Use 'get_weather' to find weather for specific cities.
- When asked for a comparison or recommendation (e.g., "which is better?", "summarize them"), you MUST use the 'compare_weather_of_queried_cities' tool to get all the data, and then provide a final answer based on that data.
- If the user asks for weather but provides NO city, use 'recall_favorite_city' to suggest a city they might be interested in, and then get the weather for that city.`,
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	sessionService := session.InMemoryService()
	memoryService := memory.InMemoryService()

	if len(os.Args) > 1 && os.Args[1] == "memory-demo" {
		demonstrateLongTermMemory(ctx, weatherAgent, sessionService, memoryService)
	} else {
		config := &launcher.Config{
			AgentLoader:    agent.NewSingleLoader(weatherAgent),
			SessionService: sessionService,
			MemoryService:  memoryService,
		}

		args := []string{"console"}
		if len(os.Args) > 1 {
			args = os.Args[1:]
		}

		l := full.NewLauncher()
		if err = l.Execute(ctx, config, args); err != nil {
			log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
		}
	}
}
