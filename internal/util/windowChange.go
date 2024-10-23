package util

import (
	"syscall"
	"time"
	"unsafe"
)

// Windows API function bindings
var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW      = user32.NewProc("GetWindowTextW")
)

// returns the GetForegroundWindowhandle of the currently focused window
func getForegroundWindow() syscall.Handle {
	hwnd, _, _ := procGetForegroundWindow.Call()
	return syscall.Handle(hwnd)
}

// GetWindowText retrieves the window title of the given window handle (HWND)
func getWindowText(hwnd syscall.Handle) string {
	buf := make([]uint16, 256)
	procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf)
}

// Monitor for window changes
func MonitorWindowChange() string {
	// Get the current foreground window
	currentWindow := getForegroundWindow()

	// Update the last window
	lastWindow := currentWindow

	for {
		// Get the current foreground window
		currentWindow = getForegroundWindow()

		// Check if the current window is different from the last window
		if currentWindow != lastWindow {
			// Retrieve the window title
			windowTitle := getWindowText(currentWindow)

			// return new title
			return windowTitle
		}

		// Poll every 500 milliseconds
		time.Sleep(500 * time.Millisecond)
	}
}
