package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

type asmLine struct {
	label   string
	text    string
	comment string
}

var (
	mulSpaces   = regexp.MustCompile(` +`)
	commaSpaces = regexp.MustCompile(`, *`)
)

var (
	insIndent     = 8
	commentIndent = 40
)

var dataPseudos = compilePseudoes()

func compilePseudoes() (ret []*regexp.Regexp) {
	raw := []string{"db", "dw", "dd", "dq", "ddq", "do", "dt"}
	for _, p := range raw {
		ret = append(ret, regexp.MustCompile(fmt.Sprintf(`(?:\s|^)%v(?:\s|")`, p)))
	}
	return
}

func parseLabel(line string) (lbl string, rest string) {
	ind := strings.Index(line, ":")
	if ind >= 0 {
		lbl = line[:ind]
		lbl = strings.TrimSpace(lbl)
		rest = line[ind+1:]
		return
	}

	for _, pseudo := range dataPseudos {
		inds := pseudo.FindStringIndex(line)
		if len(inds) > 0 {
			ind = inds[0]
			lbl = line[:ind]
			lbl = strings.TrimSpace(lbl)
			rest = line[ind:]
			return
		}
	}

	return "", line
}

func parseComment(line string) (cmt string, rest string) {
	ind := strings.Index(line, ";")
	if ind >= 0 {
		cmt = line[ind+1:]
		cmt = strings.TrimSpace(cmt)
		rest = line[:ind]
		return
	}
	return "", line
}

func parseLine(line string) *asmLine {
	l := &asmLine{}
	line = strings.TrimSpace(line)

	l.label, line = parseLabel(line)
	l.comment, line = parseComment(line)

	l.text = line
	l.text = strings.TrimSpace(l.text)
	l.text = mulSpaces.ReplaceAllString(l.text, " ")
	l.text = commaSpaces.ReplaceAllString(l.text, ", ")

	return l
}

func (l *asmLine) empty() bool {
	return l.label == "" && l.text == "" && l.comment == ""
}

func (l *asmLine) print(w io.Writer) {
	var space = []byte{' '}

	column := 0
	if l.label != "" {
		w.Write([]byte(l.label))
		w.Write([]byte{':'})
		column += utf8.RuneCountInString(l.label) + 1
		if l.text != "" {
			w.Write([]byte{'\n'})
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
		if l.label != "" || l.text != "" {
			if column < commentIndent {
				w.Write(bytes.Repeat(space, commentIndent-column-1))
			} else {
				w.Write(space)
			}
		}
		w.Write([]byte{';', ' '})
		w.Write([]byte(l.comment))
	}

	w.Write([]byte{'\n'})
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
			continue
		}
		asmLine.print(out)
		lastEmpty = asmLine.empty()
	}
	return nil
}

func format(filename string) error {
	outname := filename + "~"
	err := formatto(filename, outname)
	if err != nil {
		return err
	}
	err = os.Remove(filename)
	if err != nil {
		return err
	}
	err = os.Rename(outname, filename)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v <files...>", os.Args[0])
	}
	files := os.Args[1:]
	for _, fn := range files {
		err := format(fn)
		if err != nil {
			fmt.Println(err)
		}
	}
}
