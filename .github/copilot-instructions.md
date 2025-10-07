# AI Copilot Instructions for Photo Video Organizer

## Project Overview
This is a **Go 1.21+ Terminal UI application** that organizes media files (photos/videos) by date using EXIF metadata. The app uses **Bubble Tea** for interactive TUI with multi-screen navigation and real-time progress tracking, plus a **silent mode** for automated processing.

**Development Environment**: Primarily **Windows + PowerShell** - use PowerShell syntax for all command examples.

## Architecture & Code Organization

### Core Package Structure
```
cmd/organizer/          # Entry point with CLI parsing and Bubble Tea setup
internal/organizer/     # Business logic (processor, scanner, duplicate detection)  
internal/ui/           # Bubble Tea models and screens
internal/config/       # Configuration types and validation
internal/app/          # Silent mode runner for non-interactive execution
internal/i18n/         # Internationalization system (Chinese/English)
internal/logger/       # Structured logging system
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
```powershell
# Multi-platform builds via PowerShell script
.\build-all.ps1         # PowerShell cross-platform build script
# Or use Makefile in WSL/Linux environments if available
```

### Testing Patterns
- Unit tests in `*_test.go` focus on pure functions like `getFileType()` and `Statistics.GetSpeed()`
- Test structure: table-driven tests with `name`, `input`, `expected` fields

### Key Build Commands
```powershell
# Build and run commands
go build -o build/media-organizer.exe ./cmd/organizer  # Build for Windows
go test ./...                                         # Run all tests
.\build\media-organizer.exe                           # Execute built binary
.\build-all.ps1                                       # Cross-platform builds
```

## Module Dependencies
- **Core TUI**: `github.com/charmbracelet/bubbletea` + `lipgloss` for styling
- **EXIF**: `github.com/rwcarlsen/goexif` for photo metadata extraction
- **Module path issue**: Current go.mod uses placeholder `github.com/chiyiangel/media-organizer-v2` - should be updated to actual repo path



## Current Project Status

### Recent Developments
- ✅ **GitHub Actions CI/CD**: Complete release workflow with multi-platform builds
- ✅ **Internationalization (i18n)**: Full Chinese/English support with OS language auto-detection
- ✅ **Silent Mode**: Non-interactive execution mode for automated processing
- ✅ **Cross-platform Builds**: Windows/Linux/macOS support via Makefile and PowerShell scripts
- ✅ **Error Handling**: Proper format string handling for Go static analysis compliance

### Build System
- **PowerShell Script**: `build-all.ps1` for Windows-native multi-platform builds
- **Makefile**: Cross-platform build system with `MAIN_PATH := ./cmd/organizer`
- **Binary Name**: `media-organizer` (consistent across all build configurations)
- **GitHub Actions**: Automated testing and release creation on tag push

### Development Workflow
```powershell
# Development commands (Windows PowerShell)
go mod download; go mod tidy; go mod verify    # Install and verify dependencies
go fmt ./...                                   # Format code  
go test -v -race -coverprofile=coverage.out ./... # Run tests with race detection
go build -o build/media-organizer.exe ./cmd/organizer  # Build for current platform
.\build\media-organizer.exe                    # Execute built binary
.\build-all.ps1                                # Cross-platform builds
```

### Important Notes
- **PowerShell Priority**: Always use PowerShell syntax for command examples
- **Package Compilation**: Use `./cmd/organizer` not `cmd/organizer/main.go` to avoid `ParseCLI undefined` errors
- **i18n System**: All user-facing text supports Chinese/English via `internal/i18n` package
- **Format Strings**: Use `fmt.Errorf("%s", errorMsg)` pattern for i18n strings to satisfy Go static analyzer

## Key Features

### Internationalization (i18n)
- **Auto Language Detection**: Detects OS language (Chinese/English) via environment variables
- **Translation System**: `internal/i18n` package with `T()` and `Tf()` functions
- **UI Support**: Both Bubble Tea interface and silent mode fully localized
- **Usage Pattern**: `i18n.T("key")` for simple strings, `i18n.Tf("key", params...)` with placeholders

### Silent Mode (`internal/app/silent_runner.go`)
- **Non-interactive Processing**: Command-line execution without TUI
- **Progress Display**: Real-time progress updates with file statistics
- **Error Logging**: Structured logging to files with configurable levels
- **Configuration**: Uses same `Config` structure as interactive mode

### Build Configuration
- **PowerShell Script**: `build-all.ps1` for Windows-native multi-platform builds
- **Makefile**: Available for WSL/Linux environments  
- **GitHub Actions**: Automated CI/CD with cross-platform testing and releases
- **Binary Naming**: Consistent `media-organizer` prefix across all platforms
- **Version Embedding**: Git-based version injection during build process