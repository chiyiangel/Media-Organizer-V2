package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/yourusername/photo-video-organizer/internal/config"
	"github.com/yourusername/photo-video-organizer/internal/i18n"
)

// renderConfigScreen 渲染主配置界面
func (m Model) renderConfigScreen() string {
	var b strings.Builder

	// 标题
	b.WriteString(titleStyle.Width(m.width).Render(i18n.T("app.title")))
	b.WriteString("\n\n")

	// 源目录
	b.WriteString(labelStyle.Render(i18n.T("config.source_dir")))
	if m.config.SourceDir == "" {
		b.WriteString(hintStyle.Render(i18n.T("config.not_set")))
	} else {
		b.WriteString(textStyle.Render(m.config.SourceDir))
	}
	b.WriteString("\n")
	b.WriteString(hintStyle.Render(i18n.T("config.edit_source_hint")))
	b.WriteString("\n\n")

	// 目标目录
	b.WriteString(labelStyle.Render(i18n.T("config.target_dir")))
	if m.config.TargetDir == "" {
		b.WriteString(hintStyle.Render(i18n.T("config.not_set")))
	} else {
		b.WriteString(textStyle.Render(m.config.TargetDir))
	}
	b.WriteString("\n")
	b.WriteString(hintStyle.Render(i18n.T("config.edit_target_hint")))
	b.WriteString("\n\n")

	// 整理策略
	b.WriteString(labelStyle.Render(i18n.T("config.organize_strategy")))
	b.WriteString("\n")

	// 同文件识别
	detectionF := " "
	detectionM := " "
	if m.config.DuplicateDetection == config.DetectionFilename {
		detectionF = "●"
	} else {
		detectionM = "●"
	}

	b.WriteString(textStyle.Render(i18n.Tf("config.file_detection", detectionF, detectionM)))
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

	b.WriteString(textStyle.Render(i18n.Tf("config.duplicate_handling", strategy1, strategy2, strategy3)))
	b.WriteString("\n\n")

	// 分割线 - 调整宽度以匹配边框
	dividerWidth := m.width - 8 // 考虑边框和内边距
	if dividerWidth > 0 {
		b.WriteString(renderDivider(dividerWidth))
		b.WriteString("\n\n")
	}

	// 错误信息
	if m.err != nil {
		b.WriteString(errorStyle.Render(i18n.T("error.prefix") + m.err.Error()))
		b.WriteString("\n\n")
	}

	// 提示 - 确保文字不超出容器宽度
	hintText := i18n.T("config.start_hint")
	maxHintWidth := m.width - 8 // 考虑边框和内边距
	if len(hintText) > maxHintWidth && maxHintWidth > 0 {
		hintText = i18n.T("config.start_hint_wrapped")
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
	title := i18n.T("input.title")
	b.WriteString(titleStyle.Width(m.width).Render(title))
	b.WriteString("\n\n")

	// 提示
	prompt := i18n.T("input.prompt")
	b.WriteString(labelStyle.Render(prompt))
	b.WriteString("\n")

	// 输入框
	inputBox := fmt.Sprintf("│ %s_", m.inputValue)
	b.WriteString(textStyle.Render(inputBox))
	b.WriteString("\n\n")

	// 提示
	b.WriteString(hintStyle.Render(i18n.T("input.confirm_hint")))
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
	b.WriteString(titleStyle.Width(m.width).Render(i18n.T("progress.title")))
	b.WriteString("\n\n")

	// 当前文件
	if m.currentFile != nil {
		b.WriteString(labelStyle.Render(i18n.T("progress.current_file")))
		b.WriteString(textStyle.Render(m.currentFile.Name))
		b.WriteString("\n")

		b.WriteString(labelStyle.Render(i18n.T("progress.target_path")))
		b.WriteString(textStyle.Render(m.currentFile.TargetPath))
		b.WriteString("\n\n")
	}

	// 进度条
	progressBarWidth := m.width - 20
	if progressBarWidth < 20 {
		progressBarWidth = 20
	}
	b.WriteString(labelStyle.Render(i18n.T("progress.progress")))
	b.WriteString(renderProgressBar(m.statistics.ProcessedFiles, m.statistics.TotalFiles, progressBarWidth))
	b.WriteString("\n\n")

	// 分割线 - 统一宽度计算
	dividerWidth := m.width - 8
	if dividerWidth > 0 {
		b.WriteString(renderDivider(dividerWidth))
		b.WriteString("\n\n")
	}

	// 实时统计
	b.WriteString(labelStyle.Render(i18n.T("progress.statistics")))
	b.WriteString("\n")
	b.WriteString(textStyle.Render(i18n.Tf("progress.scanned", m.statistics.ScannedFiles) + "\n"))
	b.WriteString(textStyle.Render(i18n.Tf("progress.processed", m.statistics.ProcessedFiles) + "\n"))
	b.WriteString(textStyle.Render(i18n.Tf("progress.photos", m.statistics.PhotoCount) + "\n"))
	b.WriteString(textStyle.Render(i18n.Tf("progress.videos", m.statistics.VideoCount) + "\n"))
	b.WriteString(textStyle.Render(i18n.Tf("progress.skipped", m.statistics.SkippedCount) + "\n"))
	b.WriteString(textStyle.Render(i18n.Tf("progress.errors", m.statistics.FailedCount) + "\n"))
	b.WriteString("\n")

	// 分割线
	if dividerWidth > 0 {
		b.WriteString(renderDivider(dividerWidth))
		b.WriteString("\n\n")
	}

	// 提示
	b.WriteString(hintStyle.Render(i18n.T("progress.cancel_hint")))
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
	b.WriteString(titleStyle.Width(m.width).Render(i18n.T("summary.title")))
	b.WriteString("\n\n")

	// 汇总报告标题
	b.WriteString(labelStyle.Render(i18n.T("summary.report_title")))
	b.WriteString("\n")
	b.WriteString(renderDivider(m.width - 4))
	b.WriteString("\n\n")

	// 文件统计
	b.WriteString(labelStyle.Render(i18n.T("summary.file_statistics")))
	b.WriteString("\n")
	b.WriteString(textStyle.Render(i18n.Tf("summary.total_files", m.statistics.TotalFiles) + "\n"))
	b.WriteString(textStyle.Render(i18n.Tf("summary.total_photos", m.statistics.PhotoCount) + "\n"))
	b.WriteString(textStyle.Render(i18n.Tf("summary.total_videos", m.statistics.VideoCount) + "\n"))
	b.WriteString("\n")

	// 处理结果
	successCount := m.statistics.ProcessedFiles - m.statistics.SkippedCount - m.statistics.FailedCount
	b.WriteString(labelStyle.Render(i18n.T("summary.process_results")))
	b.WriteString("\n")
	b.WriteString(successStyle.Render(i18n.Tf("summary.success", successCount) + "\n"))
	b.WriteString(warningStyle.Render(i18n.Tf("summary.skipped", m.statistics.SkippedCount) + "\n"))
	b.WriteString(errorStyle.Render(i18n.Tf("summary.failed", m.statistics.FailedCount) + "\n"))
	b.WriteString("\n")

	// 性能数据
	b.WriteString(labelStyle.Render(i18n.T("summary.performance")))
	b.WriteString("\n")
	b.WriteString(textStyle.Render(i18n.Tf("summary.duration", m.statistics.Duration.Round(time.Second).String()) + "\n"))
	b.WriteString(textStyle.Render(i18n.Tf("summary.speed", fmt.Sprintf("%.1f", m.statistics.GetSpeed())) + "\n"))
	b.WriteString("\n")

	// 详细日志
	b.WriteString(labelStyle.Render(i18n.T("summary.log_file")))
	b.WriteString(textStyle.Render(m.logFilePath))
	b.WriteString("\n\n")

	// 分割线 - 统一宽度计算
	dividerWidth := m.width - 8
	if dividerWidth > 0 {
		b.WriteString(renderDivider(dividerWidth))
		b.WriteString("\n\n")
	}

	// 提示 - 考虑终端宽度自动换行
	hintText := i18n.T("summary.actions_hint")
	maxHintWidth := m.width - 8
	if len(hintText) > maxHintWidth && maxHintWidth > 0 {
		hintText = i18n.T("summary.actions_hint_wrapped")
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
