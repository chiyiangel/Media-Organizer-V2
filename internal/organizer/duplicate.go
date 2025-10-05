package organizer

import (
	"os"

	"github.com/chiyiangel/media-organizer-v2/internal/config"
)

// DuplicateDetector 重复文件检测器
type DuplicateDetector struct {
	config *config.Config
}

// NewDuplicateDetector 创建检测器
func NewDuplicateDetector(cfg *config.Config) *DuplicateDetector {
	return &DuplicateDetector{
		config: cfg,
	}
}

// IsDuplicate 检查是否重复
func (d *DuplicateDetector) IsDuplicate(file *FileInfo) (bool, error) {
	// 检查目标文件是否存在
	if _, err := os.Stat(file.TargetPath); os.IsNotExist(err) {
		return false, nil
	}

	switch d.config.DuplicateDetection {
	case config.DetectionFilename:
		// 文件名模式：文件存在即为重复
		return true, nil

	case config.DetectionMD5:
		// MD5模式：比较文件内容
		return d.compareByMD5(file)

	default:
		return false, nil
	}
}

// compareByMD5 通过MD5比较
func (d *DuplicateDetector) compareByMD5(file *FileInfo) (bool, error) {
	// 计算源文件MD5
	if file.MD5 == "" {
		md5, err := CalculateMD5(file.Path)
		if err != nil {
			return false, err
		}
		file.MD5 = md5
	}

	// 计算目标文件MD5
	targetMD5, err := CalculateMD5(file.TargetPath)
	if err != nil {
		return false, err
	}

	return file.MD5 == targetMD5, nil
}
