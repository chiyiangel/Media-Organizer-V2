package screens

import (
	"fmt"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// HistoryScreen represents the history screen
type HistoryScreen struct {
	app AppInterface
	
	// UI elements
	searchEntry    *widget.Entry
	historyList    *widget.List
	detailsText    *widget.RichText
	useConfigBtn   *widget.Button
	viewLogBtn     *widget.Button
	deleteBtn      *widget.Button
	clearAllBtn    *widget.Button
	
	// Data
	historyData []HistoryRecord
	selectedID  int
}

// HistoryRecord represents a history record
type HistoryRecord struct {
	Date      string
	SourceDir string
	TargetDir string
	Success   int
	Failed    int
	Duration  string
}

// NewHistoryScreen creates a new history screen
func NewHistoryScreen(app AppInterface) *HistoryScreen {
	screen := &HistoryScreen{
		app:        app,
		selectedID: -1,
	}
	
	screen.initUI()
	screen.loadHistory()
	return screen
}

// initUI initializes the UI components
func (s *HistoryScreen) initUI() {
	s.searchEntry = widget.NewEntry()
	s.searchEntry.SetPlaceHolder("搜索...")
	
	s.historyList = widget.NewList(
		func() int {
			return len(s.historyData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(s.historyData) {
				record := s.historyData[id]
				label := obj.(*widget.Label)
				label.SetText(record.Date + " | " + record.SourceDir)
			}
		},
	)
	
	s.historyList.OnSelected = func(id widget.ListItemID) {
		s.selectedID = id
		s.showDetails(id)
	}
	
	s.detailsText = widget.NewRichText()
	
	s.useConfigBtn = widget.NewButton("使用此配置", func() {
		s.onUseConfig()
	})
	
	s.viewLogBtn = widget.NewButton("查看日志", func() {
		s.onViewLog()
	})
	
	s.deleteBtn = widget.NewButton("删除记录", func() {
		s.onDelete()
	})
	s.deleteBtn.Importance = widget.DangerImportance
	
	s.clearAllBtn = widget.NewButton("清空全部历史", func() {
		s.onClearAll()
	})
	s.clearAllBtn.Importance = widget.DangerImportance
}

// Title returns the screen title
func (s *HistoryScreen) Title() string {
	return "历史记录"
}

// Content returns the screen content
func (s *HistoryScreen) Content() fyne.CanvasObject {
	// Search header
	searchBox := container.NewBorder(
		nil, nil,
		widget.NewLabel("🔍"),
		nil,
		s.searchEntry,
	)
	
	// List section
	listSection := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("📜 历史记录", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			searchBox,
		),
		nil, nil, nil,
		s.historyList,
	)
	
	// Details section
	detailsSection := container.NewBorder(
		widget.NewLabelWithStyle("详细信息", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(
			s.useConfigBtn,
			s.viewLogBtn,
			s.deleteBtn,
			s.clearAllBtn,
		),
		nil, nil,
		container.NewScroll(s.detailsText),
	)
	
	// Split view
	split := container.NewHSplit(
		listSection,
		detailsSection,
	)
	split.SetOffset(0.4)
	
	return split
}

// OnShow is called when the screen is shown
func (s *HistoryScreen) OnShow() {
	s.loadHistory()
}

// OnHide is called when the screen is hidden
func (s *HistoryScreen) OnHide() {
	// Nothing to do
}

// loadHistory loads the history records
func (s *HistoryScreen) loadHistory() {
	// TODO: Load from actual storage
	// For now, use dummy data
	s.historyData = []HistoryRecord{
		{
			Date:      "2024-03-15 17:08",
			SourceDir: "/Users/photos",
			TargetDir: "/Users/organized",
			Success:   295,
			Failed:    1,
			Duration:  "3m45s",
		},
		{
			Date:      "2024-03-10 14:30",
			SourceDir: "/Users/photos",
			TargetDir: "/Users/organized",
			Success:   180,
			Failed:    0,
			Duration:  "2m10s",
		},
		{
			Date:      "2024-03-05 09:15",
			SourceDir: "/Users/camera",
			TargetDir: "/Users/backup",
			Success:   520,
			Failed:    3,
			Duration:  "8m30s",
		},
	}
	
	s.historyList.Refresh()
}

// showDetails shows the details of a history record
func (s *HistoryScreen) showDetails(id widget.ListItemID) {
	if id < 0 || id >= len(s.historyData) {
		return
	}
	
	record := s.historyData[id]
	
	detailsMarkdown := fmt.Sprintf(`**详细信息 (%s)**

**源目录:** %s
**目标:** %s

**参数配置:**
  • 重复检测: MD5
  • 处理策略: 重命名
  • 日志级别: Info

**处理结果:**
  • 总数: %d, 成功: %d, 失败: %d
  • 耗时: %s

**失败文件:**
  • corrupt.jpg (无法读取EXIF)
`,
		record.Date,
		record.SourceDir,
		record.TargetDir,
		record.Success+record.Failed,
		record.Success,
		record.Failed,
		record.Duration,
	)
	
	s.detailsText.ParseMarkdown(detailsMarkdown)
}

// onUseConfig handles the use config button click
func (s *HistoryScreen) onUseConfig() {
	if s.selectedID < 0 {
		dialog.ShowInformation("提示", "请先选择一条历史记录", s.app.Window())
		return
	}
	
	dialog.ShowInformation("使用配置", "配置已应用到配置页面", s.app.Window())
	s.app.NavigateTo(0)
}

// onViewLog handles the view log button click
func (s *HistoryScreen) onViewLog() {
	if s.selectedID < 0 {
		dialog.ShowInformation("提示", "请先选择一条历史记录", s.app.Window())
		return
	}
	
	dialog.ShowInformation("查看日志", "查看日志功能开发中...", s.app.Window())
}

// onDelete handles the delete button click
func (s *HistoryScreen) onDelete() {
	if s.selectedID < 0 {
		dialog.ShowInformation("提示", "请先选择一条历史记录", s.app.Window())
		return
	}
	
	dialog.ShowConfirm(
		"确认删除",
		"确定要删除这条历史记录吗？",
		func(confirmed bool) {
			if confirmed {
				// TODO: Implement delete logic
				dialog.ShowInformation("删除", "历史记录已删除", s.app.Window())
			}
		},
		s.app.Window(),
	)
}

// onClearAll handles the clear all button click
func (s *HistoryScreen) onClearAll() {
	dialog.ShowConfirm(
		"确认清空",
		"确定要清空所有历史记录吗？此操作不可恢复。",
		func(confirmed bool) {
			if confirmed {
				// TODO: Implement clear all logic
				s.historyData = []HistoryRecord{}
				s.historyList.Refresh()
				s.detailsText.ParseMarkdown("")
				dialog.ShowInformation("清空", "所有历史记录已清空", s.app.Window())
			}
		},
		s.app.Window(),
	)
}
