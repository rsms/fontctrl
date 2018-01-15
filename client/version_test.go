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
    if err := ParseVersion(input, &v); err == nil {
      actual := v.String()
      if actual != expected {
        t.Errorf("(\"%s\") => \"%s\" ; expected \"%s\"\n", input, actual, expected)
      }
    } else {
      t.Errorf("(\"%s\") => error %v\n", input, err)
    }
  }
}
