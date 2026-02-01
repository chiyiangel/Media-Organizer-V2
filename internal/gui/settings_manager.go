package gui

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Settings represents GUI settings
type Settings struct {
	WindowWidth  int    `json:"window_width"`
	WindowHeight int    `json:"window_height"`
	Theme        string `json:"theme"` // "自动 (跟随系统)", "明亮模式", "暗黑模式"
	
	// Notification settings
	NotifyOnComplete bool `json:"notify_on_complete"`
	NotifyOnError    bool `json:"notify_on_error"`
	
	// Update settings
	AutoCheckUpdate bool `json:"auto_check_update"`
}

// SettingsManager manages GUI settings
type SettingsManager struct {
	storePath string
	settings  *Settings
}

// NewSettingsManager creates a new settings manager
func NewSettingsManager(storePath string) (*SettingsManager, error) {
	if storePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		storePath = filepath.Join(homeDir, ".media-organizer", "gui-settings.json")
	}
	
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(storePath), 0755); err != nil {
		return nil, err
	}
	
	manager := &SettingsManager{
		storePath: storePath,
		settings:  getDefaultSettings(),
	}
	
	// 加载设置
	if err := manager.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	
	return manager, nil
}

// getDefaultSettings returns default settings
func getDefaultSettings() *Settings {
	return &Settings{
		WindowWidth:      900,
		WindowHeight:     600,
		Theme:            "自动 (跟随系统)",
		NotifyOnComplete: true,
		NotifyOnError:    true,
		AutoCheckUpdate:  true,
	}
}

// Get returns the current settings
func (m *SettingsManager) Get() *Settings {
	return m.settings
}

// Save saves settings to file
func (m *SettingsManager) Save() error {
	data, err := json.MarshalIndent(m.settings, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(m.storePath, data, 0644)
}

// Load loads settings from file
func (m *SettingsManager) Load() error {
	data, err := os.ReadFile(m.storePath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, m.settings)
}

// Reset resets settings to defaults
func (m *SettingsManager) Reset() error {
	m.settings = getDefaultSettings()
	return m.Save()
}
