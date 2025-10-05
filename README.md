# 📸 Media Organizer

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Release](https://img.shields.io/github/v/release/chiyiangel/Media-Organizer-V2?style=for-the-badge)](https://github.com/chiyiangel/Media-Organizer-V2/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/chiyiangel/Media-Organizer-V2/build.yml?style=for-the-badge)](https://github.com/chiyiangel/Media-Organizer-V2/actions)

**A powerful command-line tool for organizing photos and videos with an elegant Terminal UI**

[Features](#-features) •
[Installation](#-installation) •
[Usage](#-usage) •
[Documentation](#-documentation) •
[Contributing](#-contributing)

</div>

---

## 🌟 Features

### 🎯 Smart Organization
- **� Date-based Structure**: Automatically organizes files by `YYYY/MM/MM-DD` format
- **📖 EXIF Support**: Extracts shooting date from photo metadata
- **🎥 Video Support**: Handles video files using creation/modification timestamps
- **🔍 Duplicate Detection**: Filename or MD5-based duplicate checking

### 🎨 Modern Interface
- **✨ Beautiful TUI**: Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework
- **📊 Real-time Progress**: Live updates during file processing
- **⌨️ Keyboard Navigation**: Intuitive keyboard shortcuts
- **🎨 Styled Components**: Professional styling with [Lipgloss](https://github.com/charmbracelet/lipgloss)

### ⚙️ Flexible Configuration
- **🔄 Duplicate Strategies**: Skip, overwrite, or rename duplicates
- ** Comprehensive Logging**: Detailed operation logs with timestamps
- **️ Customizable**: Configurable source and destination directories
- **🤫 Silent Mode**: Non-interactive execution for automation
- **📁 Configuration Files**: JSON-based configuration management
- **🔧 CLI Parameters**: Command-line argument support

## 🚀 Installation

### Prerequisites
- Go 1.21 or higher

### Option 1: Download Binary (Recommended)
```bash
# Download the latest release for your platform
curl -L https://github.com/chiyiangel/Media-Organizer-V2/releases/latest/download/media-organizer-linux-amd64 -o media-organizer
chmod +x media-organizer
```

### Option 2: Build from Source
```bash
# Clone the repository
git clone https://github.com/chiyiangel/Media-Organizer-V2.git
cd Media-Organizer-V2

# Download dependencies
go mod download

# Build the application
make build
```

### Option 3: Go Install
```bash
go install github.com/chiyiangel/Media-Organizer-V2/cmd/organizer@latest
```

## 📖 Usage

### Quick Start
```bash
# Launch the application
./media-organizer

# On Windows
media-organizer.exe
```

### Interactive Interface

The application provides an intuitive keyboard-driven interface:

#### Main Configuration Screen
- `S` - Set source directory
- `D` - Set destination directory  
- `F` - Select filename-based duplicate detection
- `M` - Select MD5-based duplicate detection
- `1/2/3` - Choose duplicate handling strategy (Skip/Overwrite/Rename)
- `Enter` - Start organization process
- `Q` - Quit application

#### Progress Screen
- `C` / `Esc` - Cancel operation
- Real-time progress bar and statistics

#### Summary Screen
- `R` - Restart organization
- `O` - Open destination folder
- `Q` - Quit application

### Silent Mode & CLI Usage

For automated execution and scripting, the application supports silent mode operation:

#### Command Line Options
```bash
# Core options
-source string      Source directory path
-target string      Target directory path
-detection string   Duplicate detection strategy (filename, md5)
-strategy string    Duplicate handling strategy (skip, overwrite, rename)

# Silent mode options
-mode string        Operation mode (interactive, silent)
-silent             Enable silent mode (equivalent to --mode silent)
-config string      Configuration file path
-log-level string   Log level (debug, info, warning, error)

# Information options
-help               Show this help message
-version            Show version information
```

#### Usage Examples
```bash
# Silent mode with CLI parameters
./media-organizer -source ./photos -target ./organized -mode silent

# Silent mode with configuration file
./media-organizer -config config.json -silent

# Using short flags
./media-organizer -source ./photos -target ./organized -silent

# Show help
./media-organizer -help

# Show version
./media-organizer -version
```

#### Configuration File Support

Create a JSON configuration file for repeated use:

**config.json**
```json
{
  "sourceDir": "./photos",
  "targetDir": "./organized",
  "duplicateDetection": "md5",
  "duplicateStrategy": "rename",
  "mode": "silent",
  "logLevel": "info"
}
```

Configuration precedence: **CLI arguments > Configuration file > Default values**

#### Configuration File Locations
The application searches for configuration files in these locations (in order):
1. Current directory: `media-organizer.json`
2. User config directory: `%APPDATA%\media-organizer\config.json` (Windows)
3. Home directory: `~/.media-organizer.json`

### Supported File Types

| Type | Extensions |
|------|------------|
| **Photos** | `.arw`, `.jpg`, `.jpeg`, `.png`, `.heic`, `.gif`, `.bmp`, `.raw` |
| **Videos** | `.mp4`, `.mov`, `.avi`, `.mkv`, `.flv`, `.wmv` |

### Organization Structure

Files are organized using the following structure:
```
destination/
├── 2024/
│   ├── 01/
│   │   ├── 01-15/    # January 15th
│   │   └── 01-28/    # January 28th
│   └── 03/
│       └── 03-22/    # March 22nd
└── 2025/
    └── 10/
        └── 10-04/    # October 4th
```

### Duplicate Handling

#### Detection Methods
- **Filename**: Compares file names only
- **MD5**: Compares file content using MD5 hash

#### Handling Strategies  
- **Skip**: Keep existing file, ignore duplicate
- **Overwrite**: Replace existing file with new one
- **Rename**: Add suffix like `image(1).jpg`, `image(2).jpg`

### Date Extraction Priority
1. **Photos**: EXIF DateTimeOriginal → File modification time
2. **Videos**: File creation time → File modification time

## 📊 Example Output

```
🔄 Processing files...

Current file: IMG_2024.jpg
Target path: /photos/2024/03/03-15/IMG_2024.jpg

Progress: [████████████████████░░░░] 8/10 (80%)

────────────────────────────────────

📊 Real-time Statistics:
    Scanned:     10 files
    Processed:    8 files
    ├─ Photos:    6 files
    ├─ Videos:    2 files
    ├─ Skipped:   1 file (duplicate)
    └─ Failed:    0 files

────────────────────────────────────

Press [C/Esc] to cancel
```

## 🏗️ Architecture

### Project Structure
```
photo-video-organizer/
├── cmd/
│   └── organizer/          # Application entry point
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── logger/            # Logging functionality
│   ├── organizer/         # Core business logic
│   │   ├── duplicate.go   # Duplicate detection
│   │   ├── metadata.go    # EXIF/metadata extraction
│   │   ├── processor.go   # File processing
│   │   ├── scanner.go     # Directory scanning
│   │   └── types.go       # Data structures
│   └── ui/                # Terminal user interface
│       ├── messages.go    # UI messages
│       ├── model.go       # Bubble Tea model
│       ├── screens.go     # Screen rendering
│       └── styles.go      # UI styling
├── docs/                  # Documentation
├── build/                 # Build artifacts
├── Makefile              # Build automation
├── go.mod                # Go modules
└── README.md
```

### Tech Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Language** | [Go](https://golang.org) 1.21+ | Core application |
| **TUI Framework** | [Bubble Tea](https://github.com/charmbracelet/bubbletea) | Interactive terminal UI |
| **Styling** | [Lipgloss](https://github.com/charmbracelet/lipgloss) | Terminal styling |
| **UI Components** | [Bubbles](https://github.com/charmbracelet/bubbles) | Reusable UI components |
| **EXIF Processing** | [go-exif](https://github.com/rwcarlsen/goexif) | Metadata extraction |

## 🔧 Development

### Prerequisites
- Go 1.21+
- Make (optional, for build automation)

### Setup
```bash
# Clone the repository
git clone https://github.com/chiyiangel/Media-Organizer-V2.git
cd photo-video-organizer

# Install dependencies
go mod download

# Run tests
go test ./...

# Build and run
make build && make run
```

### Available Make Commands
```bash
make build      # Build the application
make run        # Build and run
make clean      # Clean build artifacts
make fmt        # Format code
make deps       # Download dependencies
make build-all  # Cross-platform builds
```

### Quick Contributing Steps
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Reporting Issues
- Use the [issue tracker](https://github.com/chiyiangel/Media-Organizer-V2/issues)
- Provide detailed description and steps to reproduce
- Include system information (OS, Go version, etc.)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [go-exif](https://github.com/rwcarlsen/goexif) - EXIF metadata processing

## 📧 Contact

- **Author**: Liu Jun
- **Email**: chiyiangel@gmail.com
- **GitHub**: [@chiyiangel](https://github.com/chiyiangel)
- **Issues**: [GitHub Issues](https://github.com/chiyiangel/Media-Organizer-V2/issues)

---

<div align="center">

⭐ **Star this repository if you find it useful!** ⭐

</div>
