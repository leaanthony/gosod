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