package screens

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/chiyiangel/media-organizer-v2/internal/organizer"
)

// ProgressScreen represents the progress screen
type ProgressScreen struct {
	app AppInterface
	
	// UI elements
	currentFileLabel  *widget.Label
	targetPathLabel   *widget.Label
	progressBar       *widget.ProgressBar
	progressLabel     *widget.Label
	statsLabel        *widget.RichText
	logEntry          *widget.Entry
	pauseButton       *widget.Button
	cancelButton      *widget.Button
	openFolderButton  *widget.Button
	
	// State
	processManager *organizer.ProcessManager
	stats          *organizer.Statistics
	isPaused       bool
	isCancelled    bool
	startTime      time.Time
	logBuffer      []string
	maxLogLines    int
}

// NewProgressScreen creates a new progress screen
func NewProgressScreen(app AppInterface) *ProgressScreen {
	screen := &ProgressScreen{
		app:         app,
		maxLogLines: 10,
		logBuffer:   make([]string, 0, 10),
	}
	
	screen.initUI()
	return screen
}

// initUI initializes the UI components
func (s *ProgressScreen) initUI() {
	s.currentFileLabel = widget.NewLabel("等待开始...")
	s.currentFileLabel.Wrapping = fyne.TextWrapWord
	
	s.targetPathLabel = widget.NewLabel("")
	s.targetPathLabel.Wrapping = fyne.TextWrapWord
	
	s.progressBar = widget.NewProgressBar()
	s.progressBar.Min = 0
	s.progressBar.Max = 100
	
	s.progressLabel = widget.NewLabel("0%")
	
	s.statsLabel = widget.NewRichText()
	s.statsLabel.ParseMarkdown("**已扫描:** 0 个文件")
	
	s.logEntry = widget.NewMultiLineEntry()
	s.logEntry.Disable()
	s.logEntry.SetPlaceHolder("处理日志将显示在这里...")
	
	s.pauseButton = widget.NewButton("暂停", func() {
		s.onPause()
	})
	
	s.cancelButton = widget.NewButton("取消", func() {
		s.onCancel()
	})
	s.cancelButton.Importance = widget.DangerImportance
	
	s.openFolderButton = widget.NewButton("查看目标文件夹", func() {
		s.onOpenFolder()
	})
}

// Title returns the screen title
func (s *ProgressScreen) Title() string {
	return "进度"
}

