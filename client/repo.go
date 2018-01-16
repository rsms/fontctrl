package main

import (
  "encoding/json"
  "net/http"
  "fmt"
  "time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

type Repo struct {
  Url   string `json:"url"`
  Index RepoIndex
}

// RepoIndex corresponds to repo/index.json
//
type RepoIndex struct {
  Fonts map[string]*FontIndex `json:"fonts"`
}

// FontIndex corresponds to entries in "fonts" of repo/index.json
//
type FontIndex struct {
  Repo     *Repo      `json:"-"`  // pointer to owning Repo
  Id       string     `json:"id"` // ignored when parsed in RepoIndex
  Family   string     `json:"name"`
  Versions []*Version `json:"versions"`  // sorted latest -> oldest

  vinfo    []*FontVersionInfo `json:"-"` // lazy-loaded; order==.Versions
}

// FontVersionInfo corresponds to repo/<fontname>/<fontname>-<version>.json
//
type FontVersionInfo struct {
  Version     *Version `json:"version"`
  Checksum    string   `json:"checksum"`
  Name        string   `json:"name"`
  Styles      []string `json:"styles"`

  // optional
  ArchiveUrl  string   `json:"archive_url"`
  Description string   `json:"description"`
  InfoUrl     string   `json:"info_url"`
  Authors     []string `json:"authors"`
  License     string   `json:"license"`
}


func fetchJson(url string, v interface{}) error {
  res, err := httpClient.Get(url)
  if err != nil {
    return err
  }
  defer res.Body.Close()
  return json.NewDecoder(res.Body).Decode(v)
}


// GetVersionInfoAt returns info for the corresponding version in f.Versions
//
func (f *FontIndex) GetVersionInfoAt(i int) (*FontVersionInfo, error) {
  if f.vinfo == nil {
    if len(f.Versions) == 0 {
      return nil, nil
    }
    f.vinfo = make([]*FontVersionInfo, len(f.Versions))
  }

  fvi := f.vinfo[i]  // ok to crash on out of bounds
  if fvi != nil {
    return fvi, nil
  }

  ver := f.Versions[i]
  url := fmt.Sprintf("%s%s/%s-%s.json", f.Repo.Url, f.Id, f.Id, ver)
  L.Printf("fetching %s", url)

  fvi = &FontVersionInfo{}
  if err := fetchJson(url, fvi); err != nil {
    return nil, err
  }
  f.vinfo[i] = fvi

  return fvi, nil
}


func (r *Repo) String() string {
  if r == nil {
    return "<nil Repo>"
  }
  return r.Url
}


func (r *Repo) Update() error {
  if err := fetchJson(r.Url + "index.json", &r.Index); err != nil {
    return err
  }

  for id, f := range r.Index.Fonts {
    f.Id = id
    f.Repo = r
    SortVersions(f.Versions)
  }

  return nil
}
