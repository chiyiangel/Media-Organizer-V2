package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/yourusername/photo-video-organizer/internal/config"
)

// renderConfigScreen 渲染主配置界面
func (m Model) renderConfigScreen() string {
	var b strings.Builder

	// 标题
	b.WriteString(titleStyle.Width(m.width).Render("📸 照片视频整理工具 v1.0"))
	b.WriteString("\n\n")

	// 源目录
	b.WriteString(labelStyle.Render("📁 源目录: "))
	if m.config.SourceDir == "" {
		b.WriteString(hintStyle.Render("未设置"))
	} else {
		b.WriteString(textStyle.Render(m.config.SourceDir))
	}
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("           按 [S] 编辑路径"))
	b.WriteString("\n\n")

	// 目标目录
	b.WriteString(labelStyle.Render("📂 目标目录: "))
	if m.config.TargetDir == "" {
		b.WriteString(hintStyle.Render("未设置"))
	} else {
		b.WriteString(textStyle.Render(m.config.TargetDir))
	}
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("           按 [D] 编辑路径"))
	b.WriteString("\n\n")

	// 整理策略
	b.WriteString(labelStyle.Render("⚙️  整理策略:"))
	b.WriteString("\n")

	// 同文件识别
	detectionF := " "
	detectionM := " "
	if m.config.DuplicateDetection == config.DetectionFilename {
		detectionF = "●"
	} else {
		detectionM = "●"
	}

	b.WriteString(textStyle.Render(fmt.Sprintf("    同文件识别: [F] 文件名 %s  [M] MD5哈希 %s", detectionF, detectionM)))
	b.WriteString("\n")

	// 重复处理
	strategy1 := " "
	strategy2 := " "
	strategy3 := " "
	switch m.config.DuplicateStrategy {
	case config.StrategySkip:
		strategy1 = "●"
	case config.StrategyOverwrite:
		strategy2 = "●"
	case config.StrategyRename:
		strategy3 = "●"
	}

	b.WriteString(textStyle.Render(fmt.Sprintf("    重复处理:   [1] 跳过 %s  [2] 覆盖 %s  [3] 重命名 %s", strategy1, strategy2, strategy3)))
	b.WriteString("\n\n")

	// 分割线 - 调整宽度以匹配边框
	dividerWidth := m.width - 8 // 考虑边框和内边距
	if dividerWidth > 0 {
		b.WriteString(renderDivider(dividerWidth))
		b.WriteString("\n\n")
	}

	// 错误信息
	if m.err != nil {
		b.WriteString(errorStyle.Render("错误: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	// 提示 - 确保文字不超出容器宽度
	hintText := "按 [Enter] 开始整理  |  按 [Q/Esc] 退出程序"
	maxHintWidth := m.width - 8 // 考虑边框和内边距
	if len(hintText) > maxHintWidth && maxHintWidth > 0 {
		hintText = "按 [Enter] 开始整理\n按 [Q/Esc] 退出程序"
	}
	b.WriteString(hintStyle.Render(hintText))
	b.WriteString("\n")

	// 确保容器宽度合理
	containerWidth := m.width - 2
	if containerWidth < 40 {
		containerWidth = 40 // 最小宽度
	}

	return borderStyle.Width(containerWidth).Render(b.String())
}

// renderInputScreen 渲染路径输入界面
func (m Model) renderInputScreen() string {
	var b strings.Builder

	// 标题
	title := "📝 编辑路径"
	b.WriteString(titleStyle.Width(m.width).Render(title))
	b.WriteString("\n\n")

	// 提示
	prompt := "请输入目录路径:"
	b.WriteString(labelStyle.Render(prompt))
	b.WriteString("\n")

	// 输入框
	inputBox := fmt.Sprintf("│ %s_", m.inputValue)
	b.WriteString(textStyle.Render(inputBox))
	b.WriteString("\n\n")

	// 提示
	b.WriteString(hintStyle.Render("按 [Enter] 确认  |  按 [Esc] 取消"))
	b.WriteString("\n")

	// 确保容器宽度合理
	containerWidth := m.width - 2
	if containerWidth < 40 {
		containerWidth = 40
	}

	return borderStyle.Width(containerWidth).Render(b.String())
}

// renderProgressScreen 渲染整理进度界面
func (m Model) renderProgressScreen() string {
	var b strings.Builder

	// 标题
	b.WriteString(titleStyle.Width(m.width).Render("🔄 正在整理文件..."))
	b.WriteString("\n\n")

	// 当前文件
	if m.currentFile != nil {
		b.WriteString(labelStyle.Render("当前文件: "))
		b.WriteString(textStyle.Render(m.currentFile.Name))
		b.WriteString("\n")

		b.WriteString(labelStyle.Render("目标路径: "))
		b.WriteString(textStyle.Render(m.currentFile.TargetPath))
		b.WriteString("\n\n")
	}

	// 进度条
	progressBarWidth := m.width - 20
	if progressBarWidth < 20 {
		progressBarWidth = 20
	}
	b.WriteString(labelStyle.Render("进度: "))
	b.WriteString(renderProgressBar(m.statistics.ProcessedFiles, m.statistics.TotalFiles, progressBarWidth))
	b.WriteString("\n\n")

	// 分割线 - 统一宽度计算
	dividerWidth := m.width - 8
	if dividerWidth > 0 {
		b.WriteString(renderDivider(dividerWidth))
		b.WriteString("\n\n")
	}

	// 实时统计
	b.WriteString(labelStyle.Render("📊 实时统计:"))
	b.WriteString("\n")
	b.WriteString(textStyle.Render(fmt.Sprintf("    已扫描:  %d 个文件\n", m.statistics.ScannedFiles)))
	b.WriteString(textStyle.Render(fmt.Sprintf("    已处理:  %d 个文件\n", m.statistics.ProcessedFiles)))
	b.WriteString(textStyle.Render(fmt.Sprintf("    ├─ 照片: %d 张\n", m.statistics.PhotoCount)))
	b.WriteString(textStyle.Render(fmt.Sprintf("    ├─ 视频: %d 个\n", m.statistics.VideoCount)))
	b.WriteString(textStyle.Render(fmt.Sprintf("    ├─ 跳过: %d 个 (重复)\n", m.statistics.SkippedCount)))
	b.WriteString(textStyle.Render(fmt.Sprintf("    └─ 错误: %d 个\n", m.statistics.FailedCount)))
	b.WriteString("\n")

	// 分割线
	if dividerWidth > 0 {
		b.WriteString(renderDivider(dividerWidth))
		b.WriteString("\n\n")
	}

	// 提示
	b.WriteString(hintStyle.Render("按 [C/Esc] 取消整理"))
	b.WriteString("\n")

	// 确保容器宽度合理
	containerWidth := m.width - 2
	if containerWidth < 40 {
		containerWidth = 40
	}

	return borderStyle.Width(containerWidth).Render(b.String())
}

// renderSummaryScreen 渲染完成汇总界面
func (m Model) renderSummaryScreen() string {
	var b strings.Builder

	// 标题
	b.WriteString(titleStyle.Width(m.width).Render("✅ 整理完成!"))
	b.WriteString("\n\n")

	// 汇总报告标题
	b.WriteString(labelStyle.Render("📊 整理汇总报告"))
	b.WriteString("\n")
	b.WriteString(renderDivider(m.width - 4))
	b.WriteString("\n\n")

	// 文件统计
	b.WriteString(labelStyle.Render("文件统计:"))
	b.WriteString("\n")
	b.WriteString(textStyle.Render(fmt.Sprintf("    总文件数:      %d 个\n", m.statistics.TotalFiles)))
	b.WriteString(textStyle.Render(fmt.Sprintf("    ├─ 照片:       %d 张\n", m.statistics.PhotoCount)))
	b.WriteString(textStyle.Render(fmt.Sprintf("    └─ 视频:       %d 个\n", m.statistics.VideoCount)))
	b.WriteString("\n")

	// 处理结果
	successCount := m.statistics.ProcessedFiles - m.statistics.SkippedCount - m.statistics.FailedCount
	b.WriteString(labelStyle.Render("处理结果:"))
	b.WriteString("\n")
	b.WriteString(successStyle.Render(fmt.Sprintf("    ✓ 成功整理:    %d 个\n", successCount)))
	b.WriteString(warningStyle.Render(fmt.Sprintf("    ⊘ 跳过(重复):  %d 个\n", m.statistics.SkippedCount)))
	b.WriteString(errorStyle.Render(fmt.Sprintf("    ✗ 失败:        %d 个\n", m.statistics.FailedCount)))
	b.WriteString("\n")

	// 性能数据
	b.WriteString(labelStyle.Render("性能数据:"))
	b.WriteString("\n")
	b.WriteString(textStyle.Render(fmt.Sprintf("    耗时:          %s\n", m.statistics.Duration.Round(time.Second))))
	b.WriteString(textStyle.Render(fmt.Sprintf("    处理速度:      %.1f 文件/秒\n", m.statistics.GetSpeed())))
	b.WriteString("\n")

	// 详细日志
	b.WriteString(labelStyle.Render("💾 详细日志: "))
	b.WriteString(textStyle.Render(m.logFilePath))
	b.WriteString("\n\n")

	// 分割线 - 统一宽度计算
	dividerWidth := m.width - 8
	if dividerWidth > 0 {
		b.WriteString(renderDivider(dividerWidth))
		b.WriteString("\n\n")
	}

	// 提示 - 考虑终端宽度自动换行
	hintText := "按 [R] 重新整理  |  按 [O] 打开目标目录  |  按 [Q/Esc] 退出"
	maxHintWidth := m.width - 8
	if len(hintText) > maxHintWidth && maxHintWidth > 0 {
		hintText = "按 [R] 重新整理  |  按 [O] 打开目标目录\n按 [Q/Esc] 退出"
	}
	b.WriteString(hintStyle.Render(hintText))
	b.WriteString("\n")

	// 确保容器宽度合理
	containerWidth := m.width - 2
	if containerWidth < 40 {
		containerWidth = 40
	}

	return borderStyle.Width(containerWidth).Render(b.String())
}
