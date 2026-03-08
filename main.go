package main

import (
	"fmt"

	"github.com/wootaiklee/side-project/dict"
)

func main() {
	dictionary := dict.Dictionary{"first": "first word"}
	definition, err := dictionary.Search("second")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(definition)
}	