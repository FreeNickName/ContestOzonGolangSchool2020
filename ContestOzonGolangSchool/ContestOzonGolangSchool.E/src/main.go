package main

import (
	"sync"
	"time"
	"math/rand"
	"fmt"
	// E "syncMap"
	// E "syncCounter"
	// E "asyncFEachIn"
	// E "asyncF"
	E "syncWG"
	// "math"
)

func elapsed(what string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s took %v\n", what, time.Since(start))
    }
}

func f(x int) int {
	sum := 0
	max := 100000000
	for i := 0; i < max; i++ {
		sum += i
	}
	sum -= max
	return f_fast(x)
}

func f_fast(x int) int {
	return x * x
}

func pull(in chan<- int, x int, s int) {
	in <- x
}

func main() {
	max := 1000

	defer elapsed(fmt.Sprintf("[syncMap] max: %d", max))()

	// var t sync.Map
	// t.Store(1,1)
	// e,ok := t.Load(0)
	// println(ok, e)

	in1 := make(chan int, max)
	in2 := make(chan int, max)
	ans := make(chan int, max)
	out := make(chan int, max)

	rand.Seed(time.Now().UnixNano())
	go func() {
		for i := 0; i < max; i++ {
			x1 := rand.Intn(1000)
			pull(in1, x1, 5)
			x2 := rand.Intn(24242)
			pull(in2, x2, 3)
			a := f_fast(x2) + f_fast(x1)
			pull(ans, a, 3)
			// println("a(", i, "):", a, "=", x1, "+", x2)
		}
	}()

	E.Merge2Channels(f, in1, in2, out, max-3)

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
				// println("M(", cnt,"):", res, "=", a)
				cnt++	
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
			}
		}
	}()
	wg.Wait()
}