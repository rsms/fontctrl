package main

import "testing"

func TestParseVersion(t *testing.T) {
  successCases := [][]string{
    // well-formed semver
    []string{"2.3.4",        "2.3.4"},
    []string{"2.3.4+xy123",  "2.3.4+xy123"},
    []string{"2.3.4+1.2.3b", "2.3.4+1.2.3b"},
    []string{"2.3.4-beta",   "2.3.4-beta"},

    // common
    []string{"2.003", "2.3.0"},
    []string{"2.003;xyz", "2.3.0+xyz"},
    []string{"2.4;1b5054a", "2.4.0+1b5054a"},
    []string{"2.003 ; xyz", "2.3.0+xyz"},
    []string{"2.003-next", "2.3.0-next"},
    []string{"13.0d3e20", "13.0.0+0d3e20"},
    []string{"13.xd3e20", "13.0.0+xd3e20"},
    []string{"Version 2.003", "2.3.0"},
    []string{"version 2.003", "2.3.0"},
    []string{"version  2.003 ", "2.3.0"},
    []string{"Version 2912.010", "2912.10.0"},

    // uncommon (from real font files)
    []string{"1", "1.0.0"},
    []string{"0", "0.0.0"},
    []string{"001", "1.0.0"},
    []string{"001.001", "1.1.0"},
    []string{"Version 1.06 uh", "1.6.0+uh"},
    []string{
      "Version 2.000;GOOG;noto-source:20170915:90ef993387c0",
      "2.0.0+GOOG"},
    []string{
      "Version 1.00 August 22, 2017, initial release",
      "1.0.0+August"},
    []string{
      "Version 001.003;Core 1.0.01;otf.5.02.2298;42.06W",
      "1.3.0+Core"},
    []string{"Version 009.014; wf-rip", "9.14.0+wf-rip"},
    []string{
      "Version 3.000;PS 1.000;hotconv 1.0.50;makeotf.lib2.0.16970",
      "3.0.0+PS"},
    []string{
      "OTF 1.022;PS 001.001;Core 1.0.31;makeotf.lib1.4.1585",
      "1.22.0+PS"},
  }
  for _, c := range successCases {
    input := c[0]
    expected := c[1]
    var v Version
    if err := v.Parse(input); err == nil {
      actual := v.String()
      if actual != expected {
        t.Errorf("(\"%s\") => \"%s\" ; expected \"%s\"\n",
          input, actual, expected)
      }
    } else {
      t.Errorf("(\"%s\") => error %v\n", input, err)
    }
  }
}

func TestParse1Version(t *testing.T) {
  successCases := [][]string{
    []string{"2",            "2"},
    []string{"2.3",          "2.3"},
    []string{"2.*",          "2"},
    []string{"2.3.*",        "2.3"},
    []string{"2.*.*",        "2"},
    []string{"2.*.4",        "2.*.4"},
    []string{"*.3",          "*.3"},
    []string{"*.*.4",        "*.*.4"},
    []string{"2.3.*+xy123",  "2.3+xy123"},
    []string{"2.3.*-xy123",  "2.3-xy123"},
  }
  for _, c := range successCases {
    input := c[0]
    expected := c[1]
    var v Version
    if err := v.Parse1(input, -1); err == nil {
      actual := v.String()
      if actual != expected {
        t.Errorf("(\"%s\") => \"%s\" ; expected \"%s\"\n",
          input, actual, expected)
      } else {
        t.Logf("(\"%s\") => \"%s\" ; OK\n", input, actual)
      }
    } else {
      t.Errorf("(\"%s\") => error %v\n", input, err)
    }
  }
}

func TestParseVersionPattern(t *testing.T) {
  successCases := [][]string{
    []string{"*",             "*"},

    []string{"=2.3.4+xy123",  "2.3.4+xy123"},
    []string{"2",             "2"},
    []string{"2.3",           "2.3"},
    []string{"2.3.*",         "2.3"},
    []string{"2.3.4",         "2.3.4"},
    []string{"2.3.4-beta",    "2.3.4-beta"},

    []string{">=2",           ">=2"},
    []string{">=2.3",         ">=2.3"},
    []string{">=2.3.4",       ">=2.3.4"},

    []string{">2",            ">2"},
    []string{">2.3",          ">2.3"},
    []string{">2.3.4",        ">2.3.4"},

    []string{"<=2",           "<=2"},
    []string{"<=2.3",         "<=2.3"},
    []string{"<=2.3.4",       "<=2.3.4"},

    []string{"<2",            "<2"},
    []string{"<2.3",          "<2.3"},
    []string{"<2.3.4",        "<2.3.4"},

    // whitespace
    []string{"  = 2. 3.4+xy123", "2.3.4+xy123"},
    []string{"  2 ",             "2"},
    []string{"  2 .3",           "2.3"},
    []string{"  2 .3.*",         "2.3"},
    []string{"  2 .3.4  ",       "2.3.4"},
    []string{"  2 .3.4-beta",    "2.3.4-beta"},

    []string{"  >= 2",           ">=2"},
    []string{"  >= 2.3",         ">=2.3"},
    []string{"  >= 2.3.4",       ">=2.3.4"},

    []string{"  >2 ",            ">2"},
    []string{"  >2 .3",          ">2.3"},
    []string{"  >2 .3.4",        ">2.3.4"},

    []string{"  <=  2",          "<=2"},
    []string{"  <=  2. 3",       "<=2.3"},
    []string{"  <=  2. 3.4",     "<=2.3.4"},

    []string{"  < 2",            "<2"},
    []string{"  < 2.3",          "<2.3"},
    []string{"  < 2.3.4",        "<2.3.4"},

  }
  for _, c := range successCases {
    input := c[0]
    expected := c[1]
    var p VersionPattern
    if err := p.Parse(input); err == nil {
      actual := p.String()
      if actual != expected {
        t.Errorf("(\"%s\") => \"%s\" ; expected \"%s\"\n",
          input, actual, expected)
      } else {
        t.Logf("(\"%s\") => \"%s\" ; OK\n", input, actual)
      }
    } else {
      t.Errorf("(\"%s\") => error %v\n", input, err)
    }
  }
}


