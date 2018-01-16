package main

import (
  "io/ioutil"
  "os"
  "os/user"
  "path/filepath"
  "strings"
  "gopkg.in/yaml.v2"
)

type FontSubscription struct {
  VersionPattern   `json:"version,omitempty" yaml:"version,omitempty"`
  Styles  []string `json:"repos,omitempty" yaml:"repos,omitempty"`
}

// similar type used only for YAML encoding
type fontSubscription2 struct {
  Version *VersionPattern `yaml:"version"`
  Styles []string         `yaml:"styles"`
}

func (p *FontSubscription) UnmarshalYAML(u func(interface{}) error) error {
  // handle two different forms:
  // - <string> -- e.g. ">=1.2"
  // - { version?: v string, styles?: []string }
  if err := p.VersionPattern.UnmarshalYAML(u); err != nil {
    st := fontSubscription2{ Version: &p.VersionPattern }
    if err := u(&st); err != nil {
      return err
    }
    p.Styles = st.Styles
  }

  return nil
}

func (p *FontSubscription) MarshalYAML() (interface{}, error) {
  if p.Styles == nil || len(p.Styles) == 0 {
    return p.VersionPattern, nil
  }
  return fontSubscription2{
    Version: &p.VersionPattern,
    Styles: p.Styles,
  }, nil
}


type Config struct {
  File    string  `json:"-" yaml:"-"`
  FontDir string  `json:"font_dir,omitempty" yaml:"font-dir,omitempty"`
  Repos   []*Repo `json:"repos,omitempty" yaml:"repos,omitempty"`
  Fonts  map[string]FontSubscription `json:"fonts" yaml:"fonts"`
}

const defaultRepoURL = "https://fontctrl.org/fonts/"
var homeDir string

func init() {
  usr, err := user.Current()
  if err == nil {
    homeDir = usr.HomeDir
  }
}

func defaultFontDir() string {
  return systemFontDir()
}

// findFontIndex finds the FontIndex for the font identified by fid.
// It searches c.Repos in order and returns the first match,
// or nil if not found.
//
func (c *Config) FindFontIndex(fid string) *FontIndex {
  for _, r := range c.Repos {
    if findex, ok := r.Index.Fonts[fid]; ok {
      return findex
    }
  }
  return nil
}

// InitDefault initializes a config struct to the state of the "built in"
// configuration.
//
func (c *Config) InitDefault() {
  c.File = "<builtin>"
  c.Repos = make([]*Repo, 1)
  c.Repos[0] = &Repo{ Url: defaultRepoURL }
  c.FontDir = defaultFontDir()
  c.Fonts = nil
  c.init2()
}


// LoadFile loads a configuration from a YAML file
//
func (c *Config) LoadFile(filename string) error {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    return err
  }
  if err := yaml.Unmarshal(data, &c); err != nil {
    return err
  }
  c.File = filename
  c.init2()
  return nil
}

// LoadBestFile loads a configuration from a JSON file located in one of a
// few predefined locations:
//  1. .fontctrl.yml
//  2. fontctrl.yml
//  3. ~/.fontctrl.{yml,yaml} (macOS, linux)
//     %USERPROFILE%\AppData\Local\fontctrl\fontctrl.{yml,yaml} (windows)
//  4. <built-in> -- calls c.InitDefault()
//
func (c *Config) LoadBestFile() error {
  filename := ".fontctrl.yml"
  attempt := 0

  retry:
  attempt++
  err := c.LoadFile(filename)
  if err == nil {
    return nil
  }
  
  if pe, ok := err.(*os.PathError); !ok || pe == nil {
    return err
  }
  
  // case: file not found
  switch attempt {
    case 1:
      filename = "fontctrl.yml"
      goto retry
    case 2:
      filename = systemConfigFile() + ".yml"
      goto retry
    case 3:
      filename = systemConfigFile() + ".yaml"
      goto retry
    default:
      // Not found -- fall back on built-in default config
      // return fmt.Errorf(
      //   "file not found at any of the following locations:\n - %s",
      //   strings.Join(attempts, "\n - "))
      c.InitDefault()
  }

  return nil
}


// init2 is run after a config has been initialized from user data
// in order to check and normalize the config state.
//
func (c *Config) init2() {
  // check to make sure there're no null repos
  var nullIndices []int
  for i, r := range c.Repos {
    if r == nil {
      nullIndices = append(nullIndices, i)
    }
  }
  if nullIndices != nil {
    for _, i := range nullIndices {
      c.Repos = append(c.Repos[:i], c.Repos[i+1:]...)
    }
  }

  // no repo? use default
  if c.Repos == nil {
    c.Repos = make([]*Repo, 1)
    c.Repos[0] = &Repo{ Url: defaultRepoURL }
  } else if len(c.Repos) == 0 {
    c.Repos[0] = &Repo{ Url: defaultRepoURL }
  }

  // fontdir
  if len(c.FontDir) == 0 {
    c.FontDir = defaultFontDir()
  } else {
    p := strings.Index(c.FontDir, "~/")
    if p == -1 {
      p = strings.Index(c.FontDir, "~\\") // windows
    }
    if p != -1 {
      c.FontDir = c.FontDir[:p] + homeDir + c.FontDir[p+1:]
    }
    c.FontDir = filepath.Clean(c.FontDir)
  }
}
