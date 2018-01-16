package main

import "path/filepath"

func systemFontDir() string {
  return filepath.Join(homeDir, ".fonts", "truetype")

  // Note: is the "truetype" subdirectory really needed?

  // May need to parse fonts.conf
  // See https://www.freedesktop.org/software/fontconfig/fontconfig-user.html
}
