package main

import (
	"errors"
	"fmt"
	"sync"
)

type Account struct {
	mutex sync.Mutex
	ID      uint
	Name 	string
	Email 	string
	Wallet 	Wallet
}

type Wallet struct {
	ID.     uint
	Balance float64
}

var accounts []Account

type Repo struct {}

func (r Repo) Get(ID uint) *Account {
	for _, val := range accounts {
		if val.ID == ID {
			return &val
		}
	}

	return &Account{}
}

type AccountRepository interface {
	Save(account *Account) (*Account, error)
	FindAll() (*[]Account, error)
	Get(ID uint) *Account
}

type TransactionRepository interface {
	Transfer(transactionNO string, sender *Account, receiver *Account, amount float64, wg *sync.WaitGroup, ch chan error)
}

func NewTransactionRepository() TransactionRepository  {
	return &Repo{}
}

func (r Repo) Save(account *Account) (*Account, error) {
	accounts = append(accounts, *account)
	return account, nil
}

func (r Repo) FindAll() (*[]Account, error) {
	return &accounts, nil
}

func NewAccountRepository() AccountRepository  {
	return &Repo{}
}

func (r Repo) Transfer(transactionNO string, sender *Account, receiver *Account, amount float64, wg *sync.WaitGroup, ch chan error) {
	sender.mutex.Lock()
	defer sender.mutex.Unlock()

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	defer wg.Done()

	senderBalanceToBe := sender.Wallet.Balance - amount
	if senderBalanceToBe < 0 {
		ch <- errors.New(transactionNO + " balance is not enough")
	}else {
		sender.Wallet.Balance = senderBalanceToBe
		receiver.Wallet.Balance += amount
		ch <- nil
	}
}

func main(){

	repoTransaction := NewTransactionRepository()
	repoAccounts := NewAccountRepository()

	sam := Account{ID: 1, Name: "sam", Email: "sam@gmail.com", Wallet: Wallet{Balance: 200_000}}
	dev := Account{ID: 2, Name: "dev", Email: "dev@gmail.com", Wallet: Wallet{Balance: 200_000}}
	sammi := Account{ID: 3, Name: "dev", Email: "dev@gmail.com", Wallet: Wallet{Balance: 200_000}}

	ch := make(chan error)
	var wg sync.WaitGroup

	wg.Add(1)
	go repoTransaction.Transfer("TO1", &sam, &dev, 1, &wg, ch)
	wg.Add(1)
	go repoTransaction.Transfer("TO2", &dev, &sam, 1, &wg, ch)
	wg.Add(1)
	go repoTransaction.Transfer("TO3", &sam, &dev, 1, &wg, ch)
	wg.Add(1)
	go repoTransaction.Transfer("TO4", &dev, &sam, 1, &wg, ch)

	go func() {
		wg.Wait()
		close(ch)
	}()

	for err := range ch {
		if err == nil {
		}else {
			fmt.Println(err)
		}
	}

	repoAccounts.Save(&sam)
	repoAccounts.Save(&dev)
	repoAccounts.Save(&sammi)

	fmt.Println(*repoAccounts.Get(1))
	fmt.Println(*repoAccounts.Get(2))
}
