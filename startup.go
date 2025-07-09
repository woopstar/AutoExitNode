package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// getStartupShortcutPath returns the path to the startup shortcut.
func getStartupShortcutPath() string {
	return filepath.Join(startupDir, "AutoExitNode.lnk")
}

// isStartupEnabled checks if the startup shortcut exists.
func isStartupEnabled() bool {
	_, err := os.Stat(getStartupShortcutPath())
	return err == nil
}

// addStartupShortcut creates a Windows shortcut for autostart.
func addStartupShortcut() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	shellObj, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		fmt.Println("CreateObject error:", err)
		return
	}
	defer shellObj.Release()

	wshell, err := shellObj.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		fmt.Println("QueryInterface error:", err)
		return
	}
	defer wshell.Release()

	shortcutPath := getStartupShortcutPath()
	sc, err := oleutil.CallMethod(wshell, "CreateShortcut", shortcutPath)
	if err != nil {
		fmt.Println("CreateShortcut error:", err)
		return
	}
	defer sc.Clear()

	shortcut := sc.ToIDispatch()
	_, _ = oleutil.PutProperty(shortcut, "TargetPath", exePath)
	_, _ = oleutil.PutProperty(shortcut, "WorkingDirectory", filepath.Dir(exePath))
	_, _ = oleutil.PutProperty(shortcut, "WindowStyle", 7)
	_, _ = oleutil.CallMethod(shortcut, "Save")
}

// removeStartupShortcut deletes the autostart shortcut.
func removeStartupShortcut() {
	path := getStartupShortcutPath()
	if err := os.Remove(path); err != nil {
		fmt.Println("Failed to remove startup shortcut:", err)
	}
}
