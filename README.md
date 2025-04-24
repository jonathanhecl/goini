# goini
## INI File manager package on Golang

[![Go](https://github.com/jonathanhecl/goini/actions/workflows/go.yml/badge.svg)](https://github.com/jonathanhecl/goini/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jonathanhecl/goini)](https://goreportcard.com/report/github.com/jonathanhecl/goini)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

With goini package you can create or read INI files, preserving comments and types, very quickly and easily.

> coverage: 94.0% of statements
>

> go get github.com/jonathanhecl/goini
>

## 🫴 Features:

* You can get and set values easily.
* The sections and keys are created dynamically.
* Preserve all the comments.
* Preserve empty lines and blank lines.
* Works with big and small files quickly.

## 🔨 Example:
```
package main


import (
    "github.com/jonathanhecl/goini"
)

func main() {

    // Load an existent file
	ini, _ := goini.Load("./test.ini", nil)
    /*
        // New file
        ini := goini.New(&goini.TOptions{CaseSensitive: false})
    */
    // Read a key
	value := ini.Get("section", "key").String()
    // Set a key
	ini.Set("section", "key", goini.String("test"))
    // Save a file
	ini.Save("./new.ini")

}
```

## 🛠️ Types supported:

| Type | Notes |
| --- | --- |
| Byte | |
| String | |
| StringArray | separated with comma (,) |
| Bool | works with 1/0 and true/false |
| Int | |
| Int8 | |
| Int16 | |
| Int32 | |
| Int64 | |
| Uint64 | |
| Float32 | |
| Float64 | |