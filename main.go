package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	insIndent     int
	commentIndent int
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [params] [file1 [file2 ...]]\nParameters:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.IntVar(&insIndent, "ii", 8, "Indentation for instructions in spaces")
	flag.IntVar(&commentIndent, "ci", 40, "Indentation for comments in spaces")
}

type asmLine struct {
	label   string
	text    string
	comment string
}

var (
	mulSpaces   = regexp.MustCompile(` +`)
	commaSpaces = regexp.MustCompile(`, *`)
)

var dataPseudos = compilePseudoes()

func compilePseudoes() (ret []*regexp.Regexp) {
	var (
		raw = []string{"db", "dw", "dd", "dq", "ddq", "do", "dt"}
		re  = `(?:\s|^)%v(?:\s|")`
	)
	for _, p := range raw {
		ret = append(ret, regexp.MustCompile(fmt.Sprintf(re, p)))
	}
	return
}

func parseLabel(line string) (lbl string, rest string) {
	ind := strings.Index(line, ":")
	if ind >= 0 {
		lbl = strings.TrimSpace(line[:ind])
		rest = line[ind+1:]
		return
	}

	for _, pseudo := range dataPseudos {
		inds := pseudo.FindStringIndex(line)
		if len(inds) > 0 {
			ind := inds[0]
			lbl = strings.TrimSpace(line[:ind])
			rest = line[ind:]
			return
		}
	}

	return "", line
}

func parseComment(line string) (cmt string, rest string) {
	ind := strings.Index(line, ";")
	if ind >= 0 {
		cmt = strings.TrimSpace(line[ind+1:])
		rest = line[:ind]
		return
	}
	return "", line
}

func parseLine(line string) *asmLine {
	l := &asmLine{}

	l.label, line = parseLabel(line)
	l.comment, line = parseComment(line)

	l.text = strings.TrimSpace(line)
	l.text = mulSpaces.ReplaceAllString(l.text, " ")
	l.text = commaSpaces.ReplaceAllString(l.text, ", ")

	return l
}

func (l *asmLine) empty() bool {
	return l.label == "" && l.text == "" && l.comment == ""
}

func (l *asmLine) print(w io.Writer) {
	var (
		space  = []byte{' '}
		newl   = []byte{'\n'}
		column = 0
	)

	if l.label != "" {
		w.Write([]byte(l.label))
		w.Write([]byte{':'})
		column += utf8.RuneCountInString(l.label) + 1
		if l.text != "" {
			w.Write(newl)
			column = 0
		}
	}

	if l.text != "" {
		w.Write(bytes.Repeat(space, insIndent))
		w.Write([]byte(l.text))
		column += insIndent
		column += utf8.RuneCountInString(l.text)
	}

	if l.comment != "" {
		if column != 0 {
			if column < commentIndent-1 {
				w.Write(bytes.Repeat(space, commentIndent-column-1))
			} else {
				w.Write(space)
			}
		}
		w.Write([]byte{';', ' '})
		w.Write([]byte(l.comment))
	}

	w.Write(newl)
}

func formatto(filename, outname string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	out, err := os.Create(outname)
	if err != nil {
		return err
	}
	defer out.Close()

	scanner := bufio.NewScanner(file)
	lastEmpty := false
	for scanner.Scan() {
		line := scanner.Text()
		asmLine := parseLine(line)
		if lastEmpty && asmLine.empty() {
			continue // Output no more than one empty line
		}
		asmLine.print(out)
		lastEmpty = asmLine.empty()
	}
	return nil
}

func format(filename string) error {
	outname := filename + "~"
	if err := formatto(filename, outname); err != nil {
		return err
	}
	if err := os.Rename(outname, filename); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	files := flag.Args()
	for _, fn := range files {
		err := format(fn)
		if err != nil {
			fmt.Println(err)
		}
	}
}
