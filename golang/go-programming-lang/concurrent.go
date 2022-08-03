package main

import (
	"fmt"
	"go-programming-lang/concurrent/bank1"
)

func main() {
	bank1.Deposit(100)
	bank1.Deposit(200)
	bank1.Withdraw(150)
	fmt.Println(bank1.Balance())
	bank1.Withdraw(200)
	fmt.Println(bank1.Balance())
}
