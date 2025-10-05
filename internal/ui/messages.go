package ui

import "github.com/chiyiangel/media-organizer-v2/internal/organizer"

// FileScanCompleteMsg 文件扫描完成消息
type FileScanCompleteMsg struct {
	Files []*organizer.FileInfo
}

// FileProcessedMsg 文件处理完成消息
type FileProcessedMsg struct {
	Record    *organizer.ProcessRecord
	FileIndex int
	Total     int
}

// OrganizeCompleteMsg 整理完成消息
type OrganizeCompleteMsg struct {
	Statistics *organizer.Statistics
	LogPath    string
}

// OrganizeErrorMsg 整理错误消息
type OrganizeErrorMsg struct {
	Err error
}

// ProgressUpdateMsg 进度更新消息
type ProgressUpdateMsg struct {
	Current  int
	Total    int
	FileName string
}
