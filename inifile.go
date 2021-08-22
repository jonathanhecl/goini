package inifile

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

/*
*
* Easy INI File manager for Golang
* Jonathan G. Hecl
* https://github.com/jonathanhecl/inifile
*
 */

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

var _Section []byte = []byte{91, 93}
var _KeyValueDiff byte = byte(61)
var _FlagComments []byte = []byte{35, 39, 47, 96}
var _IgnoredSpaces []byte = []byte{9, 10, 32}

type TValue struct {
	Value string
}

type TINIFile struct {
	lines      []_TLine
	Filename   string
	TotalLines int
	options    TOptions
}

type TOptions struct {
	Debug         bool
	CaseSensitive bool
}

func (t *TINIFile) Options(o TOptions) {
	(*t).options = o
}

func New() *TINIFile {
	t := TINIFile{}
	t.lines = []_TLine{}
	t.Filename = ""
	t.TotalLines = 0
	t.options = TOptions{
		Debug:         false,
		CaseSensitive: false,
	}
	return &t
}

func Load(Path string) (*TINIFile, error) {
	t := TINIFile{}
	t.lines = []_TLine{}
	t.Filename = Path
	t.TotalLines = 0
	t.options = TOptions{
		Debug:         false,
		CaseSensitive: false,
	}
	if f, err := os.Open(t.Filename); err != nil {
		return nil, err
	} else {
		defer f.Close()
		s := bufio.NewScanner(f)
		lineNumber := 0
		for s.Scan() {
			l := strings.TrimSpace(s.Text())
			if lineNumber == 0 {
				t.lines = append(t.lines, processLine(l, _TLine{}))
			} else {
				t.lines = append(t.lines, processLine(l, t.lines[lineNumber-1]))
			}
			lineNumber++
		}
		t.TotalLines = lineNumber
	}
	return &t, nil
}

