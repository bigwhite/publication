package main

import (
	"ch23/config/stage1/config"
	"fmt"
	"log"
)

func main() {
	// 为了能直接运行，你可能需要调整配置文件路径，并确保文件存在
	// 在实际项目中，这个路径通常来自命令行参数或固定位置
	err := config.LoadGlobalConfig("app.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if config.GlobalAppConfig != nil {
		fmt.Printf("Accessing config: AppName is '%s', Server Port is %d\n",
			config.GlobalAppConfig.AppName, config.GlobalAppConfig.Server.Port)
		fmt.Printf("Database DSN: %s\n", config.GlobalAppConfig.Database.DSN)
	} else {
		fmt.Println("Config was not loaded.")
	}
}
