package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chiyiangel/media-organizer-v2/internal/i18n"
	"github.com/chiyiangel/media-organizer-v2/internal/ui"
)

func main() {
	// 初始化多语言系统，自动检测系统语言
	_ = i18n.GetLocalizer() // 这会触发语言检测

	// 创建模型
	m := ui.NewModel()

	// 创建程序
	p := tea.NewProgram(m, tea.WithAltScreen())

	// 运行程序
	if _, err := p.Run(); err != nil {
		fmt.Printf(i18n.T("error.prefix")+"%v\n", err)
		os.Exit(1)
	}
}
