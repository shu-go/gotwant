package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/shu-go/gli/v2"
)

type globalCmd struct {
	Monochrome bool `cli:"m,mono,monochrome" default:"false"`

	Debug bool

	debug func(string, ...interface{})
}

func Output(format string, a ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, a...)
}

type state uint8

const (
	searchingGot state = iota
	readingGot
	readingWant
)

func (c *globalCmd) Before() {
	if c.Debug {
		c.debug = Output
	} else {
		c.debug = func(string, ...interface{}) {}
	}
}

func (c globalCmd) Run() error {
	buf := &bytes.Buffer{}

	r := bufio.NewReader(os.Stdin)

	var got, want string

	gwRE := regexp.MustCompile(`^(\s*)(got:|want:)(\s*)`)

	s := searchingGot
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}

		c.debug("*****")
		c.debug("line=%q", line)

		matches := gwRE.FindStringSubmatch(line)
		c.debug("matches=%#v", matches)
		if len(matches) != 0 {
			if strings.HasPrefix(matches[2], "got") {
				c.debug("GOT")
				s = readingGot
				got = line[len(matches[0]):]
				want = ""
				continue
			}
			if strings.HasPrefix(matches[2], "want") {
				c.debug("WANT")
				s = readingWant
				want = line[len(matches[0]):]
				continue
			}
		}

		c.debug("mode=%d", s)

		trimline := strings.TrimLeft(line, " \t")
		if !strings.HasPrefix(trimline, "FAIL") && !strings.HasPrefix(trimline, "---") {
			if s == readingGot {
				if got != "" {
					got += strings.Repeat(" ", 8+len("got:  "))
				}
				got += trimline
				continue
			} else if s == readingWant {
				if want != "" {
					want += strings.Repeat(" ", 8+len("want: "))
				}
				want += trimline
				continue
			}
		}

		c.debug("got=%q", got)
		c.debug("want=%q", want)

		// colorize
		if s == readingWant {
			c.debug("OUTPUT")
			dmp := diffmatchpatch.New()
			diffs := splitDiff(dmp.DiffMain(got, want, false))
			for i, d := range diffs {
				c.debug("%d diff=%+v", i, d)
			}

			buf.WriteString("        got:  ")
			buf.WriteString(got)

			buf.WriteString("        want: ")
			if c.Monochrome {
				buf.WriteString(want)
			} else {
				//buf.WriteString(dmp.DiffPrettyText(diffs))
				minus := color.New(color.FgGreen, color.Bold)
				minusS := color.New(color.FgGreen, color.Underline, color.Bold)
				plus := color.New(color.FgRed, color.Bold)
				plusS := color.New(color.FgRed, color.Underline, color.Bold)

				for _, d := range diffs {
					switch d.Type {
					case diffmatchpatch.DiffEqual:
						buf.WriteString(d.Text)
					case diffmatchpatch.DiffDelete:
						if d.isSpace {
							plusS.Fprint(buf, d.Text)
						} else {
							plus.Fprint(buf, d.Text)
						}
					case diffmatchpatch.DiffInsert:
						if d.isSpace {
							minusS.Fprint(buf, d.Text)
						} else {
							minus.Fprint(buf, d.Text)
						}
					default:
					}
				}

			}
		}
		s = searchingGot

		buf.WriteString(line)
	}

	io.Copy(os.Stdout, buf)

	return nil
}

type diff struct {
	diffmatchpatch.Diff
	isSpace bool
}

func splitDiff(diffs []diffmatchpatch.Diff) []diff {
	results := make([]diff, 0, len(diffs))

	for _, d := range diffs {
		if len(d.Text) == 0 {
			continue
		}

		var prevSpace bool
		var s []rune
		for i, r := range d.Text {
			space := !unicode.IsPrint(r) || unicode.IsSpace(r)
			if i == 0 {
				prevSpace = space
			}

			if space != prevSpace {
				newD := d
				newD.Text = string(s)
				results = append(results, diff{
					Diff:    newD,
					isSpace: prevSpace,
				})
				s = s[:0]
			}
			s = append(s, r)
			prevSpace = space
		}
		if len(s) > 0 {
			newD := d
			newD.Text = string(s)
			results = append(results, diff{
				Diff:    newD,
				isSpace: prevSpace,
			})
		}
	}

	return results
}

// Version is app version
var Version string

func main() {
	app := gli.NewWith(&globalCmd{})
	app.AutoNoBoolOptions = false
	app.Name = "gotwant"
	app.Desc = "colorise and align got-want style test results"
	app.Version = Version
	app.Usage = ``
	app.Copyright = "(C) 2024 Shuhei Kubota"
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
