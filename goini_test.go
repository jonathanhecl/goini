package goini

import (
	"fmt"
	"os"
	"testing"
)

type TestValue struct {
	String      string
	Bool        bool
	Byte        byte
	Int         int
	Int8        int8
	Int16       int16
	Int32       int32
	Int64       int64
	UInt64      uint64
	Float32     float32
	Float64     float64
	StringArray []string
}

var testValues = []TestValue{
	{String: "test"},
	{Bool: true},
	{Bool: false},
	{Byte: byte(255)},
	{Int: 1234567},
	{Int8: 123},
	{Int16: 12345},
	{Int32: 1234567890},
	{Int64: 1234567890123456789},
	{Float32: 1234567890.1234567890},
	{Float64: 1234567890123456789.1234567890123456789},
	{UInt64: 12345678901234567890},
	{String: "https://www.jonathanhecl.com/"},
	{String: "//not a comment"},
	{StringArray: []string{"test", "test2"}},
}
var specialString = "it is a string=with slash // no comment"

func TestCreateNewFile(t *testing.T) {
	ini := New(&TOptions{Debug: true})
	for i, v := range testValues {
		if len(v.StringArray) > 0 {
			ini.Set("Test", fmt.Sprintf("%dstringarray", i), StringArray(v.StringArray))
		} else if len(v.String) > 0 {
			ini.Set("Test", fmt.Sprintf("%dstring", i), String(v.String))
		} else if v.UInt64 > 0 {
			ini.Set("Test", fmt.Sprintf("%duint64", i), Uint64(v.UInt64))
		} else if v.Float64 > 0 {
			ini.Set("Test", fmt.Sprintf("%dfloat64", i), Float64(v.Float64))
		} else if v.Float32 > 0 {
			ini.Set("Test", fmt.Sprintf("%dfloat32", i), Float32(v.Float32))
		} else if v.Int64 > 0 {
			ini.Set("Test", fmt.Sprintf("%dint64", i), Int64(v.Int64))
		} else if v.Int32 > 0 {
			ini.Set("Test", fmt.Sprintf("%dint32", i), Int32(v.Int32))
		} else if v.Int16 > 0 {
			ini.Set("Test", fmt.Sprintf("%dint16", i), Int16(v.Int16))
		} else if v.Int8 > 0 {
			ini.Set("Test", fmt.Sprintf("%dint8", i), Int8(v.Int8))
		} else if v.Int > 0 {
			ini.Set("Test", fmt.Sprintf("%dint", i), Int(v.Int))
		} else if v.Byte > 0 {
			ini.Set("Test", fmt.Sprintf("%dbyte", i), Byte(v.Byte))
		} else {
			ini.Set("Test", fmt.Sprintf("%dbool", i), Bool(v.Bool, v.Bool))
		}
	}
	ini.Set("Test", "specialString", String(specialString))
	err := ini.Save("test.ini")
	if err != nil {
		t.Error(err)
	}
}

