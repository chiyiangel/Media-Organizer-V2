package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/photo-video-organizer/internal/ui"
)

func main() {
	// 创建模型
	m := ui.NewModel()

	// 创建程序
	p := tea.NewProgram(m, tea.WithAltScreen())

	// 运行程序
	if _, err := p.Run(); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
}
