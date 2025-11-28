package organizer

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chiyiangel/media-organizer-v2/internal/config"
)

func TestProcessorCreatedDirsCache(t *testing.T) {
	// Test that the processor properly caches created directories
	// Use os.TempDir() for cross-platform compatibility
	tempDir := os.TempDir()
	cfg := &config.Config{
		SourceDir:          filepath.Join(tempDir, "source"),
		TargetDir:          filepath.Join(tempDir, "target"),
		DuplicateDetection: config.DetectionFilename,
		DuplicateStrategy:  config.StrategySkip,
	}

	processor := NewProcessor(cfg)

	// Verify createdDirs is initialized
	if processor.createdDirs == nil {
		t.Error("createdDirs should be initialized")
	}

	// Verify it's empty initially
	if len(processor.createdDirs) != 0 {
		t.Errorf("createdDirs should be empty initially, got %d entries", len(processor.createdDirs))
	}
}

func TestPreCreateDirectories(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "processor_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source directory with test files
	sourceDir := filepath.Join(tmpDir, "source")
	targetDir := filepath.Join(tmpDir, "target")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}

	cfg := &config.Config{
		SourceDir:          sourceDir,
		TargetDir:          targetDir,
		DuplicateDetection: config.DetectionFilename,
		DuplicateStrategy:  config.StrategySkip,
	}

	processor := NewProcessor(cfg)

	// Create test files with known dates
	testDate := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	files := []*FileInfo{
		{
			Path: filepath.Join(sourceDir, "photo1.jpg"),
			Name: "photo1.jpg",
			Type: FileTypePhoto,
			Date: testDate,
		},
		{
			Path: filepath.Join(sourceDir, "photo2.jpg"),
			Name: "photo2.jpg",
			Type: FileTypePhoto,
			Date: testDate, // Same date, same directory
		},
		{
			Path: filepath.Join(sourceDir, "video1.mp4"),
			Name: "video1.mp4",
			Type: FileTypeVideo,
			Date: time.Date(2024, 5, 20, 14, 0, 0, 0, time.UTC), // Different date
		},
	}

	// Create actual test files
	for _, f := range files {
		if err := os.WriteFile(f.Path, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", f.Path, err)
		}
		// Set the modification time
		if err := os.Chtimes(f.Path, f.Date, f.Date); err != nil {
			t.Fatalf("Failed to set file time for %s: %v", f.Path, err)
		}
	}

	// Call PreCreateDirectories
	err = processor.PreCreateDirectories(files)
	if err != nil {
		t.Fatalf("PreCreateDirectories failed: %v", err)
	}

	// Verify that directories are cached
	if len(processor.createdDirs) == 0 {
		t.Error("createdDirs should have entries after PreCreateDirectories")
	}

	// Verify target directories were created
	expectedDirs := []string{
		filepath.Join(targetDir, "2024", "03", "03-15"),
		filepath.Join(targetDir, "2024", "05", "05-20"),
	}

	for _, dir := range expectedDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Expected directory %s to be created", dir)
		}
		if !processor.createdDirs[dir] {
			t.Errorf("Expected directory %s to be in createdDirs cache", dir)
		}
	}
}

func TestCopyBufferSize(t *testing.T) {
	// Verify the buffer size constant is set to 64KB
	expectedSize := 64 * 1024
	if copyBufferSize != expectedSize {
		t.Errorf("copyBufferSize = %d, want %d", copyBufferSize, expectedSize)
	}
}
