package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/chiyiangel/media-organizer-v2/internal/config"
	"github.com/chiyiangel/media-organizer-v2/internal/gui/screens"
)

// App represents the GUI application
type App struct {
	fyneApp fyne.App
	window  fyne.Window
	config  *config.Config
	
	// UI components
	sidebar   *widget.List
	content   *fyne.Container
	currentScreen screens.Screen
}

// NewApp creates a new GUI application
func NewApp(cfg *config.Config) *App {
	a := app.New()
	w := a.NewWindow("Media Organizer V2")
	
	guiApp := &App{
		fyneApp: a,
		window:  w,
		config:  cfg,
	}
	
	guiApp.setupUI()
	return guiApp
}

// setupUI initializes the user interface
func (a *App) setupUI() {
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
	a.window.Resize(fyne.NewSize(900, 600))
	
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

// Window returns the main window
func (a *App) Window() fyne.Window {
	return a.window
}

// NavigateTo switches to a different screen (exported for screens package)
func (a *App) NavigateTo(screenID widget.ListItemID) {
	a.navigateTo(screenID)
}