func (t *TINIFile) Save(Path string) error {
	if f, err := os.Create(Path); err != nil {
		return err
	} else {
		defer f.Close()
		lineBreak := "\r"
		if runtime.GOOS == "windows" {
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

func processLine(s string, prevLine _TLine) _TLine {
	r := _TLine{
		Mode: IGNORED,
		Line: s,
	}
	ignoringBeginning := true
	ignoringComment := false
	capturingSection := false
	capturingKey := false
	capturingValue := false
	tempReading := ""
	for i := range s {
		if ignoringBeginning && !bytes.Contains(_IgnoredSpaces, []byte{s[i]}) {
			ignoringBeginning = false
			capturingKey = true
		}
		if !ignoringBeginning {
			if !ignoringComment && bytes.Contains(_FlagComments, []byte{s[i]}) {
				ignoringComment = true
				capturingKey = false
				break
			}
			if (capturingSection || capturingKey) && !capturingValue && bytes.Contains(_IgnoredSpaces, []byte{s[i]}) {
				capturingSection = false
				capturingKey = false
				break
			}
			if !capturingSection && _Section[0] == s[i] {
				capturingSection = true
				capturingKey = false
				continue
			} else if capturingSection && _Section[1] == s[i] {
				r.Mode = SECTION
				r.Section = tempReading
				r.Key = ""
				r.Value = ""
				capturingSection = false
				break
			}
			if capturingKey && _KeyValueDiff == s[i] {
				r.Mode = KEY
				r.Section = prevLine.Section
				r.Key = tempReading
				r.Value = ""
				tempReading = ""
				capturingValue = true
				continue
			}
			if !ignoringComment {
				tempReading = tempReading + string(s[i])
				if capturingValue {
					r.Value = strings.TrimSpace(tempReading)
				}
			}
		}
	}
	return r
}

func (t *TINIFile) Set(section string, key string, value TValue) {
	sectionFound := -1
	for i := range t.lines {
		if (!t.options.CaseSensitive && strings.EqualFold(t.lines[i].Section, section)) ||
			(t.options.CaseSensitive && t.lines[i].Section == section) {
			sectionFound = i
			if (!t.options.CaseSensitive && strings.EqualFold(t.lines[i].Key, key)) ||
				(t.options.CaseSensitive && t.lines[i].Key == key) {
				if t.options.Debug {
					fmt.Println("EDIT VALUE: ", section, "->", key, "=", value.Value)
				}
				tempKey := t.lines[i].Line[:strings.Index(t.lines[i].Line, key)+len(key+string(_KeyValueDiff))]
				tempRest := t.lines[i].Line[len(tempKey):]
				tempNonValue := ""
				if strings.Index(tempRest, t.lines[i].Value) > len(tempRest) {
					tempNonValue = tempRest[strings.Index(tempRest, t.lines[i].Value)+len(t.lines[i].Value):]
				}
				t.lines[i].Value = value.Value
				t.lines[i].Line = tempKey + value.Value + tempNonValue
				return
			}
		}
	}
	if len(value.Value) > 0 {
		if sectionFound >= 0 {
			if t.options.Debug {
				fmt.Println("NEW KEY: ", section, "->", key, "=", value.Value)
			}
			newLine := _TLine{
				Mode:    KEY,
				Section: section,
				Key:     key,
				Value:   value.Value,
				Line:    key + string(_KeyValueDiff) + value.Value,
			}
			tempLines := t.lines[:sectionFound]
			tempLines = append(tempLines, newLine)
			tempLines = append(tempLines, t.lines[sectionFound+1:]...)
			(*t).lines = tempLines
		} else {
			if t.options.Debug {
				fmt.Println("NEW SECTION: ", section, "->", key, "=", value.Value)
			}
			newLines := []_TLine{
				{
					Mode:    SECTION,
					Section: section,
					Line:    string(_Section[0]) + section + string(_Section[1]),
				},
				{
					Mode:    KEY,
					Section: section,
					Key:     key,
					Value:   value.Value,
					Line:    key + string(_KeyValueDiff) + value.Value,
				},
			}
			tempLines := t.lines
			tempLines = append(tempLines, newLines...)
			(*t).lines = tempLines
		}
	}
}

func (t *TINIFile) Get(section string, key string) TValue {
	for i := range t.lines {
		if (!t.options.CaseSensitive && strings.EqualFold(t.lines[i].Section, section)) ||
			(t.options.CaseSensitive && t.lines[i].Section == section) {
			if (!t.options.CaseSensitive && strings.EqualFold(t.lines[i].Key, key)) ||
				(t.options.CaseSensitive && t.lines[i].Key == key) {
				return TValue{
					Value: t.lines[i].Value,
				}
			}
		}
	}
	return TValue{}
}

// Convertions

func String(s string) TValue {
	return TValue{Value: strings.TrimSpace(s)}
}

func (t TValue) String() string {
	return t.Value
}

func Bool(b bool, isInt bool) TValue {
	s := ""
	if isInt {
		s = "0"
		if b {
			s = "1"
		}
	} else {
		s = "false"
		if b {
			s = "true"
		}
	}
	return TValue{Value: s}
}

func (t TValue) Bool() bool {
	b := false
	if strings.EqualFold(t.Value, "true") ||
		t.Value == "1" {
		b = true
	}
	return b
}

func Int(i int) TValue {
	s := strconv.Itoa(i)
	return TValue{Value: s}
}

func (t TValue) Int() int {
	i, _ := strconv.Atoi(t.Value)
	return i
}

func Int8(i int8) TValue {
	s := strconv.Itoa(int(i))
	return TValue{Value: s}
}

func (t TValue) Int8() int8 {
	i, _ := strconv.Atoi(t.Value)
	return int8(i)
}

func Int16(i int16) TValue {
	s := strconv.Itoa(int(i))
	return TValue{Value: s}
}

func (t TValue) Int16() int16 {
	i, _ := strconv.Atoi(t.Value)
	return int16(i)
}

func Int32(i int32) TValue {
	s := strconv.Itoa(int(i))
	return TValue{Value: s}
}

func (t TValue) Int32() int32 {
	i, _ := strconv.Atoi(t.Value)
	return int32(i)
}

func Int64(i int64) TValue {
	s := strconv.Itoa(int(i))
	return TValue{Value: s}
}

func (t TValue) Int64() int64 {
	i, _ := strconv.Atoi(t.Value)
	return int64(i)
}

func Float64(i float64) TValue {
	s := fmt.Sprint(i)
	return TValue{Value: s}
}

func (t TValue) Float64() float64 {
	i, _ := strconv.ParseFloat(t.Value, 64)
	return i
}

func Uint64(i uint64) TValue {
	s := strconv.FormatUint(i, 10)
	return TValue{Value: s}
}

func (t TValue) UInt64() uint64 {
	i, _ := strconv.ParseUint(t.Value, 10, 64)
	return i
}
