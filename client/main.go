package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "os/user"
  "path/filepath"
  "strings"
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


func defaultFontDir(homeDir string) string {
  // TODO: different paths on different operating systems
  return filepath.Join(homeDir, "/Library/Fonts")
}


func cmd_update(args []string) {
  if len(args) > 0 {
    L.Fatalf("'%s update' does not accept any arguments\n", progname)
  }
  // opt := flag.NewFlagSet(progname + " update", flag.ExitOnError)
  // opt.Parse(args)
  for _, r := range config.Repos {
    L.Printf("updating repo %s\n", r)
    if err := r.Update(); err != nil {
      L.Fatal(err)
    }
    L.Printf("updated index: %+v", r.Index)
  }

  if err := scanFonts(config.FontDir); err != nil {
    L.Fatal(err)
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
  // resolve computer user
  usr, err := user.Current()
  if err != nil {
    L.Fatal(err)
  }
  if len(usr.HomeDir) == 0 {
    L.Fatal("no home directory for current user")
  }

  // resolve file system paths
  configFile := filepath.Join(usr.HomeDir, ".fontctrl.json")

  // parse CLI options
  flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s [options] <command>\n", progname)
    fmt.Fprintf(os.Stderr, "\nCommands:\n")
    fmt.Fprintf(os.Stderr, "  update   Update the index of all repositories\n")
    fmt.Fprintf(os.Stderr, "  version  Print version and exit\n")
    fmt.Fprintf(os.Stderr, "\nOptions:\n")
    flag.PrintDefaults()
  }
  flag.StringVar(&configFile, "config", configFile, "Config file")
  flag.Parse()

  if flag.NArg() == 0 { // no <command>
    flag.Usage()
    os.Exit(1)
  }

  // load configuration
  if err := loadConfig(&config, configFile); err != nil {
    L.Fatalf("failed to parse config file %s: %v", configFile, err)
  }

  // FontDir
  if len(config.FontDir) == 0 {
    config.FontDir = defaultFontDir(usr.HomeDir)
  } else {
    p := strings.Index(config.FontDir, "~/") // TODO: Windows
    if p != -1 {
      config.FontDir = config.FontDir[:p] + usr.HomeDir + config.FontDir[p+1:]
    }
  }
  config.FontDir = filepath.Clean(config.FontDir)
  // L.Printf("config: %+v", config)

  // dispatch to command function
  cmd := flag.Arg(0)
  args := flag.Args()[1:]
  
  switch cmd {
    case "update": cmd_update(args)
    case "version": cmd_version(args)
    default:
      L.Fatalf("Unknown command %s\nSee %s -h for help\n", cmd, progname)
  }
}
