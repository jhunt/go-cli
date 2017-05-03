package cli_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"github.com/jhunt/go-cli"
)

func TestBooleans(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CLI Suite")
}

func ll(args ...string) []string {
	return args
}

var _ = Describe("CLI", func() {
	var (
		cmd      string
		leftover []string
		err      error
	)

	BeforeEach(func() {
		cmd = ""
		leftover = nil
		err = nil
	})

	Describe("Validation", func() { // {{{
		Context("With single-level command structure", func() {
			It("Complains when the same short option is used more than once", func() {
				var opt = struct {
					Help bool   `cli:"-h, --help"`
					Host string `cli:"-h, --host"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("-h.*reused.*global"))
			})

			It("Complains when the same long option is used more than once", func() {
				var opt = struct {
					TimeIt  bool   `cli:"-t, --time"`
					TheTime string `cli:"--time"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("--time.*reused.*global"))
			})

			It("Allows short options to be repeated inside a single field tag", func() {
				var opt = struct {
					Test bool   `cli:"-t, --time, -t"`
					Time string `cli:"-T"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("Allows long options to be repeated inside a single field tag", func() {
				var opt = struct {
					Test bool   `cli:"--time, -t, --time"`
					Time string `cli:"--test, --test"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("Allows short and long options to be repeated inside a single field tag", func() {
				var opt = struct {
					Test bool   `cli:"-t, --time, -t, --time"`
					Time string `cli:"--test, --test"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("With two-level command structure", func() {
			It("Complains about a sub-command reusing a global short option", func() {
				var opt = struct {
					Help bool `cli:"-h, --help"`
					Do   struct {
						Host string `cli:"-h, --host"`
					} `cli:"do"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("-h.*reused.*do.*command"))
			})

			It("Complains about a sub-command reusing a global long option", func() {
				var opt = struct {
					Help bool `cli:"-h, --help"`
					Do   struct {
						Help string `cli:"-H, --help"`
					} `cli:"do"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("--help.*reused.*do.*command"))
			})

			It("Allows sub-commands to use the same short options", func() {
				var opt = struct {
					Help bool `cli:"-h, --help"`
					Do   struct {
						Verbose bool `cli:"-v, --verbose"`
					} `cli:"do"`
					List struct {
						Verbose bool `cli:"-v"`
					} `cli:"list"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("Allows sub-commands to use the same long options", func() {
				var opt = struct {
					Help bool `cli:"-h, --help"`
					Do   struct {
						Verbose bool `cli:"-v, --verbose"`
					} `cli:"do"`
					List struct {
						Verbose bool `cli:"--verbose"`
					} `cli:"list"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("With n-level command structure", func() {
			It("Complains about a sub-sub-command reusing a global short option", func() {
				var opt = struct {
					Help  bool `cli:"-h, --help"`
					Thing struct {
						Create struct {
							Host string `cli:"-h, --host"`
						} `cli:"create"`
					} `cli:"thing"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("-h.*reused.*thing create.*command"))
			})

			It("Complains about a sub-sub-command reusing a global long option", func() {
				var opt = struct {
					Help  bool `cli:"-h, --help"`
					Thing struct {
						Create struct {
							Help bool `cli:"-H, --help"`
						} `cli:"create"`
					} `cli:"thing"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("--help.*reused.*thing create.*command"))
			})

			It("Complains about a sub-sub-command reusing a sub-command short option", func() {
				var opt = struct {
					Version string `cli:"--version"`
					Thing   struct {
						Help   bool `cli:"-h, --help"`
						Create struct {
							Host string `cli:"-h, --host"`
						} `cli:"create"`
					} `cli:"thing"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("-h.*reused.*thing create.*command"))
			})

			It("Complains about a sub-sub-command reusing a global long option", func() {
				var opt = struct {
					Version string `cli:"--version"`
					Thing   struct {
						Help   bool `cli:"-h, --help"`
						Create struct {
							Help bool `cli:"-H, --help"`
						} `cli:"create"`
					} `cli:"thing"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("--help.*reused.*thing create.*command"))
			})

			It("Allows sub-commands to use the same short options", func() {
				var opt = struct {
					Version string `cli:"--version"`
					Thing   struct {
						Create struct {
							Help bool `cli:"-h"`
						} `cli:"create"`

						List struct {
							Help bool `cli:"-h"`
						}
					} `cli:"thing"`

					Other struct {
						Help bool `cli:"-h, --help"`
					} `cli:"other"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("Allows sub-commands to use the same long options", func() {
				var opt = struct {
					Version string `cli:"--version"`
					Thing   struct {
						Create struct {
							Help bool `cli:"--help"`
						} `cli:"create"`

						List struct {
							Help bool `cli:"--help"`
						}
					} `cli:"thing"`

					Other struct {
						Help bool `cli:"--help"`
					} `cli:"other"`
				}{}

				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})

	// }}}
	Describe("Simple Boolean Flags", func() { // {{{
		Context("With go-supplied default values", func() {
			var opt = struct {
				Short bool `cli:"-s"`
				Long  bool `cli:"--long"`
				Both  bool `cli:"-b, --both"`
			}{}

			BeforeEach(func() {
				/* emulate what go initialized to */
				var flag bool
				opt.Short = flag
				opt.Long = flag
				opt.Both = flag
			})

			It("Leaves defaults intact", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Short).Should(BeFalse())
				Ω(opt.Long).Should(BeFalse())
				Ω(opt.Both).Should(BeFalse())
			})

			It("Overrides defaults for provided flags", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("-s", "--long"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Both).Should(BeFalse())

				Ω(opt.Short).Should(BeTrue())
				Ω(opt.Long).Should(BeTrue())
			})

			It("Allows repeat flags", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("-s", "--long", "-s", "-s"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Both).Should(BeFalse())

				Ω(opt.Short).Should(BeTrue())
				Ω(opt.Long).Should(BeTrue())
			})

			It("Handles flags after positional arguments", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("some", "argument", "-s", "--long"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).ShouldNot(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Both).Should(BeFalse())

				Ω(opt.Short).Should(BeTrue())
				Ω(opt.Long).Should(BeTrue())
			})

			It("Sets a mixed option field for the short flag", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("-b"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Short).Should(BeFalse())
				Ω(opt.Long).Should(BeFalse())

				Ω(opt.Both).Should(BeTrue())
			})

			It("Sets a mixed option field for the long flag", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("--both"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Short).Should(BeFalse())
				Ω(opt.Long).Should(BeFalse())

				Ω(opt.Both).Should(BeTrue())
			})

			It("Handles bundled short options", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("-sb"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Long).Should(BeFalse())

				Ω(opt.Short).Should(BeTrue())
				Ω(opt.Both).Should(BeTrue())
			})
		})

		Context("With overridden default values", func() {
			var opt = struct {
				Verify bool `cli:"-v, --verify, --no-verify"`
			}{}

			BeforeEach(func() {
				opt.Verify = true
			})

			It("Handles --no-<option> boolean variation", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--no-verify"))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(opt.Verify).Should(BeFalse())
			})
		})
	})

	// }}}
	Describe("Simple String Flags", func() { // {{{
		Context("With go-supplied default values", func() {
			var opt = struct {
				Short string `cli:"-s"`
				Long  string `cli:"--long"`
				Both  string `cli:"-b, --both"`
			}{}

			BeforeEach(func() {
				/* emulate what go initialized to */
				var flag string
				opt.Short = flag
				opt.Long = flag
				opt.Both = flag
			})

			It("Leaves defaults intact", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Short).Should(Equal(""))
				Ω(opt.Long).Should(Equal(""))
				Ω(opt.Both).Should(Equal(""))
			})

			It("Overrides defaults for provided flags", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("-s", "change", "--long", "override"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Both).Should(Equal(""))

				Ω(opt.Short).Should(Equal("change"))
				Ω(opt.Long).Should(Equal("override"))
			})

			It("Allows repeat flags", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("-s", "A", "--long", "B", "-s", "C", "-s", "D"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Both).Should(Equal(""))

				Ω(opt.Short).Should(Equal("D"))
				Ω(opt.Long).Should(Equal("B"))
			})

			It("Handles flags after positional arguments", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("some", "argument", "-s", "of course", "--long", "why not"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).ShouldNot(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Both).Should(Equal(""))

				Ω(opt.Short).Should(Equal("of course"))
				Ω(opt.Long).Should(Equal("why not"))
			})

			It("Sets a mixed option field for the short flag", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("-b", "sure"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Short).Should(Equal(""))
				Ω(opt.Long).Should(Equal(""))

				Ω(opt.Both).Should(Equal("sure"))
			})

			It("Sets a mixed option field for the long flag", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("--both", "yes"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Short).Should(Equal(""))
				Ω(opt.Long).Should(Equal(""))

				Ω(opt.Both).Should(Equal("yes"))
			})

			It("Handles bundled short options", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("-sbundled", "-byes"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Long).Should(Equal(""))

				Ω(opt.Short).Should(Equal("bundled"))
				Ω(opt.Both).Should(Equal("yes"))
			})

			It("Handles missing arguments to short options", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("-s"))
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("required .*-s"))
			})

			It("Handles missing arguments to long options", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("--long"))
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("required .*--long"))
			})
		})

		Context("With overridden default values", func() {
			var opt = struct {
				Type string `cli:"-t, --type"`
			}{}

			BeforeEach(func() {
				opt.Type = "normal"
			})

			It("Provides the default if no option is given", func() {
				_, _, err = cli.ParseArgs(&opt, ll())
				Ω(err).ShouldNot(HaveOccurred())
				Ω(opt.Type).Should(Equal("normal"))
			})

			It("Overrides the default if an option is given", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--type", "extra-special"))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(opt.Type).Should(Equal("extra-special"))
			})
		})
	})

	// }}}
	Describe("Simple Numeric Flags", func() { // {{{
		Context("With go-supplied default values", func() {
			var opt = struct {
				Int     int     `cli:"--int"`
				Int8    int8    `cli:"--int8"`
				Int16   int16   `cli:"--int16"`
				Int32   int32   `cli:"--int32"`
				Int64   int64   `cli:"--int64"`
				Uint    uint    `cli:"--uint"`
				Uint8   uint8   `cli:"--uint8"`
				Uint16  uint16  `cli:"--uint16"`
				Uint32  uint32  `cli:"--uint32"`
				Uint64  uint64  `cli:"--uint64"`
				Float32 float32 `cli:"--float32"`
				Float64 float64 `cli:"--float64"`
			}{}

			BeforeEach(func() {
				/* emulate what go initialized to */
				var (
					i   int
					i8  int8
					i16 int16
					i32 int32
					i64 int64
					u   uint
					u8  uint8
					u16 uint16
					u32 uint32
					u64 uint64
					f32 float32
					f64 float64
				)

				opt.Int = i
				opt.Int8 = i8
				opt.Int16 = i16
				opt.Int32 = i32
				opt.Int64 = i64
				opt.Uint = u
				opt.Uint8 = u8
				opt.Uint16 = u16
				opt.Uint32 = u32
				opt.Uint64 = u64
				opt.Float32 = f32
				opt.Float64 = f64
			})

			It("Leaves defaults intact", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.Int).Should(Equal(int(0)))
				Ω(opt.Int8).Should(Equal(int8(0)))
				Ω(opt.Int16).Should(Equal(int16(0)))
				Ω(opt.Int32).Should(Equal(int32(0)))
				Ω(opt.Int64).Should(Equal(int64(0)))
				Ω(opt.Uint).Should(Equal(uint(0)))
				Ω(opt.Uint8).Should(Equal(uint8(0)))
				Ω(opt.Uint16).Should(Equal(uint16(0)))
				Ω(opt.Uint32).Should(Equal(uint32(0)))
				Ω(opt.Uint64).Should(Equal(uint64(0)))
				Ω(opt.Float32).Should(Equal(float32(0.0)))
				Ω(opt.Float64).Should(Equal(float64(0.0)))
			})

			It("Handles numeric values that are within range", func() {
				var (
					i   int     = -42
					i8  int8    = 127
					i16 int16   = 32767
					i32 int32   = 2147483647
					i64 int64   = 9223372036854775807
					u   uint    = 42
					u8  uint8   = 255
					u16 uint16  = 65535
					u32 uint32  = 4294967295
					u64 uint64  = 18446744073709551615
					f32 float32 = 1.2345
					f64 float64 = 123456789.123456789123456789
				)
				cmd, leftover, err = cli.ParseArgs(&opt, ll(
					"--int", fmt.Sprintf("%v", i),
					"--int8", fmt.Sprintf("%v", i8),
					"--int16", fmt.Sprintf("%v", i16),
					"--int32", fmt.Sprintf("%v", i32),
					"--int64", fmt.Sprintf("%v", i64),
					"--uint", fmt.Sprintf("%v", u),
					"--uint8", fmt.Sprintf("%v", u8),
					"--uint16", fmt.Sprintf("%v", u16),
					"--uint32", fmt.Sprintf("%v", u32),
					"--uint64", fmt.Sprintf("%v", u64),
					"--float32", fmt.Sprintf("%v", f32),
					"--float64", fmt.Sprintf("%v", f64),
				))

				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())

				Ω(opt.Int).Should(Equal(i))
				Ω(opt.Int8).Should(Equal(i8))
				Ω(opt.Int16).Should(Equal(i16))
				Ω(opt.Int32).Should(Equal(i32))
				Ω(opt.Int64).Should(Equal(i64))

				Ω(opt.Uint).Should(Equal(u))
				Ω(opt.Uint8).Should(Equal(u8))
				Ω(opt.Uint16).Should(Equal(u16))
				Ω(opt.Uint32).Should(Equal(u32))
				Ω(opt.Uint64).Should(Equal(u64))

				Ω(opt.Float32).Should(Equal(f32))
				Ω(opt.Float64).Should(Equal(f64))
			})
		})

		Context("With non-numeric arguments", func() {
			var opt = struct {
				Int     int     `cli:"--int,-i"`
				Int8    int8    `cli:"--int8"`
				Int16   int16   `cli:"--int16"`
				Int32   int32   `cli:"--int32"`
				Int64   int64   `cli:"--int64"`
				Uint    uint    `cli:"--uint,-u"`
				Uint8   uint8   `cli:"--uint8"`
				Uint16  uint16  `cli:"--uint16"`
				Uint32  uint32  `cli:"--uint32"`
				Uint64  uint64  `cli:"--uint64"`
				Float32 float32 `cli:"--float32"`
				Float64 float64 `cli:"--float64"`
			}{}

			It("Complains about non-numeric int value arguments", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--int", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric int value arguments (short)", func() {
				_, _, err = cli.ParseArgs(&opt, ll("-i", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric int value arguments (short+bundled)", func() {
				_, _, err = cli.ParseArgs(&opt, ll("-iBAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric int8 value arguments", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--int8", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric int16 value arguments", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--int16", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric int32 value arguments", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--int32", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric int64 value arguments", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--int64", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric uint value arguments", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--uint", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric uint value arguments (short)", func() {
				_, _, err = cli.ParseArgs(&opt, ll("-u", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric uint value arguments (short+bundled)", func() {
				_, _, err = cli.ParseArgs(&opt, ll("-uBAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric uint8 value arguments", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--uint8", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric uint16 value arguments", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--uint16", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric int32 value arguments", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--int32", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric uint64 value arguments", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--uint64", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric float32 value arguments", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--float32", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})

			It("Complains about non-numeric float64 value arguments", func() {
				_, _, err = cli.ParseArgs(&opt, ll("--float64", "BAD"))
				Ω(opt.Int).Should(Equal(0))
				Ω(err).Should(HaveOccurred())
			})
		})
	})

	// }}}
	Describe("Simple List Flags", func() { // {{{
		Context("With go-supplied default values", func() {
			var opt = struct {
				List []string `cli:"--list"`
			}{}

			BeforeEach(func() {
				/* we have to initialize to make reflect happy... */
				opt.List = make([]string, 0)
			})

			It("Leaves defaults intact", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll())
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.List).Should(BeEmpty())
			})

			It("Handles a single argument", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("--list", "a"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.List).ShouldNot(BeEmpty())
				Ω(len(opt.List)).Should(Equal(1))
				Ω(opt.List[0]).Should(Equal("a"))
			})

			It("Handles multiple arguments", func() {
				cmd, leftover, err = cli.ParseArgs(&opt, ll("--list", "a", "--list", "b", "--list", "c"))
				Ω(cmd).Should(Equal(""))
				Ω(leftover).Should(BeEmpty())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.List).ShouldNot(BeEmpty())
				Ω(len(opt.List)).Should(Equal(3))
				Ω(opt.List[0]).Should(Equal("a"))
				Ω(opt.List[1]).Should(Equal("b"))
				Ω(opt.List[2]).Should(Equal("c"))
			})
		})

		Context("With overridden default values", func() {
			var opt = struct {
				List []int `cli:"-N, --numbers"`
			} {}

			BeforeEach(func() {
				opt.List = []int{4, 8, 15, 16}
			})

			It("Uses the default list if no options are given", func() {
				_, _, err = cli.ParseArgs(&opt, ll())
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.List).ShouldNot(BeEmpty())
				Ω(len(opt.List)).Should(Equal(4))
				Ω(opt.List[0]).Should(Equal(4))
				Ω(opt.List[1]).Should(Equal(8))
				Ω(opt.List[2]).Should(Equal(15))
				Ω(opt.List[3]).Should(Equal(16))
			})

			It("Replaces the default list if options are given", func() {
				_, _, err = cli.ParseArgs(&opt, ll("-N", "23", "-N", "42"))
				Ω(err).ShouldNot(HaveOccurred())

				Ω(opt.List).ShouldNot(BeEmpty())
				Ω(len(opt.List)).Should(Equal(2))
				Ω(opt.List[0]).Should(Equal(23))
				Ω(opt.List[1]).Should(Equal(42))
			})
		})
	})

	// }}}
	Describe("Sub-commands", func() { // {{{
		var opt = struct {
			Help    bool   `cli:"-h, -?, --help"`
			Version bool   `cli:"-v, --version"`
			Target  string `cli:"-t, --target"`

			Hosts struct {
				Raw  bool   `cli:"-R, --raw"`
				Host string `cli:"-H, --host"`
			} `cli:"hosts"`

			Users struct {
				List struct {
					All bool `cli"-a, --all"`
				} `cli:"list"`
			} `cli:"users"`
		}{}

		BeforeEach(func() {
			var (
				dBool   bool
				dString string
			)

			opt.Help = dBool
			opt.Version = dBool
			opt.Target = dString

			opt.Hosts.Raw = dBool
			opt.Hosts.Host = dString

			opt.Users.List.All = dBool
		})

		It("Just works", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("-v", "--target", "prod", "hosts", "--raw", "active"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cmd).Should(Equal("hosts"))

			Ω(leftover).ShouldNot(BeEmpty())
			Ω(len(leftover)).Should(Equal(1))
			Ω(leftover[0]).Should(Equal("active"))

			Ω(opt.Help).Should(BeFalse())
			Ω(opt.Version).Should(BeTrue())
			Ω(opt.Target).Should(Equal("prod"))
			Ω(opt.Hosts.Raw).Should(BeTrue())
			Ω(opt.Hosts.Host).Should(Equal(""))
			Ω(opt.Users.List.All).Should(BeFalse())
		})

		It("Allows global options to be specified after sub-commands", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("-v", "hosts", "--raw", "--target", "qa", "active"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cmd).Should(Equal("hosts"))

			Ω(leftover).ShouldNot(BeEmpty())
			Ω(len(leftover)).Should(Equal(1))
			Ω(leftover[0]).Should(Equal("active"))

			Ω(opt.Help).Should(BeFalse())
			Ω(opt.Version).Should(BeTrue())
			Ω(opt.Target).Should(Equal("qa"))
			Ω(opt.Hosts.Raw).Should(BeTrue())
			Ω(opt.Hosts.Host).Should(Equal(""))
			Ω(opt.Users.List.All).Should(BeFalse())
		})

		It("Allows global options to be specified after positional arguments", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("-v", "hosts", "--raw", "active", "--target", "dev1"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cmd).Should(Equal("hosts"))

			Ω(leftover).ShouldNot(BeEmpty())
			Ω(len(leftover)).Should(Equal(1))
			Ω(leftover[0]).Should(Equal("active"))

			Ω(opt.Help).Should(BeFalse())
			Ω(opt.Version).Should(BeTrue())
			Ω(opt.Target).Should(Equal("dev1"))
			Ω(opt.Hosts.Raw).Should(BeTrue())
			Ω(opt.Hosts.Host).Should(Equal(""))
			Ω(opt.Users.List.All).Should(BeFalse())
		})

		It("Stops processing sub-commands after an empty argument", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("-v", "--target", "prod", "", "hosts", "--help"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cmd).Should(Equal(""))

			Ω(leftover).ShouldNot(BeEmpty())
			Ω(len(leftover)).Should(Equal(2))
			Ω(leftover[0]).Should(Equal(""))
			Ω(leftover[1]).Should(Equal("hosts"))

			Ω(opt.Help).Should(BeTrue())
			Ω(opt.Version).Should(BeTrue())
			Ω(opt.Target).Should(Equal("prod"))
			Ω(opt.Hosts.Raw).Should(BeFalse())
			Ω(opt.Hosts.Host).Should(Equal(""))
			Ω(opt.Users.List.All).Should(BeFalse())
		})

		It("Stops processing all arguments after the first `--`", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("-v", "--target", "dev", "--", "hosts", "--help"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cmd).Should(Equal(""))

			Ω(leftover).ShouldNot(BeEmpty())
			Ω(len(leftover)).Should(Equal(2))
			Ω(leftover[0]).Should(Equal("hosts"))
			Ω(leftover[1]).Should(Equal("--help"))

			Ω(opt.Help).Should(BeFalse())
			Ω(opt.Version).Should(BeTrue())
			Ω(opt.Target).Should(Equal("dev"))
			Ω(opt.Hosts.Raw).Should(BeFalse())
			Ω(opt.Hosts.Host).Should(Equal(""))
			Ω(opt.Users.List.All).Should(BeFalse())
		})
	})

	// }}}
	Describe("Sub-command aliasing", func() { // {{{
		var opt = struct {
			Help bool `cli:"-h, -?, --help"`

			List struct {
				All bool `cli:"--all"`
			} `cli:"list, ls"`
		}{}

		BeforeEach(func() {
			var (
				dBool bool
			)

			opt.Help = dBool
			opt.List.All = dBool
		})

		It("Handles the first alias", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("list", "--all"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cmd).Should(Equal("list"))

			Ω(leftover).Should(BeEmpty())

			Ω(opt.Help).Should(BeFalse())
			Ω(opt.List.All).Should(BeTrue())
		})

		It("Treats the first alias as the canonical alias", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("ls", "--all"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(cmd).Should(Equal("list"))

			Ω(leftover).Should(BeEmpty())

			Ω(opt.Help).Should(BeFalse())
			Ω(opt.List.All).Should(BeTrue())
		})
	})

	// }}}
	Describe("Unrecognized arguments", func() { // {{{
		var opt = struct {
			Help bool `cli:"-h, -?, --help"`

			Sub struct {
				Host string `cli:"-H, --host"`
			} `cli:"sub"`
		}{}

		It("Complains about unrecognized global short options", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("-h", "-x"))
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(MatchRegexp("unrecognized.*-x"))
		})

		It("Complains about unrecognized global long options", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("-h", "--exclude", "whatever"))
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(MatchRegexp("unrecognized.*--exclude"))
		})

		It("Complains about unrecognized sub-command short options", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("-h", "sub", "-H", "test", "-x"))
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(MatchRegexp("unrecognized.*-x"))
		})

		It("Complains about unrecognized sub-command long options", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("-h", "sub", "-H", "test", "--exclude", "whatever"))
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(MatchRegexp("unrecognized.*--exclude"))
		})

		It("Does not recognized sub-command options at the global level", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("-H", "test", "sub"))
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(MatchRegexp("unrecognized.*-H"))
		})
	})

	// }}}
	Describe("Invalid object types", func() { // {{{
		Context("With non-struct types", func() {
			var (
				i   int
				i8  int8
				i16 int16
				i32 int32
				i64 int64

				u   uint
				u8  uint8
				u16 uint16
				u32 uint32
				u64 uint64

				f32 float32
				f64 float64

				tf  bool
				tfl []bool

				s string
			)

			It("Fails to work with an int", func() {
				_, _, err = cli.ParseArgs(&i, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with an int8", func() {
				_, _, err = cli.ParseArgs(&i8, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with an int16", func() {
				_, _, err = cli.ParseArgs(&i16, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with an int32", func() {
				_, _, err = cli.ParseArgs(&i32, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with an int64", func() {
				_, _, err = cli.ParseArgs(&i64, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with a uint", func() {
				_, _, err = cli.ParseArgs(&u, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with a uint8", func() {
				_, _, err = cli.ParseArgs(&u8, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with a uint16", func() {
				_, _, err = cli.ParseArgs(&u16, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with a uint32", func() {
				_, _, err = cli.ParseArgs(&u32, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with a uint64", func() {
				_, _, err = cli.ParseArgs(&u64, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with a float32", func() {
				_, _, err = cli.ParseArgs(&f32, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with a float64", func() {
				_, _, err = cli.ParseArgs(&f64, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with a bool", func() {
				_, _, err = cli.ParseArgs(&tf, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with a list of bool", func() {
				_, _, err = cli.ParseArgs(&tfl, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})

			It("Fails to work with a string", func() {
				_, _, err = cli.ParseArgs(&s, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("only operates on structures"))
			})
		})

		Context("With non-scalar types embedded in a struct", func() {
			It("Fails to work with an interface field", func() {
				var opt = struct {
					Bad interface{} `cli:"--interface"`
				}{}
				_, _, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("cannot operate on this type"))
			})
		})
	})

	// }}}
	Describe("Malformed tags", func() { // {{{
		Context("In the global context", func() {
			It("Complains about oddball short options", func() {
				var opt = struct {
					Bad bool `cli:"-%"`
				}{}
				_, _, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("invalid.*-%"))
			})

			It("Complains about oddball long options", func() {
				var opt = struct {
					Bad bool `cli:"--earned-%"`
				}{}
				_, _, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("invalid.*--earned-%"))
			})

			It("Complains about long options that are too short", func() {
				var opt = struct {
					Bad bool `cli:"--y"`
				}{}
				_, _, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("invalid.*--y"))
			})

			It("Complains about long options that start with a hyphen", func() {
				var opt = struct {
					Bad bool `cli:"---dashes-yo"`
				}{}
				_, _, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("invalid.*---dashes-yo"))
			})
		})

		Context("In a sub-command context", func() {
			It("Complains about oddball short options", func() {
				var opt = struct {
					Sub struct {
						Bad bool `cli:"-%"`
					} `cli:"sub"`
				}{}
				_, _, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("invalid.*-%"))
			})

			It("Complains about oddball long options", func() {
				var opt = struct {
					Sub struct {
						Bad bool `cli:"--earned-%"`
					} `cli:"sub"`
				}{}
				_, _, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("invalid.*--earned-%"))
			})

			It("Complains about long options that are too short", func() {
				var opt = struct {
					Sub struct {
						Bad bool `cli:"--y"`
					} `cli:"sub"`
				}{}
				_, _, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("invalid.*--y"))
			})

			It("Complains about long options that start with a hyphen", func() {
				var opt = struct {
					Sub struct {
						Bad bool `cli:"---dashes-yo"`
					} `cli:"sub"`
				}{}
				_, _, err = cli.ParseArgs(&opt, ll())
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(MatchRegexp("invalid.*---dashes-yo"))
			})
		})
	})

	// }}}
	Describe("Chained Commands", func() { // {{{
		var opt = struct {
			Help     bool   `cli:"-h, -?, --help"`
			Insecure bool   `cli:"-k, --insecure"`
			Target   string `cli:"-t, --target"`

			Sub struct {
				Host string `cli:"-H, --host"`
			} `cli:"sub"`

			List struct {
			} `cli:"list"`
		}{}

		BeforeEach(func() {
			var (
				b bool
				s string
			)
			opt.Help = b
			opt.Insecure = b
			opt.Target = s
			opt.Sub.Host = s
		})

		It("Works with no arguments", func() {
			p, err := cli.NewParser(&opt, ll())
			Ω(err).ShouldNot(HaveOccurred())

			Ω(p.Next()).Should(BeFalse())
		})

		It("Works with just global options", func() {
			p, err := cli.NewParser(&opt, ll("--insecure", "-t", "my-target"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(p.Next()).Should(BeFalse())
			Ω(opt.Help).Should(BeFalse())
			Ω(opt.Insecure).Should(BeTrue())
			Ω(opt.Target).Should(Equal("my-target"))
			Ω(opt.Sub.Host).Should(Equal(""))
		})

		It("Works with a single sub-command", func() {
			p, err := cli.NewParser(&opt, ll("--insecure", "-t", "my-target", "sub"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(opt.Help).Should(BeFalse())
			Ω(opt.Insecure).Should(BeTrue())
			Ω(opt.Target).Should(Equal("my-target"))

			Ω(p.Next()).Should(BeTrue())
			Ω(p.Command).Should(Equal("sub"))
			Ω(opt.Sub.Host).Should(Equal(""))

			Ω(p.Next()).Should(BeFalse())
		})

		It("Works with two sub-commands", func() {
			p, err := cli.NewParser(&opt, ll("--insecure", "-t", "my-target", "list", "--", "sub"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(opt.Help).Should(BeFalse())
			Ω(opt.Insecure).Should(BeTrue())
			Ω(opt.Target).Should(Equal("my-target"))

			Ω(p.Next()).Should(BeTrue())
			Ω(p.Command).Should(Equal("list"))

			Ω(p.Next()).Should(BeTrue())
			Ω(p.Command).Should(Equal("sub"))
			Ω(opt.Sub.Host).Should(Equal(""))

			Ω(p.Next()).Should(BeFalse())
		})

		It("Keeps sub-command arguments separate", func() {
			p, err := cli.NewParser(&opt, ll("--insecure", "-t", "my-target", "sub", "--host", "prod", "--", "sub", "-H", "dev"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(opt.Help).Should(BeFalse())
			Ω(opt.Insecure).Should(BeTrue())
			Ω(opt.Target).Should(Equal("my-target"))

			Ω(p.Next()).Should(BeTrue())
			Ω(p.Command).Should(Equal("sub"))
			Ω(opt.Sub.Host).Should(Equal("prod"))

			Ω(p.Next()).Should(BeTrue())
			Ω(p.Command).Should(Equal("sub"))
			Ω(opt.Sub.Host).Should(Equal("dev"))

			Ω(p.Next()).Should(BeFalse())
		})

		It("Allows global option overrides on a per sub-command basis", func() {
			p, err := cli.NewParser(&opt, ll("-k", "-t", "x", "sub", "--target", "my-target", "--", "sub", "-k"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(opt.Help).Should(BeFalse())
			Ω(opt.Insecure).Should(BeTrue())
			Ω(opt.Target).Should(Equal("x"))

			Ω(p.Next()).Should(BeTrue())
			Ω(p.Command).Should(Equal("sub"))
			Ω(opt.Target).Should(Equal("my-target"))
			Ω(opt.Sub.Host).Should(Equal(""))

			Ω(p.Next()).Should(BeTrue())
			Ω(p.Command).Should(Equal("sub"))
			Ω(opt.Target).Should(Equal("x"))
			Ω(opt.Sub.Host).Should(Equal(""))

			Ω(p.Next()).Should(BeFalse())
		})
	})

	// }}}
	Describe("Unchained -- behavior", func() { // {{{
		var opt = struct {
			Help  bool `cli:"-h, -?, --help"`
			Debug bool `cli:"-D, --debug"`

			Merge struct {
				Prune string `cli:"-d, --prune"`
			} `cli:"merge"`
		}{}

		It("Concatenates arguments after -- with positionals before it", func() {
			cmd, leftover, err = cli.ParseArgs(&opt, ll("merge", "a.yml", "-D", "--prune", "x.y.z", "--", "--dash.yml"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cmd).Should(Equal("merge"))

			Ω(opt.Debug).Should(BeTrue())
			Ω(opt.Merge.Prune).Should(Equal("x.y.z"))

			Ω(len(leftover)).Should(Equal(2))
			Ω(leftover[0]).Should(Equal("a.yml"))
			Ω(leftover[1]).Should(Equal("--dash.yml"))
		})
	})

	// }}}
	Describe("Full-stop behavior", func() { // {{{
		var opt = struct {
			Help  bool `cli:"-h, -?, --help"`
			Debug bool `cli:"-D, --debug"`

			FullStop struct {
			} `cli:"stop!"`
		}{}

		It("Treats non-options after 'stop' as arguments", func() {
			cmd, leftover, err := cli.ParseArgs(&opt, ll("-D", "stop", "all"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cmd).Should(Equal("stop"))

			Ω(opt.Debug).Should(BeTrue())
			Ω(len(leftover)).Should(Equal(1))
			Ω(leftover[0]).Should(Equal("all"))
		})

		It("Treats options after 'stop' as arguments", func() {
			cmd, leftover, err := cli.ParseArgs(&opt, ll("-D", "stop", "--hard", "all"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cmd).Should(Equal("stop"))

			Ω(opt.Debug).Should(BeTrue())
			Ω(len(leftover)).Should(Equal(2))
			Ω(leftover[0]).Should(Equal("--hard"))
			Ω(leftover[1]).Should(Equal("all"))
		})
	})

	// }}}
})
