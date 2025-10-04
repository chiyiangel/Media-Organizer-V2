package config

import (
	"fmt"
	"os"
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
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() *Config {
	return &Config{
		SourceDir:          "",
		TargetDir:          "",
		DuplicateDetection: DetectionFilename,
		DuplicateStrategy:  StrategySkip,
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.SourceDir == "" {
		return fmt.Errorf("源目录不能为空")
	}
	if c.TargetDir == "" {
		return fmt.Errorf("目标目录不能为空")
	}
	// 检查目录是否存在
	if _, err := os.Stat(c.SourceDir); os.IsNotExist(err) {
		return fmt.Errorf("源目录不存在: %s", c.SourceDir)
	}
	return nil
}
