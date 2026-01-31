package organizer

import (
	"fmt"
	"time"
)

// ProcessCallback 定义处理过程中的回调接口
type ProcessCallback interface {
	OnStart(totalFiles int)
	OnProgress(current int, total int, file *FileInfo)
	OnComplete(stats *Statistics)
	OnError(err error)
}

// ProcessManager 管理文件处理流程
type ProcessManager struct {
	scanner   *Scanner
	processor *Processor
	callback  ProcessCallback
	
	// Control
	stopChan  chan bool
	pauseChan chan bool
	isPaused  bool
	isStopped bool
}

// NewProcessManager 创建处理管理器
func NewProcessManager(scanner *Scanner, processor *Processor, callback ProcessCallback) *ProcessManager {
	return &ProcessManager{
		scanner:   scanner,
		processor: processor,
		callback:  callback,
		stopChan:  make(chan bool),
		pauseChan: make(chan bool),
	}
}

// Start 开始处理
func (pm *ProcessManager) Start() error {
	// 扫描文件
	files, err := pm.scanner.Scan()
	if err != nil {
		pm.callback.OnError(err)
		return err
	}
	
	totalFiles := len(files)
	if totalFiles == 0 {
		return fmt.Errorf("未找到需要处理的文件")
	}
	
	// 通知开始
	pm.callback.OnStart(totalFiles)
	
	// 初始化统计信息
	stats := &Statistics{
		TotalFiles:   totalFiles,
		StartTime:    time.Now(),
	}
	
	// 处理文件
	for i, file := range files {
		// 检查是否停止
		select {
		case <-pm.stopChan:
			pm.isStopped = true
			stats.EndTime = time.Now()
			stats.Duration = stats.EndTime.Sub(stats.StartTime)
			pm.callback.OnComplete(stats)
			return nil
		default:
		}
		
		// 检查暂停
		for pm.isPaused {
			select {
			case <-pm.pauseChan:
				pm.isPaused = false
			case <-pm.stopChan:
				pm.isStopped = true
				stats.EndTime = time.Now()
				stats.Duration = stats.EndTime.Sub(stats.StartTime)
				pm.callback.OnComplete(stats)
				return nil
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
		
		// 更新扫描数
		stats.ScannedFiles = i + 1
		
		// 处理文件
		record, err := pm.processor.Process(file)
		if err == nil && record != nil {
			stats.ProcessedFiles++
			
			// 更新统计
			switch file.Type {
			case FileTypePhoto:
				stats.PhotoCount++
			case FileTypeVideo:
				stats.VideoCount++
			}
			
			switch record.Result {
			case ResultSkipped:
				stats.SkippedCount++
			case ResultFailed:
				stats.FailedCount++
			}
		} else {
			stats.FailedCount++
		}
		
		// 更新进度
		stats.Duration = time.Since(stats.StartTime)
		pm.callback.OnProgress(i+1, totalFiles, file)
	}
	
	// 完成
	stats.EndTime = time.Now()
	stats.Duration = stats.EndTime.Sub(stats.StartTime)
	pm.callback.OnComplete(stats)
	
	return nil
}

// Pause 暂停处理
func (pm *ProcessManager) Pause() {
	pm.isPaused = true
}

// Resume 恢复处理
func (pm *ProcessManager) Resume() {
	if pm.isPaused {
		pm.pauseChan <- true
	}
}

// Stop 停止处理
func (pm *ProcessManager) Stop() {
	pm.stopChan <- true
}

// IsPaused 是否已暂停
func (pm *ProcessManager) IsPaused() bool {
	return pm.isPaused
}

// IsStopped 是否已停止
func (pm *ProcessManager) IsStopped() bool {
	return pm.isStopped
}
