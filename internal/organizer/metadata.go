package organizer

import (
	"fmt"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

// MetadataExtractor 元数据提取器
type MetadataExtractor struct{}

// NewMetadataExtractor 创建元数据提取器
func NewMetadataExtractor() *MetadataExtractor {
	return &MetadataExtractor{}
}

// ExtractDate 提取日期
func (e *MetadataExtractor) ExtractDate(file *FileInfo) (time.Time, error) {
	switch file.Type {
	case FileTypePhoto:
		return e.extractPhotoDate(file.Path)
	case FileTypeVideo:
		return e.extractVideoDate(file.Path)
	default:
		return time.Time{}, fmt.Errorf("不支持的文件类型")
	}
}

// extractPhotoDate 提取照片日期
func (e *MetadataExtractor) extractPhotoDate(path string) (time.Time, error) {
	// 尝试读取EXIF
	f, err := os.Open(path)
	if err != nil {
		return e.getFileCreationTime(path)
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return e.getFileCreationTime(path)
	}

	// 尝试获取拍摄时间
	dateTime, err := x.DateTime()
	if err == nil {
		return dateTime, nil
	}

	// 尝试获取原始拍摄时间
	tag, err := x.Get(exif.DateTimeOriginal)
	if err == nil {
		if dateStr, err := tag.StringVal(); err == nil {
			if t, err := time.Parse("2006:01:02 15:04:05", dateStr); err == nil {
				return t, nil
			}
		}
	}

	// 回退到文件创建时间
	return e.getFileCreationTime(path)
}

// extractVideoDate 提取视频日期
func (e *MetadataExtractor) extractVideoDate(path string) (time.Time, error) {
	// 视频直接使用文件创建时间
	return e.getFileCreationTime(path)
}

// getFileCreationTime 获取文件创建时间
func (e *MetadataExtractor) getFileCreationTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}
