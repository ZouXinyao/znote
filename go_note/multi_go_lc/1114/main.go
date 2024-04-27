package main

import (
	"fmt"
	"sync"
)

/*
三个不同的线程 A、B、C 将会共用一个 Foo 实例。

线程 A 将会调用 first() 方法
线程 B 将会调用 second() 方法
线程 C 将会调用 third() 方法
请设计修改程序，以确保 second() 方法在 first() 方法之后被执行，third() 方法在 second() 方法之后被执行。
提示：

尽管输入中的数字似乎暗示了顺序，但是我们并不保证线程在操作系统中的调度顺序。
你看到的输入格式主要是为了确保测试的全面性。
示例 1：
输入：nums = [1,2,3] 输出："firstsecondthird" 解释： 有三个线程会被异步启动。输入 [1,2,3] 表示线程 A 将会调用 first() 方法，线程 B 将会调用 second() 方法，线程 C 将会调用 third() 方法。正确的输出是 "firstsecondthird"。
示例 2：
输入：nums = [1,3,2] 输出："firstsecondthird" 解释： 输入 [1,3,2] 表示线程 A 将会调用 first() 方法，线程 B 将会调用 third() 方法，线程 C 将会调用 second() 方法。正确的输出是 "firstsecondthird"。
*/

func first(ch1 chan int, wg *sync.WaitGroup) {
	// first不阻塞，ch1中传入数据
	defer wg.Done()
	fmt.Print("first")
	ch1 <- 1
}

func second(ch1, ch2 chan int, wg *sync.WaitGroup) {
	// second等待first(接收ch1中的数)，向ch2传数
	defer wg.Done()
	<-ch1
	fmt.Print("second")
	ch2 <- 1
}

func third(ch2 chan int, wg *sync.WaitGroup) {
	// third等待second(接收ch2中的数)
	defer wg.Done()
	<-ch2
	fmt.Print("thrid")
}

func main() {
	// 使用两个chan实现3个goroutine的顺序执行。
	var ch1 = make(chan int)
	var ch2 = make(chan int)
	var wg sync.WaitGroup
	wg.Add(3)
	go first(ch1, &wg)
	go second(ch1, ch2, &wg)
	go third(ch2, &wg)
	wg.Wait()
}
