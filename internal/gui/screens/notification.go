package screens

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Notification represents a system notification
type Notification struct {
	Title   string
	Message string
	Sound   bool
}

// sendNotification sends a system notification
func sendNotification(n *Notification) error {
	switch runtime.GOOS {
	case "darwin": // macOS
		return sendMacNotification(n)
	case "linux":
		return sendLinuxNotification(n)
	case "windows":
		return sendWindowsNotification(n)
	default:
		// Unsupported platform
		return fmt.Errorf("notifications not supported on %s", runtime.GOOS)
	}
}

// sendMacNotification sends a notification on macOS using osascript
func sendMacNotification(n *Notification) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, n.Message, n.Title)
	if n.Sound {
		script += ` sound name "default"`
	}
	
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// sendLinuxNotification sends a notification on Linux using notify-send
func sendLinuxNotification(n *Notification) error {
	args := []string{n.Title, n.Message}
	if !n.Sound {
		args = append(args, "--hint=string:sound-name:none")
	}
	
	cmd := exec.Command("notify-send", args...)
	return cmd.Run()
}

// sendWindowsNotification sends a notification on Windows using PowerShell
func sendWindowsNotification(n *Notification) error {
	script := fmt.Sprintf(`
		[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
		[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
		[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null
		
		$template = @"
		<toast>
			<visual>
				<binding template="ToastText02">
					<text id="1">%s</text>
					<text id="2">%s</text>
				</binding>
			</visual>
		</toast>
"@
		
		$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
		$xml.LoadXml($template)
		$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
		[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("Media Organizer").Show($toast)
	`, n.Title, n.Message)
	
	cmd := exec.Command("powershell", "-Command", script)
	return cmd.Run()
}
