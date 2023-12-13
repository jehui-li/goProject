// 生产者产生随机数，消费者求数字各个位的和
// 生产者---itemChan--->消费者---resultChan--->输出，遍历，打印结果
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
sync.WaitGroup 用于等待一组 goroutine 完成执行。
通过 Add 方法增加计数器，Done 方法减少计数器，
Wait 方法用于等待计数器归零。
*/
var sw sync.WaitGroup

// 随机数
type obj struct {
	id  int64
	num int64
}

// 随机数的和
type objSum struct {
	obj    *obj
	sumObj int64
}

// 传送随机数，*obj表示传送指向 obj 类型的指针
var itemChan chan *obj

// 传送随机数，*item表示传送指向 item 类型的指针
var resultChan chan *objSum

// 生产者
func producer(itemChan chan *obj) {
	var id int64
	for i := 0; i < 10000; i++ {
		id++
		number := rand.Int63()
		tmp := &obj{
			id:  id,
			num: number,
		}
		itemChan <- tmp
	}
	close(itemChan)
}

// 计算求和
func calc(num int64) int64 {
	var sum int64
	for num > 0 {
		sum = sum + num%10
		num = num / 10
	}
	return sum
}

// 消费者
func consumer(itemChan chan *obj, resultChan chan *objSum) {
	//计算器减1
	defer sw.Done()
	for tmp := range itemChan {
		sum := calc(tmp.num)
		objSum := &objSum{
			obj:    tmp,
			sumObj: sum,
		}
		resultChan <- objSum
	}
}

// 打印
func printResult(resultChan chan *objSum) {
	//迭代通道中的值，直到通道关闭
	for ret := range resultChan {
		//%v 会根据值的具体类型自动选择合适的格式
		fmt.Printf("id: %v, num: %v, sum: %v \n", ret.obj.id, ret.obj.num, ret.sumObj)
		//将 goroutine 的执行挂起，展示效果
		time.Sleep(time.Second)
	}
}

func startWorker(n int, itemChan chan *obj, resultChan chan *objSum) {
	for i := 0; i < n; i++ {
		go consumer(itemChan, resultChan)
	}
}

func main() {
	//该通道的缓冲区大小为100。缓冲区满时，生产者会阻塞
	itemChan = make(chan *obj, 10000)
	resultChan = make(chan *objSum, 10000)

	//异步调用 producer 函数，使用 go 关键字启动一个新的 goroutine，不会等待 producer 函数的完成，而是会继续执行后续的代码
	go producer(itemChan)
	//启用计数器
	sw.Add(30)
	//消费者goroutine
	startWorker(30, itemChan, resultChan)
	//等待goroutine结束
	sw.Wait()
	//关闭resultChan通道
	close(resultChan)
	//调用输出函数
	printResult(resultChan)
}
