package goini

import (
	"fmt"
	"strconv"
	"strings"
)

// Conversions

func String(s string) TValue {
	return TValue{Value: []byte(strings.TrimSpace(s))}
}

func (t TValue) String() string {
	return string(ValueToRead(t.Value))
}

func StringArray(s []string) TValue {
	return TValue{Value: []byte(strings.Join(s, string(_ArraySeparator)))}
}

func (t TValue) StringArray() []string {
	return strings.Split(string(t.Value), string(_ArraySeparator))
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
	return TValue{Value: []byte(s)}
}

func (t TValue) Bool() bool {
	b := false
	if strings.EqualFold(string(t.Value), "true") ||
		string(t.Value) == "1" {
		b = true
	}
	return b
}

func Byte(i byte) TValue {
	s := strconv.Itoa(int(i))
	return TValue{Value: []byte(s)}
}

func (t TValue) Byte() byte {
	i, _ := strconv.Atoi(string(t.Value))
	return byte(i)
}

func Int(i int) TValue {
	s := strconv.Itoa(i)
	return TValue{Value: []byte(s)}
}

func (t TValue) Int() int {
	i, _ := strconv.Atoi(string(t.Value))
	return i
}

func Int8(i int8) TValue {
	s := strconv.Itoa(int(i))
	return TValue{Value: []byte(s)}
}

func (t TValue) Int8() int8 {
	i, _ := strconv.Atoi(string(t.Value))
	return int8(i)
}

func Int16(i int16) TValue {
	s := strconv.Itoa(int(i))
	return TValue{Value: []byte(s)}
}

func (t TValue) Int16() int16 {
	i, _ := strconv.Atoi(string(t.Value))
	return int16(i)
}

func Int32(i int32) TValue {
	s := strconv.Itoa(int(i))
	return TValue{Value: []byte(s)}
}

func (t TValue) Int32() int32 {
	i, _ := strconv.Atoi(string(t.Value))
	return int32(i)
}

func Int64(i int64) TValue {
	s := strconv.Itoa(int(i))
	return TValue{Value: []byte(s)}
}

func (t TValue) Int64() int64 {
	i, _ := strconv.Atoi(string(t.Value))
	return int64(i)
}

func Float32(i float32) TValue {
	s := fmt.Sprint(i)
	return TValue{Value: []byte(s)}
}

func (t TValue) Float32() float32 {
	i, _ := strconv.ParseFloat(string(t.Value), 32)
	return float32(i)
}

func Float64(i float64) TValue {
	s := fmt.Sprint(i)
	return TValue{Value: []byte(s)}
}

func (t TValue) Float64() float64 {
	i, _ := strconv.ParseFloat(string(t.Value), 64)
	return i
}

func Uint64(i uint64) TValue {
	s := strconv.FormatUint(i, 10)
	return TValue{Value: []byte(s)}
}

func (t TValue) Uint64() uint64 {
	i, _ := strconv.ParseUint(string(t.Value), 10, 64)
	return i
}
