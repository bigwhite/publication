package config

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// AppConfig, ServerConfig, DatabaseConfig 结构体定义同stage1
// 但需要使用 `mapstructure` 标签以供 Viper Unmarshal
type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	Timeout string `mapstructure:"timeout"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type AppConfig struct {
	AppName  string         `mapstructure:"appName"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ViperInstance 是一个导出的 viper 实例，方便在应用中其他地方按需获取配置
// 或者，你也可以将 Load 返回的 AppConfig 实例通过 DI 传递
var ViperInstance *viper.Viper

func init() {
	ViperInstance = viper.New()
}

// LoadConfigWithViper initializes and loads configuration using Viper.
func LoadConfigWithViper(configPath string, configName string, configType string) (*AppConfig, error) {
	v := ViperInstance // Use the global instance or a new one

	// 1. 设置默认值
	v.SetDefault("server.port", 8080)
	v.SetDefault("appName", "DefaultViperAppFromCode")

	// 2. 绑定命令行参数 (使用 pflag)
	// pflag 的定义通常在 main 包的 init 中，或者一个集中的 flag 定义文件
	// 这里为了示例完整性，假设已定义并 Parse
	if pflag.Parsed() { // Ensure flags are parsed before binding
		err := v.BindPFlags(pflag.CommandLine)
		if err != nil {
			return nil, fmt.Errorf("stage4: failed to bind pflags: %w", err)
		}
	} else {
		// This case might happen if LoadConfigWithViper is called before pflag.Parse()
		// For a robust setup, ensure pflag.Parse() is called in main before this.
		// For this example, we'll assume flags are available if pflag was touched.
		// A better approach is to pass the *pflag.FlagSet to this function.
		fmt.Println("[Stage4 Config] pflag not parsed, skipping BindPFlags. Ensure pflag.Parse() is called in main.")
	}

	// 3. 绑定环境变量
	v.SetEnvPrefix("MYAPP")                            // e.g., MYAPP_SERVER_PORT, MYAPP_DATABASE_DSN
	v.AutomaticEnv()                                   // Automatically read matching env variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // server.port -> SERVER_PORT

	// 4. 设置配置文件路径和类型
	if configPath != "" {
		v.AddConfigPath(configPath) // 如 "./configs"
	}
	v.AddConfigPath("$HOME/.myapp") // HOME目录
	v.AddConfigPath(".")            // 当前工作目录
	v.SetConfigName(configName)     // "app" (不带扩展名)
	v.SetConfigType(configType)     // "yaml"

	// 5. 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到是可接受的，可能依赖环境变量或默认值
			fmt.Printf("[Stage4 Config] Config file '%s.%s' not found in search paths; relying on other sources.\n", configName, configType)
		} else {
			// 其他读取配置文件的错误
			return nil, fmt.Errorf("stage4: failed to read config file: %w", err)
		}
	} else {
		fmt.Printf("[Stage4 Config] Using config file: %s\n", v.ConfigFileUsed())
	}

	// 6. Unmarshal到结构体
	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("stage4: failed to unmarshal config to struct: %w", err)
	}

	fmt.Printf("[Stage4 Config] Successfully loaded and unmarshalled. AppName: %s, ServerPort: %d, DBDsn: %s\n",
		cfg.AppName, cfg.Server.Port, cfg.Database.DSN)
	return &cfg, nil
}
