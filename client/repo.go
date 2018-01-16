package main

import (
  "encoding/json"
  Path "path"
  "net/http"
  "fmt"
  "time"
  "strings"
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
  if res.StatusCode < 200 || res.StatusCode > 299 {
    return fmt.Errorf("%d %s (GET %s)",
      res.StatusCode, http.StatusText(res.StatusCode), url)
  }
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
  url, err := f.Repo.GetUrl(fmt.Sprintf("%s/%s-%s.json", f.Id, f.Id, ver))
  if err != nil {
    return nil, err
  }

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


func withTrailingSlash(s string) string {
  if s[len(s)-1] != '/' {
    return s + "/"
  }
  return s
}

func parseGithubRepoUrl(s string, path string) (string, error) {
  // expect "user/repo/path?#branch?"
  br := "master"
  var ps string
  i, p := 0, 0
  for i < len(s) {
    if s[i] == '/' {
      if p == 1 {
        ps = s[i:]
        s = s[:i]
        if x := strings.IndexByte(ps, '#'); x > -1 {
          br = ps[x+1:]
          ps = ps[:x]
        }
        break
      }
      p++
    }
    i++
  }

  if p == 0 {
    // no slash found
    return s, fmt.Errorf(
      "invalid github repo url \"%s\"; expected user/repo", s)
  }

  if len(ps) == 0 {
    if p = strings.IndexByte(s, '#'); p > -1 {
      br = s[p+1:]
      s = s[:p]
    }
  }

  s = "https://raw.githubusercontent.com/" + Path.Join(s, br, ps, path)

  return s, nil
}


func (r *Repo) GetUrl(path string) (string, error) {
  s := r.Url
  p := strings.IndexByte(s, ':')
  if p == -1 {
    return s, fmt.Errorf("invalid repo url \"%s\"; missing prototcol", s)
  }
  proto := s[:p]
  switch proto {
    case "http": // pass
    case "https": // pass
    case "github": return parseGithubRepoUrl(s[p+1:], path)
    default:
      return s, fmt.Errorf(
        "can not understand repo url \"%s\"; unknown protocol", s)
  }
  return withTrailingSlash(s) + path, nil
}


func (r *Repo) Update() error {
  url, err := r.GetUrl("index.json")
  if err != nil {
    return err
  }

  L.Printf("url: %s\n", url)
  if err := fetchJson(url, &r.Index); err != nil {
    return err
  }

  for id, f := range r.Index.Fonts {
    f.Id = id
    f.Repo = r
    SortVersions(f.Versions)
  }

  return nil
}
