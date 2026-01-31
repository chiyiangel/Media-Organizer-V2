package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// ExportConfig 导出配置到文件
func ExportConfig(cfg *Config, filepath string) error {
	// 创建导出配置（只包含用户可配置的字段）
	exportCfg := map[string]interface{}{
		"source_dir":          cfg.SourceDir,
		"target_dir":          cfg.TargetDir,
		"duplicate_detection": cfg.DuplicateDetection,
		"duplicate_strategy":  cfg.DuplicateStrategy,
		"log_level":           cfg.LogLevel,
	}
	
	// 序列化为 JSON
	data, err := json.MarshalIndent(exportCfg, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	
	// 写入文件
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	
	return nil
}

// ImportConfig 从文件导入配置
func ImportConfig(filepath string) (*Config, error) {
	// 读取文件
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	
	// 解析 JSON
	var rawConfig map[string]interface{}
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	
	// 创建配置对象并填充
	cfg := NewDefaultConfig()
	
	if val, ok := rawConfig["source_dir"].(string); ok {
		cfg.SourceDir = val
	}
	if val, ok := rawConfig["target_dir"].(string); ok {
		cfg.TargetDir = val
	}
	if val, ok := rawConfig["duplicate_detection"].(string); ok {
		cfg.DuplicateDetection = DuplicateDetection(val)
	}
	if val, ok := rawConfig["duplicate_strategy"].(string); ok {
		cfg.DuplicateStrategy = DuplicateStrategy(val)
	}
	if val, ok := rawConfig["log_level"].(string); ok {
		cfg.LogLevel = val
	}
	
	// 验证导入的配置
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}
	
	return cfg, nil
}

// SaveConfig 保存配置到默认配置文件
func SaveConfig(cfg *Config) error {
	if cfg.ConfigFile == "" {
		cfg.ConfigFile = "media-organizer.json"
	}
	return ExportConfig(cfg, cfg.ConfigFile)
}

// LoadConfig 从默认配置文件加载配置
func LoadConfig(filepath string) (*Config, error) {
	return ImportConfig(filepath)
}
