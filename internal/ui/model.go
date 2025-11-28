package ui

import (
	"os/exec"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chiyiangel/media-organizer-v2/internal/config"
	"github.com/chiyiangel/media-organizer-v2/internal/logger"
	"github.com/chiyiangel/media-organizer-v2/internal/organizer"
)

// Screen 界面类型
type Screen int

const (
	ScreenConfig   Screen = iota // 主配置界面
	ScreenInput                  // 路径输入弹窗
	ScreenProgress               // 整理进度界面
	ScreenSummary                // 完成汇总界面
)

// InputMode 输入模式
type InputMode int

const (
	InputNone   InputMode = iota // 无输入
	InputSource                  // 输入源目录
	InputTarget                  // 输入目标目录
)

// Model Bubble Tea 主模型
type Model struct {
	// 当前界面
	currentScreen Screen

	// 配置
	config *config.Config

	// 输入状态
	inputMode   InputMode
	inputValue  string
	inputCursor int

	// 整理状态
	isOrganizing bool
	currentFile  *organizer.FileInfo
	statistics   *organizer.Statistics
	records      []organizer.ProcessRecord
	allFiles     []*organizer.FileInfo

	// 处理器（复用实例，避免每个文件都创建新的处理器）
	processor *organizer.Processor

	// 日志记录器
	logger      *logger.Logger
	logFilePath string

	// 窗口尺寸
	width  int
	height int

	// 错误信息
	err error

	// 取消标志
	cancelled bool
}

// NewModel 创建新模型
func NewModel() Model {
	return Model{
		currentScreen: ScreenConfig,
		config:        config.NewDefaultConfig(),
		inputMode:     InputNone,
		statistics:    &organizer.Statistics{},
		records:       make([]organizer.ProcessRecord, 0),
		width:         80,
		height:        24,
	}
}

// Init 初始化
func (m Model) Init() tea.Cmd {
	return nil
}

// Update 更新模型
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case FileScanCompleteMsg:
		return m.handleFileScanComplete(msg)

	case FileProcessedMsg:
		return m.handleFileProcessed(msg)

	case OrganizeCompleteMsg:
		return m.handleOrganizeComplete(msg)

	case OrganizeErrorMsg:
		m.err = msg.Err
		m.isOrganizing = false
		return m, nil

	case ProgressUpdateMsg:
		if m.currentScreen == ScreenProgress {
			m.statistics.ProcessedFiles = msg.Current
			m.statistics.TotalFiles = msg.Total
		}
		return m, nil
	}

	return m, nil
}

// View 渲染视图
func (m Model) View() string {
	switch m.currentScreen {
	case ScreenConfig:
		return m.renderConfigScreen()
	case ScreenInput:
		return m.renderInputScreen()
	case ScreenProgress:
		return m.renderProgressScreen()
	case ScreenSummary:
		return m.renderSummaryScreen()
	default:
		return ""
	}
}

// handleKeyPress 处理按键
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 全局快捷键
	switch msg.String() {
	case "ctrl+c":
		if m.logger != nil {
			m.logger.Close()
		}
		return m, tea.Quit
	case "esc":
		return m.handleEscape()
	}

	// 根据当前界面处理按键
	switch m.currentScreen {
	case ScreenConfig:
		return m.handleConfigKeys(msg)
	case ScreenInput:
		return m.handleInputKeys(msg)
	case ScreenProgress:
		return m.handleProgressKeys(msg)
	case ScreenSummary:
		return m.handleSummaryKeys(msg)
	}

	return m, nil
}

// handleEscape 处理Esc键
func (m Model) handleEscape() (tea.Model, tea.Cmd) {
	switch m.currentScreen {
	case ScreenInput:
		m.currentScreen = ScreenConfig
		m.inputMode = InputNone
		return m, nil
	case ScreenProgress:
		m.cancelled = true
		m.currentScreen = ScreenConfig
		m.isOrganizing = false
		return m, nil
	case ScreenConfig, ScreenSummary:
		if m.logger != nil {
			m.logger.Close()
		}
		return m, tea.Quit
	}
	return m, nil
}

// handleFileScanComplete 处理文件扫描完成
func (m Model) handleFileScanComplete(msg FileScanCompleteMsg) (tea.Model, tea.Cmd) {
	m.allFiles = msg.Files

	// 初始化统计信息
	m.statistics = &organizer.Statistics{
		StartTime:    time.Now(),
		TotalFiles:   len(msg.Files),
		ScannedFiles: len(msg.Files),
	}

	// 统计照片和视频数量
	for _, file := range msg.Files {
		if file.Type == organizer.FileTypePhoto {
			m.statistics.PhotoCount++
		} else if file.Type == organizer.FileTypeVideo {
			m.statistics.VideoCount++
		}
	}

	// 预创建目录（方案2优化）
	if m.processor != nil && len(msg.Files) > 0 {
		_ = m.processor.PreCreateDirectories(msg.Files)
	}

	// 切换到进度界面
	m.currentScreen = ScreenProgress
	m.isOrganizing = true

	// 开始处理第一个文件
	return m, m.processNextFileCmd(0)
}

// handleFileProcessed 处理文件处理完成
func (m Model) handleFileProcessed(msg FileProcessedMsg) (tea.Model, tea.Cmd) {
	m.records = append(m.records, *msg.Record)
	m.currentFile = msg.Record.File

	// 更新统计
	m.statistics.ProcessedFiles++

	switch msg.Record.Result {
	case organizer.ResultSuccess:
		// 成功计数已在ProcessedFiles中
	case organizer.ResultSkipped:
		m.statistics.SkippedCount++
	case organizer.ResultFailed:
		m.statistics.FailedCount++
	}

	// 记录到日志
	if m.logger != nil {
		m.logger.LogRecord(msg.Record)
	}

	// 处理下一个文件
	return m, m.processNextFileCmd(msg.FileIndex + 1)
}

