package syncMap

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
	f1 := make(chan int, 5)
	f2 := make(chan int, 5)
	b1 := CreateMap()
	b2 := CreateMap()
	s1 := make(chan bool, 2)
	s2 := make(chan bool, 2)

	go ChannelToBuff(f, in1, n, b1, s1)
	go ChannelToBuff(f, in2, n, b2, s2)

	go MapToChanByOrder(b1, f1, n, s1)
	go MapToChanByOrder(b2, f2, n, s2)

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

func ChannelToBuff(f func(int) int, in <-chan int, max int, b *syncMap, done chan<- bool) {
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
		go UseFToMap(f, i, res, b, done)
		i++
	}
}

func UseFToMap(f func(int) int, key int, val int, b *syncMap, done chan<- bool) {
	b.Store(key, f(val))
	// println("f:", key, "=", val)
	done <- true
}

// Запись мапа в канал по порядку ключей в asc, signal указывает наполнение мапа 
func MapToChanByOrder(b *syncMap, out chan<- int, max int, signal <-chan bool) {
	i := 0
	for  {
		if i == max {
			break
		}
		select {
			case _, oks := <-signal:
				if (!oks) {
					println("signal closed")
					return
				}
				for {
					v, ok := b.Load(i)
					// println("b", i, "=", v)
					if !ok {
						break
					}	
					out <- v
					b.Delete(i)
					i++
				}
		}
	}
	close(out)
}


