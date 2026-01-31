package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/chiyiangel/media-organizer-v2/internal/config"
	"github.com/chiyiangel/media-organizer-v2/internal/organizer"
)

// Record 历史记录
type Record struct {
	ID        string                 `json:"id"`
	Date      time.Time              `json:"date"`
	SourceDir string                 `json:"source_dir"`
	TargetDir string                 `json:"target_dir"`
	Config    map[string]interface{} `json:"config"`
	Stats     *organizer.Statistics  `json:"stats"`
}

// Manager 历史记录管理器
type Manager struct {
	storePath string
	records   []*Record
}

// NewManager 创建历史记录管理器
func NewManager(storePath string) (*Manager, error) {
	if storePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		storePath = filepath.Join(homeDir, ".media-organizer", "history.json")
	}
	
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(storePath), 0755); err != nil {
		return nil, err
	}
	
	manager := &Manager{
		storePath: storePath,
		records:   make([]*Record, 0),
	}
	
	// 加载历史记录
	if err := manager.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	
	return manager, nil
}

// Add 添加历史记录
func (m *Manager) Add(cfg *config.Config, stats *organizer.Statistics) error {
	record := &Record{
		ID:        time.Now().Format("20060102150405"),
		Date:      time.Now(),
		SourceDir: cfg.SourceDir,
		TargetDir: cfg.TargetDir,
		Config: map[string]interface{}{
			"duplicate_detection": cfg.DuplicateDetection,
			"duplicate_strategy":  cfg.DuplicateStrategy,
			"log_level":           cfg.LogLevel,
		},
		Stats: stats,
	}
	
	m.records = append([]*Record{record}, m.records...)
	
	// 限制历史记录数量（保留最近100条）
	if len(m.records) > 100 {
		m.records = m.records[:100]
	}
	
	return m.Save()
}

// GetAll 获取所有历史记录
func (m *Manager) GetAll() []*Record {
	return m.records
}

// GetByID 根据ID获取历史记录
func (m *Manager) GetByID(id string) *Record {
	for _, record := range m.records {
		if record.ID == id {
			return record
		}
	}
	return nil
}

// Delete 删除历史记录
func (m *Manager) Delete(id string) error {
	for i, record := range m.records {
		if record.ID == id {
			m.records = append(m.records[:i], m.records[i+1:]...)
			return m.Save()
		}
	}
	return nil
}

// Clear 清空所有历史记录
func (m *Manager) Clear() error {
	m.records = make([]*Record, 0)
	return m.Save()
}

// Save 保存历史记录到文件
func (m *Manager) Save() error {
	data, err := json.MarshalIndent(m.records, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(m.storePath, data, 0644)
}

// Load 从文件加载历史记录
func (m *Manager) Load() error {
	data, err := os.ReadFile(m.storePath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, &m.records)
}
