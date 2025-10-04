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
		{"MOV video", "/path/to/video.mov", FileTypeVideo},
		{"Other file", "/path/to/document.pdf", FileTypeOther},
		{"Uppercase extension", "/path/to/photo.JPG", FileTypePhoto},
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

func TestStatisticsGetSpeed(t *testing.T) {
	stats := &Statistics{
		ProcessedFiles: 100,
		Duration:       10 * 1000000000, // 10 seconds in nanoseconds
	}

	speed := stats.GetSpeed()
	expected := 10.0 // 100 files / 10 seconds

	if speed != expected {
		t.Errorf("GetSpeed() = %v, want %v", speed, expected)
	}
}

func TestStatisticsGetSpeedZeroDuration(t *testing.T) {
	stats := &Statistics{
		ProcessedFiles: 100,
		Duration:       0,
	}

	speed := stats.GetSpeed()
	expected := 0.0

	if speed != expected {
		t.Errorf("GetSpeed() with zero duration = %v, want %v", speed, expected)
	}
}
