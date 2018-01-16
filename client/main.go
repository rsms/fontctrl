package main

import (
  "flag"
  "fmt"
  "log"
  "os"
)

// set at compile time
var version string = "0.0.0"
var versionGit string = "?"

var progname string
var L *log.Logger
var config Config


func init() {
  progname = os.Args[0]
  L = log.New(os.Stdout, "", log.Ltime)
  // L = log.New(os.Stdout, "", log.Ldate | log.Ltime | log.LUTC)
}


func updateRepos() {
  for _, r := range config.Repos {
    L.Printf("updating repo %s\n", r)
    if err := r.Update(); err != nil {
      L.Fatal(err)
    }
  }
}


func cmd_sync(args []string) {
  if len(args) > 0 {
    L.Fatalf("'%s sync' does not accept any arguments\n", progname)
  }
  // opt := flag.NewFlagSet(progname + " sync", flag.ExitOnError)
  // opt.Parse(args)

  updateRepos()

  L.Printf("scanning fonts in %s\n", config.FontDir)
  var local LocalFontIndex
  if err := local.Scandir(config.FontDir); err != nil {
    if pe, ok := err.(*os.PathError); ok && pe != nil {
      // not found -- continue
    } else {
      L.Fatal(err)
    }
  }

  for fid, vpattern := range config.Fonts {
    findex := config.FindFontIndex(fid)
    if findex == nil {
      L.Printf("error: font \"%s\" not found in any repository\n", fid)
      continue
    }

    L.Printf("found %s (%s) => %+v in repo %s\n",
      fid, vpattern.String(), findex, findex.Repo)

    // matching version
    i, latever := vpattern.Match(findex.Versions)
    if i == -1 {
      L.Printf("no matching version for %s\n", fid)
      continue
    }

    L.Printf("latest version for %s => %s\n", fid, latever)

    // get font info
    finfo, err := findex.GetVersionInfoAt(i)
    if err != nil {
      L.Fatal(err)
    }
    L.Printf("findex.GetInfo() => %+v\n", finfo)

    // find local
    locals := local.FindFamily(findex.Family)
    for _, lf := range locals {
      L.Printf("local font: %+v\n", lf.Style)
    }
  }
}


func cmd_version(_ []string) {
  fmt.Fprintf(
    os.Stderr,
    "Font Control v%s (%s) https://fontctrl.org/\n",
    version, versionGit)
  os.Exit(0)
}


func main() {
  // parse CLI options
  flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s [options] <command>\n", progname)
    fmt.Fprintf(os.Stderr, "\nCommands:\n")
    fmt.Fprintf(os.Stderr, "  sync     Sync repositories and update fonts\n")
    fmt.Fprintf(os.Stderr, "  version  Print version and exit\n")
    fmt.Fprintf(os.Stderr, "\nOptions:\n")
    flag.PrintDefaults()
  }
  var configFile string
  flag.StringVar(&configFile, "config", "", "Config file")
  flag.Parse()

  if flag.NArg() == 0 { // no <command>
    flag.Usage()
    os.Exit(1)
  }

  // load configuration
  var err error
  if len(configFile) > 0 {
    err = config.LoadFile(configFile)
  } else {
    err = config.LoadBestFile()
  }
  if err != nil {
    L.Fatalf("failed to read config file: %v", err)
  }
  // L.Printf("config: %+v", config)

  // dispatch to command function
  cmd := flag.Arg(0)
  args := flag.Args()[1:]
  
  switch cmd {
    case "sync":    cmd_sync(args)
    case "version": cmd_version(args)
    default:
      L.Fatalf("Unknown command %s\nSee %s -h for help\n", cmd, progname)
  }
}
