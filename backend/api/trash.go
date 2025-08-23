package api

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var runtimeGOOS = runtime.GOOS

// moveToTrash moves the given file to the operating system's trash/recycle bin.
// It attempts to use platform-native mechanisms on Windows and macOS, and falls
// back to the freedesktop.org trash specification on other systems (e.g. Linux).
// If an error occurs, it is returned to the caller for handling.
func moveToTrash(path string) error {
	switch runtimeGOOS {
	case "windows":
		cmd := exec.Command("powershell", "-NoProfile", "-Command",
			fmt.Sprintf(`Add-Type -AssemblyName Microsoft.VisualBasic; [Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile(%q, [Microsoft.VisualBasic.FileIO.UIOption]::OnlyErrorDialogs, [Microsoft.VisualBasic.FileIO.RecycleOption]::SendToRecycleBin)`, path))
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("powershell recycle failed: %v: %s", err, strings.TrimSpace(string(out)))
		}
		return nil
	case "darwin":
		script := fmt.Sprintf(`tell application \"Finder\" to delete POSIX file %q`, path)
		cmd := exec.Command("osascript", "-e", script)
		return cmd.Run()
	default:
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		filesDir := filepath.Join(home, ".local/share/Trash/files")
		infoDir := filepath.Join(home, ".local/share/Trash/info")
		if err := os.MkdirAll(filesDir, 0o755); err != nil {
			return err
		}
		if err := os.MkdirAll(infoDir, 0o755); err != nil {
			return err
		}
		base := filepath.Base(abs)
		dest := filepath.Join(filesDir, base)
		for i := 1; ; i++ {
			if _, err := os.Stat(dest); os.IsNotExist(err) {
				break
			}
			dest = filepath.Join(filesDir, fmt.Sprintf("%s.%d", base, i))
		}
		if err := os.Rename(abs, dest); err != nil {
			return err
		}
		infoPath := filepath.Join(infoDir, filepath.Base(dest)+".trashinfo")
		u := url.PathEscape(abs)
		ts := time.Now().Format("2006-01-02T15:04:05")
		content := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n", u, ts)
		return os.WriteFile(infoPath, []byte(content), 0o644)
	}
}
