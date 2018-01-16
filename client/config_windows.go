package main

import (
  "os"
  "path/filepath"
)

func getWinDir() string {
  windir := os.Getenv("windir")
  if len(windir) == 0 {
    windir := os.Getenv("SYSTEMROOT")
    if len(windir) == 0 {
      windir := os.Getenv("SYSTEMDRIVE")
      if len(windir) == 0 {
        windir := "C:\\Windows"
      }
    }
  }
  return windir
}

func systemConfigFile() string {
  // %USERPROFILE%\AppData\Local\fontctrl\fontctrl
  return filepath.Join(homeDir,"AppData","Local","fontctrl","fontctrl")
}

func systemFontDir() string {
  // %windir%\fonts
  return filepath.Join(getWinDir(), "fonts")
}
