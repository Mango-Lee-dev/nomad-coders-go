package main

import (
	"fmt"

	"github.com/wootaiklee/side-project/banking"
)

func main() {
	account := banking.NewAccount("John Doe")
	fmt.Println(account)
}	