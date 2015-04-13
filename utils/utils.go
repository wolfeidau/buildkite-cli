package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/github/hub/utils"
	"github.com/mattn/go-isatty"
)

// Check simple error assertion which exits the app
func Check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stdout, format, a...)
}

func Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stdout, a...)
}

func ToMap(columns []string, values []interface{}) map[string]interface{} {

	m := make(map[string]interface{})

	if len(columns) != len(values) {
		utils.Check(fmt.Errorf("Miss match in columns and values"))
	}

	for i, c := range columns {
		m[c] = values[i]
	}

	return m
}

func IsTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd)
}

func BrowserLauncher() ([]string, error) {
	browser := os.Getenv("BROWSER")
	if browser == "" {
		browser = searchBrowserLauncher(runtime.GOOS)
	}

	if browser == "" {
		return nil, errors.New("Please set $BROWSER to a web launcher")
	}

	return strings.Split(browser, " "), nil
}

func searchBrowserLauncher(goos string) (browser string) {
	switch goos {
	case "darwin":
		browser = "open"
	case "windows":
		browser = "cmd /c start"
	default:
		candidates := []string{"xdg-open", "cygstart", "x-www-browser", "firefox",
			"opera", "mozilla", "netscape"}
		for _, b := range candidates {
			path, err := exec.LookPath(b)
			if err == nil {
				browser = path
				break
			}
		}
	}

	return browser
}
