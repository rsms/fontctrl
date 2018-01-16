package main

import (
  "errors"
  "fmt"
  "regexp"
  "strconv"
  "strings"
  "sort"
  "encoding/json"
)

var versionRegExp *regexp.Regexp

func init() {
  versionRegExp = regexp.MustCompile(
    `(?i)(?:^\s*\*|\b([\d\*]+))` +  // 1 major
    `(?:\s*\.\s*([\d\*]+)` +  // 2 minor
      `(?:\s*\.\s*([\d\*]+))?` +  // 3 build
    `)?` +
    `(?:` +
      `(?:(?:\b|)+([^A-Za-z0-9]+)([A-Za-z0-9\-\.]*)|([A-Za-z0-9\-\.]+))` +
        // 4 tag prefix?, 5 build or prerel, or 6 build compl to 3
      `|$|[^\d\s]|\b` +
    `)?` +
    `\s*`,
  )
}

type VersionOperator int
const (
  Any = VersionOperator(iota)
          // * | ""
  Eq      // =
  Gt      // >
  GtEq    // >=
  Lt      // <
  LtEq    // <=
  Latest  // latest -- includes any Prerel and/or Build
)

// VersionSpec represents a pattern that matches certain versions
//
type VersionPattern struct {
  Version *Version
  Op      VersionOperator
}


// Match finds the most recent version in versions that matches p.
// Assumes versions are sorted from most recent to least recent.
// Returns -1,nil if none matches.
//
func (p *VersionPattern) Match(versions []*Version) (int, *Version) {
  if len(versions) == 0 {
    return -1, nil
  }

  if p.Op == Latest {
    return 0, versions[0]
  }

  if p.Op == Any && p.Version == nil {
    for i, v := range versions {
      if len(v.Prerel) == 0 {
        return i, v
      }
    }
    return -1, nil
  }

  pv := p.Version  // non-nil when p.Op!=Any
  for i, v := range versions {
    switch v.Compare(pv) {
      case 0: // v == pv
        if p.Op == Eq || p.Op == LtEq || p.Op == GtEq {
          return i, v
        }
      case -1: // v < pv
        if p.Op == Lt || p.Op == LtEq {
          return i, v
        }
      case 1: // v > pv
        if p.Op == Gt || p.Op == GtEq {
          return i, v
        }
    }
  }

  return -1, nil
}


func (p *VersionPattern) UnmarshalJSON(b []byte) error {
  var s string
  if err := json.Unmarshal(b, &s); err != nil {
    return err
  }
  return p.Parse(s)
}

func (p VersionPattern) MarshalJSON() ([]byte, error) {
  return json.Marshal(p.String())
}


func (p *VersionPattern) UnmarshalYAML(u func(interface{}) error) error {
  var s string
  if err := u(&s); err != nil {
    return err
  }
  return p.Parse(s)
}

func (p *VersionPattern) MarshalYAML() (interface{}, error) {
  return p.String(), nil
}


func (p *VersionPattern) Parse(s string) error {
  p.Version = nil

  i, z := 0, len(s)
  if z == 0 || (z == 1 && s[i] == '*') {
    p.Op = Any
    return nil
  }

  if s == "latest" {
    p.Op = Latest
    return nil
  }

  p.Op = Eq  // default when operator is absent

  parse_op:
  switch s[i] {
    case '=': p.Op = Eq
    // case '*', '-', '+': p.Op = Any
    case '>':
      i++
      p.Op = Gt
      if len(s) > i && s[i] == '=' {
        p.Op = GtEq
      }
    case '<':
      i++
      p.Op = Lt
      if len(s) > i && s[i] == '=' {
        p.Op = LtEq
      }
    case ' ': // skip whitespace
      i++
      if i < z {
        goto parse_op
      }
  }

  if i < len(s) {
    p.Version = &Version{}
    p.Version.Parse1(s[i:], -1)
  } else if p.Op != Any {
    return fmt.Errorf(
      "invalid version pattern \"%s\"; expecting version number or tag", s)
  }

  return nil
}


func (p *VersionPattern) String() string {
  var ops string
  switch p.Op {
    case Any:    return "*"
    case Latest: return "latest"
    // case Eq:  ops = ""
    case Gt:     ops = ">"
    case GtEq:   ops = ">="
    case Lt:     ops = "<"
    case LtEq:   ops = "<="
  }
  if p.Version != nil {
    return fmt.Sprintf("%s%s", ops, p.Version.String())
  }
  return ops
}


// Match compares a list of version strings and returns the offset to
// the version that matches v, or -1 if none matches.
//
func (p *VersionPattern) Matches(v *Version) bool {
  // v := Version{}

  // for i, vstr := range vv {
  //   v
  // }

  return true
}



type Version struct {
  Major  int32  `json:"major"`
  Minor  int32  `json:"minor"`
  Patch  int32  `json:"patch"`
  Prerel string `json:"prerel"`
  Build  string `json:"build"`
}

