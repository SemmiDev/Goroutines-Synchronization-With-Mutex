package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Account struct {
	mutex sync.Mutex
	Name 	string
	Balance float64
}

func transferSync(transactionNO string,
		sender *Account,
		receiver *Account,
		amount float64,
		wg *sync.WaitGroup,
		ch chan error)  {
	sender.mutex.Lock()
	defer sender.mutex.Unlock()

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	defer wg.Done()

	senderBalanceToBe := sender.Balance - amount
	//sleep(500)
	if senderBalanceToBe < 0 {
		 ch <- errors.New(transactionNO + " balance is not enough")
	}else {
		sender.Balance = senderBalanceToBe
		receiver.Balance += amount
		ch <- nil
	}
}

func sleep(duration int) {
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

func main(){
	sam := Account{Name: "sam", Balance: 200_000}
	dev := Account{Name: "dev", Balance: 200_000}

	ch := make(chan error)
	var wg sync.WaitGroup

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go transferSync("TO1", &sam, &dev, 2000, &wg, ch)
		wg.Add(1)
		go transferSync("TO1", &dev, &sam,5000, &wg, ch)
	}

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

	fmt.Println(sam)
	fmt.Println(dev)
}
