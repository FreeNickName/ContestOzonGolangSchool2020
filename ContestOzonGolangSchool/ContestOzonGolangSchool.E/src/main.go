package main

import (
	"sync"
	"time"
	"math/rand"
	"fmt"
	"reflect"
	E "syncMap"
	// E "syncCounter"
	// E "asyncFEachIn"
	// E "only1Goroutine"
	// E "asyncF"
	// E "sync1WG"
	// E "sync2WG"
	// E "sync4WG"
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
	overMax := 50
	testName := reflect.TypeOf(E.PkgName{}).PkgPath()
	defer elapsed(fmt.Sprintf("[%s] max: %d", testName, max))()

	in1 := make(chan int, max+overMax)
	in2 := make(chan int, max+overMax)
	ans := make(chan int, max+overMax)
	out := make(chan int, max+overMax)

	rand.Seed(time.Now().UnixNano())
	go func() {
		for i := 0; i < max + overMax; i++ {
			x1 := rand.Intn(1000)
			pull(in1, x1, 5)
			x2 := rand.Intn(24242)
			pull(in2, x2, 3)
			a := f_fast(x2) + f_fast(x1)
			pull(ans, a, 3)
			// println("a(", i, "):", a, "=", x1, "+", x2)
		}
	}()

	E.Merge2Channels(f, in1, in2, out, max)

	var wg sync.WaitGroup
	wg.Add(1)

	cnt := 0
	go func() {
		defer wg.Done()
		
		// println("M:wait")
		for {
			// println(cnt)
			if cnt == max {
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
					panic(fmt.Sprintf("step %d: %d <> %d", cnt, res, a))
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
				// time.Sleep(1)
			}
		}
	}()
	wg.Wait()
	if cnt != max {
		panic(fmt.Sprintf("digits overhead: %d", cnt))
	}
}