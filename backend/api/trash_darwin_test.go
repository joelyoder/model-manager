//go:build trash_darwin
// +build trash_darwin

package api

import "testing"

func TestMoveToTrashDarwin(t *testing.T) {
	old := runtimeGOOS
	runtimeGOOS = "darwin"
	defer func() { runtimeGOOS = old }()
	testDarwinTrash(t)
}
