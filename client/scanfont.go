package main

import (
  // "crypto/sha1"
  "os"
  "path/filepath"
  "strings"
)

func visitFont(dir string, file os.FileInfo) error {
  // Note: may run on different OS threads

  filename := filepath.Join(dir, file.Name())
  fext := strings.ToLower(filepath.Ext(filename))

  if _, ok := extToFontType[fext]; !ok {
    // unknown file type -- skip
    return nil
  }

  // reldir := dir
  // if len(dir) > len(config.FontDir) {
  //   reldir = dir[len(config.FontDir)+1:]
  // }
  // L.Printf("- %s/%s\n", reldir, file.Name())

  // contents, err := ioutil.ReadFile(filename)
  // if err != nil {
  //   return err
  // }
  // checksum := sha1.Sum(contents)
  // L.Printf("  checksum: %x\n", checksum)
  // ... = parseFont(bytes.NewReader(contents))

  f := FontFile{}

  err := f.ParseFile(filename)
  if err != nil {
    L.Printf("\nfailed to parse %s: %s\n\n", filename, err)
    return nil
  }

  // L.Printf("font: %+v\n", font)

  return nil
}


func scanFonts(dir string) error {
  return NewFSScanner(dir, visitFont).Scan()
}
