package screens

import (
	"os/exec"
	"runtime"
)

// openFolder opens a folder in the system file manager
func openFolder(path string) error {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("explorer", path)
	case "linux":
		// Try common Linux file managers
		cmd = exec.Command("xdg-open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	
	return cmd.Start()
}
