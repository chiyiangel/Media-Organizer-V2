package organizer

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	// 照片扩展名（小写）
	photoExtensions = []string{".arw", ".jpg", ".jpeg", ".png", ".heic", ".gif", ".bmp", ".raw"}

	// 视频扩展名（小写）
	videoExtensions = []string{".mp4", ".mov", ".avi", ".mkv", ".flv", ".wmv"}
)

// Scanner 文件扫描器
type Scanner struct {
	sourceDir string
}

// NewScanner 创建扫描器
func NewScanner(sourceDir string) *Scanner {
	return &Scanner{
		sourceDir: sourceDir,
	}
}

// Scan 扫描文件
func (s *Scanner) Scan() ([]*FileInfo, error) {
	var files []*FileInfo

	err := filepath.Walk(s.sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查文件类型
		fileType := getFileType(path)
		if fileType == FileTypeOther {
			return nil
		}

		// 创建文件信息
		fileInfo := &FileInfo{
			Path: path,
			Name: info.Name(),
			Type: fileType,
			Size: info.Size(),
		}

		files = append(files, fileInfo)
		return nil
	})

	return files, err
}

// getFileType 获取文件类型
func getFileType(path string) FileType {
	ext := strings.ToLower(filepath.Ext(path))

	for _, photoExt := range photoExtensions {
		if ext == photoExt {
			return FileTypePhoto
		}
	}

	for _, videoExt := range videoExtensions {
		if ext == videoExt {
			return FileTypeVideo
		}
	}

	return FileTypeOther
}
