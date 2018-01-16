package main

import (
  "os"
  "path/filepath"
  "sync"
  "strings"
)

type LocalFontIndex struct {
  fontsmu       sync.RWMutex  // protects access to fonts and fontsByFamily
  fonts         []*FontFile   // all font files
  fontsByFamily map[string][]*FontFile  // keyed by family name
}


func (l *LocalFontIndex) Scandir(dir string) error {
  return NewFSScanner(dir, l.visitFile).Scan()
}


func (l *LocalFontIndex) FindFamily(family string) []*FontFile {
  l.fontsmu.RLock()
  defer l.fontsmu.RUnlock()
  if l.fontsByFamily != nil {
    if fv, ok := l.fontsByFamily[family]; ok {
      return fv
    }
  }
  return make([]*FontFile, 0)
}


func (l *LocalFontIndex) visitFile(dir string, file os.FileInfo) error {
  // Note: may run on different OS threads

  // L.Printf("%s: ModTime: %v, Size: %d\n",
  //   file.Name(), file.ModTime(), file.Size())

  filename := filepath.Join(dir, file.Name())
  fext := strings.ToLower(filepath.Ext(filename))

  if _, ok := extToFontType[fext]; !ok {
    // unknown file type -- skip
    return nil
  }

  f := &FontFile{}

  err := f.ParseFile(filename)
  if err != nil {
    L.Printf("\nfailed to parse %s: %s\n\n", filename, err)
    return nil
  }

  // if strings.Contains(f.Family, "Inter UI") {
  //   L.Printf("+ \"%s\", \"%s\"\n", f.Family, f.Style)
  // }

  l.fontsmu.Lock()
  defer l.fontsmu.Unlock()

  l.fonts = append(l.fonts, f)

  if l.fontsByFamily == nil {
    l.fontsByFamily = make(map[string][]*FontFile)
  }
  v, _ := l.fontsByFamily[f.Family]
  l.fontsByFamily[f.Family] = append(v, f)

  return nil
}

