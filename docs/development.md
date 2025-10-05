# 照片视频整理App - 开发文档

## 📋 技术栈

### 核心技术
- **开发语言**: Go 1.21+
- **TUI框架**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) - 现代化的Go TUI框架
- **UI组件**: [Bubbles](https://github.com/charmbracelet/bubbles) - Bubble Tea的UI组件库
- **样式**: [Lipgloss](https://github.com/charmbracelet/lipgloss) - 终端样式库
- **EXIF处理**: [go-exif](https://github.com/rwcarlsen/goexif) - EXIF信息读取
- **文件操作**: Go标准库 `os`, `io`, `path/filepath`
- **哈希计算**: Go标准库 `crypto/md5`

---

## 🏗️ 项目结构

```
photo-video-organizer/
├── cmd/
│   └── organizer/
│       └── main.go                 # 程序入口
├── internal/
│   ├── app/
│   │   └── app.go                  # 应用主逻辑
│   ├── ui/
│   │   ├── model.go                # Bubble Tea主模型
│   │   ├── styles.go               # UI样式定义
│   │   ├── keys.go                 # 快捷键绑定
│   │   ├── messages.go             # 消息类型定义
│   │   └── screens/
│   │       ├── config.go           # 主配置界面
│   │       ├── input.go            # 路径输入弹窗
│   │       ├── progress.go         # 整理进度界面
│   │       └── summary.go          # 完成汇总界面
│   ├── organizer/
│   │   ├── scanner.go              # 文件扫描器
│   │   ├── metadata.go             # 元数据提取器
│   │   ├── processor.go            # 文件处理器
│   │   ├── duplicate.go            # 重复文件检测
│   │   └── types.go                # 类型定义
│   ├── config/
│   │   └── config.go               # 配置管理
│   └── logger/
│       └── logger.go               # 日志记录器
├── pkg/
│   └── fileutil/
│       └── fileutil.go             # 文件工具函数
├── docs/
│   ├── design.md                   # 功能设计文档
│   └── development.md              # 开发文档(本文件)
├── go.mod
├── go.sum
├── README.md
├── Makefile
└── .gitignore
```

### 目录说明

- **cmd/**: 应用程序入口点
- **internal/**: 私有应用代码，不可被外部导入
  - **app/**: 应用主逻辑和初始化
  - **ui/**: 用户界面相关代码（Bubble Tea模型和视图）
  - **organizer/**: 核心业务逻辑（扫描、处理、元数据提取）
  - **config/**: 配置管理
  - **logger/**: 日志记录
- **pkg/**: 可被外部导入的公共库代码

---

## 📦 核心依赖

### go.mod 依赖项

```go
module github.com/chiyiangel/media-organizer-v2

go 1.21

require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/bubbles v0.18.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd
)
```

### 安装命令

```bash
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/lipgloss
go get github.com/rwcarlsen/goexif
```

---

## 🔧 核心数据结构

### 1. 配置结构 (internal/config/config.go)

```go
package config

// DuplicateDetection 重复文件识别策略
type DuplicateDetection string

const (
    DetectionFilename DuplicateDetection = "filename" // 文件名
    DetectionMD5      DuplicateDetection = "md5"      // MD5哈希
)

// DuplicateStrategy 重复文件处理策略
type DuplicateStrategy string

const (
    StrategySkip     DuplicateStrategy = "skip"     // 跳过
    StrategyOverwrite DuplicateStrategy = "overwrite" // 覆盖
    StrategyRename    DuplicateStrategy = "rename"   // 重命名
)

// Config 应用配置
type Config struct {
    SourceDir          string             // 源目录
    TargetDir          string             // 目标目录
    DuplicateDetection DuplicateDetection // 重复识别策略
    DuplicateStrategy  DuplicateStrategy  // 重复处理策略
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() *Config {
    return &Config{
        SourceDir:          "",
        TargetDir:          "",
        DuplicateDetection: DetectionFilename,
        DuplicateStrategy:  StrategySkip,
    }
}

// Validate 验证配置
func (c *Config) Validate() error {
    if c.SourceDir == "" {
        return fmt.Errorf("源目录不能为空")
    }
    if c.TargetDir == "" {
        return fmt.Errorf("目标目录不能为空")
    }
    // 检查目录是否存在
    if _, err := os.Stat(c.SourceDir); os.IsNotExist(err) {
        return fmt.Errorf("源目录不存在: %s", c.SourceDir)
    }
    return nil
}
```

### 2. 文件信息结构 (internal/organizer/types.go)

```go
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
    Path         string    // 文件路径
    Name         string    // 文件名
    Type         FileType  // 文件类型
    Size         int64     // 文件大小
    Date         time.Time // 日期（来自EXIF或创建时间）
    MD5          string    // MD5哈希（按需计算）
    TargetPath   string    // 目标路径
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
    TotalFiles    int           // 总文件数
    ScannedFiles  int           // 已扫描文件数
    ProcessedFiles int          // 已处理文件数
    PhotoCount    int           // 照片数量
    VideoCount    int           // 视频数量
    SkippedCount  int           // 跳过数量
    FailedCount   int           // 失败数量
    StartTime     time.Time     // 开始时间
    EndTime       time.Time     // 结束时间
    Duration      time.Duration // 耗时
}

// GetSpeed 计算处理速度（文件/秒）
func (s *Statistics) GetSpeed() float64 {
    if s.Duration.Seconds() == 0 {
        return 0
    }
    return float64(s.ProcessedFiles) / s.Duration.Seconds()
}
```

### 3. Bubble Tea 模型 (internal/ui/model.go)

```go
package ui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/chiyiangel/media-organizer-v2/internal/config"
    "github.com/chiyiangel/media-organizer-v2/internal/organizer"
)

// Screen 界面类型
type Screen int

const (
    ScreenConfig   Screen = iota // 主配置界面
    ScreenInput                   // 路径输入弹窗
    ScreenProgress                // 整理进度界面
    ScreenSummary                 // 完成汇总界面
)

// InputMode 输入模式
type InputMode int

const (
    InputNone   InputMode = iota // 无输入
    InputSource                   // 输入源目录
    InputTarget                   // 输入目标目录
)

// Model Bubble Tea 主模型
type Model struct {
    // 当前界面
    currentScreen Screen
    
    // 配置
    config *config.Config
    
    // 输入状态
    inputMode     InputMode
    inputValue    string
    inputCursor   int
    
    // 整理状态
    isOrganizing  bool
    currentFile   *organizer.FileInfo
    statistics    *organizer.Statistics
    records       []organizer.ProcessRecord
    
    // 日志文件路径
    logFilePath   string
    
    // 窗口尺寸
    width         int
    height        int
    
    // 错误信息
    err           error
}

// NewModel 创建新模型
func NewModel() Model {
    return Model{
        currentScreen: ScreenConfig,
        config:        config.NewDefaultConfig(),
        inputMode:     InputNone,
        statistics:    &organizer.Statistics{},
        records:       make([]organizer.ProcessRecord, 0),
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
    case FileProcessedMsg:
        // 文件处理完成消息
        return m.handleFileProcessed(msg)
    case OrganizeCompleteMsg:
        // 整理完成消息
        return m.handleOrganizeComplete(msg)
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
```

---

## 🎨 UI实现要点

### 1. 样式定义 (internal/ui/styles.go)

```go
package ui

import (
    "github.com/charmbracelet/lipgloss"
)

var (
    // 边框样式
    borderStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("63")).
        Padding(1, 2)
    
    // 标题样式
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205")).
        Align(lipgloss.Center)
    
    // 标签样式
    labelStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("86")).
        Bold(true)
    
    // 普通文本样式
    textStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("252"))
    
    // 高亮样式
    highlightStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("226")).
        Bold(true)
    
    // 成功样式
    successStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("42"))
    
    // 错误样式
    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("196"))
    
    // 警告样式
    warningStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("214"))
    
    // 分割线样式
    dividerStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("238"))
    
    // 提示样式
    hintStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        Italic(true)
)

// 进度条渲染
func renderProgressBar(current, total int, width int) string {
    if total == 0 {
        return ""
    }
    
    percentage := float64(current) / float64(total)
    filled := int(float64(width) * percentage)
    
    bar := ""
    for i := 0; i < width; i++ {
        if i < filled {
            bar += "█"
        } else {
            bar += "░"
        }
    }
    
    return fmt.Sprintf("%s %.0f%% (%d/%d)", bar, percentage*100, current, total)
}
```

### 2. 快捷键处理 (internal/ui/keys.go)

```go
package ui

import (
    tea "github.com/charmbracelet/bubbletea"
)

// handleKeyPress 处理按键
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // 全局快捷键
    switch msg.String() {
    case "ctrl+c":
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

// handleConfigKeys 主配置界面按键
func (m Model) handleConfigKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "s", "S":
        // 编辑源目录
        m.currentScreen = ScreenInput
        m.inputMode = InputSource
        m.inputValue = m.config.SourceDir
        return m, nil
        
    case "d", "D":
        // 编辑目标目录
        m.currentScreen = ScreenInput
        m.inputMode = InputTarget
        m.inputValue = m.config.TargetDir
        return m, nil
        
    case "f", "F":
        // 选择文件名识别
        m.config.DuplicateDetection = config.DetectionFilename
        return m, nil
        
    case "m", "M":
        // 选择MD5识别
        m.config.DuplicateDetection = config.DetectionMD5
        return m, nil
        
    case "1":
        // 选择跳过策略
        m.config.DuplicateStrategy = config.StrategySkip
        return m, nil
        
    case "2":
        // 选择覆盖策略
        m.config.DuplicateStrategy = config.StrategyOverwrite
        return m, nil
        
    case "3":
        // 选择重命名策略
        m.config.DuplicateStrategy = config.StrategyRename
        return m, nil
        
    case "enter":
        // 开始整理
        if err := m.config.Validate(); err != nil {
            m.err = err
            return m, nil
        }
        return m.startOrganizing()
        
    case "q", "Q":
        return m, tea.Quit
    }
    
    return m, nil
}

// handleInputKeys 输入界面按键
func (m Model) handleInputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "enter":
        // 确认输入
        switch m.inputMode {
        case InputSource:
            m.config.SourceDir = m.inputValue
        case InputTarget:
            m.config.TargetDir = m.inputValue
        }
        m.currentScreen = ScreenConfig
        m.inputMode = InputNone
        return m, nil
        
    case "esc":
        // 取消输入
        m.currentScreen = ScreenConfig
        m.inputMode = InputNone
        return m, nil
        
    case "backspace":
        // 删除字符
        if len(m.inputValue) > 0 {
            m.inputValue = m.inputValue[:len(m.inputValue)-1]
        }
        return m, nil
        
    default:
        // 输入字符
        m.inputValue += msg.String()
        return m, nil
    }
}

// handleProgressKeys 进度界面按键
func (m Model) handleProgressKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "c", "C":
        // 取消整理（需要确认）
        return m.cancelOrganizing()
    }
    return m, nil
}

// handleSummaryKeys 汇总界面按键
func (m Model) handleSummaryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "r", "R":
        // 重新整理
        m.currentScreen = ScreenConfig
        m.isOrganizing = false
        m.statistics = &organizer.Statistics{}
        m.records = make([]organizer.ProcessRecord, 0)
        return m, nil
        
    case "o", "O":
        // 打开目标目录
        return m, openDirectory(m.config.TargetDir)
        
    case "q", "Q":
        return m, tea.Quit
    }
    return m, nil
}
```

---

## 🔨 核心功能实现

### 1. 文件扫描器 (internal/organizer/scanner.go)

```go
package organizer

import (
    "os"
    "path/filepath"
    "strings"
)

var (
    // 照片扩展名
    photoExtensions = []string{".jpg", ".jpeg", ".png", ".heic", ".gif", ".bmp", ".raw"}
    
    // 视频扩展名
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
```

### 2. 元数据提取器 (internal/organizer/metadata.go)

```go
package organizer

import (
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
```

### 3. 文件处理器 (internal/organizer/processor.go)

```go
package organizer

import (
    "crypto/md5"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
    
    "github.com/chiyiangel/media-organizer-v2/internal/config"
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
func (p *Processor) generateTargetPath(file *FileInfo) string {
    year := file.Date.Format("2006")
    month := file.Date.Format("01")
    day := file.Date.Format("02")
    
    targetDir := filepath.Join(p.config.TargetDir, year, month, day)
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
```

### 4. 重复文件检测 (internal/organizer/duplicate.go)

```go
package organizer

import (
    "os"
    "path/filepath"
    
    "github.com/chiyiangel/media-organizer-v2/internal/config"
)

// DuplicateDetector 重复文件检测器
type DuplicateDetector struct {
    config *config.Config
}

// NewDuplicateDetector 创建检测器
func NewDuplicateDetector(cfg *config.Config) *DuplicateDetector {
    return &DuplicateDetector{
        config: cfg,
    }
}

// IsDuplicate 检查是否重复
func (d *DuplicateDetector) IsDuplicate(file *FileInfo) (bool, error) {
    // 检查目标文件是否存在
    if _, err := os.Stat(file.TargetPath); os.IsNotExist(err) {
        return false, nil
    }
    
    switch d.config.DuplicateDetection {
    case config.DetectionFilename:
        // 文件名模式：文件存在即为重复
        return true, nil
        
    case config.DetectionMD5:
        // MD5模式：比较文件内容
        return d.compareByMD5(file)
        
    default:
        return false, nil
    }
}

// compareByMD5 通过MD5比较
func (d *DuplicateDetector) compareByMD5(file *FileInfo) (bool, error) {
    // 计算源文件MD5
    if file.MD5 == "" {
        md5, err := CalculateMD5(file.Path)
        if err != nil {
            return false, err
        }
        file.MD5 = md5
    }
    
    // 计算目标文件MD5
    targetMD5, err := CalculateMD5(file.TargetPath)
    if err != nil {
        return false, err
    }
    
    return file.MD5 == targetMD5, nil
}
```

---

## 🚀 构建与运行

### Makefile

```makefile
# 变量定义
BINARY_NAME=photo-organizer
BUILD_DIR=build
MAIN_PATH=cmd/organizer/main.go

# Go相关变量
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)

# 默认目标
.PHONY: all
all: clean build

# 构建
.PHONY: build
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(GOBIN)/$(BINARY_NAME)"

# 运行
.PHONY: run
run: build
	@$(GOBIN)/$(BINARY_NAME)

# 清理
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean complete"

# 测试
.PHONY: test
test:
	@echo "Testing..."
	@go test -v ./...

# 安装依赖
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# 代码格式化
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 代码检查
.PHONY: lint
lint:
	@echo "Linting code..."
	@golangci-lint run

# 跨平台构建
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(GOBIN)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build -o $(GOBIN)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build -o $(GOBIN)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=linux GOARCH=amd64 go build -o $(GOBIN)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "Build complete for all platforms"
```

### 构建命令

```bash
# 安装依赖
make deps

# 构建
make build

# 运行
make run

# 测试
make test

# 代码格式化
make fmt

# 跨平台构建
make build-all
```

---

## 📝 代码规范

### Effective Go 规范要点

1. **命名规范**
   - 使用驼峰命名法（camelCase/PascalCase）
   - 导出的名称首字母大写
   - 包名使用小写单个单词
   - 接口名使用 `-er` 后缀（如 `Reader`, `Writer`）

2. **注释规范**
   - 每个导出的函数、类型、常量都需要注释
   - 注释以名称开头，如 `// NewModel creates a new model`
   - 包注释在 package 语句之前

3. **错误处理**
   - 优先返回错误而非panic
   - 使用 `fmt.Errorf` 包装错误信息
   - 错误信息小写开头，不以标点结尾

4. **并发安全**
   - 使用 channel 进行goroutine通信
   - 避免共享内存，使用消息传递
   - 必要时使用 mutex 保护共享状态

5. **代码组织**
   - 相关功能放在同一个包中
   - 避免循环依赖
   - 使用 internal 包防止外部导入

---

## 🧪 测试策略

### 单元测试示例

```go
// internal/organizer/scanner_test.go
package organizer

import (
    "testing"
)

func TestGetFileType(t *testing.T) {
    tests := []struct {
        name     string
        path     string
        expected FileType
    }{
        {"JPEG photo", "/path/to/photo.jpg", FileTypePhoto},
        {"PNG photo", "/path/to/photo.png", FileTypePhoto},
        {"MP4 video", "/path/to/video.mp4", FileTypeVideo},
        {"Other file", "/path/to/document.pdf", FileTypeOther},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := getFileType(tt.path)
            if result != tt.expected {
                t.Errorf("getFileType(%s) = %v, want %v", 
                    tt.path, result, tt.expected)
            }
        })
    }
}
```

---

## 📚 参考资源

- [Effective Go](https://go.dev/doc/effective_go)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Examples](https://github.com/charmbracelet/lipgloss/tree/master/examples)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
