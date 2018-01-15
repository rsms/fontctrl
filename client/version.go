package main

import (
  "errors"
  "fmt"
  "regexp"
  "strconv"
  "strings"
)

type Version struct {
  Major  uint32 `json:"major"`
  Minor  uint32 `json:"minor"`
  Patch  uint32 `json:"patch"`
  Prerel string `json:"prerel"`
  Build  string `json:"build"`
}

var versionRegExp *regexp.Regexp

func init() {
  versionRegExp = regexp.MustCompile(
    `(?i)\b(\d+)` +  // 1 major
    `(?:\s*\.\s*(\d+)` +  // 2 minor
      `(?:\s*\.\s*(\d+))?` +  // 3 build
    `)?` +
    `(?:` +
      `(?:\b+([^A-Za-z0-9]*)([A-Za-z0-9\-\.]+)|([A-Za-z0-9\-\.]+))` +
        // 4 tag prefix?, 5 build or prerel, or 6 build compl to 3
      `|$|[^\d\s]|\b` +
    `)?` +
    `\s*`,
  )
}


func (v *Version) String() string {
  hasPrerel := len(v.Prerel) > 0
  hasBuild := len(v.Build) > 0
  if hasPrerel && hasBuild {
    return fmt.Sprintf(
      "%d.%d.%d-%s+%s", v.Major, v.Minor, v.Patch, v.Prerel, v.Build)
  } else if hasPrerel {
    return fmt.Sprintf("%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.Prerel)
  } else if hasBuild {
    return fmt.Sprintf("%d.%d.%d+%s", v.Major, v.Minor, v.Patch, v.Build)
  } else {
    return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
  }
}


func ParseVersion(version string, v *Version) error {
  m := versionRegExp.FindStringSubmatch(version)
  if len(m) == 0 {
    return errors.New("invalid format " + version)
  }
  // fmt.Printf("\"%s\" => \"%s\"\n", version, strings.Join(m, "\", \""))

  var n uint64
  var err error

  n, err = strconv.ParseUint(m[1], 10, 32)
  if err != nil {
    return err
  }
  v.Major = uint32(n)
  v.Prerel = ""
  v.Build = ""

  if len(m[2]) > 0 {
    n, err = strconv.ParseUint(m[2], 10, 32)
    if err != nil {
      return err
    }
    v.Minor = uint32(n)

    m6 := strings.TrimSpace(m[6])

    if len(m[3]) > 0 {
      n, err = strconv.ParseUint(m[3], 10, 32)
      if err != nil {
        return err
      }
      v.Patch = uint32(n)

      if len(m6) > 0 {
        // e.g. 1.2.0df73 => (1, 2, 0, "0df73")
        v.Patch = 0
        v.Build = m[3] + m6
      }
    } else {
      v.Patch = 0
      if len(m6) > 0 {
        // e.g. 1.0df73 => (1, 0, 0, "0df73")
        v.Minor = 0
        v.Build = m[2] + m6
      }
    }
  } else {
    v.Minor = 0
  }

  m5 := strings.TrimSpace(m[5])
  if len(m5) > 0 {
    if len(m[4]) == 1 && m[4][0] == '-' {
      // e.g. "1.2.0-xyz"
      v.Prerel = m5
    } else {
      // e.g. "1.2.0+xyz", "1.2.0;xyz", etc
      v.Build = m5
    }
  }

  return nil
}
