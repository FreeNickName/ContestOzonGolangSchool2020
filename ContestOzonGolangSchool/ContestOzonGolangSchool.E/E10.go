package main

import (
	"sync"
	"time"
	"math/rand"
	"fmt"
	// "math"
)
// import "math"
func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	f1 := make(chan int, 10)
	f2 := make(chan int, 10)
	// var b1 sync.Map
	b1 := CreateMap()
	// var b2 sync.Map
	b2 := CreateMap()
	s1 := make(chan bool, 10)
	s2 := make(chan bool, 10)

	go UseFChannel(f, in1, f1, n, b1, s1)
	go UseFChannel(f, in2, f2, n, b2, s2)
	go BuffToChannel(f1, n, b1, s1)
	go BuffToChannel(f2, n, b2, s2)

	go SumChannels(f1, f2, out)
}

// func fast(x int) int {
// 	return int(math.Pow(float64(x), 2))
// }

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

func UseFChannel(f func(int) int, in1 <-chan int, out1 chan<- int, n int, b *buf, s chan<- bool) {
// func UseFChannel(f func(int) int, in1 <-chan int, out1 chan<- int, n int, b *sync.Map, s chan<- bool) {
	i := 0
	for  {
		if i == n {
			break
		}
		i++
		res, ok := <-in1
		if !ok {
			println("in1 closed")
			break
		}
		go RunF(f, i-1, res, b, s)
		// out1 <- RunF(f, i, res, b)
	}
	// close(out1)
}

type buf struct{
	sync.RWMutex
	m map[int]int
}

func CreateMap() *buf {
	return &buf{m: make(map[int]int)}
}

func (b *buf) Load(key int) (int, bool) {
    b.Lock()
    defer b.Unlock()
    val, ok := b.m[key]
    return val, ok
}

func (b *buf) Store(key int, value int) {
    b.Lock()
    defer b.Unlock()
    b.m[key] = value
}

func (b *buf) Delete(key int) {
    b.Lock()
    defer b.Unlock()
    delete(b.m, key)
}

// func BuffToChannel(out1 chan<- int, n int, b *sync.Map, signal <-chan bool) {
func BuffToChannel(out1 chan<- int, n int, b *buf, signal <-chan bool) {
	i := 0
	for  {
		if i == n {
			break
		}
		select {
			case _, oks := <-signal:
				if (!oks) {
					// return
				}
				// b.RLock()
				for {
					v, ok := b.Load(i)
					// b.RUnlock()
					// println("b:", i, "=", v)
					if !ok {
						break
					}	
					out1 <- v //.(int)
					// b.Lock()
					// delete(b.m, i)
					b.Delete(i)
					// b.Unlock()
					i++
					
				}
		}
	}
	close(out1)
}

// func RunF(f func(int) int, n int, x int, b *sync.Map, s chan<- bool) {
func RunF(f func(int) int, n int, x int, b *buf, s chan<- bool) {
	b.Store(n, f(x))
	// println("f:", n, "=", x)
	s <- true
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

func pull(in1 chan<- int, x int, s int) {
	in1 <- x
}

func main() {
	max := 1000

	// var t sync.Map
	// t.Store(1,1)
	// e,ok := t.Load(0)
	// println(ok, e)

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
				// println("M(", cnt,"):", res, "=", a)
				cnt++	
			default:
				// println("M:...")
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
