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

	// Register the filename to ignore
	basic.SetTemplateFilters([]string{".template", ".tmpl"})

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
