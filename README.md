go-cli
======

Rite-of-passage-style Go command line parser.

CLIs are awesome.  Most options libraries aren't.  `go-cli`
doesn't attempt to change that, it just tries to focus on doing
one thing well.

Things you will not find in `go-cli`:

  - Magical bash/zsh Auto-completion support
  - Usage generation
  - Option help string
  - Option defaults

Things you will find in `go0cli`:

  - A dead-simple, tagged-struct approach to options
  - A rudimentary sub-command recognizer
  - A flexible argument processor

Usage
=====

Should be pretty simple.

```
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
    fmt.Printf("generating a password %d characters long", options.Gen.Length)
    // ...
  }
}
```

Reusing Flags
=============

You can reuse option flags, both short and long, as long as it is
provable unambiguous where and when callers can use the flag.
Practically, this means:

  1. You cannot reuse flags defined "above" you
  2. You cannot reuse flags on the same level as you

This allows `go-cli` to recognize arguments for a single level
(global, sub-command, sub-sub-command, ad infinitum) at any point
after the "beginning" of that level.

Let's look at some examples, shall we?

```
type Options struct {
  Help   bool           `cli:"-h, --help"`

  List struct {
    LongForm   bool     `cli:"-l, --long"`
    All        bool     `cli:"-a, --all"`
  } `cli:"list"`

  Create struct {
    Archive    string   `cli:"-a, "--archive"`
    Name       string   `cli:"-n, "--name"`
  } `cli:"new"`
}
```

Here, `-h` / `--help` is _global option_.  It can appear anywhere
in the command line invocation, and it has the same semantics
everywhere (namely, to show the help or something).

On the contrary, `-l` / `--long` only makes sense after the `list`
sub-command.  If encountered before `list`, it's an unrecognized
flag.

The up-shot of this is that a user of your CLI can do this:

```
$ ./foo -h
$ ./foo list -h
$ ./foo list -h --all
$ ./foo -h list -h --all -h
```

This is why you can't override the `-h` / `--help` flag on a
per-command basis -- it's just too confusing to end users
(including the author of `go-cli`).

If you look closely, you'll notice that both `list` and `new`
define a `-a` short option.  What gives?  Didn't this guy _just
get done saying that you can't override flags??_

It's cool.  It's going to be alright.  There's not much chance of
a user conflating the two `-a` use cases - `list -a` lists
everything, but `new -a name` sets an archive name.  And since
`-a` doesn't exist at the global level ("above"), you can't do
this:

```
$ ./foo -a list                 # this is bad
```

So, without any ambiguity, `go-cli` is perfectly happy to let you
overload the meaning of `-a`.  Whether you _should_, is entirely
up to you.

Contributing
============

This code is licensed MIT.  Enjoy.

If you find a bug, please raise a [Github Issue][issues] first,
before submitting a PR.

When you do work up a patch, keep in mind that we have a fairly
extensive test suite in `cli_test.go`.  I don't care _all that
much_ about code coverage, but we do have >90% C0 code coverage on
the current tests suite, and I'd like to keep it that way.

(That's not to say we've caught 90% of the bugs, but it's better
than nothin')

Happy Hacking!

