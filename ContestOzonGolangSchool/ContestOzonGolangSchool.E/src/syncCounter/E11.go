package syncCounter

import (
	"sync"
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
