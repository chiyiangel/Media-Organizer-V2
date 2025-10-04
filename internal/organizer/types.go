package organizer

import "time"

// FileType 文件类型
type FileType string

const (
	FileTypePhoto FileType = "photo" // 照片
	FileTypeVideo FileType = "video" // 视频
	FileTypeOther FileType = "other" // 其他
)

// FileInfo 文件信息
type FileInfo struct {
	Path       string    // 文件路径
	Name       string    // 文件名
	Type       FileType  // 文件类型
	Size       int64     // 文件大小
	Date       time.Time // 日期（来自EXIF或创建时间）
	MD5        string    // MD5哈希（按需计算）
	TargetPath string    // 目标路径
}

// ProcessResult 处理结果
type ProcessResult string

const (
	ResultSuccess ProcessResult = "success" // 成功
	ResultSkipped ProcessResult = "skipped" // 跳过
	ResultFailed  ProcessResult = "failed"  // 失败
)

// ProcessRecord 处理记录
type ProcessRecord struct {
	File    *FileInfo     // 文件信息
	Result  ProcessResult // 处理结果
	Message string        // 消息（错误信息等）
}

// Statistics 统计信息
type Statistics struct {
	TotalFiles     int           // 总文件数
	ScannedFiles   int           // 已扫描文件数
	ProcessedFiles int           // 已处理文件数
	PhotoCount     int           // 照片数量
	VideoCount     int           // 视频数量
	SkippedCount   int           // 跳过数量
	FailedCount    int           // 失败数量
	StartTime      time.Time     // 开始时间
	EndTime        time.Time     // 结束时间
	Duration       time.Duration // 耗时
}

// GetSpeed 计算处理速度（文件/秒）
func (s *Statistics) GetSpeed() float64 {
	if s.Duration.Seconds() == 0 {
		return 0
	}
	return float64(s.ProcessedFiles) / s.Duration.Seconds()
}
