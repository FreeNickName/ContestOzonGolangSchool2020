package main

import (
	"sync"
	"time"
	"math/rand"
	"fmt"
	// "math"
)

type syncMap struct{
	sync.RWMutex
	m map[int]int
}

func CreateMap() *syncMap {
	return &syncMap{m: make(map[int]int)}
}

func (b *syncMap) Load(key int) (int, bool) {
    b.Lock()
    defer b.Unlock()
    val, ok := b.m[key]
    return val, ok
}

func (b *syncMap) Store(key int, value int) {
    b.Lock()
    defer b.Unlock()
    b.m[key] = value
}

func (b *syncMap) Delete(key int) {
    b.Lock()
    defer b.Unlock()
    delete(b.m, key)
}

func Merge2Channels(f func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	f1 := make(chan int, 10)
	f2 := make(chan int, 10)
	b1 := CreateMap()
	b2 := CreateMap()
	// s1 := make(chan bool, 10)
	// s2 := make(chan bool, 10)
	c1 := CreateCounter(n)
	c2 := CreateCounter(n)

	go ChannelToBuff(f, in1, n, b1, f1, c1)
	go ChannelToBuff(f, in2, n, b2, f2, c2)

	// go MapToChanByOrder(b1, f1, n, s1)
	// go MapToChanByOrder(b2, f2, n, s2)

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

func ChannelToBuff(f func(int) int, in <-chan int, max int, b *syncMap, out chan<- int, c *counter) {
	i := 0
	for  {
		if i == max {
			break
		}
		res, ok := <-in
		if !ok {
			println("in closed")
			break
		}
		go UseFToMap(f, i, res, b, out, c)
		i++
	}
}

func UseFToMap(f func(int) int, key int, val int, b *syncMap, out chan<- int, c *counter) {
	b.Store(key, f(val))
	// println("f:", key, "=", val)
	MapToChanByOrder(b, out, c)
}

type counter struct {
	sync.Mutex
	i int
	max int
}

func (c *counter) Get() int {
    return c.i
}

func (c *counter) Inc() {
    c.i++
}

func (c *counter) IsDone() bool {
    return c.i == c.max
}

// func Counter() func(x int) int {
// 	i := 0
// 	return func() (ret int) {
// 		for _, num := range params {
// 			sum += num
// 		}
// 		ret = sum
// 		return
// 	}
// }

func CreateCounter(max int) *counter {
	return &counter{i: 0, max: max}
}

// Запись мапа в канал по порядку ключей в asc, signal указывает наполнение мапа 
func MapToChanByOrder(b *syncMap, out chan<- int, c *counter) {
	defer c.Unlock()
	c.Lock()
	if c.IsDone() {
		// close(out)
		return
	}
	for {
		i := c.Get()
		v, ok := b.Load(i)
		// println("b", i, "=", v)
		if !ok {
			return
		}	
		out <- v
		b.Delete(i)
		c.Inc()
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
	return x*x
}

func pull(in chan<- int, x int, s int) {
	in <- x
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
