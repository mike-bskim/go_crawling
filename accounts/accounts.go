package accounts

import (
	"errors"
	"fmt"
)

// Account struct
type Account struct {
	owner   string
	balance int
}

func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

func (theAccount *Account) Deposit(amount int) {
	theAccount.balance += amount
}

func (theAccount Account) Balance() int {
	return theAccount.balance
}

func (theAccount *Account) Withdraw(amount int) error {
	if theAccount.balance <= amount {
		return errors.New("Can not withdraw, you are poor")
	}
	theAccount.balance -= amount
	return nil
}

func (theAccount *Account) ChangeOwner(new string) {
	theAccount.owner = new
}

func (theAccount Account) Owner() string {
	return theAccount.owner
}

func (theAccount Account) String() string {
	return fmt.Sprintf("Owner is %s, Balance: %d", theAccount.Owner(), theAccount.Balance())
}
