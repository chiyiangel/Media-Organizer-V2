package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/chiyiangel/media-organizer-v2/internal/config"
	"github.com/chiyiangel/media-organizer-v2/internal/gui/screens"
	"github.com/chiyiangel/media-organizer-v2/internal/history"
)

// App represents the GUI application
type App struct {
	fyneApp fyne.App
	window  fyne.Window
	config  *config.Config
	history *history.Manager
	settings *SettingsManager
	
	// UI components
	sidebar   *widget.List
	content   *fyne.Container
	currentScreen screens.Screen
}

// NewApp creates a new GUI application
func NewApp(cfg *config.Config) *App {
	a := app.New()
	w := a.NewWindow("Media Organizer V2")
	
	// 创建历史记录管理器
	historyMgr, err := history.NewManager("")
	if err != nil {
		// 如果创建失败，使用nil，界面会优雅降级
		historyMgr = nil
	}
	
	// 创建设置管理器
	settingsMgr, err := NewSettingsManager("")
	if err != nil {
		// 如果创建失败，使用默认设置
		settingsMgr = nil
	}
	
	guiApp := &App{
		fyneApp:  a,
		window:   w,
		config:   cfg,
		history:  historyMgr,
		settings: settingsMgr,
	}
	
	guiApp.setupUI()
	return guiApp
}

// SetTheme sets the application theme
func (a *App) SetTheme(themeName string) {
	switch themeName {
	case "明亮模式":
		a.fyneApp.Settings().SetTheme(&customLightTheme{})
	case "暗黑模式":
		a.fyneApp.Settings().SetTheme(&customDarkTheme{})
	default: // "自动 (跟随系统)"
		a.fyneApp.Settings().SetTheme(nil) // Use system default
	}
	
	// 保存主题设置
	if a.settings != nil {
		settings := a.settings.Get()
		settings.Theme = themeName
		a.settings.Save()
	}
}

// setupUI initializes the user interface
func (a *App) setupUI() {
	// 应用保存的设置
	if a.settings != nil {
		settings := a.settings.Get()
		
		// 应用主题
		a.SetTheme(settings.Theme)
		
		// 应用窗口大小
		a.window.Resize(fyne.NewSize(float32(settings.WindowWidth), float32(settings.WindowHeight)))
	} else {
		// 使用默认大小
		a.window.Resize(fyne.NewSize(900, 600))
	}
	
	// Create sidebar navigation
	a.sidebar = widget.NewList(
		func() int { return 4 }, // Number of items
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			switch id {
			case 0:
				label.SetText("🏠 配置")
			case 1:
				label.SetText("📊 进度")
			case 2:
				label.SetText("📜 历史")
			case 3:
				label.SetText("⚙️  设置")
			}
		},
	)
	
	// Handle sidebar selection
	a.sidebar.OnSelected = func(id widget.ListItemID) {
		a.navigateTo(id)
	}
	
	// Create initial content area
	a.content = container.NewStack()
	
	// Create main layout with sidebar and content
	split := container.NewHSplit(
		container.NewBorder(nil, nil, nil, nil, a.sidebar),
		a.content,
	)
	split.SetOffset(0.2) // 20% for sidebar
	
	// Set window content
	a.window.SetContent(split)
	
	// 保存窗口关闭时的设置
	a.window.SetOnClosed(func() {
		if a.settings != nil {
			settings := a.settings.Get()
			size := a.window.Canvas().Size()
			settings.WindowWidth = int(size.Width)
			settings.WindowHeight = int(size.Height)
			a.settings.Save()
		}
	})
	
	// Show config screen by default
	a.navigateTo(0)
}

// navigateTo switches to a different screen
func (a *App) navigateTo(screenID widget.ListItemID) {
	if a.currentScreen != nil {
		a.currentScreen.OnHide()
	}
	
	var newScreen screens.Screen
	
	switch screenID {
	case 0:
		newScreen = screens.NewConfigScreen(a)
	case 1:
		newScreen = screens.NewProgressScreen(a)
	case 2:
		newScreen = screens.NewHistoryScreen(a)
	case 3:
		newScreen = screens.NewSettingsScreen(a)
	default:
		newScreen = screens.NewConfigScreen(a)
	}
	
	a.currentScreen = newScreen
	a.content.Objects = []fyne.CanvasObject{newScreen.Content()}
	a.content.Refresh()
	
	newScreen.OnShow()
}

// Run starts the GUI application
func (a *App) Run() {
	a.window.ShowAndRun()
}

// Config returns the application configuration
func (a *App) Config() *config.Config {
	return a.config
}

// History returns the history manager
func (a *App) History() *history.Manager {
	return a.history
}

// Settings returns the settings manager
func (a *App) Settings() *SettingsManager {
	return a.settings
}

// Window returns the main window
func (a *App) Window() fyne.Window {
	return a.window
}

// NavigateTo switches to a different screen (exported for screens package)
func (a *App) NavigateTo(screenID widget.ListItemID) {
	a.navigateTo(screenID)
}
