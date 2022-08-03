// Package bank1 provides a concurrency-safe bank with one account.
package bank1

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance
var withdraws = make(chan struct {
	amount  int
	success chan bool
}) // withdraws statuses

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }
func Withdraw(amount int) bool {
	var success = make(chan bool)
	withdraws <- struct {
		amount  int
		success chan bool
	}{amount: amount, success: success}
	return <-success
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case withdraw := <-withdraws:
			if withdraw.amount > balance {
				withdraw.success <- false
			} else {
				balance -= withdraw.amount
				withdraw.success <- true
			}
		case balances <- balance:
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}
