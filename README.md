# goini
## INI File manager package on Golang

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)

With this package you can create or read INI files, preserving comments and types, very quick & easy.

> coverage: 83.2% of statements
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
## Updated 25 Dec 2021