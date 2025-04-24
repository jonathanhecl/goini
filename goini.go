package goini

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

/*
*
* goini
* INI File manager package on Golang
* Created by Jonathan G. Hecl
* https://github.com/jonathanhecl/goini
*
 */

const (
	IsWindows = runtime.GOOS == "windows"
)

type _EType int8

const (
	SECTION _EType = iota
	KEY
	IGNORED
)

type _TLine struct {
	Mode    _EType
	Section string
	Key     string
	Value   string
	Line    string
}

type _TSection struct {
	Section string
	Begin   int
	End     int
}

var _Section []byte = []byte{91, 93}                  // [ ]
var _ArraySeparator []byte = []byte{44}               // 44 is the ascii code for comma
var _FlagComments []byte = []byte{35, 39, 47, 59, 96} // 47 double
var _IgnoredSpaces []byte = []byte{9, 10, 32}         // Bool returns the value as a boolean.
var _KeyValueDiff byte = byte(61)                     // 61 is the ascii code for =
var _FlagQuoting byte = byte(34)                      // 34 is the ascii code for "

type TValue struct {
	Value []byte
}

type TINIFile struct {
	lines      []_TLine
	sections   []_TSection
	Filename   string
	TotalLines int
	options    *TOptions
}

type TOptions struct {
	Debug                  bool
	CaseSensitive          bool
	DontPreserveEmptyLines bool
	ForceSaveWithoutQuotes bool
}

var timeMark time.Time

func (t *TINIFile) Options(o *TOptions) {
	(*t).options = o
}

func New(o *TOptions) *TINIFile {
	t := TINIFile{}
	t.lines = []_TLine{}
	t.sections = []_TSection{}
	t.Filename = ""
	t.TotalLines = 0
	t.options = o
	if t.options == nil {
		t.options = &TOptions{
			Debug:                  false,
			CaseSensitive:          false,
			DontPreserveEmptyLines: false,
			ForceSaveWithoutQuotes: false,
		}
	}
	return &t
}

