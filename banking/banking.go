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