package banking

import (
	"errors"
)

// Account is a struct that represents a bank account
type Account struct {
	owner string
	balance int
}

var NoMoney = errors.New("Cannot withdraw, insufficient balance")

// NewAccount creates a new bank account
func NewAccount(owner string) *Account {
	return &Account{owner: owner, balance: 0}
}

// Deposit adds money to the account
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

// Balance returns the current balance of the account
func (a *Account) Balance() int {
	return a.balance
}

// Withdraw removes money from the account
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return NoMoney
	}
	a.balance -= amount
	return nil
}

// ChangeOwner changes the owner of the account
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

// Owner returns the owner of the account
func (a *Account) Owner() string {
	return a.owner
}