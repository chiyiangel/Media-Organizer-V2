package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// 边框样式 - 减少内边距避免布局冲突
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 1)

	// 标题样式
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Align(lipgloss.Center)

	// 标签样式
	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	// 普通文本样式
	textStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	// 高亮样式
	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)

	// 成功样式
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	// 错误样式
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	// 警告样式
	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	// 分割线样式
	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238"))

	// 提示样式
	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

	// 选中样式
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)
)

// renderProgressBar 进度条渲染
func renderProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	filled := int(float64(width) * percentage)

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return fmt.Sprintf("%s %.0f%% (%d/%d)", bar, percentage*100, current, total)
}

// renderDivider 渲染分割线
func renderDivider(width int) string {
	return dividerStyle.Render(strings.Repeat("─", width))
}

// renderBox 渲染边框盒子
func renderBox(title, content string, width int) string {
	titleStr := titleStyle.Width(width).Render(title)
	contentStr := textStyle.Width(width).Render(content)

	return borderStyle.Width(width).Render(titleStr + "\n" + contentStr)
}
