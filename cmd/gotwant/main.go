package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/shu-go/gli/v2"
)

type globalCmd struct {
	Monochrome bool `cli:"m,mono,monochrome" default:"false"`

	Efficiency       bool `cli:"e,efficiency" help:"reduces the number of edits by eliminating operationally trivial equalities"`
	Merge            bool `cli:"merge" help:"Any edit section can move as long as it doesn't cross an equality"`
	Semantic         bool `cli:"s,semantic" help:"reduces the number of edits by eliminating semantically trivial equalities"`
	SemanticLossless bool `cli:"sl,semantic-lossless" help:"looks for single edits surrounded on both sides by equalities which can be shifted sideways to align the edit to a word boundary"`

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
	gwIndent := 0

	gwRE := regexp.MustCompile(`^(\s*)(got:|want:)\s( *)`)

	s := searchingGot
	for {
		line, err := r.ReadString('\n')
		if err != nil && line == "" {
			break
		}
		indent := countIndent(line)

		c.debug("*****")
		c.debug("line=%q", line)

		matches := gwRE.FindStringSubmatch(line)
		c.debug("matches=%#v", matches)
		if len(matches) != 0 {
			if strings.HasPrefix(matches[2], "got") {
				c.debug("GOT")
				s = readingGot
				got = strings.TrimRight(line[len(matches[0]):], "\n")
				want = ""
				gwIndent = len(matches[1])
				c.debug("gwIndent=%d", gwIndent)
				continue
			}
			if strings.HasPrefix(matches[2], "want") {
				c.debug("WANT")
				s = readingWant
				want = strings.TrimRight(line[len(matches[0]):], "\n")
				gwIndent = len(matches[1])
				continue
			}
		}

		c.debug("mode=%d", s)

		trimline := strings.TrimRight(line, "\n")
		if gwIndent == 0 {
			// output from Example
			// nop
		}
		outputIndentStr := strings.Repeat(" ", gwIndent+6)
		if strings.HasPrefix(trimline, outputIndentStr) {
			trimline = trimline[len(outputIndentStr):]
		}
		if !strings.HasPrefix(trimline, "FAIL") && !strings.HasPrefix(trimline, "---") && gwIndent <= indent {
			if s == readingGot {
				if got != "" {
					got += "\n"
				}
				got += trimline
				continue
			} else if s == readingWant {
				if want != "" {
					want += "\n"
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
			dmpdiffs := dmp.DiffMain(got, want, true)
			c.debug("BEFORE")
			for i, d := range dmpdiffs {
				c.debug("%d dmpdiff=%v:%q", i, d.Type, d.Text)
			}
			switch true {
			case c.Efficiency:
				dmpdiffs = dmp.DiffCleanupEfficiency(dmpdiffs)
			case c.Merge:
				dmpdiffs = dmp.DiffCleanupMerge(dmpdiffs)
			case c.Semantic:
				dmpdiffs = dmp.DiffCleanupSemantic(dmpdiffs)
			case c.SemanticLossless:
				dmpdiffs = dmp.DiffCleanupSemanticLossless(dmpdiffs)
			default:
			}
			if c.Efficiency || c.Merge || c.Semantic || c.SemanticLossless {
				c.debug("AFTER CLEANUP")
				for i, d := range dmpdiffs {
					c.debug("%d dmpdiff=%v:%q", i, d.Type, d.Text)
				}
			}
			dmpdiffs = splitByNewline(dmpdiffs)
			dmpdiffs = addIndents(dmpdiffs, outputIndentStr)
			c.debug("AFTER INDENTATION")
			for i, d := range dmpdiffs {
				c.debug("%d dmpdiff=%v:%q", i, d.Type, d.Text)
			}
			diffs := splitDiff(suppressPrefixUnderline(dmpdiffs))
			c.debug("AFTER SUPPRESSION")
			for i, d := range diffs {
				c.debug("%d diff=%v:%q (%v)", i, d.Type, d.Text, d.isSpace)
			}

			buf.WriteString(strings.Repeat(" ", gwIndent))
			buf.WriteString("got:  ")
			buf.WriteString(strings.ReplaceAll(got, "\n", "\n"+outputIndentStr))
			if !strings.HasSuffix(got, "\n") {
				buf.WriteByte('\n')
			}

			buf.WriteString(strings.Repeat(" ", gwIndent))
			buf.WriteString("want: ")
			if c.Monochrome {
				buf.WriteString(strings.ReplaceAll(want, "\n", "\n"+outputIndentStr))
				buf.WriteByte('\n')
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
				buf.WriteByte('\n')

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

func suppressPrefixUnderline(diffs []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	result := make([]diffmatchpatch.Diff, 0, len(diffs)*2)

	result = slices.Insert(result, 0, diffs...)

	suppressLen := 14
	linePrefix := strings.Repeat(" ", suppressLen)

	index := len(result) - 1
	for index >= 0 {
		if result[index].Type == diffmatchpatch.DiffEqual {
			index--
			continue
		}

		insIndex := index
		for {
			//fmt.Fprintf(os.Stderr, " [%d] %s %q\n", insIndex, result[insIndex].Type, result[insIndex].Text)
			pos := strings.Index(result[insIndex].Text, linePrefix)
			if pos == -1 || result[insIndex].Text == linePrefix {
				break
			}
			//fmt.Fprintf(os.Stderr, " %d\n", pos)

			d := result[insIndex]

			newD := d
			newD.Type = diffmatchpatch.DiffEqual
			newD.Text = linePrefix
			//fmt.Fprintf(os.Stderr, "  [%d(%d)] %s %q\n", insIndex, 0, newD.Type, newD.Text)
			result[insIndex] = newD

			back := 0

			if pos+suppressLen < len(d.Text) {
				newD.Type = d.Type
				newD.Text = d.Text[pos+suppressLen:]
				//fmt.Fprintf(os.Stderr, "  [%d(%d)] %s %q\n", insIndex+1, +1, newD.Type, newD.Text)
				result = slices.Insert(result, insIndex+1, newD)

				back++
			}
			//debug
			//for i, d := range result {
			//	fmt.Fprintf(os.Stderr, "  %d: %v:%q\n", i, d.Type, d.Text)
			//}
			//println("---")

			if pos > 0 {
				newD := d
				newD.Text = d.Text[:pos]
				//fmt.Fprintf(os.Stderr, "  [%d(%d)] %s %q\n", insIndex, -1, newD.Type, newD.Text)
				result = slices.Insert(result, insIndex, newD)

				back++
			}
			//debug
			//for i, d := range result {
			//	fmt.Fprintf(os.Stderr, "  %d: %v:%q\n", i, d.Type, d.Text)
			//}
			//println("---")

			//println("  back", back, insIndex)
			insIndex += back
			//println(" >back", back, insIndex)
		}

		index--
	}

	return result
}

func splitByNewline(diffs []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	results := make([]diffmatchpatch.Diff, 0, len(diffs))

	for _, d := range diffs {
		dText := d.Text
		for {
			pos := strings.Index(dText, "\n")
			if pos == -1 {
				break
			}

			newD := d
			newD.Text = dText[:pos+1]
			results = append(results, newD)

			dText = dText[pos+1:]
		}

		if dText != "" {
			newD := d
			newD.Text = dText
			results = append(results, newD)
		}
	}

	return results
}

func addIndents(diffs []diffmatchpatch.Diff, indent string) []diffmatchpatch.Diff {
	newlined := false

	results := make([]diffmatchpatch.Diff, 0, len(diffs))

	for _, d := range diffs {
		dText := d.Text
		if newlined {
			newD := diffmatchpatch.Diff{
				Type: diffmatchpatch.DiffEqual,
				Text: indent,
			}
			results = append(results, newD)
		}

		results = append(results, d)

		newlined = strings.Contains(dText, "\n")
	}

	return results
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

func countIndent(line string) int {
	for i := range len(line) {
		if line[i] != ' ' {
			return i
		}
	}
	return 0
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
