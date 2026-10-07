//go:build production

package token

// DevToken is empty in a production build. DESKTOP_DEV_TOKEN is ignored.
func DevToken() string {
	return ""
}
