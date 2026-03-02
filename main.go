package main

import (
	"fmt"

	"github.com/wootaiklee/side-project/banking"
)

func main() {
	account := banking.NewAccount("John Doe")
	account.Deposit(100)
	account.Deposit(100)
	err := account.Withdraw(250)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(account.Balance())
}	