func TestReadFile(t *testing.T) {
	ini, err := Load("test.ini", &TOptions{Debug: true})
	if err != nil {
		t.Error(err)
	}
	for i, v := range testValues {
		if len(v.StringArray) > 0 {
			stra := ini.Get("Test", fmt.Sprintf("%dstringarray", i)).StringArray()
			if len(stra) != len(v.StringArray) {
				t.Errorf("Expected %s, got %s", v.StringArray, ini.Get("Test", fmt.Sprintf("%dstringarray", i)).StringArray())
			}
		} else if len(v.String) > 0 {
			if ini.Get("Test", fmt.Sprintf("%dstring", i)).String() != v.String {
				t.Errorf("Expected %s, got %s", v.String, ini.Get("Test", fmt.Sprintf("%dstring", i)).String())
			}
		} else if v.UInt64 > 0 {
			if ini.Get("Test", fmt.Sprintf("%duint64", i)).Uint64() != v.UInt64 {
				t.Errorf("Expected %d, got %d", v.UInt64, ini.Get("Test", fmt.Sprintf("%duint64", i)).Uint64())
			}
		} else if v.Float64 > 0 {
			if ini.Get("Test", fmt.Sprintf("%dfloat64", i)).Float64() != v.Float64 {
				t.Errorf("Expected %f, got %f", v.Float64, ini.Get("Test", fmt.Sprintf("%dfloat64", i)).Float64())
			}
		} else if v.Float32 > 0 {
			if ini.Get("Test", fmt.Sprintf("%dfloat32", i)).Float32() != v.Float32 {
				t.Errorf("Expected %f, got %f", v.Float32, ini.Get("Test", fmt.Sprintf("%dfloat32", i)).Float32())
			}
		} else if v.Int64 > 0 {
			if ini.Get("Test", fmt.Sprintf("%dint64", i)).Int64() != v.Int64 {
				t.Errorf("Expected %d, got %d", v.Int64, ini.Get("Test", fmt.Sprintf("%dint64", i)).Int64())
			}
		} else if v.Int32 > 0 {
			if ini.Get("Test", fmt.Sprintf("%dint32", i)).Int32() != v.Int32 {
				t.Errorf("Expected %d, got %d", v.Int32, ini.Get("Test", fmt.Sprintf("%dint32", i)).Int32())
			}
		} else if v.Int16 > 0 {
			if ini.Get("Test", fmt.Sprintf("%dint16", i)).Int16() != v.Int16 {
				t.Errorf("Expected %d, got %d", v.Int16, ini.Get("Test", fmt.Sprintf("%dint16", i)).Int16())
			}
		} else if v.Int8 > 0 {
			if ini.Get("Test", fmt.Sprintf("%dint8", i)).Int8() != v.Int8 {
				t.Errorf("Expected %d, got %d", v.Int8, ini.Get("Test", fmt.Sprintf("%dint8", i)).Int8())
			}
		} else if v.Int > 0 {
			if ini.Get("Test", fmt.Sprintf("%dint", i)).Int() != v.Int {
				t.Errorf("Expected %d, got %d", v.Int, ini.Get("Test", fmt.Sprintf("%dint", i)).Int())
			}
		} else if v.Byte > 0 {
			if ini.Get("Test", fmt.Sprintf("%dbyte", i)).Byte() != v.Byte {
				t.Errorf("Expected %d, got %d", v.Byte, ini.Get("Test", fmt.Sprintf("%dbyte", i)).Byte())
			}
		} else {
			if ini.Get("Test", fmt.Sprintf("%dbool", i)).Bool() != v.Bool {
				t.Errorf("Expected %t, got %t", v.Bool, ini.Get("Test", fmt.Sprintf("%dbool", i)).Bool())
			}
		}
	}
	if ini.Get("Test", "specialString").String() != specialString {
		t.Errorf("Expected %s, got %s", specialString, ini.Get("Test", "specialString").String())
	}
}

func TestSpecial2(t *testing.T) {
	content := []byte(`[Test]
change=4 ' comment
ignore=I'will ignore this
same=Never change this	// comment`)
	err := os.WriteFile("test2.ini", content, 0644)
	if err != nil {
		t.Errorf("Error creating test file: %s", err)
	}

	ini, err := Load("test2.ini", &TOptions{Debug: true, ForceSaveWithoutQuotes: true})
	if err != nil {
		t.Error(err)
	}

	if ini.Get("Test", "change").Int() != 4 {
		t.Errorf("Expected 4, got %d", ini.Get("Test", "change").Int())
	}

	if ini.Get("Test", "ignore").String() != "I'will ignore this" {
		t.Errorf("Expected I'will ignore this, got %s", ini.Get("Test", "ignore").String())
	}

	if ini.Get("Test", "same").String() != "Never change this" {
		t.Errorf("Expected Never change this, got %s", ini.Get("Test", "same").String())
	}

	// Test save
	ini.Set("Test", "change", Int(5))
	ini.Set("Test", "ignore", String("I'will change this"))
	ini.Set("Test", "same", String("Never change this"))
	err = ini.Save("test2.ini")
	if err != nil {
		t.Error(err)
	}

	ini, err = Load("test2.ini", &TOptions{Debug: true})
	if err != nil {
		t.Error(err)
	}

	if ini.Get("Test", "change").Int() != 5 {
		t.Errorf("Expected 5, got %d", ini.Get("Test", "change").Int())
	}

	if ini.Get("Test", "ignore").String() != "I'will change this" {
		t.Errorf("Expected I'will change this, got %s", ini.Get("Test", "ignore").String())
	}

	if ini.Get("Test", "same").String() != "Never change this" {
		t.Errorf("Expected Never change this, got %s", ini.Get("Test", "same").String())
	}
}
