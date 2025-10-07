package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chiyiangel/media-organizer-v2/internal/organizer"
)

// Logger 日志记录器
type Logger struct {
	file *os.File
	path string
}

// NewLogger 创建日志记录器
func NewLogger() (*Logger, error) {
	// 生成日志文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("organize_log_%s.txt", timestamp)

	// 创建日志文件
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		file: file,
		path: filename,
	}

	// 写入头部信息
	logger.writeHeader()

	return logger, nil
}

// GetPath 获取日志文件路径
func (l *Logger) GetPath() string {
	absPath, _ := filepath.Abs(l.path)
	return absPath
}

// writeHeader 写入日志头部
func (l *Logger) writeHeader() {
	header := fmt.Sprintf("照片视频整理日志 - %s\n", time.Now().Format("2006-01-02 15:04:05"))
	header += "=" + string(make([]byte, 60)) + "\n\n"
	l.file.WriteString(header)
}

// LogRecord 记录处理记录
func (l *Logger) LogRecord(record *organizer.ProcessRecord) {
	timestamp := time.Now().Format("15:04:05")
	var status string

	switch record.Result {
	case organizer.ResultSuccess:
		status = "✓ 成功"
	case organizer.ResultSkipped:
		status = "⊘ 跳过"
	case organizer.ResultFailed:
		status = "✗ 失败"
	}

	line := fmt.Sprintf("[%s] %s | %s -> %s | %s\n",
		timestamp,
		status,
		record.File.Name,
		record.File.TargetPath,
		record.Message,
	)

	l.file.WriteString(line)
}

// LogStatistics 记录统计信息
func (l *Logger) LogStatistics(stats *organizer.Statistics) {
	summary := "\n" + "=" + string(make([]byte, 60)) + "\n"
	summary += "整理完成汇总\n"
	summary += "=" + string(make([]byte, 60)) + "\n\n"

	summary += fmt.Sprintf("文件统计:\n")
	summary += fmt.Sprintf("  总文件数:     %d 个\n", stats.TotalFiles)
	summary += fmt.Sprintf("  ├─ 照片:      %d 张\n", stats.PhotoCount)
	summary += fmt.Sprintf("  └─ 视频:      %d 个\n\n", stats.VideoCount)

	summary += fmt.Sprintf("处理结果:\n")
	summary += fmt.Sprintf("  ✓ 成功整理:   %d 个\n", stats.ProcessedFiles-stats.SkippedCount-stats.FailedCount)
	summary += fmt.Sprintf("  ⊘ 跳过(重复): %d 个\n", stats.SkippedCount)
	summary += fmt.Sprintf("  ✗ 失败:       %d 个\n\n", stats.FailedCount)

	summary += fmt.Sprintf("性能数据:\n")
	summary += fmt.Sprintf("  开始时间:     %s\n", stats.StartTime.Format("2006-01-02 15:04:05"))
	summary += fmt.Sprintf("  结束时间:     %s\n", stats.EndTime.Format("2006-01-02 15:04:05"))
	summary += fmt.Sprintf("  耗时:         %s\n", stats.Duration.Round(time.Second))
	summary += fmt.Sprintf("  处理速度:     %.2f 文件/秒\n", stats.GetSpeed())

	l.file.WriteString(summary)
}

// LogError 记录错误信息
func (l *Logger) LogError(message string) {
	timestamp := time.Now().Format("15:04:05")
	line := fmt.Sprintf("[%s] ERROR | %s\n", timestamp, message)
	l.file.WriteString(line)
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
