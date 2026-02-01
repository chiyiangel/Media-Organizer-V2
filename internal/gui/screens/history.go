package screens

import (
	"fmt"
	"time"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/chiyiangel/media-organizer-v2/internal/history"
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
	historyData []*history.Record
	selectedID  int
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
				dateStr := record.Date.Format("2006-01-02 15:04")
				label.SetText(dateStr + " | " + record.SourceDir)
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
	histMgr := s.app.History()
	if histMgr == nil {
		// 历史记录管理器不可用
		s.historyData = make([]*history.Record, 0)
		return
	}
	
	s.historyData = histMgr.GetAll()
	s.historyList.Refresh()
}

// showDetails shows the details of a history record
func (s *HistoryScreen) showDetails(id widget.ListItemID) {
	if id < 0 || id >= len(s.historyData) {
		return
	}
	
	record := s.historyData[id]
	
	// 从配置中提取信息
	detection := "未知"
	if val, ok := record.Config["duplicate_detection"].(string); ok {
		detection = val
	}
	strategy := "未知"
	if val, ok := record.Config["duplicate_strategy"].(string); ok {
		strategy = val
	}
	logLevel := "info"
	if val, ok := record.Config["log_level"].(string); ok {
		logLevel = val
	}
	
	// 统计信息
	stats := record.Stats
	successCount := stats.ProcessedFiles - stats.SkippedCount - stats.FailedCount
	
	detailsMarkdown := fmt.Sprintf(`**详细信息 (%s)**

**源目录:** %s
**目标:** %s

**参数配置:**
  • 重复检测: %s
  • 处理策略: %s
  • 日志级别: %s

**处理结果:**
  • 总数: %d, 成功: %d, 跳过: %d, 失败: %d
  • 照片: %d, 视频: %d
  • 耗时: %s
  • 速度: %.1f 个/秒
`,
		record.Date.Format("2006-01-02 15:04:05"),
		record.SourceDir,
		record.TargetDir,
		detection,
		strategy,
		logLevel,
		stats.TotalFiles,
		successCount,
		stats.SkippedCount,
		stats.FailedCount,
		stats.PhotoCount,
		stats.VideoCount,
		stats.Duration.Round(time.Second),
		stats.GetSpeed(),
	)
	
	s.detailsText.ParseMarkdown(detailsMarkdown)
}

// onUseConfig handles the use config button click
func (s *HistoryScreen) onUseConfig() {
	if s.selectedID < 0 || s.selectedID >= len(s.historyData) {
		dialog.ShowInformation("提示", "请先选择一条历史记录", s.app.Window())
		return
	}
	
	record := s.historyData[s.selectedID]
	cfg := s.app.Config()
	
	// 应用配置
	cfg.SourceDir = record.SourceDir
	cfg.TargetDir = record.TargetDir
	
	// 从 map 中恢复配置
	if val, ok := record.Config["duplicate_detection"].(string); ok {
		cfg.DuplicateDetection = val
	}
	if val, ok := record.Config["duplicate_strategy"].(string); ok {
		cfg.DuplicateStrategy = val
	}
	if val, ok := record.Config["log_level"].(string); ok {
		cfg.LogLevel = val
	}
	
	dialog.ShowInformation("使用配置", "配置已应用，请前往配置页面查看", s.app.Window())
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
	if s.selectedID < 0 || s.selectedID >= len(s.historyData) {
		dialog.ShowInformation("提示", "请先选择一条历史记录", s.app.Window())
		return
	}
	
	record := s.historyData[s.selectedID]
	
	dialog.ShowConfirm(
		"确认删除",
		"确定要删除这条历史记录吗？",
		func(confirmed bool) {
			if confirmed {
				histMgr := s.app.History()
				if histMgr != nil {
					if err := histMgr.Delete(record.ID); err != nil {
						dialog.ShowError(err, s.app.Window())
						return
					}
				}
				
				// 重新加载列表
				s.loadHistory()
				s.selectedID = -1
				s.detailsText.ParseMarkdown("")
				
				dialog.ShowInformation("删除", "历史记录已删除", s.app.Window())
			}
		},
		s.app.Window(),
	)
}

// onClearAll handles the clear all button click
func (s *HistoryScreen) onClearAll() {
	if len(s.historyData) == 0 {
		dialog.ShowInformation("提示", "没有历史记录可清空", s.app.Window())
		return
	}
	
	dialog.ShowConfirm(
		"确认清空",
		"确定要清空所有历史记录吗？此操作不可恢复。",
		func(confirmed bool) {
			if confirmed {
				histMgr := s.app.History()
				if histMgr != nil {
					if err := histMgr.Clear(); err != nil {
						dialog.ShowError(err, s.app.Window())
						return
					}
				}
				
				// 重新加载列表
				s.loadHistory()
				s.selectedID = -1
				s.detailsText.ParseMarkdown("")
				
				dialog.ShowInformation("清空", "所有历史记录已清空", s.app.Window())
			}
		},
		s.app.Window(),
	)
}
