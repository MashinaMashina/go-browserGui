package programs

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		panic(err)
	}

}

func OpenBrowserDelay(url string, delay time.Duration) {
	go func(url string, delay time.Duration) {
		time.Sleep(delay)
		OpenBrowser(url)
	}(url, delay)
}
