package main

import (
  "fmt"
  "os"

  "github.com/jhunt/go-cli"
)

type Options struct {
  Help     bool   `cli:"-h, --help"`
  Version  bool   `cli:"-v, --version"`
  Insecure bool   `cli:"-k, --insecure"`
  URL      string `cli:"-U, --url"`

  Gen struct {
    Length int     `cli:"-l, --length"`
    Policy string  `cli:"-p, --policy"`
  } `cli:"gen"`
}

func main() {
  var options Options
  options.Gen.Length = 48 // a default

  command, args, err := cli.Parse(&options)
  if err != nil {
    fmt.Fprintf(os.Stderr, "!!! %s\n", err)
    os.Exit(1)
  }

  if command == "gen" {
    path := ""
    if len(args) > 0 {
      path = args[0]
    }
    fmt.Printf("generating a password %d characters long", options.Gen.Length)
    fmt.Printf("in the vault at %s\n", path)
    // ...
  }
}
