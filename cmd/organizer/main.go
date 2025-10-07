package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chiyiangel/media-organizer-v2/internal/app"
	"github.com/chiyiangel/media-organizer-v2/internal/config"
	"github.com/chiyiangel/media-organizer-v2/internal/i18n"
	"github.com/chiyiangel/media-organizer-v2/internal/ui"
)

func main() {
	// Parse CLI arguments first
	cliConfig, err := ParseCLI()
	if err != nil {
		errorMsg := i18n.Tf("cli.error.cli_parse", err) + "\n"
		fmt.Print(errorMsg)
		os.Exit(1)
	}

	// Load full configuration with precedence: CLI > File > Defaults
	finalConfig, err := config.LoadFullConfig(cliConfig)
	if err != nil {
		errorMsg := i18n.Tf("cli.error.config_load", err) + "\n"
		fmt.Print(errorMsg)
		os.Exit(1)
	}

	// Validate configuration
	if err := finalConfig.Validate(); err != nil {
		errorMsg := i18n.Tf("cli.error.config_validate", err) + "\n"
		fmt.Print(errorMsg)
		os.Exit(1)
	}

	// Route based on operation mode
	switch finalConfig.Mode {
	case config.ModeSilent:
		runSilentMode(finalConfig)
	case config.ModeInteractive:
		fallthrough
	default:
		runInteractiveMode(finalConfig)
	}
}

// runSilentMode executes the application in silent mode
func runSilentMode(config *config.Config) {
	// Create silent runner
	runner, err := app.NewSilentRunner(config)
	if err != nil {
		errorMsg := i18n.Tf("cli.error.silent_runner", err) + "\n"
		fmt.Print(errorMsg)
		os.Exit(1)
	}

	// Execute in silent mode
	if err := runner.Run(); err != nil {
		errorMsg := i18n.Tf("cli.error.silent_exec", err) + "\n"
		fmt.Print(errorMsg)
		os.Exit(1)
	}
}

// runInteractiveMode executes the application in interactive TUI mode
func runInteractiveMode(config *config.Config) {
	// 初始化多语言系统，自动检测系统语言
	_ = i18n.GetLocalizer() // 这会触发语言检测

	// 创建模型
	m := ui.NewModel()

	// 创建程序
	p := tea.NewProgram(m, tea.WithAltScreen())

	// 运行程序
	if _, err := p.Run(); err != nil {
		errorMsg := i18n.T("error.prefix") + fmt.Sprintf("%v\n", err)
		fmt.Print(errorMsg)
		os.Exit(1)
	}
}