// Content returns the screen content
func (s *ProgressScreen) Content() fyne.CanvasObject {
	// Header section
	headerSection := container.NewVBox(
		widget.NewLabelWithStyle("🔄 正在处理文件...", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
		widget.NewLabel("当前文件:"),
		s.currentFileLabel,
		widget.NewLabel("目标路径:"),
		s.targetPathLabel,
		widget.NewLabel(""),
	)
	
	// Progress section
	progressBox := container.NewBorder(nil, nil, nil, s.progressLabel, s.progressBar)
	progressSection := container.NewVBox(
		widget.NewLabel("进度:"),
		progressBox,
		widget.NewSeparator(),
	)
	
	// Statistics section
	statsSection := container.NewVBox(
		widget.NewLabelWithStyle("📊 实时统计", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		s.statsLabel,
		widget.NewSeparator(),
	)
	
	// Log section
	logSection := container.NewVBox(
		widget.NewLabelWithStyle("📝 处理日志 (最近10条)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewScroll(s.logEntry),
	)
	
	// Button section
	buttonBox := container.NewHBox(
		s.pauseButton,
		s.cancelButton,
		s.openFolderButton,
	)
	
	// Main content
	content := container.NewBorder(
		container.NewVBox(
			headerSection,
			progressSection,
			statsSection,
		),
		buttonBox,
		nil,
		nil,
		logSection,
	)
	
	return content
}

// OnShow is called when the screen is shown
func (s *ProgressScreen) OnShow() {
	s.startTime = time.Now()
	s.isPaused = false
	s.isCancelled = false
	s.logBuffer = make([]string, 0, s.maxLogLines)
	
	// 创建处理器
	cfg := s.app.Config()
	scanner := organizer.NewScanner(cfg.SourceDir)
	processor := organizer.NewProcessor(cfg)
	s.processManager = organizer.NewProcessManager(scanner, processor, s)
	
	// 启动处理
	go func() {
		if err := s.processManager.Start(); err != nil {
			dialog.ShowError(err, s.app.Window())
		}
	}()
}

// OnHide is called when the screen is hidden
func (s *ProgressScreen) OnHide() {
	// Cleanup if needed
	if s.processManager != nil && !s.processManager.IsStopped() {
		s.processManager.Stop()
	}
}

// ProcessCallback 接口实现

// OnStart 处理开始时调用
func (s *ProgressScreen) OnStart(totalFiles int) {
	s.currentFileLabel.SetText("准备开始...")
	s.progressBar.Max = float64(totalFiles)
	s.addLog(fmt.Sprintf("开始处理，共 %d 个文件", totalFiles))
}

// OnProgress 处理进度更新时调用
func (s *ProgressScreen) OnProgress(current int, total int, file *organizer.FileInfo) {
	percent := float64(current) / float64(total)
	
	s.currentFileLabel.SetText(file.Name)
	s.targetPathLabel.SetText(file.TargetPath)
	s.progressBar.SetValue(float64(current))
	s.progressLabel.SetText(fmt.Sprintf("%d%%", int(percent*100)))
	
	// 更新统计信息
	if s.stats != nil {
		s.updateStats()
	}
	
	// 添加日志
	s.addLog(fmt.Sprintf("✓ %s → %s", file.Name, file.TargetPath))
}

// OnComplete 处理完成时调用
func (s *ProgressScreen) OnComplete(stats *organizer.Statistics) {
	s.stats = stats
	s.updateStats()
	
	// 保存到历史记录
	histMgr := s.app.History()
	if histMgr != nil {
		cfg := s.app.Config()
		if err := histMgr.Add(cfg, stats); err != nil {
			s.addLog(fmt.Sprintf("⚠️  保存历史记录失败: %v", err))
		}
	}
	
	// 发送系统通知
	successCount := stats.ProcessedFiles - stats.SkippedCount - stats.FailedCount
	notification := &Notification{
		Title:   "Media Organizer - 处理完成",
		Message: fmt.Sprintf("成功处理 %d 个文件", successCount),
		Sound:   true,
	}
	
	if err := sendNotification(notification); err != nil {
		s.addLog(fmt.Sprintf("⚠️  发送通知失败: %v", err))
	}
	
	s.showCompletionDialog()
}

// OnError 发生错误时调用
func (s *ProgressScreen) OnError(err error) {
	s.addLog(fmt.Sprintf("❌ 错误: %v", err))
	
	// 发送错误通知
	notification := &Notification{
		Title:   "Media Organizer - 处理错误",
		Message: fmt.Sprintf("错误: %v", err),
		Sound:   true,
	}
	
	if notifyErr := sendNotification(notification); notifyErr != nil {
		s.addLog(fmt.Sprintf("⚠️  发送通知失败: %v", notifyErr))
	}
	
	dialog.ShowError(err, s.app.Window())
}

// Helper methods

// addLog 添加日志条目
func (s *ProgressScreen) addLog(message string) {
	timestamp := time.Now().Format("15:04:05")
	logLine := fmt.Sprintf("[%s] %s", timestamp, message)
	
	s.logBuffer = append(s.logBuffer, logLine)
	if len(s.logBuffer) > s.maxLogLines {
		s.logBuffer = s.logBuffer[1:]
	}
	
	logText := ""
	for _, line := range s.logBuffer {
		logText += line + "\n"
	}
	s.logEntry.SetText(logText)
}

// updateStats 更新统计信息显示
func (s *ProgressScreen) updateStats() {
	if s.stats == nil {
		return
	}
	
	elapsed := s.stats.Duration
	if elapsed == 0 {
		elapsed = time.Since(s.startTime)
	}
	
	speed := float64(s.stats.ProcessedFiles) / elapsed.Seconds()
	remaining := 0
	if speed > 0 {
		remaining = int(float64(s.stats.TotalFiles-s.stats.ProcessedFiles) / speed)
	}
	
	statsText := fmt.Sprintf(`**已扫描:** %d 个文件
**已处理:** %d 个文件
  ├─ 📷 照片: %d 个
  ├─ 🎥 视频: %d 个
  ├─ ⏭️  跳过: %d 个 (重复)
  └─ ❌ 失败: %d 个

**剩余时间:** 约 %d 秒
**处理速度:** %.1f 个/秒`,
		s.stats.ScannedFiles,
		s.stats.ProcessedFiles,
		s.stats.PhotoCount,
		s.stats.VideoCount,
		s.stats.SkippedCount,
		s.stats.FailedCount,
		remaining,
		speed,
	)
	s.statsLabel.ParseMarkdown(statsText)
}

// onPause handles the pause button click
func (s *ProgressScreen) onPause() {
	if s.processManager == nil {
		return
	}
	
	if s.processManager.IsPaused() {
		s.processManager.Resume()
		s.pauseButton.SetText("暂停")
		s.addLog("继续处理...")
	} else {
		s.processManager.Pause()
		s.pauseButton.SetText("继续")
		s.addLog("已暂停")
	}
}

// onCancel handles the cancel button click
func (s *ProgressScreen) onCancel() {
	dialog.ShowConfirm(
		"确认取消",
		"确定要取消当前处理吗？",
		func(confirmed bool) {
			if confirmed {
				if s.processManager != nil {
					s.processManager.Stop()
				}
				s.isCancelled = true
				s.addLog("已取消处理")
				time.Sleep(500 * time.Millisecond)
				s.app.NavigateTo(0) // Go back to config screen
			}
		},
		s.app.Window(),
	)
}

// onOpenFolder handles the open folder button click
func (s *ProgressScreen) onOpenFolder() {
	cfg := s.app.Config()
	if cfg.TargetDir == "" {
		dialog.ShowError(fmt.Errorf("目标目录未设置"), s.app.Window())
		return
	}
	
	// macOS: use 'open' command
	// TODO: Add support for other platforms
	err := openFolder(cfg.TargetDir)
	if err != nil {
		dialog.ShowError(fmt.Errorf("打开文件夹失败: %v", err), s.app.Window())
	}
}

// showCompletionDialog shows the completion dialog
func (s *ProgressScreen) showCompletionDialog() {
	successMsg := fmt.Sprintf(`处理完成！

总文件数: %d
成功: %d
跳过: %d
失败: %d

耗时: %s`,
		s.stats.TotalFiles,
		s.stats.ProcessedFiles - s.stats.SkippedCount - s.stats.FailedCount,
		s.stats.SkippedCount,
		s.stats.FailedCount,
		s.stats.Duration.Round(time.Second),
	)
	
	dialog.ShowInformation(
		"处理完成",
		successMsg,
		s.app.Window(),
	)
	
	// Navigate to results/summary screen (we'll use history screen for now)
	s.app.NavigateTo(2)
}