func ReadFile(Path string, EmptyLines bool) ([]string, error) {
	f, err := os.Open(Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var (
		buf   []byte = make([]byte, 32*1024)
		lines []string
		line  []byte = []byte{}
	)
	for {
		n, err := f.Read(buf)
		if n > 0 {
			for i := 0; i < n; i++ {
				if buf[i] == 13 {
					// ignore \r
					continue
				} else if buf[i] == 10 {
					if len(line) > 0 || (EmptyLines && buf[i] == 10) {
						lines = append(lines, string(line))
						line = []byte{}
					}
				} else {
					line = append(line, buf[i])
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read %d bytes: %v", n, err)
		}
	}
	if len(line) > 0 || EmptyLines {
		lines = append(lines, string(line))
	}

	return lines, nil
}

func Load(Path string, o *TOptions) (*TINIFile, error) {
	t := TINIFile{}
	t.lines = []_TLine{}
	t.sections = []_TSection{}
	t.Filename = Path
	t.TotalLines = 0
	t.options = o
	if t.options == nil {
		t.options = &TOptions{
			CaseSensitive:          false,
			Debug:                  false,
			ForceSaveWithoutQuotes: false,
			DontPreserveEmptyLines: false,
		}
	}
	if t.options.Debug {
		timeMark = time.Now()
	}
	if lines, err := ReadFile(Path, !t.options.DontPreserveEmptyLines); err == nil {
		lineNumber := 0
		if t.options.Debug {
			fmt.Println("Total lines: ", len(lines))
			for i := range lines {
				fmt.Println("LINE ", i, " -> ", lines[i])
			}
		}

		for i := range lines {
			l := strings.TrimSpace(lines[i])
			if lineNumber == 0 {
				t.lines = append(t.lines, t.processLine(l, _TLine{}))
			} else {
				t.lines = append(t.lines, t.processLine(l, t.lines[lineNumber-1]))
			}
			lineNumber++
		}
	} else {
		return nil, err
	}
	if t.options.Debug {
		fmt.Println("File loaded on ", time.Since(timeMark))
	}
	return &t, nil
}

func (t *TINIFile) Save(Path string) error {
	if f, err := os.Create(Path); err != nil {
		return err
	} else {
		defer f.Close()

		lineBreak := "\r"
		if IsWindows {
			lineBreak = "\r\n"
		}
		for i := range t.lines {
			if _, err := f.Write([]byte(t.lines[i].Line + lineBreak)); err != nil {
				panic(err)
			}
		}
	}
	return nil
}

// Logic

func (t *TINIFile) processLine(line string, prevLine _TLine) _TLine {
	r := _TLine{
		Mode:    IGNORED,
		Section: prevLine.Section,
		Line:    line,
	}
	ignoringBeginning := true
	possibleComment := false
	ignoringComment := false
	possibleQuoting := false
	endingQuoting := 0
	capturingSection := false
	capturingKey := false
	capturingValue := false
	tempReading := []byte{}

	if len(line) == 0 {
		ignoringBeginning = true
	} else {
		for i := range line {
			if t.options.Debug {
				flagsStr := ""
				if ignoringBeginning {
					flagsStr += "ignoringBeginning "
				}
				if ignoringComment {
					flagsStr += "ignoringComment "
				}
				if capturingSection {
					flagsStr += "capturingSection "
				}
				if capturingKey {
					flagsStr += "capturingKey "
				}
				if capturingValue {
					flagsStr += "capturingValue "
				}
				if possibleComment {
					flagsStr += "possibleComment "
				}
				if possibleQuoting {
					flagsStr += "possibleQuoting "
				}
				if endingQuoting > 0 {
					flagsStr += fmt.Sprintf("endingQuoting(%d) ", endingQuoting)
				}
				if len(flagsStr) > 0 {
					flagsStr = flagsStr[:len(flagsStr)-1] // remove last space
				}
				fmt.Println(fmt.Sprintf("Previous flags: (%s) - Current character: %s", flagsStr, string(line[i])))
			}

			if ignoringBeginning && !bytes.Contains(_IgnoredSpaces, []byte{byte(line[i])}) {
				ignoringBeginning = false
				capturingKey = true
			}

			if !ignoringBeginning {
				if !ignoringComment && !possibleQuoting &&
					possibleComment && bytes.Contains(_FlagComments, []byte{byte(line[i])}) {
					isComment := true
					possibleComment = false
					if byte(line[i]) == 47 && len(line) > i { // 47 special
						if byte(line[i+1]) != 47 {
							isComment = false
						}
					}
					if isComment {
						ignoringComment = true
						capturingKey = false
						if t.options.Debug {
							fmt.Println("Ignoring Comments")
						}
						break
					}
				}

				if (capturingSection || capturingKey) &&
					!capturingValue && bytes.Contains(_IgnoredSpaces, []byte{byte(line[i])}) {
					capturingSection = false
					capturingKey = false
					if t.options.Debug {
						fmt.Println("End of key")
					}
					break
				}

				if !capturingSection &&
					_Section[0] == byte(line[i]) &&
					!capturingValue {
					capturingSection = true
					capturingKey = false

					if t.options.Debug {
						fmt.Println("Start of section")
					}
					continue
				} else if capturingSection && _Section[1] == byte(line[i]) {
					r.Mode = SECTION
					r.Section = string(tempReading)
					r.Key = ""
					r.Value = ""

					sectionKey := string(tempReading)
					if !t.options.CaseSensitive {
						sectionKey = strings.ToUpper(sectionKey)
					}

					sec := t.getSection(sectionKey)
					if sec == nil {
						t.sections = append(t.sections, _TSection{
							Section: sectionKey,
							Begin:   len(t.lines) + 1,
							End:     len(t.lines) + 1,
						})
					} else {
						sec.End = len(t.lines) + 1
					}

					capturingSection = false
					if t.options.Debug {
						fmt.Println("End of section")
					}
					break
				}

				if capturingKey && _KeyValueDiff == byte(line[i]) {
					r.Mode = KEY
					r.Section = prevLine.Section
					r.Key = string(tempReading)
					r.Value = ""
					tempReading = []byte{}
					capturingValue = true

					sectionKey := string(prevLine.Section)
					if !t.options.CaseSensitive {
						sectionKey = strings.ToUpper(sectionKey)
					}

					sec := t.getSection(sectionKey)
					if sec != nil {
						sec.End = len(t.lines) + 1
					}

					if t.options.Debug {
						fmt.Println("Start of key")
					}
					capturingKey = false
					continue
				}

				if !ignoringComment {
					tempReading = append(tempReading, byte(line[i]))
					if bytes.Contains(_IgnoredSpaces, []byte{byte(line[i])}) && // 9 tab
						!possibleQuoting {
						possibleComment = true
					} else {
						possibleComment = false
					}
					if _FlagQuoting == byte(line[i]) &&
						len(r.Value) == 0 {
						endingQuoting = strings.LastIndex(string(line[i:]), string(_FlagQuoting))
						if endingQuoting != i {
							possibleQuoting = true
						}
					} else if endingQuoting == i {
						endingQuoting = 0
						possibleQuoting = false
					}
					if capturingValue {
						r.Value = strings.TrimSpace(string(tempReading))
					}
				}
			}
		}
	}

	if t.options.Debug {
		fmt.Println("Line analyzed: ", string(line))
		fmt.Println("Line information: ", r)
	}

	return r
}

func (t *TINIFile) getSection(sectionKey string) *_TSection {
	for i := range t.sections {
		if t.sections[i].Section == sectionKey {
			return &t.sections[i]
		}
	}

	return nil
}

func (t *TINIFile) Set(section string, key string, value TValue) {
	sectionKey := section
	if !t.options.CaseSensitive {
		sectionKey = strings.ToUpper(sectionKey)
	}

	// Check if section does not exist, if so, create it
	sec := t.getSection(sectionKey)
	if sec == nil {
		if t.options.Debug {
			fmt.Println(fmt.Sprintf("Creating section [%s] with key [%s] and value [%s]", section, key, string(ValueToSave(value.Value, t.options.ForceSaveWithoutQuotes))))
		}

		t.sections = append(t.sections, _TSection{
			Section: sectionKey,
			Begin:   len(t.lines) + 1,
			End:     len(t.lines) + 2,
		})

		newLines := []_TLine{
			{
				Mode: IGNORED, // empty line
			},
			{
				Mode:    SECTION,
				Section: section,
				Line:    string(_Section[0]) + section + string(_Section[1]),
			},
			{
				Mode:    KEY,
				Section: section,
				Key:     key,
				Value:   string(ValueToSave(value.Value, t.options.ForceSaveWithoutQuotes)),
				Line:    key + string(_KeyValueDiff) + string(ValueToSave(value.Value, t.options.ForceSaveWithoutQuotes)),
			},
		}
		t.lines = append(t.lines, newLines...)
		(*t).lines = t.lines
		return
	}

	// if section exists, check if key exists, if so, change value
	var prevLine _TLine
	for i := sec.Begin; i <= sec.End && i < len(t.lines); i++ {
		if t.lines[i].Mode == KEY {
			prevLine = t.lines[i]
			if (!t.options.CaseSensitive && strings.EqualFold(t.lines[i].Key, key)) ||
				(t.options.CaseSensitive && t.lines[i].Key == key) {
				if t.lines[i].Value == string(ValueToSave(value.Value, t.options.ForceSaveWithoutQuotes)) {
					if t.options.Debug {
						fmt.Println(fmt.Sprintf("Ignoring value of key [%s] in section [%s], value is the same: [%s]", key, section, t.lines[i].Value))
					}
					return
				}

				if t.options.Debug {
					fmt.Println(fmt.Sprintf("Changing value of key [%s] in section [%s], previous value: [%s], new value: [%s]", key, section, t.lines[i].Value, string(ValueToSave(value.Value, t.options.ForceSaveWithoutQuotes))))
				}

				key = t.lines[i].Key
				tempKey := []byte(t.lines[i].Line[:strings.Index(t.lines[i].Line, key)+len(key+string(_KeyValueDiff))])
				tempRest := []byte(t.lines[i].Line[len(tempKey):])
				tempNonValue := []byte{}
				if len(t.lines[i].Value)+1 < len(tempRest) {
					tempNonValue = append([]byte{32}, tempRest[len(t.lines[i].Value)+1:]...)
				}

				(*t).lines[i].Value = string(ValueToSave(value.Value, t.options.ForceSaveWithoutQuotes))
				(*t).lines[i].Line = string(tempKey) + t.lines[i].Value + string(tempNonValue)
				if t.options.Debug {
					fmt.Println(fmt.Sprintf("Line changed, previous line: [%s], new line: [%s]", prevLine.Line, t.lines[i].Line))
				}
				return
			}
		}
	}

	// if section exists, check if key exists, if not, create it
	if len(value.Value) > 0 {
		if t.options.Debug {
			fmt.Println(fmt.Sprintf("Creating key [%s] in section [%s] with value [%s]", key, section, string(ValueToSave(value.Value, t.options.ForceSaveWithoutQuotes))))
		}

		sec := t.getSection(sectionKey)
		if sec == nil {
			return
		}
		sec.End++
		moving := false
		for i := range t.sections {
			if t.sections[i].Section == sectionKey {
				moving = true
			} else if moving {
				t.sections[i].Begin++
				t.sections[i].End++
			}
		}

		newLine := _TLine{
			Mode:    KEY,
			Section: section,
			Key:     key,
			Value:   string(ValueToSave(value.Value, t.options.ForceSaveWithoutQuotes)),
			Line:    key + string(_KeyValueDiff) + string(ValueToSave(value.Value, t.options.ForceSaveWithoutQuotes)),
		}

		t.lines = append(t.lines, _TLine{})
		copy(t.lines[sec.End-1:], t.lines[sec.End-2:])
		t.lines[sec.End-1] = newLine
		(*t).lines = t.lines
	}
}

func (t *TINIFile) Get(section string, key string) TValue {
	sectionKey := section
	if !t.options.CaseSensitive {
		sectionKey = strings.ToUpper(sectionKey)
	}

	sec := t.getSection(sectionKey)
	if sec == nil {
		return TValue{}
	}
	for i := sec.Begin; i < sec.End; i++ {
		if (!t.options.CaseSensitive && strings.EqualFold(t.lines[i].Section, section)) ||
			(t.options.CaseSensitive && t.lines[i].Section == section) {
			if (!t.options.CaseSensitive && strings.EqualFold(t.lines[i].Key, key)) ||
				(t.options.CaseSensitive && t.lines[i].Key == key) {
				return TValue{
					Value: []byte(t.lines[i].Value),
				}
			}
		}
	}

	return TValue{}
}

func ValueToSave(value []byte, forceWithoutQuotes bool) []byte {
	if !forceWithoutQuotes {
		flagQuote := false
		for i := range value {
			if value[i] == '\n' {
				value[i] = ' '
			} else if bytes.Contains(_FlagComments, []byte{byte(value[i])}) {
				flagQuote = true
			}
		}

		if flagQuote && len(value) > 0 {
			if value[0] != _FlagQuoting {
				value = append([]byte{_FlagQuoting}, value...)
				value = append(value, _FlagQuoting)
			}
		}
	}

	return value
}

func ValueToRead(value []byte) []byte {
	if len(value) > 0 {
		if value[0] == _FlagQuoting &&
			value[len(value)-1] == _FlagQuoting {
			value = value[1 : len(value)-1]
		}
	}

	return value
}
