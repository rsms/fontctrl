package main

import "path/filepath"

func systemFontDir() string {
  return filepath.Join(homeDir, "Library", "Fonts")
}