// handleOrganizeComplete 处理整理完成
func (m Model) handleOrganizeComplete(msg OrganizeCompleteMsg) (tea.Model, tea.Cmd) {
	m.statistics = msg.Statistics
	m.logFilePath = msg.LogPath
	m.isOrganizing = false
	m.currentScreen = ScreenSummary

	// 记录统计到日志
	if m.logger != nil {
		m.logger.LogStatistics(m.statistics)
		m.logger.Close()
	}

	return m, nil
}

// startOrganizing 开始整理
func (m Model) startOrganizing() (tea.Model, tea.Cmd) {
	m.isOrganizing = true
	m.currentScreen = ScreenProgress
	m.statistics = &organizer.Statistics{
		StartTime: time.Now(),
	}
	m.records = make([]organizer.ProcessRecord, 0)
	m.cancelled = false

	// 创建处理器实例（复用，避免每个文件都创建新实例）
	m.processor = organizer.NewProcessor(m.config)

	// 创建日志记录器
	log, err := logger.NewLogger()
	if err != nil {
		m.err = err
		m.isOrganizing = false
		m.currentScreen = ScreenConfig
		return m, nil
	}
	m.logger = log
	m.logFilePath = log.GetPath()

	// 启动整理任务
	return m, m.organizeCmd()
}

// organizeCmd 整理命令 - 只扫描文件并启动处理
func (m Model) organizeCmd() tea.Cmd {
	return func() tea.Msg {
		// 扫描文件
		scanner := organizer.NewScanner(m.config.SourceDir)
		files, err := scanner.Scan()
		if err != nil {
			return OrganizeErrorMsg{Err: err}
		}

		return FileScanCompleteMsg{
			Files: files,
		}
	}
}

// processNextFileCmd 处理下一个文件
func (m Model) processNextFileCmd(fileIndex int) tea.Cmd {
	if fileIndex >= len(m.allFiles) || m.cancelled {
		// 所有文件处理完成
		m.statistics.EndTime = time.Now()
		m.statistics.Duration = m.statistics.EndTime.Sub(m.statistics.StartTime)

		return func() tea.Msg {
			return OrganizeCompleteMsg{
				Statistics: m.statistics,
				LogPath:    m.logFilePath,
			}
		}
	}

	// 捕获处理器引用用于闭包
	processor := m.processor
	logger := m.logger
	allFiles := m.allFiles

	return func() tea.Msg {
		file := allFiles[fileIndex]

		// 使用共享的处理器实例（方案4优化）
		record, _ := processor.Process(file)

		// 记录到日志
		if logger != nil {
			logger.LogRecord(record)
		}

		return FileProcessedMsg{
			Record:    record,
			FileIndex: fileIndex,
			Total:     len(allFiles),
		}
	}
}

// openDirectory 打开目录
func openDirectory(path string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("explorer", path)
		case "darwin":
			cmd = exec.Command("open", path)
		default: // linux
			cmd = exec.Command("xdg-open", path)
		}
		cmd.Run()
		return nil
	}
}

// 简化的输入处理
func (m Model) handleInputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		switch m.inputMode {
		case InputSource:
			m.config.SourceDir = m.inputValue
		case InputTarget:
			m.config.TargetDir = m.inputValue
		}
		m.currentScreen = ScreenConfig
		m.inputMode = InputNone
		return m, nil

	case "backspace":
		if len(m.inputValue) > 0 {
			m.inputValue = m.inputValue[:len(m.inputValue)-1]
		}
		return m, nil

	default:
		// 过滤特殊键
		if len(msg.String()) == 1 {
			m.inputValue += msg.String()
		}
		return m, nil
	}
}

func (m Model) handleConfigKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "s":
		m.currentScreen = ScreenInput
		m.inputMode = InputSource
		m.inputValue = m.config.SourceDir
		return m, nil
	case "d":
		m.currentScreen = ScreenInput
		m.inputMode = InputTarget
		m.inputValue = m.config.TargetDir
		return m, nil
	case "f":
		m.config.DuplicateDetection = config.DetectionFilename
		return m, nil
	case "m":
		m.config.DuplicateDetection = config.DetectionMD5
		return m, nil
	case "1":
		m.config.DuplicateStrategy = config.StrategySkip
		return m, nil
	case "2":
		m.config.DuplicateStrategy = config.StrategyOverwrite
		return m, nil
	case "3":
		m.config.DuplicateStrategy = config.StrategyRename
		return m, nil
	case "enter":
		if err := m.config.Validate(); err != nil {
			m.err = err
			return m, nil
		}
		return m.startOrganizing()
	case "q":
		if m.logger != nil {
			m.logger.Close()
		}
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleProgressKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if strings.ToLower(msg.String()) == "c" {
		m.cancelled = true
		m.isOrganizing = false
		m.currentScreen = ScreenConfig
		return m, nil
	}
	return m, nil
}

func (m Model) handleSummaryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "r":
		m.currentScreen = ScreenConfig
		m.isOrganizing = false
		m.statistics = &organizer.Statistics{}
		m.records = make([]organizer.ProcessRecord, 0)
		m.err = nil
		return m, nil
	case "o":
		return m, openDirectory(m.config.TargetDir)
	case "q":
		if m.logger != nil {
			m.logger.Close()
		}
		return m, tea.Quit
	}
	return m, nil
}
