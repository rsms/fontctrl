package main

import (
  "errors"
  "os"
  "regexp"
  // "strings"
  "github.com/ConradIrwin/font/sfnt"
)

type FontType int

const (
  FontTypeTTF = FontType(iota)
  FontTypeOTF
)

type FontFile struct {
  Family  string   `json:"family"`
  Style   string   `json:"style"`
  Version Version  `json:"version"`
  FontID  string   `json:"uid"`  // == fontId record of name table
}


var extToFontType map[string]FontType
var uidRegExpEnd, uidRegExpAny *regexp.Regexp

func init() {
  extToFontType = make(map[string]FontType)
  extToFontType[".ttf"] = FontTypeTTF
  extToFontType[".ttx"] = FontTypeTTF
  extToFontType[".otf"] = FontTypeOTF

  uidRegExpEnd = regexp.MustCompile(`(?i)\b([A-Fa-f0-9][A-Fa-f0-9\-.]*)\s*$`)
  uidRegExpAny = regexp.MustCompile(`(?i)\b([A-Fa-f0-9][A-Fa-f0-9\-.]*)\b`)
}


func (f *FontFile) ParseFile(filename string) error {
  fp, err := os.Open(filename)
  if err != nil {
    return err
  }
  defer fp.Close()
  return f.Parse(fp)
}


func (f *FontFile) Parse(file sfnt.File) error {
  font, err := sfnt.Parse(file)
  if err != nil {
    return err
  }

  namet := font.NameTable()
  if namet == nil {
    return errors.New("missing name table")
  }

  var version, uid string

  for _, ent := range namet.List() {
    switch ent.NameID {
    
    case sfnt.NamePreferredFamily:
      if len(f.Family) == 0 {
        // L.Printf("- family: %s\n", ent.String())
        f.Family = ent.String()
      }

    case sfnt.NamePreferredSubfamily:
      if len(f.Style) == 0 {
        // L.Printf("- style: %s\n", ent.String())
        f.Style = ent.String()
      }

    case sfnt.NameVersion:
      version = ent.String()

    case sfnt.NameUniqueIdentifier:
      uid = ent.String()

    }
  }

  if len(f.Family) == 0 {
    // maybe font is missing typoPreferredFamily
    f.Family = findFontNameValue(namet, sfnt.NameFontFamily)
    if len(f.Family) == 0 {
      return errors.New("no family name")
    }
  }

  if len(f.Style) == 0 {
    // maybe font is missing typoPreferredSubfamily
    f.Style = findFontNameValue(namet, sfnt.NameFontSubfamily)
    if len(f.Style) == 0 {
      return errors.New("no subfamily/style name")
    }
  }

  // parse version
  err = parseFontVersion(version, uid, &f.Version)
  if err != nil {
    return err
  }

  f.FontID = uid  

  // if strings.Contains(uid, "Inter UI") {
  //   L.Printf("- uid: %s; font = %+v\n", uid, f)
  // }

  return nil
}


func findFontNameValue(namet *sfnt.TableName, nameID sfnt.NameID) string {
  for _, ent := range namet.List() {
    if ent.NameID == nameID {
      return ent.String()
    }
  }
  return ""
}


func parseFontVersion(version, uid string, v *Version) error {
  if err := v.Parse(version); err != nil {
    return err
  }

  if len(v.Build) == 0 && len(uid) > 0 {
    // try finding a version build metadata in the font's unique id
    m := uidRegExpEnd.FindStringSubmatch(uid)
    if len(m) == 0 {
      m = uidRegExpAny.FindStringSubmatch(uid)
    }
    if len(m) > 0 {
      v.Build = m[1]
    }
  }

  return nil
}

