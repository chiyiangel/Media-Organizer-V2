package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chiyiangel/media-organizer-v2/internal/config"
	"github.com/chiyiangel/media-organizer-v2/internal/i18n"
)

// CLIParser handles command line argument parsing
type CLIParser struct {
	flags       *flag.FlagSet
	config      *config.Config
	silentFlag  bool
	showHelp    bool
	showVersion bool
}

// NewCLIParser creates a new CLI parser
func NewCLIParser() *CLIParser {
	parser := &CLIParser{
		config: &config.Config{},
	}
	parser.setupFlags()
	return parser
}

// setupFlags configures all CLI flags
func (p *CLIParser) setupFlags() {
	p.flags = flag.NewFlagSet("organizer", flag.ExitOnError)

	// Core configuration flags
	p.flags.StringVar(&p.config.SourceDir, "source", "", i18n.T("cli.option.source"))
	p.flags.StringVar(&p.config.TargetDir, "target", "", i18n.T("cli.option.target"))
	p.flags.StringVar((*string)(&p.config.DuplicateDetection), "detection", "", i18n.T("cli.option.detection"))
	p.flags.StringVar((*string)(&p.config.DuplicateStrategy), "strategy", "", i18n.T("cli.option.strategy"))

	// Silent mode and configuration flags
	p.flags.StringVar((*string)(&p.config.Mode), "mode", "", i18n.T("cli.option.mode"))
	p.flags.BoolVar(&p.silentFlag, "silent", false, i18n.T("cli.option.silent"))
	p.flags.StringVar(&p.config.ConfigFile, "config", "", i18n.T("cli.option.config"))
	p.flags.StringVar(&p.config.LogLevel, "log-level", "", i18n.T("cli.option.log_level"))

	// Information flags
	p.flags.BoolVar(&p.showHelp, "help", false, i18n.T("cli.option.help"))
	p.flags.BoolVar(&p.showVersion, "version", false, i18n.T("cli.option.version"))
}

// Parse parses command line arguments
func (p *CLIParser) Parse(args []string) (*config.Config, error) {
	if err := p.flags.Parse(args); err != nil {
		return nil, err
	}

	// Handle information flags first
	if p.showHelp {
		p.ShowHelp()
		os.Exit(0)
	}

	if p.showVersion {
		p.ShowVersion()
		os.Exit(0)
	}

	// Handle silent flag
	if p.silentFlag {
		p.config.Mode = config.ModeSilent
	}

	// Validate the parsed configuration
	if err := p.validate(); err != nil {
		return nil, err
	}

	return p.config, nil
}

// validate validates the parsed CLI configuration
func (p *CLIParser) validate() error {
	// Validate operation mode if specified
	if p.config.Mode != "" && p.config.Mode != config.ModeInteractive && p.config.Mode != config.ModeSilent {
		return fmt.Errorf(i18n.Tf("cli.error.invalid_mode", p.config.Mode))
	}

	// Validate duplicate detection if specified
	if p.config.DuplicateDetection != "" &&
		p.config.DuplicateDetection != config.DetectionFilename &&
		p.config.DuplicateDetection != config.DetectionMD5 {
		return fmt.Errorf(i18n.Tf("cli.error.invalid_detection", p.config.DuplicateDetection))
	}

	// Validate duplicate strategy if specified
	if p.config.DuplicateStrategy != "" &&
		p.config.DuplicateStrategy != config.StrategySkip &&
		p.config.DuplicateStrategy != config.StrategyOverwrite &&
		p.config.DuplicateStrategy != config.StrategyRename {
		return fmt.Errorf(i18n.Tf("cli.error.invalid_strategy", p.config.DuplicateStrategy))
	}

	// Validate log level if specified
	if p.config.LogLevel != "" {
		validLogLevels := map[string]bool{
			"debug":   true,
			"info":    true,
			"warning": true,
			"error":   true,
		}
		if !validLogLevels[p.config.LogLevel] {
			return fmt.Errorf(i18n.Tf("cli.error.invalid_log_level", p.config.LogLevel))
		}
	}

	return nil
}

// ShowHelp displays usage information
func (p *CLIParser) ShowHelp() {
	fmt.Printf(i18n.Tf("cli.help.title", getVersion()) + "\n\n")
	fmt.Println(i18n.T("cli.help.usage"))
	fmt.Println()
	fmt.Println(i18n.T("cli.options.core"))
	fmt.Println("  -source string      " + i18n.T("cli.option.source"))
	fmt.Println("  -target string      " + i18n.T("cli.option.target"))
	fmt.Println("  -detection string   " + i18n.T("cli.option.detection"))
	fmt.Println("  -strategy string    " + i18n.T("cli.option.strategy"))
	fmt.Println()
	fmt.Println(i18n.T("cli.options.silent"))
	fmt.Println("  -mode string        " + i18n.T("cli.option.mode"))
	fmt.Println("  -silent             " + i18n.T("cli.option.silent"))
	fmt.Println("  -config string      " + i18n.T("cli.option.config"))
	fmt.Println("  -log-level string   " + i18n.T("cli.option.log_level"))
	fmt.Println()
	fmt.Println(i18n.T("cli.options.info"))
	fmt.Println("  -help               " + i18n.T("cli.option.help"))
	fmt.Println("  -version            " + i18n.T("cli.option.version"))
	fmt.Println()
	fmt.Println(i18n.T("cli.examples"))
	fmt.Println("  " + i18n.T("cli.example.silent_mode"))
	fmt.Println("  " + i18n.T("cli.example.config_file"))
	fmt.Println("  " + i18n.T("cli.example.help"))
}

// ShowVersion displays version information
func (p *CLIParser) ShowVersion() {
	fmt.Printf(i18n.Tf("cli.version", getVersion()) + "\n")
}

// getVersion returns the application version
func getVersion() string {
	// This should be set during build, using a placeholder for now
	return "2.0.0"
}

// ParseCLI is the main entry point for CLI parsing
func ParseCLI() (*config.Config, error) {
	parser := NewCLIParser()
	return parser.Parse(os.Args[1:])
}
