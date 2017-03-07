package cli

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

/* Parse looks through os.Args, and returns the sub-command name
   (or "" for none), the remaining positional arguments, and any
   error that was encountered. */
func Parse(thing interface{}) (string, []string, error) {
	return ParseArgs(thing, os.Args[1:])
}

/* ParseArgs is like Parse(), except that it operates on an explicit
   list of arguments, instead of implicitly using os.Args. */
func ParseArgs(thing interface{}, args []string) (string, []string, error) {
	remnants := make([]string, 0)

	c, err := reflectOnIt(thing)
	if err != nil {
		return "", remnants, err
	}

	/* make sure we didn't do anything semantically invalid... */
	if err := validate(c); err != nil {
		return "", remnants, err
	}

	/* process the arguments */
	stack := make([]string, 0) /* command stack */
	insub := true              /* we are still parsing sub-commands */
	level := c
	for len(args) > 0 {
		arg := args[0]
		args = args[1:]

		if len(arg) == 0 {
			/* an empty argument, while strange, is not illegal
			   it means that we are done with sub-command processing,
			   and that all future non-options are remnants */
			insub = false
			remnants = append(remnants, arg)
			continue
		}

		if arg == "--" {
			/* immediately stop any and all argument processing. */
			remnants = append(remnants, args...)
			break
		}

		if arg[0] != '-' || len(arg) == 1 { /* handle '-' for stdin... */
			if insub {
				/* possible sub-command.  see if one was defined,
				   and either descend to its level, or flip insub
				   and start treating all non-options as remnants */
				if next, ok := level.Subs[arg]; ok {
					stack = append(stack, arg)
					level = next
					continue
				}
				insub = false
			}
			remnants = append(remnants, arg)
			continue
		}

		if arg[1] == '-' { /* long option! */
			name := arg[2:]
			opt, err := c.findLong(stack, name)
			if err != nil {
				return "", remnants, err
			}

			/* now we need to determine if we have a value arg or not.
			   `cli` uses a simple heuristic that works well in practice:

			     - bool receivers do not take value args
			     - everything else takes a value arg
			*/
			if opt.Kind == reflect.Bool {
				opt.enable()

			} else {
				if len(args) == 0 {
					return "", remnants, fmt.Errorf("missing required value for `%s` flag", arg)
				}
				err = opt.set(args[0])
				if err != nil {
					return "", remnants, err
				}
				args = args[1:]
			}

		} else { /* short option(s)! */
			arg = arg[1:]
			for len(arg) > 0 {
				name := arg[0:1]
				arg = arg[1:]

				opt, err := c.findShort(stack, name)
				if err != nil {
					return "", remnants, err
				}
				if opt.Kind == reflect.Bool {
					opt.enable()

				} else {
					/* attempt to use the rest of the short block, if there is one,
					   as the value arg... */
					if len(arg) > 0 {
						err = opt.set(arg)
						if err != nil {
							return "", remnants, err
						}
						break
					}
					/* otherwise, we need the next argument in the arg list... */
					if len(args) == 0 {
						return "", remnants, fmt.Errorf("missing required value for `-%s` flag", name)
					}
					err = opt.set(args[0])
					if err != nil {
						return "", remnants, err
					}
					args = args[1:]
					break
				}
			}
		}
	}

	return strings.Join(stack, " "), remnants, nil
}
