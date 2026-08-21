package util

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
)

// OpenURL opens a trusted HTTPS URL in the user's default browser.
func OpenURL(rawURL string) (err error) {
	var (
		parsedURL *url.URL
		command   string
		arguments []string
	)

	if parsedURL, err = url.Parse(rawURL); err != nil {
		err = fmt.Errorf("parse browser URL: %w", err)
		return
	}

	if parsedURL.Scheme != "https" {
		err = fmt.Errorf("refusing to open non-HTTPS URL")
		return
	}

	switch runtime.GOOS {
	case "windows":
		command = "rundll32.exe"
		arguments = []string{"url.dll,FileProtocolHandler", rawURL}
	case "darwin":
		command = "open"
		arguments = []string{rawURL}
	case "linux":
		command = "xdg-open"
		arguments = []string{rawURL}
	default:
		err = fmt.Errorf("opening a browser is unsupported on %s", runtime.GOOS)
		return
	}

	// #nosec G204 -- command is selected from the fixed platform allowlist above.
	if err = exec.Command(command, arguments...).Start(); err != nil {
		err = fmt.Errorf("open default browser: %w", err)
	}

	return
}
