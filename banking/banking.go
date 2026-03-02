package banking

// Account is a struct that represents a bank account
type Account struct {
	owner string
	balance int
}

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