//https://github.com/FreeNickName/ContestOzonGolangSchool2020
package main

import (
	"sync"
	"time"
	"math/rand"
	"fmt"
	"reflect"
	E9 "map2out"
	E8 "mapAndChan"
	E10 "mapAndCache"
	E11 "syncCounter"
	E3 "asyncChannels"
	E4 "asyncFAndSum"
	E1 "asyncF"
	E2 "asyncF2"
	E5 "sync1WG"
	E6 "sync2WG"
	E7 "sync4WG"
	// "math"
)

type process struct {
	process func (f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int)
	name string
}

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
	max := 300
	complexity := 1000 * 1000 * 40
	extraSize := 70
	
	f := GetSlowF(complexity)
	variants := []process {
		process{process: E1.Merge2Channels, name: reflect.TypeOf(E1.PkgName{}).PkgPath()},
		process{process: E2.Merge2Channels, name: reflect.TypeOf(E2.PkgName{}).PkgPath()},
		process{process: E3.Merge2Channels, name: reflect.TypeOf(E3.PkgName{}).PkgPath()},
		process{process: E4.Merge2Channels, name: reflect.TypeOf(E4.PkgName{}).PkgPath()},
		process{process: E5.Merge2Channels, name: reflect.TypeOf(E5.PkgName{}).PkgPath()},
		process{process: E6.Merge2Channels, name: reflect.TypeOf(E6.PkgName{}).PkgPath()},
		process{process: E7.Merge2Channels, name: reflect.TypeOf(E7.PkgName{}).PkgPath()},
		process{process: E8.Merge2Channels, name: reflect.TypeOf(E8.PkgName{}).PkgPath()},
		process{process: E9.Merge2Channels, name: reflect.TypeOf(E9.PkgName{}).PkgPath()},
		process{process: E10.Merge2Channels, name: reflect.TypeOf(E10.PkgName{}).PkgPath()},
		process{process: E11.Merge2Channels, name: reflect.TypeOf(E11.PkgName{}).PkgPath()},
	}
	digitsPoolSize := max + extraSize
	ain1 := make([]int, digitsPoolSize)
	ain2 := make([]int, digitsPoolSize)
	aans := make([]int, digitsPoolSize)

	rand.Seed(time.Now().UnixNano())
	func() {
		for i := 0; i < digitsPoolSize; i++ {
			x1 := rand.Intn(1000)
			ain1 = append(ain1, x1)
			x2 := rand.Intn(24242)
			ain2 = append(ain2, x2)
			a := f_fast(x2) + f_fast(x1)
			aans = append(aans, a)
			// println("a(", i, "):", a, "=", x1, "+", x2)
			// if i > 100 {
			// 	close(in1)
			// 	break
			// }
		}
	}()

	var wg sync.WaitGroup
	for _, E := range variants {
		in1 := make(chan int, digitsPoolSize)
		in2 := make(chan int, digitsPoolSize)
		ans := make(chan int, digitsPoolSize)
		for i := 0; i < digitsPoolSize; i++ {
			ToChannel(in1, ain1[i])
			ToChannel(in2, ain2[i])
			ToChannel(ans, aans[i])
		}
		out := make(chan int)
		testName := E.name
		cnt := 0
		wg.Add(1)
		func () {
			defer elapsed(fmt.Sprintf("[%s] max: %d complexity: %d", testName, max, complexity))()
			E.process(f, in1, in2, out, max)
			go func() {
				defer wg.Done()
				for {
					if cnt == max {
						// close(in1)
						// close(in2)
						// println("M: iterations is done")
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
			// fmt.Printf("M: Report:\n Iterations: %d\n Last in1: %d\n Last in2: %d\n Last out: %d\n", cnt, len(in1), len(in2), len(out))
		}()
	}
}
