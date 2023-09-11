# goini
## INI File manager package on Golang

[![Go](https://github.com/jonathanhecl/goini/actions/workflows/go.yml/badge.svg)](https://github.com/jonathanhecl/goini/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jonathanhecl/goini)](https://goreportcard.com/report/github.com/jonathanhecl/goini)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

With this package you can create or read INI files, preserving comments and types, very quick & easy.

> coverage: 83.9% of statements
>

> go get github.com/jonathanhecl/goini
> 

## Features:

* Get & Set values.
* Create sections & keys dinamicaly.
* Preserve all the comments.
* Works with big & small files.

## Example:
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
	value := ini.Get("section", "key").String())
    // Set a key
	ini.Set("section", "key", goini.String("test"))
    // Save a file
	ini.Save("./new.ini")

}
```

## Types supported:

* Byte
* String
* StringArray _(separated with comma)_
* Bool _(works with 1/0 and true/false)_
* Int
* Int8
* Int16
* Int32
* Uint64
* Float32
* Float64

## All made on day 21 Agu 2021
## Updated 11 Sep 2023