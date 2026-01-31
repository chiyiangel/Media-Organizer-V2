package screens

import "fyne.io/fyne/v2"

// Screen interface for different views
type Screen interface {
	Title() string
	Content() fyne.CanvasObject
	OnShow()
	OnHide()
}
