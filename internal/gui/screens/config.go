package screens

import (
	"fmt"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/chiyiangel/media-organizer-v2/internal/config"
)

// AppInterface defines the interface for accessing the main app
type AppInterface interface {
	Config() *config.Config
	Window() fyne.Window
	NavigateTo(id widget.ListItemID)
}

// ConfigScreen represents the configuration screen
type ConfigScreen struct {
	app AppInterface
	
	// UI elements
	sourceEntry     *widget.Entry
	targetEntry     *widget.Entry
	detectionGroup  *widget.RadioGroup
	strategyGroup   *widget.RadioGroup
	logLevelSelect  *widget.Select
	startButton     *widget.Button
	importButton    *widget.Button
	exportButton    *widget.Button
}

// NewConfigScreen creates a new configuration screen
func NewConfigScreen(app AppInterface) *ConfigScreen {
	screen := &ConfigScreen{
		app: app,
	}
	
	screen.initUI()
	return screen
}

// initUI initializes the UI components
func (s *ConfigScreen) initUI() {
	// Initialize components
	s.sourceEntry = widget.NewEntry()
	s.sourceEntry.SetPlaceHolder("/Users/photos")
	
	s.targetEntry = widget.NewEntry()
	s.targetEntry.SetPlaceHolder("/Users/organized")
	
	s.detectionGroup = widget.NewRadioGroup(
		[]string{"按文件名", "按MD5哈希"},
		func(value string) {},
	)
	s.detectionGroup.SetSelected("按文件名")
	s.detectionGroup.Horizontal = true
	
	s.strategyGroup = widget.NewRadioGroup(
		[]string{"跳过 (保留原文件)", "重命名 (添加序号)", "覆盖 (替换原文件)"},
		func(value string) {},
	)
	s.strategyGroup.SetSelected("重命名 (添加序号)")
	
	s.logLevelSelect = widget.NewSelect(
		[]string{"Debug", "Info", "Warning", "Error"},
		func(value string) {},
	)
	s.logLevelSelect.SetSelected("Info")
	
	s.startButton = widget.NewButton("开始整理", func() {
		s.onStart()
	})
	s.startButton.Importance = widget.HighImportance
	
	s.importButton = widget.NewButton("导入配置", func() {
		s.onImportConfig()
	})
	
	s.exportButton = widget.NewButton("导出配置", func() {
		s.onExportConfig()
	})
}

// Title returns the screen title
func (s *ConfigScreen) Title() string {
	return "配置"
}

// Content returns the screen content
func (s *ConfigScreen) Content() fyne.CanvasObject {
	// Directory selection section
	sourceBrowseBtn := widget.NewButton("浏览...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				s.sourceEntry.SetText(uri.Path())
			}
		}, s.app.Window())
	})
	
	targetBrowseBtn := widget.NewButton("浏览...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				s.targetEntry.SetText(uri.Path())
			}
		}, s.app.Window())
	})
	
	sourceBox := container.NewBorder(nil, nil, nil, sourceBrowseBtn, s.sourceEntry)
	targetBox := container.NewBorder(nil, nil, nil, targetBrowseBtn, s.targetEntry)
	
	dirSection := container.NewVBox(
		widget.NewLabelWithStyle("📁 目录设置", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
		widget.NewLabel("源目录:"),
		sourceBox,
		widget.NewLabel(""),
		widget.NewLabel("目标目录:"),
		targetBox,
		widget.NewSeparator(),
	)
	
	// Detection method section
	detectionSection := container.NewVBox(
		widget.NewLabelWithStyle("🔍 重复检测方式", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		s.detectionGroup,
		widget.NewSeparator(),
	)
	
	// Strategy section
	strategySection := container.NewVBox(
		widget.NewLabelWithStyle("📋 重复文件处理", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		s.strategyGroup,
		widget.NewSeparator(),
	)
	
	// Log level section
	logSection := container.NewVBox(
		widget.NewLabel("📝 日志级别:"),
		s.logLevelSelect,
		widget.NewSeparator(),
	)
	
	// Button section
	buttonBox := container.NewHBox(
		s.startButton,
		s.importButton,
		s.exportButton,
	)
	
	// Main content
	content := container.NewVBox(
		dirSection,
		detectionSection,
		strategySection,
		logSection,
		widget.NewLabel(""),
		buttonBox,
	)
	
	// Wrap in a scrollable container
	return container.NewScroll(content)
}

// OnShow is called when the screen is shown
func (s *ConfigScreen) OnShow() {
	// Load configuration into UI
	cfg := s.app.Config()
	if cfg.SourceDir != "" {
		s.sourceEntry.SetText(cfg.SourceDir)
	}
	if cfg.TargetDir != "" {
		s.targetEntry.SetText(cfg.TargetDir)
	}
	
	// Set detection method
	if cfg.DuplicateDetection == "md5" {
		s.detectionGroup.SetSelected("按MD5哈希")
	} else {
		s.detectionGroup.SetSelected("按文件名")
	}
	
	// Set strategy
	switch cfg.DuplicateStrategy {
	case "skip":
		s.strategyGroup.SetSelected("跳过 (保留原文件)")
	case "overwrite":
		s.strategyGroup.SetSelected("覆盖 (替换原文件)")
	default:
		s.strategyGroup.SetSelected("重命名 (添加序号)")
	}
}

// OnHide is called when the screen is hidden
func (s *ConfigScreen) OnHide() {
	// Save configuration from UI
	cfg := s.app.Config()
	cfg.SourceDir = s.sourceEntry.Text
	cfg.TargetDir = s.targetEntry.Text
	
	// Save detection method
	if s.detectionGroup.Selected == "按MD5哈希" {
		cfg.DuplicateDetection = "md5"
	} else {
		cfg.DuplicateDetection = "filename"
	}
	
	// Save strategy
	switch s.strategyGroup.Selected {
	case "跳过 (保留原文件)":
		cfg.DuplicateStrategy = "skip"
	case "覆盖 (替换原文件)":
		cfg.DuplicateStrategy = "overwrite"
	default:
		cfg.DuplicateStrategy = "rename"
	}
}

// onStart handles the start button click
func (s *ConfigScreen) onStart() {
	// Validate configuration
	if s.sourceEntry.Text == "" {
		dialog.ShowError(fmt.Errorf("请选择源目录"), s.app.Window())
		return
	}
	
	if s.targetEntry.Text == "" {
		dialog.ShowError(fmt.Errorf("请选择目标目录"), s.app.Window())
		return
	}
	
	// Save current configuration
	s.OnHide()
	
	// Switch to progress screen
	s.app.NavigateTo(1)
}

// onImportConfig handles the import config button click
func (s *ConfigScreen) onImportConfig() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, s.app.Window())
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()
		
		// TODO: Implement config import logic
		dialog.ShowInformation("导入配置", "配置导入功能开发中...", s.app.Window())
	}, s.app.Window())
}

// onExportConfig handles the export config button click
func (s *ConfigScreen) onExportConfig() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, s.app.Window())
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()
		
		// TODO: Implement config export logic
		dialog.ShowInformation("导出配置", "配置导出功能开发中...", s.app.Window())
	}, s.app.Window())
}
