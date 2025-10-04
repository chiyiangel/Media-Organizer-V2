# AI Copilot Instructions for Photo Video Organizer

## Project Overview
This is a **Go 1.21+ Terminal UI application** that organizes photos/videos by date using EXIF metadata. The app uses **Bubble Tea** for interactive TUI with multi-screen navigation and real-time progress tracking.

## Architecture & Code Organization

### Core Package Structure
```
cmd/organizer/          # Entry point with Bubble Tea setup
internal/organizer/     # Business logic (processor, scanner, duplicate detection)  
internal/ui/           # Bubble Tea models and screens
internal/config/       # Configuration types and validation
```

### Key Data Flow
1. **Scanner** (`scanner.go`) → discovers files and determines `FileType` (photo/video/other)
2. **MetadataExtractor** (`metadata.go`) → extracts dates from EXIF or file timestamps
3. **DuplicateDetector** (`duplicate.go`) → handles filename/MD5-based duplicate checking
4. **Processor** (`processor.go`) → orchestrates the pipeline: date extraction → path generation → duplicate handling → file copying

### Critical Types in `internal/organizer/types.go`
- `FileInfo`: Contains `Path`, `Type`, `Date`, `MD5`, `TargetPath` 
- `Statistics`: Real-time counters with `GetSpeed()` method for performance display
- `ProcessRecord`: Result tracking (`success`/`skipped`/`failed`) with messages

## Development Patterns

### Bubble Tea UI Architecture
- **Screen-based navigation**: `ScreenConfig` → `ScreenProgress` → `ScreenSummary` via `currentScreen` field
- **Message-driven updates**: Custom messages like `FileScanCompleteMsg` drive state transitions
- **Input handling**: Modal input dialogs for paths using `InputMode` enum (`InputSource`/`InputTarget`)

### File Organization Logic
- **Date-based paths**: `YYYY/MM/MM-DD` structure (e.g., `2024/03/03-15/`)
- **Priority order**: EXIF DateTimeOriginal → file modification time → file creation time
- **Duplicate strategies**: Skip/Overwrite/Rename with unique suffix generation

### Error Handling Convention
Functions return `(*ProcessRecord, error)` - the record contains user-friendly messages while error is for technical details. Always populate both for failed operations.

## Build & Development

### Cross-Platform Building
```bash
# Multi-platform builds via PowerShell script or Makefile
make build-all          # Uses Makefile with GOOS/GOARCH detection
.\build-all.ps1         # PowerShell equivalent for Windows
```

### Testing Patterns
- Unit tests in `*_test.go` focus on pure functions like `getFileType()` and `Statistics.GetSpeed()`
- Test structure: table-driven tests with `name`, `input`, `expected` fields

### Key Build Commands
```bash
make build      # Single platform build with version from git tags
make run        # Build and execute  
go test ./...   # Run all tests
```

## Module Dependencies
- **Core TUI**: `github.com/charmbracelet/bubbletea` + `lipgloss` for styling
- **EXIF**: `github.com/rwcarlsen/goexif` for photo metadata extraction
- **Module path issue**: Current go.mod uses placeholder `github.com/yourusername/photo-video-organizer` - should be updated to actual repo path

## Configuration & State Management

### Config Structure (`internal/config/config.go`)
```go
type Config struct {
    SourceDir          string
    TargetDir          string  
    DuplicateDetection DuplicateDetection  // "filename" | "md5"
    DuplicateStrategy  DuplicateStrategy   // "skip" | "overwrite" | "rename"
}
```

### UI State (`internal/ui/model.go`)
The main model tracks: current screen, config, input state, organizing progress, statistics, and file records. State transitions are message-driven following Bubble Tea patterns.

## File Processing Pipeline

1. **Scanning**: Recursively find files, classify by extension into `FileType`
2. **Metadata**: Extract date from EXIF (photos) or file timestamps (videos)  
3. **Path Generation**: Create `YYYY/MM/MM-DD` structure based on extracted date
4. **Duplicate Check**: Compare by filename or MD5 hash based on config
5. **Processing**: Copy/move files with progress tracking and error recording

When adding features, maintain this pipeline structure and ensure proper error propagation through `ProcessRecord` objects.