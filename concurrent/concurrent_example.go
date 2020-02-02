package concurrent

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

/*
	多线程交替输出：
	新开两个子线程，分别输出1,3,5,7,9...和2,4,6,8,10...，主线程接受子线程的值，输出1,2,3,4,5...
	主要问题点是如何协调两个生产者的步调，因为协程的调度是不可控的，所以需要额外的机制来协调
*/

/*
  1.使用多个无缓冲channel实现，每个生产者由channel触发生产
*/
func AlternateOutputViaChannel() {
	out := make(chan int)
	c1 := make(chan struct{})
	c2 := make(chan struct{})
	runtime.GC()

	// 生产者1
	go func() {
		i := 1
		for {
			if _, ok := <-c1; !ok {
				return
			}
			fmt.Println("1, out=", i)
			out <- i
			i += 2
		}
	}()

	// 生产者2
	go func() {
		i := 2
		for {
			if _, ok := <-c2; !ok {
				return
			}
			fmt.Println("2, out=", i)
			out <- i
			i += 2
		}
	}()

	// fmt.Println("stage 1, go num=", runtime.NumGoroutine())
	c1 <- struct{}{} // 启动g1，这里因为c1是堵塞队列，所以在g1未准备好的情况下也会堵塞主线程，保证步调一致
	for {
		i, _ := <-out
		if i >= 100 {
			close(c1)
			close(c2)
			fmt.Println("close")
			break
		}
		if i%2 == 1 {
			c2 <- struct{}{} // 启动c2
		} else {
			c1 <- struct{}{} // 	启动c1
		}
		// fmt.Println("main, out=", i)
	}
	time.Sleep(2 * time.Second)
	fmt.Println("stage 2, go num=", runtime.NumGoroutine())
}

/*
	使用原子自增变量协调，就是将上面的i%2逻辑挪到生产者
*/
func AlternateOutputViaAtomic() {
	sig := int32(1)
	out := make(chan int32) // 简单起见，用channel来收集输出

	// 生产者1
	go func() {
		i := int32(1)
		for i <= 100 {
			// 这里不能直接使用cas，因为cas之后另外一个线程的条件也会马上满足，从而导致out输出结果的不确定性
			// 所以应该先判断值，然后输出out，最后调整原子变量
			if atomic.LoadInt32(&sig)%2 != 1 {
				time.Sleep(1 * time.Millisecond)
				continue
			}
			fmt.Println("1, out=", i)
			out <- i
			atomic.StoreInt32(&sig, i+1)
			i += 2
		}
	}()

	// 生产者2
	go func() {
		i := int32(2)
		for i <= 100 {
			if atomic.LoadInt32(&sig)%2 != 0 {
				time.Sleep(1 * time.Millisecond)
				continue
			}
			fmt.Println("2, out=", i)
			out <- i
			atomic.StoreInt32(&sig, i+1)
			i += 2
		}
	}()

	m := int32(100)
	fmt.Println("stage 1, go num=", runtime.NumGoroutine())
	for {
		i, _ := <-out
		// fmt.Println("main, out=", i)
		if i >= m {
			fmt.Println("finish")
			break
		}
	}
	fmt.Println("stage 2, go num=", runtime.NumGoroutine())
}

/*
	使用条件信号处理
*/
func AlternateOutputViaCond() {
	var m sync.Mutex
	c := sync.NewCond(&m)
	quit := make(chan bool)

	// 生产者1
	go func() {
		i := int32(1)
		for i <= 100 {
			c.L.Lock()
			// fmt.Println("1, waiting")
			c.Wait()
			fmt.Println("1, out=", i)
			// out <- i
			i += 2
			c.L.Unlock()
			c.Broadcast()
		}
	}()

	// 生产者2
	go func() {
		i := int32(2)
		for i <= 100 {
			c.L.Lock()
			// fmt.Println("2, waiting")
			c.Wait()
			fmt.Println("2, out=", i)
			// out <- i
			i += 2
			c.L.Unlock()
			c.Broadcast()
		}
		quit <- true
	}()

	// 正常状态下，请求写锁是公平的，不会区分先来后到，但是当协程等待锁的时间超过一定时间之后会改用fifo调度
	time.Sleep(100 * time.Millisecond)
	c.Signal()
	<-quit
}

// Singleton 单例
type Singleton struct {
	Value int64
}

var s *Singleton

var _ = new(sync.Once)
var once = new(Once)

// GetGlobalSingle 使用once来实现全局单例
func GetGlobalSingle() *Singleton {
	once.Do(func() {
		s = &Singleton{
			Value: time.Now().Unix(),
		}
	})
	return s
}
