package main

/*
编写一个可以从 1 到 n 输出代表这个数字的字符串的程序，但是：

如果这个数字可以被 3 整除，输出 "fizz"。
如果这个数字可以被 5 整除，输出 "buzz"。
如果这个数字可以同时被 3 和 5 整除，输出 "fizzbuzz"。
例如，当 n = 15，输出： 1, 2, fizz, 4, buzz, fizz, 7, 8, fizz, buzz, 11, fizz, 13, 14, fizzbuzz。
请你实现一个有四个线程的多线程版 FizzBuzz， 同一个 FizzBuzz 实例会被如下四个线程使用：

线程A将调用 fizz() 来判断是否能被 3 整除，如果可以，则输出 fizz。
线程B将调用 buzz() 来判断是否能被 5 整除，如果可以，则输出 buzz。
线程C将调用 fizzbuzz() 来判断是否同时能被 3 和 5 整除，如果可以，则输出 fizzbuzz。
线程D将调用 number() 来实现输出既不能被 3 整除也不能被 5 整除的数字。
*/

import (
	"fmt"
	"sync"
)

type FizzBuzz struct {
	n          int      // 最大数
	numCh      chan int // 通知number()
	fizzCh     chan int // 通知fizz()
	buzzCh     chan int // 通知buzz()
	fizzBuzzCh chan int // 通知fizzbuzz()
	q          chan int // 通知所有goroutine停止
}

// n%3 == 0
func (f *FizzBuzz) fizz(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// 收到f.fizzCh中的数据，就可以打印fizz
		// 收到f.q的数据，就可以退出fizz
		// f.numCh <- 1通知number()开始执行
		select {
		case <-f.fizzCh:
			fmt.Println("fizz")
			f.numCh <- 1
		case <-f.q:
			return
		}
	}
}

// n%5 == 0
func (f *FizzBuzz) buzz(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-f.buzzCh:
			fmt.Println("buzz")
			f.numCh <- 1
		case <-f.q:
			return
		}
	}
}

// n%3 == 0 && n%5 == 0
func (f *FizzBuzz) fizzbuzz(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-f.fizzBuzzCh:
			fmt.Println("fizzbuzz")
			f.numCh <- 1
		case <-f.q:
			return
		}
	}
}

func (f *FizzBuzz) number(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= f.n; i++ {
		// 读取numCh可以开始运行。
		<-f.numCh

		// 根据不同i通知不同的goroutine执行。
		if i%3 == 0 && i%5 == 0 {
			f.fizzBuzzCh <- 1
			continue
		}
		if i%3 == 0 {
			f.fizzCh <- 1
			continue
		}
		if i%5 == 0 {
			f.buzzCh <- 1
			continue
		}
		fmt.Println(i)
		//i为特殊值(3、5有关时阻塞)，f.numCh <- 1保证number()继续执行
		f.numCh <- 1
	}

	// 执行完循环，退出其他3个goroutine
	f.q <- 1
	f.q <- 1
	f.q <- 1
}

func NewFizzBuzz(n int) *FizzBuzz {
	return &FizzBuzz{
		n:          n,
		numCh:      make(chan int, 1), //多一个缓冲，避免下轮循环死锁
		fizzCh:     make(chan int),
		buzzCh:     make(chan int),
		fizzBuzzCh: make(chan int),
		q:          make(chan int, 3), //用来通知另外三个线程退出，避免泄露
	}
}

func main() {
	fb := NewFizzBuzz(15)

	var wg sync.WaitGroup
	wg.Add(4)

	go fb.fizz(&wg)
	go fb.buzz(&wg)
	go fb.fizzbuzz(&wg)
	go fb.number(&wg)
	fb.numCh <- 1 // 避免number()阻塞

	wg.Wait()
}
