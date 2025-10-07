package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chiyiangel/media-organizer-v2/internal/config"
	"github.com/chiyiangel/media-organizer-v2/internal/logger"
	"github.com/chiyiangel/media-organizer-v2/internal/organizer"
)

// SilentRunner handles non-interactive execution of media organization
type SilentRunner struct {
	config    *config.Config
	logger    *logger.Logger
	processor *organizer.Processor
}

// NewSilentRunner creates a new silent mode runner
func NewSilentRunner(config *config.Config) (*SilentRunner, error) {
	// Create logger
	log, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Create processor with configuration
	processor := organizer.NewProcessor(config)

	return &SilentRunner{
		config:    config,
		logger:    log,
		processor: processor,
	}, nil
}

// Run executes the media organization in silent mode
func (r *SilentRunner) Run() error {
	fmt.Println("开始静默模式媒体整理")
	fmt.Printf("源目录: %s\n", r.config.SourceDir)
	fmt.Printf("目标目录: %s\n", r.config.TargetDir)
	fmt.Printf("重复识别策略: %s\n", r.config.DuplicateDetection)
	fmt.Printf("重复处理策略: %s\n", r.config.DuplicateStrategy)
	fmt.Println()

	// Set up interrupt handling
	r.handleInterrupt()

	// Start processing
	startTime := time.Now()
	fmt.Println("开始扫描文件...")

	// Create scanner and scan files
	scanner := organizer.NewScanner(r.config.SourceDir)
	files, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("文件扫描失败: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("未找到支持的媒体文件")
		return nil
	}

	fmt.Printf("找到 %d 个媒体文件，开始处理...\n", len(files))

	// Initialize statistics
	stats := &organizer.Statistics{
		TotalFiles:     len(files),
		ProcessedFiles: 0,
		PhotoCount:     0,
		VideoCount:     0,
		SkippedCount:   0,
		FailedCount:    0,
		StartTime:      startTime,
	}

	// Process each file
	var records []organizer.ProcessRecord
	for i, file := range files {
		// Update statistics based on file type
		if file.Type == organizer.FileTypePhoto {
			stats.PhotoCount++
		} else if file.Type == organizer.FileTypeVideo {
			stats.VideoCount++
		}

		// Process the file
		record, err := r.processor.Process(file)
		if err != nil {
			r.logger.LogError(fmt.Sprintf("处理文件失败: %s, 错误: %v", file.Path, err))
		}

		if record != nil {
			records = append(records, *record)

			// Update statistics based on result
			switch record.Result {
			case organizer.ResultSuccess:
				// Success count is calculated as ProcessedFiles - SkippedCount - FailedCount
			case organizer.ResultSkipped:
				stats.SkippedCount++
			case organizer.ResultFailed:
				stats.FailedCount++
				r.logger.LogError(fmt.Sprintf("文件处理失败: %s, 原因: %s", file.Path, record.Message))
			}
		}

		stats.ProcessedFiles++

		// Print progress every 10 files or on last file
		if (i+1)%10 == 0 || i == len(files)-1 {
			r.printProgress(stats)
		}
	}

	// Finalize statistics
	stats.EndTime = time.Now()
	stats.Duration = time.Since(startTime)

	fmt.Println() // Add newline after progress

	// Log all processing records to file
	for _, record := range records {
		r.logger.LogRecord(&record)
	}

	// Print final summary
	r.printSummary(stats)

	elapsed := time.Since(startTime)
	fmt.Printf("处理完成，耗时: %v\n", elapsed)
	fmt.Printf("详细日志已保存到: %s\n", r.logger.GetPath())

	// Log statistics to file
	r.logger.LogStatistics(stats)

	// Close logger
	r.logger.Close()

	return nil
}

// printProgress displays progress updates
func (r *SilentRunner) printProgress(stats *organizer.Statistics) {
	if stats.TotalFiles == 0 {
		return
	}

	percentage := float64(stats.ProcessedFiles) / float64(stats.TotalFiles) * 100
	successful := stats.ProcessedFiles - stats.SkippedCount - stats.FailedCount
	fmt.Printf("\r进度: %d/%d (%.1f%%) | 成功: %d | 失败: %d | 跳过: %d",
		stats.ProcessedFiles, stats.TotalFiles, percentage,
		successful, stats.FailedCount, stats.SkippedCount)
}

// printSummary displays final summary
func (r *SilentRunner) printSummary(stats *organizer.Statistics) {
	fmt.Println("\n")
	fmt.Println("=== 处理摘要 ===")
	fmt.Printf("总文件数: %d\n", stats.TotalFiles)
	fmt.Printf("照片数量: %d\n", stats.PhotoCount)
	fmt.Printf("视频数量: %d\n", stats.VideoCount)
	fmt.Printf("成功处理: %d\n", stats.ProcessedFiles-stats.SkippedCount-stats.FailedCount)
	fmt.Printf("处理失败: %d\n", stats.FailedCount)
	fmt.Printf("跳过文件: %d\n", stats.SkippedCount)

	if stats.FailedCount > 0 {
		fmt.Println("\n注意: 有文件处理失败，请查看日志文件了解详情")
	}

	fmt.Printf("\n重复文件处理策略: %s\n", r.config.DuplicateStrategy)
}

// handleInterrupt sets up graceful interrupt handling
func (r *SilentRunner) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\n\n接收到中断信号，正在停止...")
		// TODO: Implement proper stop mechanism when processor supports it
		os.Exit(1)
	}()
}
