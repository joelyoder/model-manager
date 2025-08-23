//go:build trash_windows
// +build trash_windows

package api

import "testing"

func TestMoveToTrashWindows(t *testing.T) {
	old := runtimeGOOS
	runtimeGOOS = "windows"
	defer func() { runtimeGOOS = old }()
	testWindowsTrash(t)
}
