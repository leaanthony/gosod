# gosod

An application scaffold library.

## Installation 

`go get github.com/leaanthony/gosod`

## Usage

  1. Define a template directory
  2. Define some data
  3. Extract to a target directory

```
package main

import (
	"log"

	"github.com/leaanthony/gosod"
)

type config struct {
	Name string
}

func main() {

	// Define a new Template directory
	basic, err := gosod.TemplateDir("./myTemplate")
	if err != nil {
		log.Fatal(err)
	}

	// Make some config data
	myConfig := &config{
		Name: "Mat",
	}

	// Create a new directory using the template and config
	err = basic.Extract("./generated", myConfig)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Template Directories

A template directory is simply a directory structure contianing files you wish to copy. The algorithm for copying is:

  * Categorise all files into one of: directory, standard file and template files
	* Create the directory structure
	* Copy standard files
	* Copy template files, assembled using the given data

Template files, by default, are any file with ".tmpl" in their filename. To change this, use `SetTemplateFilters([]string)`. This allows you to set any number of filters.

Files may also be ignored by using the `IgnoreFilename(string)` method.