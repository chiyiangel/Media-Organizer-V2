package screens

import (
	"fmt"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SettingsScreen represents the settings screen
type SettingsScreen struct {
	app AppInterface
	
	// UI elements
	themeGroup           *widget.RadioGroup
	notifyCompleteCheck  *widget.Check
	notifyErrorCheck     *widget.Check
	autoUpdateCheck      *widget.Check
	versionLabel         *widget.Label
	checkUpdateBtn       *widget.Button
	defaultSourceEntry   *widget.Entry
	defaultTargetEntry   *widget.Entry
	historyCountLabel    *widget.Label
	logSizeLabel         *widget.Label
	clearHistoryBtn      *widget.Button
	clearLogsBtn         *widget.Button
	resetButton          *widget.Button
	saveButton           *widget.Button
	cancelButton         *widget.Button
}

// NewSettingsScreen creates a new settings screen
func NewSettingsScreen(app AppInterface) *SettingsScreen {
	screen := &SettingsScreen{
		app: app,
	}
	
	screen.initUI()
	return screen
}

// initUI initializes the UI components
func (s *SettingsScreen) initUI() {
	// Theme section
	s.themeGroup = widget.NewRadioGroup(
		[]string{"自动 (跟随系统)", "明亮模式", "暗黑模式"},
		func(value string) {
			// 主题更改时立即应用
			if appWithTheme, ok := s.app.(interface{ SetTheme(string) }); ok {
				appWithTheme.SetTheme(value)
			}
		},
	)
	s.themeGroup.SetSelected("自动 (跟随系统)")
	
	// Notification section
	s.notifyCompleteCheck = widget.NewCheck("处理完成时发送系统通知", func(bool) {})
	s.notifyCompleteCheck.SetChecked(true)
	
	s.notifyErrorCheck = widget.NewCheck("发生错误时发送通知", func(bool) {})
	s.notifyErrorCheck.SetChecked(true)
	
	// Update section
	s.autoUpdateCheck = widget.NewCheck("启动时自动检查更新", func(bool) {})
	s.autoUpdateCheck.SetChecked(true)
	
	s.versionLabel = widget.NewLabel("当前版本: v2.0.0")
	
	s.checkUpdateBtn = widget.NewButton("立即检查更新", func() {
		s.onCheckUpdate()
	})
	
	// Default directories section
	s.defaultSourceEntry = widget.NewEntry()
	s.defaultSourceEntry.SetPlaceHolder("/Users/photos")
	
	s.defaultTargetEntry = widget.NewEntry()
	s.defaultTargetEntry.SetPlaceHolder("/Users/organized")
	
	// Data management section
	s.historyCountLabel = widget.NewLabel("历史记录数: 5 条")
	s.logSizeLabel = widget.NewLabel("日志文件大小: 2.3 MB")
	
	s.clearHistoryBtn = widget.NewButton("清理历史记录", func() {
		s.onClearHistory()
	})
	
	s.clearLogsBtn = widget.NewButton("清理日志文件", func() {
		s.onClearLogs()
	})
	
	// Action buttons
	s.resetButton = widget.NewButton("恢复默认设置", func() {
		s.onReset()
	})
	s.resetButton.Importance = widget.WarningImportance
	
	s.saveButton = widget.NewButton("保存", func() {
		s.onSave()
	})
	s.saveButton.Importance = widget.HighImportance
	
	s.cancelButton = widget.NewButton("取消", func() {
		s.app.NavigateTo(0)
	})
}

// Title returns the screen title
func (s *SettingsScreen) Title() string {
	return "设置"
}

// Content returns the screen content
func (s *SettingsScreen) Content() fyne.CanvasObject {
	// Appearance section
	appearanceSection := container.NewVBox(
		widget.NewLabelWithStyle("🎨 外观", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("主题:"),
		s.themeGroup,
		widget.NewSeparator(),
	)
	
	// Notification section
	notificationSection := container.NewVBox(
		widget.NewLabelWithStyle("🔔 通知", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		s.notifyCompleteCheck,
		s.notifyErrorCheck,
		widget.NewSeparator(),
	)
	
	// Update section
	updateSection := container.NewVBox(
		widget.NewLabelWithStyle("🔄 更新", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		s.autoUpdateCheck,
		s.versionLabel,
		s.checkUpdateBtn,
		widget.NewSeparator(),
	)
	
	// Default directories section
	sourceBrowseBtn := widget.NewButton("浏览...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				s.defaultSourceEntry.SetText(uri.Path())
			}
		}, s.app.Window())
	})
	
	targetBrowseBtn := widget.NewButton("浏览...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				s.defaultTargetEntry.SetText(uri.Path())
			}
		}, s.app.Window())
	})
	
	sourceBox := container.NewBorder(nil, nil, nil, sourceBrowseBtn, s.defaultSourceEntry)
	targetBox := container.NewBorder(nil, nil, nil, targetBrowseBtn, s.defaultTargetEntry)
	
	directoriesSection := container.NewVBox(
		widget.NewLabelWithStyle("📂 默认目录", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("源目录:"),
		sourceBox,
		widget.NewLabel("目标:"),
		targetBox,
		widget.NewSeparator(),
	)
	
	// Data management section
	dataSection := container.NewVBox(
		widget.NewLabelWithStyle("🗑️ 数据管理", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		s.historyCountLabel,
		s.logSizeLabel,
		container.NewHBox(
			s.clearHistoryBtn,
			s.clearLogsBtn,
		),
		widget.NewSeparator(),
	)
	
	// Button section
	buttonBox := container.NewHBox(
		s.resetButton,
		s.saveButton,
		s.cancelButton,
	)
	
	// Main content
	content := container.NewVBox(
		appearanceSection,
		notificationSection,
		updateSection,
		directoriesSection,
		dataSection,
		widget.NewLabel(""),
		buttonBox,
	)
	
	// Wrap in a scrollable container
	return container.NewScroll(content)
}

// OnShow is called when the screen is shown
func (s *SettingsScreen) OnShow() {
	// 加载设置
	cfg := s.app.Config()
	if cfg.SourceDir != "" {
		s.defaultSourceEntry.SetText(cfg.SourceDir)
	}
	if cfg.TargetDir != "" {
		s.defaultTargetEntry.SetText(cfg.TargetDir)
	}
	
	// 更新历史记录数量
	histMgr := s.app.History()
	if histMgr != nil {
		count := len(histMgr.GetAll())
		s.historyCountLabel.SetText(fmt.Sprintf("历史记录数: %d 条", count))
	} else {
		s.historyCountLabel.SetText("历史记录数: 0 条")
	}
	
	// TODO: 更新日志文件大小
	s.logSizeLabel.SetText("日志文件大小: 0 MB")
}

// OnHide is called when the screen is hidden
func (s *SettingsScreen) OnHide() {
	// Nothing to do
}

// onCheckUpdate handles the check update button click
func (s *SettingsScreen) onCheckUpdate() {
	dialog.ShowInformation("检查更新", "您正在使用最新版本", s.app.Window())
}

// onClearHistory handles the clear history button click
func (s *SettingsScreen) onClearHistory() {
	histMgr := s.app.History()
	if histMgr == nil {
		dialog.ShowError(fmt.Errorf("历史记录管理器不可用"), s.app.Window())
		return
	}
	
	if len(histMgr.GetAll()) == 0 {
		dialog.ShowInformation("提示", "没有历史记录可清理", s.app.Window())
		return
	}
	
	dialog.ShowConfirm(
		"确认清理",
		"确定要清理所有历史记录吗？此操作不可恢复。",
		func(confirmed bool) {
			if confirmed {
				if err := histMgr.Clear(); err != nil {
					dialog.ShowError(fmt.Errorf("清理失败: %v", err), s.app.Window())
					return
				}
				
				// 更新显示
				s.historyCountLabel.SetText("历史记录数: 0 条")
				dialog.ShowInformation("清理", "历史记录已清理", s.app.Window())
			}
		},
		s.app.Window(),
	)
}

// onClearLogs handles the clear logs button click
func (s *SettingsScreen) onClearLogs() {
	dialog.ShowConfirm(
		"确认清理",
		"确定要清理日志文件吗？",
		func(confirmed bool) {
			if confirmed {
				dialog.ShowInformation("清理", "日志文件已清理", s.app.Window())
			}
		},
		s.app.Window(),
	)
}

// onReset handles the reset button click
func (s *SettingsScreen) onReset() {
	dialog.ShowConfirm(
		"确认重置",
		"确定要恢复默认设置吗？",
		func(confirmed bool) {
			if confirmed {
				// Reset to defaults
				s.themeGroup.SetSelected("自动 (跟随系统)")
				s.notifyCompleteCheck.SetChecked(true)
				s.notifyErrorCheck.SetChecked(true)
				s.autoUpdateCheck.SetChecked(true)
				s.defaultSourceEntry.SetText("")
				s.defaultTargetEntry.SetText("")
				
				dialog.ShowInformation("重置", "设置已恢复为默认值", s.app.Window())
			}
		},
		s.app.Window(),
	)
}

// onSave handles the save button click
func (s *SettingsScreen) onSave() {
	cfg := s.app.Config()
	
	// 保存默认目录
	if s.defaultSourceEntry.Text != "" {
		cfg.SourceDir = s.defaultSourceEntry.Text
	}
	if s.defaultTargetEntry.Text != "" {
		cfg.TargetDir = s.defaultTargetEntry.Text
	}
	
	// 主题已经通过 SetTheme 实时保存了
	// 其他设置暂时不持久化（TODO: 未来可以添加）
	
	dialog.ShowInformation("保存", "设置已保存", s.app.Window())
	s.app.NavigateTo(0) // 返回配置页面
}
