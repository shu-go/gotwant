package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/shu-go/gli/v2"
)

type globalCmd struct {
}

type state uint8

const (
	searchingGot state = iota
	readingGot
	readingWant
)

func (c globalCmd) Run() error {
	buf := &bytes.Buffer{}

	r := bufio.NewReader(os.Stdin)

	var got, want string

	s := searchingGot
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}

		if strings.HasPrefix(line, "        got:  ") {
			s = readingGot
			got = line[14:]
			want = ""
			continue
		}
		if strings.HasPrefix(line, "        want: ") {
			s = readingWant
			want = line[14:]
			continue
		}
		if strings.HasPrefix(line, "        ") {
			if s == readingGot {
				got += line
				continue
			} else if s == readingWant {
				want += line
				continue
			}
		}

		// colorize
		if s == readingWant {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(got, want, false)

			buf.WriteString("        got:  ")
			buf.WriteString(got)
			//minus := color.New(color.FgGreen, color.Bold)
			//plus := color.New(color.FgRed, color.CrossedOut, color.Bold)
			/*
				for _, d := range diffs {
					switch d.Type {
					case diffmatchpatch.DiffEqual:
						buf.WriteString(d.Text)
					case diffmatchpatch.DiffDelete:
						plus.Fprint(buf, d.Text)
					case diffmatchpatch.DiffInsert:
						minus.Fprint(buf, d.Text)
					default:
					}
				}
			*/

			buf.WriteString("        want: ")
			//buf.WriteString(want)
			buf.WriteString(dmp.DiffPrettyText(diffs))

			// output
		}
		s = searchingGot

		buf.WriteString(line)
	}

	io.Copy(os.Stdout, buf)

	return nil
}

// Version is app version
var Version string

func main() {
	app := gli.NewWith(&globalCmd{})
	app.Name = "gotwant"
	app.Desc = "colorise"
	app.Version = Version
	app.Usage = ``
	app.Copyright = "(C) 2024 Shuhei Kubota"
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
