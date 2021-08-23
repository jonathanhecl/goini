# inifile
Easy INI File manager for Golang

> go get github.com/jonathanhecl/inifile
> 

Features:

* Get & Set values
* Create sections & keys dinamicaly
* Preserve all the comments

Example:
```
package main


import (
	"github.com/jonathanhecl/inifile"
)

func main() {

    // Load an existent file
	ini, _ := inifile.Load("./test.ini", nil)
    /*
        // New file
        ini := inifile.New(&inifile.TOptions{CaseSensitive: false})
    */
    // Read a key
	value := ini.Get("section", "key").String())
    // Set a key
	ini.Set("section", "key", inifile.String("test"))
    // Save a file
	ini.Save("./new.ini")

}
```

Types predesigned:
* String
* Int
* Bool _(works with 1/0 and true/false)_
* Int8
* Int16
* Int32
* Uint64
* Float64

# All made on day 21 Agu 2021