// Compare compares two versions a <=> b and returns
//
// -1 if a < b
//  1 if a > b
//  0 if a == b
//
// Handles wildcard versions too, e.g.
//   "2" == "2.0.0"
//   "2" < "1.2.3"
//   "2.0-beta" < "2.0.0"
//   "2.*.5" < "2.0.6"
//
func (a *Version) Compare(b *Version) int {
  if a.Major >= 0 && b.Major >= 0 {
    if a.Major < b.Major {
      return -1
    }
    if b.Major < a.Major {
      return 1
    }
  }

  if a.Minor >= 0 && b.Minor >= 0 {
    if a.Minor < b.Minor {
      return -1
    }
    if b.Minor < a.Minor {
      return 1
    }
  }

  if a.Patch >= 0 && b.Patch >= 0 {
    if a.Patch < b.Patch {
      return -1
    }
    if b.Patch < a.Patch {
      return 1
    }
  }
  
  // non-empty prerel has lower precedence than empty prerel
  if len(b.Prerel) == 0 && len(a.Prerel) > 0 {
    // e.g. "1.2.3-beta > 1.2.3" => -1
    return -1  // a is lesser than b
  }
  if len(a.Prerel) == 0 && len(b.Prerel) > 0 {
    // e.g. "1.2.3 > 1.2.3-beta" => 1
    return 1  // a is greater than b
  }
  if len(a.Prerel) > 0 && len(b.Prerel) > 0 {
    if a.Prerel < b.Prerel {
      return -1
    }
    if b.Prerel < a.Prerel {
      return 1
    }
  }

  // string comparison for build
  if a.Build < b.Build {
    // e.g. "1.2.3+abc" < "1.2.3+def"
    return -1
  }
  if b.Build < a.Build {
    return 1
  }

  // equivalent
  return 0
}

// Sortable version lists
type VersionList []*Version
func (a VersionList) Len() int           { return len(a) }
func (a VersionList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VersionList) Less(i, j int) bool { return a[i].Compare(a[j]) < 0 }


func SortVersions(v []*Version) {
  sort.Sort(VersionList(v))
}


func (v *Version) UnmarshalJSON(b []byte) error {
  var s string
  if err := json.Unmarshal(b, &s); err != nil {
    return err
  }
  return v.Parse(s)
}


func (v Version) MarshalJSON() ([]byte, error) {
  return json.Marshal(v.String())
}


func ParseVersion(version string) (*Version, error) {
  v := Version{}
  if err := v.Parse(version); err != nil {
    return nil, err
  }
  return &v, nil
}


func (v *Version) String() string {
  var s string

  if v.Major > -1 {
    if v.Minor > -1 {
      if v.Patch > -1 {
        s = fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
      } else {
        s = fmt.Sprintf("%d.%d", v.Major, v.Minor)
      }
    } else if v.Patch > -1 {
      s = fmt.Sprintf("%d.*.%d", v.Major, v.Patch)
    } else {
      s = fmt.Sprintf("%d", v.Major)
    }
  } else if v.Minor > -1 {
    if v.Patch > -1 {
      s = fmt.Sprintf("*.%d.%d", v.Minor, v.Patch)
    } else {
      s = fmt.Sprintf("*.%d", v.Minor)
    }
  } else if v.Patch > -1 {
    s = fmt.Sprintf("*.*.%d", v.Patch)
  }

  hasPrerel := len(v.Prerel) > 0
  hasBuild := len(v.Build) > 0

  if hasPrerel && hasBuild {
    return fmt.Sprintf("%s-%s+%s", s, v.Prerel, v.Build)
  } else if hasPrerel {
    return fmt.Sprintf("%s-%s", s, v.Prerel)
  } else if hasBuild {
    return fmt.Sprintf("%s+%s", s, v.Build)
  } else {
    return s
  }
}


func (v *Version) Parse(version string) error {
  return v.Parse1(version, 0)
}


func (v *Version) Parse1(version string, defaultv int32) error {
  m := versionRegExp.FindStringSubmatch(version)
  if len(m) == 0 {
    return errors.New("invalid format " + version)
  }
  // fmt.Printf("\"%s\" => \"%s\"\n", version, strings.Join(m, "\", \""))

  var n uint64
  var err error

  if len(m[1]) > 0 {  // Note: "*." prefix means m[1] is empty
    n, err = strconv.ParseUint(m[1], 10, 31)
    if err != nil {
      return err
    }
    v.Major = int32(n)
  } else {
    v.Major = defaultv
  }
  v.Minor = defaultv
  v.Patch = defaultv
  v.Prerel = ""
  v.Build = ""

  if len(m[2]) > 0 {
    if m[2][0] != '*' {
      n, err = strconv.ParseUint(m[2], 10, 31)
      if err != nil {
        return err
      }
      v.Minor = int32(n)
    }

    m6 := strings.TrimSpace(m[6])

    if len(m[3]) > 0 {
      if m[3][0] != '*' {
        n, err = strconv.ParseUint(m[3], 10, 31)
        if err != nil {
          return err
        }
        v.Patch = int32(n)
      }

      if len(m6) > 0 {
        // e.g. 1.2.0df73 => (1, 2, 0, "0df73")
        v.Patch = defaultv
        v.Build = m[3] + m6
      }
    } else {
      if len(m6) > 0 {
        // e.g. 1.0df73 => (1, 0, 0, "0df73")
        v.Minor = defaultv
        v.Build = m[2] + m6
      }
    }
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
