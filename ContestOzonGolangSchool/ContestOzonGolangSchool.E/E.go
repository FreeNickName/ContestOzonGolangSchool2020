package main

import (
	"sync"
	"time"
	"math/rand"
	"fmt"
)

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	// println()
	go Sum2ChannelsN(f, in1, in2, out, n)
}

func SumChannels(f func(int) int, in1 int, in2 int, out chan<- int, key int, wg *sync.WaitGroup, pwg *sync.WaitGroup) {
	res := f(in1) + f(in2)
    if key > 0 {
		pwg.Wait()
	}
	out <- res
	defer wg.Done()
	// println(key, ":", res)
}

func Sum2ChannelsN(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	prevChanWg := new(sync.WaitGroup)
	prevCalcWg := new(sync.WaitGroup)
	i := 0
	for  {
		currChanWg := new(sync.WaitGroup)
		currChanWg.Add(1)
		currCalcWg := new(sync.WaitGroup)
		currCalcWg.Add(1)

		go func(key int, wg *sync.WaitGroup, pwg *sync.WaitGroup, cwg *sync.WaitGroup, pcwg *sync.WaitGroup) {
			if key > 0 {
				pwg.Wait()
			}
			x1 := <-in1
			x2 := <-in2
			wg.Done()
			go SumChannels(f, x1, x2, out, key, cwg, pcwg)
		}(i, currChanWg, prevChanWg, currCalcWg, prevCalcWg)
		prevCalcWg = currCalcWg
		prevChanWg = currChanWg
		i++
	}
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

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < max; i++ {
		x1 := rand.Intn(1000)
		pull(in1, x1, 5)
		x2 := rand.Intn(24242)
		pull(in2, x2, 3)
		a := f_fast(x2) + f_fast(x1)
		pull(ans, a, 3)
		// println("a(", i, "):", a, "=", x1, "+", x2)
	}

	Merge2Channels(f, in1, in2, out, max)

	cnt := 0
	// println("M:wait")
	for {
		
		// select {
		// case 
		res, ok := <-out
			if !ok {
				return
			}
			cnt++
			a := <-ans
			if (a != res) {
				panic(fmt.Sprintf("%d: %d <> %d", cnt, res, a))
			}
			println("M(", cnt,"):", res, "=", a)
			if cnt >= max {
				return
			}	
		// default:
			// println("M:...")
			// rand.Seed(time.Now().UnixNano())
			// x1 := rand.Intn(1000)
			// pull(in1, x1, 5)
			// x2 := rand.Intn(24242)
			// pull(in2, x2, 3)
			// a := f_fast(x2) + f_fast(x1)
			// pull(ans, a, 3)
			// time.Sleep(2)
		// }
	}
}
