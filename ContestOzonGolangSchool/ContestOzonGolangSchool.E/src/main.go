//https://github.com/FreeNickName/ContestOzonGolangSchool2020
package main

import (
	"sync"
	"time"
	"math/rand"
	"fmt"
	"reflect"
	E "map2out"
	// E "mapAndChan"
	// E "mapAndCache"
	// E "syncCounter"
	// E "asyncChannels"
	// E "asyncFAndSum"
	// E "asyncF"
	// E "sync1WG"
	// E "sync2WG"
	// E "sync4WG"
	// "math"
)

func elapsed(what string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s. elapsed %v\n", what, time.Since(start))
    }
}

func GetSlowF(amountIterations int) func(x int) int {
	return func(x int) int {
		// sum := 0
		for i := 0; i < amountIterations; i++ {
			// sum += i
		}
		// sum -= amountIterations
		return f_fast(x)
	}
}

func f_fast(x int) int {
	return x * x
}

// func f_fast(x int) int {
// 	return int(math.Pow(float64(x), 2))
// }

func ToChannel(in chan<- int, x int) {
	in <- x
}

func main() {
	max := 30
	complexity := 1000 * 1000 * 40

	extraSize := 70
	digitsPoolSize := max + extraSize
	testName := reflect.TypeOf(E.PkgName{}).PkgPath()
	defer elapsed(fmt.Sprintf("[%s] max: %d complexity: %d", testName, max, complexity))()

	in1 := make(chan int, digitsPoolSize)
	in2 := make(chan int, digitsPoolSize)
	ans := make(chan int, digitsPoolSize)
	out := make(chan int)
	f := GetSlowF(complexity)

	rand.Seed(time.Now().UnixNano())
	go func() {
		for i := 0; i < digitsPoolSize; i++ {
			x1 := rand.Intn(1000)
			ToChannel(in1, x1)
			x2 := rand.Intn(24242)
			ToChannel(in2, x2)
			a := f_fast(x2) + f_fast(x1)
			ToChannel(ans, a)
			// println("a(", i, "):", a, "=", x1, "+", x2)
			// if i > 100 {
			// 	close(in1)
			// 	break
			// }
		}
	}()

	E.Merge2Channels(f, in1, in2, out, max)

	var wg sync.WaitGroup
	wg.Add(1)

	cnt := 0
	go func() {
		defer wg.Done()
		for {
			if cnt == max {
				// close(in1)
				// close(in2)
				println("M: iterations is done")
				return
			}	
			select {
			case res, ok := <-out:
				if !ok {
					println("M: out is closed")
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
				// ToChannel(in1, x1)
				// x2 := rand.Intn(24242)
				// ToChannel(in2, x2)
				// a := f_fast(x2) + f_fast(x1)
				// ToChannel(ans, a)
				// time.Sleep(1)
			}
		}
	}()
	wg.Wait()
	if cnt != max {
		panic(fmt.Sprintf("amount iterations wrong: %d <> %d", max, cnt))
	}
	if extraSize != len(in1) {
		panic(fmt.Sprintf("amount digits in in1 is wrong: %d <> %d", extraSize, len(in1)))
	}
	if extraSize != len(in2) {
		panic(fmt.Sprintf("amount digits in in2 is wrong: %d <> %d", extraSize, len(in2)))
	}
	if 0 != len(out) {
		panic(fmt.Sprintf("amount digits in out is wrong: %d <> %d", 0, len(out)))
	}
	fmt.Printf("M: Report:\n Iterations: %d\n Last in1: %d\n Last in2: %d\n Last out: %d\n", cnt, len(in1), len(in2), len(out))
}