func TestCompareVersions(t *testing.T) {
  type Sample struct {
    a, b string
    expected int
  }
  successCases := []Sample{
    // a == b
    Sample{"2.0.0",         "2.0.0",        0},
    Sample{"2.3.0",         "2.3.0",        0},
    Sample{"2.3.0",         "2.3.0",        0},
    Sample{"2.3.4",         "2.3.4",        0},
    Sample{"2.3.4-beta",    "2.3.4-beta",   0},
    Sample{"2.3.4+xy123",   "2.3.4+xy123",  0},

    // a is less than b
    Sample{"2.0.0",         "3.0.0",       -1},
    Sample{"2.3.0",         "2.4.0",       -1},
    Sample{"2.3.4",         "2.3.5",       -1},
    Sample{"2.3.4-alpha",   "2.3.4",       -1},
    Sample{"2.3.4-alpha",   "2.3.4-beta",  -1},

    // a is greater than b
    Sample{"3.0.0",         "2.0.0",        1},
    Sample{"3.0.0",         "2.3.0",        1},
    Sample{"2.4.0",         "2.3.0",        1},
    Sample{"2.3.5",         "2.3.4",        1},
    Sample{"2.3.4",         "2.3.4-alpha",  1},
    Sample{"2.3.4-beta",    "2.3.4-alpha",  1},

    // wildcard versions

    // a == b
    Sample{"2",             "2.0.0",        0},
    Sample{"2",             "2.1.0",        0},
    Sample{"2",             "2.1.2",        0},
    Sample{"2.3",           "2.3.0",        0},
    Sample{"2.3",           "2.3.1",        0},
    Sample{"2.*.4",         "2.3.4",        0},
    Sample{"2.*.4",         "2.1.4",        0},
    Sample{"*.*.4",         "1.1.4",        0},
    Sample{"*.*.4",         "2.2.4",        0},
    Sample{"*.*.4",         "99.99.4",      0},
    Sample{"2-beta",        "2.3.4-beta",   0},
    Sample{"2-beta",        "2.4.0-beta",   0},
    Sample{"2-beta",        "2.30.40-beta", 0},
    Sample{"*-beta",        "2.3.4-beta",   0},
    Sample{"2+xy123",       "2.0.0+xy123",  0},
    Sample{"2+xy123",       "2.1.1+xy123",  0},
    Sample{"2+xy123",       "2.2.2+xy123",  0},

    // a is less than b
    Sample{"2",             "3.0.0",       -1},
    Sample{"2-beta",        "2.0.0",       -1},
    Sample{"2-beta",        "2.9.9",       -1},
    Sample{"2.3",           "3.0.0",       -1},
    Sample{"2.3",           "2.4.0",       -1},
    Sample{"2.3",           "2.4.9",       -1},
    Sample{"2.*.4",         "2.0.5",       -1},
    Sample{"2.*.4",         "2.9.5",       -1},
    Sample{"*.*.4",         "0.0.5",       -1},
    Sample{"*.*.4",         "0.9.5",       -1},
    Sample{"*.*.4",         "9.0.5",       -1},
    Sample{"*.*.4",         "9.9.5",       -1},

    // a is greater than b
    Sample{"3",             "2.0.0",        1},
    Sample{"3",             "3.0.0-beta",   1},
    Sample{"3",             "3.9.9-beta",   1},
    Sample{"3.3",           "2.0.0",        1},
    Sample{"2.4",           "2.3.0",        1},
    Sample{"2.4",           "2.3.9",        1},
    Sample{"3.0.0",         "2.3.0",        1},
    Sample{"2.*.4",         "2.0.3",        1},
    Sample{"2.*.4",         "2.9.3",        1},
    Sample{"*.*.4",         "0.0.3",        1},
    Sample{"*.*.4",         "0.9.3",        1},
    Sample{"*.*.4",         "9.0.3",        1},
    Sample{"*.*.4",         "9.9.3",        1},
  }

  for _, c := range successCases {
    var a, b Version
    if err := a.Parse1(c.a, -1); err != nil {
      t.Errorf("Parse(\"%s\") => error %v\n", c.a, err)
      break
    }

    if err := b.Parse(c.b); err != nil {
      t.Errorf("Parse1(\"%s\") => error %v\n", c.b, err)
      break
    }

    actual := a.Compare(&b)
    if actual != c.expected {
      t.Errorf("\"%s\" <> \"%s\" => %d ; expected %d\n",
        c.a, c.b, actual, c.expected)
    }
    // else { t.Logf("\"%s\" <> \"%s\" => %d ; OK\n", c.a, c.b, actual) }

    // reversed
    actual = b.Compare(&a)
    expected := -c.expected
    if actual != expected {
      t.Errorf("\"%s\" <> \"%s\" => %d (R) ; expected %d\n",
        c.b, c.a, actual, expected)
    }
    // else { t.Logf("\"%s\" <> \"%s\" => %d (R) ; OK\n", c.b, c.a, actual) }
  }
}


// TODO: v.Compare(*Version)
