package organizer

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/yourusername/photo-video-organizer/internal/config"
)

// Processor 文件处理器
type Processor struct {
	config            *config.Config
	metadataExtractor *MetadataExtractor
	duplicateDetector *DuplicateDetector
}

// NewProcessor 创建处理器
func NewProcessor(cfg *config.Config) *Processor {
	return &Processor{
		config:            cfg,
		metadataExtractor: NewMetadataExtractor(),
		duplicateDetector: NewDuplicateDetector(cfg),
	}
}

// Process 处理文件
func (p *Processor) Process(file *FileInfo) (*ProcessRecord, error) {
	// 提取日期
	date, err := p.metadataExtractor.ExtractDate(file)
	if err != nil {
		return &ProcessRecord{
			File:    file,
			Result:  ResultFailed,
			Message: fmt.Sprintf("无法提取日期: %v", err),
		}, err
	}
	file.Date = date

	// 生成目标路径
	targetPath := p.generateTargetPath(file)
	file.TargetPath = targetPath

	// 检查重复
	isDuplicate, err := p.duplicateDetector.IsDuplicate(file)
	if err != nil {
		return &ProcessRecord{
			File:    file,
			Result:  ResultFailed,
			Message: fmt.Sprintf("检查重复失败: %v", err),
		}, err
	}

	if isDuplicate {
		switch p.config.DuplicateStrategy {
		case config.StrategySkip:
			return &ProcessRecord{
				File:    file,
				Result:  ResultSkipped,
				Message: "重复文件，已跳过",
			}, nil

		case config.StrategyOverwrite:
			// 继续处理，覆盖文件

		case config.StrategyRename:
			// 重命名文件
			targetPath = p.generateUniqueTargetPath(file)
			file.TargetPath = targetPath
		}
	}

	// 复制文件
	if err := p.copyFile(file.Path, file.TargetPath); err != nil {
		return &ProcessRecord{
			File:    file,
			Result:  ResultFailed,
			Message: fmt.Sprintf("复制文件失败: %v", err),
		}, err
	}

	return &ProcessRecord{
		File:    file,
		Result:  ResultSuccess,
		Message: "成功",
	}, nil
}

// generateTargetPath 生成目标路径
// 目录结构: YYYY/MM/MM-DD (月份-日期)
// 例如: 2025/10/10-01, 2025/10/10-25, 2025/12/12-31
func (p *Processor) generateTargetPath(file *FileInfo) string {
	year := file.Date.Format("2006")
	month := file.Date.Format("01")
	day := file.Date.Format("02")

	// 格式化为 MM-DD 格式
	monthDay := fmt.Sprintf("%s-%s", month, day)
	targetDir := filepath.Join(p.config.TargetDir, year, month, monthDay)
	return filepath.Join(targetDir, file.Name)
}

// generateUniqueTargetPath 生成唯一目标路径
func (p *Processor) generateUniqueTargetPath(file *FileInfo) string {
	basePath := p.generateTargetPath(file)
	ext := filepath.Ext(file.Name)
	nameWithoutExt := file.Name[:len(file.Name)-len(ext)]

	counter := 1
	newPath := basePath

	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		newPath = filepath.Join(
			filepath.Dir(basePath),
			fmt.Sprintf("%s(%d)%s", nameWithoutExt, counter, ext),
		)
		counter++
	}
}

// copyFile 复制文件
func (p *Processor) copyFile(src, dst string) error {
	// 创建目标目录
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制
	_, err = io.Copy(dstFile, srcFile)
	return err
}

// CalculateMD5 计算文件MD5
func CalculateMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
