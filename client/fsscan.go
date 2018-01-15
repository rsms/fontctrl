package main

import (
  "path/filepath"
  "io/ioutil"
  "os"
  "sync"
)

type FSVisitor func(dir string, f os.FileInfo)error

type FSScanner struct {
  rootDir string
  visitor FSVisitor
  visited map[string]struct{}
  errch   chan error
  resch   chan struct{}
}


func NewFSScanner(dir string, f FSVisitor) *FSScanner {
  return &FSScanner{
    rootDir: dir,
    visitor: f,
    visited: make(map[string]struct{}),
    errch: make(chan error),
    resch: make(chan struct{}),
  }
}


func (s *FSScanner) Scan() error {
  return s.scandir(s.rootDir)
}


func (s *FSScanner) scandir(dir string) error {
  files, err := ioutil.ReadDir(dir)
  if err != nil {
    // TODO: ignore "no read access" errors for subdirectories
    return err
  }

  var wg sync.WaitGroup
  var errch chan error

  for _, f := range files {
    if f.IsDir() {
      dir2 := filepath.Join(dir, f.Name())
      if _, ok := s.visited[dir2]; ok {
        // we've visited this directory already -- skip
        continue
      }
      s.visited[dir2] = struct{}{} // mark as visited

      if errch == nil {
        errch = make(chan error, 1)
      }

      wg.Add(1)
      go func() {
        if err := s.scandir(dir2); err != nil {
          select {
            case errch <- err:
            default:
          }
        }
        wg.Done()
      }()

    } else {
      s.visitor(dir, f)
    }
  }

  if errch != nil {
    wg.Wait()
    // grab first error that occured, if any
    select {
      case err := <- errch:
        return err
      default: // no error
    }
  }

  return nil
}

