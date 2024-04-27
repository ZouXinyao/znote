package main

/*
两个不同的线程将会共用一个 FooBar 实例：

线程 A 将会调用 foo() 方法，而
线程 B 将会调用 bar() 方法
请设计修改程序，以确保 "foobar" 被输出 n 次。
示例 1：
输入：n = 1 输出："foobar" 解释：这里有两个线程被异步启动。其中一个调用 foo() 方法, 另一个调用 bar() 方法，"foobar" 将被输出一次。
示例 2：
输入：n = 2 输出："foobarfoobar" 解释："foobar" 将被输出两次。
*/

import (
	"fmt"
	"sync"
)

var fooMsg = make(chan int, 1)
var barMsg = make(chan int, 1)

// 使用两个channel进行goroutine通信，实现同步
// 使用waitGroup保证两个线程必定被执行完
// chan和wg可以像上个题写成局部变量
var wg = &sync.WaitGroup{}

func foo(n int) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		<-fooMsg
		fmt.Print("foo")
		barMsg <- 1
	}
}

func bar(n int) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		<-barMsg
		fmt.Print("bar")
		fooMsg <- 1
	}
}

func main() {
	n := 5
	fooMsg <- 1
	//触发执行
	wg.Add(2)
	go foo(n)
	go bar(n)
	wg.Wait()
}
