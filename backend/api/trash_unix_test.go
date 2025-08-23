//go:build !trash_windows && !trash_darwin
// +build !trash_windows,!trash_darwin

package api

import "testing"

func TestMoveToTrashUnix(t *testing.T) {
	old := runtimeGOOS
	runtimeGOOS = "linux"
	defer func() { runtimeGOOS = old }()
	testUnixTrash(t)
}
