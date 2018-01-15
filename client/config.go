package main

import (
  "encoding/json"
  "errors"
  "fmt"
  "os"
)

type Config struct {
  Repos   []*Repo `json:"repos"`
  FontDir string `json:"font_dir"`
}


func loadConfig(c *Config, filename string) error {
  f, err := os.Open(filename)
  if err != nil {
    return err
  }
  defer f.Close()
  if err := json.NewDecoder(f).Decode(c); err != nil {
    if se, ok := err.(*json.SyntaxError); ok {
      return errors.New(fmt.Sprintf("%s at offset %d", se.Error(), se.Offset))
    }
    return err
  }

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

  return nil
}
