package organizer

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/chiyiangel/media-organizer-v2/internal/config"
	"github.com/chiyiangel/media-organizer-v2/internal/i18n"
)

// copyBufferSize 复制文件时使用的缓冲区大小 (64KB)
const copyBufferSize = 64 * 1024

// Processor 文件处理器
type Processor struct {
	config            *config.Config
	metadataExtractor *MetadataExtractor
	duplicateDetector *DuplicateDetector
	createdDirs       map[string]bool // 缓存已创建的目录，避免重复调用 os.MkdirAll
}

// NewProcessor 创建处理器
func NewProcessor(cfg *config.Config) *Processor {
	return &Processor{
		config:            cfg,
		metadataExtractor: NewMetadataExtractor(),
		duplicateDetector: NewDuplicateDetector(cfg),
		createdDirs:       make(map[string]bool),
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
			Message: i18n.Tf("error.extract_date", err.Error()),
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
			Message: i18n.Tf("error.check_duplicate", err.Error()),
		}, err
	}

	if isDuplicate {
		switch p.config.DuplicateStrategy {
		case config.StrategySkip:
			return &ProcessRecord{
				File:    file,
				Result:  ResultSkipped,
				Message: i18n.T("message.duplicate_skipped"),
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
			Message: i18n.Tf("error.copy_file", err.Error()),
		}, err
	}

	return &ProcessRecord{
		File:    file,
		Result:  ResultSuccess,
		Message: i18n.T("message.success"),
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

// PreCreateDirectories 预创建目录（方案2优化）
// 在批量处理文件前，先收集所有需要创建的目录并批量创建
// 这样可以避免在处理每个文件时重复检查和创建目录
func (p *Processor) PreCreateDirectories(files []*FileInfo) error {
	// 收集所有唯一的目标目录
	dirs := make(map[string]bool)
	for _, file := range files {
		// 提取日期
		date, err := p.metadataExtractor.ExtractDate(file)
		if err != nil {
			continue // 跳过无法提取日期的文件
		}
		file.Date = date

		// 生成目标路径
		targetPath := p.generateTargetPath(file)
		dir := filepath.Dir(targetPath)
		dirs[dir] = true
	}

	// 批量创建目录
	for dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		p.createdDirs[dir] = true
	}

	return nil
}

// copyFile 复制文件（优化版本）
// 使用缓冲IO和目录缓存来提高小文件处理效率
func (p *Processor) copyFile(src, dst string) error {
	// 检查并创建目标目录（使用缓存避免重复调用）
	dir := filepath.Dir(dst)
	if !p.createdDirs[dir] {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		p.createdDirs[dir] = true
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

	// 使用缓冲写入器提高小文件写入效率
	bufferedWriter := bufio.NewWriterSize(dstFile, copyBufferSize)
	defer bufferedWriter.Flush()

	// 复制文件内容
	_, err = io.Copy(bufferedWriter, srcFile)
	if err != nil {
		return err
	}

	// 确保所有数据写入磁盘
	return bufferedWriter.Flush()
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
