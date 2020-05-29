package main

import (
	"sync"
	"time"
	"math/rand"
	"fmt"
)

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	f1 := make(chan int, 10)
	f2 := make(chan int, 10)

	go UseFChannels(f, in1, in2, f1, f2, n)

	go SumChannels(f1, f2, out)
}

func SumChannels(in1 <-chan int, in2 <-chan int, out chan<- int) {
	for  {
		sum := 0
		ok := false
		select {
			case sum, ok = <-in1:
				sum += <-in2
			case sum, ok = <-in2:
				sum += <-in1
		}
		if !ok {
			// println("push to out is done")
			// close(out)
			return
		}
		out <- sum
	}
}

func UseFChannels(f func(int) int, in1 <-chan int, in2 <-chan int, out1 chan<- int, out2 chan<- int, n int) {
	i := 0
	for  {
		if i == n {
			break
		}
		i++
		select {
			case res, ok := <-in1:
				if !ok {
					println("in1 closed")
					break
				} 
				out1 <- f(res)
				out2 <- f(<-in2)
			case res, ok := <-in2:
				if !ok {
					println("in2 closed")
					break
				} 
				out2 <- f(res)
				out1 <- f(<-in1)
		}
	}
	close(out1)
	close(out2)
}

func f(x int) int {
	sum := 0
	max := 10000000
	for i := 0; i < max; i++ {
		sum += i
	}
	sum -= max
	return f_fast(x)
}

func f_fast(x int) int {
	return x
}

func pull(in1 chan<- int, x int, s int) {
	in1 <- x
}

func main() {
	max := 10

	in1 := make(chan int, max)
	in2 := make(chan int, max)
	ans := make(chan int, max)
	out := make(chan int, max)

	// rand.Seed(time.Now().UnixNano())
	// for i := 0; i < max; i++ {
	// 	x1 := rand.Intn(1000)
	// 	pull(in1, x1, 5)
	// 	x2 := rand.Intn(24242)
	// 	pull(in2, x2, 3)
	// 	a := f_fast(x2) + f_fast(x1)
	// 	pull(ans, a, 3)
	// 	// println("a(", i, "):", a, "=", x1, "+", x2)
	// }

	Merge2Channels(f, in1, in2, out, max-3)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		cnt := 0
		// println("M:wait")
		for {
			// println(cnt)
			if cnt == max - 3 {
				close(in1)
				close(in2)
				println("main closed")
				return
			}	
			
			select {
			case res, ok := <-out:
				if !ok {
					println("main exit")
					return
				}
				
				a := <-ans
				if (a != res) {
					panic(fmt.Sprintf("%d: %d <> %d", cnt, res, a))
				}
				println("M(", cnt,"):", res, "=", a)
				cnt++	
			default:
				println("M:...")
				rand.Seed(time.Now().UnixNano())
				x1 := rand.Intn(1000)
				pull(in1, x1, 5)
				x2 := rand.Intn(24242)
				pull(in2, x2, 3)
				a := f_fast(x2) + f_fast(x1)
				pull(ans, a, 3)
				time.Sleep(2)
			}
		}
	}()
	wg.Wait()
}
