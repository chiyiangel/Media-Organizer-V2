package config

import (
	"fmt"
	"os"
)

// OperationMode defines how the application runs
type OperationMode string

const (
	ModeInteractive OperationMode = "interactive" // Default TUI mode
	ModeSilent      OperationMode = "silent"      // Non-interactive mode
)

// ConfigSource defines where configuration comes from
type ConfigSource string

const (
	SourceDefault ConfigSource = "default" // Default values
	SourceFile    ConfigSource = "file"    // Configuration file
	SourceCLI     ConfigSource = "cli"     // Command line arguments
)

// DuplicateDetection 重复文件识别策略
type DuplicateDetection string

const (
	DetectionFilename DuplicateDetection = "filename" // 文件名
	DetectionMD5      DuplicateDetection = "md5"      // MD5哈希
)

// DuplicateStrategy 重复文件处理策略
type DuplicateStrategy string

const (
	StrategySkip      DuplicateStrategy = "skip"      // 跳过
	StrategyOverwrite DuplicateStrategy = "overwrite" // 覆盖
	StrategyRename    DuplicateStrategy = "rename"    // 重命名
)

// Config 应用配置
type Config struct {
	SourceDir          string             // 源目录
	TargetDir          string             // 目标目录
	DuplicateDetection DuplicateDetection // 重复识别策略
	DuplicateStrategy  DuplicateStrategy  // 重复处理策略

	// New fields for silent mode and configuration management
	Mode       OperationMode // Operation mode (interactive/silent)
	ConfigFile string        // Path to configuration file
	LogLevel   string        // Log level for silent mode
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() *Config {
	return &Config{
		SourceDir:          "",
		TargetDir:          "",
		DuplicateDetection: DetectionFilename,
		DuplicateStrategy:  StrategySkip,
		Mode:               ModeInteractive,
		ConfigFile:         "",
		LogLevel:           "info",
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	// Validate operation mode
	if c.Mode != ModeInteractive && c.Mode != ModeSilent {
		return fmt.Errorf("无效的操作模式: %s", c.Mode)
	}

	// Mode-specific validation
	if c.Mode == ModeSilent {
		// In silent mode, source and target directories are required
		if c.SourceDir == "" {
			return fmt.Errorf("在静默模式下，源目录不能为空")
		}
		if c.TargetDir == "" {
			return fmt.Errorf("在静默模式下，目标目录不能为空")
		}
	} else {
		// In interactive mode, directories are optional (can be set in TUI)
		if c.SourceDir != "" && c.TargetDir != "" {
			// If both are provided, validate them
			if _, err := os.Stat(c.SourceDir); os.IsNotExist(err) {
				return fmt.Errorf("源目录不存在: %s", c.SourceDir)
			}
		}
	}

	// Validate config file path if specified
	if c.ConfigFile != "" {
		if _, err := os.Stat(c.ConfigFile); os.IsNotExist(err) {
			return fmt.Errorf("配置文件不存在: %s", c.ConfigFile)
		}
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug":   true,
		"info":    true,
		"warning": true,
		"error":   true,
	}
	if c.LogLevel != "" && !validLogLevels[c.LogLevel] {
		return fmt.Errorf("无效的日志级别: %s (有效值: debug, info, warning, error)", c.LogLevel)
	}

	return nil
}
