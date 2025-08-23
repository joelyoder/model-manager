package api

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testUnixTrash(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	file := filepath.Join(tmpHome, "example.txt")
	if err := os.WriteFile(file, []byte("data"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	if err := moveToTrash(file); err != nil {
		t.Fatalf("moveToTrash: %v", err)
	}
	dest := filepath.Join(tmpHome, ".local/share/Trash/files/example.txt")
	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("expected file moved: %v", err)
	}
	info := filepath.Join(tmpHome, ".local/share/Trash/info/example.txt.trashinfo")
	b, err := os.ReadFile(info)
	if err != nil {
		t.Fatalf("read info: %v", err)
	}
	if !strings.Contains(string(b), "Path="+url.PathEscape(file)) {
		t.Fatalf("info missing path: %s", b)
	}
	if !strings.Contains(string(b), "DeletionDate=") {
		t.Fatalf("info missing deletion date: %s", b)
	}
}

func testWindowsTrash(t *testing.T) {
	tmp := t.TempDir()
	logPath := filepath.Join(tmp, "log")
	psPath := filepath.Join(tmp, "powershell")
	script := fmt.Sprintf("#!/bin/sh\necho $@ > %s\n", logPath)
	if err := os.WriteFile(psPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write stub: %v", err)
	}
	t.Setenv("PATH", tmp+":"+os.Getenv("PATH"))
	file := filepath.Join(tmp, "example.txt")
	if err := os.WriteFile(file, []byte("data"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := moveToTrash(file); err != nil {
		t.Fatalf("moveToTrash: %v", err)
	}
	b, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !strings.Contains(string(b), file) {
		t.Fatalf("powershell not called with file: %s", b)
	}
}

func testDarwinTrash(t *testing.T) {
	tmp := t.TempDir()
	logPath := filepath.Join(tmp, "log")
	osaPath := filepath.Join(tmp, "osascript")
	script := fmt.Sprintf("#!/bin/sh\necho $@ > %s\n", logPath)
	if err := os.WriteFile(osaPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write stub: %v", err)
	}
	t.Setenv("PATH", tmp+":"+os.Getenv("PATH"))
	file := filepath.Join(tmp, "example.txt")
	if err := os.WriteFile(file, []byte("data"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := moveToTrash(file); err != nil {
		t.Fatalf("moveToTrash: %v", err)
	}
	b, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !strings.Contains(string(b), file) {
		t.Fatalf("osascript not called with file: %s", b)
	}
}
