package cli

import (
	"flag"
	"fmt"
	"os"
)

// CLI 命令行接口
type CLI struct {
	name   string
	args   []string
	output *Output
}

// NewCLI 创建新的CLI实例
func NewCLI(name string, args []string) *CLI {
	return &CLI{
		name: name,
		args: args,
		output: &Output{
			Writer: os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}
}

// Output 输出配置
type Output struct {
	Writer      interface{}
	ErrorWriter interface{}
}

// Command 命令定义
type Command struct {
	Name        string
	Description string
	Flags       *flag.FlagSet
	Run         func(*Context) error
}

// Context 命令执行上下文
type Context struct {
	Args   []string
	Flags  map[string]string
	Output *Output
}

// CommandRegistry 命令注册表
type CommandRegistry struct {
	commands map[string]*Command
}

// NewCommandRegistry 创建命令注册表
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]*Command),
	}
}

// Register 注册命令
func (cr *CommandRegistry) Register(cmd *Command) error {
	if _, exists := cr.commands[cmd.Name]; exists {
		return fmt.Errorf("command %s already registered", cmd.Name)
	}
	cr.commands[cmd.Name] = cmd
	return nil
}

// Get 获取命令
func (cr *CommandRegistry) Get(name string) (*Command, bool) {
	cmd, exists := cr.commands[name]
	return cmd, exists
}

// List 列出所有命令
func (cr *CommandRegistry) List() []*Command {
	commands := make([]*Command, 0, len(cr.commands))
	for _, cmd := range cr.commands {
		commands = append(commands, cmd)
	}
	return commands
}

// Error CLI错误类型
type Error struct {
	Message string
	Code    int
}

func (e *Error) Error() string {
	return e.Message
}

// NewError 创建CLI错误
func NewError(message string, code int) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}