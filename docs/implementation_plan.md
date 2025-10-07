# Implementation Plan

Add silent mode operation and configuration file/CLI parameter support to enable automated execution without TUI interface.

The current application is a Go-based media organizer with a TUI interface built on Bubble Tea. This implementation will add non-interactive operation modes while maintaining backward compatibility with the existing TUI interface. The solution will provide multiple configuration sources with clear precedence rules and a clean separation between interactive and automated modes.

## Types  
Extend configuration system with new types for silent mode and configuration sources.

### New Type Definitions
```go
// OperationMode defines how the application runs
type OperationMode string

const (
    ModeInteractive OperationMode = "interactive"  // Default TUI mode
    ModeSilent      OperationMode = "silent"       // Non-interactive mode
)

// ConfigSource defines where configuration comes from
type ConfigSource string

const (
    SourceDefault ConfigSource = "default"     // Default values
    SourceFile    ConfigSource = "file"        // Configuration file
    SourceCLI     ConfigSource = "cli"         // Command line arguments
)

// Extended Config structure
type Config struct {
    // Existing fields
    SourceDir          string             
    TargetDir          string             
    DuplicateDetection DuplicateDetection 
    DuplicateStrategy  DuplicateStrategy  
    
    // New fields
    Mode        OperationMode  // Operation mode (interactive/silent)
    ConfigFile  string         // Path to configuration file
    LogLevel    string         // Log level for silent mode
}
```

### Validation Rules
- `Mode`: Must be "interactive" or "silent"
- `ConfigFile`: Must be valid file path when specified
- `SourceDir` and `TargetDir`: Required in silent mode, optional in interactive mode
- All existing validation rules maintained

## Files
Create new files and modify existing files to support silent mode and configuration management.

### New Files to Create
- `cmd/organizer/cli.go` - CLI argument parsing and help system
- `internal/config/loader.go` - Configuration file loading and merging logic
- `internal/app/silent_runner.go` - Silent mode execution engine
- `config.example.json` - Example configuration file

### Existing Files to Modify
- `cmd/organizer/main.go` - Add mode detection and routing logic
- `internal/config/config.go` - Extend Config struct with new fields
- `internal/ui/model.go` - Add silent mode detection to prevent TUI initialization
- `internal/organizer/processor.go` - Add progress reporting for silent mode

### Configuration File Updates
- `go.mod` - Add dependencies for CLI parsing (cobra or flag package)
- Update build scripts to include new files

## Functions
Add new functions and modify existing ones to support configuration management and silent execution.

### New Functions
- `cmd/organizer/cli.go`:
  - `ParseCLI() (*Config, error)` - Parse command line arguments
  - `ShowHelp()` - Display usage information
  - `ShowVersion()` - Display version information

- `internal/config/loader.go`:
  - `LoadConfigFromFile(path string) (*Config, error)` - Load config from JSON file
  - `MergeConfigs(defaults, file, cli *Config) *Config` - Merge configs with precedence
  - `FindConfigFile() (string, error)` - Search for config file in standard locations

- `internal/app/silent_runner.go`:
  - `NewSilentRunner(config *Config) *SilentRunner` - Create silent mode runner
  - `Run() error` - Execute organization in silent mode
  - `printProgress(statistics *organizer.Statistics)` - Display progress updates
  - `printSummary(statistics *organizer.Statistics)` - Display final summary

### Modified Functions
- `cmd/organizer/main.go`:
  - `main()` - Add mode detection and routing to TUI or silent runner

- `internal/config/config.go`:
  - `NewDefaultConfig()` - Initialize with ModeInteractive
  - `Validate()` - Add mode-specific validation rules

- `internal/ui/model.go`:
  - `NewModel()` - Check if running in silent mode
  - `Init()` - Skip TUI initialization in silent mode

## Classes
Extend existing classes and create new ones for configuration management.

### New Classes
- `SilentRunner` (internal/app/silent_runner.go):
  - Fields: `config *Config`, `logger *logger.Logger`, `processor *organizer.Processor`
  - Methods: `Run()`, `printProgress()`, `printSummary()`, `handleInterrupt()`

- `CLIParser` (cmd/organizer/cli.go):
  - Fields: `flags *flag.FlagSet`, `config *Config`
  - Methods: `Parse()`, `validate()`, `showUsage()`

### Modified Classes
- `Config` (internal/config/config.go):
  - Add new fields: `Mode`, `ConfigFile`, `LogLevel`
  - Add validation for new fields

- `Model` (internal/ui/model.go):
  - Add silent mode detection logic
  - Skip TUI initialization when in silent mode

## Dependencies
Add minimal dependencies for CLI parsing and configuration management.

### New Dependencies
- Add `flag` package (standard library) for CLI argument parsing
- Consider `github.com/spf13/cobra` for advanced CLI features (optional)
- Use standard `encoding/json` for configuration file parsing

### Integration Requirements
- Maintain compatibility with existing Bubble Tea TUI
- No breaking changes to existing API
- Backward compatibility with current behavior

## Testing
Implement comprehensive testing for new functionality while maintaining existing test coverage.

### Test File Requirements
- `cmd/organizer/cli_test.go` - Test CLI argument parsing
- `internal/config/loader_test.go` - Test configuration loading and merging
- `internal/app/silent_runner_test.go` - Test silent mode execution
- Update existing tests to account for new configuration fields

### Test Scenarios
- CLI argument parsing with valid and invalid inputs
- Configuration file loading from different locations
- Config merging with different precedence scenarios
- Silent mode execution with various configurations
- Mode detection and routing logic
- Backward compatibility with existing TUI mode

### Validation Strategies
- Unit tests for individual components
- Integration tests for full workflow
- Error handling and edge case testing
- Performance testing for large file sets

## Implementation Order
Execute changes in logical sequence to minimize conflicts and ensure successful integration.

1. **Phase 1: Configuration System Extension**
   - Extend Config struct with new fields
   - Add validation for new configuration options
   - Create configuration loader with file support
   - Implement configuration merging logic

2. **Phase 2: CLI Argument Parsing**
   - Add CLI flag parsing for all configuration options
   - Implement help and version commands
   - Add mode detection logic (--silent flag)
   - Test CLI argument parsing and validation

3. **Phase 3: Silent Mode Execution Engine**
   - Create SilentRunner class for non-interactive execution
   - Implement progress reporting for silent mode
   - Add summary display functionality
   - Handle interrupt signals gracefully

4. **Phase 4: Main Application Routing**
   - Modify main() to detect operation mode
   - Route to TUI or silent runner based on configuration
   - Add configuration file discovery logic
   - Implement configuration precedence rules

5. **Phase 5: Integration and Testing**
   - Create example configuration file
   - Update documentation with new features
   - Test full workflow with various configurations
   - Verify backward compatibility

6. **Phase 6: Polish and Documentation**
   - Add comprehensive help text
   - Create usage examples
   - Update README with new features
   - Final testing and bug fixes

### Priority Rules
- CLI arguments override configuration file
- Configuration file overrides default values
- Silent mode disables TUI interface
- Required validation occurs before execution
