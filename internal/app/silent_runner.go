package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chiyiangel/media-organizer-v2/internal/config"
	"github.com/chiyiangel/media-organizer-v2/internal/i18n"
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
	fmt.Println(i18n.T("silent.start"))
	fmt.Println(i18n.Tf("silent.source_dir", r.config.SourceDir))
	fmt.Println(i18n.Tf("silent.target_dir", r.config.TargetDir))
	fmt.Println(i18n.Tf("silent.duplicate_detection", r.config.DuplicateDetection))
	fmt.Println(i18n.Tf("silent.duplicate_strategy", r.config.DuplicateStrategy))
	fmt.Println()

	// Set up interrupt handling
	r.handleInterrupt()

	// Start processing
	startTime := time.Now()
	fmt.Println(i18n.T("silent.scan_start"))

	// Create scanner and scan files
	scanner := organizer.NewScanner(r.config.SourceDir)
	files, err := scanner.Scan()
	if err != nil {
		errorMsg := i18n.Tf("silent.scan_failed", err.Error())
		return fmt.Errorf(errorMsg)
	}

	if len(files) == 0 {
		fmt.Println(i18n.T("silent.no_media_files"))
		return nil
	}

	fmt.Println(i18n.Tf("silent.files_found", len(files)))

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
			r.logger.LogError(i18n.Tf("silent.process_file_failed", file.Path, err))
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
				r.logger.LogError(i18n.Tf("silent.file_process_failed", file.Path, record.Message))
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
	fmt.Println(i18n.Tf("silent.completed", elapsed))
	fmt.Println(i18n.Tf("silent.log_saved", r.logger.GetPath()))

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
	progressMsg := i18n.Tf("silent.progress",
		stats.ProcessedFiles, stats.TotalFiles, fmt.Sprintf("%.1f", percentage),
		successful, stats.FailedCount, stats.SkippedCount)
	fmt.Print("\r" + progressMsg)
}

// printSummary displays final summary
func (r *SilentRunner) printSummary(stats *organizer.Statistics) {
	fmt.Println()
	fmt.Println(i18n.T("silent.summary_title"))
	fmt.Println(i18n.Tf("silent.total_files", stats.TotalFiles))
	fmt.Println(i18n.Tf("silent.photo_count", stats.PhotoCount))
	fmt.Println(i18n.Tf("silent.video_count", stats.VideoCount))
	fmt.Println(i18n.Tf("silent.success_count", stats.ProcessedFiles-stats.SkippedCount-stats.FailedCount))
	fmt.Println(i18n.Tf("silent.failed_count", stats.FailedCount))
	fmt.Println(i18n.Tf("silent.skipped_count", stats.SkippedCount))

	if stats.FailedCount > 0 {
		fmt.Println("\n" + i18n.T("silent.failed_notice"))
	}

	fmt.Println("\n" + i18n.Tf("silent.strategy_used", string(r.config.DuplicateStrategy)))
}

// handleInterrupt sets up graceful interrupt handling
func (r *SilentRunner) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\n\n" + i18n.T("silent.interrupt_received"))
		// TODO: Implement proper stop mechanism when processor supports it
		os.Exit(1)
	}()
}
