package main

import "path/filepath"

func systemConfigFile() string {
  return filepath.Join(homeDir, ".fontctrl")
}
