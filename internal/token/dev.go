//go:build !production

package token

import "os"

// DevToken is the DESKTOP_DEV_TOKEN override. wails build sets the
// production tag and compiles this function out.
func DevToken() string {
	return os.Getenv("DESKTOP_DEV_TOKEN")
}
