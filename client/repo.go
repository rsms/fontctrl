package main

import (
  "encoding/json"
  "net/http"
  // "fmt"
  "time"
)

type Repo struct {
  Url string `json:"url"`
  Index RepoIndex
}


type RepoIndex struct {
  Fonts map[string]RepoIndexFont `json:"fonts"`
}

type RepoIndexFont struct {
  Name string `json:"name"`
  Versions []string `json:"versions"`
}


var httpClient = &http.Client{Timeout: 30 * time.Second}


func (r *Repo) String() string {
  if r == nil {
    return "<nil Repo>"
  }
  return r.Url
}


func (r *Repo) Update() error {
  indexUrl := r.Url + "/index.json"

  res, err := httpClient.Get(indexUrl)
  if err != nil {
    return err
  }

  defer res.Body.Close()

  return json.NewDecoder(res.Body).Decode(&r.Index)
}
