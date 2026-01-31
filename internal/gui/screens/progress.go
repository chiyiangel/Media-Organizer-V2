package screens

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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
	isPaused    bool
	isCancelled bool
	startTime   time.Time
}

// NewProgressScreen creates a new progress screen
func NewProgressScreen(app AppInterface) *ProgressScreen {
	screen := &ProgressScreen{
		app: app,
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
	
	// Start processing
	go s.startProcessing()
}

// OnHide is called when the screen is hidden
func (s *ProgressScreen) OnHide() {
	// Cleanup if needed
}

// startProcessing starts the file processing
func (s *ProgressScreen) startProcessing() {
	// Simulate processing for demo
	// TODO: Replace with actual processing logic
	
	totalFiles := 100
	
	for i := 1; i <= totalFiles; i++ {
		if s.isCancelled {
			break
		}
		
		// Check if paused
		for s.isPaused {
			time.Sleep(100 * time.Millisecond)
			if s.isCancelled {
				return
			}
		}
		
		// Update UI
		s.updateProgress(i, totalFiles, fmt.Sprintf("IMG_%04d.jpg", i))
		
		// Simulate processing time
		time.Sleep(50 * time.Millisecond)
	}
	
	if !s.isCancelled {
		// Processing complete, show summary
		s.showCompletionDialog()
	}
}

// updateProgress updates the progress display
func (s *ProgressScreen) updateProgress(current, total int, filename string) {
	percent := float64(current) / float64(total)
	
	s.currentFileLabel.SetText(filename)
	s.targetPathLabel.SetText(fmt.Sprintf("/Users/organized/2024/01/01-15/%s", filename))
	s.progressBar.SetValue(percent)
	s.progressLabel.SetText(fmt.Sprintf("%d%%", int(percent*100)))
	
	// Update statistics
	statsText := fmt.Sprintf(`**已扫描:** %d 个文件
**已处理:** %d 个文件
  ├─ 📷 照片: %d 个
  ├─ 🎥 视频: %d 个
  ├─ ⏭️  跳过: %d 个 (重复)
  └─ ❌ 失败: %d 个

**剩余时间:** 约 %d 秒
**处理速度:** %.1f 个/秒`,
		total,
		current,
		int(float64(current) * 0.7),
		int(float64(current) * 0.25),
		int(float64(current) * 0.04),
		int(float64(current) * 0.01),
		(total-current)/10,
		10.0,
	)
	s.statsLabel.ParseMarkdown(statsText)
	
	// Add log entry
	timestamp := time.Now().Format("15:04:05")
	logLine := fmt.Sprintf("[%s] ✓ %s → 2024/01/01-15/\n", timestamp, filename)
	s.logEntry.SetText(s.logEntry.Text + logLine)
}

// onPause handles the pause button click
func (s *ProgressScreen) onPause() {
	s.isPaused = !s.isPaused
	
	if s.isPaused {
		s.pauseButton.SetText("继续")
	} else {
		s.pauseButton.SetText("暂停")
	}
}

// onCancel handles the cancel button click
func (s *ProgressScreen) onCancel() {
	dialog.ShowConfirm(
		"确认取消",
		"确定要取消当前处理吗？",
		func(confirmed bool) {
			if confirmed {
				s.isCancelled = true
				s.app.NavigateTo(0) // Go back to config screen
			}
		},
		s.app.Window(),
	)
}

// onOpenFolder handles the open folder button click
func (s *ProgressScreen) onOpenFolder() {
	// TODO: Implement open folder logic
	dialog.ShowInformation("打开文件夹", "打开文件夹功能开发中...", s.app.Window())
}

// showCompletionDialog shows the completion dialog
func (s *ProgressScreen) showCompletionDialog() {
	dialog.ShowInformation(
		"处理完成",
		"文件整理已完成！\n\n点击确定查看结果详情。",
		s.app.Window(),
	)
	
	// Navigate to results/summary screen (we'll use history screen for now)
	s.app.NavigateTo(2)
}
